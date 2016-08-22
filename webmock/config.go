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
	port      = 8080
	cacheDir  = "cache"
	syncCache = false
	masterURL = ""
)

type Config struct {
	record    bool
	local     bool
	port      int
	cacheDir  string
	syncCache bool
	masterURL string
}

func NewConfig() (*Config, error) {
	if os.Getenv("WEBMOCK_PROXY_RECORD") == "1" {
		record = true
	}
	if os.Getenv("WEBMOCK_PROXY_LOCAL") == "1" {
		local = true
		if os.Getenv("WEBMOCK_PROXY_SYNC_CACHE") == "1" {
			if url := os.Getenv("WEBMOCK_PROXY_MASTER_URL"); url != "" {
				masterURL = url
				syncCache = true
			} else {
				return nil, fmt.Errorf("[ERROR] Prease set $WEBMOCK_PROXY_MASTER_URL")
			}
		}
	}
	if portStr := os.Getenv("WEBMOCK_PROXY_PORT"); portStr != "" {
		var err error
		port, err = strconv.Atoi(portStr)
		if err != nil {
			return nil, fmt.Errorf("[ERROR] Illegal value in $WEBMOCK_PROXY_PORT: %v", err)
		}
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
		syncCache: syncCache,
		masterURL: masterURL,
	}, nil
}
