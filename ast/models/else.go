package models

import "strings"

type Else struct {
	Tok   Tok
	Block Block
}

func (elseast Else) String() string {
	var cxx strings.Builder
	cxx.WriteString("else ")
	cxx.WriteString(elseast.Block.String())
	return cxx.String()
}
