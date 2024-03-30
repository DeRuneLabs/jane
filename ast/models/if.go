package models

import "strings"

type If struct {
	Tok   Tok
	Expr  Expr
	Block *Block
}

func (ifast If) String() string {
	var cpp strings.Builder
	cpp.WriteString("if (")
	cpp.WriteString(ifast.Expr.String())
	cpp.WriteString(") ")
	cpp.WriteString(ifast.Block.String())
	return cpp.String()
}
