package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/elazarl/goproxy"
)

func main() {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true

	// MEMO(munisystem):
	// Request Body を取得できるタイミングが OnRequest 内しか存在しない。
	// チャネル作って OnResponse 内に送る。
	// make でチャネルを作るときに第二引数を与えないと buffer size 0 となり、
	// ロックされてレスポンスが一生帰ってこなくなる。
	c := make(chan string, 1)

	// MEMO(munisystem):
	// 外部にリクエストが飛ばないように、環境変数でモードを切り替えるようにした
	// WEBMOCK_PROXY_RECORD = true
	// の場合は API を cache として取得するようになる。
	//
	// 通常では cache の有無を確認して、存在しない場合は 418 を返す。
	// しかしこれでは API のアップデートをしたい場合に一々起動しなおさなきゃいけなくなる。
	env := os.Getenv("WEBMOCK_PROXY_RECORD")
	if env == "true" {
		fmt.Println("webmock-proxy run record mode.")
		proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)
		proxy.OnRequest().DoFunc(
			func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
				body := readRequestBody(r)
				c <- body
				r.Body = ioutil.NopCloser(bytes.NewBufferString(body))
				return r, nil
			})

		proxy.OnResponse().Do(
			goproxy.HandleBytes(
				func(b []byte, ctx *goproxy.ProxyCtx) []byte {
					respBody := <-c
					createCacheFile(respBody, b, ctx)
					return b
				}))
	} else {
		proxy.OnRequest().DoFunc(
			func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
				body := readRequestBody(r)
				var resp *http.Response
				if validateRequest(r, body) {
					resp = newResponse(r)
					ctx.Logf("webmock-proxy use http request cache!!")
				} else {
					resp = newErrorResponse(r)
				}
				r.Body = ioutil.NopCloser(bytes.NewBufferString(body))
				return r, resp
			})
	}
	log.Fatal(http.ListenAndServe(":8080", proxy))
}
