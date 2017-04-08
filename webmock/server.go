package webmock

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

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
	db, err := initDB(config)
	if err != nil {
		return nil, err
	}

	return &Server{
		config: config,
		db:     db,
		proxy:  goproxy.NewProxyHttpServer(),
		body:   "",
		head:   make(map[string][]string),
	}, nil
}

func initDB(config *Config) (*gorm.DB, error) {
	if config.local == false {
		db, err := NewDBConnection()
		if err != nil {
			return nil, fmt.Errorf("[ERROR] Faild to connect database: %v", err)
		}
		log.Println("[INFO] Use database.")
		return db, nil
	}
	log.Println("[INFO] Use local cache files.")
	return nil, nil
}

func (s *Server) connectionCacheHandler() {
	s.proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)
	s.proxy.OnRequest().DoFunc(
		func(req *http.Request, pctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			log.Printf("[INFO] req %s %s", pctx.Req.Method, pctx.Req.URL.Host+pctx.Req.URL.Path)

			body, err := ioutil.ReadAll(req.Body)
			if err != nil {
				log.Printf("failed to copy request: %v", err)
			}
			req.Body.Close()
			pctx.UserData = &Context{RequestBody: body}
			req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
			return req, nil
		})
	s.proxy.OnResponse().Do(
		goproxy.HandleBytes(
			func(b []byte, pctx *goproxy.ProxyCtx) []byte {
				log.Printf("[INFO] resp %s", pctx.Resp.Status)
				ctx := pctx.UserData.(*Context)
				err := createCache(string(ctx.RequestBody), b, pctx.Req, pctx.Resp, s)
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
					resp, err = createHttpResponse(req, conn)
					if err != nil {
						log.Printf("[ERROR] %v", err)
					}
					req.Body = ioutil.NopCloser(bytes.NewBufferString(reqBody))
					return req, resp
				}
			}
			resp, err = createHttpErrorResponse(req)
			if err != nil {
				log.Printf("[ERROR] %v", err)
			}
			req.Body = ioutil.NopCloser(bytes.NewBufferString(reqBody))
			return req, resp
		})
}

func (s *Server) Start() {
	if s.config.record == true {
		log.Println("[INFO] All HTTP/S request and response is cached.")
		s.connectionCacheHandler()
	} else {
		s.mockOnlyHandler()
	}
	log.Println("[INFO] Running...")
	http.ListenAndServe(":8080", s.proxy)
}
