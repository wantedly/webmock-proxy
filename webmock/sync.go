package webmock

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
)

func sync(config *Config) error {
	proxyURL, _ := url.Parse(config.masterURL)
	http.DefaultTransport = &http.Transport{
		Proxy:              http.ProxyURL(proxyURL),
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
	}
	files, err := readFilePaths(config.cacheDir)
	if err != nil {
		return err
	}
	for _, file := range files {
		dst := filepath.Join(config.cacheDir, file)
		b, err := readFile(dst)
		if err != nil {
			return err
		}
		var conn Connection
		err = jsonToStruct(b, &conn)
		if err != nil {
			return err
		}
		var header interface{}
		b = []byte(conn.Request.Header)
		err = jsonToStruct(b, &header)
		if err != nil {
			return err
		}
		dir := filepath.Dir(file)
		base := filepath.Base(file)
		method := strings.TrimRight(base, ".json")

		// e.g. example.com:443/api/users => [example.com 443/api/users]
		x := strings.Split(dir, ":")

		// e.g. [example.com 443/api/users] => [443 api/users]
		x2 := strings.Split(x[1], "/")
		port := x2[0]
		url := filepath.Join(x[0], x2[1])
		switch port {
		case "443":
			url = "https://" + url
		default:
			url = "http//" + base
		}
		req, err := http.NewRequest(
			method,
			url,
			bytes.NewBuffer([]byte(conn.Request.String)),
		)
		for k, v := range header.(map[string]interface{}) {
			for _, vv := range v.([]interface{}) {
				req.Header.Set(k, vv.(string))
			}
		}
		client := new(http.Client)
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		respStruct, err := responseStruct(body, resp)
		if err != nil {
			return err
		}
		conn = Connection{Request: conn.Request, Response: respStruct, RecordedAt: resp.Header.Get("Date")}
		byteArr, err := structToJSON(conn)
		if err != nil {
			return err
		}
		if err = writeFile(string(byteArr), dst); err != nil {
			return err
		}
	}
	return nil
}
