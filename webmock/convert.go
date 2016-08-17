package webmock

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/elazarl/goproxy"
)

type Connection struct {
	ID         uint     `gorm:"primary_key;AUTO_INCREMENT" json:"-"`
	EndpointID uint     `json:"-"`
	Request    Request  `json:"request"`
	Response   Response `json:"response"`
	RecordedAt string   `json:"recorded_at"`
}

type Request struct {
	ID           uint   `gorm:"primary_key;AUTO_INCREMENT" json:"-"`
	ConnectionID uint   `json:"-"`
	Header       Header `json:"header"`
	String       string `json:"string"`
	Method       string `json:"method"`
	URL          string `json:"url"`
}

type Response struct {
	ID           uint   `gorm:"primary_key;AUTO_INCREMENT" json:"-"`
	ConnectionID uint   `json:"-"`
	Status       string `json:"status"`
	Header       Header `json:"header"`
	String       string `json:"string"`
}

type Header struct {
	ID            uint   `gorm:"primary_key;AUTO_INCREMENT" json:"-"`
	RequestID     uint   `json:"-"`
	ResponseID    uint   `json:"-"`
	Status        string `json:"status"`
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

func jsonToStruct(b []byte) (*Connection, error) {
	var conn Connection
	err := json.Unmarshal(b, &conn)
	if err != nil {
		return &conn, err
	}
	return &conn, nil
}

func createReqStruct(body string, ctx *goproxy.ProxyCtx) Request {
	contentType := ctx.Req.Header.Get("Content-Type")
	contentLength := ctx.Req.Header.Get("Content-Length")
	header := Header{ContentType: contentType, ContentLength: contentLength}
	method := ctx.Req.Method
	host := ctx.Req.URL.Host
	path := ctx.Req.URL.Path

	return Request{Header: header, String: body, Method: method, URL: host + path}
}

func createRespStruct(b []byte, ctx *goproxy.ProxyCtx) Response {
	contentType := ctx.Resp.Header.Get("Content-Type")
	contentLength := ctx.Resp.Header.Get("Content-Length")
	header := Header{ContentType: contentType, ContentLength: contentLength}
	body := strings.TrimRight(string(b), "\n")

	return Response{Status: ctx.Resp.Status, Header: header, String: body}
}

func createErrorMessage(str string) (string, error) {
	mes := "Not found webmock-proxy cache. URL: " + str
	body := &ResponseBody{Message: mes}
	return structToJSON(body)
}
