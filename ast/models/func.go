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
	var cpp strings.Builder
	cpp.WriteByte('(')
	if len(f.Params) > 0 {
		for _, p := range f.Params {
			if p.Variadic {
				cpp.WriteString("...")
			}
			cpp.WriteString(p.Type.Kind)
			cpp.WriteString(", ")
		}
		cxxStr := cpp.String()[:cpp.Len()-2]
		cpp.Reset()
		cpp.WriteString(cxxStr)
	}
	cpp.WriteByte(')')
	if f.RetType.Type.MultiTyped {
		cpp.WriteByte('[')
		for _, t := range f.RetType.Type.Tag.([]DataType) {
			cpp.WriteString(t.Kind)
			cpp.WriteByte(',')
		}
		return cpp.String()[:cpp.Len()-1] + "]"
	} else if f.RetType.Type.Id != jntype.Void {
		cpp.WriteString(f.RetType.Type.Kind)
	}
	return cpp.String()
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
