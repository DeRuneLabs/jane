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
	"strconv"
	"strings"

	"github.com/DeRuneLabs/jane/lexer"
)

type Break struct {
	Token      lexer.Token
	LabelToken lexer.Token
	Label      string
}

func (b Break) String() string {
	return "goto " + b.Label + ";"
}

type Continue struct {
	Token     lexer.Token
	LoopLabel lexer.Token
	Label     string
}

func (c Continue) String() string {
	return "goto " + c.Label + ";"
}

type Iter struct {
	Token   lexer.Token
	Block   *Block
	Parent  *Block
	Profile IterProfile
}

func (i *Iter) BeginLabel() string {
	var cpp strings.Builder
	cpp.WriteString("iter_begin_")
	cpp.WriteString(strconv.Itoa(i.Token.Row))
	cpp.WriteString(strconv.Itoa(i.Token.Column))
	return cpp.String()
}

func (i *Iter) EndLabel() string {
	var cpp strings.Builder
	cpp.WriteString("iter_end_")
	cpp.WriteString(strconv.Itoa(i.Token.Row))
	cpp.WriteString(strconv.Itoa(i.Token.Column))
	return cpp.String()
}

func (i *Iter) NextLabel() string {
	var cpp strings.Builder
	cpp.WriteString("iter_next_")
	cpp.WriteString(strconv.Itoa(i.Token.Row))
	cpp.WriteString(strconv.Itoa(i.Token.Column))
	return cpp.String()
}

func (i *Iter) infinityString() string {
	var cpp strings.Builder
	indent := IndentString()
	begin := i.BeginLabel()
	cpp.WriteString(begin)
	cpp.WriteString(":;\n")
	cpp.WriteString(indent)
	cpp.WriteString(i.Block.String())
	cpp.WriteByte('\n')
	cpp.WriteString(indent)
	cpp.WriteString(i.NextLabel())
	cpp.WriteString(":;\n")
	cpp.WriteString(indent)
	cpp.WriteString("goto ")
	cpp.WriteString(begin)
	cpp.WriteString(";\n")
	cpp.WriteString(indent)
	cpp.WriteString(i.EndLabel())
	cpp.WriteString(":;")
	return cpp.String()
}

func (i Iter) String() string {
	if i.Profile == nil {
		return i.infinityString()
	}
	return i.Profile.String(&i)
}
