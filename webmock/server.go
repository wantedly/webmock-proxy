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
	bodyCh chan string
	headCh chan map[string][]string
}

func initDB(config *Config) (*gorm.DB, error) {
	if config.local == false {
		db, err := NewDBConnection()
		if err != nil {
			return nil, fmt.Errorf("Faild to connect database: %v", err)
		}
		log.Println("Use db.")
		return db, nil
	}
	log.Println("Use local cache files.")
	return nil, nil
}

func NewServer(config *Config) (*Server, error) {
	db, err := initDB(config)
	if err != nil {
		return nil, err
	}
	var bodyCh chan string
	var headCh chan map[string][]string
	if config.record == true {
		log.Println("All http/s request is cached.")
		bodyCh = make(chan string, 1)
		headCh = make(chan map[string][]string, 1)
	}

	return &Server{
		config: config,
		db:     db,
		proxy:  goproxy.NewProxyHttpServer(),
		bodyCh: bodyCh,
		headCh: headCh,
	}, nil
}

func (s *Server) connectionCacheHandler() {
	s.proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)
	s.proxy.OnRequest().DoFunc(
		func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {

			// DeepCopy *http.Request.Header (type: map[string][]string)
			reqHeader := make(map[string][]string, len(req.Header))
			for k, v := range req.Header {
				reqHeader[k] = v
			}

			reqBody, err := ioReader(req.Body)
			if err != nil {
				log.Println(err)
			}
			s.bodyCh <- reqBody
			s.headCh <- reqHeader
			req.Body = ioutil.NopCloser(bytes.NewBufferString(reqBody))
			return req, nil
		})
	s.proxy.OnResponse().Do(
		goproxy.HandleBytes(
			func(b []byte, ctx *goproxy.ProxyCtx) []byte {
				reqBody := <-s.bodyCh
				reqHeader := <-s.headCh
				ctx.Req.Header = reqHeader
				err := createCache(reqBody, b, ctx.Req, ctx.Resp, s)
				if err != nil {
					log.Println(err)
				}
				return b
			}))
}

func (s *Server) mockOnlyHandler() {
	s.proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)
	s.proxy.OnRequest().DoFunc(
		func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			reqBody, err := ioReader(req.Body)
			if err != nil {
				log.Println(err)
			}
			conn, err := NewConnection(req, s)
			if err != nil {
				log.Println(err)
			}
			var resp *http.Response
			if conn != nil {
				is, err := validateRequest(req, conn, reqBody)
				if err != nil {
					log.Println(err)
				}
				if is == true {
					resp, err = createHttpResponse(req, conn)
					if err != nil {
						log.Println(err)
					}
					req.Body = ioutil.NopCloser(bytes.NewBufferString(reqBody))
					return req, resp
				}
			}
			resp, err = createHttpErrorResponse(req)
			if err != nil {
				log.Println(err)
			}
			req.Body = ioutil.NopCloser(bytes.NewBufferString(reqBody))
			return req, resp
		})
}

func (s *Server) Start() {
	if s.config.record == true {
		s.connectionCacheHandler()
	} else {
		s.mockOnlyHandler()
	}
	log.Println("Running...")
	http.ListenAndServe(":8080", s.proxy)
}
