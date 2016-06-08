package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

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

	// TODO cache からレスポンス作るときはこのチャネル使う
	// useCache := make(chan bool, 1)

	proxy.OnRequest().DoFunc(
		func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			defer r.Body.Close()
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				// TODO: Add error handling
				fmt.Println("Oops...")
			}
			c <- string(body)
			r.Body = ioutil.NopCloser(bytes.NewBufferString(string(body)))
			return r, nil
		})

	proxy.OnResponse().Do(
		goproxy.HandleBytes(
			func(b []byte, ctx *goproxy.ProxyCtx) []byte {
				respBody := <-c
				ConvertJsonFile(respBody, b, ctx)
				return b
			}))
	log.Fatal(http.ListenAndServe(":8080", proxy))
}
