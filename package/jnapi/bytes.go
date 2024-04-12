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
	"strconv"
	"strings"
	"unicode/utf8"
)

func ToStr(bytes []byte) string {
	var cpp strings.Builder
	cpp.WriteString("str_jnt{")
	btoa := bytesToStr(bytes)
	if btoa != "" {
		cpp.WriteByte('{')
		cpp.WriteString(btoa)
		cpp.WriteByte('}')
	}
	cpp.WriteString("}")
	return cpp.String()
}

func ToRawStr(bytes []byte) string {
	return ToStr(bytes)
}

func ToChar(b byte) string {
	return btoa(b)
}

func ToRune(bytes []byte) string {
	if len(bytes) == 0 {
		return ""
	} else if bytes[0] == '\\' {
		if len(bytes) > 1 && (bytes[1] == 'u' || bytes[1] == 'U') {
			bytes = bytes[2:]
			i, _ := strconv.ParseInt(string(bytes), 16, 32)
			return "0x" + strconv.FormatInt(i, 16)
		}
	}
	r, _ := utf8.DecodeRune(bytes)
	return "0x" + strconv.FormatInt(int64(r), 16)
}

func btoa(b byte) string {
	return "0x" + strconv.FormatUint(uint64(b), 16)
}

func bytesToStr(bytes []byte) string {
	if len(bytes) == 0 {
		return ""
	}
	var str strings.Builder
	for i := 0; i < len(bytes); i++ {
		b := bytes[i]
		if b == '\\' {
			i++
			switch bytes[i] {
			case 'u':
				rc, _ := strconv.ParseUint(string(bytes[i+1:i+5]), 16, 32)
				r := rune(rc)
				str.WriteString(bytesToStr([]byte(string(r))))
				i += 4
			case 'U':
				rc, _ := strconv.ParseUint(string(bytes[i+1:i+9]), 16, 32)
				r := rune(rc)
				str.WriteString(bytesToStr([]byte(string(r))))
				i += 8
			case 'x':
				str.WriteByte('0')
				str.Write(bytes[i : i+3])
				i += 2
			default:
				str.Write(bytes[i : i+3])
				i += 2
			}
		} else {
			str.WriteString(btoa(b))
		}
		str.WriteByte(',')
	}
	return str.String()[:str.Len()-1]
}
