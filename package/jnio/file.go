package jnio

import (
	"path/filepath"
	"runtime"
	"strings"

	"github.com/DeRuneLabs/jane/package/jn"
)

type File struct {
	Path string
	Data []rune
}

func checkPlatform(path string) (ok bool, exist bool) {
	ok = false
	exist = true
	switch path {
	case jn.PlatformWindows:
		ok = runtime.GOOS == "windows"
	case jn.PlatformDarwin:
		ok = runtime.GOOS == "darwin"
	case jn.PlatformLinux:
		ok = runtime.GOOS == "linux"
	default:
		ok = true
		exist = false
	}
	return
}

func checkArch(path string) (ok bool, exist bool) {
	ok = false
	exist = true
	switch path {
	case jn.ArchI386:
		ok = runtime.GOARCH == "386"
	case jn.ArchAmd64:
		ok = runtime.GOARCH == "amd64"
	case jn.ArchArm:
		ok = runtime.GOARCH == "arm"
	case jn.ArchArm64:
		ok = runtime.GOARCH == "arm64"
	default:
		ok = true
		exist = false
	}
	return
}

func IsUseable(path string) bool {
	path = filepath.Base(path)
	path = path[:len(path)-len(filepath.Ext(path))]
	index := strings.LastIndexByte(path, '_')
	if index == -1 {
		return true
	}
	path = path[index+1:]
	ok, exist := checkPlatform(path)
	if exist {
		return ok
	}
	ok, _ = checkArch(path)
	return ok
}
