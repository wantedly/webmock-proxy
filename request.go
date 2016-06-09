package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func getReqStruct(f file) request {
	return parseReqStruct(convertJSONToStruct(readFile(f.Path)))
}

func parseReqStruct(con connection) request {
	return con.Request
}

func readRequestBody(r *http.Request) string {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		// TODO: Add error handling
		fmt.Println("Oops...")
	}
	return string(body)
}
