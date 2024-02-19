package ast

import (
	"fmt"
	"github.com/De-Rune/jane/lexer"
	"strings"
)

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
	Content []StatementAST
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

type ExpressionAST struct {
	Content []ExpressionNode
	Type    uint8
}

func (e ExpressionAST) string() string {
	var sb strings.Builder
	for _, node := range e.Content {
		sb.WriteString(node.String() + " ")
	}
	return sb.String()[:sb.Len()-1]
}

type ExpressionNode struct {
	Content interface{}
	Type    uint8
}

func (n ExpressionNode) String() string {
	return fmt.Sprint(n.Content)
}

type ValueAST struct {
	Token lexer.Token
	Value string
	Type  uint8
}

func (v ValueAST) String() string {
	return v.Value
}

type BraceAST struct {
	Token lexer.Token
	Value string
}

func (b BraceAST) String() string {
	return b.Value
}

type OperatorAST struct {
	Token lexer.Token
	Value string
}

func (o OperatorAST) String() string {
	return o.Value
}

type ReturnAST struct {
	Token      lexer.Token
	Expression ExpressionAST
}

func (rast ReturnAST) String() string {
	if rast.Expression.Type != NA {
		return rast.Token.Value + " " + rast.Expression.string()
	}
	return rast.Token.Value
}
