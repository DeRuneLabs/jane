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

package parser

import (
	"strings"

	"github.com/DeRuneLabs/jane/ast/models"
	"github.com/DeRuneLabs/jane/package/jn"
	"github.com/DeRuneLabs/jane/package/jnapi"
)

type function struct {
	Ast          *Func
	Desc         string
	used         bool
	checked      bool
	isEntryPoint bool
}

func (f *function) outId() string {
	if f.isEntryPoint {
		return jnapi.OutId(f.Ast.Id, nil)
	}
	return f.Ast.OutId()
}

func (f function) String() string {
	var cpp strings.Builder
	cpp.WriteString(f.Head())
	cpp.WriteByte(' ')
	block := f.Ast.Block
	vars := f.Ast.RetType.Vars()
	if vars != nil {
		statements := make([]models.Statement, len(vars))
		for i, v := range vars {
			statements[i] = models.Statement{Tok: v.Token, Data: *v}
		}
		block.Tree = append(statements, block.Tree...)
	}
	if f.Ast.Receiver != nil && !typeIsPtr(*f.Ast.Receiver) {
		s := f.Ast.Receiver.Tag.(*jnstruct)
		self := s.selfVar(*f.Ast.Receiver)
		statements := make([]models.Statement, 1)
		statements[0] = models.Statement{Tok: s.Ast.Tok, Data: self}
		block.Tree = append(statements, block.Tree...)
	}
	cpp.WriteString(block.String())
	return cpp.String()
}

func (f *function) Head() string {
	var cpp strings.Builder
	cpp.WriteString(f.declHead())
	cpp.WriteString(paramsToCpp(f.Ast.Params))
	return cpp.String()
}

func (f *function) declHead() string {
	var cpp strings.Builder
	cpp.WriteString(genericsToCpp(f.Ast.Generics))
	if cpp.Len() > 0 {
		cpp.WriteByte('\n')
		cpp.WriteString(models.IndentString())
	}
	cpp.WriteString(attributesToString(f.Ast.Attributes))
	cpp.WriteString(f.Ast.RetType.String())
	cpp.WriteByte(' ')
	cpp.WriteString(f.outId())
	return cpp.String()
}

func (f *function) Prototype() string {
	var cpp strings.Builder
	cpp.WriteString(f.declHead())
	cpp.WriteString(f.PrototypeParams())
	cpp.WriteByte(';')
	return cpp.String()
}

func (f *function) PrototypeParams() string {
	if len(f.Ast.Params) == 0 {
		return "(void)"
	}
	var cpp strings.Builder
	cpp.WriteByte('(')
	for _, p := range f.Ast.Params {
		cpp.WriteString(p.Prototype())
		cpp.WriteByte(',')
	}
	return cpp.String()[:cpp.Len()-1] + ")"
}

func isOutableAttribute(kind string) bool {
	return kind == jn.Attribute_Inline
}

func attributesToString(attributes []models.Attribute) string {
	var cpp strings.Builder
	for _, attr := range attributes {
		if isOutableAttribute(attr.Tag) {
			cpp.WriteString(attr.String())
			cpp.WriteByte(' ')
		}
	}
	return cpp.String()
}

func paramsToCpp(params []Param) string {
	if len(params) == 0 {
		return "(void)"
	}
	var cpp strings.Builder
	cpp.WriteByte('(')
	for _, p := range params {
		cpp.WriteString(p.String())
		cpp.WriteByte(',')
	}
	return cpp.String()[:cpp.Len()-1] + ")"
}
