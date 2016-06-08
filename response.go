package main

import (
	"github.com/elazarl/goproxy"
	"net/http"
)

func NewResponse(r *http.Request) *http.Response {
	resp := response()
	contentType := resp.Header.ContentType
	statusCode := resp.Status.Code
	body := resp.String
	return goproxy.NewResponse(r, contentType, statusCode, body)
}

func response() Response {
	filename := "webmock-cache.json"
	return RespStructParser(ConvertStruct(ReadFile(filename)))
}

func RespStructParser(httpInt HttpInteractions) Response {
	return httpInt.Connection[0].Response
}

func ReqStructParser(httpInt HttpInteractions) Request {
	return httpInt.Connection[0].Request
}
