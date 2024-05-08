// Copyright (c) 2024 arfy slowy - DeRuneLabs
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package ast

import (
	"strconv"
	"strings"

	"github.com/DeRuneLabs/jane/build"
	"github.com/DeRuneLabs/jane/lexer"
)

type Arg struct {
	Token    lexer.Token
	TargetId string
	Expr     Expr
}

func (a Arg) String() string {
	var cpp strings.Builder
	if a.TargetId != "" {
		cpp.WriteString(a.TargetId)
		cpp.WriteString(": ")
	}
	cpp.WriteString(a.Expr.String())
	return cpp.String()
}

type Args struct {
	Src                      []Arg
	Targeted                 bool
	Generics                 []Type
	DynamicGenericAnnotation bool
	NeedsPureType            bool
}

type AssignLeft struct {
	Var    Var
	Expr   Expr
	Ignore bool
}

type Assign struct {
	Setter      lexer.Token
	Left        []AssignLeft
	Right       []Expr
	IsExpr      bool
	MultipleRet bool
}

type Attribute struct {
	Token lexer.Token
	Tag   string
}

func HasAttribute(kind string, attributes []Attribute) bool {
	for i := range attributes {
		attribute := attributes[i]
		if attribute.Tag == kind {
			return true
		}
	}
	return false
}

type RecoverCall struct {
	Try     *Block
	Handler *Fn
}

type Block struct {
	IsUnsafe bool
	Deferred bool
	Parent   *Block
	SubIndex int
	Tree     []St
	Func     *Fn

	Gotos  *Gotos
	Labels *Labels
}

type ConcurrentCall struct {
	Token lexer.Token
	Expr  Expr
}

type Comment struct {
	Token   lexer.Token
	Content string
}

func (c Comment) String() string {
	return "// " + c.Content
}

type If struct {
	Token lexer.Token
	Expr  Expr
	Block *Block
}

type Else struct {
	Token lexer.Token
	Block *Block
}
type Conditional struct {
	If      *If
	Elifs   []*If
	Default *Else
}

type CppLinkFn struct {
	Token lexer.Token
	Link  *Fn
}

type CppLinkVar struct {
	Token lexer.Token
	Link  *Var
}

type CppLinkStruct struct {
	Token lexer.Token
	Link  Struct
}

type CppLinkAlias struct {
	Token lexer.Token
	Link  TypeAlias
}

type Data struct {
	Token    lexer.Token
	Value    string
	DataType Type
}

func (d Data) String() string {
	return d.Value
}

type EnumItem struct {
	Token   lexer.Token
	Id      string
	Expr    Expr
	ExprTag any
}

type Enum struct {
	Pub      bool
	Token    lexer.Token
	Id       string
	DataType Type
	Items    []*EnumItem
	Used     bool
	Doc      string
}

func (e *Enum) ItemById(id string) *EnumItem {
	for _, item := range e.Items {
		if item.Id == id {
			return item
		}
	}
	return nil
}

type BinopExpr struct {
	Tokens []lexer.Token
}

type Binop struct {
	L  any
	R  any
	Op lexer.Token
}

type Expr struct {
	Tokens []lexer.Token
	Op     any
	Model  ExprModel
}

func (e *Expr) IsNotBinop() bool {
	switch e.Op.(type) {
	case Binop:
		return true
	default:
		return false
	}
}

func (e *Expr) IsEmpty() bool {
	return e.Op == nil
}

func (e Expr) String() string {
	if e.Model != nil {
		return e.Model.String()
	}
	return ""
}

type Fn struct {
	Public        bool
	IsUnsafe      bool
	IsEntryPoint  bool
	Used          bool
	Token         lexer.Token
	Id            string
	Generics      []*GenericType
	Combines      *[][]Type
	Attributes    []Attribute
	Params        []Param
	RetType       RetType
	Block         *Block
	Receiver      *Var
	Owner         any
	BuiltinCaller any
	Doc           string
}

func (f *Fn) IsConstructor() bool {
	if f.RetType.DataType.Id != struct_t {
		return false
	}
	s := f.RetType.DataType.Tag.(*Struct)
	return s.Id == f.Id
}

func (f *Fn) plain_type_kind() string {
	var s strings.Builder
	if len(f.Generics) > 0 {
		s.WriteByte('[')
		for i, t := range f.Generics {
			s.WriteString(t.Id)
			if i+1 < len(f.Generics) {
				s.WriteByte(',')
			}
		}
		s.WriteByte(']')
	}
	s.WriteByte('(')
	n := len(f.Params)
	if f.Receiver != nil {
		s.WriteString(f.Receiver.ReceiverTypeString())
		if n > 0 {
			s.WriteString(", ")
		}
	}
	if n > 0 {
		for _, p := range f.Params {
			if p.Variadic {
				s.WriteString("...")
			}
			s.WriteString(p.TypeString())
			s.WriteString(", ")
		}
		cppStr := s.String()[:s.Len()-2]
		s.Reset()
		s.WriteString(cppStr)
	}
	s.WriteByte(')')
	if f.RetType.DataType.MultiTyped {
		s.WriteByte('(')
		for _, t := range f.RetType.DataType.Tag.([]Type) {
			s.WriteString(t.Kind)
			s.WriteByte(',')
		}
		return s.String()[:s.Len()-1] + ")"
	} else if f.RetType.DataType.Id != void_t {
		s.WriteString(f.RetType.DataType.Kind)
	}
	return s.String()
}

func (f *Fn) CppKind(pure bool) string {
	f.RetType.DataType.Pure = pure
	var cpp strings.Builder
	cpp.WriteString(f.RetType.String())
	cpp.WriteByte('(')
	if len(f.Params) > 0 {
		for _, param := range f.Params {
			param.DataType.Pure = pure
			cpp.WriteString(param.Prototype())
			cpp.WriteByte(',')
		}
		cppStr := cpp.String()[:cpp.Len()-1]
		cpp.Reset()
		cpp.WriteString(cppStr)
	} else {
		cpp.WriteString("void")
	}
	cpp.WriteByte(')')
	return cpp.String()
}

func (f *Fn) TypeKind() string {
	var cpp strings.Builder
	if f.IsUnsafe {
		cpp.WriteString("unsafe ")
	}
	cpp.WriteString("fn")
	cpp.WriteString(f.plain_type_kind())
	return cpp.String()
}

func (f *Fn) OutId() string {
	if f.IsEntryPoint {
		return build.OutId(f.Id, 0)
	}
	if f.Receiver != nil {
		return "_method_" + f.Id
	}
	return build.OutId(f.Id, f.Token.File.Addr())
}

func (f *Fn) DefineString() string {
	var s strings.Builder
	if f.IsUnsafe {
		s.WriteString("unsafe ")
	}
	s.WriteString("fn ")
	s.WriteString(f.Id)
	s.WriteString(f.plain_type_kind())
	return s.String()
}

func (f *Fn) PrototypeParams() string {
	if len(f.Params) == 0 {
		return "(void)"
	}
	var cpp strings.Builder
	cpp.WriteByte('(')
	for _, p := range f.Params {
		cpp.WriteString(p.Prototype())
		cpp.WriteByte(',')
	}
	return cpp.String()[:cpp.Len()-1] + ")"
}

type GenericType struct {
	Token lexer.Token
	Id    string
}

func (gt *GenericType) OutId() string {
	return build.AsId(gt.Id)
}

func (gt GenericType) String() string {
	var cpp strings.Builder
	cpp.WriteString("typename ")
	cpp.WriteString(gt.OutId())
	return cpp.String()
}

type (
	Labels []*Label
	Gotos  []*Goto
)

type Label struct {
	Token lexer.Token
	Label string
	Index int
	Used  bool
	Block *Block
}

func (l Label) String() string {
	return l.Label + ":;"
}

type Goto struct {
	Token lexer.Token
	Label string
	Index int
	Block *Block
}

func (gt Goto) String() string {
	var cpp strings.Builder
	cpp.WriteString("goto ")
	cpp.WriteString(gt.Label)
	cpp.WriteByte(';')
	return cpp.String()
}

type Impl struct {
	Base   lexer.Token
	Target Type
	Tree   []Node
}

type Genericable interface {
	GetGenerics() []Type
	SetGenerics([]Type)
}

type ExprModel interface {
	String() string
}

type IterForeach struct {
	KeyA     Var
	KeyB     Var
	InToken  lexer.Token
	Expr     Expr
	ExprType Type
}

type IterWhile struct {
	Expr Expr
	Next St
}

type Break struct {
	Token      lexer.Token
	LabelToken lexer.Token
	Label      string
}

func (b Break) String() string {
	return "goto " + b.Label + ";"
}

type Continue struct {
	Token     lexer.Token
	LoopLabel lexer.Token
	Label     string
}

func (c Continue) String() string {
	return "goto " + c.Label + ";"
}

type Iter struct {
	Token   lexer.Token
	Block   *Block
	Parent  *Block
	Profile any
}

func (i *Iter) BeginLabel() string {
	var cpp strings.Builder
	cpp.WriteString("iter_begin_")
	cpp.WriteString(strconv.Itoa(i.Token.Row))
	cpp.WriteString(strconv.Itoa(i.Token.Column))
	return cpp.String()
}

func (i *Iter) EndLabel() string {
	var cpp strings.Builder
	cpp.WriteString("iter_end_")
	cpp.WriteString(strconv.Itoa(i.Token.Row))
	cpp.WriteString(strconv.Itoa(i.Token.Column))
	return cpp.String()
}

func (i *Iter) NextLabel() string {
	var cpp strings.Builder
	cpp.WriteString("iter_next_")
	cpp.WriteString(strconv.Itoa(i.Token.Row))
	cpp.WriteString(strconv.Itoa(i.Token.Column))
	return cpp.String()
}

type Fall struct {
	Token lexer.Token
	Case  *Case
}

type Case struct {
	Token lexer.Token
	Exprs []Expr
	Block *Block
	Match *Match
	Next  *Case
}

func (c *Case) BeginLabel() string {
	var cpp strings.Builder

	cpp.WriteString("case_begin_")
	cpp.WriteString(strconv.Itoa(c.Token.Row))
	cpp.WriteString(strconv.Itoa(c.Token.Column))
	return cpp.String()
}

func (c *Case) EndLabel() string {
	var cpp strings.Builder

	cpp.WriteString("case_end_")
	cpp.WriteString(strconv.Itoa(c.Token.Row))
	cpp.WriteString(strconv.Itoa(c.Token.Column))
	return cpp.String()
}

type Match struct {
	Token     lexer.Token
	Expr      Expr
	ExprType  Type
	Default   *Case
	TypeMatch bool
	Cases     []Case
}

func (m *Match) EndLabel() string {
	var cpp strings.Builder
	cpp.WriteString("match_end_")
	cpp.WriteString(strconv.Itoa(m.Token.Row))
	cpp.WriteString(strconv.Itoa(m.Token.Column))
	return cpp.String()
}

type Namespace struct {
	Token   lexer.Token
	Id      string
	Defines *Defmap
}

type Node struct {
	Token lexer.Token
	Data  any
}

type Param struct {
	Token    lexer.Token
	Id       string
	Variadic bool
	Mutable  bool
	DataType Type
	Default  Expr
}

func (p *Param) TypeString() string {
	var ts strings.Builder
	if p.Mutable {
		ts.WriteString(lexer.KND_MUT + " ")
	}
	if p.Variadic {
		ts.WriteString(lexer.KND_TRIPLE_DOT)
	}
	ts.WriteString(p.DataType.Kind)
	return ts.String()
}

func (p *Param) OutId() string {
	return as_local_id(p.Token.Row, p.Token.Column, p.Id)
}

func (p Param) String() string {
	var cpp strings.Builder
	cpp.WriteString(p.Prototype())
	if p.Id != "" && !lexer.IsIgnoreId(p.Id) && p.Id != lexer.ANONYMOUS_ID {
		cpp.WriteByte(' ')
		cpp.WriteString(p.OutId())
	}
	return cpp.String()
}

func (p *Param) Prototype() string {
	var cpp strings.Builder
	if p.Variadic {
		cpp.WriteString(build.AsTypeId("slice"))
		cpp.WriteByte('<')
		cpp.WriteString(p.DataType.String())
		cpp.WriteByte('>')
	} else {
		cpp.WriteString(p.DataType.String())
	}
	return cpp.String()
}

type RetType struct {
	DataType    Type
	Identifiers []lexer.Token
}

func (rt RetType) String() string {
	return rt.DataType.String()
}

func (rt *RetType) AnyVar() bool {
	for _, tok := range rt.Identifiers {
		if !lexer.IsIgnoreId(tok.Kind) {
			return true
		}
	}
	return false
}

func (rt *RetType) Vars(owner *Block) []*Var {
	get := func(tok lexer.Token, t Type) *Var {
		v := new(Var)
		v.Token = tok
		if lexer.IsIgnoreId(tok.Kind) {
			v.Id = lexer.IGNORE_ID
		} else {
			v.Id = tok.Kind
		}
		v.DataType = t
		v.Owner = owner
		v.Mutable = true
		return v
	}
	if !rt.DataType.MultiTyped {
		if len(rt.Identifiers) > 0 {
			v := get(rt.Identifiers[0], rt.DataType)
			if v == nil {
				return nil
			}
			return []*Var{v}
		}
		return nil
	}
	var vars []*Var
	types := rt.DataType.Tag.([]Type)
	for i, tok := range rt.Identifiers {
		v := get(tok, types[i])
		if v != nil {
			vars = append(vars, v)
		}
	}
	return vars
}

type Ret struct {
	Token lexer.Token
	Expr  Expr
}

type St struct {
	Token          lexer.Token
	Data           any
	WithTerminator bool
}

type ExprSt struct {
	Expr Expr
}

type Struct struct {
	Token       lexer.Token
	Id          string
	Pub         bool
	Fields      []*Var
	Attributes  []Attribute
	Generics    []*GenericType
	Owner       any
	Origin      *Struct
	Traits      []*Trait
	Defines     *Defmap
	Used        bool
	Doc         string
	CppLinked   bool
	Constructor *Fn
	Depends     []*Struct
	Order       int
	_generics   []Type
}

func (s *Struct) IsSameBase(s2 *Struct) bool {
	return s.Origin == s2.Origin
}

func (s *Struct) IsDependedTo(s2 *Struct) bool {
	for _, d := range s.Origin.Depends {
		if s2.IsSameBase(d) {
			return true
		}
	}
	return false
}

func (s *Struct) OutId() string {
	if s.CppLinked {
		return s.Id
	}
	return build.OutId(s.Id, s.Token.File.Addr())
}

func (s *Struct) GetGenerics() []Type {
	return s._generics
}

func (s *Struct) SetGenerics(generics []Type) {
	s._generics = generics
}

func (s *Struct) SelfVar(receiver *Var) *Var {
	v := new(Var)
	v.Token = s.Token
	v.DataType = receiver.DataType
	v.DataType.Tag = s
	v.DataType.Id = struct_t
	v.Mutable = receiver.Mutable
	v.Id = lexer.KND_SELF
	return v
}

func (s *Struct) AsTypeKind() string {
	var dts strings.Builder
	dts.WriteString(s.Id)
	if len(s.Generics) > 0 {
		dts.WriteByte('[')
		var gs strings.Builder
		if len(s._generics) > 0 {
			for _, generic := range s.GetGenerics() {
				gs.WriteString(generic.Kind)
				gs.WriteByte(',')
			}
		} else {
			for _, generic := range s.Generics {
				gs.WriteString(generic.Id)
				gs.WriteByte(',')
			}
		}
		dts.WriteString(gs.String()[:gs.Len()-1])
		dts.WriteByte(']')
	}
	return dts.String()
}

func (s *Struct) HasTrait(t *Trait) bool {
	for _, st := range s.Origin.Traits {
		if t == st {
			return true
		}
	}
	return false
}

func (s *Struct) GetSelfRefVarType() Type {
	var t Type
	t.Id = struct_t
	t.Kind = lexer.KND_AMPER + s.Id
	t.Tag = s
	t.Token = s.Token
	return t
}

type Trait struct {
	Pub     bool
	Token   lexer.Token
	Id      string
	Desc    string
	Used    bool
	Funcs   []*Fn
	Defines *Defmap
}

func (t *Trait) FindFunc(id string) *Fn {
	for _, f := range t.Defines.Fns {
		if f.Id == id {
			return f
		}
	}
	return nil
}

func (t *Trait) OutId() string {
	return build.OutId(t.Id, t.Token.File.Addr())
}

type TypeAlias struct {
	Owner      *Block
	Pub        bool
	Token      lexer.Token
	Id         string
	TargetType Type
	Doc        string
	Used       bool
	Generic    bool
}

type Size = int

type TypeSize struct {
	N         Size
	Expr      Expr
	AutoSized bool
}

type Type struct {
	Token         lexer.Token
	Id            uint8
	Original      any
	Kind          string
	MultiTyped    bool
	ComponentType *Type
	Size          TypeSize
	Tag           any
	Pure          bool
	Generic       bool
	CppLinked     bool
}

func (t *Type) InitValue() string {
	if t.Id != enum_t {
		return build.CPP_DEFAULT_EXPR
	}
	return "{" + t.Tag.(*Enum).Items[0].Expr.String() + "}"
}

func (dt *Type) Copy() Type {
	copy := *dt
	if dt.ComponentType != nil {
		copy.ComponentType = new(Type)
		*copy.ComponentType = dt.ComponentType.Copy()
	}
	return copy
}

func (dt *Type) KindWithOriginalId() string {
	if dt.Original == nil {
		return dt.Kind
	}
	_, prefix := dt.KindId()
	original := dt.Original.(Type)
	id, _ := original.KindId()
	return prefix + id
}

func (dt *Type) OriginalKindId() string {
	if dt.Original == nil {
		return ""
	}
	t := dt.Original.(Type)
	id, _ := t.KindId()
	return id
}

func (dt *Type) KindId() (id, prefix string) {
	if dt.Id == map_t || dt.Id == fn_t {
		return dt.Kind, ""
	}
	id = dt.Kind
	runes := []rune(dt.Kind)
	for i, r := range dt.Kind {
		if r == '_' || lexer.IsLetter(r) {
			id = string(runes[i:])
			prefix = string(runes[:i])
			break
		}
	}
	for _, dt := range type_map {
		if dt == id {
			return
		}
	}
	runes = []rune(id)
	for i := 0; i < len(runes); i++ {
		r := runes[i]
		if r == ':' && i+1 < len(runes) && runes[i+1] == ':' {
			i++
			continue
		}
		if r != '_' && !lexer.IsLetter(r) && !lexer.IsDecimal(byte(r)) {
			id = string(runes[:i])
			break
		}
	}
	return
}

func is_necessary_type(id uint8) bool {
	return id == trait_t
}

func (dt *Type) set_to_original_cpp_linked() {
	if dt.Original == nil {
		return
	}
	if dt.Id == struct_t {
		id := dt.Id
		tag := dt.Tag
		*dt = dt.Original.(Type)
		dt.Id = id
		dt.Tag = tag
		return
	}
	*dt = dt.Original.(Type)
}

func (dt *Type) SetToOriginal() {
	if dt.CppLinked {
		dt.set_to_original_cpp_linked()
		return
	} else if dt.Pure || dt.Original == nil {
		return
	}
	kind := dt.KindWithOriginalId()
	id := dt.Id
	tok := dt.Token
	generic := dt.Generic
	*dt = dt.Original.(Type)
	dt.Kind = kind
	dt.Token = tok
	dt.Generic = generic
	if is_necessary_type(id) {
		dt.Id = id
	}
	tag := dt.Tag
	switch tag.(type) {
	case Genericable:
		dt.Tag = tag
	}
}

func (dt *Type) Modifiers() string {
	for i, r := range dt.Kind {
		if r != '*' && r != '&' {
			return dt.Kind[:i]
		}
	}
	return ""
}

func (dt *Type) Pointers() string {
	for i, r := range dt.Kind {
		if r != '*' {
			return dt.Kind[:i]
		}
	}
	return ""
}

func (dt *Type) References() string {
	for i, r := range dt.Kind {
		if r != '&' {
			return dt.Kind[:i]
		}
	}
	return ""
}

func (dt Type) String() (s string) {
	dt.SetToOriginal()
	if dt.MultiTyped {
		return dt.multi_type_str()
	}
	i := strings.LastIndex(dt.Kind, lexer.KND_DBLCOLON)
	if i != -1 {
		dt.Kind = dt.Kind[i+len(lexer.KND_DBLCOLON):]
	}
	modifiers := dt.Modifiers()
	defer func() {
		var cpp strings.Builder
		for _, r := range modifiers {
			if r == '&' {
				cpp.WriteString(build.AsTypeId("ref"))
				cpp.WriteByte('<')
			}
		}
		cpp.WriteString(s)
		for _, r := range modifiers {
			if r == '&' {
				cpp.WriteByte('>')
			}
		}
		for _, r := range modifiers {
			if r == '*' {
				cpp.WriteByte('*')
			}
		}
		s = cpp.String()
	}()
	dt.Kind = dt.Kind[len(modifiers):]
	switch dt.Id {
	case slice_t:
		return dt.slice_str()
	case array_t:
		return dt.array_str()
	case map_t:
		return dt.map_str()
	}
	switch dt.Tag.(type) {
	case *Struct:
		return dt.struct_str()
	}
	switch dt.Id {
	case id_t:
		if dt.CppLinked {
			return dt.Kind
		}
		if dt.Generic {
			return build.AsId(dt.Kind)
		}
		return build.OutId(dt.Kind, dt.Token.File.Addr())
	case enum_t:
		e := dt.Tag.(*Enum)
		return e.DataType.String()
	case trait_t:
		return dt.trait_str()
	case struct_t:
		return dt.struct_str()
	case fn_t:
		return dt.FnString()
	default:
		return cpp_id(dt.Id)
	}
}

func (dt *Type) slice_str() string {
	var cpp strings.Builder
	cpp.WriteString(build.AsTypeId("slice"))
	cpp.WriteByte('<')
	dt.ComponentType.Pure = dt.Pure
	cpp.WriteString(dt.ComponentType.String())
	cpp.WriteByte('>')
	return cpp.String()
}

func (dt *Type) array_str() string {
	var cpp strings.Builder
	cpp.WriteString(build.AsTypeId("array"))
	cpp.WriteByte('<')
	dt.ComponentType.Pure = dt.Pure
	cpp.WriteString(dt.ComponentType.String())
	cpp.WriteByte(',')
	cpp.WriteString(dt.Size.Expr.String())
	cpp.WriteByte('>')
	return cpp.String()
}

func (dt *Type) map_str() string {
	var cpp strings.Builder
	types := dt.Tag.([]Type)
	cpp.WriteString(build.AsTypeId("map"))
	cpp.WriteByte('<')
	key := types[0]
	key.Pure = dt.Pure
	cpp.WriteString(key.String())
	cpp.WriteByte(',')
	value := types[1]
	value.Pure = dt.Pure
	cpp.WriteString(value.String())
	cpp.WriteByte('>')
	return cpp.String()
}

func (dt *Type) trait_str() string {
	var cpp strings.Builder
	id, _ := dt.KindId()
	cpp.WriteString(build.AsTypeId("trait"))
	cpp.WriteByte('<')
	cpp.WriteString(build.OutId(id, dt.Token.File.Addr()))
	cpp.WriteByte('>')
	return cpp.String()
}

func (dt *Type) struct_str() string {
	var cpp strings.Builder
	s := dt.Tag.(*Struct)
	if s.CppLinked && !HasAttribute(build.ATTR_TYPEDEF, s.Attributes) {
		cpp.WriteString("struct ")
	}
	cpp.WriteString(s.OutId())
	types := s.GetGenerics()
	if len(types) == 0 {
		return cpp.String()
	}
	cpp.WriteByte('<')
	for _, t := range types {
		t.Pure = dt.Pure
		cpp.WriteString(t.String())
		cpp.WriteByte(',')
	}
	return cpp.String()[:cpp.Len()-1] + ">"
}

func (dt *Type) FnString() string {
	var cpp strings.Builder
	cpp.WriteString(build.AsTypeId("fn"))
	cpp.WriteByte('<')
	f := dt.Tag.(*Fn)
	cpp.WriteString(f.CppKind(dt.Pure))
	cpp.WriteByte('>')
	return cpp.String()
}

func (dt *Type) multi_type_str() string {
	var cpp strings.Builder
	cpp.WriteString("std::tuple<")
	types := dt.Tag.([]Type)
	for _, t := range types {
		if !t.Pure {
			t.Pure = dt.Pure
		}
		cpp.WriteString(t.String())
		cpp.WriteByte(',')
	}
	return cpp.String()[:cpp.Len()-1] + ">" + dt.Modifiers()
}

func (dt *Type) MapKind() string {
	types := dt.Tag.([]Type)
	var kind strings.Builder
	kind.WriteByte('[')
	kind.WriteString(types[0].Kind)
	kind.WriteByte(':')
	kind.WriteString(types[1].Kind)
	kind.WriteByte(']')
	return kind.String()
}

type UseDecl struct {
	Token      lexer.Token
	Path       string
	Cpp        bool
	LinkString string
	FullUse    bool
	Selectors  []lexer.Token
	Defines    *Defmap
}

type Var struct {
	Owner     *Block
	Public    bool
	Mutable   bool
	Token     lexer.Token
	SetterTok lexer.Token
	Id        string
	DataType  Type
	Expr      Expr
	Constant  bool
	New       bool
	Tag       any
	ExprTag   any
	Doc       string
	Used      bool
	IsField   bool
	CppLinked bool
}

func (v *Var) IsLocal() bool {
	return v.Owner != nil
}

func as_local_id(row, column int, id string) string {
	id = strconv.Itoa(row) + strconv.Itoa(column) + "_" + id
	return build.AsId(id)
}

func (v *Var) OutId() string {
	switch {
	case v.CppLinked:
		return v.Id
	case v.Id == lexer.KND_SELF:
		return "self"
	case v.IsLocal():
		return as_local_id(v.Token.Row, v.Token.Column, v.Id)
	case v.IsField:
		return "_field_" + v.Id
	default:
		return build.OutId(v.Id, v.Token.File.Addr())
	}
}

func (v Var) String() string {
	if lexer.IsIgnoreId(v.Id) {
		return ""
	}
	if v.Constant {
		return ""
	}
	var cpp strings.Builder
	cpp.WriteString(v.DataType.String())
	cpp.WriteByte(' ')
	cpp.WriteString(v.OutId())
	expr := v.Expr.String()
	if expr != "" {
		cpp.WriteString(" = ")
		cpp.WriteString(v.Expr.String())
	} else {
		cpp.WriteString(v.DataType.InitValue())
	}
	cpp.WriteByte(';')
	return cpp.String()
}

func (v *Var) FieldString() string {
	var cpp strings.Builder
	if v.Constant {
		cpp.WriteString("const ")
	}
	cpp.WriteString(v.DataType.String())
	cpp.WriteByte(' ')
	cpp.WriteString(v.OutId())
	cpp.WriteString(v.DataType.InitValue())
	cpp.WriteByte(';')
	return cpp.String()
}

func (v *Var) ReceiverTypeString() string {
	var s strings.Builder
	if v.Mutable {
		s.WriteString("mut ")
	}
	if v.DataType.Kind != "" && v.DataType.Kind[0] == '&' {
		s.WriteByte('&')
	}
	s.WriteString("self")
	return s.String()
}
