package parser

import (
	"strconv"
	"strings"

	"github.com/DeRuneLabs/jane/ast"
	"github.com/DeRuneLabs/jane/package/jnapi"
)

type jnstruct struct {
	Ast         Struct
	Defs        *Defmap
	Used        bool
	Desc        string
	constructor *Func
	generics    []DataType
}

func (s *jnstruct) declString() string {
	var cxx strings.Builder
	cxx.WriteString(genericsToCxx(s.Ast.Generics))
	cxx.WriteByte('\n')
	cxx.WriteString("struct ")
	cxx.WriteString(jnapi.OutId(s.Ast.Id, s.Ast.Tok.File))
	cxx.WriteString(" {\n")
	ast.AddIndent()
	for _, g := range s.Defs.Globals {
		cxx.WriteString(ast.IndentString())
		cxx.WriteString(g.FieldString())
		cxx.WriteByte('\n')
	}
	ast.DoneIndent()
	cxx.WriteString(ast.IndentString())
	cxx.WriteString("};")
	return cxx.String()
}

func (s *jnstruct) ostreams() string {
	var cxx strings.Builder
	var generics string
	if len(s.Ast.Generics) > 0 {
		var gb strings.Builder
		gb.WriteByte('<')
		for i := range s.Ast.Generics {
			gb.WriteByte('T')
			gb.WriteString(strconv.Itoa(i))
			gb.WriteByte(',')
		}
		generics = gb.String()[:gb.Len()-1] + ">"
		cxx.WriteString("template<typename ")
		cxx.WriteString(generics[1:])
		cxx.WriteByte('\n')
	}
	cxx.WriteString("std::ostream &operator<<(std::ostream &_Stream, const ")
	cxx.WriteString(jnapi.OutId(s.Ast.Id, s.Ast.Tok.File))
	cxx.WriteString(generics)
	cxx.WriteString(" &_Src) {\n")
	ast.AddIndent()
	cxx.WriteString(ast.IndentString())
	cxx.WriteString(`_Stream << "`)
	cxx.WriteString(s.Ast.Id)
	cxx.WriteString("{\";\n")
	for i, field := range s.Ast.Fields {
		cxx.WriteString(ast.IndentString())
		cxx.WriteString(`_Stream << "`)
		cxx.WriteString(field.Id)
		cxx.WriteString(`:" << _Src.`)
		cxx.WriteString(jnapi.OutId(field.Id, s.Ast.Tok.File))
		if i+1 < len(s.Ast.Fields) {
			cxx.WriteString(" << \", \"")
		}
		cxx.WriteString(";\n")
	}
	cxx.WriteString(ast.IndentString())
	cxx.WriteString("_Stream << \"}\";\n")
	cxx.WriteString(ast.IndentString())
	cxx.WriteString("return _Stream;\n")
	ast.DoneIndent()
	cxx.WriteString(ast.IndentString())
	cxx.WriteString("}")
	return cxx.String()
}

func (s jnstruct) String() string {
	var cxx strings.Builder
	cxx.WriteString(s.declString())
	cxx.WriteString("\n\n")
	cxx.WriteString(ast.IndentString())
	cxx.WriteString(s.ostreams())
	return cxx.String()
}

func (s *jnstruct) Generics() []DataType {
	return s.generics
}

func (s *jnstruct) SetGenerics(generics []DataType) {
	s.generics = generics
}

func (s *jnstruct) dataTypeString() string {
	var dts strings.Builder
	dts.WriteString(s.Ast.Id)
	if len(s.Ast.Generics) > 0 {
		dts.WriteByte('[')
		var gs strings.Builder
		if len(s.generics) > 0 {
			for _, generic := range s.generics {
				gs.WriteString(generic.String())
				gs.WriteByte(',')
			}
		} else {
			for _, generic := range s.Ast.Generics {
				gs.WriteString(generic.Id)
				gs.WriteByte(',')
			}
		}
		dts.WriteString(gs.String()[:gs.Len()-1])
		dts.WriteByte(']')
	}
	return dts.String()
}
