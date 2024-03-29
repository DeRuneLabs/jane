package parser

import (
	"strconv"
	"strings"

	"github.com/DeRuneLabs/jane/ast/models"
	"github.com/DeRuneLabs/jane/lexer/tokens"
	"github.com/DeRuneLabs/jane/package/jnapi"
	"github.com/DeRuneLabs/jane/package/jntype"
)

type jnstruct struct {
	Ast         Struct
	Defs        *Defmap
	Used        bool
	Desc        string
	constructor *Func
	traits      []*trait
	generics    []DataType
}

func (s *jnstruct) hasTrait(t *trait) bool {
	for _, st := range s.traits {
		if t == st {
			return true
		}
	}
	return false
}

func (s *jnstruct) cxxGenerics() (def string, serie string) {
	if len(s.Ast.Generics) == 0 {
		return "", ""
	}
	var cxxDef strings.Builder
	cxxDef.WriteString("template<typename ")
	var cxxSerie strings.Builder
	cxxSerie.WriteByte('<')
	for i := range s.Ast.Generics {
		cxxSerie.WriteByte('T')
		cxxSerie.WriteString(strconv.Itoa(i))
		cxxSerie.WriteByte(',')
	}
	serie = cxxSerie.String()[:cxxSerie.Len()-1] + ">"
	cxxDef.WriteString(serie[1:])
	cxxDef.WriteByte('\n')
	return cxxDef.String(), serie
}

// OutId returns jnapi.OutId of struct.
//
// This function is should be have this function
// for CompiledStruct interface of ast package.
func (s *jnstruct) OutId() string {
	return jnapi.OutId(s.Ast.Id, s.Ast.Tok.File)
}

func (s *jnstruct) operators() string {
	outid := s.OutId()
	genericsDef, genericsSerie := s.cxxGenerics()
	var cxx strings.Builder
	cxx.WriteString(models.IndentString())
	if l, _ := cxx.WriteString(genericsDef); l > 0 {
		cxx.WriteString(models.IndentString())
	}
	cxx.WriteString("inline bool operator==(const ")
	cxx.WriteString(outid)
	cxx.WriteString(genericsSerie)
	cxx.WriteString(" &_Src) {")
	if len(s.Defs.Globals) > 0 {
		models.AddIndent()
		cxx.WriteByte('\n')
		cxx.WriteString(models.IndentString())
		var expr strings.Builder
		expr.WriteString("return ")
		models.AddIndent()
		for _, g := range s.Defs.Globals {
			expr.WriteByte('\n')
			expr.WriteString(models.IndentString())
			expr.WriteString("this->")
			gid := g.OutId()
			expr.WriteString(gid)
			expr.WriteString(" == _Src.")
			expr.WriteString(gid)
			expr.WriteString(" &&")
		}
		models.DoneIndent()
		cxx.WriteString(expr.String()[:expr.Len()-3])
		cxx.WriteString(";\n")
		models.DoneIndent()
		cxx.WriteString(models.IndentString())
		cxx.WriteByte('}')
	} else {
		cxx.WriteString(" return true; }")
	}
	cxx.WriteString("\n\n")
	cxx.WriteString(models.IndentString())
	if l, _ := cxx.WriteString(genericsDef); l > 0 {
		cxx.WriteString(models.IndentString())
	}
	cxx.WriteString("inline bool operator!=(const ")
	cxx.WriteString(outid)
	cxx.WriteString(genericsSerie)
	cxx.WriteString(" &_Src) { return !this->operator==(_Src); }")
	return cxx.String()
}

func (s *jnstruct) cxxConstructor() string {
	var cxx strings.Builder
	cxx.WriteString(models.IndentString())
	cxx.WriteString(s.OutId())
	cxx.WriteString(paramsToCxx(s.constructor.Params))
	cxx.WriteString(" noexcept {")
	if len(s.Defs.Globals) > 0 {
		models.AddIndent()
		for i, g := range s.Defs.Globals {
			cxx.WriteByte('\n')
			cxx.WriteString(models.IndentString())
			cxx.WriteString(g.OutId())
			cxx.WriteString(" = ")
			cxx.WriteString(s.constructor.Params[i].OutId())
			cxx.WriteByte(';')
		}
		models.DoneIndent()
		cxx.WriteByte('\n')
	}
	cxx.WriteString(models.IndentString())
	cxx.WriteByte('}')
	return cxx.String()
}

func (s *jnstruct) cxxTraits() string {
	if len(s.traits) == 0 {
		return ""
	}
	var cxx strings.Builder
	cxx.WriteString(": ")
	for _, t := range s.traits {
		cxx.WriteString("public ")
		cxx.WriteString(t.OutId())
		cxx.WriteByte(',')
	}
	return cxx.String()[:cxx.Len()-1]
}

func (s *jnstruct) prototype() string {
	var cxx strings.Builder
	cxx.WriteString(genericsToCxx(s.Ast.Generics))
	cxx.WriteString(" struct ")
	cxx.WriteString(s.OutId())
	cxx.WriteByte(';')
	return cxx.String()
}

func (s *jnstruct) decldefString() string {
	var cxx strings.Builder
	cxx.WriteString(genericsToCxx(s.Ast.Generics))
	cxx.WriteByte('\n')
	cxx.WriteString("struct ")
	cxx.WriteString(s.OutId())
	cxx.WriteString(s.cxxTraits())
	cxx.WriteString(" {\n")
	models.AddIndent()
	if len(s.Defs.Globals) > 0 {
		for _, g := range s.Defs.Globals {
			cxx.WriteString(models.IndentString())
			cxx.WriteString(g.FieldString())
			cxx.WriteByte('\n')
		}
		cxx.WriteString("\n\n")
		cxx.WriteString(s.cxxConstructor())
		cxx.WriteString("\n\n")
	}
	cxx.WriteString(models.IndentString())
	cxx.WriteString(s.OutId())
	cxx.WriteString("(void) noexcept {}\n\n")
	for _, f := range s.Defs.Funcs {
		if f.used {
			cxx.WriteString(models.IndentString())
			cxx.WriteString(f.String())
			cxx.WriteString("\n\n")
		}
	}
	cxx.WriteString(s.operators())
	cxx.WriteByte('\n')
	models.DoneIndent()
	cxx.WriteString(models.IndentString())
	cxx.WriteString("};")
	return cxx.String()
}

func (s *jnstruct) ostream() string {
	var cxx strings.Builder
	genericsDef, genericsSerie := s.cxxGenerics()
	cxx.WriteString(models.IndentString())
	if l, _ := cxx.WriteString(genericsDef); l > 0 {
		cxx.WriteString(models.IndentString())
	}
	cxx.WriteString("std::ostream &operator<<(std::ostream &_Stream, const ")
	cxx.WriteString(s.OutId())
	cxx.WriteString(genericsSerie)
	cxx.WriteString(" &_Src) {\n")
	models.AddIndent()
	cxx.WriteString(models.IndentString())
	cxx.WriteString(`_Stream << "`)
	cxx.WriteString(s.Ast.Id)
	cxx.WriteString("{\";\n")
	for i, field := range s.Ast.Fields {
		cxx.WriteString(models.IndentString())
		cxx.WriteString(`_Stream << "`)
		cxx.WriteString(field.Id)
		cxx.WriteString(`:" << _Src.`)
		cxx.WriteString(field.OutId())
		if i+1 < len(s.Ast.Fields) {
			cxx.WriteString(" << \", \"")
		}
		cxx.WriteString(";\n")
	}
	cxx.WriteString(models.IndentString())
	cxx.WriteString("_Stream << \"}\";\n")
	cxx.WriteString(models.IndentString())
	cxx.WriteString("return _Stream;\n")
	models.DoneIndent()
	cxx.WriteString(models.IndentString())
	cxx.WriteString("}")
	return cxx.String()
}

func (s jnstruct) String() string {
	var cxx strings.Builder
	cxx.WriteString(s.decldefString())
	cxx.WriteString("\n\n")
	cxx.WriteString(s.ostream())
	return cxx.String()
}

func (s *jnstruct) Generics() []DataType {
	return s.generics
}

func (s *jnstruct) SetGenerics(generics []DataType) {
	s.generics = generics
}

func (s *jnstruct) selfVar(receiver DataType) *Var {
	v := new(models.Var)
	v.IdTok = s.Ast.Tok
	v.Type = receiver
	v.Type.Id = jntype.Struct
	v.Id = tokens.SELF
	if typeIsPtr(receiver) {
		v.Expr.Model = exprNode{jnapi.CxxSelf}
	} else {
		v.Expr.Model = exprNode{tokens.STAR + jnapi.CxxSelf}
	}
	return v
}

func (s *jnstruct) dataTypeString() string {
	var dts strings.Builder
	dts.WriteString(s.Ast.Id)
	if len(s.Ast.Generics) > 0 {
		dts.WriteByte('[')
		var gs strings.Builder
		if len(s.generics) > 0 {
			for _, generic := range s.generics {
				gs.WriteString(generic.Kind)
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
