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
	Val       Expr
	Const     bool
	Volatile  bool
	New       bool
	Tag       any
	Desc      string
	Used      bool
}

func (v Var) String() string {
	var cxx strings.Builder
	if v.Volatile {
		cxx.WriteString("volatile ")
	}
	if v.Const {
		cxx.WriteString("const ")
	}
	cxx.WriteString(v.Type.String())
	cxx.WriteByte(' ')
	cxx.WriteString(jnapi.OutId(v.Id, v.IdTok.File))
	cxx.WriteByte('{')
	if v.Val.Processes != nil {
		cxx.WriteString(v.Val.String())
	}
	cxx.WriteByte('}')
	cxx.WriteByte(';')
	return cxx.String()
}

func (v *Var) FieldString() string {
	var cxx strings.Builder
	if v.Volatile {
		cxx.WriteString("volatile ")
	}
	if v.Const {
		cxx.WriteString("const ")
	}
	cxx.WriteString(v.Type.String())
	cxx.WriteByte(' ')
	cxx.WriteString(jnapi.OutId(v.Id, v.IdTok.File))
	cxx.WriteByte(';')
	return cxx.String()
}
