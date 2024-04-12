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

	"github.com/DeRuneLabs/jane/package/jnapi"
	"github.com/DeRuneLabs/jane/package/jntype"
)

type IterForeach struct {
	KeyA     Var
	KeyB     Var
	InTok    Tok
	Expr     Expr
	ExprType DataType
}

func (f IterForeach) String(iter Iter) string {
	if !jnapi.IsIgnoreId(f.KeyA.Id) {
		return f.ForeachString(iter)
	}
	return f.IterationString(iter)
}

func (f *IterForeach) ClassicString(iter Iter) string {
	var cpp strings.Builder
	cpp.WriteString("foreach<")
	cpp.WriteString(f.ExprType.String())
	cpp.WriteByte(',')
	cpp.WriteString(f.KeyA.Type.String())
	if !jnapi.IsIgnoreId(f.KeyB.Id) {
		cpp.WriteByte(',')
		cpp.WriteString(f.KeyB.Type.String())
	}
	cpp.WriteString(">(")
	cpp.WriteString(f.Expr.String())
	cpp.WriteString(", [&](")
	cpp.WriteString(f.KeyA.Type.String())
	cpp.WriteByte(' ')
	cpp.WriteString(f.KeyA.OutId())
	if !jnapi.IsIgnoreId(f.KeyB.Id) {
		cpp.WriteByte(',')
		cpp.WriteString(f.KeyB.Type.String())
		cpp.WriteByte(' ')
		cpp.WriteString(f.KeyB.OutId())
	}
	cpp.WriteString(") -> void ")
	cpp.WriteString(iter.Block.String())
	cpp.WriteString(");")
	return cpp.String()
}

func (f *IterForeach) MapString(iter Iter) string {
	var cpp strings.Builder
	cpp.WriteString("foreach<")
	types := f.ExprType.Tag.([]DataType)
	cpp.WriteString(types[0].String())
	cpp.WriteByte(',')
	cpp.WriteString(types[1].String())
	cpp.WriteString(">(")
	cpp.WriteString(f.Expr.String())
	cpp.WriteString(", [&](")
	cpp.WriteString(f.KeyA.Type.String())
	cpp.WriteByte(' ')
	cpp.WriteString(f.KeyA.OutId())
	if !jnapi.IsIgnoreId(f.KeyB.Id) {
		cpp.WriteByte(',')
		cpp.WriteString(f.KeyB.Type.String())
		cpp.WriteByte(' ')
		cpp.WriteString(f.KeyB.OutId())
	}
	cpp.WriteString(") -> void ")
	cpp.WriteString(iter.Block.String())
	cpp.WriteString(");")
	return cpp.String()
}

func (f *IterForeach) ForeachString(iter Iter) string {
	switch f.ExprType.Id {
	case jntype.Str, jntype.Slice, jntype.Array:
		return f.ClassicString(iter)
	case jntype.Map:
		return f.MapString(iter)
	}
	return ""
}

func (f IterForeach) IterationString(iter Iter) string {
	var cpp strings.Builder
	cpp.WriteString("for (auto ")
	cpp.WriteString(f.KeyB.OutId())
	cpp.WriteString(" : ")
	cpp.WriteString(f.Expr.String())
	cpp.WriteString(") ")
	cpp.WriteString(iter.Block.String())
	return cpp.String()
}
