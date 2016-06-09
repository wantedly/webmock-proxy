package main

import (
	"net/http"

	"github.com/elazarl/goproxy"
)

func newResponse(r *http.Request) *http.Response {
	file := getFileStruct(r)
	resp := getRespStruct(file)
	contentType := resp.Header.ContentType
	statusCode := resp.Status.Code
	body := resp.String
	return goproxy.NewResponse(r, contentType, statusCode, body)
}

func getRespStruct(file File) Response {
	return parseRespStruct(convertJSONToStruct(readFile(file.Path)))
}

func parseRespStruct(httpInt HttpInteractions) Response {
	return httpInt.Connection[0].Response
}
