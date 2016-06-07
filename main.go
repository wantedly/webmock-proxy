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
		func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			return r, goproxy.NewResponse(r, "application/json", http.StatusForbidden, "{\"sample\":\"json\"}")
		})
	proxy.OnResponse().Do(
		goproxy.HandleBytes(
			func(b []byte, ctx *goproxy.ProxyCtx) []byte {
				ConvertYaml(b, ctx)
				return b
			}))
	log.Fatal(http.ListenAndServe(":8080", proxy))
}
