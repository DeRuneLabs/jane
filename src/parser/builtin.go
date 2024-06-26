// Copyright (c) 2024 arfy slowy - DeRuneLabs
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

package parser

import (
	"strconv"

	"github.com/DeRuneLabs/jane/ast"
	"github.com/DeRuneLabs/jane/lexer"
	"github.com/DeRuneLabs/jane/types"
)

type builtin_caller = func(*Parser, *Fn, call_data, *expr_model) value

const (
	maxI8  = 127
	minI8  = -128
	maxI16 = 32767
	minI16 = -32768
	maxI32 = 2147483647
	minI32 = -2147483648
	maxI64 = 9223372036854775807
	minI64 = -9223372036854775808
	maxU8  = 255
	maxU16 = 65535
	maxU32 = 4294967295
	maxU64 = 18446744073709551615
)

var i8statics = &ast.Defmap{
	Globals: []*Var{
		{
			Public:   true,
			Constant: true,
			Id:       "max",
			DataType: Type{Id: types.I8, Kind: types.TYPE_MAP[types.I8]},
			ExprTag:  int64(maxI8),
			Expr: ast.Expr{
				Model: exprNode{types.CppId(types.I8) + "{" + strconv.FormatInt(maxI8, 10) + "}"},
			},
		},
		{
			Public:   true,
			Constant: true,
			Id:       "min",
			DataType: Type{Id: types.I8, Kind: types.TYPE_MAP[types.I8]},
			ExprTag:  int64(minI8),
			Expr: ast.Expr{
				Model: exprNode{types.CppId(types.I8) + "{" + strconv.FormatInt(minI8, 10) + "}"},
			},
		},
	},
}

var i16statics = &ast.Defmap{
	Globals: []*Var{
		{
			Public:   true,
			Constant: true,
			Id:       "max",
			DataType: Type{Id: types.I16, Kind: types.TYPE_MAP[types.I16]},
			ExprTag:  int64(maxI16),
			Expr: ast.Expr{
				Model: exprNode{types.CppId(types.I16) + "{" + strconv.FormatInt(maxI16, 10) + "}"},
			},
		},
		{
			Public:   true,
			Constant: true,
			Id:       "min",
			DataType: Type{Id: types.I16, Kind: types.TYPE_MAP[types.I16]},
			ExprTag:  int64(minI16),
			Expr: ast.Expr{
				Model: exprNode{types.CppId(types.I16) + "{" + strconv.FormatInt(minI16, 10) + "}"},
			},
		},
	},
}

var i32statics = &ast.Defmap{
	Globals: []*Var{
		{
			Public:   true,
			Constant: true,
			Id:       "max",
			DataType: Type{Id: types.I32, Kind: types.TYPE_MAP[types.I32]},
			ExprTag:  int64(maxI32),
			Expr: ast.Expr{
				Model: exprNode{types.CppId(types.I32) + "{" + strconv.FormatInt(maxI32, 10) + "}"},
			},
		},
		{
			Public:   true,
			Constant: true,
			Id:       "min",
			DataType: Type{Id: types.I32, Kind: types.TYPE_MAP[types.I32]},
			ExprTag:  int64(minI32),
			Expr: ast.Expr{
				Model: exprNode{types.CppId(types.I32) + "{" + strconv.FormatInt(minI32, 10) + "}"},
			},
		},
	},
}

var i64statics = &ast.Defmap{
	Globals: []*Var{
		{
			Public:   true,
			Constant: true,
			Id:       "max",
			DataType: Type{Id: types.I64, Kind: types.TYPE_MAP[types.I64]},
			ExprTag:  int64(maxI64),
			Expr: ast.Expr{
				Model: exprNode{
					types.CppId(types.I64) + "{" + strconv.FormatInt(maxI64, 10) + "LL}",
				},
			},
		},
		{
			Public:   true,
			Constant: true,
			Id:       "min",
			DataType: Type{Id: types.I64, Kind: types.TYPE_MAP[types.I64]},
			ExprTag:  int64(minI64),
			Expr: ast.Expr{
				Model: exprNode{
					types.CppId(types.I64) + "{" + strconv.FormatInt(minI64, 10) + "LL}",
				},
			},
		},
	},
}

var u8statics = &ast.Defmap{
	Globals: []*Var{
		{
			Public:   true,
			Constant: true,
			Id:       "max",
			DataType: Type{Id: types.U8, Kind: types.TYPE_MAP[types.U8]},
			ExprTag:  uint64(maxU8),
			Expr: ast.Expr{
				Model: exprNode{types.CppId(types.U8) + "{" + strconv.FormatUint(maxU8, 10) + "}"},
			},
		},
	},
}

var u16statics = &ast.Defmap{
	Globals: []*Var{
		{
			Public:   true,
			Constant: true,
			Id:       "max",
			DataType: Type{Id: types.U16, Kind: types.TYPE_MAP[types.U16]},
			ExprTag:  uint64(maxU16),
			Expr: ast.Expr{
				Model: exprNode{
					types.CppId(types.U16) + "{" + strconv.FormatUint(maxU16, 10) + "}",
				},
			},
		},
	},
}

var u32statics = &ast.Defmap{
	Globals: []*Var{
		{
			Public:   true,
			Constant: true,
			Id:       "max",
			DataType: Type{Id: types.U32, Kind: types.TYPE_MAP[types.U32]},
			ExprTag:  uint64(maxU32),
			Expr: ast.Expr{
				Model: exprNode{
					types.CppId(types.U32) + "{" + strconv.FormatUint(maxU32, 10) + "}",
				},
			},
		},
	},
}

var u64statics = &ast.Defmap{
	Globals: []*Var{
		{
			Public:   true,
			Constant: true,
			Id:       "max",
			DataType: Type{Id: types.U64, Kind: types.TYPE_MAP[types.U64]},
			ExprTag:  uint64(maxU64),
			Expr: ast.Expr{
				Model: exprNode{
					types.CppId(types.U64) + "{" + strconv.FormatUint(maxU64, 10) + "ULL}",
				},
			},
		},
	},
}

var uintStatics = &ast.Defmap{
	Globals: []*Var{
		{
			Public:   true,
			Constant: true,
			Id:       "max",
			DataType: Type{Id: types.UINT, Kind: types.TYPE_MAP[types.UINT]},
		},
	},
}

var intStatics = &ast.Defmap{
	Globals: []*Var{
		{
			Constant: true,
			Id:       "max",
			DataType: Type{Id: types.INT, Kind: types.TYPE_MAP[types.INT]},
		},
		{
			Constant: true,
			Id:       "min",
			DataType: Type{Id: types.INT, Kind: types.TYPE_MAP[types.INT]},
		},
	},
}

const (
	maxF32 = 0x1p127 * (1 + (1 - 0x1p-23))
	minF32 = 1.17549435082228750796873653722224568e-38
)

var min_modelF32 = exprNode{types.CppId(types.F32) + "{1.17549435082228750796873653722224568e-38F}"}

var f32statics = &ast.Defmap{
	Globals: []*Var{
		{
			Public:   true,
			Constant: true,
			Id:       "max",
			DataType: Type{Id: types.F32, Kind: types.TYPE_MAP[types.F32]},
			ExprTag:  float64(maxF32),
			Expr:     ast.Expr{Model: exprNode{strconv.FormatFloat(maxF32, 'e', -1, 32) + "F"}},
		},
		{
			Public:   true,
			Constant: true,
			Id:       "min",
			DataType: Type{Id: types.F32, Kind: types.TYPE_MAP[types.F32]},
			ExprTag:  float64(minF32),
			Expr:     ast.Expr{Model: min_modelF32},
		},
	},
}

const (
	maxF64 = 0x1p1023 * (1 + (1 - 0x1p-52))
	minF64 = 2.22507385850720138309023271733240406e-308
)

var min_modelF64 = exprNode{types.CppId(types.F64) + "{2.22507385850720138309023271733240406e-308}"}

var f64statics = &ast.Defmap{
	Globals: []*Var{
		{
			Public:   true,
			Constant: true,
			Id:       "max",
			DataType: Type{Id: types.F64, Kind: types.TYPE_MAP[types.F64]},
			ExprTag:  float64(maxF64),
			Expr:     ast.Expr{Model: exprNode{strconv.FormatFloat(maxF64, 'e', -1, 64)}},
		},
		{
			Public:   true,
			Constant: true,
			Id:       "min",
			DataType: Type{Id: types.F64, Kind: types.TYPE_MAP[types.F64]},
			ExprTag:  float64(minF64),
			Expr:     ast.Expr{Model: min_modelF64},
		},
	},
}

var errorTrait = &ast.Trait{
	Id: "Error",
	Defines: &ast.Defmap{
		Fns: []*Fn{
			{
				Public: true,
				Id:     "error",
				RetType: ast.RetType{
					DataType: Type{Id: types.STR, Kind: types.TYPE_MAP[types.STR]},
				},
			},
		},
	},
}

var errorType = Type{
	Id:   types.TRAIT,
	Kind: errorTrait.Id,
	Tag:  errorTrait,
	Pure: true,
}

var panicFunc = &Fn{
	Public: true,
	Id:     "panic",
	Params: []ast.Param{
		{
			Id:       "error",
			DataType: Type{Id: types.ANY, Kind: types.TYPE_MAP[types.ANY]},
		},
	},
}

var errorHandlerFunc = &Fn{
	Id: "handler",
	Params: []ast.Param{
		{
			Id:       "error",
			DataType: errorType,
		},
	},
	RetType: ast.RetType{
		DataType: ast.Type{
			Id:   types.VOID,
			Kind: types.TYPE_MAP[types.VOID],
		},
	},
}

var recoverFunc = &Fn{
	Public: true,
	Id:     "recover",
	Params: []ast.Param{
		{
			Id: "handler",
			DataType: ast.Type{
				Id:   types.FN,
				Kind: errorHandlerFunc.TypeKind(),
				Tag:  errorHandlerFunc,
			},
		},
	},
}

var out_fn = &Fn{
	Public: true,
	Id:     "print",
	RetType: RetType{
		DataType: Type{Id: types.VOID, Kind: types.TYPE_MAP[types.VOID]},
	},
	Params: []Param{{
		Id:       "expr",
		DataType: Type{Id: types.ANY, Kind: types.TYPE_MAP[types.ANY]},
	}},
}

var (
	make_fn   = &Fn{Public: true, Id: "make"}
	copy_fn   = &Fn{Public: true, Id: "copy"}
	append_fn = &Fn{Public: true, Id: "append"}
	new_fn    = &Fn{Public: true, Id: "new"}
	drop_fn   = &Fn{Public: true, Id: "drop"}
	real_fn   = &Fn{Public: true, Id: "real"}
)

var outln_fn *Fn

var builtinFile = &Parser{}

var Builtin = &ast.Defmap{
	Types: []*ast.TypeAlias{
		{
			Pub:        true,
			Id:         "byte",
			TargetType: Type{Id: types.U8, Kind: types.TYPE_MAP[types.U8]},
		},
		{
			Pub:        true,
			Id:         "rune",
			TargetType: Type{Id: types.I32, Kind: types.TYPE_MAP[types.I32]},
		},
	},
	Fns: []*Fn{
		out_fn,
		panicFunc,
		recoverFunc,
		new_fn,
		drop_fn,
		real_fn,
		make_fn,
		copy_fn,
		append_fn,
	},
	Traits: []*ast.Trait{
		errorTrait,
	},
}

var strDefines = &ast.Defmap{
	Globals: []*Var{
		{
			Public:   true,
			Id:       "len",
			DataType: Type{Id: types.INT, Kind: types.TYPE_MAP[types.INT]},
			Tag:      "_len()",
		},
	},
	Fns: []*Fn{
		{
			Public: true,
			Id:     "has_prefix",
			Params: []Param{
				{Id: "sub", DataType: Type{Id: types.STR, Kind: types.TYPE_MAP[types.STR]}},
			},
			RetType: RetType{DataType: Type{Id: types.BOOL, Kind: types.TYPE_MAP[types.BOOL]}},
		},
		{
			Public: true,
			Id:     "has_suffix",
			Params: []Param{
				{Id: "sub", DataType: Type{Id: types.STR, Kind: types.TYPE_MAP[types.STR]}},
			},
			RetType: RetType{DataType: Type{Id: types.BOOL, Kind: types.TYPE_MAP[types.BOOL]}},
		},
		{
			Public: true,
			Id:     "find",
			Params: []Param{
				{Id: "sub", DataType: Type{Id: types.STR, Kind: types.TYPE_MAP[types.STR]}},
			},
			RetType: RetType{DataType: Type{Id: types.INT, Kind: types.TYPE_MAP[types.INT]}},
		},
		{
			Public: true,
			Id:     "rfind",
			Params: []Param{
				{Id: "sub", DataType: Type{Id: types.STR, Kind: types.TYPE_MAP[types.STR]}},
			},
			RetType: RetType{DataType: Type{Id: types.INT, Kind: types.TYPE_MAP[types.INT]}},
		},
		{
			Public: true,
			Id:     "trim",
			Params: []Param{
				{Id: "bytes", DataType: Type{Id: types.STR, Kind: types.TYPE_MAP[types.STR]}},
			},
			RetType: RetType{
				DataType: Type{Id: types.STR, Kind: types.TYPE_MAP[types.STR]},
			},
		},
		{
			Public: true,
			Id:     "rtrim",
			Params: []Param{
				{Id: "bytes", DataType: Type{Id: types.STR, Kind: types.TYPE_MAP[types.STR]}},
			},
			RetType: RetType{DataType: Type{Id: types.STR, Kind: types.TYPE_MAP[types.STR]}},
		},
		{
			Public: true,
			Id:     "split",
			Params: []Param{
				{Id: "sub", DataType: Type{Id: types.STR, Kind: types.TYPE_MAP[types.STR]}},
				{
					Id:       "n",
					DataType: Type{Id: types.INT, Kind: types.TYPE_MAP[types.INT]},
				},
			},
			RetType: RetType{
				DataType: Type{Id: types.STR, Kind: lexer.PREFIX_SLICE + types.TYPE_MAP[types.STR]},
			},
		},
		{
			Public: true,
			Id:     "replace",
			Params: []Param{
				{Id: "sub", DataType: Type{Id: types.STR, Kind: types.TYPE_MAP[types.STR]}},
				{Id: "new", DataType: Type{Id: types.STR, Kind: types.TYPE_MAP[types.STR]}},
				{
					Id:       "n",
					DataType: Type{Id: types.INT, Kind: types.TYPE_MAP[types.INT]},
				},
			},
			RetType: RetType{DataType: Type{Id: types.STR, Kind: types.TYPE_MAP[types.STR]}},
		},
	},
}

func readyStrDefines(s value) {
	lenVar := strDefines.Globals[0]
	lenVar.Constant = s.constant
	if lenVar.Constant {
		lenVar.ExprTag = int64(len(s.expr.(string)))
		lenVar.Expr.Model = get_const_expr_model(value{
			expr: lenVar.ExprTag,
			data: ast.Data{DataType: lenVar.DataType},
		})
	}
}

var sliceDefines = &ast.Defmap{
	Globals: []*Var{
		{
			Public:   true,
			Id:       "len",
			DataType: Type{Id: types.INT, Kind: types.TYPE_MAP[types.INT]},
			Tag:      "_len()",
		},
		{
			Public:   true,
			Id:       "cap",
			DataType: Type{Id: types.INT, Kind: types.TYPE_MAP[types.INT]},
			Tag:      "_cap()",
		},
	},
}

var arrayDefines = &ast.Defmap{
	Globals: []*Var{
		{
			Public:   true,
			Constant: true,
			Id:       "len",
			DataType: Type{Id: types.INT, Kind: types.TYPE_MAP[types.INT]},
		},
	},
}

func readyArrayDefines(s value) {
	lenVar := arrayDefines.Globals[0]
	lenVar.ExprTag = int64(s.data.DataType.Size.N)
	lenVar.Expr = s.data.DataType.Size.Expr
}

var mapDefines = &ast.Defmap{
	Globals: []*Var{
		{
			Public:   true,
			Id:       "len",
			DataType: Type{Id: types.INT, Kind: types.TYPE_MAP[types.INT]},
			Tag:      "_len()",
		},
	},
	Fns: []*Fn{
		{
			Public: true,
			Id:     "clear",
		},
		{
			Public: true,
			Id:     "keys",
		},
		{
			Public: true,
			Id:     "values",
		},
		{
			Public:  true,
			Id:      "has",
			Params:  []Param{{Id: "key"}},
			RetType: RetType{DataType: Type{Id: types.BOOL, Kind: types.TYPE_MAP[types.BOOL]}},
		},
		{
			Public: true,
			Id:     "del",
			Params: []Param{{Id: "key"}},
		},
	},
}

func readyMapDefines(mapt Type) {
	types := mapt.Tag.([]Type)
	keyt := types[0]
	valt := types[1]

	keysFunc, _, _ := mapDefines.FnById("keys", nil)
	keysFunc.RetType.DataType = keyt
	keysFunc.RetType.DataType.Kind = lexer.PREFIX_SLICE + keysFunc.RetType.DataType.Kind

	valuesFunc, _, _ := mapDefines.FnById("values", nil)
	valuesFunc.RetType.DataType = valt
	valuesFunc.RetType.DataType.Kind = lexer.PREFIX_SLICE + valuesFunc.RetType.DataType.Kind

	hasFunc, _, _ := mapDefines.FnById("has", nil)
	hasFunc.Params[0].DataType = keyt

	delFunc, _, _ := mapDefines.FnById("del", nil)
	delFunc.Params[0].DataType = keyt
}

func init() {
	out_fn.BuiltinCaller = caller_out
	outln_fn = new(Fn)
	*outln_fn = *out_fn
	outln_fn = new(ast.Fn)
	*outln_fn = *out_fn
	outln_fn.Id = "println"
	Builtin.Fns = append(Builtin.Fns, outln_fn)

	make_fn.BuiltinCaller = caller_make
	copy_fn.BuiltinCaller = caller_copy
	append_fn.BuiltinCaller = caller_append

	new_fn.BuiltinCaller = caller_new
	drop_fn.BuiltinCaller = caller_drop
	real_fn.BuiltinCaller = caller_real

	receiver := new(Var)
	receiver.Mutable = false
	for _, f := range errorTrait.Defines.Fns {
		f.Receiver = receiver
		f.Receiver.Tag = errorTrait
		f.Owner = builtinFile
	}

	intMax := intStatics.Globals[0]
	intMin := intStatics.Globals[1]
	uintMax := uintStatics.Globals[0]
	switch types.BIT_SIZE {
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

func caller_out(p *Parser, f *Fn, data call_data, m *expr_model) (v value) {
	errtok := data.args[0]
	v.data.DataType = f.RetType.DataType
	data.args = data.args[1 : len(data.args)-1]
	arg, model := p.evalToks(data.args, nil)
	if types.IsFn(arg.data.DataType) {
		p.pusherrtok(errtok, "invalid_expr")
	}
	m.append_sub(exprNode{"(" + model.String() + ")"})
	arg.constant = false
	return v
}

func caller_make(p *Parser, _ *Fn, data call_data, m *expr_model) (v value) {
	errtok := data.args[0]
	args := p.get_args(data.args, false)
	if len(args.Src) == 0 {
		p.pusherrtok(errtok, "missing_expr")
		return
	}
	type_tokens := args.Src[0].Expr.Tokens
	r := new_builder(nil)
	i := 0
	t, ok := r.DataType(type_tokens, &i, true)
	if !ok {
		p.pusherrs(r.Errors...)
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
	case types.IsSlice(t):
		return fn_make(p, m, t, args, errtok)
	default:
		p.pusherrtok(errtok, "invalid_type")
	}
	return
}

func caller_copy(p *Parser, _ *Fn, data call_data, m *expr_model) (v value) {
	v.data.DataType.Id = types.INT
	v.data.DataType.Kind = types.TYPE_MAP[v.data.DataType.Id]
	v.data.Value = " "

	errtok := data.args[0]
	args := p.get_args(data.args, false)
	if len(args.Src) < 1 {
		p.pusherrtok(errtok, "missing_expr_for", "dest")
		return
	} else if len(args.Src) < 2 {
		p.pusherrtok(errtok, "missing_expr_for", "src")
		return
	} else if len(args.Src) > 2 {
		p.pusherrtok(errtok, "argument_overflow")
	}

	dest_expr := args.Src[0].Expr
	dest_v, dest_expr_model := p.eval_expr(dest_expr, nil)
	if !types.IsSlice(dest_v.data.DataType) {
		p.pusherrtok(errtok, "invalid_type")
		return
	}

	src_expr := args.Src[1].Expr
	src_v, src_expr_model := p.eval_expr(src_expr, nil)

	if !types.Equals(dest_v.data.DataType, src_v.data.DataType) {
		t := dest_v.data.DataType.ComponentType
		if !types.IsPure(*t) || t.Id != types.U8 {
			p.pusherrtok(errtok, "incompatible_types", t.Kind, src_v.data.DataType.Kind)
			return
		}

		if !types.IsPure(src_v.data.DataType) || src_v.data.DataType.Id != types.STR {
			p.pusherrtok(errtok, "incompatible_types", t.Kind, src_v.data.DataType.Kind)
			return
		}
		m.append_sub(exprNode{"<"})
		m.append_sub(exprNode{t.String()})
		m.append_sub(exprNode{">"})
	}

	m.append_sub(exprNode{"("})
	m.append_sub(dest_expr_model)
	m.append_sub(exprNode{","})
	m.append_sub(src_expr_model)
	m.append_sub(exprNode{")"})
	return v
}

func caller_append(p *Parser, _ *Fn, data call_data, m *expr_model) (v value) {
	errtok := data.args[0]
	args := p.get_args(data.args, false)
	if len(args.Src) < 1 {
		p.pusherrtok(errtok, "missing_expr_for", "src")
		return
	} else if len(args.Src) < 2 {
		p.pusherrtok(errtok, "missing_expr_for", "values")
		return
	}

	src_expr := args.Src[0].Expr
	src_v, src_expr_model := p.eval_expr(src_expr, nil)
	if !types.IsSlice(src_v.data.DataType) {
		p.pusherrtok(errtok, "invalid_type")
		return
	}
	t := src_v.data.DataType.ComponentType
	v.data.DataType = src_v.data.DataType.Copy()
	v.data.Value = " "

	v.mutable = true

	m.append_sub(exprNode{"("})
	m.append_sub(src_expr_model)
	m.append_sub(exprNode{",{"})

	for _, arg := range args.Src[1:] {
		arg_expr := arg.Expr
		arg_v, arg_expr_model := p.eval_expr(arg_expr, nil)

		p.check_assign_type(*t, arg_v, errtok)

		m.append_sub(arg_expr_model)
		m.append_sub(exprNode{","})
	}

	m.append_sub(exprNode{"})"})
	return v
}

func caller_drop(p *Parser, _ *Fn, data call_data, m *expr_model) (v value) {
	errtok := data.args[0]
	args := p.get_args(data.args, false)
	if len(args.Src) < 1 {
		p.pusherrtok(errtok, "missing_expr_for", "ref")
		return
	} else if len(args.Src) > 1 {
		p.pusherrtok(errtok, "argument_overflow")
	}
	ref_expr := args.Src[0].Expr
	ref_v, ref_expr_model := p.eval_expr(ref_expr, nil)
	if !types.IsRef(ref_v.data.DataType) {
		p.pusherrtok(errtok, "invalid_type")
		return
	}
	v.data.DataType.Id = types.VOID
	v.data.DataType.Kind = types.TYPE_MAP[v.data.DataType.Id]
	v.data.Value = " "
	m.append_sub(exprNode{"("})
	m.append_sub(ref_expr_model)
	m.append_sub(exprNode{")"})
	return v
}

func caller_real(p *Parser, _ *Fn, data call_data, m *expr_model) (v value) {
	v.data.DataType.Id = types.BOOL
	v.data.DataType.Kind = types.TYPE_MAP[v.data.DataType.Id]

	errtok := data.args[0]
	args := p.get_args(data.args, false)
	if len(args.Src) < 1 {
		p.pusherrtok(errtok, "missing_expr_for", "ref")
		return
	} else if len(args.Src) > 1 {
		p.pusherrtok(errtok, "argument_overflow")
	}
	ref_expr := args.Src[0].Expr
	ref_v, ref_expr_model := p.eval_expr(ref_expr, nil)
	if !types.IsRef(ref_v.data.DataType) {
		p.pusherrtok(errtok, "invalid_type")
		return
	}
	v.data.Value = " "
	m.append_sub(exprNode{"("})
	m.append_sub(ref_expr_model)
	m.append_sub(exprNode{")"})
	return v
}

func caller_new(p *Parser, _ *Fn, data call_data, m *expr_model) (v value) {
	errtok := data.args[0]
	args := p.get_args(data.args, false)
	if len(args.Src) < 1 {
		p.pusherrtok(errtok, "missing_expr_for", "type")
		return
	} else if len(args.Src) > 2 {
		p.pusherrtok(errtok, "argument_overflow")
	}

	r := new_builder(nil)
	i := 0
	t, ok := r.DataType(args.Src[0].Expr.Tokens, &i, true)
	if !ok {
		p.pusherrs(r.Errors...)
		return
	}
	if i+1 < len(args.Src[0].Expr.Tokens) {
		p.pusherrtok(args.Src[0].Expr.Tokens[i+1], "invalid_syntax")
	}
	t, _ = p.realType(t, true)
	if !types.ValidForRef(t) {
		p.pusherrtok(errtok, "invalid_type")
	}
	if types.IsStruct(t) {
		s := t.Tag.(*ast.Struct)
		for _, f := range s.Defines.Globals {
			if types.IsRef(f.DataType) {
				p.pusherrtok(errtok, "ref_used_struct_used_at_new_fn")
				break
			}
		}
	}
	if len(args.Src) == 1 {
		m.append_sub(exprNode{"<" + t.String() + ">()"})
	} else {
		data_expr := args.Src[1].Expr
		data_v, data_expr_model := p.eval_expr(data_expr, nil)
		p.check_type(t, data_v.data.DataType, false, true, errtok)
		m.append_sub(exprNode{"<" + t.String() + ">("})
		m.append_sub(data_expr_model)
		m.append_sub(exprNode{")"})
	}
	t.Kind = "&" + t.Kind
	v.data.DataType = t
	v.data.Value = t.Kind
	v.mutable = true
	return
}

var std_mem_builtin = &ast.Defmap{
	Fns: []*Fn{
		{Id: "size_of", BuiltinCaller: caller_mem_size_of},
		{Id: "align_of", BuiltinCaller: caller_mem_align_of},
	},
}

func caller_mem_size_of(p *Parser, _ *Fn, data call_data, m *expr_model) (v value) {
	data.args = data.args[1 : len(data.args)-1]
	v.data.DataType = Type{
		Id:   types.UINT,
		Kind: types.TYPE_MAP[types.UINT],
	}
	nodes := m.nodes[m.index].nodes
	node := &nodes[len(nodes)-1]
	r := new_builder(nil)
	i := 0
	t, ok := r.DataType(data.args, &i, true)
	if !ok {
		v, model := p.evalToks(data.args, nil)
		*node = exprNode{"sizeof(" + model.String() + ")"}
		v.constant = false
		return v
	}
	if i+1 < len(data.args) {
		p.pusherrtok(data.args[i+1], "invalid_syntax")
	}
	p.pusherrs(r.Errors...)
	t, _ = p.realType(t, true)
	*node = exprNode{"sizeof(" + t.String() + ")"}
	return
}

func caller_mem_align_of(p *Parser, _ *Fn, data call_data, m *expr_model) (v value) {
	data.args = data.args[1 : len(data.args)-1]
	v.data.DataType = Type{
		Id:   types.UINT,
		Kind: types.TYPE_MAP[types.UINT],
	}
	nodes := m.nodes[m.index].nodes
	node := &nodes[len(nodes)-1]
	r := new_builder(nil)
	i := 0
	t, ok := r.DataType(data.args, &i, true)
	if !ok {
		v, model := p.evalToks(data.args, nil)
		*node = exprNode{"alignof(" + model.String() + ")"}
		v.constant = false
		return v
	}
	if i+1 < len(data.args) {
		p.pusherrtok(data.args[i+1], "invalid_syntax")
	}
	p.pusherrs(r.Errors...)
	t, _ = p.realType(t, true)
	*node = exprNode{"alignof(" + t.String() + ")"}
	return
}

var std_builtin_defines = map[string]*ast.Defmap{
	"std::mem": std_mem_builtin,
}

func fn_make(p *Parser, m *expr_model, t ast.Type, args *ast.Args, errtok lexer.Token) (v value) {
	v.data.DataType = t
	v.data.Value = " "
	v.mutable = true
	if len(args.Src) < 2 {
		p.pusherrtok(errtok, "missing_expr_for", "len")
		return
	} else if len(args.Src) > 2 {
		p.pusherrtok(errtok, "argument_overflow")
	}
	len_expr := args.Src[1].Expr
	len_v, len_expr_model := p.eval_expr(len_expr, nil)
	err_key := check_value_for_indexing(len_v)
	if err_key != "" {
		p.pusherrtok(errtok, err_key)
	}

	write_model := true
	if len_v.constant {
		switch len_v.expr.(type) {
		case float64:
			write_model = len_v.expr.(float64) != 0
		case int64:
			write_model = len_v.expr.(int64) != 0
		}
	}
	m.nodes[m.index].nodes[0] = nil
	m.append_sub(exprNode{t.String()})
	m.append_sub(exprNode{"("})
	if write_model {
		m.append_sub(len_expr_model)
	}
	m.append_sub(exprNode{")"})
	return
}
