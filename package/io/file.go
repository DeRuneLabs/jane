package io

import "os"

type FILE struct {
	Path    string
	Content []rune
}

func WriteFileTruncate(path string, content []byte) error {
	if _, err := os.Open(path); err == nil {
		os.Remove(path)
	}
	return os.WriteFile(path, content, 0x025E)
}
