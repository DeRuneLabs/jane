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
	"github.com/DeRuneLabs/jane/package/jntype"
)

type foreach_setter interface {
	setup_vars(key_a, key_b Var) string
	next_steps(ket_a, key_b Var, begin string) string
}

type index_setter struct{}

func (index_setter) setup_vars(key_a, key_b Var) string {
	var cpp strings.Builder
	indent := IndentString()
	if !jnapi.IsIgnoreId(key_a.Id) {
		if key_a.New {
			cpp.WriteString(key_a.String())
			cpp.WriteByte(' ')
		}
		cpp.WriteString(key_a.OutId())
		cpp.WriteString(" = 0;\n")
		cpp.WriteString(indent)
	}
	if !jnapi.IsIgnoreId(key_b.Id) {
		if key_b.New {
			cpp.WriteString(key_b.String())
			cpp.WriteByte(' ')
		}
		cpp.WriteString(key_b.OutId())
		cpp.WriteString(" = *__jnc_foreach_begin;\n")
		cpp.WriteString(indent)
	}
	return cpp.String()
}

func (index_setter) next_steps(key_a, key_b Var, begin string) string {
	var cpp strings.Builder
	indent := IndentString()
	cpp.WriteString("++__jnc_foreach_begin;\n")
	cpp.WriteString(indent)
	cpp.WriteString("if (__jnc_foreach_begin != __jnc_foreach_end) { ")
	if !jnapi.IsIgnoreId(key_a.Id) {
		cpp.WriteString("++")
		cpp.WriteString(key_a.OutId())
		cpp.WriteString("; ")
	}
	if !jnapi.IsIgnoreId(key_b.Id) {
		cpp.WriteString(key_b.OutId())
		cpp.WriteString(" = *__jnc_foreach_begin; ")
	}
	cpp.WriteString("goto ")
	cpp.WriteString(begin)
	cpp.WriteString("; }\n")
	return cpp.String()
}

type map_setter struct{}

func (map_setter) setup_vars(key_a, key_b Var) string {
	var cpp strings.Builder
	indent := IndentString()
	if !jnapi.IsIgnoreId(key_a.Id) {
		if key_a.New {
			cpp.WriteString(key_a.String())
			cpp.WriteByte(' ')
		}
		cpp.WriteString(key_a.OutId())
		cpp.WriteString(" = __jnc_foreach_begin->first;\n")
		cpp.WriteString(indent)
	}
	if !jnapi.IsIgnoreId(key_b.Id) {
		if key_b.New {
			cpp.WriteString(key_b.String())
			cpp.WriteByte(' ')
		}
		cpp.WriteString(key_b.OutId())
		cpp.WriteString(" = __jnc_foreach_begin->second;\n")
		cpp.WriteString(indent)
	}
	return cpp.String()
}

func (map_setter) next_steps(key_a, key_b Var, begin string) string {
	var cpp strings.Builder
	indent := IndentString()
	cpp.WriteString("++__jnc_foreach_begin;\n")
	cpp.WriteString(indent)
	cpp.WriteString("if (__jnc_foreach_begin != __jnc_foreach_end) { ")
	if !jnapi.IsIgnoreId(key_a.Id) {
		cpp.WriteString(key_a.OutId())
		cpp.WriteString(" = __jnc_foreach_begin->first; ")
	}
	if !jnapi.IsIgnoreId(key_b.Id) {
		cpp.WriteString(key_b.OutId())
		cpp.WriteString(" = __jnc_foreach_begin->second; ")
	}
	cpp.WriteString("goto ")
	cpp.WriteString(begin)
	cpp.WriteString("; }\n")
	return cpp.String()
}

type IterForeach struct {
	KeyA     Var
	KeyB     Var
	InToken  lexer.Token
	Expr     Expr
	ExprType Type
}

func (f IterForeach) String(i *Iter) string {
	switch f.ExprType.Id {
	case jntype.Str, jntype.Slice, jntype.Array:
		return f.IterationString(i, index_setter{})
	case jntype.Map:
		return f.IterationString(i, map_setter{})
	}
	return ""
}

func (f *IterForeach) IterationString(i *Iter, setter foreach_setter) string {
	var cpp strings.Builder
	cpp.WriteString("{\n")
	AddIndent()
	indent := IndentString()
	cpp.WriteString(indent)
	cpp.WriteString("auto __jnc_foreach_expr = ")
	cpp.WriteString(f.Expr.String())
	cpp.WriteString(";\n")
	cpp.WriteString(indent)
	cpp.WriteString("if (__jnc_foreach_expr.begin() != __jnc_foreach_expr.end()) {\n")
	AddIndent()
	indent = IndentString()
	cpp.WriteString(indent)
	cpp.WriteString("auto __jnc_foreach_begin = __jnc_foreach_expr.begin();\n")
	cpp.WriteString(indent)
	cpp.WriteString("const auto __jnc_foreach_end = __jnc_foreach_expr.end();\n")
	cpp.WriteString(indent)
	cpp.WriteString(setter.setup_vars(f.KeyA, f.KeyB))
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
	cpp.WriteString(setter.next_steps(f.KeyA, f.KeyB, begin))
	cpp.WriteString(indent)
	cpp.WriteString(i.EndLabel())
	cpp.WriteString(":;")
	cpp.WriteByte('\n')
	DoneIndent()
	cpp.WriteString(IndentString())
	cpp.WriteString("}\n")
	DoneIndent()
	cpp.WriteString(IndentString())
	cpp.WriteByte('}')
	return cpp.String()
}
