package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func apiRequest(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return string(body)
}

func main() {
	fmt.Println(apiRequest("https://api.github.com"))
}
