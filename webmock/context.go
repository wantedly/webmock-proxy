package webmock

import (
	"io/ioutil"
	"net/http"
	"net/url"
)

type Context struct {
	Req *Req
}

type Req struct {
	Method string
	URL    *url.URL
	Proto  string
	Header http.Header
	Body   []byte
}

func NewReq(r *http.Request) (*Req, error) {
	req := new(Req)
	req.Method = r.Method
	req.URL = r.URL
	req.Proto = r.Proto
	req.Header = make(http.Header, len(r.Header))
	for k, vs := range r.Header {
		for _, v := range vs {
			req.Header.Add(k, v)
		}
	}
	var err error
	req.Body, err = ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	return req, nil
}
