package models

import (
	"fmt"
	"strings"
)

type Statement struct {
	Tok            Tok
	Data           any
	WithTerminator bool
}

func (s Statement) String() string {
	return fmt.Sprint(s.Data)
}

type ExprStatement struct {
	Expr Expr
}

func (be ExprStatement) String() string {
	var cpp strings.Builder
	cpp.WriteString(be.Expr.String())
	cpp.WriteByte(';')
	return cpp.String()
}
