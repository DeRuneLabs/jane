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
	var cxx strings.Builder
	cxx.WriteString("goto ")
	cxx.WriteString(gt.Label)
	cxx.WriteByte(';')
	return cxx.String()
}
