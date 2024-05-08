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
	"strings"

	"github.com/DeRuneLabs/jane/ast"
	"github.com/DeRuneLabs/jane/types"
)

func check_value_for_indexing(v value) (err_key string) {
	switch {
	case !types.IsPure(v.data.DataType):
		return "invalid_expr"
	case !types.IsInteger(v.data.DataType.Id):
		return "invalid_expr"
	case v.constant && to_num_signed(v.expr) < 0:
		return "overflow_limits"
	default:
		return ""
	}
}

func get_indexing_expr_model(v value, i ast.ExprModel) ast.ExprModel {
	if i == nil {
		return i
	}
	if v.constant {
		return i
	}
	var model strings.Builder
	model.WriteString("static_cast<")
	model.WriteString(types.CppId(types.INT))
	model.WriteString(">(")
	model.WriteString(i.String())
	model.WriteByte(')')
	return exprNode{model.String()}
}

func is_enum_type(v value) bool {
	return v.is_type && types.IsEnum(v.data.DataType)
}

func is_bool_expr(v value) bool {
	return types.IsPure(v.data.DataType) && v.data.DataType.Id == types.BOOL
}

func can_get_ptr(v value) bool {
	if !v.lvalue || v.constant {
		return false
	}
	switch v.data.DataType.Id {
	case types.FN, types.ENUM:
		return false
	default:
		return true
	}
}

func is_struct_ins(val value) bool {
	return !val.is_type && types.IsStruct(val.data.DataType)
}

func is_trait_ins(val value) bool {
	return !val.is_type && types.IsTrait(val.data.DataType)
}

func is_foreach_iter_expr(val value) bool {
	switch {
	case types.IsSlice(val.data.DataType),
		types.IsArray(val.data.DataType),
		types.IsMap(val.data.DataType):
		return true
	case !types.IsPure(val.data.DataType):
		return false
	}
	code := val.data.DataType.Id
	return code == types.STR
}

func check_float_bit(v ast.Data, bit int) bool {
	if bit == 0 {
		return false
	}
	return types.CheckBitFloat(v.Value, bit)
}

func is_valid_for_const(v value) bool {
	return v.constant
}

func is_ok_for_shifting(v value) bool {
	if !types.IsPure(v.data.DataType) ||
		!types.IsInteger(v.data.DataType.Id) {
		return false
	}
	if !v.constant {
		return true
	}
	switch t := v.expr.(type) {
	case int64:
		return t >= 0
	case uint64:
		return true
	}
	return false
}
