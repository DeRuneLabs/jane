package models

import (
	"strings"

	"github.com/DeRuneLabs/jane/package/jnapi"
)

type Var struct {
	Pub       bool
	DefTok    Tok
	IdTok     Tok
	SetterTok Tok
	Id        string
	Type      DataType
	Expr      Expr
	Const     bool
	New       bool
	Tag       any
	ExprTag   any
	Desc      string
	Used      bool
}

func (v *Var) OutId() string {
	return jnapi.OutId(v.Id, v.IdTok.File)
}

func (v Var) String() string {
	if v.Const {
		return ""
	}
	var cxx strings.Builder
	cxx.WriteString(v.Type.String())
	cxx.WriteByte(' ')
	cxx.WriteString(v.OutId())
	expr := v.Expr.String()
	if expr != "" {
		cxx.WriteString(" = ")
		cxx.WriteString(v.Expr.String())
	} else {
		cxx.WriteString(jnapi.DefaultExpr)
	}
	cxx.WriteByte(';')
	return cxx.String()
}

func (v *Var) FieldString() string {
	var cxx strings.Builder
	if v.Const {
		cxx.WriteString("const ")
	}
	cxx.WriteString(v.Type.String())
	cxx.WriteByte(' ')
	cxx.WriteString(v.OutId())
	cxx.WriteString(jnapi.DefaultExpr)
	cxx.WriteByte(';')
	return cxx.String()
}
