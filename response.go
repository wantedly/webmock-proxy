package main

import (
	"github.com/elazarl/goproxy"
	"net/http"
)

type RespVar struct {
	ContentType string
	StatusCode  int
	Body        string
}

func NewResponse(r *http.Request) *http.Response {
	respVar := respVar()
	return goproxy.NewResponse(r, respVar.ContentType, respVar.StatusCode, respVar.Body)
}

func respVar() RespVar {
	filename := "webmock-cache.json"
	return structParser(ConvertStruct(ReadFile(filename)))
}

func structParser(httpInt HttpInteractions) RespVar {
	contentType := httpInt.Connection[0].Response.Header.ContentType
	statusCode := httpInt.Connection[0].Response.Status.Code
	body := httpInt.Connection[0].Response.String
	return RespVar{contentType, statusCode, body}
}
