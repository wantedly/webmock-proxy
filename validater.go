package main

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/elazarl/goproxy"
)

func ValidateRequest(r *http.Request, body string) bool {
	filename := "webmock-cache.json"
	request := ReqStructParser(ConvertStruct(ReadFile(filename)))
	if r.Header.Get("Content-Type") == request.Header.ContentType && body == request.String && r.Method == request.Method && r.URL.Host+r.URL.Path == request.Url {
		return true
	}
	return false
}

func ValidateResponse(ctx *goproxy.ProxyCtx, body string) bool {
	resp := response()
	contentType := resp.Header.ContentType
	statusCode := resp.Status.Code
	cacheBody := resp.String

	statusArray := strings.Fields(ctx.Resp.Status)
	codeStr := statusArray[0]
	code, _ := strconv.Atoi(codeStr)

	if ctx.Resp.Header.Get("Content-Type") == contentType && body == cacheBody && code == statusCode {
		return true
	}
	return false
}
