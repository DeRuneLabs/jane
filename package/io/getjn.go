package io

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/De-Rune/jane/package/jn"
)

func GetJn(path string) (*FILE, error) {
	if filepath.Ext(path) != jn.Extension {
		return nil, errors.New(jn.Errors[`file_not_jn`] + path)
	}
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	f := new(FILE)
	f.Path = path
	f.Content = []rune(string(bytes))
	return f, nil
}