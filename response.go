package main

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/elazarl/goproxy"
)

func newResponse(r *http.Request) *http.Response {
	file := getFileStruct(r)
	resp := getRespStruct(file)
	contentType := resp.Header.ContentType
	arr := strings.Fields(resp.Status)
	code, _ := strconv.Atoi(arr[0])
	body := resp.String
	return goproxy.NewResponse(r, contentType, code, body)
}

func newErrorResponse(r *http.Request) *http.Response {
	body := createErrorMessage(r.URL.Host + r.URL.Path)
	return goproxy.NewResponse(r, "application/json", http.StatusTeapot, body)
}

func getRespStruct(f file) response {
	return parseRespStruct(convertJSONToStruct(readFile(f.Path)))
}

func parseRespStruct(con connection) response {
	return con.Response
}
