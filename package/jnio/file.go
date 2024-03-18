package jnio

import "path/filepath"

type File struct {
	Dir  string
	Name string
	Data []rune
}

func (f *File) Path() string {
	return filepath.Join(f.Dir, f.Name)
}
