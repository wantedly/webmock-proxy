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

	c := make(chan string, 1)
	env := os.Getenv("WEBMOCK_PROXY_RECORD")
	if env == "1" {
		log.Println("webmock-proxy run record mode.")
		proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)
		proxy.OnRequest().DoFunc(
			func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
				body, err := readRequestBody(r)
				if err != nil {
					log.Println(err)
				}
				c <- body
				r.Body = ioutil.NopCloser(bytes.NewBufferString(body))
				return r, nil
			})
		proxy.OnResponse().Do(
			goproxy.HandleBytes(
				func(b []byte, ctx *goproxy.ProxyCtx) []byte {
					body := <-c
					err := createCache(body, b, ctx, db)
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
