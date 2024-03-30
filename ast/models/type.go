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
	var cpp strings.Builder
	cpp.WriteString("typedef ")
	cpp.WriteString(t.Type.String())
	cpp.WriteByte(' ')
	cpp.WriteString(jnapi.OutId(t.Id, t.Tok.File))
	cpp.WriteByte(';')
	return cpp.String()
}
