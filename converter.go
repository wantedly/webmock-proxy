package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/elazarl/goproxy"
)

type HttpInteractions struct {
	Connection []Connection `json:"http_interactions"`
}

type Connection struct {
	Request    Request  `json:"request"`
	Response   Response `json:"response"`
	RecordedAt string   `json:"recorded_at"`
}

type Request struct {
	Header Header `json:"header"`
	String string `json:"string"`
	Method string `json:"method"`
	Url    string `json:"url"`
}

type Response struct {
	Status Status `json:"status"`
	Header Header `json:"header"`
	String string `json:"string"`
}

type Status struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Header struct {
	ContentType   string `json:"Content-Type"`
	ContentLength string `json:"Content-Length"`
}

func createCacheFile(respBody string, b []byte, ctx *goproxy.ProxyCtx) {
	file := getFileStruct(ctx.Req)
	req := createReqStruct(respBody, ctx)
	resp := createRespStruct(b, ctx)
	con := Connection{req, resp, ctx.Resp.Header.Get("Date")}
	arr := []Connection{con}
	httpInt := HttpInteractions{arr}
	recursiveWriteFile(convertStructToJSON(httpInt), file)
}

func convertStructToJSON(httpInt HttpInteractions) string {
	jsonBytes, err := json.Marshal(httpInt)
	if err != nil {
		// TODO: Add error handling
		fmt.Println("JSON Marshal Error: ", err)
		return "error"
	}
	out := new(bytes.Buffer)
	json.Indent(out, jsonBytes, "", "    ")
	return out.String()
}

func createReqStruct(respBody string, ctx *goproxy.ProxyCtx) Request {
	contentType := ctx.Req.Header.Get("Content-Type")
	contentLength := ctx.Req.Header.Get("Content-Length")
	header := Header{contentType, contentLength}

	method := ctx.Req.Method

	host := ctx.Req.URL.Host
	path := ctx.Req.URL.Path

	request := Request{header, respBody, method, host + path}
	return request
}

func createRespStruct(b []byte, ctx *goproxy.ProxyCtx) Response {
	statusArray := strings.Fields(ctx.Resp.Status)
	codeStr, message := statusArray[0], statusArray[1]
	code, _ := strconv.Atoi(codeStr)
	status := Status{code, message}

	contentType := ctx.Resp.Header.Get("Content-Type")
	contentLength := ctx.Resp.Header.Get("Content-Length")
	header := Header{contentType, contentLength}

	body := strings.TrimRight(string(b), "\n")

	response := Response{status, header, body}
	return response
}

func convertJSONToStruct(b []byte) HttpInteractions {
	var httpInt HttpInteractions
	err := json.Unmarshal(b, &httpInt)
	if err != nil {
		// TODO: Add error handling
		fmt.Println("JSON Marshal Error: ")
		return httpInt
	}
	return httpInt
}
