package models

import (
	"strings"

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
		tok := as.Expr.Toks[0]
		return jnapi.OutId(tok.Kind, tok.File)
	case as.Ignore:
		return jnapi.CxxIgnore
	}
	return as.Expr.String()
}

type Assign struct {
	Setter      Tok
	Left        []AssignLeft
	Right       []Expr
	IsExpr      bool
	MultipleRet bool
}

func (a *Assign) cxxSingleAssign() string {
	expr := a.Left[0]
	if expr.Var.New {
		expr.Var.Val = a.Right[0]
		s := expr.Var.String()
		return s[:len(s)-1]
	}
	var cxx strings.Builder
	if len(expr.Expr.Toks) != 1 ||
		!jnapi.IsIgnoreId(expr.Expr.Toks[0].Kind) {
		cxx.WriteString(expr.String())
		cxx.WriteString(a.Setter.Kind)
	}
	cxx.WriteString(a.Right[0].String())
	return cxx.String()
}

func (a *Assign) hasSelector() bool {
	for _, s := range a.Left {
		if !s.Ignore {
			return true
		}
	}
	return false
}

func (a *Assign) cxxMultipleAssign() string {
	var cxx strings.Builder
	if !a.hasSelector() {
		for _, expr := range a.Right {
			cxx.WriteString(expr.String())
			cxx.WriteByte(';')
		}
		return cxx.String()[:cxx.Len()-1]
	}
	cxx.WriteString(a.cxxNewDefines())
	cxx.WriteString("std::tie(")
	var expCxx strings.Builder
	expCxx.WriteString("std::make_tuple(")
	for i, selector := range a.Left {
		cxx.WriteString(selector.String())
		cxx.WriteByte(',')
		expCxx.WriteString(a.Right[i].String())
		expCxx.WriteByte(',')
	}
	str := cxx.String()[:cxx.Len()-1] + ")"
	cxx.Reset()
	cxx.WriteString(str)
	cxx.WriteString(a.Setter.Kind)
	cxx.WriteString(expCxx.String()[:expCxx.Len()-1] + ")")
	return cxx.String()
}

func (a *Assign) cxxMultipleReturn() string {
	var cxx strings.Builder
	cxx.WriteString(a.cxxNewDefines())
	cxx.WriteString("std::tie(")
	for _, selector := range a.Left {
		if selector.Ignore {
			cxx.WriteString(jnapi.CxxIgnore)
			cxx.WriteByte(',')
			continue
		}
		cxx.WriteString(selector.String())
		cxx.WriteByte(',')
	}
	str := cxx.String()[:cxx.Len()-1]
	cxx.Reset()
	cxx.WriteString(str)
	cxx.WriteByte(')')
	cxx.WriteString(a.Setter.Kind)
	cxx.WriteString(a.Right[0].String())
	return cxx.String()
}

func (a *Assign) cxxNewDefines() string {
	var cxx strings.Builder
	for _, selector := range a.Left {
		if selector.Ignore || !selector.Var.New {
			continue
		}
		cxx.WriteString(selector.Var.String() + " ")
	}
	return cxx.String()
}

func (a Assign) String() string {
	var cxx strings.Builder
	switch {
	case a.MultipleRet:
		cxx.WriteString(a.cxxMultipleReturn())
	case len(a.Left) == 1:
		cxx.WriteString(a.cxxSingleAssign())
	default:
		cxx.WriteString(a.cxxMultipleAssign())
	}
	if !a.IsExpr {
		cxx.WriteByte(';')
	}
	return cxx.String()
}
