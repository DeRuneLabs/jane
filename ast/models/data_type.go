package models

import (
	"strings"
	"unicode"

	"github.com/DeRuneLabs/jane/package/jnapi"
	"github.com/DeRuneLabs/jane/package/jntype"
)

type genericableTypes struct {
	types []DataType
}

func (gt genericableTypes) Generics() []DataType {
	return gt.types
}

func (gt genericableTypes) SetGenerics([]DataType) {}

type DataType struct {
	Tok             Tok
	Id              uint8
	Original        any
	Kind            string
	MultiTyped      bool
	Tag             any
	DontUseOriginal bool
}

func (dt *DataType) KindWithOriginalId() string {
	if dt.Original == nil {
		return dt.Kind
	}
	_, prefix := dt.KindId()
	original := dt.Original.(DataType)
	return prefix + original.Tok.Kind
}

func (dt *DataType) OriginalKindId() string {
	if dt.Original == nil {
		return ""
	}
	t := dt.Original.(DataType)
	id, _ := t.KindId()
	return id
}

func (dt *DataType) KindId() (id, prefix string) {
	id = dt.Kind
	runes := []rune(dt.Kind)
	for i, r := range dt.Kind {
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

func (dt *DataType) SetToOriginal() {
	if dt.DontUseOriginal || dt.Original == nil {
		return
	}
	val := dt.KindWithOriginalId()
	tok := dt.Tok
	*dt = dt.Original.(DataType)
	dt.Kind = val
	dt.Tok = tok
}

func (dt *DataType) Pointers() string {
	for i, run := range dt.Kind {
		if run != '*' {
			return dt.Kind[:i]
		}
	}
	return ""
}

func (dt DataType) String() string {
	dt.SetToOriginal()
	if dt.MultiTyped {
		return dt.MultiTypeString()
	}
	pointers := dt.Pointers()
	dt.Kind = dt.Kind[len(pointers):]
	if dt.Kind != "" {
		switch {
		case strings.HasPrefix(dt.Kind, "[]"):
			return dt.ArrayString() + pointers
		case dt.Id == jntype.Map && dt.Kind[0] == '[':
			return dt.MapString() + pointers
		}
	}
	if dt.Tag != nil {
		switch t := dt.Tag.(type) {
		case []DataType:
			dt.Tag = genericableTypes{t}
			return dt.StructString()
		case Genericable:
			return dt.StructString()
		}
	}
	switch dt.Id {
	case jntype.Id, jntype.Enum:
		return jnapi.OutId(dt.Kind, dt.Tok.File) + pointers
	case jntype.Struct:
		return dt.StructString() + pointers
	case jntype.Func:
		return dt.FuncString() + pointers
	default:
		return jntype.CxxTypeIdFromType(dt.Id) + pointers
	}
}

func (dt DataType) ArrayString() string {
	var cxx strings.Builder
	cxx.WriteString("array<")
	dt.Kind = dt.Kind[2:]
	cxx.WriteString(dt.String())
	cxx.WriteByte('>')
	return cxx.String()
}

func (dt *DataType) MapString() string {
	var cxx strings.Builder
	types := dt.Tag.([]DataType)
	cxx.WriteString("map<")
	key := types[0]
	key.DontUseOriginal = dt.DontUseOriginal
	cxx.WriteString(key.String())
	cxx.WriteByte(',')
	value := types[1]
	value.DontUseOriginal = dt.DontUseOriginal
	cxx.WriteString(value.String())
	cxx.WriteByte('>')
	return cxx.String()
}

func (dt *DataType) StructString() string {
	var cxx strings.Builder
	id, _ := dt.KindId()
	cxx.WriteString(jnapi.OutId(id, dt.Tok.File))
	s := dt.Tag.(Genericable)
	types := s.Generics()
	if len(types) == 0 {
		return cxx.String()
	}
	cxx.WriteByte('<')
	for _, t := range types {
		t.DontUseOriginal = dt.DontUseOriginal
		cxx.WriteString(t.String())
		cxx.WriteByte(',')
	}
	return cxx.String()[:cxx.Len()-1] + ">"
}

func (dt *DataType) FuncString() string {
	var cxx strings.Builder
	cxx.WriteString("func<")
	f := dt.Tag.(*Func)
	f.RetType.Type.DontUseOriginal = dt.DontUseOriginal
	cxx.WriteString(f.RetType.String())
	cxx.WriteByte('(')
	if len(f.Params) > 0 {
		for _, param := range f.Params {
			param.Type.DontUseOriginal = dt.DontUseOriginal
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
		t.DontUseOriginal = dt.DontUseOriginal
		cxx.WriteString(t.String())
		cxx.WriteByte(',')
	}
	return cxx.String()[:cxx.Len()-1] + ">" + dt.Pointers()
}
