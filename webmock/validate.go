package webmock

import (
	"net/http"
	"reflect"
)

func validateRequest(req *http.Request, conn *Connection, body string) (bool, error) {
	var header interface{}
	b := []byte(conn.Request.Header)
	err := jsonToStruct(b, &header)
	if err != nil {
		return false, err
	}
	if (body == conn.Request.String) &&
		(reflect.DeepEqual(mapToMapInterface(req.Header), header) == true) &&
		(req.Method == conn.Request.Method) {
		return true, nil
	}
	return false, nil
}
