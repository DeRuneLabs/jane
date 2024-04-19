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
	"github.com/DeRuneLabs/jane/lexer/tokens"
	"github.com/DeRuneLabs/jane/package/jnapi"
)

type Var struct {
	Owner     *Block
	Pub       bool
	Mutable   bool
	Token     lexer.Token
	SetterTok lexer.Token
	Id        string
	Type      Type
	Expr      Expr
	Const     bool
	New       bool
	Tag       any
	ExprTag   any
	Desc      string
	Used      bool
	IsField   bool
	CppLinked bool
}

func (v *Var) IsLocal() bool { return v.Owner != nil }

func as_local_id(row, column int, id string) string {
	id = strconv.Itoa(row) + strconv.Itoa(column) + "_" + id
	return jnapi.AsId(id)
}

func (v *Var) OutId() string {
	switch {
	case v.CppLinked:
		return v.Id
	case v.Id == tokens.SELF:
		return "self"
	case v.IsLocal():
		return as_local_id(v.Token.Row, v.Token.Column, v.Id)
	case v.IsField:
		return "__jnc_field_" + jnapi.AsId(v.Id)
	default:
		return jnapi.OutId(v.Id, v.Token.File)
	}
}

func (v Var) String() string {
	if jnapi.IsIgnoreId(v.Id) {
		return ""
	}
	if v.Const {
		return ""
	}
	var cpp strings.Builder
	cpp.WriteString(v.Type.String())
	cpp.WriteByte(' ')
	cpp.WriteString(v.OutId())
	expr := v.Expr.String()
	if expr != "" {
		cpp.WriteString(" = ")
		cpp.WriteString(v.Expr.String())
	} else {
		cpp.WriteString(jnapi.DefaultExpr)
	}
	cpp.WriteByte(';')
	return cpp.String()
}

func (v *Var) FieldString() string {
	var cpp strings.Builder
	if v.Const {
		cpp.WriteString("const ")
	}
	cpp.WriteString(v.Type.String())
	cpp.WriteByte(' ')
	cpp.WriteString(v.OutId())
	cpp.WriteString(jnapi.DefaultExpr)
	cpp.WriteByte(';')
	return cpp.String()
}

func (v *Var) ReceiverTypeString() string {
	var s strings.Builder
	if v.Mutable {
		s.WriteString("mut ")
	}
	if v.Type.Kind != "" && v.Type.Kind[0] == '&' {
		s.WriteByte('&')
	}
	s.WriteString("self")
	return s.String()
}
