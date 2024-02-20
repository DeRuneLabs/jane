package parser

import (
	"github.com/De-Rune/jane/lexer"
)

type variable struct {
	Name  string
	Token lexer.Token
	Type  uint8
}
