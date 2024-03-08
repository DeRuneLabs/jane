package ast

import (
	"fmt"
	"strings"
	"sync/atomic"

	"github.com/DeRuneLabs/jane/lexer"
	"github.com/DeRuneLabs/jane/package/jn"
	"github.com/DeRuneLabs/jane/package/jnapi"
)

type Obj struct {
	Token lexer.Token
	Value interface{}
}

type Statement struct {
	Token          lexer.Token
	Val            interface{}
	WithTerminator bool
}

func (s Statement) String() string {
	return fmt.Sprint(s.Val)
}

type BlockAST struct {
	Tree []Statement
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
	for _, s := range b.Tree {
		cxx.WriteByte('\n')
		cxx.WriteString(strings.Repeat(" ", indent*IndentSpace))
		cxx.WriteString(s.String())
	}
	cxx.WriteByte('\n')
	cxx.WriteString(strings.Repeat(" ", (indent-1)*IndentSpace) + "}")
	return cxx.String()
}

type DataType struct {
	Token      lexer.Token
	Id         uint8
	Val        string
	MultiTyped bool
	Tag        interface{}
}

func (dt DataType) String() string {
	var cxx strings.Builder
	for i, run := range dt.Val {
		if run == '*' {
			cxx.WriteRune(run)
			continue
		}
		dt.Val = dt.Val[i:]
		break
	}
	if dt.MultiTyped {
		return dt.MultiTypeString() + cxx.String()
	}
	if dt.Val != "" {
		switch {
		case strings.HasPrefix(dt.Val, "[]"):
			pointers := cxx.String()
			cxx.Reset()
			cxx.WriteString("array<")
			dt.Val = dt.Val[2:]
			cxx.WriteString(dt.String())
			cxx.WriteByte('>')
			cxx.WriteString(pointers)
			return cxx.String()
		case dt.Id == jn.Map && dt.Val[0] == '[':
			pointers := cxx.String()
			types := dt.Tag.([]DataType)
			cxx.Reset()
			cxx.WriteString("map<")
			cxx.WriteString(types[0].String())
			cxx.WriteByte(',')
			cxx.WriteString(types[1].String())
			cxx.WriteByte('>')
			cxx.WriteString(pointers)
			return cxx.String()
		}
	}
	switch dt.Id {
	case jn.Id:
		return jnapi.AsId(dt.Token.Kind) + cxx.String()
	case jn.Func:
		return dt.FunctionString() + cxx.String()
	default:
		return jn.CxxTypeIdFromType(dt.Id) + cxx.String()
	}
}

func (dt DataType) FunctionString() string {
	var cxx strings.Builder
	cxx.WriteString("func<")
	fun := dt.Tag.(Func)
	cxx.WriteString(fun.RetType.String())
	cxx.WriteByte('(')
	if len(fun.Params) > 0 {
		for _, param := range fun.Params {
			cxx.WriteString(param.Type.String())
			cxx.WriteString(", ")
		}
		cxxStr := cxx.String()[:cxx.Len()-1]
		cxx.Reset()
		cxx.WriteString(cxxStr)
	} else {
		cxx.WriteString("void")
	}
	cxx.WriteString(")>")
	return cxx.String()
}

func (dt DataType) MultiTypeString() string {
	types := dt.Tag.([]DataType)
	var cxx strings.Builder
	cxx.WriteString("std::tuple<")
	for _, t := range types {
		cxx.WriteString(t.String())
		cxx.WriteByte(',')
	}
	return cxx.String()[:cxx.Len()-1] + ">"
}

type Type struct {
	Pub   bool
	Token lexer.Token
	Id    string
	Type  DataType
	Desc  string
}

func (t Type) String() string {
	var cxx strings.Builder
	cxx.WriteString("typedef ")
	cxx.WriteString(t.Type.String())
	cxx.WriteByte(' ')
	cxx.WriteString(jnapi.AsId(t.Id))
	cxx.WriteByte(';')
	return cxx.String()
}

type Func struct {
	Pub     bool
	Token   lexer.Token
	Id      string
	Params  []Parameter
	RetType DataType
	Block   BlockAST
}

func (fc Func) DataTypeString() string {
	var cxx strings.Builder
	cxx.WriteByte('(')
	if len(fc.Params) > 0 {
		for _, param := range fc.Params {
			cxx.WriteString(param.Type.String())
			cxx.WriteString(", ")
		}
		cxxStr := cxx.String()[:cxx.Len()-2]
		cxx.Reset()
		cxx.WriteString(cxxStr)
	}
	cxx.WriteByte(')')
	if fc.RetType.Id != jn.Void {
		cxx.WriteString(fc.RetType.String())
	}
	return cxx.String()
}

type Parameter struct {
	Token    lexer.Token
	Id       string
	Const    bool
	Volatile bool
	Variadic bool
	Type     DataType
}

func (p Parameter) String() string {
	var cxx strings.Builder
	cxx.WriteString(p.Prototype())
	if p.Id != "" {
		cxx.WriteByte(' ')
		cxx.WriteString(jnapi.AsId(p.Id))
	}
	if p.Variadic {
		cxx.WriteString(" =array<")
		cxx.WriteString(p.Type.String())
		cxx.WriteString(">()")
	}
	return cxx.String()
}

func (p Parameter) Prototype() string {
	var cxx strings.Builder
	if p.Volatile {
		cxx.WriteString("volatile ")
	}
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

type Arg struct {
	Token lexer.Token
	Expr  Expr
}

func (a Arg) String() string {
	return a.Expr.String()
}

type Expr struct {
	Tokens    []lexer.Token
	Processes [][]lexer.Token
	Model     IExprModel
}

type IExprModel interface {
	String() string
}

func (e Expr) String() string {
	if e.Model != nil {
		return e.Model.String()
	}
	var expr strings.Builder
	for _, process := range e.Processes {
		for _, tok := range process {
			switch tok.Id {
			case lexer.Id:
				expr.WriteString(jnapi.AsId(tok.Kind))
			default:
				expr.WriteString(tok.Kind)
			}
		}
	}
	return expr.String()
}

type ExprStatement struct {
	Expr Expr
}

func (be ExprStatement) String() string {
	var cxx strings.Builder
	cxx.WriteString(be.Expr.String())
	cxx.WriteByte(';')
	return cxx.String()
}

type Value struct {
	Token lexer.Token
	Data  string
	Type  DataType
}

func (v Value) String() string {
	return v.Data
}

type Ret struct {
	Token lexer.Token
	Expr  Expr
}

func (r Ret) String() string {
	var cxx strings.Builder
	cxx.WriteString("return ")
	cxx.WriteString(r.Expr.String())
	cxx.WriteByte(';')
	return cxx.String()
}

type Attribute struct {
	Token lexer.Token
	Tag   lexer.Token
}

func (a Attribute) String() string {
	return a.Tag.Kind
}

type Var struct {
	Pub         bool
	DefToken    lexer.Token
	IdToken     lexer.Token
	SetterToken lexer.Token
	Id          string
	Type        DataType
	Val         Expr
	Const       bool
	Volatile    bool
	New         bool
	Tag         interface{}
	Desc        string
}

func (v Var) String() string {
	var cxx strings.Builder
	if v.Volatile {
		cxx.WriteString("volatile ")
	}
	if v.Const {
		cxx.WriteString("const ")
	}
	cxx.WriteString(v.Type.String())
	cxx.WriteByte(' ')
	cxx.WriteString(jnapi.AsId(v.Id))
	cxx.WriteByte('{')
	if v.Val.Processes != nil {
		cxx.WriteString(v.Val.String())
	}
	cxx.WriteByte('}')
	cxx.WriteByte(';')
	return cxx.String()
}

type AssignSelector struct {
	Var    Var
	Expr   Expr
	Ignore bool
}

func (vs AssignSelector) String() string {
	if vs.Var.New {
		return jnapi.AsId(vs.Expr.Tokens[0].Kind)
	}
	return vs.Expr.String()
}

type Assign struct {
	Setter      lexer.Token
	SelectExprs []AssignSelector
	ValueExprs  []Expr
	IsExpr      bool
	MultipleRet bool
}

func (vs Assign) cxxSingleAssign() string {
	expr := vs.SelectExprs[0]
	if expr.Var.New {
		expr.Var.Val = vs.ValueExprs[0]
		s := expr.Var.String()
		return s[:len(s)-1]
	}
	var cxx strings.Builder
	if len(expr.Expr.Tokens) != 1 ||
		!jnapi.IsIgnoreId(expr.Expr.Tokens[0].Kind) {
		cxx.WriteString(expr.String())
		cxx.WriteString(vs.Setter.Kind)
	}
	cxx.WriteString(vs.ValueExprs[0].String())
	return cxx.String()
}

func (vs Assign) cxxMultipleAssign() string {
	var cxx strings.Builder
	cxx.WriteString(vs.cxxNewDefines())
	cxx.WriteString("std::tie(")
	var expCxx strings.Builder
	expCxx.WriteString("std::make_tuple(")
	for i, selector := range vs.SelectExprs {
		if selector.Ignore {
			continue
		}
		cxx.WriteString(selector.String())
		cxx.WriteByte(',')
		expCxx.WriteString(vs.ValueExprs[i].String())
		expCxx.WriteByte(',')
	}
	str := cxx.String()[:cxx.Len()-1] + ")"
	cxx.Reset()
	cxx.WriteString(str)
	cxx.WriteString(vs.Setter.Kind)
	cxx.WriteString(expCxx.String()[:expCxx.Len()-1] + ")")
	return cxx.String()
}

func (vs Assign) cxxMultipleReturn() string {
	var cxx strings.Builder
	cxx.WriteString(vs.cxxNewDefines())
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
	return cxx.String()
}

func (vs Assign) cxxNewDefines() string {
	var cxx strings.Builder
	for _, selector := range vs.SelectExprs {
		if selector.Ignore || !selector.Var.New {
			continue
		}
		cxx.WriteString(selector.Var.String() + " ")
	}
	return cxx.String()
}

func (vs Assign) String() string {
	var cxx strings.Builder
	switch {
	case vs.MultipleRet:
		cxx.WriteString(vs.cxxMultipleReturn())
	case len(vs.SelectExprs) == 1:
		cxx.WriteString(vs.cxxSingleAssign())
	default:
		cxx.WriteString(vs.cxxMultipleAssign())
	}
	if !vs.IsExpr {
		cxx.WriteByte(';')
	}
	return cxx.String()
}

type Free struct {
	Token lexer.Token
	Expr  Expr
}

func (f Free) String() string {
	var cxx strings.Builder
	cxx.WriteString("delete ")
	cxx.WriteString(f.Expr.String())
	cxx.WriteByte(';')
	return cxx.String()
}

type IterProfile interface {
	String(iter Iter) string
}

type WhileProfile struct {
	Expr Expr
}

func (wp WhileProfile) String(iter Iter) string {
	var cxx strings.Builder
	cxx.WriteString("while (")
	cxx.WriteString(wp.Expr.String())
	cxx.WriteString(") ")
	cxx.WriteString(iter.Block.String())
	return cxx.String()
}

type ForeachProfile struct {
	KeyA     Var
	KeyB     Var
	InToken  lexer.Token
	Expr     Expr
	ExprType DataType
}

func (fp ForeachProfile) String(iter Iter) string {
	if !jnapi.IsIgnoreId(fp.KeyA.Id) {
		return fp.ForeachString(iter)
	}
	return fp.IterationSring(iter)
}

func (fp ForeachProfile) ClassicString(iter Iter) string {
	var cxx strings.Builder
	cxx.WriteString("foreach<")
	cxx.WriteString(fp.ExprType.String())
	cxx.WriteByte(',')
	cxx.WriteString(fp.KeyA.Type.String())
	if !jnapi.IsIgnoreId(fp.KeyB.Id) {
		cxx.WriteByte(',')
		cxx.WriteString(fp.KeyB.Type.String())
	}
	cxx.WriteString(">(")
	cxx.WriteString(fp.Expr.String())
	cxx.WriteString(", [&](")
	cxx.WriteString(fp.KeyA.Type.String())
	cxx.WriteByte(' ')
	cxx.WriteString(jnapi.AsId(fp.KeyA.Id))
	if !jnapi.IsIgnoreId(fp.KeyB.Id) {
		cxx.WriteByte(',')
		cxx.WriteString(fp.KeyB.Type.String())
		cxx.WriteByte(' ')
		cxx.WriteString(jnapi.AsId(fp.KeyB.Id))
	}
	cxx.WriteString(") -> void ")
	cxx.WriteString(iter.Block.String())
	cxx.WriteString(");")
	return cxx.String()
}

func (fp ForeachProfile) MapString(iter Iter) string {
	var cxx strings.Builder
	cxx.WriteString("for (auto ")
	cxx.WriteString(jnapi.AsId(fp.KeyB.Id))
	cxx.WriteString(" : ")
	cxx.WriteString(fp.Expr.String())
	cxx.WriteString(") ")
	cxx.WriteString(iter.Block.String())
	return cxx.String()
}

func (fp *ForeachProfile) ForeachString(iter Iter) string {
	switch {
	case fp.ExprType.Val == "str",
		strings.HasPrefix(fp.ExprType.Val, "[]"):
		return fp.ClassicString(iter)
	case fp.ExprType.Val[0] == '[':
		return fp.MapString(iter)
	}
	return ""
}

func (fp ForeachProfile) IterationSring(iter Iter) string {
	var cxx strings.Builder
	cxx.WriteString("for (auto ")
	cxx.WriteString(jnapi.AsId(fp.KeyB.Id))
	cxx.WriteString(" : ")
	cxx.WriteString(fp.Expr.String())
	cxx.WriteString(") ")
	cxx.WriteString(iter.Block.String())
	return cxx.String()
}

type Iter struct {
	Token   lexer.Token
	Block   BlockAST
	Profile IterProfile
}

func (iter Iter) String() string {
	if iter.Profile == nil {
		var cxx strings.Builder
		cxx.WriteString("while (true) ")
		cxx.WriteString(iter.Block.String())
		return cxx.String()
	}
	return iter.Profile.String(iter)
}

type Break struct {
	Token lexer.Token
}

func (b Break) String() string {
	return "break;"
}

type Continue struct {
	Token lexer.Token
}

func (c Continue) String() string {
	return "continue;"
}

type If struct {
	Token lexer.Token
	Expr  Expr
	Block BlockAST
}

func (ifast If) String() string {
	var cxx strings.Builder
	cxx.WriteString("if (")
	cxx.WriteString(ifast.Expr.String())
	cxx.WriteString(") ")
	cxx.WriteString(ifast.Block.String())
	return cxx.String()
}

type ElseIf struct {
	Token lexer.Token
	Expr  Expr
	Block BlockAST
}

func (elif ElseIf) String() string {
	var cxx strings.Builder
	cxx.WriteString("else if (")
	cxx.WriteString(elif.Expr.String())
	cxx.WriteString(") ")
	cxx.WriteString(elif.Block.String())
	return cxx.String()
}

type Else struct {
	Token lexer.Token
	Block BlockAST
}

func (elseast Else) String() string {
	var cxx strings.Builder
	cxx.WriteString("else ")
	cxx.WriteString(elseast.Block.String())
	return cxx.String()
}

type Comment struct {
	Content string
}

func (c Comment) String() string {
	var cxx strings.Builder
	cxx.WriteString("// ")
	cxx.WriteString(c.Content)
	return cxx.String()
}

type Use struct {
	Token lexer.Token
	Path  string
}

type CxxEmbed struct {
	Content string
}

func (ce CxxEmbed) String() string {
	return ce.Content
}

type Preprocessor struct {
	Token   lexer.Token
	Command interface{}
}

func (pp Preprocessor) String() string {
	return fmt.Sprint(pp.Command)
}

type Directive struct {
	Command interface{}
}

func (d Directive) String() string {
	return fmt.Sprint(d.Command)
}

type EnofiDirective struct{}

func (EnofiDirective) String() string {
	return ""
}
