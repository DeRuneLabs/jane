package models

import (
	"strings"

	"github.com/DeRuneLabs/jane/lexer/tokens"
)

type Try struct {
	Tok   Tok
	Block Block
	Catch Catch
}

func (t Try) String() string {
	var cxx strings.Builder
	cxx.WriteString("try ")
	cxx.WriteString(t.Block.String())
	if t.Catch.Tok.Id == tokens.NA {
		cxx.WriteString(" catch(...) {}")
	} else {
		cxx.WriteByte(' ')
		cxx.WriteString(t.Catch.String())
	}
	return cxx.String()
}
