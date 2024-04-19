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

package jnapi

import (
	"fmt"
	"strings"
	"unsafe"

	"github.com/DeRuneLabs/jane/package/jnio"
)

// ignore operator
const Ignore = "_"

// initializer caller identifier
const InitializerCaller = "__jnc_call_package_initializers"

// type extension
const typeExtension = "_jnt"

// reporting identifier is ignore or not
func IsIgnoreId(id string) bool {
	return id == Ignore
}

// return specified identifier as jn Cpp identifier
func AsId(id string) string {
	// convert to "JNC_ID(" + id + ")"
	return "_" + id
}

func getPtrAsId(ptr unsafe.Pointer) string {
	address := fmt.Sprintf("%p", ptr)
	address = address[3:]
	for i, r := range address {
		if r != '0' {
			address = address[i:]
			break
		}
	}
	return address
}

func OutId(id string, f *jnio.File) string {
	if f != nil {
		var out strings.Builder
		out.WriteByte('f')
		out.WriteString(getPtrAsId(unsafe.Pointer(f)))
		out.WriteByte('_')
		out.WriteString(id)
		return out.String()
	}
	return AsId(id)
}

func AsTypeId(id string) string {
	return id + typeExtension
}
