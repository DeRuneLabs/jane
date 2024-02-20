package parser

import (
	"fmt"

	"github.com/De-Rune/jane/ast"
	"github.com/De-Rune/jane/lexer"
	"github.com/De-Rune/jane/package/jane"
)

const entryPointStandard = `
  setlocale(0x0, "");
`

type Function struct {
	Token      lexer.Token
	Name       string
	ReturnType uint8
	Params     []ast.ParameterAST
	Block      ast.BlockAST
}

func (f Function) String() string {
	code := ""
	code += jane.CxxTypeNameFromType(f.ReturnType)
	code += " "
	code += f.Name
	code += "("
	if len(f.Params) > 0 {
		any := false
		for _, p := range f.Params {
			code += p.String()
			code += ","
			if !any {
				any = p.Type.Type == jane.Any
			}
		}
		code = code[:len(code)-1]
		if any {
			code = "templat <typename any>\n" + code
		}
	}
	code += ") {"
	code += getFunctionStandardCode(f.Name)
	for _, s := range f.Block.Content {
		code += "\n"
		code += " " + fmt.Sprint(s.Value)
		code += ";"
	}
	code += "\n}"
	return code
}

func getFunctionStandardCode(name string) string {
	switch name {
	case jane.EntryPoint:
		return entryPointStandard
	}
	return ""
}
