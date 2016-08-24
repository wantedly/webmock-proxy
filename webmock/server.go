package webmock

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/url"
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
		log.Println("[INFO] Use local cache files.")
	} else {
		db, err = NewDBConnection()
		if err != nil {
			return nil, fmt.Errorf("[ERROR] Faild to connect database: %v", err)
		}
		log.Println("[INFO] Use database.")
	}
	proxy := goproxy.NewProxyHttpServer()
	if config.masterURL != "" {
		proxyURL, err := url.Parse(config.masterURL)
		if err != nil {
			return nil, fmt.Errorf("[ERROR] Faild to parse webmock-proxy master url: %v", err)
		}
		proxy.Tr.Proxy = http.ProxyURL(proxyURL)
		proxy.Tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		proxy.Tr.DisableCompression = true
	}

	return &Server{
		config: config,
		db:     db,
		proxy:  proxy,
		body:   "",
		head:   make(map[string][]string),
	}, nil
}

func (s *Server) connectionCacheHandler() {
	s.proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)
	s.proxy.OnRequest().DoFunc(
		func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			log.Printf("[INFO] req %s %s", ctx.Req.Method, ctx.Req.URL.Host+ctx.Req.URL.Path)
			req.Header.Del("Proxy-Connection")

			// DeepCopy *http.Request.Header (type: map[string][]string)
			reqHeader := make(map[string][]string, len(req.Header))
			for k, v := range req.Header {
				reqHeader[k] = v
			}

			reqBody, err := ioReader(req.Body)
			if err != nil {
				log.Printf("[ERROR] %v", err)
			}
			s.body = string(reqBody)
			s.head = reqHeader
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
			req.Header.Del("Proxy-Connection")

			reqBody, err := ioReader(req.Body)
			if err != nil {
				log.Printf("[ERROR] %v", err)
			}
			conn, err := NewConnection(req, s)
			if err != nil {
				log.Printf("[ERROR] %v", err)
			}
			if conn != nil {
				is, err := validateRequest(req, conn, string(reqBody))
				if err != nil {
					log.Printf("[ERROR] %v", err)
				}
				if is == true {
					log.Printf("[INFO] Create HTTP/S response using connection cache.")
					resp, err := createHttpResponse(req, conn)
					if err != nil {
						log.Printf("[ERROR] %v", err)
					}
					return req, resp
				}
			}
			log.Printf("[INFO] Not match http connection cache.")
			resp, err := createHttpErrorResponse(req)
			if err != nil {
				log.Printf("[ERROR] %v", err)
			}
			return req, resp
		})
}

func (s *Server) NonProxyHandler(config *Config) {
	s.proxy.NonproxyHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		pattern := req.URL.Path
		switch pattern {
		case "/":
			if req.Method != "POST" {
				createHttpResponseWriter(w, "Not Found", 404)
				return
			}

			reqBody, err := ioReader(req.Body)
			if err != nil {
				fmt.Fprintln(w, err)
				return
			}
			var jsonReqBody interface{}
			err = jsonToStruct([]byte(reqBody), &jsonReqBody)
			if err != nil {
				createHttpResponseWriter(w, "Bad request", 404)
				return
			}

			callbackURL := jsonReqBody.(map[string]interface{})["callback"].(string)
			endpointURL := jsonReqBody.(map[string]interface{})["endpoint"].(string)
			method := jsonReqBody.(map[string]interface{})["method"].(string)
			if (callbackURL == "") && (endpointURL == "") && (method == "") {
				createHttpResponseWriter(w, "Bad request", 404)
				return
			}

			//e.g. https://api.example.com/users => [https api.example.com/users]
			x := strings.Split(endpointURL, "://")

			// e.g. api.example.com/users => [api.example.com users]
			x2 := strings.SplitN(x[1], "/", 2)
			callbackURL = callbackURL + "/" + x2[1]
			scheme := x[0]
			var url string
			switch scheme {
			case "http":
				url = filepath.Join(x2[0], x2[1])
			case "https":
				url = filepath.Join(x2[0]+":443", x2[1])
			}

			var conn Connection
			if s.config.local == true {
				dst := filepath.Join(config.cacheDir, url, method+".json")
				b, err := readFile(dst)
				if err != nil {
					createHttpResponseWriter(w, "Don't exist http connection cache", 404)
					return
				}
				err = jsonToStruct(b, &conn)
				if err != nil {
					createHttpResponseWriter(w, "Faild to create callback request body", 500)
					return
				}
			} else {
				endpoint := findEndpoint(req.Method, url, s.db)
				if len(endpoint.Connections) == 0 {
					createHttpResponseWriter(w, "Don't exist http connection cache", 404)
					return
				}
			}
			var header interface{}
			b := []byte(conn.Request.Header)
			err = jsonToStruct(b, &header)
			if err != nil {
				createHttpResponseWriter(w, "Faild to create callback request body", 500)
				return
			}

			callbackReq, err := http.NewRequest(
				method,
				callbackURL,
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
				createHttpResponseWriter(w, "Faild to callback", 500)
				return
			}
			body, err := ioReader(resp.Body)
			if err != nil {
				createHttpResponseWriter(w, "Faild to read callback response body", 500)
				return
			}
			respStruct, err := responseStruct(body, resp)
			if err != nil {
				createHttpResponseWriter(w, "Faild to read callback response", 500)
				return
			}
			var mes string
			if respStruct.String == conn.Response.String {
				mes = "Valid"
			} else {
				mes = "Invalid"
			}
			createHttpResponseWriter(w, mes, 200)
			return

		default:
			createHttpResponseWriter(w, "Not Found", 404)
			return
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
	log.Println("[INFO] Serving webmock-proxy on " + s.config.port)
	log.Fatal("[ERROR] ", http.ListenAndServe(s.config.port, s.proxy))
}
