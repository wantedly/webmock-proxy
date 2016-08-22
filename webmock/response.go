package webmock

import (
	"fmt"
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
	contentType := header.(map[string]interface{})["Content-Type"].([]interface{})[0].(string)
	arr := strings.Fields(resp.Status)
	code, _ := strconv.Atoi(arr[0])
	body := resp.String
	log.Printf("[INFO] Create HTTP/S response using connection cache.")
	fmt.Println(body)
	return goproxy.NewResponse(req, contentType, code, body), nil
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
