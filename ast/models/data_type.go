package models

import (
	"strings"
	"unicode"

	"github.com/DeRuneLabs/jane/lexer/tokens"
	"github.com/DeRuneLabs/jane/package/jn"
	"github.com/DeRuneLabs/jane/package/jnapi"
	"github.com/DeRuneLabs/jane/package/jntype"
)

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
	id, _ := original.KindId()
	return prefix + id
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
	if dt.Id == jntype.Map || dt.Id == jntype.Func {
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

func (dt *DataType) SetToOriginal() {
	if dt.DontUseOriginal || dt.Original == nil {
		return
	}
	tag := dt.Tag
	switch tag.(type) {
	case Genericable:
		defer func() { dt.Tag = tag }()
	}
	kind := dt.KindWithOriginalId()
	tok := dt.Tok
	*dt = dt.Original.(DataType)
	dt.Kind = kind
	dt.Tok = tok
	if strings.HasPrefix(dt.Kind, jn.Prefix_Array) {
		dt.Tag = tag
	}
}

func (dt *DataType) Pointers() string {
	for i, run := range dt.Kind {
		if run != '*' {
			return dt.Kind[:i]
		}
	}
	return ""
}

func (dt DataType) String() (s string) {
	dt.SetToOriginal()
	if dt.MultiTyped {
		return dt.MultiTypeString()
	}
	i := strings.LastIndex(dt.Kind, tokens.DOUBLE_COLON)
	if i != -1 {
		dt.Kind = dt.Kind[i+len(tokens.DOUBLE_COLON):]
	}
	pointers := dt.Pointers()

	defer func() {
		var cxx strings.Builder
		for range pointers {
			cxx.WriteString("ptr<")
		}
		cxx.WriteString(s)
		for range pointers {
			cxx.WriteString(">")
		}
		s = cxx.String()
	}()
	dt.Kind = dt.Kind[len(pointers):]
	if dt.Kind != "" {
		switch {
		case strings.HasPrefix(dt.Kind, jn.Prefix_Slice):
			return dt.SliceString()
		case strings.HasPrefix(dt.Kind, jn.Prefix_Array):
			return dt.ArrayString()
		case dt.Id == jntype.Map && dt.Kind[0] == '[' && dt.Kind[len(dt.Kind)-1] == ']':
			return dt.MapString()
		}
	}
	switch dt.Tag.(type) {
	case CompiledStruct:
		return dt.StructString()
	}
	switch dt.Id {
	case jntype.Id, jntype.Enum:
		return jnapi.OutId(dt.Kind, dt.Tok.File)
	case jntype.Trait:
		return dt.TraitString()
	case jntype.Struct:
		return dt.StructString()
	case jntype.Func:
		return dt.FuncString()
	default:
		return jntype.CxxId(dt.Id)
	}
}

func (dt DataType) SliceString() string {
	var cxx strings.Builder
	cxx.WriteString("slice<")
	dt.Kind = dt.Kind[len(jn.Prefix_Slice):]
	cxx.WriteString(dt.String())
	cxx.WriteByte('>')
	return cxx.String()
}

func (dt DataType) ArrayComponent() DataType {
	dt.Kind = dt.Kind[len(jn.Prefix_Array):]
	exprs := dt.Tag.([][]any)[1:]
	dt.Tag = exprs
	return dt
}

func (dt DataType) ArrayString() string {
	var cxx strings.Builder
	cxx.WriteString("array<")
	exprs := dt.Tag.([][]any)
	expr := exprs[0][1].(Expr)
	cxx.WriteString(dt.ArrayComponent().String())
	cxx.WriteByte(',')
	cxx.WriteString(expr.String())
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

func (dt *DataType) TraitString() string {
	var cxx strings.Builder
	id, _ := dt.KindId()
	cxx.WriteString("trait<")
	cxx.WriteString(jnapi.OutId(id, dt.Tok.File))
	cxx.WriteByte('>')
	return cxx.String()
}

func (dt *DataType) StructString() string {
	var cxx strings.Builder
	s := dt.Tag.(CompiledStruct)
	cxx.WriteString(s.OutId())
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
	cxx.WriteString("std::function<")
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

func (dt *DataType) MapKind() string {
	types := dt.Tag.([]DataType)
	var kind strings.Builder
	kind.WriteByte('[')
	kind.WriteString(types[0].Kind)
	kind.WriteByte(':')
	kind.WriteString(types[1].Kind)
	kind.WriteByte(']')
	return kind.String()
}
