package webmock

import (
	"io/ioutil"
	"net/http"
	"os"
	"time"
	"unsafe"

	"github.com/elazarl/goproxy"
	"github.com/jinzhu/gorm"
)

type File struct {
	URL  string
	Root string
	Path string
	Dir  string
	Name string
}

type Endpoint struct {
	ID          uint `gorm:"primary_key;AUTO_INCREMENT" json:"-"`
	URL         string
	Connections []Connection
	Update      time.Time
}

func createCache(body string, b []byte, ctx *goproxy.ProxyCtx, db *gorm.DB) error {
	file := getFileStruct(ctx.Req)
	req := createReqStruct(body, ctx)
	resp := createRespStruct(b, ctx)
	conn := Connection{Request: req, Response: resp, RecordedAt: ctx.Resp.Header.Get("Date")}
	jsonStr, err := structToJSON(conn)
	if err != nil {
		return err
	}
	if err = writeFile(jsonStr, file); err != nil {
		return err
	}
	var conns []Connection
	conns = append(conns, conn)
	endpoint := &Endpoint{
		URL:         file.URL,
		Connections: conns,
		Update:      time.Now(),
	}
	if err := insertCache(endpoint, db); err != nil {
		return err
	}
	return nil
}

func writeFile(str string, f *File) error {
	if !fileExists(f.Path) && !fileExists(f.Dir) {
		err := mkdir(f.Dir)
		if err != nil {
			return err
		}
	}
	b := *(*[]byte)(unsafe.Pointer(&str))
	err := ioutil.WriteFile(f.Path, b, 0644)
	if err != nil {
		return err
	}
	return nil
}

func readFile(path string) ([]byte, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return []byte(""), err
	}
	return b, nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func getFileStruct(r *http.Request) *File {
	root := "webmock-cache/"
	url := r.URL.Host + r.URL.Path
	name := "cache.json"
	dir := root + url
	path := dir + "/" + name
	return &File{url, root, path, dir, name}
}

func mkdir(dir string) error {
	return os.MkdirAll(dir, os.ModePerm)
}