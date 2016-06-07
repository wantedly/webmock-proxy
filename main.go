package main

import (
	"log"
	"net/http"

	"github.com/elazarl/goproxy"
)

func main() {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true
  proxy.OnRequest().DoFunc(
    func(r *http.Request, ctx *goproxy.ProxyCtx)(*http.Request, *http.Response) {
			return r, goproxy.NewResponse(r, "application/json", http.StatusForbidden, "{\"sample\":\"json\"}")
	})
	log.Fatal(http.ListenAndServe(":8080", proxy))
}
