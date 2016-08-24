package webmock

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/elazarl/goproxy"
	"github.com/jinzhu/gorm"
)

type Server struct {
	config *Config
	db     *gorm.DB
	proxy  *goproxy.ProxyHttpServer
	body   string
	head   map[string][]string
}

func NewServer(config *Config) (*Server, error) {
	var db *gorm.DB
	var err error
	if config.local == true {
		if config.syncCache == true {
			err := sync(config)
			if err != nil {
				return nil, fmt.Errorf("[ERROR] Faild to sync cache: %v", err)
			}
		}
		log.Println("[INFO] Use local cache files.")
	} else {
		db, err = NewDBConnection()
		if err != nil {
			return nil, fmt.Errorf("[ERROR] Faild to connect database: %v", err)
		}
		log.Println("[INFO] Use database.")
	}

	return &Server{
		config: config,
		db:     db,
		proxy:  goproxy.NewProxyHttpServer(),
		body:   "",
		head:   make(map[string][]string),
	}, nil
}

func (s *Server) connectionCacheHandler() {
	s.proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)
	s.proxy.OnRequest().DoFunc(
		func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			log.Printf("[INFO] req %s %s", ctx.Req.Method, ctx.Req.URL.Host+ctx.Req.URL.Path)

			// DeepCopy *http.Request.Header (type: map[string][]string)
			reqHeader := make(map[string][]string, len(req.Header))
			for k, v := range req.Header {
				reqHeader[k] = v
			}

			reqBody, err := ioReader(req.Body)
			if err != nil {
				log.Printf("[ERROR] %v", err)
			}
			s.body = reqBody
			s.head = reqHeader
			req.Body = ioutil.NopCloser(bytes.NewBufferString(reqBody))
			return req, nil
		})
	s.proxy.OnResponse().Do(
		goproxy.HandleBytes(
			func(b []byte, ctx *goproxy.ProxyCtx) []byte {
				log.Printf("[INFO] resp %s", ctx.Resp.Status)
				reqBody := s.body
				reqHeader := s.head
				ctx.Req.Header = reqHeader
				err := createCache(reqBody, b, ctx.Req, ctx.Resp, s)
				if err != nil {
					log.Printf("[ERROR] %v", err)
				}
				fmt.Println(string(b))
				return b
			}))
}

func (s *Server) mockOnlyHandler() {
	s.proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)
	s.proxy.OnRequest().DoFunc(
		func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			log.Printf("[INFO] req %s %s", ctx.Req.Method, ctx.Req.URL.Host+ctx.Req.URL.Path)
			reqBody, err := ioReader(req.Body)
			if err != nil {
				log.Printf("[ERROR] %v", err)
			}
			conn, err := NewConnection(req, s)
			if err != nil {
				log.Printf("[ERROR] %v", err)
			}
			var resp *http.Response
			if conn != nil {
				is, err := validateRequest(req, conn, reqBody)
				if err != nil {
					log.Printf("[ERROR] %v", err)
				}
				if is == true {
					log.Printf("[INFO] Create HTTP/S response using connection cache.")
					resp, err = createHttpResponse(req, conn)
					if err != nil {
						log.Printf("[ERROR] %v", err)
					}
					req.Body = ioutil.NopCloser(bytes.NewBufferString(reqBody))
					return req, resp
				}
			}
			log.Printf("[INFO] Not match http connection cache.")
			resp, err = createHttpErrorResponse(req)
			if err != nil {
				log.Printf("[ERROR] %v", err)
			}
			req.Body = ioutil.NopCloser(bytes.NewBufferString(reqBody))
			return req, resp
		})
}

// {callback:"http://localhost:3000/api/users" endpoint:"https://example.com/api/users", method:"GET"}
func (s *Server) NonProxyHandler(config *Config) {
	s.proxy.NonproxyHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		pattern := req.URL.Path
		switch pattern {
		case "/":
			if req.Method != "POST" {
				mes := "Not Found"
				respBody := &responseBody{Message: mes}
				byteArr, err := structToJSON(respBody)
				if err != nil {
					fmt.Fprintln(w, err)
				}
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(404)
				fmt.Fprintf(w, string(byteArr))
				return
			}

			reqBody, err := ioReader(req.Body)
			if err != nil {
				fmt.Fprintln(w, err)
				return
			}
			var rbody interface{}
			err = jsonToStruct([]byte(reqBody), &rbody)
			if err != nil {
				fmt.Fprintln(w, err)
				return
			}

			callbackURL := rbody.(map[string]interface{})["callback"].(string)
			endpoint := rbody.(map[string]interface{})["endpoint"].(string)
			method := rbody.(map[string]interface{})["method"].(string)
			if (callbackURL == "") && (endpoint == "") && (method == "") {
				mes := "Not Found"
				respBody := &responseBody{Message: mes}
				byteArr, err := structToJSON(respBody)
				if err != nil {
					fmt.Fprintln(w, err)
				}
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(404)
				fmt.Fprintf(w, string(byteArr))
			}

			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(200)

			//e.g. https://api.example.com/users => [https api.example.com/users]
			x := strings.Split(endpoint, "://")

			// e.g. api.example.com/users => [api.example.com users]
			x2 := strings.SplitN(x[1], "/", 2)
			scheme := x[0]
			var url string
			switch scheme {
			case "http":
				url = filepath.Join(x2[0], x2[1])
			case "https":
				url = filepath.Join(x2[0]+":443", x2[1])
			}
			dst := filepath.Join(config.cacheDir, url, method+".json")
			b, err := readFile(dst)
			if err != nil {
				return
			}
			var conn Connection
			err = jsonToStruct(b, &conn)
			if err != nil {
				return
			}
			var header interface{}
			b = []byte(conn.Request.Header)
			err = jsonToStruct(b, &header)
			if err != nil {
				return
			}

			callbackReq, err := http.NewRequest(
				method,
				endpoint,
				bytes.NewBuffer([]byte(conn.Request.String)),
			)
			for k, v := range header.(map[string]interface{}) {
				for _, vv := range v.([]interface{}) {
					callbackReq.Header.Set(k, vv.(string))
				}
			}
			client := new(http.Client)
			resp, err := client.Do(callbackReq)
			if err != nil {
				return
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return
			}
			respStruct, err := responseStruct(body, resp)
			if err != nil {
				return
			}
			var mes string
			if respStruct.String == conn.Response.String {
				mes = "Valid"
			} else {
				mes = "Invalid"
			}

			nnbody := &responseBody{Message: mes}
			byteArr, err := structToJSON(nnbody)
			if err != nil {
				fmt.Fprintln(w, err)
			}
			fmt.Fprintf(w, string(byteArr))
		default:
			http.Error(w, "Not Found", 404)
		}
	})
}

func (s *Server) Start() {
	if s.config.record == true {
		log.Println("[INFO] All HTTP/S request and response is cached.")
		s.connectionCacheHandler()
	} else {
		s.mockOnlyHandler()
	}
	s.NonProxyHandler(s.config)
	log.Println("[INFO] Running...")
	http.ListenAndServe(":8080", s.proxy)
}
