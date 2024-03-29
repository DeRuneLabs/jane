package models

import (
	"strings"

	"github.com/DeRuneLabs/jane/lexer"
	"github.com/DeRuneLabs/jane/package/jnapi"
)

type (
	Tok  = lexer.Tok
	Toks = []Tok
)

type Type struct {
	Pub  bool
	Tok  Tok
	Id   string
	Type DataType
	Desc string
	Used bool
}

func (t Type) String() string {
	var cxx strings.Builder
	cxx.WriteString("typedef ")
	cxx.WriteString(t.Type.String())
	cxx.WriteByte(' ')
	cxx.WriteString(jnapi.OutId(t.Id, t.Tok.File))
	cxx.WriteByte(';')
	return cxx.String()
}
