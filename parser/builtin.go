package parser

import (
	"github.com/De-Rune/jane/ast"
	"github.com/De-Rune/jane/package/jn"
)

var builtinFunctions = []*function{
	{
		Ast: ast.FunctionAST{
			Name: "_print",
			ReturnType: ast.DataTypeAST{
				Code: jn.Void,
			},
			Params: []ast.ParameterAST{{
				Name: "v",
				Type: ast.DataTypeAST{
					Value: "any",
					Code:  jn.Any,
				},
			}},
		},
	},
	{
		Ast: ast.FunctionAST{
			Name: "_println",
			ReturnType: ast.DataTypeAST{
				Code: jn.Void,
			},
			Params: []ast.ParameterAST{{
				Name: "v",
				Type: ast.DataTypeAST{
					Value: "any",
					Code:  jn.Any,
				},
			}},
		},
	},
}
