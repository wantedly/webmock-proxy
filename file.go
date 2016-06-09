package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"unsafe"
)

type File struct {
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
		fmt.Println("Cannot create file")
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func recursiveWriteFile(str string, file File) {
	if fileExists(file.Path) || fileExists(file.Dir) {
		writeFile(str, file.Path)
		return
	}
	createDir(file.Dir)
	writeFile(str, file.Path)
}

func getFileStruct(r *http.Request) File {
	rootDir := "cache/"
	url := r.URL.Host + r.URL.Path
	path := rootDir + url
	arr := strings.Split(path, "/")
	name := arr[len(arr)-1]
	dir := strings.TrimRight(path, name)
	return File{url, rootDir, path, dir, name}
}
