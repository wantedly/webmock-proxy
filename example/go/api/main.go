package main

import (
	"github.com/wantedly/webmock-proxy/example/go/api/db"
	"github.com/wantedly/webmock-proxy/example/go/api/server"
)

// main ...
func main() {
	database := db.Connect()
	s := server.Setup(database)
	s.Run(":8080")
}
