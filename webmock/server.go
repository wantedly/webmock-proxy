package webmock

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/elazarl/goproxy"
)

func Server() {
	proxy := goproxy.NewProxyHttpServer()

	db, err := Connect()
	if err != nil {
		log.Fatal(err)
	}

	bCh := make(chan string, 1)
	hCh := make(chan map[string][]string, 1)
	env := os.Getenv("WEBMOCK_PROXY_RECORD")
	if env == "1" {
		log.Println("webmock-proxy run record mode.")
		proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)
		proxy.OnRequest().DoFunc(
			func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
				header := make(map[string][]string)

				// DeepCopy *http.Request.Header (type: map[string][]string)
				for k, v := range req.Header {
					header[k] = v
				}
				body, err := readRequestBody(req)
				if err != nil {
					log.Println(err)
				}
				bCh <- body
				hCh <- header
				req.Body = ioutil.NopCloser(bytes.NewBufferString(body))
				return req, nil
			})
		proxy.OnResponse().Do(
			goproxy.HandleBytes(
				func(b []byte, ctx *goproxy.ProxyCtx) []byte {
					reqBody := <-bCh
					reqHeader := <-hCh
					ctx.Req.Header = reqHeader
					err = createCache(reqBody, b, ctx.Req, ctx.Resp, db)
					if err != nil {
						log.Println(err)
					}
					return b
				}))
	} else {
		proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)
		proxy.OnRequest().DoFunc(
			func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
				body, err := readRequestBody(r)
				if err != nil {
					log.Println(err)
				}
				var resp *http.Response
				if validateRequest(r, body, db) {
					// resp, err = newResponse(r)
					resp, err = newResponseFromDB(db, r)
					if err != nil {
						log.Println(err)
					}
				} else {
					resp, err = newErrorResponse(r)
					if err != nil {
						log.Println(err)
					}
				}
				r.Body = ioutil.NopCloser(bytes.NewBufferString(body))
				return r, resp
			})
	}
	http.ListenAndServe(":8080", proxy)
}
