// Copyright (c) 2024 - DeRuneLabs
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package transpiler

import (
	"strconv"

	"github.com/DeRuneLabs/jane/ast"
	"github.com/DeRuneLabs/jane/ast/models"
	"github.com/DeRuneLabs/jane/lexer/tokens"
	"github.com/DeRuneLabs/jane/package/jn"
	"github.com/DeRuneLabs/jane/package/jntype"
)

type BuiltinCaller = func(*Transpiler, *Func, callData, *exprModel) value

const maxI8 = 127
const minI8 = -128
const maxI16 = 32767
const minI16 = -32768
const maxI32 = 2147483647
const minI32 = -2147483648
const maxI64 = 9223372036854775807
const minI64 = -9223372036854775808
const maxU8 = 255
const maxU16 = 65535
const maxU32 = 4294967295
const maxU64 = 18446744073709551615

var i8statics = &DefineMap{
	Globals: []*Var{
		{
			Pub:     true,
			Const:   true,
			Id:      "max",
			Type:    Type{Id: jntype.I8, Kind: tokens.I8},
			ExprTag: int64(maxI8),
			Expr: models.Expr{
				Model: exprNode{jntype.CppId(jntype.I8) + "{" + strconv.FormatInt(maxI8, 10) + "}"},
			},
		},
		{
			Pub:     true,
			Const:   true,
			Id:      "min",
			Type:    Type{Id: jntype.I8, Kind: tokens.I8},
			ExprTag: int64(minI8),
			Expr: models.Expr{
				Model: exprNode{jntype.CppId(jntype.I8) + "{" + strconv.FormatInt(minI8, 10) + "}"},
			},
		},
	},
}

var i16statics = &DefineMap{
	Globals: []*Var{
		{
			Pub:     true,
			Const:   true,
			Id:      "max",
			Type:    Type{Id: jntype.I16, Kind: tokens.I16},
			ExprTag: int64(maxI16),
			Expr: models.Expr{
				Model: exprNode{jntype.CppId(jntype.I16) + "{" + strconv.FormatInt(maxI16, 10) + "}"},
			},
		},
		{
			Pub:     true,
			Const:   true,
			Id:      "min",
			Type:    Type{Id: jntype.I16, Kind: tokens.I16},
			ExprTag: int64(minI16),
			Expr: models.Expr{
				Model: exprNode{jntype.CppId(jntype.I16) + "{" + strconv.FormatInt(minI16, 10) + "}"},
			},
		},
	},
}

var i32statics = &DefineMap{
	Globals: []*Var{
		{
			Pub:     true,
			Const:   true,
			Id:      "max",
			Type:    Type{Id: jntype.I32, Kind: tokens.I32},
			ExprTag: int64(maxI32),
			Expr: models.Expr{
				Model: exprNode{jntype.CppId(jntype.I32) + "{" + strconv.FormatInt(maxI32, 10) + "}"},
			},
		},
		{
			Pub:     true,
			Const:   true,
			Id:      "min",
			Type:    Type{Id: jntype.I32, Kind: tokens.I32},
			ExprTag: int64(minI32),
			Expr: models.Expr{
				Model: exprNode{jntype.CppId(jntype.I32) + "{" + strconv.FormatInt(minI32, 10) + "}"},
			},
		},
	},
}

var i64statics = &DefineMap{
	Globals: []*Var{
		{
			Pub:     true,
			Const:   true,
			Id:      "max",
			Type:    Type{Id: jntype.I64, Kind: tokens.I64},
			ExprTag: int64(maxI64),
			Expr: models.Expr{
				Model: exprNode{jntype.CppId(jntype.I64) + "{" + strconv.FormatInt(maxI64, 10) + "LL}"},
			},
		},
		{
			Pub:     true,
			Const:   true,
			Id:      "min",
			Type:    Type{Id: jntype.I64, Kind: tokens.I64},
			ExprTag: int64(minI64),
			Expr: models.Expr{
				Model: exprNode{jntype.CppId(jntype.I64) + "{" + strconv.FormatInt(minI64, 10) + "LL}"},
			},
		},
	},
}

var u8statics = &DefineMap{
	Globals: []*Var{
		{
			Pub:     true,
			Const:   true,
			Id:      "max",
			Type:    Type{Id: jntype.U8, Kind: tokens.U8},
			ExprTag: uint64(maxU8),
			Expr: models.Expr{
				Model: exprNode{jntype.CppId(jntype.U8) + "{" + strconv.FormatUint(maxU8, 10) + "}"},
			},
		},
	},
}

var u16statics = &DefineMap{
	Globals: []*Var{
		{
			Pub:     true,
			Const:   true,
			Id:      "max",
			Type:    Type{Id: jntype.U16, Kind: tokens.U16},
			ExprTag: uint64(maxU16),
			Expr: models.Expr{
				Model: exprNode{jntype.CppId(jntype.U16) + "{" + strconv.FormatUint(maxU16, 10) + "}"},
			},
		},
	},
}

var u32statics = &DefineMap{
	Globals: []*Var{
		{
			Pub:     true,
			Const:   true,
			Id:      "max",
			Type:    Type{Id: jntype.U32, Kind: tokens.U32},
			ExprTag: uint64(maxU32),
			Expr: models.Expr{
				Model: exprNode{jntype.CppId(jntype.U32) + "{" + strconv.FormatUint(maxU32, 10) + "}"},
			},
		},
	},
}

var u64statics = &DefineMap{
	Globals: []*Var{
		{
			Pub:     true,
			Const:   true,
			Id:      "max",
			Type:    Type{Id: jntype.U64, Kind: tokens.U64},
			ExprTag: uint64(maxU64),
			Expr: models.Expr{
				Model: exprNode{jntype.CppId(jntype.U64) + "{" + strconv.FormatUint(maxU64, 10) + "ULL}"},
			},
		},
	},
}

var uintStatics = &DefineMap{
	Globals: []*Var{
		{
			Pub:   true,
			Const: true,
			Id:    "max",
			Type:  Type{Id: jntype.UInt, Kind: tokens.UINT},
		},
	},
}

var intStatics = &DefineMap{
	Globals: []*Var{
		{
			Const: true,
			Id:    "max",
			Type:  Type{Id: jntype.Int, Kind: tokens.INT},
		},
		{
			Const: true,
			Id:    "min",
			Type:  Type{Id: jntype.Int, Kind: tokens.INT},
		},
	},
}

const maxF32 = 0x1p127 * (1 + (1 - 0x1p-23))
const minF32 = 1.17549435082228750796873653722224568e-38

var min_modelF32 = exprNode{jntype.CppId(jntype.F32) + "{1.17549435082228750796873653722224568e-38F}"}

var f32statics = &DefineMap{
	Globals: []*Var{
		{
			Pub:     true,
			Const:   true,
			Id:      "max",
			Type:    Type{Id: jntype.F32, Kind: tokens.F32},
			ExprTag: float64(maxF32),
			Expr:    models.Expr{Model: exprNode{strconv.FormatFloat(maxF32, 'e', -1, 32) + "F"}},
		},
		{
			Pub:     true,
			Const:   true,
			Id:      "min",
			Type:    Type{Id: jntype.F32, Kind: tokens.F32},
			ExprTag: float64(minF32),
			Expr:    models.Expr{Model: min_modelF32},
		},
	},
}

const maxF64 = 0x1p1023 * (1 + (1 - 0x1p-52))
const minF64 = 2.22507385850720138309023271733240406e-308

var min_modelF64 = exprNode{jntype.CppId(jntype.F64) + "{2.22507385850720138309023271733240406e-308}"}

var f64statics = &DefineMap{
	Globals: []*Var{
		{
			Pub:     true,
			Const:   true,
			Id:      "max",
			Type:    Type{Id: jntype.F64, Kind: tokens.F64},
			ExprTag: float64(maxF64),
			Expr:    models.Expr{Model: exprNode{strconv.FormatFloat(maxF64, 'e', -1, 64)}},
		},
		{
			Pub:     true,
			Const:   true,
			Id:      "min",
			Type:    Type{Id: jntype.F64, Kind: tokens.F64},
			ExprTag: float64(minF64),
			Expr:    models.Expr{Model: min_modelF64},
		},
	},
}

var strDefaultFunc = Func{
	Pub:     true,
	Id:      "str",
	Params:  []Param{{Id: "obj", Type: Type{Id: jntype.Any, Kind: tokens.ANY}}},
	RetType: RetType{Type: Type{Id: jntype.Str, Kind: tokens.STR}},
}

var errorTrait = &trait{
	Ast: &models.Trait{
		Id: "Error",
	},
	Defines: &DefineMap{
		Funcs: []*Fn{
			{Ast: &models.Fn{
				Pub:     true,
				Id:      "error",
				RetType: models.RetType{Type: Type{Id: jntype.Str, Kind: tokens.STR}},
			}},
		},
	},
}

var errorType = Type{
	Id:   jntype.Trait,
	Kind: errorTrait.Ast.Id,
	Tag:  errorTrait,
	Pure: true,
}

var panicFunc = &Fn{
	Ast: &models.Fn{
		Pub: true,
		Id:  "panic",
		Params: []models.Param{
			{
				Id:   "error",
				Type: Type{Id: jntype.Any, Kind: jntype.TypeMap[jntype.Any]},
			},
		},
	},
}

var errorHandlerFunc = &models.Fn{
	Id: "handler",
	Params: []models.Param{
		{
			Id:   "error",
			Type: errorType,
		},
	},
	RetType: models.RetType{
		Type: models.Type{
			Id:   jntype.Void,
			Kind: jntype.TypeMap[jntype.Void],
		},
	},
}

var recoverFunc = &Fn{
	Ast: &models.Fn{
		Pub: true,
		Id:  "recover",
		Params: []models.Param{
			{
				Id: "handler",
				Type: models.Type{
					Id:   jntype.Fn,
					Kind: errorHandlerFunc.DataTypeString(),
					Tag:  errorHandlerFunc,
				},
			},
		},
	},
}

var out_fn = &Fn{Ast: &Func{
	Pub: true,
	Id:  "print",
	RetType: RetType{
		Type: Type{Id: jntype.Void, Kind: jntype.TypeMap[jntype.Void]},
	},
	Params: []Param{{
		Id:   "expr",
		Type: Type{Id: jntype.Any, Kind: tokens.ANY},
	}},
}}

var make_fn = &Fn{Ast: &Func{
	Pub: true,
	Id:  "make",
}}

var outln_fn *Fn
var builtinFile = &Transpiler{}
var Builtin = &DefineMap{
	Types: []*models.TypeAlias{
		{
			Pub:  true,
			Id:   "byte",
			Type: Type{Id: jntype.U8, Kind: jntype.TypeMap[jntype.U8]},
		},
		{
			Pub:  true,
			Id:   "rune",
			Type: Type{Id: jntype.I32, Kind: jntype.TypeMap[jntype.I32]},
		},
	},
	Funcs: []*Fn{
		out_fn,
		panicFunc,
		recoverFunc,
		make_fn,
		{Ast: &Func{
			Pub:   true,
			Id:    "new",
			Owner: builtinFile,
		}},
		{Ast: &Func{
			Pub:      true,
			Id:       "copy",
			Owner:    builtinFile,
			Generics: []*GenericType{{Id: "Item"}},
			RetType:  models.RetType{Type: Type{Id: jntype.Int, Kind: jntype.TypeMap[jntype.Int]}},
			Params: []models.Param{
				{
					Mutable: true,
					Id:      "dest",
					Type: Type{
						Id:            jntype.Slice,
						Kind:          jn.Prefix_Slice + "Item",
						ComponentType: &Type{Id: jntype.Id, Kind: "Item"},
					},
				},
				{
					Id: "src",
					Type: Type{
						Id:            jntype.Slice,
						Kind:          jn.Prefix_Slice + "Item",
						ComponentType: &Type{Id: jntype.Id, Kind: "Item"},
					},
				},
			},
		}},
		{Ast: &Func{
			Pub:      true,
			Id:       "append",
			Owner:    builtinFile,
			Generics: []*GenericType{{Id: "Item"}},
			RetType: models.RetType{
				Type: Type{
					Id:            jntype.Slice,
					Kind:          jn.Prefix_Slice + "Item",
					ComponentType: &Type{Id: jntype.Id, Kind: "Item"},
				},
			},
			Params: []models.Param{
				{
					Id: "src",
					Type: Type{
						Id:            jntype.Slice,
						Kind:          jn.Prefix_Slice + "Item",
						ComponentType: &Type{Id: jntype.Id, Kind: "Item"},
					},
				},
				{
					Id:       "components",
					Type:     Type{Id: jntype.Id, Kind: "Item"},
					Variadic: true,
				},
			},
		}},
	},
	Traits: []*trait{
		errorTrait,
	},
}

var strDefines = &DefineMap{
	Globals: []*Var{
		{
			Pub:  true,
			Id:   "len",
			Type: Type{Id: jntype.Int, Kind: tokens.INT},
			Tag:  "len()",
		},
	},
	Funcs: []*Fn{
		{Ast: &Func{
			Pub:     true,
			Id:      "empty",
			RetType: RetType{Type: Type{Id: jntype.Bool, Kind: tokens.BOOL}},
		}},
		{Ast: &Func{
			Pub:     true,
			Id:      "has_prefix",
			Params:  []Param{{Id: "sub", Type: Type{Id: jntype.Str, Kind: tokens.STR}}},
			RetType: RetType{Type: Type{Id: jntype.Bool, Kind: tokens.BOOL}},
		}},
		{Ast: &Func{
			Pub:     true,
			Id:      "has_suffix",
			Params:  []Param{{Id: "sub", Type: Type{Id: jntype.Str, Kind: tokens.STR}}},
			RetType: RetType{Type: Type{Id: jntype.Bool, Kind: tokens.BOOL}},
		}},
		{Ast: &Func{
			Pub:     true,
			Id:      "find",
			Params:  []Param{{Id: "sub", Type: Type{Id: jntype.Str, Kind: tokens.STR}}},
			RetType: RetType{Type: Type{Id: jntype.Int, Kind: tokens.INT}},
		}},
		{Ast: &Func{
			Pub:     true,
			Id:      "rfind",
			Params:  []Param{{Id: "sub", Type: Type{Id: jntype.Str, Kind: tokens.STR}}},
			RetType: RetType{Type: Type{Id: jntype.Int, Kind: tokens.INT}},
		}},
		{Ast: &Func{
			Pub:     true,
			Id:      "trim",
			Params:  []Param{{Id: "bytes", Type: Type{Id: jntype.Str, Kind: tokens.STR}}},
			RetType: RetType{Type: Type{Id: jntype.Str, Kind: tokens.STR}},
		}},
		{Ast: &Func{
			Pub:     true,
			Id:      "rtrim",
			Params:  []Param{{Id: "bytes", Type: Type{Id: jntype.Str, Kind: tokens.STR}}},
			RetType: RetType{Type: Type{Id: jntype.Str, Kind: tokens.STR}},
		}},
		{Ast: &Func{
			Pub: true,
			Id:  "split",
			Params: []Param{
				{Id: "sub", Type: Type{Id: jntype.Str, Kind: tokens.STR}},
				{
					Id:   "n",
					Type: Type{Id: jntype.Int, Kind: tokens.INT},
				},
			},
			RetType: RetType{Type: Type{Id: jntype.Str, Kind: jn.Prefix_Slice + tokens.STR}},
		}},
		{Ast: &Func{
			Pub: true,
			Id:  "replace",
			Params: []Param{
				{Id: "sub", Type: Type{Id: jntype.Str, Kind: tokens.STR}},
				{Id: "new", Type: Type{Id: jntype.Str, Kind: tokens.STR}},
				{
					Id:   "n",
					Type: Type{Id: jntype.Int, Kind: tokens.INT},
				},
			},
			RetType: RetType{Type: Type{Id: jntype.Str, Kind: tokens.STR}},
		}},
	},
}

func readyStrDefines(s value) {
	lenVar := strDefines.Globals[0]
	lenVar.Const = s.constExpr
	if lenVar.Const {
		lenVar.ExprTag = int64(len(s.expr.(string)))
		lenVar.Expr.Model = getModel(value{
			expr: lenVar.ExprTag,
			data: models.Data{Type: lenVar.Type},
		})
	}
}

var sliceDefines = &DefineMap{
	Globals: []*Var{
		{
			Pub:  true,
			Id:   "len",
			Type: Type{Id: jntype.Int, Kind: tokens.INT},
			Tag:  "len()",
		},
	},
	Funcs: []*Fn{
		{Ast: &Func{
			Pub:     true,
			Id:      "empty",
			RetType: RetType{Type: Type{Id: jntype.Bool, Kind: tokens.BOOL}},
		}},
	},
}

var arrayDefines = &DefineMap{
	Globals: []*Var{
		{
			Pub:  true,
			Id:   "len",
			Type: Type{Id: jntype.Int, Kind: tokens.INT},
			Tag:  "len()",
		},
	},
	Funcs: []*Fn{
		{Ast: &Func{
			Pub:     true,
			Id:      "empty",
			RetType: RetType{Type: Type{Id: jntype.Bool, Kind: tokens.BOOL}},
		}},
	},
}

var mapDefines = &DefineMap{
	Globals: []*Var{
		{
			Pub:  true,
			Id:   "len",
			Type: Type{Id: jntype.Int, Kind: tokens.INT},
			Tag:  "len()",
		},
	},
	Funcs: []*Fn{
		{Ast: &Func{
			Pub: true,
			Id:  "clear",
		}},
		{Ast: &Func{
			Pub: true,
			Id:  "keys",
		}},
		{Ast: &Func{
			Pub: true,
			Id:  "values",
		}},
		{Ast: &Func{
			Pub:     true,
			Id:      "empty",
			RetType: RetType{Type: Type{Id: jntype.Bool, Kind: tokens.BOOL}},
		}},
		{Ast: &Func{
			Pub:     true,
			Id:      "has",
			Params:  []Param{{Id: "key"}},
			RetType: RetType{Type: Type{Id: jntype.Bool, Kind: tokens.BOOL}},
		}},
		{Ast: &Func{
			Pub:    true,
			Id:     "del",
			Params: []Param{{Id: "key"}},
		}},
	},
}

func readyMapDefines(mapt Type) {
	types := mapt.Tag.([]Type)
	keyt := types[0]
	valt := types[1]

	keysFunc, _, _ := mapDefines.funcById("keys", nil)
	keysFunc.Ast.RetType.Type = keyt
	keysFunc.Ast.RetType.Type.Kind = jn.Prefix_Slice + keysFunc.Ast.RetType.Type.Kind

	valuesFunc, _, _ := mapDefines.funcById("values", nil)
	valuesFunc.Ast.RetType.Type = valt
	valuesFunc.Ast.RetType.Type.Kind = jn.Prefix_Slice + valuesFunc.Ast.RetType.Type.Kind

	hasFunc, _, _ := mapDefines.funcById("has", nil)
	hasFunc.Ast.Params[0].Type = keyt

	delFunc, _, _ := mapDefines.funcById("del", nil)
	delFunc.Ast.Params[0].Type = keyt
}

func init() {
	out_fn.Ast.BuiltinCaller = caller_out
	outln_fn = new(Fn)
	*outln_fn = *out_fn
	outln_fn.Ast = new(models.Fn)
	*outln_fn.Ast = *out_fn.Ast
	outln_fn.Ast.Id = "println"
	Builtin.Funcs = append(Builtin.Funcs, outln_fn)
	make_fn.Ast.BuiltinCaller = caller_make

	fn_new, _, _ := Builtin.funcById("new", nil)
	fn_new.Ast.BuiltinCaller = caller_new

	receiver := new(Var)
	receiver.Mutable = false
	for _, f := range errorTrait.Defines.Funcs {
		f.Ast.Receiver = receiver
		f.Ast.Receiver.Tag = errorTrait
		f.Ast.Owner = builtinFile
	}

	intMax := intStatics.Globals[0]
	intMin := intStatics.Globals[1]
	uintMax := uintStatics.Globals[0]
	switch jntype.BitSize {
	case 32:
		intMax.Expr = i32statics.Globals[0].Expr
		intMax.ExprTag = i32statics.Globals[0].ExprTag
		intMin.Expr = i32statics.Globals[1].Expr
		intMin.ExprTag = i32statics.Globals[1].ExprTag

		uintMax.Expr = u32statics.Globals[0].Expr
		uintMax.ExprTag = u32statics.Globals[0].ExprTag
	case 64:
		intMax.Expr = i64statics.Globals[0].Expr
		intMax.ExprTag = i64statics.Globals[0].ExprTag
		intMin.Expr = i64statics.Globals[1].Expr
		intMin.ExprTag = i64statics.Globals[1].ExprTag

		uintMax.Expr = u64statics.Globals[0].Expr
		uintMax.ExprTag = u64statics.Globals[0].ExprTag
	}
}

func caller_out(p *Transpiler, f *Func, data callData, m *exprModel) (v value) {
	errtok := data.args[0]
	v.data.Type = f.RetType.Type
	data.args = data.args[1 : len(data.args)-1]
	arg, model := p.evalToks(data.args)
	if typeIsFunc(arg.data.Type) {
		p.pusherrtok(errtok, "invalid_expr")
	}
	m.appendSubNode(exprNode{"(" + model.String() + ")"})
	arg.constExpr = false
	return v
}

func caller_make(p *Transpiler, _ *Func, data callData, m *exprModel) (v value) {
	errtok := data.args[0]
	args := p.getArgs(data.args, false)
	if len(args.Src) == 0 {
		p.pusherrtok(errtok, "missing_expr")
		return
	}
	type_tokens := args.Src[0].Expr.Tokens
	b := ast.NewBuilder(nil)
	i := 0
	t, ok := b.DataType(type_tokens, &i, false, true)
	b.Wait()
	if !ok {
		p.pusherrs(b.Errors...)
		return
	}
	if i+1 < len(type_tokens) {
		p.pusherrtok(type_tokens[i+1], "invalid_syntax")
	}
	t, ok = p.realType(t, true)
	if !ok {
		return
	}
	switch {
	case typeIsSlice(t):
		return make_slice(p, m, t, args, errtok)
	default:
		p.pusherrtok(errtok, "invalid_type")
	}
	return
}

func caller_new(p *Transpiler, _ *Func, data callData, m *exprModel) (v value) {
	errtok := data.args[0]
	data.args = data.args[1 : len(data.args)-1]
	b := ast.NewBuilder(nil)
	i := 0
	t, ok := b.DataType(data.args, &i, true, true)
	b.Wait()
	if !ok {
		p.pusherrs(b.Errors...)
		return
	}
	if i+1 < len(data.args) {
		p.pusherrtok(data.args[i+1], "invalid_syntax")
	}
	t, _ = p.realType(t, true)
	if !is_valid_type_for_reference(t) {
		p.pusherrtok(errtok, "invalid_type")
	}
	if typeIsStruct(t) {
		s := t.Tag.(*structure)
		for _, f := range s.Defines.Globals {
			if typeIsRef(f.Type) {
				p.pusherrtok(errtok, "ref_used_struct_used_at_new_fn")
				break
			}
		}
	}
	m.appendSubNode(exprNode{"<" + t.String() + ">()"})
	t.Kind = "&" + t.Kind
	v.data.Type = t
	v.data.Value = t.Kind
	return
}

var std_mem_builtin = &DefineMap{
	Funcs: []*Fn{
		{Ast: &models.Fn{Id: "size_of", BuiltinCaller: caller_mem_size_of}},
		{Ast: &models.Fn{Id: "align_of", BuiltinCaller: caller_mem_align_of}},
	},
}

func caller_mem_size_of(p *Transpiler, _ *Func, data callData, m *exprModel) (v value) {

	data.args = data.args[1 : len(data.args)-1]
	v.data.Type = Type{
		Id:   jntype.UInt,
		Kind: jntype.TypeMap[jntype.UInt],
	}
	nodes := m.nodes[m.index].nodes
	node := &nodes[len(nodes)-1]
	b := ast.NewBuilder(nil)
	i := 0
	t, ok := b.DataType(data.args, &i, true, true)
	b.Wait()
	if !ok {
		v, model := p.evalToks(data.args)
		*node = exprNode{"sizeof(" + model.String() + ")"}
		v.constExpr = false
		return v
	}
	if i+1 < len(data.args) {
		p.pusherrtok(data.args[i+1], "invalid_syntax")
	}
	p.pusherrs(b.Errors...)
	t, _ = p.realType(t, true)
	*node = exprNode{"sizeof(" + t.String() + ")"}
	return
}

func caller_mem_align_of(p *Transpiler, _ *Func, data callData, m *exprModel) (v value) {
	data.args = data.args[1 : len(data.args)-1]
	v.data.Type = Type{
		Id:   jntype.UInt,
		Kind: jntype.TypeMap[jntype.UInt],
	}
	nodes := m.nodes[m.index].nodes
	node := &nodes[len(nodes)-1]
	b := ast.NewBuilder(nil)
	i := 0
	t, ok := b.DataType(data.args, &i, true, true)
	b.Wait()
	if !ok {
		v, model := p.evalToks(data.args)
		*node = exprNode{"alignof(" + model.String() + ")"}
		v.constExpr = false
		return v
	}
	if i+1 < len(data.args) {
		p.pusherrtok(data.args[i+1], "invalid_syntax")
	}
	p.pusherrs(b.Errors...)
	t, _ = p.realType(t, true)
	*node = exprNode{"alignof(" + t.String() + ")"}
	return
}

var std_builtin_defines = map[string]*DefineMap{
	"std::mem": std_mem_builtin,
}
