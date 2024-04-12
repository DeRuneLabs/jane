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

package models

import (
	"strings"
	"sync/atomic"

	"github.com/DeRuneLabs/jane/package/jn"
)

type Block struct {
	Parent   *Block
	SubIndex int
	Tree     []Statement
	Gotos    *Gotos
	Labels   *Labels
	Func     *Func
}

func (b Block) String() string {
	AddIndent()
	defer func() { DoneIndent() }()
	return ParseBlock(b)
}

func ParseBlock(b Block) string {
	var cpp strings.Builder
	cpp.WriteByte('{')
	for _, s := range b.Tree {
		if s.Data == nil {
			continue
		}
		cpp.WriteByte('\n')
		cpp.WriteString(IndentString())
		cpp.WriteString(s.String())
	}
	cpp.WriteByte('\n')
	indent := strings.Repeat(jn.Set.Indent, int(Indent-1)*jn.Set.IndentCount)
	cpp.WriteString(indent)
	cpp.WriteByte('}')
	return cpp.String()
}

var Indent uint32 = 0

func IndentString() string {
	return strings.Repeat(jn.Set.Indent, int(Indent)*jn.Set.IndentCount)
}

func AddIndent() {
	atomic.AddUint32(&Indent, 1)
}

func DoneIndent() {
	atomic.SwapUint32(&Indent, Indent-1)
}
