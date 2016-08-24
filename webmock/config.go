package webmock

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	record    = false
	local     = false
	port      = ":8080"
	cacheDir  = "cache"
	masterURL = ""
)

type Config struct {
	record    bool
	local     bool
	port      string
	cacheDir  string
	masterURL string
}

func NewConfig() (*Config, error) {
	if os.Getenv("WEBMOCK_PROXY_RECORD") == "1" {
		record = true
	}
	if os.Getenv("WEBMOCK_PROXY_LOCAL") == "1" {
		local = true
		if url := os.Getenv("WEBMOCK_PROXY_MASTER_URL"); url != "" {
			masterURL = url
		}
	}
	if portStr := os.Getenv("WEBMOCK_PROXY_PORT"); portStr != "" {
		portInt, err := strconv.Atoi(portStr)
		if err != nil {
			return nil, fmt.Errorf("[ERROR] Illegal value in $WEBMOCK_PROXY_PORT: %v", err)
		}
		port = ":" + strconv.Itoa(portInt)
	}
	if str := os.Getenv("WEBMOCK_PROXY_CACHE_DIR"); str != "" {
		str := strings.TrimRight(str, "/")
		cacheDir = str
	}
	return &Config{
		record:    record,
		local:     local,
		port:      port,
		cacheDir:  cacheDir,
		masterURL: masterURL,
	}, nil
}
