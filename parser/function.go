package parser

import (
	"strings"

	"github.com/DeRuneLabs/jane/ast/models"
	"github.com/DeRuneLabs/jane/package/jn"
	"github.com/DeRuneLabs/jane/package/jnapi"
)

type function struct {
	Ast     *Func
	Desc    string
	used    bool
	checked bool
}

func (f *function) outId() string {
	if f.Ast.Id == jn.EntryPoint {
		return jnapi.OutId(f.Ast.Id, nil)
	}
	return jnapi.OutId(f.Ast.Id, f.Ast.Tok.File)
}

func (f function) String() string {
	var cxx strings.Builder
	cxx.WriteString(f.Head())
	cxx.WriteByte(' ')
	block := f.Ast.Block
	vars := f.Ast.RetType.Vars()
	if vars != nil {
		statements := make([]models.Statement, len(vars))
		for i, v := range vars {
			statements[i] = models.Statement{
				Tok: v.IdTok,
				Val: *v,
			}
		}
		block.Tree = append(statements, block.Tree...)
	}
	cxx.WriteString(block.String())
	return cxx.String()
}

func (f *function) Head() string {
	var cxx strings.Builder
	cxx.WriteString(f.declHead())
	cxx.WriteString(paramsToCxx(f.Ast.Params))
	return cxx.String()
}

func (f *function) declHead() string {
	var cxx strings.Builder
	cxx.WriteString(genericsToCxx(f.Ast.Generics))
	if cxx.Len() > 0 {
		cxx.WriteByte('\n')
	}
	cxx.WriteString(attributesToString(f.Ast.Attributes))
	cxx.WriteString(f.Ast.RetType.String())
	cxx.WriteByte(' ')
	cxx.WriteString(f.outId())
	return cxx.String()
}

func (f *function) Prototype() string {
	var cxx strings.Builder
	cxx.WriteString(f.declHead())
	cxx.WriteString(f.PrototypeParams())
	cxx.WriteByte(';')
	return cxx.String()
}

func (f *function) PrototypeParams() string {
	if len(f.Ast.Params) == 0 {
		return "(void)"
	}
	var cxx strings.Builder
	cxx.WriteByte('(')
	for _, p := range f.Ast.Params {
		cxx.WriteString(p.Prototype())
		cxx.WriteByte(',')
	}
	return cxx.String()[:cxx.Len()-1] + ")"
}

func isOutableAttribute(kind string) bool {
	return kind == jn.Attribute_Inline
}

func attributesToString(attributes []Attribute) string {
	var cxx strings.Builder
	for _, attr := range attributes {
		if isOutableAttribute(attr.Tag.Kind) {
			cxx.WriteString(attr.String())
			cxx.WriteByte(' ')
		}
	}
	return cxx.String()
}

func paramsToCxx(params []Param) string {
	if len(params) == 0 {
		return "(void)"
	}
	var cxx strings.Builder
	cxx.WriteByte('(')
	for _, p := range params {
		cxx.WriteString(p.String())
		cxx.WriteByte(',')
	}
	return cxx.String()[:cxx.Len()-1] + ")"
}
