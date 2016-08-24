package webmock

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type responseBody struct {
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

func requestStruct(body string, req *http.Request) (Request, error) {
	method := req.Method
	host := req.URL.Host
	path := req.URL.Path
	header, err := structToJSON(req.Header, false)
	if err != nil {
		return Request{}, err
	}
	return Request{Header: string(header), String: body, Method: method, URL: host + path}, nil
}

func responseStruct(b []byte, resp *http.Response) (Response, error) {
	header, err := structToJSON(resp.Header, false)
	if err != nil {
		return Response{}, err
	}
	return Response{Status: resp.Status, Header: string(header), String: string(b)}, nil
}
