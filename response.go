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

func getRespStruct(f file) response {
	return parseRespStruct(convertJSONToStruct(readFile(f.Path)))
}

func parseRespStruct(con connection) response {
	return con.Response
}
