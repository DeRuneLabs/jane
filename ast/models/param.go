package models

import (
	"strings"

	"github.com/DeRuneLabs/jane/lexer/tokens"
	"github.com/DeRuneLabs/jane/package/jn"
	"github.com/DeRuneLabs/jane/package/jnapi"
)

type Param struct {
	Tok       Tok
	Id        string
	Const     bool
	Volatile  bool
	Variadic  bool
	Reference bool
	Type      DataType
	Default   Expr
}

func (p *Param) TypeString() string {
	var ts strings.Builder
	if p.Variadic {
		ts.WriteString(tokens.TRIPLE_DOT)
	}
	if p.Reference {
		ts.WriteString(tokens.AMPER)
	}
	ts.WriteString(p.Type.Kind)
	return ts.String()
}

func (p *Param) OutId() string {
	return jnapi.OutId(p.Id, p.Tok.File)
}

func (p Param) String() string {
	var cxx strings.Builder
	cxx.WriteString(p.Prototype())
	if p.Id != "" && !jnapi.IsIgnoreId(p.Id) && p.Id != jn.Anonymous {
		cxx.WriteByte(' ')
		cxx.WriteString(p.OutId())
	}
	return cxx.String()
}

func (p *Param) Prototype() string {
	var cxx strings.Builder
	if p.Variadic {
		cxx.WriteString("slice<")
		cxx.WriteString(p.Type.String())
		cxx.WriteByte('>')
	} else {
		cxx.WriteString(p.Type.String())
	}
	if p.Reference {
		cxx.WriteByte('&')
	}
	return cxx.String()
}
