package controllers_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/wantedly/webmock-proxy/example/go/api/middleware"
	"github.com/wantedly/webmock-proxy/example/go/api/models"
	"github.com/wantedly/webmock-proxy/example/go/api/router"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func initialize() (*gin.Engine, error) {
	outDir, err := ioutil.TempDir("", "test")
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open("sqlite3", outDir+"/database.db")
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(
		&models.User{},
	)

	r := gin.Default()
	r.Use(middleware.SetDBtoContext(db))
	router.Initialize(r)
	return r, nil
}

func ioReader(i io.ReadCloser) ([]byte, error) {
	defer i.Close()
	b, err := ioutil.ReadAll(i)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func httpRequest(t *testing.T, method, url, body string) *http.Response {
	req, err := http.NewRequest(
		method,
		url,
		bytes.NewBuffer([]byte(body)),
	)

	if err != nil {
		t.Errorf("Faild to create %v request: %v", method, err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "*/*")

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Faild to http connection: %v", err)
	}
	return resp
}

func TestUser(t *testing.T) {
	s, err := initialize()
	if err != nil {
		t.Errorf("Faild to initialize server: %v", err)
	}
	ts := httptest.NewServer(s)
	defer ts.Close()

	// POST: /users
	body := `
{"name":"wantedly","age":4,"email":"wantedly@wantedly.com"}
`
	resp := httpRequest(t, "POST", ts.URL+"/users", body)
	b, err := ioReader(resp.Body)
	if err != nil {
		t.Errorf("Faild to read http response body: %v", err)
	}

	expect := `{"id":1,"name":"wantedly","age":4,"email":"wantedly@wantedly.com"}
`
	actual := string(b)
	if expect != actual {
		t.Fatalf("Incorrect response body. expeted: %s, actual: %s", expect, actual)
	}

	// GET /users
	resp = httpRequest(t, "GET", ts.URL+"/users", "")
	b, err = ioReader(resp.Body)
	if err != nil {
		t.Errorf("Faild to read http response body: %v", err)
	}

	expect = `[{"age":4,"email":"wantedly@wantedly.com","id":1,"name":"wantedly"}]
`
	actual = string(b)
	if expect != actual {
		t.Fatalf("Incorrect response body. expeted: %s, actual: %s", expect, actual)
	}

	// GET /users/:id
	resp = httpRequest(t, "GET", ts.URL+"/users/1", "")
	b, err = ioReader(resp.Body)
	if err != nil {
		t.Errorf("Faild to read http response body: %v", err)
	}

	expect = `{"age":4,"email":"wantedly@wantedly.com","id":1,"name":"wantedly"}
`
	actual = string(b)
	if expect != actual {
		t.Fatalf("Incorrect response body. expeted: %s, actual: %s", expect, actual)
	}

	// PUT /users/:id
	body = `
{"name":"sync","age":4,"email":"sync@wantedly.com"}
`
	resp = httpRequest(t, "PUT", ts.URL+"/users/1", body)
	b, err = ioReader(resp.Body)
	if err != nil {
		t.Errorf("Faild to read http response body: %v", err)
	}

	expect = `{"id":1,"name":"sync","age":4,"email":"sync@wantedly.com"}
`
	actual = string(b)
	if expect != actual {
		t.Fatalf("Incorrect response body. expeted: %s, actual: %s", expect, actual)
	}

	// DELETE /users/:id
	resp = httpRequest(t, "DELETE", ts.URL+"/users/1", "")
	b, err = ioReader(resp.Body)
	if err != nil {
		t.Errorf("Faild to read http response body: %v", err)
	}

	expect = ""
	actual = string(b)
	if expect != actual {
		t.Fatalf("Incorrect response body. expeted: %s, actual: %s", expect, actual)
	}
}
