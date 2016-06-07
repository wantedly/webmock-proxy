package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/elazarl/goproxy"
)

type Response struct {
	Status         Status         `json:"status"`
	ResponseHeader ResponseHeader `json:"response"`
	ResponseBody   ResponseBody   `json:"body"`
}

type Status struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ResponseHeader struct {
	ContentType   string `jdon:"Content-Type"`
	ContentLength string `jdon:"Content-Length"`
}

type ResponseBody struct {
	String string `jdon:"string"`
}

func ConvertJson(b []byte, ctx *goproxy.ProxyCtx) {
	statusArray := strings.Fields(ctx.Resp.Status)
	codeStr, message := statusArray[0], statusArray[1]
	code, _ := strconv.Atoi(codeStr)
	status := Status{code, message}

	contentType := ctx.Resp.Header.Get("Content-Type")
	contentLength := ctx.Resp.Header.Get("Content-Length")
	responseHeader := ResponseHeader{contentType, contentLength}

	bodyStr := string(b)
	responseBody := ResponseBody{bodyStr}

	response := Response{status, responseHeader, responseBody}

	// Marshal
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		fmt.Println("JSON Marshal Error: ", err)
		return
	}
	out := new(bytes.Buffer)
	json.Indent(out, jsonBytes, "", "    ")
	fmt.Println(out.String())
}
