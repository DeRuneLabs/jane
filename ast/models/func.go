package models

import (
	"strings"

	"github.com/DeRuneLabs/jane/package/jnapi"
	"github.com/DeRuneLabs/jane/package/jntype"
)

type Func struct {
	Pub        bool
	Tok        Tok
	Id         string
	Generics   []*GenericType
	Combines   *[][]DataType
	Attributes []Attribute
	Params     []Param
	RetType    RetType
	Block      *Block
	Receiver   *DataType
	Owner      any
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
			cxx.WriteString(p.Type.Kind)
			cxx.WriteString(", ")
		}
		cxxStr := cxx.String()[:cxx.Len()-2]
		cxx.Reset()
		cxx.WriteString(cxxStr)
	}
	cxx.WriteByte(')')
	if f.RetType.Type.MultiTyped {
		cxx.WriteByte('[')
		for _, t := range f.RetType.Type.Tag.([]DataType) {
			cxx.WriteString(t.Kind)
			cxx.WriteByte(',')
		}
		return cxx.String()[:cxx.Len()-1] + "]"
	} else if f.RetType.Type.Id != jntype.Void {
		cxx.WriteString(f.RetType.Type.Kind)
	}
	return cxx.String()
}

func (f *Func) OutId() string {
	if f.Receiver != nil {
		return f.Id
	}
	return jnapi.OutId(f.Id, f.Tok.File)
}

func (f *Func) DefString() string {
	return f.Id + f.DataTypeString()
}
