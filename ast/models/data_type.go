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
		var cpp strings.Builder
		for range pointers {
			cpp.WriteString("ptr<")
		}
		cpp.WriteString(s)
		for range pointers {
			cpp.WriteString(">")
		}
		s = cpp.String()
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
		return jntype.CppId(dt.Id)
	}
}

func (dt DataType) SliceString() string {
	var cpp strings.Builder
	cpp.WriteString("slice<")
	dt.Kind = dt.Kind[len(jn.Prefix_Slice):]
	cpp.WriteString(dt.String())
	cpp.WriteByte('>')
	return cpp.String()
}

func (dt DataType) ArrayComponent() DataType {
	dt.Kind = dt.Kind[len(jn.Prefix_Array):]
	exprs := dt.Tag.([][]any)[1:]
	dt.Tag = exprs
	return dt
}

func (dt DataType) ArrayString() string {
	var cpp strings.Builder
	cpp.WriteString("array<")
	exprs := dt.Tag.([][]any)
	expr := exprs[0][1].(Expr)
	cpp.WriteString(dt.ArrayComponent().String())
	cpp.WriteByte(',')
	cpp.WriteString(expr.String())
	cpp.WriteByte('>')
	return cpp.String()
}

func (dt *DataType) MapString() string {
	var cpp strings.Builder
	types := dt.Tag.([]DataType)
	cpp.WriteString("map<")
	key := types[0]
	key.DontUseOriginal = dt.DontUseOriginal
	cpp.WriteString(key.String())
	cpp.WriteByte(',')
	value := types[1]
	value.DontUseOriginal = dt.DontUseOriginal
	cpp.WriteString(value.String())
	cpp.WriteByte('>')
	return cpp.String()
}

func (dt *DataType) TraitString() string {
	var cpp strings.Builder
	id, _ := dt.KindId()
	cpp.WriteString("trait<")
	cpp.WriteString(jnapi.OutId(id, dt.Tok.File))
	cpp.WriteByte('>')
	return cpp.String()
}

func (dt *DataType) StructString() string {
	var cpp strings.Builder
	s := dt.Tag.(CompiledStruct)
	cpp.WriteString(s.OutId())
	types := s.Generics()
	if len(types) == 0 {
		return cpp.String()
	}
	cpp.WriteByte('<')
	for _, t := range types {
		t.DontUseOriginal = dt.DontUseOriginal
		cpp.WriteString(t.String())
		cpp.WriteByte(',')
	}
	return cpp.String()[:cpp.Len()-1] + ">"
}

func (dt *DataType) FuncString() string {
	var cpp strings.Builder
	cpp.WriteString("std::function<")
	f := dt.Tag.(*Func)
	f.RetType.Type.DontUseOriginal = dt.DontUseOriginal
	cpp.WriteString(f.RetType.String())
	cpp.WriteByte('(')
	if len(f.Params) > 0 {
		for _, param := range f.Params {
			param.Type.DontUseOriginal = dt.DontUseOriginal
			cpp.WriteString(param.Prototype())
			cpp.WriteByte(',')
		}
		cxxStr := cpp.String()[:cpp.Len()-1]
		cpp.Reset()
		cpp.WriteString(cxxStr)
	} else {
		cpp.WriteString("void")
	}
	cpp.WriteString(")>")
	return cpp.String()
}

func (dt *DataType) MultiTypeString() string {
	types := dt.Tag.([]DataType)
	var cpp strings.Builder
	cpp.WriteString("std::tuple<")
	for _, t := range types {
		t.DontUseOriginal = dt.DontUseOriginal
		cpp.WriteString(t.String())
		cpp.WriteByte(',')
	}
	return cpp.String()[:cpp.Len()-1] + ">" + dt.Pointers()
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
