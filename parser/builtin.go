package parser

import (
	"github.com/De-Rune/jane/ast"
	"github.com/De-Rune/jane/package/jn"
)

var builtinFuncs = []*function{
	{
		Ast: ast.Func{
			Id: "out",
			RetType: ast.DataType{
				Id: jn.Void,
			},
			Params: []ast.Parameter{{
				Id:    "v",
				Const: true,
				Type: ast.DataType{
					Val: "any",
					Id:  jn.Any,
				},
			}},
		},
	},
	{
		Ast: ast.Func{
			Id: "println",
			RetType: ast.DataType{
				Id: jn.Void,
			},
			Params: []ast.Parameter{{
				Id:    "v",
				Const: true,
				Type: ast.DataType{
					Val: "any",
					Id:  jn.Any,
				},
			}},
		},
	},
}

var strDefs = &defmap{
	Globals: []ast.Var{
		{
			Id:    "len",
			Const: true,
			Type:  ast.DataType{Id: jn.Size, Val: "size"},
			Tag:   "length()",
		},
	},
}

var arrDefs = &defmap{
	Globals: []ast.Var{
		{
			Id:    "len",
			Const: true,
			Type:  ast.DataType{Id: jn.Size, Val: "size"},
			Tag:   "_buffer.size()",
		},
	},
}
