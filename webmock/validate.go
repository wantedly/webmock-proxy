package webmock

import (
	"net/http"
	"reflect"

	"github.com/elazarl/goproxy"
	"github.com/jinzhu/gorm"
)

func validateRequest(r *http.Request, body string, db *gorm.DB) bool {
	file := getFileStruct(r)
	if !fileExists(file.Path) {
		return false
	}

	// req, err := fileToReqStruct(file)
	// if err != nil {
	// 	return false
	// }

	endpoint := selectCache(db, r, file)
	req := endpoint.Connections[0].Request
	var header interface{}
	b := []byte(req.Header)
	err := jsonToStruct(b, &header)
	if err != nil {
		return false
	}

	if (body == endpoint.Connections[0].Request.String) &&
		(reflect.DeepEqual(mapToMapInterface(r.Header), header) == true) {
		return true
	}

	// if (reflect.DeepEqual(mapToMapInterface(r.Header), header) == true) &&
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
	var header interface{}
	b := []byte(resp.Header)
	err = jsonToStruct(b, &header)
	if err != nil {
		return false
	}
	contentType := header.(map[string]interface{})["Content-Type"].([]interface{})[0].(string)
	statusCode := resp.Status
	cacheBody := resp.String

	if (ctx.Resp.Header.Get("Content-Type") == contentType) &&
		(body == cacheBody) &&
		(ctx.Resp.Status == statusCode) {
		return true
	}
	return false
}
