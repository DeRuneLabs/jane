package lexer

import "github.com/De-Rune/jane/package/io"

type Lexer struct {
	File     *io.FILE
	Position int
	Column   int
	Line     int
	Errors   []string
}

func New(f *io.FILE) *Lexer {
	lex := new(Lexer)
	lex.File = f
	lex.Line = 1
	lex.Column = 1
	lex.Position = 0
	return lex
}
