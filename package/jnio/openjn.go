package jnio

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/DeRuneLabs/jane/package/jn"
)

func OpenJn(path string) (*File, error) {
	if filepath.Ext(path) != jn.SrcExt {
		return nil, errors.New(jn.GetErr("file_not_jn", path))
	}
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	f := new(File)
	f.Path = path
	f.Data = []rune(string(bytes))
	return f, nil
}
