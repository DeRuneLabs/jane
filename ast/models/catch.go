package models

import (
	"strings"

	"github.com/DeRuneLabs/jane/package/jnapi"
)

type Catch struct {
	Tok   Tok
	Var   Var
	Block Block
}

func (c Catch) String() string {
	var cxx strings.Builder
	cxx.WriteString("catch (")
	if c.Var.Id == "" {
		cxx.WriteString("...")
	} else {
		cxx.WriteString(c.Var.Type.String())
		cxx.WriteByte(' ')
		cxx.WriteString(jnapi.OutId(c.Var.Id, c.Tok.File))
	}
	cxx.WriteString(") ")
	cxx.WriteString(c.Block.String())
	return cxx.String()
}
