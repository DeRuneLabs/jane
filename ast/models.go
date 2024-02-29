package ast

import (
	"fmt"
	"strings"

	"github.com/De-Rune/jane/lexer"
	"github.com/De-Rune/jane/package/jn"
)

type Object struct {
	Token lexer.Token
	Value interface{}
}

type StatementAST struct {
	Token lexer.Token
	Value interface{}
}

func (s StatementAST) String() string {
	return fmt.Sprint(s.Value)
}

type RangeAST struct {
	Type    uint8
	Content []Object
}

type BlockAST struct {
	Statements []StatementAST
}

var Indent = 1

func (b BlockAST) String() string {
	Indent := 1
	return ParseBlock(b, Indent)
}

const IndentSpace = 2

func ParseBlock(b BlockAST, indent int) string {
	var cxx strings.Builder
	cxx.WriteByte('{')
	for _, s := range b.Statements {
		cxx.WriteByte('\n')
		cxx.WriteString(strings.Repeat(" ", indent*IndentSpace))
		cxx.WriteString(s.String())
	}
	cxx.WriteByte('\n')
	cxx.WriteString(strings.Repeat(" ", (indent-1)*IndentSpace) + "}")
	return cxx.String()
}

type DataTypeAST struct {
	Token      lexer.Token
	Code       uint8
	Value      string
	MultiTyped bool
	Tag        interface{}
}

func (dt DataTypeAST) String() string {
	var cxx strings.Builder
	for index, run := range dt.Value {
		if run == '*' {
			cxx.WriteRune(run)
			continue
		}
		dt.Value = dt.Value[index:]
		break
	}
	if dt.MultiTyped {
		return dt.MultiTypeString() + cxx.String()
	}
	if dt.Value != "" && dt.Value[0] == '[' {
		pointers := cxx.String()
		cxx.Reset()
		cxx.WriteString("array<")
		dt.Value = dt.Value[2:]
		cxx.WriteString(dt.String())
		cxx.WriteByte('>')
		cxx.WriteString(pointers)
		return cxx.String()
	}
	switch dt.Code {
	case jn.Name:
		return dt.Token.Kind + cxx.String()
	case jn.Function:
		return dt.FunctionString() + cxx.String()
	}
	return jn.CxxTypeNameFromType(dt.Code) + cxx.String()
}

func (dt DataTypeAST) FunctionString() string {
	cxx := "function<"
	fun := dt.Tag.(FunctionAST)
	cxx += fun.ReturnType.String()
	cxx += "("
	if len(fun.Params) > 0 {
		for _, param := range fun.Params {
			cxx += param.Type.String() + ", "
		}
		cxx = cxx[:len(cxx)-2]
	}
	cxx += ")>"
	return cxx
}

func (dt DataTypeAST) MultiTypeString() string {
	types := dt.Tag.([]DataTypeAST)
	var cxx strings.Builder
	cxx.WriteString("std::tuple<")
	for _, t := range types {
		cxx.WriteString(t.String())
		cxx.WriteByte(',')
	}
	return cxx.String()[:cxx.Len()-1] + ">"
}

type TypeAST struct {
	Token lexer.Token
	Name  string
	Type  DataTypeAST
}

func (t TypeAST) String() string {
	return "typedef " + t.Type.String() + " " + t.Name + ";"
}

type FunctionAST struct {
	Token      lexer.Token
	Name       string
	Params     []ParameterAST
	ReturnType DataTypeAST
	Block      BlockAST
}

type ParameterAST struct {
	Token lexer.Token
	Name  string
	Const bool
	Type  DataTypeAST
}

func (p ParameterAST) String() string {
	var cxx strings.Builder
	if p.Const {
		cxx.WriteString("const ")
	}
	cxx.WriteString(p.Type.String())
	if p.Name != "" {
		return cxx.String() + " " + p.Name
	}
	return cxx.String()
}

func (fc FunctionAST) DataTypeString() string {
	dt := "("
	if len(fc.Params) > 0 {
		for _, param := range fc.Params {
			dt += param.Type.Value + ", "
		}
		dt = dt[:len(dt)-2]
	}
	dt += ")"
	if fc.ReturnType.Code != jn.Void {
		dt += fc.ReturnType.Value
	}
	return dt
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
	Model     Expressionmodel
}

type Expressionmodel interface {
	String() string
}

func (e ExpressionAST) String() string {
	if e.Model != nil {
		return e.Model.String()
	}
	var sb strings.Builder
	for _, process := range e.Processes {
		if len(process) == 1 && process[0].Id == lexer.Operator {
			sb.WriteByte(' ')
			sb.WriteString(process[0].Kind)
			sb.WriteByte(' ')
			continue
		}
		for _, token := range process {
			sb.WriteString(token.Kind)
		}
	}
	return sb.String()
}

type BlockExpressionAST struct {
	Expression ExpressionAST
}

func (be BlockExpressionAST) String() string {
	return be.Expression.String() + ";"
}

type ValueAST struct {
	Token lexer.Token
	Value string
	Type  DataTypeAST
}

func (v ValueAST) String() string {
	return v.Value
}

type BraceAST struct {
	Token lexer.Token
}

func (b BraceAST) String() string {
	return b.Token.Kind
}

type OperatorAST struct {
	Token lexer.Token
}

func (o OperatorAST) String() string {
	return o.Token.Kind
}

type ReturnAST struct {
	Token      lexer.Token
	Expression ExpressionAST
}

func (r ReturnAST) String() string {
	switch r.Token.Id {
	case lexer.Operator:
		return "return " + r.Expression.String() + ";"
	}
	return "return " + r.Expression.String() + ";"
}

type AttributeAST struct {
	Token lexer.Token
	Tag   lexer.Token
}

func (a AttributeAST) String() string {
	return a.Tag.Kind[1:]
}

type VariableAST struct {
	DefineToken lexer.Token
	NameToken   lexer.Token
	SetterToken lexer.Token
	Name        string
	Type        DataTypeAST
	Value       ExpressionAST
	Tag         interface{}
}

func (v VariableAST) String() string {
	var sb strings.Builder
	switch v.DefineToken.Id {
	case lexer.Const:
		sb.WriteString("const ")
	}
	sb.WriteString(v.StringType())
	sb.WriteByte(' ')
	sb.WriteString(v.Name)
	if v.Value.Processes != nil {
		sb.WriteString(" = ")
		sb.WriteString(v.Value.String())
	}
	sb.WriteByte(';')
	return sb.String()
}

func (v VariableAST) StringType() string {
	if v.Type.Code == jn.Void {
		return "auto"
	}
	return v.Type.String()
}

type VarsetSelector struct {
	NewVariable bool
	Variable    VariableAST
	Expression  ExpressionAST
	Ignore      bool
}

func (vs VarsetSelector) String() string {
	if vs.NewVariable {
		return vs.Expression.Tokens[0].Kind
	}
	return vs.Expression.String()
}

type VariableSetAST struct {
	Setter            lexer.Token
	SelectExpressions []VarsetSelector
	ValueExpressions  []ExpressionAST
	JustDeclare       bool
	MultipleReturn    bool
}

func (vs VariableSetAST) cxxSingleSet(cxx *strings.Builder) string {
	cxx.WriteString(vs.SelectExpressions[0].String())
	cxx.WriteString(vs.Setter.Kind)
	cxx.WriteString(vs.ValueExpressions[0].String())
	cxx.WriteByte(';')
	return cxx.String()
}

func (vs VariableSetAST) cxxMultipleSet(cxx *strings.Builder) string {
	cxx.WriteString("std::tie(")
	var expCxx strings.Builder
	expCxx.WriteString("std::make_tuple(")
	for index, selector := range vs.SelectExpressions {
		if selector.Ignore {
			continue
		}
		cxx.WriteString(selector.String())
		cxx.WriteByte(',')
		expCxx.WriteString(vs.ValueExpressions[index].String())
		expCxx.WriteByte(',')
	}
	str := cxx.String()[:cxx.Len()-1] + ")"
	cxx.Reset()
	cxx.WriteString(str)
	cxx.WriteString(vs.Setter.Kind)
	cxx.WriteString(expCxx.String()[:expCxx.Len()-1] + ")")
	cxx.WriteByte(';')
	return cxx.String()
}

func (vs VariableSetAST) cxxMultipleReturn(cxx *strings.Builder) string {
	cxx.WriteString("std::tie(")
	for _, selector := range vs.SelectExpressions {
		if selector.Ignore {
			cxx.WriteString("std::ignore,")
			continue
		}
		cxx.WriteString(selector.String())
		cxx.WriteByte(',')
	}
	str := cxx.String()[:cxx.Len()-1]
	cxx.Reset()
	cxx.WriteString(str)
	cxx.WriteByte(')')
	cxx.WriteString(vs.Setter.Kind)
	cxx.WriteString(vs.ValueExpressions[0].String())
	cxx.WriteByte(';')
	return cxx.String()
}

func (vs VariableSetAST) cxxNewDefines(cxx *strings.Builder) {
	for _, selector := range vs.SelectExpressions {
		if selector.Ignore || !selector.NewVariable {
			continue
		}
		cxx.WriteString(selector.Variable.String() + " ")
	}
}

func (vs VariableSetAST) String() string {
	var cxx strings.Builder
	vs.cxxNewDefines(&cxx)
	if vs.JustDeclare {
		return cxx.String()[:cxx.Len()-1]
	}
	if vs.MultipleReturn {
		return vs.cxxMultipleReturn(&cxx)
	} else if len(vs.SelectExpressions) == 1 {
		return vs.cxxSingleSet(&cxx)
	}
	return vs.cxxMultipleSet(&cxx)
}

type FreeAST struct {
	Token      lexer.Token
	Expression ExpressionAST
}

func (f FreeAST) String() string {
	return "delete " + f.Expression.String() + ";"
}
