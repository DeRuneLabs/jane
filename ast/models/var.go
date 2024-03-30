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
