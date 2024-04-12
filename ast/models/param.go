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

	"github.com/DeRuneLabs/jane/lexer/tokens"
	"github.com/DeRuneLabs/jane/package/jn"
	"github.com/DeRuneLabs/jane/package/jnapi"
)

type Param struct {
	Tok       Tok
	Id        string
	Const     bool
	Volatile  bool
	Variadic  bool
	Reference bool
	Type      DataType
	Default   Expr
}

func (p *Param) TypeString() string {
	var ts strings.Builder
	if p.Variadic {
		ts.WriteString(tokens.TRIPLE_DOT)
	}
	if p.Reference {
		ts.WriteString(tokens.AMPER)
	}
	ts.WriteString(p.Type.Kind)
	return ts.String()
}

func (p *Param) OutId() string {
	return jnapi.AsId(p.Id)
}

func (p Param) String() string {
	var cpp strings.Builder
	cpp.WriteString(p.Prototype())
	if p.Id != "" && !jnapi.IsIgnoreId(p.Id) && p.Id != jn.Anonymous {
		cpp.WriteByte(' ')
		cpp.WriteString(p.OutId())
	}
	return cpp.String()
}

func (p *Param) Prototype() string {
	var cpp strings.Builder
	if p.Variadic {
		cpp.WriteString("slice<")
		cpp.WriteString(p.Type.String())
		cpp.WriteByte('>')
	} else {
		cpp.WriteString(p.Type.String())
	}
	if p.Reference {
		cpp.WriteByte('&')
	}
	return cpp.String()
}
