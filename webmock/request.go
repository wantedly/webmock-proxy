package webmock

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func getReqStruct(f *File) Request {
	b, err := readFile(f.Path)
	if err != nil {
		fmt.Println(err)
	}
	return parseReqStruct(jsonToStruct(b))
}

func parseReqStruct(conn *Connection) Request {
	return conn.Request
}

func readRequestBody(r *http.Request) (string, error) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return string(body), err
	}
	return string(body), nil
}
