// Copyright (c) 2024 - DeRuneLabs
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

package models

import (
	"strings"
	"unicode"

	"github.com/DeRuneLabs/jane/lexer"
	"github.com/DeRuneLabs/jane/lexer/tokens"
	"github.com/DeRuneLabs/jane/package/jnapi"
	"github.com/DeRuneLabs/jane/package/jntype"
)

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
	if dt.Id == jntype.Map || dt.Id == jntype.Fn {
		return dt.Kind, ""
	}
	id = dt.Kind
	runes := []rune(dt.Kind)
	for i, r := range dt.Kind {
		if r == '_' || unicode.IsLetter(r) {
			id = string(runes[i:])
			prefix = string(runes[:i])
			break
		}
	}
	for _, dt := range jntype.TypeMap {
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
		if r != '_' && !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			id = string(runes[:i])
			break
		}
	}
	return
}

func is_necessary_type(id uint8) bool {
	return id == jntype.Trait
}

func (dt *Type) SetToOriginal() {
	if dt.Pure || dt.Original == nil {
		return
	}
	tag := dt.Tag
	switch tag.(type) {
	case Genericable:
		defer func() { dt.Tag = tag }()
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
		return dt.MultiTypeString()
	}
	i := strings.LastIndex(dt.Kind, tokens.DOUBLE_COLON)
	if i != -1 {
		dt.Kind = dt.Kind[i+len(tokens.DOUBLE_COLON):]
	}
	modifiers := dt.Modifiers()
	defer func() {
		var cpp strings.Builder
		for _, r := range modifiers {
			if r == '&' {
				cpp.WriteString("jnc_ref<")
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
	case jntype.Slice:
		return dt.SliceString()
	case jntype.Array:
		return dt.ArrayString()
	case jntype.Map:
		return dt.MapString()
	}
	switch dt.Tag.(type) {
	case CompiledStruct:
		return dt.StructString()
	}
	switch dt.Id {
	case jntype.Id:
		if dt.CppLinked {
			return dt.Kind
		}
		if dt.Generic {
			return jnapi.AsId(dt.Kind)
		}
		return jnapi.OutId(dt.Kind, dt.Token.File)
	case jntype.Enum:
		e := dt.Tag.(*Enum)
		return e.Type.String()
	case jntype.Trait:
		return dt.TraitString()
	case jntype.Struct:
		return dt.StructString()
	case jntype.Fn:
		return dt.FuncString()
	default:
		return jntype.CppId(dt.Id)
	}
}

func (dt *Type) SliceString() string {
	var cpp strings.Builder
	cpp.WriteString("slice<")
	dt.ComponentType.Pure = dt.Pure
	cpp.WriteString(dt.ComponentType.String())
	cpp.WriteByte('>')
	return cpp.String()
}

func (dt *Type) ArrayString() string {
	var cpp strings.Builder
	cpp.WriteString("array<")
	dt.ComponentType.Pure = dt.Pure
	cpp.WriteString(dt.ComponentType.String())
	cpp.WriteByte(',')
	cpp.WriteString(dt.Size.Expr.String())
	cpp.WriteByte('>')
	return cpp.String()
}

func (dt *Type) MapString() string {
	var cpp strings.Builder
	types := dt.Tag.([]Type)
	cpp.WriteString("map<")
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

func (dt *Type) TraitString() string {
	var cpp strings.Builder
	id, _ := dt.KindId()
	cpp.WriteString("trait<")
	cpp.WriteString(jnapi.OutId(id, dt.Token.File))
	cpp.WriteByte('>')
	return cpp.String()
}

func (dt *Type) StructString() string {
	var cpp strings.Builder
	s := dt.Tag.(CompiledStruct)
	if s.CppLinked() {
		cpp.WriteString("struct ")
	}
	cpp.WriteString(s.OutId())
	types := s.Generics()
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

func (dt *Type) FuncString() string {
	var cpp strings.Builder
	cpp.WriteString("fn<std::function<")
	f := dt.Tag.(*Fn)
	f.RetType.Type.Pure = dt.Pure
	cpp.WriteString(f.RetType.String())
	cpp.WriteByte('(')
	if len(f.Params) > 0 {
		for _, param := range f.Params {
			param.Type.Pure = dt.Pure
			cpp.WriteString(param.Prototype())
			cpp.WriteByte(',')
		}
		cppStr := cpp.String()[:cpp.Len()-1]
		cpp.Reset()
		cpp.WriteString(cppStr)
	} else {
		cpp.WriteString("void")
	}
	cpp.WriteString(")>>")
	return cpp.String()
}

func (dt *Type) MultiTypeString() string {
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
