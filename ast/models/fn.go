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

	"github.com/DeRuneLabs/jane/lexer"
	"github.com/DeRuneLabs/jane/package/jnapi"
	"github.com/DeRuneLabs/jane/package/jntype"
)

type Fn struct {
	Pub           bool
	IsUnsafe      bool
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
}

func (f *Fn) FindAttribute(kind string) *Attribute {
	for i := range f.Attributes {
		attribute := &f.Attributes[i]
		if attribute.Tag == kind {
			return attribute
		}
	}
	return nil
}

func (f *Fn) plainTypeString() string {
	var s strings.Builder
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
	if f.RetType.Type.MultiTyped {
		s.WriteByte('(')
		for _, t := range f.RetType.Type.Tag.([]Type) {
			s.WriteString(t.Kind)
			s.WriteByte(',')
		}
		return s.String()[:s.Len()-1] + ")"
	} else if f.RetType.Type.Id != jntype.Void {
		s.WriteString(f.RetType.Type.Kind)
	}
	return s.String()
}

func (f *Fn) DataTypeString() string {
	var cpp strings.Builder
	if f.IsUnsafe {
		cpp.WriteString("unsafe ")
	}
	cpp.WriteString("fn")
	cpp.WriteString(f.plainTypeString())
	return cpp.String()
}

func (f *Fn) OutId() string {
	if f.Receiver != nil {
		return f.Id
	}
	return jnapi.OutId(f.Id, f.Token.File)
}

func (f *Fn) DefString() string {
	var s strings.Builder
	if f.IsUnsafe {
		s.WriteString("unsafe ")
	}
	s.WriteString("fn ")
	s.WriteString(f.Id)
	s.WriteString(f.plainTypeString())
	return s.String()
}
