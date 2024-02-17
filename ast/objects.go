package ast

import "github.com/De-Rune/jane/lexer"

type Object struct {
	Token lexer.Token
	Value interface{}
	Type  uint8
}

type IdentifierAST struct {
	Type  uint8
	Value string
}

type StatementAST struct {
	Token lexer.Token
	Type  uint8
	Value interface{}
}

type RangeAST struct {
	Type    uint8
	Content []Object
}

type BlockAST struct {
	content []Object
}

type TypeAST struct {
	Type  uint8
	Value string
}

type FunctionAST struct {
	Token      lexer.Token
	Name       string
	ReturnType TypeAST
	Block      BlockAST
}
