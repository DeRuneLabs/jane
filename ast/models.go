package ast

import (
	"fmt"
	"strings"
	"sync/atomic"
	"unicode"

	"github.com/DeRuneLabs/jane/lexer/tokens"
	"github.com/DeRuneLabs/jane/package/jn"
	"github.com/DeRuneLabs/jane/package/jnapi"
	"github.com/DeRuneLabs/jane/package/jntype"
)

type exprNode struct{ expr string }

func (en exprNode) String() string {
	return en.expr
}

type Genericable interface {
	Generics() []DataType
	SetGenerics([]DataType)
}

type Obj struct {
	Tok   Tok
	Value any
}

type Statement struct {
	Tok            Tok
	Val            any
	WithTerminator bool
}

func (s Statement) String() string {
	return fmt.Sprint(s.Val)
}

type (
	Labels []*Label
	Gotos  []*Goto
)

type Block struct {
	Parent   *Block
	SubIndex int
	Tree     []Statement
	Gotos    *Gotos
	Labels   *Labels
	Func     *Func
}

var Indent uint32 = 0

func IndentString() string {
	return strings.Repeat(jn.Set.Indent, int(Indent)*jn.Set.IndentCount)
}

func AddIndent() {
	atomic.AddUint32(&Indent, 1)
}

func DoneIndent() {
	atomic.SwapUint32(&Indent, Indent-1)
}

func (b Block) String() string {
	AddIndent()
	defer func() { DoneIndent() }()
	return ParseBlock(b)
}

func ParseBlock(b Block) string {
	var cxx strings.Builder
	cxx.WriteByte('{')
	for _, s := range b.Tree {
		if s.Val == nil {
			continue
		}
		cxx.WriteByte('\n')
		cxx.WriteString(IndentString())
		cxx.WriteString(s.String())
	}
	cxx.WriteByte('\n')
	cxx.WriteString(strings.Repeat(jn.Set.Indent, int(Indent-1)*jn.Set.IndentCount))
	cxx.WriteByte('}')
	return cxx.String()
}

type genericableTypes struct {
	types []DataType
}

func (gt genericableTypes) Generics() []DataType {
	return gt.types
}

func (gt genericableTypes) SetGenerics([]DataType) {}

type DataType struct {
	Tok        Tok
	Id         uint8
	Original   any
	Val        string
	MultiTyped bool
	Tag        any
}

func (dt *DataType) ValWithOriginalId() string {
	if dt.Original == nil {
		return dt.Val
	}
	_, prefix := dt.GetValId()
	original := dt.Original.(DataType)
	return prefix + original.Tok.Kind
}

func (dt *DataType) OriginalValId() string {
	if dt.Original == nil {
		return ""
	}
	t := dt.Original.(DataType)
	id, _ := t.GetValId()
	return id
}

func (dt *DataType) GetValId() (id, prefix string) {
	id = dt.Val
	runes := []rune(dt.Val)
	for i, r := range dt.Val {
		if r == '_' || unicode.IsLetter(r) {
			id = string(runes[i:])
			prefix = string(runes[:i])
			break
		}
	}
	runes = []rune(id)
	for i, r := range runes {
		if r != '_' && !unicode.IsLetter(r) {
			id = string(runes[:i])
			break
		}
	}
	return
}

func (dt DataType) String() string {
	var cxx strings.Builder
	if dt.Original != nil {
		val := dt.ValWithOriginalId()
		tok := dt.Tok
		dt = dt.Original.(DataType)
		dt.Val = val
		dt.Tok = tok
	}
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
		case dt.Id == jntype.Map && dt.Val[0] == '[':
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
	if dt.Tag != nil {
		switch t := dt.Tag.(type) {
		case Genericable:
			return dt.StructString() + cxx.String()
		case []DataType:
			dt.Tag = genericableTypes{t}
			return dt.StructString() + cxx.String()
		}
	}
	switch dt.Id {
	case jntype.Id, jntype.Enum:
		return jnapi.OutId(dt.Val, dt.Tok.File) + cxx.String()
	case jntype.Struct:
		return dt.StructString() + cxx.String()
	case jntype.Func:
		return dt.FuncString() + cxx.String()
	default:
		return jntype.CxxTypeIdFromType(dt.Id) + cxx.String()
	}
}

func (dt *DataType) StructString() string {
	var cxx strings.Builder
	cxx.WriteString(jnapi.OutId(dt.Val, dt.Tok.File))
	s := dt.Tag.(Genericable)
	types := s.Generics()
	if len(types) == 0 {
		return cxx.String()
	}
	cxx.WriteByte('<')
	for _, t := range types {
		cxx.WriteString(t.String())
		cxx.WriteByte(',')
	}
	return cxx.String()[:cxx.Len()-1] + ">"
}

func (dt *DataType) FuncString() string {
	var cxx strings.Builder
	cxx.WriteString("func<")
	fun := dt.Tag.(*Func)
	cxx.WriteString(fun.RetType.String())
	cxx.WriteByte('(')
	if len(fun.Params) > 0 {
		for _, param := range fun.Params {
			cxx.WriteString(param.Prototype())
			cxx.WriteByte(',')
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

func (dt *DataType) MultiTypeString() string {
	types := dt.Tag.([]DataType)
	var cxx strings.Builder
	cxx.WriteString("std::tuple<")
	for _, t := range types {
		cxx.WriteString(t.String())
		cxx.WriteByte(',')
	}
	return cxx.String()[:cxx.Len()-1] + ">"
}

type GenericType struct {
	Tok Tok
	Id  string
}

func (gt GenericType) String() string {
	var cxx strings.Builder
	cxx.WriteString("typename ")
	cxx.WriteString(jnapi.OutId(gt.Id, gt.Tok.File))
	cxx.WriteByte('>')
	return cxx.String()
}

type Type struct {
	Pub  bool
	Tok  Tok
	Id   string
	Type DataType
	Desc string
	Used bool
}

func (t Type) String() string {
	var cxx strings.Builder
	cxx.WriteString("typedef ")
	cxx.WriteString(t.Type.String())
	cxx.WriteByte(' ')
	cxx.WriteString(jnapi.OutId(t.Id, t.Tok.File))
	cxx.WriteByte(';')
	return cxx.String()
}

type RetType struct {
	Type        DataType
	Identifiers Toks
}

func (rt RetType) String() string {
	return rt.Type.String()
}

func (rt *RetType) AnyVar() bool {
	for _, tok := range rt.Identifiers {
		if !jnapi.IsIgnoreId(tok.Kind) {
			return true
		}
	}
	return false
}

func (rt *RetType) Vars() []*Var {
	if !rt.Type.MultiTyped {
		return nil
	}
	types := rt.Type.Tag.([]DataType)
	var vars []*Var
	for i, tok := range rt.Identifiers {
		if jnapi.IsIgnoreId(tok.Kind) {
			continue
		}
		variable := new(Var)
		variable.IdTok = tok
		variable.Id = tok.Kind
		variable.Type = types[i]
		vars = append(vars, variable)
	}
	return vars
}

type Func struct {
	Pub        bool
	Tok        Tok
	Id         string
	Generics   []*GenericType
	Combines   [][]DataType
	Attributes []Attribute
	Params     []Param
	RetType    RetType
	Block      Block
}

func (f *Func) FindAttribute(kind string) *Attribute {
	for i := range f.Attributes {
		attribute := &f.Attributes[i]
		if attribute.Tag.Kind == kind {
			return attribute
		}
	}
	return nil
}

func (f *Func) DataTypeString() string {
	var cxx strings.Builder
	cxx.WriteByte('(')
	if len(f.Params) > 0 {
		for _, p := range f.Params {
			if p.Variadic {
				cxx.WriteString("...")
			}
			cxx.WriteString(p.Type.Val)
			cxx.WriteString(", ")
		}
		cxxStr := cxx.String()[:cxx.Len()-2]
		cxx.Reset()
		cxx.WriteString(cxxStr)
	}
	cxx.WriteByte(')')
	if f.RetType.Type.Id != jntype.Void {
		cxx.WriteString(f.RetType.Type.Val)
	}
	return cxx.String()
}

type Param struct {
	Tok       Tok
	Id        string
	Const     bool
	Volatile  bool
	Variadic  bool
	Reference bool
	Type      DataType
	Default   Expr
}

func (p Param) String() string {
	var cxx strings.Builder
	cxx.WriteString(p.Prototype())
	if p.Id != "" {
		cxx.WriteByte(' ')
		cxx.WriteString(jnapi.OutId(p.Id, p.Tok.File))
	}
	return cxx.String()
}

func (p *Param) Prototype() string {
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
	if p.Reference {
		cxx.WriteByte('&')
	}
	return cxx.String()
}

type Arg struct {
	Tok      Tok
	TargetId string
	Expr     Expr
}

type Args struct {
	Src      []Arg
	Targeted bool
}

func (a Arg) String() string {
	return a.Expr.String()
}

type Expr struct {
	Toks      []Tok
	Processes [][]Tok
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
			case tokens.Id:
				expr.WriteString(jnapi.OutId(tok.Kind, tok.File))
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
	Tok  Tok
	Data string
	Type DataType
}

func (v Value) String() string {
	return v.Data
}

type Ret struct {
	Tok  Tok
	Expr Expr
}

func (r Ret) String() string {
	var cxx strings.Builder
	cxx.WriteString("return ")
	cxx.WriteString(r.Expr.String())
	cxx.WriteByte(';')
	return cxx.String()
}

type Attribute struct {
	Tok Tok
	Tag Tok
}

func (a Attribute) String() string {
	return a.Tag.Kind
}

type Var struct {
	Pub       bool
	DefTok    Tok
	IdTok     Tok
	SetterTok Tok
	Id        string
	Type      DataType
	Val       Expr
	Const     bool
	Volatile  bool
	New       bool
	Tag       any
	Desc      string
	Used      bool
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
	cxx.WriteString(jnapi.OutId(v.Id, v.IdTok.File))
	cxx.WriteByte('{')
	if v.Val.Processes != nil {
		cxx.WriteString(v.Val.String())
	}
	cxx.WriteByte('}')
	cxx.WriteByte(';')
	return cxx.String()
}

func (v *Var) FieldString() string {
	var cxx strings.Builder
	if v.Volatile {
		cxx.WriteString("volatile ")
	}
	if v.Const {
		cxx.WriteString("const ")
	}
	cxx.WriteString(v.Type.String())
	cxx.WriteByte(' ')
	cxx.WriteString(jnapi.OutId(v.Id, v.IdTok.File))
	cxx.WriteByte(';')
	return cxx.String()
}

type AssignSelector struct {
	Var    Var
	Expr   Expr
	Ignore bool
}

func (as AssignSelector) String() string {
	switch {
	case as.Var.New:
		tok := as.Expr.Toks[0]
		return jnapi.OutId(tok.Kind, tok.File)
	case as.Ignore:
		return jnapi.CxxIgnore
	}
	return as.Expr.String()
}

type Assign struct {
	Setter      Tok
	SelectExprs []AssignSelector
	ValueExprs  []Expr
	IsExpr      bool
	MultipleRet bool
}

func (a *Assign) cxxSingleAssign() string {
	expr := a.SelectExprs[0]
	if expr.Var.New {
		expr.Var.Val = a.ValueExprs[0]
		s := expr.Var.String()
		return s[:len(s)-1]
	}
	var cxx strings.Builder
	if len(expr.Expr.Toks) != 1 ||
		!jnapi.IsIgnoreId(expr.Expr.Toks[0].Kind) {
		cxx.WriteString(expr.String())
		cxx.WriteString(a.Setter.Kind)
	}
	cxx.WriteString(a.ValueExprs[0].String())
	return cxx.String()
}

func (a *Assign) hasSelector() bool {
	for _, s := range a.SelectExprs {
		if !s.Ignore {
			return true
		}
	}
	return false
}

func (a *Assign) cxxMultipleAssign() string {
	var cxx strings.Builder
	if !a.hasSelector() {
		for _, expr := range a.ValueExprs {
			cxx.WriteString(expr.String())
			cxx.WriteByte(';')
		}
		return cxx.String()[:cxx.Len()-1]
	}
	cxx.WriteString(a.cxxNewDefines())
	cxx.WriteString("std::tie(")
	var expCxx strings.Builder
	expCxx.WriteString("std::make_tuple(")
	for i, selector := range a.SelectExprs {
		cxx.WriteString(selector.String())
		cxx.WriteByte(',')
		expCxx.WriteString(a.ValueExprs[i].String())
		expCxx.WriteByte(',')
	}
	str := cxx.String()[:cxx.Len()-1] + ")"
	cxx.Reset()
	cxx.WriteString(str)
	cxx.WriteString(a.Setter.Kind)
	cxx.WriteString(expCxx.String()[:expCxx.Len()-1] + ")")
	return cxx.String()
}

func (a *Assign) cxxMultipleReturn() string {
	var cxx strings.Builder
	cxx.WriteString(a.cxxNewDefines())
	cxx.WriteString("std::tie(")
	for _, selector := range a.SelectExprs {
		if selector.Ignore {
			cxx.WriteString(jnapi.CxxIgnore)
			cxx.WriteByte(',')
			continue
		}
		cxx.WriteString(selector.String())
		cxx.WriteByte(',')
	}
	str := cxx.String()[:cxx.Len()-1]
	cxx.Reset()
	cxx.WriteString(str)
	cxx.WriteByte(')')
	cxx.WriteString(a.Setter.Kind)
	cxx.WriteString(a.ValueExprs[0].String())
	return cxx.String()
}

func (a *Assign) cxxNewDefines() string {
	var cxx strings.Builder
	for _, selector := range a.SelectExprs {
		if selector.Ignore || !selector.Var.New {
			continue
		}
		cxx.WriteString(selector.Var.String() + " ")
	}
	return cxx.String()
}

func (a Assign) String() string {
	var cxx strings.Builder
	switch {
	case a.MultipleRet:
		cxx.WriteString(a.cxxMultipleReturn())
	case len(a.SelectExprs) == 1:
		cxx.WriteString(a.cxxSingleAssign())
	default:
		cxx.WriteString(a.cxxMultipleAssign())
	}
	if !a.IsExpr {
		cxx.WriteByte(';')
	}
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
	InTok    Tok
	Expr     Expr
	ExprType DataType
}

func (fp ForeachProfile) String(iter Iter) string {
	if !jnapi.IsIgnoreId(fp.KeyA.Id) {
		return fp.ForeachString(iter)
	}
	return fp.IterationString(iter)
}

func (fp *ForeachProfile) ClassicString(iter Iter) string {
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
	cxx.WriteString(jnapi.OutId(fp.KeyA.Id, fp.KeyA.IdTok.File))
	if !jnapi.IsIgnoreId(fp.KeyB.Id) {
		cxx.WriteByte(',')
		cxx.WriteString(fp.KeyB.Type.String())
		cxx.WriteByte(' ')
		cxx.WriteString(jnapi.OutId(fp.KeyB.Id, fp.KeyB.IdTok.File))
	}
	cxx.WriteString(") -> void ")
	cxx.WriteString(iter.Block.String())
	cxx.WriteString(");")
	return cxx.String()
}

func (fp *ForeachProfile) MapString(iter Iter) string {
	var cxx strings.Builder
	cxx.WriteString("foreach<")
	types := fp.ExprType.Tag.([]DataType)
	cxx.WriteString(types[0].String())
	cxx.WriteByte(',')
	cxx.WriteString(types[1].String())
	cxx.WriteString(">(")
	cxx.WriteString(fp.Expr.String())
	cxx.WriteString(", [&](")
	cxx.WriteString(fp.KeyA.Type.String())
	cxx.WriteByte(' ')
	cxx.WriteString(jnapi.OutId(fp.KeyA.Id, fp.KeyA.IdTok.File))
	if !jnapi.IsIgnoreId(fp.KeyB.Id) {
		cxx.WriteByte(',')
		cxx.WriteString(fp.KeyB.Type.String())
		cxx.WriteByte(' ')
		cxx.WriteString(jnapi.OutId(fp.KeyB.Id, fp.KeyB.IdTok.File))
	}
	cxx.WriteString(") -> void ")
	cxx.WriteString(iter.Block.String())
	cxx.WriteString(");")
	return cxx.String()
}

func (fp *ForeachProfile) ForeachProfile(iter Iter) string {
	switch {
	case fp.ExprType.Val == tokens.STR,
		strings.HasPrefix(fp.ExprType.Val, "[]"):
		return fp.ClassicString(iter)
	case fp.ExprType.Val[0] == '[':
		return fp.MapString(iter)
	}
	return ""
}

func (fp *ForeachProfile) ForeachString(iter Iter) string {
	switch {
	case fp.ExprType.Val == tokens.STR,
		strings.HasPrefix(fp.ExprType.Val, "[]"):
		return fp.ClassicString(iter)
	case fp.ExprType.Val[0] == '[':
		return fp.MapString(iter)
	}
	return ""
}

func (fp ForeachProfile) IterationString(iter Iter) string {
	var cxx strings.Builder
	cxx.WriteString("for (auto ")
	cxx.WriteString(jnapi.OutId(fp.KeyB.Id, fp.KeyB.IdTok.File))
	cxx.WriteString(" : ")
	cxx.WriteString(fp.Expr.String())
	cxx.WriteString(") ")
	cxx.WriteString(iter.Block.String())
	return cxx.String()
}

type Iter struct {
	Tok     Tok
	Block   Block
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
	Tok Tok
}

func (b Break) String() string {
	return "break;"
}

type Continue struct {
	Tok Tok
}

func (c Continue) String() string {
	return "continue;"
}

type If struct {
	Tok   Tok
	Expr  Expr
	Block Block
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
	Tok   Tok
	Expr  Expr
	Block Block
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
	Tok   Tok
	Block Block
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
	Tok  Tok
	Path string
}

type CxxEmbed struct {
	Tok     Tok
	Content string
}

func (ce CxxEmbed) String() string {
	return ce.Content
}

type Preprocessor struct {
	Tok     Tok
	Command any
}

func (pp Preprocessor) String() string {
	return fmt.Sprint(pp.Command)
}

type Directive struct {
	Command any
}

func (d Directive) String() string {
	return fmt.Sprint(d.Command)
}

type EnofiDirective struct{}

func (EnofiDirective) String() string {
	return ""
}

type Defer struct {
	Tok  Tok
	Expr Expr
}

func (d Defer) String() string {
	return jnapi.ToDeferredCall(d.Expr.String())
}

type Label struct {
	Tok   Tok
	Label string
	Index int
	Used  bool
	Block *Block
}

func (l Label) String() string {
	return l.Label + ":;"
}

type Goto struct {
	Tok   Tok
	Label string
	Index int
	Block *Block
}

func (gt Goto) String() string {
	var cxx strings.Builder
	cxx.WriteString("goto ")
	cxx.WriteString(gt.Label)
	cxx.WriteByte(';')
	return cxx.String()
}

type Namespace struct {
	Tok  Tok
	Ids  []string
	Tree []Obj
}

type EnumItem struct {
	Tok  Tok
	Id   string
	Expr Expr
}

func (ei EnumItem) String() string {
	var cxx strings.Builder
	cxx.WriteString(jnapi.OutId(ei.Id, ei.Tok.File))
	cxx.WriteString(" = ")
	cxx.WriteString(ei.Expr.String())
	return cxx.String()
}

type Enum struct {
	Pub   bool
	Tok   Tok
	Id    string
	Type  DataType
	Items []*EnumItem
	Used  bool
	Desc  string
}

func (e *Enum) ItemById(id string) *EnumItem {
	for _, item := range e.Items {
		if item.Id == id {
			return item
		}
	}
	return nil
}

func (e Enum) String() string {
	var cxx strings.Builder
	cxx.WriteString("enum ")
	cxx.WriteString(jnapi.OutId(e.Id, e.Tok.File))
	cxx.WriteByte(':')
	cxx.WriteString(e.Type.String())
	cxx.WriteString(" {\n")
	AddIndent()
	for _, item := range e.Items {
		cxx.WriteString(IndentString())
		cxx.WriteString(item.String())
		cxx.WriteString(",\n")
	}
	DoneIndent()
	cxx.WriteString("};")
	return cxx.String()
}

type Struct struct {
	Tok      Tok
	Id       string
	Pub      bool
	Fields   []*Var
	Generics []*GenericType
}

type ConcurrentCall struct {
	Tok  Tok
	Expr Expr
}

func (cc ConcurrentCall) String() string {
	return jnapi.ToConcurrentCall(cc.Expr.String())
}

type Try struct {
	Tok   Tok
	Block Block
	Catch Catch
}

func (t Try) String() string {
	var cxx strings.Builder
	cxx.WriteString("try ")
	cxx.WriteString(t.Block.String())
	if t.Catch.Tok.Id == tokens.NA {
		cxx.WriteString(" catch(...) {}")
	} else {
		cxx.WriteByte(' ')
		cxx.WriteString(t.Catch.String())
	}
	return cxx.String()
}

type Catch struct {
	Tok   Tok
	Var   Var
	Block Block
}

func (c Catch) String() string {
	var cxx strings.Builder
	cxx.WriteString("catch (")
	if c.Var.Id == "" {
		cxx.WriteString("...")
	} else {
		cxx.WriteString(c.Var.Type.String())
		cxx.WriteByte(' ')
		cxx.WriteString(jnapi.OutId(c.Var.Id, c.Tok.File))
	}
	cxx.WriteString(") ")
	cxx.WriteString(c.Block.String())
	return cxx.String()
}
