// Copyright (c) 2024 - DeRuneLabs
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package jnio

import (
	"path/filepath"
	"runtime"
	"strings"

	"github.com/DeRuneLabs/jane/package/jn"
)

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
