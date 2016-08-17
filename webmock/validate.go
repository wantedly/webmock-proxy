package webmock

import (
	"net/http"

	"github.com/elazarl/goproxy"
	"github.com/jinzhu/gorm"
)

func validateRequest(r *http.Request, body string, db *gorm.DB) bool {
	file := getFileStruct(r)
	if !fileExists(file.Path) {
		return false
	}

	endpoint := dbConnectionCacheToStruct(db, r, file)
	if body == endpoint.Connections[0].Request.String {
		return true
	}

	// req, err = fileToReqStruct(file)
	// if err != nil {
	// 	return false
	// }
	// if (r.Header.Get("Content-Type") == req.Header.ContentType) &&
	// 	(body == req.String) &&
	// 	(r.Method == req.Method) &&
	// 	(file.URL == req.URL) {
	// 	return true
	// }
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
