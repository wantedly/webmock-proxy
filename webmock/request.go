package webmock

import (
	"io/ioutil"
	"net/http"

	"github.com/jinzhu/gorm"
)

func fileToReqStruct(f *File) (Request, error) {
	b, err := readFile(f.Path)
	if err != nil {
		return Request{}, err
	}
	conn, err := jsonToStruct(b)
	if err != nil {
		return Request{}, err
	}
	return parseReqStruct(conn), nil
}

func parseReqStruct(conn *Connection) Request {
	return conn.Request
}

func dbConnectionCacheToStruct(db *gorm.DB, r *http.Request, file *File) Endpoint {
	endpoint := selectCache(db, r, file)
	return endpoint
}

func readRequestBody(r *http.Request) (string, error) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return string(body), err
	}
	return string(body), nil
}
