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

package jnlog

import (
	"fmt"
	"strings"
)

const (
	FlatError   uint8 = 0
	FlatWarning uint8 = 1
	Error       uint8 = 2
	Warning     uint8 = 3
)

const warningMark = "<!>"

type CompilerLog struct {
	Type    uint8
	Row     int
	Column  int
	Path    string
	Message string
}

func (clog *CompilerLog) flatError() string {
	return clog.Message
}

func (clog *CompilerLog) error() string {
	var log strings.Builder
	log.WriteString(clog.Path)
	log.WriteByte(':')
	log.WriteString(fmt.Sprint(clog.Row))
	log.WriteByte(':')
	log.WriteString(fmt.Sprint(clog.Column))
	log.WriteByte(' ')
	log.WriteString(clog.Message)
	return log.String()
}

func (clog *CompilerLog) flatWarning() string {
	return warningMark + " " + clog.Message
}

func (clog *CompilerLog) warning() string {
	var log strings.Builder
	log.WriteString(warningMark)
	log.WriteByte(' ')
	log.WriteString(clog.Path)
	log.WriteByte(':')
	log.WriteString(fmt.Sprint(clog.Row))
	log.WriteByte(':')
	log.WriteString(fmt.Sprint(clog.Column))
	log.WriteByte(' ')
	log.WriteString(clog.Message)
	return log.String()
}

func (clog CompilerLog) String() string {
	switch clog.Type {
	case FlatError:
		return clog.flatError()
	case Error:
		return clog.error()
	case FlatWarning:
		return clog.flatWarning()
	case Warning:
		return clog.warning()
	}
	return ""
}
