package main

import (
	"log"

	"github.com/wantedly/webmock-proxy/webmock"
)

func main() {
	config, err := webmock.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	server, err := webmock.NewServer(config)
	if err != nil {
		log.Fatal(err)
	}
	server.Start()
}
