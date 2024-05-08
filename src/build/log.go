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
// Software

package build

import (
	"strconv"
	"strings"
)

const FLAT_ERR uint8 = 0
const ERR uint8 = 1

type Log struct {
	Type   uint8
	Row    int
	Column int
	Path   string
	Text   string
}

func (l *Log) flat_err() string {
	return l.Text
}

func (l *Log) err() string {
	var log strings.Builder
	log.WriteString(l.Path)
	log.WriteByte(':')
	log.WriteString(strconv.Itoa(l.Row))
	log.WriteByte(':')
	log.WriteString(strconv.Itoa(l.Column))
	log.WriteByte(' ')
	log.WriteString(l.Text)
	return log.String()
}

func (l Log) String() string {
	switch l.Type {
	case FLAT_ERR:
		return l.flat_err()
	case ERR:
		return l.err()
	}
	return ""
}
