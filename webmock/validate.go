package webmock

import (
	"net/http"

	"github.com/elazarl/goproxy"
)

func validateRequest(r *http.Request, body string) bool {
	file := getFileStruct(r)
	if !fileExists(file.Path) {
		return false
	}

	req := getReqStruct(file)
	if (r.Header.Get("Content-Type") == req.Header.ContentType) &&
		(body == req.String) &&
		(r.Method == req.Method) &&
		(file.URL == req.URL) {
		return true
	}
	return false
}

func validateResponse(ctx *goproxy.ProxyCtx, body string) bool {
	file := getFileStruct(ctx.Req)
	resp, err := getRespStruct(file)
	if err != nil {
		return false
	}
	contentType := resp.Header.ContentType
	statusCode := resp.Status
	cacheBody := resp.String

	if (ctx.Resp.Header.Get("Content-Type") == contentType) &&
		(body == cacheBody) &&
		(ctx.Resp.Status == statusCode) {
		return true
	}
	return false
}
