package ast

import (
	"fmt"
	"strings"
	"sync/atomic"

	"github.com/De-Rune/jane/lexer"
	"github.com/De-Rune/jane/package/jn"
)

type Object struct {
	Token lexer.Token
	Value interface{}
}

type StatementAST struct {
	Token          lexer.Token
	Value          interface{}
	WithTerminator bool
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

var Indent int32 = 0

func (b BlockAST) String() string {
	atomic.SwapInt32(&Indent, Indent+1)
	defer func() { atomic.SwapInt32(&Indent, Indent-1) }()
	return ParseBlock(b, int(Indent))
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
	cxx := "std::function<"
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
	Token    lexer.Token
	Name     string
	Const    bool
	Variadic bool
	Type     DataTypeAST
}

func (p ParameterAST) String() string {
	var cxx strings.Builder
	cxx.WriteString(p.Prototype())
	if p.Name != "" {
		cxx.WriteByte(' ')
		cxx.WriteString(p.Name)
	}
	if p.Variadic {
		cxx.WriteString(" =array<")
		cxx.WriteString(p.Type.String())
		cxx.WriteString(">()")
	}
	return cxx.String()
}

func (p ParameterAST) Prototype() string {
	var cxx strings.Builder
	if p.Const {
		cxx.WriteString("const ")
	}
	if p.Variadic {
		cxx.WriteString("array<")
		cxx.WriteString(p.Type.String())
		cxx.WriteByte('>')
	} else {
		cxx.WriteString(p.Type.String())
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
	Token lexer.Token
	Expr  ExprAST
}

func (a ArgAST) String() string {
	return a.Expr.String()
}

type ExprAST struct {
	Tokens    []lexer.Token
	Processes [][]lexer.Token
	Model     IExprModel
}

type IExprModel interface {
	String() string
}

func (e ExprAST) String() string {
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

type ExprStatementAST struct {
	Expr ExprAST
}

func (be ExprStatementAST) String() string {
	return be.Expr.String() + ";"
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
	Token lexer.Token
	Expr  ExprAST
}

func (r ReturnAST) String() string {
	return "return " + r.Expr.String() + ";"
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
	Value       ExprAST
	New         bool
	Tag         interface{}
}

func (v VariableAST) String() string {
	var sb strings.Builder
	switch v.DefineToken.Id {
	case lexer.Const:
		sb.WriteString("const ")
	}
	sb.WriteString(v.Type.String())
	sb.WriteByte(' ')
	sb.WriteString(v.Name)
	if v.Value.Processes != nil {
		sb.WriteString(" = ")
		sb.WriteString(v.Value.String())
	}
	sb.WriteByte(';')
	return sb.String()
}

type VarsetSelector struct {
	NewVariable bool
	Variable    VariableAST
	Expr        ExprAST
	Ignore      bool
}

func (vs VarsetSelector) String() string {
	if vs.NewVariable {
		return vs.Expr.Tokens[0].Kind
	}
	return vs.Expr.String()
}

type VariableSetAST struct {
	Setter         lexer.Token
	SelectExprs    []VarsetSelector
	ValueExprs     []ExprAST
	JustDeclare    bool
	MultipleReturn bool
}

func (vs VariableSetAST) cxxSingleSet(cxx *strings.Builder) string {
	cxx.WriteString(vs.SelectExprs[0].String())
	cxx.WriteString(vs.Setter.Kind)
	cxx.WriteString(vs.ValueExprs[0].String())
	cxx.WriteByte(';')
	return cxx.String()
}

func (vs VariableSetAST) cxxMultipleSet(cxx *strings.Builder) string {
	cxx.WriteString("std::tie(")
	var expCxx strings.Builder
	expCxx.WriteString("std::make_tuple(")
	for index, selector := range vs.SelectExprs {
		if selector.Ignore {
			continue
		}
		cxx.WriteString(selector.String())
		cxx.WriteByte(',')
		expCxx.WriteString(vs.ValueExprs[index].String())
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
	for _, selector := range vs.SelectExprs {
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
	cxx.WriteString(vs.ValueExprs[0].String())
	cxx.WriteByte(';')
	return cxx.String()
}

func (vs VariableSetAST) cxxNewDefines(cxx *strings.Builder) {
	for _, selector := range vs.SelectExprs {
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
	} else if len(vs.SelectExprs) == 1 {
		return vs.cxxSingleSet(&cxx)
	}
	return vs.cxxMultipleSet(&cxx)
}

type FreeAST struct {
	Token lexer.Token
	Expr  ExprAST
}

func (f FreeAST) String() string {
	return "delete " + f.Expr.String() + ";"
}

type IterProfile interface {
	String(iter IterAST) string
}

type WhileProfile struct {
	Expr ExprAST
}

func (wp WhileProfile) String(iter IterAST) string {
	var cxx strings.Builder
	cxx.WriteString("while (")
	cxx.WriteString(wp.Expr.String())
	cxx.WriteByte(')')
	cxx.WriteString(iter.Block.String())
	return cxx.String()
}

type ForeachProfile struct {
	KeyA     VariableAST
	KeyB     VariableAST
	InToken  lexer.Token
	Expr     ExprAST
	ExprType DataTypeAST
}

func (fp ForeachProfile) String(iter IterAST) string {
	if !jn.IsIgnoreName(fp.KeyA.Name) {
		return fp.ForeachString(iter)
	}
	return fp.IterationSring(iter)
}

func (fp ForeachProfile) ForeachString(iter IterAST) string {
	var cxx strings.Builder
	cxx.WriteString("foreach<")
	cxx.WriteString(fp.ExprType.String())
	cxx.WriteString("," + fp.KeyA.Type.String())
	if !jn.IsIgnoreName(fp.KeyB.Name) {
		cxx.WriteString("," + fp.KeyB.Type.String())
	}
	cxx.WriteString(">(")
	cxx.WriteString(fp.Expr.String())
	cxx.WriteString(", [&](")
	cxx.WriteString(fp.KeyA.Type.String())
	cxx.WriteString(" " + fp.KeyA.Name)
	if !jn.IsIgnoreName(fp.KeyB.Name) {
		cxx.WriteString(",")
		cxx.WriteString(fp.KeyB.Type.String())
		cxx.WriteString(" " + fp.KeyB.Name)
	}
	cxx.WriteString(") -> void ")
	cxx.WriteString(iter.Block.String())
	cxx.WriteString(");")
	return cxx.String()
}

func (fp ForeachProfile) IterationSring(iter IterAST) string {
	var cxx strings.Builder
	cxx.WriteString("for (auto ")
	cxx.WriteString(fp.KeyB.Name)
	cxx.WriteString(" : ")
	cxx.WriteString(fp.Expr.String())
	cxx.WriteString(") ")
	cxx.WriteString(iter.Block.String())
	return cxx.String()
}

type IterAST struct {
	Token   lexer.Token
	Block   BlockAST
	Profile IterProfile
}

func (iter IterAST) String() string {
	if iter.Profile == nil {
		return "while (true) " + iter.Block.String()
	}
	return iter.Profile.String(iter)
}

type BreakAST struct {
	Token lexer.Token
}

func (b BreakAST) String() string {
	return "break;"
}

type ContinueAST struct {
	Token lexer.Token
}

func (c ContinueAST) String() string {
	return "continue;"
}

type IfAST struct {
	Token lexer.Token
	Expr  ExprAST
	Block BlockAST
}

func (ifast IfAST) String() string {
	var cxx strings.Builder
	cxx.WriteString("if (")
	cxx.WriteString(ifast.Expr.String())
	cxx.WriteString(") ")
	cxx.WriteString(ifast.Block.String())
	return cxx.String()
}

type ElseIfAST struct {
	Token lexer.Token
	Expr  ExprAST
	Block BlockAST
}

func (elif ElseIfAST) String() string {
	var cxx strings.Builder
	cxx.WriteString("else if (")
	cxx.WriteString(elif.Expr.String())
	cxx.WriteString(") ")
	cxx.WriteString(elif.Block.String())
	return cxx.String()
}

type ElseAST struct {
	Token lexer.Token
	Block BlockAST
}

func (elseast ElseAST) String() string {
	var cxx strings.Builder
	cxx.WriteString("else ")
	cxx.WriteString(elseast.Block.String())
	return cxx.String()
}
