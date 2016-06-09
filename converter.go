package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/elazarl/goproxy"
)

type connection struct {
	Request    request  `json:"request"`
	Response   response `json:"response"`
	RecordedAt string   `json:"recorded_at"`
}

type request struct {
	Header header `json:"header"`
	String string `json:"string"`
	Method string `json:"method"`
	URL    string `json:"url"`
}

type response struct {
	Status string `json:"status"`
	Header header `json:"header"`
	String string `json:"string"`
}

type header struct {
	ContentType   string `json:"Content-Type"`
	ContentLength string `json:"Content-Length"`
}

func createCacheFile(respBody string, b []byte, ctx *goproxy.ProxyCtx) {
	file := getFileStruct(ctx.Req)
	req := createReqStruct(respBody, ctx)
	resp := createRespStruct(b, ctx)
	con := connection{req, resp, ctx.Resp.Header.Get("Date")}
	recursiveWriteFile(convertStructToJSON(con), file)
}

func convertStructToJSON(con connection) string {
	jsonBytes, err := json.Marshal(con)
	if err != nil {
		// TODO: Add error handling
		fmt.Println("JSON Marshal Error: ", err)
		return "error"
	}
	out := new(bytes.Buffer)
	json.Indent(out, jsonBytes, "", "    ")
	return out.String()
}

func convertJSONToStruct(b []byte) connection {
	var con connection
	err := json.Unmarshal(b, &con)
	if err != nil {
		// TODO: Add error handling
		fmt.Println("JSON Marshal Error: ")
		return con
	}
	return con
}

func createReqStruct(respBody string, ctx *goproxy.ProxyCtx) request {
	contentType := ctx.Req.Header.Get("Content-Type")
	contentLength := ctx.Req.Header.Get("Content-Length")
	header := header{contentType, contentLength}
	method := ctx.Req.Method
	host := ctx.Req.URL.Host
	path := ctx.Req.URL.Path

	return request{header, respBody, method, host + path}
}

func createRespStruct(b []byte, ctx *goproxy.ProxyCtx) response {
	contentType := ctx.Resp.Header.Get("Content-Type")
	contentLength := ctx.Resp.Header.Get("Content-Length")
	header := header{contentType, contentLength}
	body := strings.TrimRight(string(b), "\n")

	return response{ctx.Resp.Status, header, body}
}
