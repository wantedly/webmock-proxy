package webmock

import (
	"bytes"
	"encoding/json"
	"log"
	"strings"

	"github.com/elazarl/goproxy"
)

type Connection struct {
	ID         uint     `gorm:"primary_key;AUTO_INCREMENT" json:"-"`
	EndpointID uint     `json:"-"`
	Request    Request  `json:"request"`
	Response   Response `json:"response"`
	RecordedAt string   `json:"recorded_at"`
}

type Request struct {
	ID           uint   `gorm:"primary_key;AUTO_INCREMENT" json:"-"`
	ConnectionID uint   `json:"-"`
	Header       string `json:"header"`
	String       string `json:"string"`
	Method       string `json:"method"`
	URL          string `json:"url"`
}

type Response struct {
	ID           uint   `gorm:"primary_key;AUTO_INCREMENT" json:"-"`
	ConnectionID uint   `json:"-"`
	Status       string `json:"status"`
	Header       string `json:"header"`
	String       string `json:"string"`
}

type ResponseBody struct {
	Message string `json:"message"`
}

func structToJSON(v interface{}, indentFlag ...bool) ([]byte, error) {
	byteArr, err := json.Marshal(v)
	if err != nil {
		return make([]byte, 0), err
	}
	if indentFlag != nil {
		return byteArr, nil
	}
	out := new(bytes.Buffer)
	json.Indent(out, byteArr, "", "    ")
	return out.Bytes(), nil
}

func jsonToStruct(b []byte, v interface{}) error {
	err := json.Unmarshal(b, v)
	if err != nil {
		return err
	}
	return nil
}

func mapToMapInterface(m map[string][]string) map[string]interface{} {
	mi := make(map[string]interface{}, len(m))
	for k, v := range m {
		mi[k] = v
	}
	return mi
}

// FIXME:
// *http.Request.Header => map[User-Agent:[Go-http-client/1.1] Accept-Encoding:[gzip]]
// *goproxy.ProxyCtx.Req.Header => map[User-Agent:[Go-http-client/1.1]]
func createReqStruct(body string, ctx *goproxy.ProxyCtx) Request {
	log.Println(ctx.Req.Header)
	method := ctx.Req.Method
	host := ctx.Req.URL.Host
	path := ctx.Req.URL.Path
	header, err := structToJSON(mapToMapInterface(ctx.Req.Header), false)
	if err != nil {
		//TODO
		log.Println(err)
	}

	return Request{Header: string(header), String: body, Method: method, URL: host + path}
}

func createRespStruct(b []byte, ctx *goproxy.ProxyCtx) (Response, error) {
	body := strings.TrimRight(string(b), "\n")
	header, err := structToJSON(mapToMapInterface(ctx.Resp.Header), false)
	if err != nil {
		return Response{}, err
	}
	return Response{Status: ctx.Resp.Status, Header: string(header), String: body}, nil
}

func createErrorMessage(str string) (string, error) {
	mes := "Not found webmock-proxy cache. URL: " + str
	body := &ResponseBody{Message: mes}
	byteArr, err := structToJSON(body)
	if err != nil {
		return "", err
	}
	return string(byteArr), nil
}
