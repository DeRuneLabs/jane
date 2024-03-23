package models

import (
	"strings"

	"github.com/DeRuneLabs/jane/package/jntype"
)

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
