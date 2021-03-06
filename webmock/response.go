package webmock

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
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

func createHttpResponse(req *http.Request, conn *Connection) (*http.Response, error) {
	resp := conn.Response
	var header interface{}
	b := []byte(resp.Header)
	err := jsonToStruct(b, &header)
	if err != nil {
		return goproxy.NewResponse(req, "application/json", http.StatusInternalServerError, errms), err
	}
	body := resp.String
	log.Printf("[INFO] Create HTTP/S response using connection cache.")
	fmt.Println(body)
	return newResponse(req, &resp, header)
}

func createHttpErrorResponse(r *http.Request) (*http.Response, error) {
	body, err := errorMessage(r.URL.Host + r.URL.Path)
	if err != nil {
		return goproxy.NewResponse(r, "application/json", http.StatusInternalServerError, errms), err
	}
	log.Printf("[INFO] Not match http connection cache.")
	fmt.Println(body)
	return goproxy.NewResponse(r, "application/json", http.StatusTeapot, body), nil
}

func errorMessage(url string) (string, error) {
	mes := "Not found webmock-proxy cache. URL: " + url
	body := &responseBody{Message: mes}
	byteArr, err := structToJSON(body)
	if err != nil {
		return "", err
	}
	return string(byteArr), nil
}

func newResponse(req *http.Request, resp *Response, header interface{}) (*http.Response, error) {
	r := &http.Response{}
	r.Request = req
	r.Header = make(http.Header)
	for k, v := range header.(map[string]interface{}) {
		for _, vv := range v.([]interface{}) {
			r.Header.Add(k, vv.(string))
		}
	}
	r.Status = resp.Status
	split := strings.Split(resp.Status, " ")
	status, err := strconv.Atoi(split[0])
	if err != nil {
		return nil, err
	}
	r.StatusCode = status
	buf := bytes.NewBufferString(resp.String)
	r.ContentLength = int64(buf.Len())
	r.Body = ioutil.NopCloser(buf)
	return r, nil
}
