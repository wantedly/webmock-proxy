package webmock

import (
	"log"
	"net/http"
	"path/filepath"
	"time"
)

func createCache(body string, b []byte, req *http.Request, resp *http.Response, s *Server) error {
	var (
		root = s.config.cacheDir + "/"
		url  = req.URL.Host + req.URL.Path
		file = req.Method + ".json"
		dst  = filepath.Join(root, url, file)
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
	if s.config.local == true {
		if err := writeFile(string(byteArr), dst); err != nil {
			return err
		}
		log.Printf("[INFO] Create HTTP/S connection cache.")
		return nil
	}
	var conns []Connection
	conns = append(conns, conn)
	endpoint := &Endpoint{
		URL:         url,
		Connections: conns,
		Update:      time.Now(),
	}
	ce := readEndpoint(url, s.db)
	if len(ce.Connections) != 0 {
		for _, v := range ce.Connections {
			deleteConnection(&v, s.db)
			if v.Request.Method == req.Method {
				continue
			}
			conns = append(conns, v)
		}
		endpoint.Connections = conns
		updateEndpoint(ce, endpoint, s.db)

		log.Printf("[INFO] Update HTTP/S connection cache.")
		return nil
	}
	if err := insertEndpoint(endpoint, s.db); err != nil {
		return err
	}
	log.Printf("[INFO] Create HTTP/S connection cache.")
	return nil
}
