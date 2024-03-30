package models

import "strings"

type Labels []*Label
type Gotos []*Goto

type Label struct {
	Tok   Tok
	Label string
	Index int
	Used  bool
	Block *Block
}

func (l Label) String() string {
	return l.Label + ":;"
}

type Goto struct {
	Tok   Tok
	Label string
	Index int
	Block *Block
}

func (gt Goto) String() string {
	var cpp strings.Builder
	cpp.WriteString("goto ")
	cpp.WriteString(gt.Label)
	cpp.WriteByte(';')
	return cpp.String()
}
