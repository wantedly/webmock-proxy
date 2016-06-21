package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"unsafe"
)

type file struct {
	URL     string
	RootDir string
	Path    string
	Dir     string
	Name    string
}

// TODO: Add error handling
func writeFile(str string, path string) {
	b := *(*[]byte)(unsafe.Pointer(&str))
	err := ioutil.WriteFile(path, b, 0644)
	if err != nil {
		fmt.Println("Cannot Write file")
	}
}

func readFile(path string) []byte {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return []byte("error")
	}
	return b
}

func createDir(dir string) {
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		fmt.Println("Cannot create dir")
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func recursiveWriteFile(str string, f file) {
	if fileExists(f.Path) || fileExists(f.Dir) {
		writeFile(str, f.Path)
		return
	}
	createDir(f.Dir)
	writeFile(str, f.Path)
}

func getFileStruct(r *http.Request) file {
	rootDir := "webmock-cache/"
	url := r.URL.Host + r.URL.Path
	name := "cache.json"
	dir := rootDir + url
	path := dir + "/" + name
	return file{url, rootDir, path, dir, name}
}
