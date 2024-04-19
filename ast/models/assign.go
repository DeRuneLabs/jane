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

type AssignLeft struct {
	Var    Var
	Expr   Expr
	Ignore bool
}

func (as AssignLeft) String() string {
	switch {
	case as.Var.New:
		return as.Var.OutId()
	case as.Ignore:
		return jnapi.CppIgnore
	}
	return as.Expr.String()
}

type Assign struct {
	Setter      lexer.Token
	Left        []AssignLeft
	Right       []Expr
	IsExpr      bool
	MultipleRet bool
}

func (a *Assign) cppSingleAssign() string {
	expr := a.Left[0]
	if expr.Var.New {
		expr.Var.Expr = a.Right[0]
		s := expr.Var.String()
		return s[:len(s)-1]
	}
	var cpp strings.Builder
	if len(expr.Expr.Tokens) != 1 ||
		!jnapi.IsIgnoreId(expr.Expr.Tokens[0].Kind) {
		cpp.WriteString(expr.String())
		cpp.WriteString(a.Setter.Kind)
	}
	cpp.WriteString(a.Right[0].String())
	return cpp.String()
}

func (a *Assign) hasLeft() bool {
	for _, s := range a.Left {
		if !s.Ignore {
			return true
		}
	}
	return false
}

func (a *Assign) cppMultipleAssign() string {
	var cpp strings.Builder
	if !a.hasLeft() {
		for _, right := range a.Right {
			cpp.WriteString(right.String())
			cpp.WriteByte(';')
		}
		return cpp.String()[:cpp.Len()-1]
	}
	cpp.WriteString(a.cppNewDefines())
	cpp.WriteString("std::tie(")
	var exprCpp strings.Builder
	exprCpp.WriteString("std::make_tuple(")
	for i, left := range a.Left {
		cpp.WriteString(left.String())
		cpp.WriteByte(',')
		exprCpp.WriteString(a.Right[i].String())
		exprCpp.WriteByte(',')
	}
	str := cpp.String()[:cpp.Len()-1] + ")"
	cpp.Reset()
	cpp.WriteString(str)
	cpp.WriteString(a.Setter.Kind)
	cpp.WriteString(exprCpp.String()[:exprCpp.Len()-1] + ")")
	return cpp.String()
}

func (a *Assign) cppMultiRet() string {
	var cpp strings.Builder
	cpp.WriteString(a.cppNewDefines())
	cpp.WriteString("std::tie(")
	for _, left := range a.Left {
		if left.Ignore {
			cpp.WriteString(jnapi.CppIgnore)
			cpp.WriteByte(',')
			continue
		}
		cpp.WriteString(left.String())
		cpp.WriteByte(',')
	}
	str := cpp.String()[:cpp.Len()-1]
	cpp.Reset()
	cpp.WriteString(str)
	cpp.WriteByte(')')
	cpp.WriteString(a.Setter.Kind)
	cpp.WriteString(a.Right[0].String())
	return cpp.String()
}

func (a *Assign) cppNewDefines() string {
	var cpp strings.Builder
	for _, left := range a.Left {
		if left.Ignore || !left.Var.New {
			continue
		}
		cpp.WriteString(left.Var.String() + " ")
	}
	return cpp.String()
}

func (a *Assign) cppPostfix() string {
	var cpp strings.Builder
	cpp.WriteString(a.Left[0].Expr.String())
	cpp.WriteString(a.Setter.Kind)
	return cpp.String()
}

func (a Assign) String() string {
	var cpp strings.Builder
	switch {
	case len(a.Right) == 0:
		cpp.WriteString(a.cppPostfix())
	case a.MultipleRet:
		cpp.WriteString(a.cppMultiRet())
	case len(a.Left) == 1:
		cpp.WriteString(a.cppSingleAssign())
	default:
		cpp.WriteString(a.cppMultipleAssign())
	}
	if !a.IsExpr {
		cpp.WriteByte(';')
	}
	return cpp.String()
}
