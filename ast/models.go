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
	Token      lexer.Token
	Name       string
	Expression ExpressionAST
	Args       []lexer.Token
}

func (fc FunctionCallAST) String() string {
	return fc.Expression.string()
}

type ExpressionAST struct {
	Tokens    []lexer.Token
	Processes [][]lexer.Token
}

func (e ExpressionAST) string() string {
	var sb strings.Builder
	for _, process := range e.Processes {
		if len(process) == 1 && process[0].Type == lexer.Operator {
			sb.WriteByte(' ')
			sb.WriteString(process[0].Value)
			sb.WriteByte(' ')
			continue
		}
		for _, token := range process {
			sb.WriteString(token.Value)
		}
	}
	return sb.String()
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

type AttributeAST struct {
	Token lexer.Token
	Value string
}

func (t AttributeAST) String() string {
	return t.Value
}
