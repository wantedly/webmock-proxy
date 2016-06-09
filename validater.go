package main

import (
	"net/http"
	"strconv"
	"strings"

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

// MEMO(munisystem):
// 現状は URL と request body が一致していた場合は response の状態にかかわらずキャッシュを使うようになっている。
// TTL を設けるか、cron で定期的に回すか...
// cron で対応する場合は response の validation を行って対応をする形になる
func validateResponse(ctx *goproxy.ProxyCtx, body string) bool {
	file := getFileStruct(ctx.Req)
	resp := getRespStruct(file)
	contentType := resp.Header.ContentType
	statusCode := resp.Status.Code
	cacheBody := resp.String
	statusArray := strings.Fields(ctx.Resp.Status)
	codeStr := statusArray[0]
	code, _ := strconv.Atoi(codeStr)

	if (ctx.Resp.Header.Get("Content-Type") == contentType) &&
		(body == cacheBody) &&
		(code == statusCode) {
		return true
	}
	return false
}
