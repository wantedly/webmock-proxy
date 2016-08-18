package webmock

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"
)

type Endpoint struct {
	ID          uint `gorm:"primary_key;AUTO_INCREMENT" json:"-"`
	URL         string
	Connections []Connection
	Update      time.Time
}

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

func NewConnection(req *http.Request, s *Server) (*Connection, error) {
	var (
		root = "webmock-cache/"
		url  = req.URL.Host + req.URL.Path
		file = "cache.json"
		dst  = filepath.Join(root, url, file)
	)
	if s.config.local == true {
		b, err := readFile(dst)
		if err != nil {
			return nil, fmt.Errorf("Faild to load cache file: %v", err)
		}
		var conn Connection
		err = jsonToStruct(b, &conn)
		if err != nil {
			return nil, fmt.Errorf("Faild to convert json into struct: %v", err)
		}
		return &conn, nil
	}
	endpoint := findEndpoint(req.Method, url, s.db)
	if len(endpoint.Connections) == 0 {
		return nil, nil
	}
	return &endpoint.Connections[0], nil
}
