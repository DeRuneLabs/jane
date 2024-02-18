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
	lexer := new(Lexer)
	lexer.File = f
	lexer.Line = 1
	lexer.Column = 1
	lexer.Position = 1
	return lexer
}
