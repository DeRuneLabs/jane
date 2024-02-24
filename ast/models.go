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

func (b BlockAST) String() string {
	return parseBlock(b, 1)
}

func parseBlock(b BlockAST, indent int) string {
	const indentSpace = 2
	var cxx strings.Builder
	for _, s := range b.Content {
		cxx.WriteByte('\n')
		cxx.WriteString(fullString(' ', indent*indentSpace))
		cxx.WriteString(fmt.Sprint(s.Value))
	}
	return cxx.String()
}

func fullString(b byte, count int) string {
	var sb strings.Builder
	for count > 0 {
		count--
		sb.WriteByte(b)
	}
	return sb.String()
}

type TypeAST struct {
	Token lexer.Token
	Code  uint8
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
	return jane.CxxTypeNameFromType(p.Type.Code) + " " + p.Name
}

type FunctionCallAST struct {
	Token lexer.Token
	Name  string
	Args  []ArgAST
}

func (fc FunctionCallAST) String() string {
	var cxx string
	cxx += fc.Name
	cxx += "("
	if len(fc.Args) > 0 {
		for _, arg := range fc.Args {
			cxx += arg.String()
			cxx += ","
		}
		cxx = cxx[:len(cxx)-1]
	}
	cxx += ");"
	return cxx
}

type ArgAST struct {
	Token      lexer.Token
	Tokens     []lexer.Token
	Expression ExpressionAST
}

func (a ArgAST) String() string {
	return a.Expression.String()
}

type ExpressionAST struct {
	Tokens    []lexer.Token
	Processes [][]lexer.Token
}

func (e ExpressionAST) String() string {
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
	return r.Token.Value + " " + r.Expression.String()
}

type AttributeAST struct {
	Token lexer.Token
	Value string
}

func (t AttributeAST) String() string {
	return t.Value
}

type VariableAST struct {
	Token lexer.Token
	Name  string
	Type  TypeAST
	Value ExpressionAST
}

func (v VariableAST) String() string {
	var sb strings.Builder
	if v.Type.Code == jane.Void {
		sb.WriteString("auto")
	} else {
		sb.WriteString(jane.CxxTypeNameFromType(v.Type.Code))
	}
	sb.WriteByte(' ')
	sb.WriteString(v.Name)
	sb.WriteString(" = ")
	sb.WriteString(v.Value.String())
	sb.WriteByte(';')
	return sb.String()
}
