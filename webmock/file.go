package webmock

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

func writeFile(str, path string) error {
	dir := filepath.Dir(path)
	if !fileExists(path) && !fileExists(dir) {
		err := mkdir(dir)
		if err != nil {
			return err
		}
	}
	b := []byte(str)
	err := ioutil.WriteFile(path, b, 0644)
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

func mkdir(dir string) error {
	return os.MkdirAll(dir, os.ModePerm)
}
