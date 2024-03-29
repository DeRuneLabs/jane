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
	var cxx strings.Builder
	cxx.WriteString(be.Expr.String())
	cxx.WriteByte(';')
	return cxx.String()
}
