package lexer

import "github.com/De-Rune/jane/package/io"

type Token struct {
	File   *io.FILE
	Line   int
	Column int
	Value  string
	Type   uint
}

const (
	NA        = 0
	Type      = 1
	Name      = 2
	Brace     = 3
	Return    = 4
	SemiColon = 5
)
