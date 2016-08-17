package webmock

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/elazarl/goproxy"
	"github.com/jinzhu/gorm"
)

const (
	errms = `
{"message": "webmock-proxy fail to create response body."}
`
)

func newResponse(r *http.Request) (*http.Response, error) {
	file := getFileStruct(r)
	resp, err := getRespStruct(file)
	if err != nil {
		return goproxy.NewResponse(r, "application/json", http.StatusInternalServerError, errms), err
	}
	var header interface{}
	b := []byte(resp.Header)
	err = jsonToStruct(b, &header)
	if err != nil {
		return goproxy.NewResponse(r, "application/json", http.StatusInternalServerError, errms), err
	}
	contentType := header.(map[string]interface{})["Content-Type"].([]interface{})[0].(string)
	arr := strings.Fields(resp.Status)
	code, _ := strconv.Atoi(arr[0])
	body := resp.String
	return goproxy.NewResponse(r, contentType, code, body), nil
}

func newResponseFromDB(db *gorm.DB, r *http.Request) (*http.Response, error) {
	file := getFileStruct(r)
	endpoint := selectCache(db, r, file)
	resp := endpoint.Connections[0].Response
	var header interface{}
	b := []byte(resp.Header)
	err := jsonToStruct(b, &header)
	if err != nil {
		return goproxy.NewResponse(r, "application/json", http.StatusInternalServerError, errms), err
	}

	contentType := header.(map[string]interface{})["Content-Type"].([]interface{})[0].(string)
	arr := strings.Fields(resp.Status)
	code, _ := strconv.Atoi(arr[0])
	body := resp.String
	return goproxy.NewResponse(r, contentType, code, body), nil
}

func newErrorResponse(r *http.Request) (*http.Response, error) {
	body, err := createErrorMessage(r.URL.Host + r.URL.Path)
	if err != nil {
		return goproxy.NewResponse(r, "application/json", http.StatusInternalServerError, errms), err
	}
	return goproxy.NewResponse(r, "application/json", http.StatusTeapot, body), nil
}

func getRespStruct(f *File) (Response, error) {
	b, err := readFile(f.Path)
	if err != nil {
		return Response{}, err
	}
	var conn Connection
	err = jsonToStruct(b, &conn)
	if err != nil {
		return Response{}, err
	}
	return parseRespStruct(&conn), nil
}

func parseRespStruct(conn *Connection) Response {
	return conn.Response
}
