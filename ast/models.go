package ast

import (
	"fmt"
	"strings"

	"github.com/De-Rune/jane/lexer"
	"github.com/De-Rune/jane/package/jane"
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
	Token lexer.Token
	Type  uint8
	Value string
}

type FunctionAST struct {
	Token      lexer.Token
	Name       string
	Params     []ParameterAST
	ReturnType TypeAST
	Block      BlockAST
}

type ParameterAST struct {
	Token lexer.Token
	Name  string
	Type  TypeAST
}

func (p ParameterAST) String() string {
	return jane.CxxTypeNameFromType(p.Type.Type) + " " + p.Name
}

type FunctionCallAST struct {
	Token lexer.Token
	Name  string
	Args  []lexer.Token
}

func (fc FunctionCallAST) String() string {
	var sb strings.Builder
	sb.WriteString(fc.Name)
	sb.WriteByte('(')
	sb.WriteString(tokensToString(fc.Args))
	sb.WriteByte(')')
	return sb.String()
}

type ExpressionAST struct {
	Tokens    []lexer.Token
	Processes [][]lexer.Token
}

func (e ExpressionAST) string() string {
	return tokensToString(e.Tokens)
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

func (r ReturnAST) String() string {
	return r.Token.Value + " " + r.Expression.string()
}

func tokensToString(tokens []lexer.Token) string {
	var sb strings.Builder
	for _, token := range tokens {
		sb.WriteString(token.Value)
		sb.WriteByte(' ')
	}
	return sb.String()
}

type TagAST struct {
	Token lexer.Token
	Value string
}

func (t TagAST) String() string {
	return t.Value
}
