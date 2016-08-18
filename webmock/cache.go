package webmock

import (
	"net/http"
	"path/filepath"
	"time"
)

func createCache(body string, b []byte, req *http.Request, resp *http.Response, s *Server) error {
	var (
		root = "webmock-cache"
		url  = req.URL.Host + req.URL.Path
		dir  = root + url
		file = "cache.json"
		dst  = filepath.Join(dir, file)
	)

	reqStruct, err := requestStruct(body, req)
	if err != nil {
		return err
	}
	respStruct, err := responseStruct(b, resp)
	if err != nil {
		return err
	}
	conn := Connection{Request: reqStruct, Response: respStruct, RecordedAt: resp.Header.Get("Date")}
	byteArr, err := structToJSON(conn)
	if err != nil {
		return err
	}
	var conns []Connection
	conns = append(conns, conn)
	endpoint := &Endpoint{
		URL:         url,
		Connections: conns,
		Update:      time.Now(),
	}
	if s.config.local == true {
		if err := writeFile(string(byteArr), dst); err != nil {
			return err
		}
	}
	if err := insertEndpoint(endpoint, s.db); err != nil {
		return err
	}
	return nil
}
