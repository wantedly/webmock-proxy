package webmock

import (
	"fmt"
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

func readFilePaths(dir string) ([]string, error) {
	var files []string
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	target := filepath.Join(wd, dir)
	err = filepath.Walk(target, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(target, path)
		if err != nil {
			return err
		}
		files = append(files, rel)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}
