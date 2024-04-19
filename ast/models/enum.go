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

	"github.com/DeRuneLabs/jane/lexer"
	"github.com/DeRuneLabs/jane/package/jnapi"
)

type EnumItem struct {
	Token   lexer.Token
	Id      string
	Expr    Expr
	ExprTag any
}

func (ei EnumItem) String() string {
	var cpp strings.Builder
	cpp.WriteString(jnapi.OutId(ei.Id, ei.Token.File))
	cpp.WriteString(" = ")
	cpp.WriteString(ei.Expr.String())
	return cpp.String()
}

type Enum struct {
	Pub   bool
	Token lexer.Token
	Id    string
	Type  Type
	Items []*EnumItem
	Used  bool
	Desc  string
}

func (e *Enum) ItemById(id string) *EnumItem {
	for _, item := range e.Items {
		if item.Id == id {
			return item
		}
	}
	return nil
}

func (e Enum) String() string {
	var cpp strings.Builder
	cpp.WriteString("enum ")
	cpp.WriteString(jnapi.OutId(e.Id, e.Token.File))
	cpp.WriteByte(':')
	cpp.WriteString(e.Type.String())
	cpp.WriteString(" {\n")
	AddIndent()
	for _, item := range e.Items {
		cpp.WriteString(IndentString())
		cpp.WriteString(item.String())
		cpp.WriteString(",\n")
	}
	DoneIndent()
	cpp.WriteString("};")
	return cpp.String()
}
