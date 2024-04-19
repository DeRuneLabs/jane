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

const (
	os_windows = "windows"
	os_darwin  = "darwin"
	os_linux   = "linux"
)

const (
	arch_i386  = "386"
	arch_amd64 = "amd64"
	arch_arm   = "arm"
	arch_arm64 = "arm64"
)

func checkPlatform(path string) (ok bool, exist bool) {
	ok = false
	exist = true
	switch path {
	case jn.PlatformWindows:
		ok = runtime.GOOS == os_windows
	case jn.PlatformDarwin:
		ok = runtime.GOOS == os_darwin
	case jn.PlatformLinux:
		ok = runtime.GOOS == os_linux
	case jn.PlatformUnix:
		switch runtime.GOOS {
		case os_darwin, os_linux:
			ok = true
		}
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
		ok = runtime.GOARCH == arch_i386
	case jn.ArchAmd64:
		ok = runtime.GOARCH == arch_amd64
	case jn.ArchArm:
		ok = runtime.GOARCH == arch_arm
	case jn.ArchArm64:
		ok = runtime.GOARCH == arch_arm64
	case jn.Arch64Bit:
		switch runtime.GOARCH {
		case arch_amd64, arch_arm64:
			ok = true
		}
	case jn.Arch32Bit:
		switch runtime.GOARCH {
		case arch_i386, arch_arm:
			ok = true
		}
	default:
		ok = true
		exist = false
	}
	return
}

// return true if file path is pass file annotation,
// return false if not
func IsPassFileAnnotation(p string) bool {
	p = filepath.Base(p)
	n := len(p)
	p = p[:n-len(filepath.Ext(p))]

	a1 := ""
	a2 := ""

	i := strings.LastIndexByte(p, '_')
	if i == -1 {
		return true
	}
	if i+1 >= n {
		return true
	}
	a1 = p[i+1:]

	p = p[:i]
	i = strings.LastIndexByte(p, '_')
	if i != -1 {
		a2 = p[i+1:]
	}
	if a2 == "" {
		ok, exist := checkPlatform(a1)
		if exist && !ok {
			return false
		}
		ok, exist = checkArch(a1)
		return !exist || ok
	}

	ok, exist := checkArch(a1)
	if exist && !ok {
		return false
	}
	ok, exist = checkPlatform(a2)
	return !exist || ok
}
