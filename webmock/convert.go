package webmock

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
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

// Webmock-proxy need to validate http response header because of using cache data.
// But, after converting Request.Header (json str) into Struct using json.Marshal(),
// Struct type has changed map[string][]interface{} (original: map[string][]string)
// When validating http request, execute mapToMapInterface("http request header cache").
func mapToMapInterface(m map[string][]string) map[string]interface{} {
	mi := make(map[string]interface{}, len(m))
	for k, v := range m {
		si := make([]interface{}, len(v))
		for vk, vv := range v {
			si[vk] = vv
		}
		mi[k] = si
	}
	return mi
}

func createReqStruct(body string, req *http.Request) (Request, error) {
	method := req.Method
	host := req.URL.Host
	path := req.URL.Path
	header, err := structToJSON(req.Header, false)
	if err != nil {
		return Request{}, err
	}
	return Request{Header: string(header), String: body, Method: method, URL: host + path}, nil
}

func createRespStruct(b []byte, resp *http.Response) (Response, error) {
	body := strings.TrimRight(string(b), "\n")
	header, err := structToJSON(resp.Header, false)
	if err != nil {
		return Response{}, err
	}
	return Response{Status: resp.Status, Header: string(header), String: body}, nil
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
