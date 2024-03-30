package models

import "strings"

type Ret struct {
	Tok  Tok
	Expr Expr
}

func (r Ret) String() string {
	var cpp strings.Builder
	cpp.WriteString("return ")
	cpp.WriteString(r.Expr.String())
	cpp.WriteByte(';')
	return cpp.String()
}
