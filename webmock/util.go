package webmock

import (
	"io"
	"io/ioutil"
)

func ioReader(io io.ReadCloser) ([]byte, error) {
	defer io.Close()
	b, err := ioutil.ReadAll(io)
	if err != nil {
		return b, err
	}
	return b, nil
}
