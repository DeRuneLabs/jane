package parser

import (
	"github.com/De-Rune/jane/ast"
	"github.com/De-Rune/jane/package/jane"
)

var builtinFunctions = []*function{
	{
		Name: "print",
		ReturnType: ast.TypeAST{
			Type: jane.Void,
		},
		Params: []ast.ParameterAST{{
			Name: "v",
			Type: ast.TypeAST{
				Value: "any",
				Type:  jane.Any,
			},
		}},
	}, {
		Name: "println",
		ReturnType: ast.TypeAST{
			Type: jane.Void,
		},
		Params: []ast.ParameterAST{{
			Name: "v",
			Type: ast.TypeAST{
				Value: "any",
				Type:  jane.Any,
			},
		}},
	},
}
