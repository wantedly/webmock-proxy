package webmock

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/elazarl/goproxy"
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
	contentType := resp.Header.ContentType
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

func getRespStruct(f *file) (*response, error) {
	b, err := readFile(f.Path)
	if err != nil {
		return nil, err
	}
	return parseRespStruct(jsonToStruct(b)), nil
}

func parseRespStruct(conn *connection) *response {
	return conn.Response
}
