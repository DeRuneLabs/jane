package models

import (
	"strings"

	"github.com/DeRuneLabs/jane/package/jnapi"
)

type EnumItem struct {
	Tok  Tok
	Id   string
	Expr Expr
}

func (ei EnumItem) String() string {
	var cxx strings.Builder
	cxx.WriteString(jnapi.OutId(ei.Id, ei.Tok.File))
	cxx.WriteString(" = ")
	cxx.WriteString(ei.Expr.String())
	return cxx.String()
}

type Enum struct {
	Pub   bool
	Tok   Tok
	Id    string
	Type  DataType
	Items []*EnumItem
	Used  bool
	Desc  string
}

func (e *Enum) ItemById(id string) *EnumItem {
	for _, item := range e.Items {
		if item.Id == id {
			return item
		}
	}
	return nil
}

func (e Enum) String() string {
	var cxx strings.Builder
	cxx.WriteString("enum ")
	cxx.WriteString(jnapi.OutId(e.Id, e.Tok.File))
	cxx.WriteByte(':')
	cxx.WriteString(e.Type.String())
	cxx.WriteString(" {\n")
	AddIndent()
	for _, item := range e.Items {
		cxx.WriteString(IndentString())
		cxx.WriteString(item.String())
		cxx.WriteString(",\n")
	}
	DoneIndent()
	cxx.WriteString("};")
	return cxx.String()
}
