package webmock

import (
	"io"
	"io/ioutil"
)

func ioReader(io io.ReadCloser) (string, error) {
	defer io.Close()
	body, err := ioutil.ReadAll(io)
	if err != nil {
		return string(body), err
	}
	return string(body), nil
}
