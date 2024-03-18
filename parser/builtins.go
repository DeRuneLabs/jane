package parser

import (
	"github.com/DeRuneLabs/jane/lexer/tokens"
	"github.com/DeRuneLabs/jane/package/jnapi"
	"github.com/DeRuneLabs/jane/package/jntype"
)

var i8statics = &Defmap{
	Globals: []*Var{
		{
			Pub:   true,
			Id:    "max",
			Const: true,
			Type:  DataType{Id: jntype.I8, Val: tokens.I8},
			Tag:   "INT8_MAX",
		},
		{
			Pub:   true,
			Id:    "min",
			Const: true,
			Type:  DataType{Id: jntype.I8, Val: tokens.I8},
			Tag:   "INT8_MIN",
		},
	},
}

var i16statics = &Defmap{
	Globals: []*Var{
		{
			Pub:   true,
			Id:    "max",
			Const: true,
			Type:  DataType{Id: jntype.I16, Val: tokens.I16},
			Tag:   "INT16_MAX",
		},
		{
			Pub:   true,
			Id:    "min",
			Const: true,
			Type:  DataType{Id: jntype.I16, Val: tokens.I16},
			Tag:   "INT16_MIN",
		},
	},
}

var i32statics = &Defmap{
	Globals: []*Var{
		{
			Pub:   true,
			Id:    "max",
			Const: true,
			Type:  DataType{Id: jntype.I32, Val: tokens.I32},
			Tag:   "INT32_MAX",
		},
		{
			Pub:   true,
			Id:    "min",
			Const: true,
			Type:  DataType{Id: jntype.I32, Val: tokens.I32},
			Tag:   "INT32_MIN",
		},
	},
}

var i64statics = &Defmap{
	Globals: []*Var{
		{
			Pub:   true,
			Id:    "max",
			Const: true,
			Type:  DataType{Id: jntype.I64, Val: tokens.I64},
			Tag:   "INT64_MAX",
		},
		{
			Pub:   true,
			Id:    "min",
			Const: true,
			Type:  DataType{Id: jntype.I64, Val: tokens.I64},
			Tag:   "INT64_MIN",
		},
	},
}

var u8statics = &Defmap{
	Globals: []*Var{
		{
			Pub:   true,
			Id:    "max",
			Const: true,
			Type:  DataType{Id: jntype.U8, Val: tokens.U8},
			Tag:   "UINT8_MAX",
		},
	},
}

var u16statics = &Defmap{
	Globals: []*Var{
		{
			Pub:   true,
			Id:    "max",
			Const: true,
			Type:  DataType{Id: jntype.U16, Val: tokens.U16},
			Tag:   "UINT16_MAX",
		},
	},
}

var u32statics = &Defmap{
	Globals: []*Var{
		{
			Pub:   true,
			Id:    "max",
			Const: true,
			Type:  DataType{Id: jntype.U32, Val: tokens.U32},
			Tag:   "UINT32_MAX",
		},
	},
}

var u64statics = &Defmap{
	Globals: []*Var{
		{
			Pub:   true,
			Id:    "max",
			Const: true,
			Type:  DataType{Id: jntype.U64, Val: tokens.U64},
			Tag:   "UINT64_MAX",
		},
	},
}

var uintStatics = &Defmap{
	Globals: []*Var{
		{
			Pub:   true,
			Id:    "max",
			Const: true,
			Type:  DataType{Id: jntype.UInt, Val: tokens.UINT},
			Tag:   "SIZE_MAX",
		},
	},
}

var intStatics = &Defmap{
	Globals: []*Var{
		{
			Id:    "max",
			Const: true,
			Type:  DataType{Id: jntype.Int, Val: tokens.INT},
			Tag:   "",
		},
		{
			Id:    "min",
			Const: true,
			Type:  DataType{Id: jntype.Int, Val: tokens.INT},
			Tag:   "",
		},
	},
}

var f32statics = &Defmap{
	Globals: []*Var{
		{
			Pub:   true,
			Id:    "max",
			Const: true,
			Type:  DataType{Id: jntype.F32, Val: tokens.F32},
			Tag:   "__FLT_MAX__",
		},
		{
			Pub:   true,
			Id:    "min",
			Const: true,
			Type:  DataType{Id: jntype.F32, Val: tokens.F32},
			Tag:   "__FLT_MIN__",
		},
	},
}

var f64statics = &Defmap{
	Globals: []*Var{
		{
			Pub:   true,
			Id:    "max",
			Const: true,
			Type:  DataType{Id: jntype.F64, Val: tokens.F64},
			Tag:   "__DBL_MAX__",
		},
		{
			Pub:   true,
			Id:    "min",
			Const: true,
			Type:  DataType{Id: jntype.F64, Val: tokens.F64},
			Tag:   "__DBL_MIN__",
		},
	},
}

var strStatics = &Defmap{
	Globals: []*Var{
		{
			Pub:   true,
			Id:    "npos",
			Const: true,
			Type:  DataType{Id: jntype.UInt, Val: tokens.UINT},
			Tag:   "std::string::npos",
		},
	},
}

var strDefaultFunc = Func{
	Pub:     true,
	Id:      "str",
	Params:  []Param{{Id: "obj", Type: DataType{Id: jntype.Any, Val: "any"}}},
	RetType: DataType{Id: jntype.Str, Val: tokens.STR},
}

var errorStruct = &jnstruct{
	Ast: Struct{
		Id: "error",
	},
	Defs: &Defmap{
		Globals: []*Var{
			{
				Pub:  true,
				Id:   "message",
				Type: DataType{Id: jntype.Str, Val: tokens.STR},
			},
		},
	},
	constructor: &Func{
		Pub: true,
		Params: []Param{
			{
				Id:      "message",
				Type:    DataType{Id: jntype.Str, Val: tokens.STR},
				Default: Expr{Model: exprNode{jnapi.ToStr(`"error: undefined error"`)}},
			},
		},
	},
}

var errorType = DataType{Id: jntype.Struct, Val: "error", Tag: errorStruct}

var Builtin = &Defmap{
	Funcs: []*function{
		{
			Ast: &Func{
				Pub:     true,
				Id:      "print",
				RetType: DataType{Id: jntype.Void, Val: jntype.VoidTypeStr},
				Params: []Param{{
					Id:      "v",
					Const:   true,
					Type:    DataType{Id: jntype.Any, Val: "any"},
					Default: Expr{Model: exprNode{`""`}},
				}},
			},
		},
		{
			Ast: &Func{
				Pub:     true,
				Id:      "println",
				RetType: DataType{Id: jntype.Void, Val: jntype.VoidTypeStr},
				Params: []Param{{
					Id:      "v",
					Const:   true,
					Type:    DataType{Id: jntype.Any, Val: "any"},
					Default: Expr{Model: exprNode{`""`}},
				}},
			},
		},
	},
	Structs: []*jnstruct{
		errorStruct,
	},
}

var strDefs = &Defmap{
	Globals: []*Var{
		{
			Pub:   true,
			Id:    "len",
			Const: true,
			Type:  DataType{Id: jntype.UInt, Val: tokens.UINT},
			Tag:   "len()",
		},
	},
	Funcs: []*function{
		{Ast: &Func{
			Pub:     true,
			Id:      "empty",
			RetType: DataType{Id: jntype.Bool, Val: tokens.BOOL},
		}},
		{Ast: &Func{
			Pub:     true,
			Id:      "has_prefix",
			Params:  []Param{{Id: "sub", Type: DataType{Id: jntype.Str, Val: tokens.STR}}},
			RetType: DataType{Id: jntype.Bool, Val: tokens.BOOL},
		}},
		{Ast: &Func{
			Pub:     true,
			Id:      "has_suffix",
			Params:  []Param{{Id: "sub", Type: DataType{Id: jntype.Str, Val: tokens.STR}}},
			RetType: DataType{Id: jntype.Bool, Val: tokens.BOOL},
		}},
		{Ast: &Func{
			Pub:     true,
			Id:      "find",
			Params:  []Param{{Id: "sub", Type: DataType{Id: jntype.Str, Val: tokens.STR}}},
			RetType: DataType{Id: jntype.UInt, Val: tokens.UINT},
		}},
		{Ast: &Func{
			Pub:     true,
			Id:      "rfind",
			Params:  []Param{{Id: "sub", Type: DataType{Id: jntype.Str, Val: tokens.STR}}},
			RetType: DataType{Id: jntype.UInt, Val: tokens.UINT},
		}},
		{Ast: &Func{
			Pub:     true,
			Id:      "trim",
			Params:  []Param{{Id: "bytes", Type: DataType{Id: jntype.Str, Val: tokens.STR}}},
			RetType: DataType{Id: jntype.Str, Val: tokens.STR},
		}},
		{Ast: &Func{
			Pub:     true,
			Id:      "rtrim",
			Params:  []Param{{Id: "bytes", Type: DataType{Id: jntype.Str, Val: tokens.STR}}},
			RetType: DataType{Id: jntype.Str, Val: tokens.STR},
		}},
		{Ast: &Func{
			Pub: true,
			Id:  "split",
			Params: []Param{
				{Id: "sub", Type: DataType{Id: jntype.Str, Val: tokens.STR}},
				{
					Id:      "n",
					Type:    DataType{Id: jntype.I64, Val: tokens.I64},
					Default: Expr{Model: exprNode{"-1"}},
				},
			},
			RetType: DataType{Id: jntype.Str, Val: "[]" + tokens.STR},
		}},
		{Ast: &Func{
			Pub: true,
			Id:  "replace",
			Params: []Param{
				{Id: "sub", Type: DataType{Id: jntype.Str, Val: tokens.STR}},
				{Id: "new", Type: DataType{Id: jntype.Str, Val: tokens.STR}},
				{
					Id:      "n",
					Type:    DataType{Id: jntype.I64, Val: tokens.I64},
					Default: Expr{Model: exprNode{"-1"}},
				},
			},
			RetType: DataType{Id: jntype.Str, Val: tokens.STR},
		}},
	},
}

var arrDefs = &Defmap{
	Globals: []*Var{
		{
			Pub:   true,
			Id:    "len",
			Const: true,
			Type:  DataType{Id: jntype.UInt, Val: tokens.UINT},
			Tag:   "len()",
		},
	},
	Funcs: []*function{
		{Ast: &Func{Pub: true, Id: "clear"}},
		{Ast: &Func{
			Pub:     true,
			Id:      "empty",
			RetType: DataType{Id: jntype.Bool, Val: tokens.BOOL},
		}},
		{Ast: &Func{
			Pub:    true,
			Id:     "find",
			Params: []Param{{Id: "value"}},
		}},
		{Ast: &Func{
			Pub:    true,
			Id:     "rfind",
			Params: []Param{{Id: "value"}},
		}},
		{Ast: &Func{
			Pub:    true,
			Id:     "erase",
			Params: []Param{{Id: "value"}},
		}},
		{Ast: &Func{
			Pub:    true,
			Id:     "erase_all",
			Params: []Param{{Id: "value"}},
		}},
		{Ast: &Func{
			Pub:    true,
			Id:     "append",
			Params: []Param{{Id: "values", Variadic: true}},
		}},
		{Ast: &Func{
			Pub: true,
			Id:  "insert",
			Params: []Param{
				{Id: "start", Type: DataType{Id: jntype.UInt, Val: tokens.UINT}},
				{Id: "values", Variadic: true},
			},
			RetType: DataType{Id: jntype.Bool, Val: tokens.BOOL},
		}},
	},
}

func readyArrDefs(arrt DataType) {
	elemType := typeOfArrayComponents(arrt)

	findFunc, _, _ := arrDefs.funcById("find", nil)
	findFunc.Ast.Params[0].Type = elemType
	findFunc.Ast.RetType = elemType
	findFunc.Ast.RetType.Val = tokens.STAR + findFunc.Ast.RetType.Val

	rfindFunc, _, _ := arrDefs.funcById("rfind", nil)
	rfindFunc.Ast.Params[0].Type = elemType
	rfindFunc.Ast.RetType = elemType
	rfindFunc.Ast.RetType.Val = tokens.STAR + rfindFunc.Ast.RetType.Val

	eraseFunc, _, _ := arrDefs.funcById("erase", nil)
	eraseFunc.Ast.Params[0].Type = elemType

	eraseAllFunc, _, _ := arrDefs.funcById("erase_all", nil)
	eraseAllFunc.Ast.Params[0].Type = elemType

	appendFunc, _, _ := arrDefs.funcById("append", nil)
	appendFunc.Ast.Params[0].Type = elemType

	insertFunc, _, _ := arrDefs.funcById("insert", nil)
	insertFunc.Ast.Params[1].Type = elemType
}

var mapDefs = &Defmap{
	Globals: []*Var{
		{
			Pub:   true,
			Id:    "len",
			Const: true,
			Type:  DataType{Id: jntype.UInt, Val: tokens.UINT},
			Tag:   "size()",
		},
	},
	Funcs: []*function{
		{Ast: &Func{Pub: true, Id: "clear"}},
		{Ast: &Func{Pub: true, Id: "keys"}},
		{Ast: &Func{Pub: true, Id: "values"}},
		{Ast: &Func{
			Pub:     true,
			Id:      "empty",
			RetType: DataType{Id: jntype.Bool, Val: tokens.BOOL},
		}},
		{Ast: &Func{
			Pub:     true,
			Id:      "has",
			Params:  []Param{{Id: "key", Const: true}},
			RetType: DataType{Id: jntype.Bool, Val: tokens.BOOL},
		}},
		{Ast: &Func{
			Pub:    true,
			Id:     "del",
			Params: []Param{{Id: "key", Const: true}},
		}},
	},
}

func readyMapDefs(mapt DataType) {
	types := mapt.Tag.([]DataType)
	keyt := types[0]
	valt := types[1]

	keysFunc, _, _ := mapDefs.funcById("keys", nil)
	keysFunc.Ast.RetType = keyt
	keysFunc.Ast.RetType.Val = "[]" + keysFunc.Ast.RetType.Val

	valuesFunc, _, _ := mapDefs.funcById("values", nil)
	valuesFunc.Ast.RetType = valt
	valuesFunc.Ast.RetType.Val = "[]" + valuesFunc.Ast.RetType.Val

	hasFunc, _, _ := mapDefs.funcById("has", nil)
	hasFunc.Ast.Params[0].Type = keyt

	delFunc, _, _ := mapDefs.funcById("del", nil)
	delFunc.Ast.Params[0].Type = keyt
}

func init() {
	intMax := intStatics.Globals[0]
	intMin := intStatics.Globals[1]
	switch jntype.BitSize {
	case 8:
		intMax.Tag = i8statics.Globals[0].Tag
		intMin.Tag = i8statics.Globals[1].Tag
	case 16:
		intMax.Tag = i16statics.Globals[0].Tag
		intMin.Tag = i16statics.Globals[1].Tag
	case 32:
		intMax.Tag = i32statics.Globals[0].Tag
		intMin.Tag = i32statics.Globals[1].Tag
	case 64:
		intMax.Tag = i64statics.Globals[0].Tag
		intMin.Tag = i64statics.Globals[1].Tag
	}

	errorStruct.constructor.Id = errorStruct.Ast.Id
	errorStruct.constructor.RetType = DataType{
		Id:  jntype.Struct,
		Val: errorStruct.Ast.Id,
		Tag: errorStruct,
	}
}
