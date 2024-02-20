package parser

import (
	"github.com/De-Rune/jane/lexer"
)

type Variable struct {
	Name  string
	Token lexer.Token
	Type  uint8
}
