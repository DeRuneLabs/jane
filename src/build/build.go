// Copyright (c) 2024 arfy slowy - DeRuneLabs
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
// SOFTWARE

package build

import (
	"path/filepath"
	"runtime"
	"strings"

	"github.com/DeRuneLabs/jane"
)

const (
	ATTR_CDEF    = "cdef"
	ATTR_TYPEDEF = "typedef"
)

var ATTRS = [...]string{
	ATTR_CDEF,
	ATTR_TYPEDEF,
}

func check_os(arg string) (ok bool, exist bool) {
	ok = false
	exist = true
	switch arg {
	case OS_WINDOWS:
		ok = IsWindows(runtime.GOOS)
	case OS_DARWIN:
		ok = IsDarwin(runtime.GOOS)
	case OS_LINUX:
		ok = IsLinux(runtime.GOOS)
	case OS_UNIX:
		ok = IsUnix(runtime.GOOS)
	default:
		ok = true
		exist = false
	}
	return
}

func check_arch(arg string) (ok bool, exist bool) {
	ok = false
	exist = false
	switch arg {
	case ARCH_I386:
		ok = IsI386(runtime.GOARCH)
	case ARCH_AMD64:
		ok = IsAmd64(runtime.GOARCH)
	case ARCH_ARM:
		ok = IsArm(runtime.GOARCH)
	case ARCH_ARM64:
		ok = IsArm64(runtime.GOARCH)
	case ARCH_64Bit:
		ok = IsX64(runtime.GOARCH)
	case ARCH_32Bit:
		ok = IsX32(runtime.GOARCH)
	default:
		ok = true
		exist = false
	}
	return
}

func IsPassFileAnnotation(p string) bool {
	p = filepath.Base(p)
	n := len(p)
	p = p[:n-len(filepath.Ext(p))]
	a1 := ""
	a2 := ""

	i := strings.LastIndexByte(p, '_')
	if i == -1 {
		ok, exist := check_os(p)
		if exist {
			return ok
		}
		ok, exist = check_arch(p)
		return !exist || ok
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
		ok, exist := check_os(a1)
		if exist {
			return ok
		}
		ok, exist = check_arch(a1)
		return !exist || ok
	}

	ok, exist := check_arch(a1)
	if exist {
		if !ok {
			return false
		}
		ok, exist = check_os(a2)
		return !exist || ok
	}
	ok, exist = check_os(a1)
	return !exist || ok
}

func IsStdHeaderPath(p string) bool {
	return p[0] == '<' && p[len(p)-1] == '>'
}

var CPP_HEADER_EXTS = []string{
	".h", ".hpp", ".hxx", ".hh",
}

func IsValidHeader(ext string) bool {
	for _, validExt := range CPP_HEADER_EXTS {
		if ext == validExt {
			return true
		}
	}
	return false
}

func IsJane(path string) bool {
	abs, err := filepath.Abs(path)
	return err == nil && filepath.Ext(abs) == jane.EXT
}
