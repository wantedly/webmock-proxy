package webmock

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/elazarl/goproxy"
)

type Connection struct {
	Request    *Request  `json:"request"`
	Response   *Response `json:"response"`
	RecordedAt string    `json:"recorded_at"`
}

type Request struct {
	Header *Header `json:"header"`
	String string  `json:"string"`
	Method string  `json:"method"`
	URL    string  `json:"url"`
}

type Response struct {
	Status string  `json:"status"`
	Header *Header `json:"header"`
	String string  `json:"string"`
}

type Header struct {
	ContentType   string `json:"Content-Type"`
	ContentLength string `json:"Content-Length"`
}

type ResponseBody struct {
	Message string `json:"message"`
}

func structToJSON(v interface{}) (string, error) {
	jsonBytes, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	out := new(bytes.Buffer)
	json.Indent(out, jsonBytes, "", "    ")
	return out.String(), nil
}

func jsonToStruct(b []byte) *Connection {
	var conn Connection
	err := json.Unmarshal(b, &conn)
	if err != nil {
		// TODO: Add error handling
		fmt.Println("JSON Marshal Error: ")
		return &conn
	}
	return &conn
}

func createReqStruct(body string, ctx *goproxy.ProxyCtx) *Request {
	contentType := ctx.Req.Header.Get("Content-Type")
	contentLength := ctx.Req.Header.Get("Content-Length")
	header := &Header{contentType, contentLength}
	method := ctx.Req.Method
	host := ctx.Req.URL.Host
	path := ctx.Req.URL.Path

	return &Request{header, body, method, host + path}
}

func createRespStruct(b []byte, ctx *goproxy.ProxyCtx) *Response {
	contentType := ctx.Resp.Header.Get("Content-Type")
	contentLength := ctx.Resp.Header.Get("Content-Length")
	header := &Header{contentType, contentLength}
	body := strings.TrimRight(string(b), "\n")

	return &Response{ctx.Resp.Status, header, body}
}

func createErrorMessage(str string) (string, error) {
	mes := "Not found webmock-proxy cache. URL: " + str
	body := &ResponseBody{mes}
	return structToJSON(body)
}
