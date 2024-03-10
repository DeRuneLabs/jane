package parser

import (
	"github.com/DeRuneLabs/jane/ast"
	"github.com/DeRuneLabs/jane/package/jn"
)

var builtinFuncs = []*function{
	{
		Ast: ast.Func{
			Id:      "print",
			RetType: ast.DataType{Id: jn.Void},
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
			Id:      "println",
			RetType: ast.DataType{Id: jn.Void},
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
	Globals: []*ast.Var{
		{
			Id:    "len",
			Const: true,
			Type:  ast.DataType{Id: jn.Size, Val: "size"},
			Tag:   "length()",
		},
	},
}

var arrDefs = &defmap{
	Globals: []*ast.Var{
		{
			Id:    "len",
			Const: true,
			Type:  ast.DataType{Id: jn.Size, Val: "size"},
			Tag:   "_buffer.size()",
		},
	},
}

var mapDefs = &defmap{
	Globals: []*ast.Var{
		{
			Id:    "len",
			Const: true,
			Type:  ast.DataType{Id: jn.Size, Val: "size"},
			Tag:   "size()",
		},
	},
	Funcs: []*function{
		{Ast: ast.Func{Id: "clear"}},
		{Ast: ast.Func{Id: "keys"}},
		{Ast: ast.Func{Id: "values"}},
		{Ast: ast.Func{
			Id:      "has",
			Params:  []ast.Parameter{{Id: "key", Const: true}},
			RetType: ast.DataType{Id: jn.Bool, Val: "bool"},
		}},
	},
}

func readyMapDefs(mapt ast.DataType) {
	types := mapt.Tag.([]ast.DataType)
	keyt := types[0]
	valt := types[1]

	keysFunc := mapDefs.funcById("keys")
	keysFunc.Ast.RetType = keyt
	keysFunc.Ast.RetType.Val = "[]" + keysFunc.Ast.RetType.Val

	valuesFunc := mapDefs.funcById("values")
	valuesFunc.Ast.RetType = valt
	valuesFunc.Ast.RetType.Val = "[]" + valuesFunc.Ast.RetType.Val

	hasFunc := mapDefs.funcById("has")
	hasFunc.Ast.Params[0].Type = keyt
}
