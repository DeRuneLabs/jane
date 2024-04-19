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
	"strings"

	"github.com/DeRuneLabs/jane/ast/models"
	"github.com/DeRuneLabs/jane/lexer/tokens"
	"github.com/DeRuneLabs/jane/package/jnbits"
	"github.com/DeRuneLabs/jane/package/jntype"
)

func check_value_for_indexing(v value) (err_key string) {
	switch {
	case !typeIsPure(v.data.Type):
		return "invalid_expr"
	case !jntype.IsInteger(v.data.Type.Id):
		return "invalid_expr"
	case v.constExpr && tonums(v.expr) < 0:
		return "overflow_limits"
	default:
		return ""
	}
}

func indexingExprModel(i iExpr) iExpr {
	if i == nil {
		return i
	}
	var model strings.Builder
	model.WriteString("static_cast<")
	model.WriteString(jntype.CppId(jntype.Int))
	model.WriteString(">(")
	model.WriteString(i.String())
	model.WriteByte(')')
	return exprNode{model.String()}
}

func isstr(s string) bool {
	return s != "" && (s[0] == '"' || israwstr(s))
}

func israwstr(s string) bool {
	return s != "" && s[0] == '`'
}

func ischar(s string) bool {
	return s != "" && s[0] == '\''
}

func isnil(s string) bool {
	return s == tokens.NIL
}

func isbool(s string) bool {
	return s == tokens.TRUE || s == tokens.FALSE
}

func valIsEnumType(v value) bool {
	return v.is_type && typeIsEnum(v.data.Type)
}

func isBoolExpr(v value) bool {
	return typeIsPure(v.data.Type) && v.data.Type.Id == jntype.Bool
}

func isfloat(s string) bool {
	if strings.HasPrefix(s, "0x") {
		return strings.ContainsAny(s, ".pP")
	}
	return strings.ContainsAny(s, ".eE")
}

func canGetPtr(v value) bool {
	if !v.lvalue || v.constExpr {
		return false
	}
	switch v.data.Type.Id {
	case jntype.Fn, jntype.Enum:
		return false
	default:
		return v.data.Token.Id == tokens.Id
	}
}

func valIsStructIns(val value) bool {
	return !val.is_type && typeIsStruct(val.data.Type)
}

func valIsTraitIns(val value) bool {
	return !val.is_type && typeIsTrait(val.data.Type)
}

func isForeachIterExpr(val value) bool {
	switch {
	case typeIsSlice(val.data.Type),
		typeIsArray(val.data.Type),
		typeIsMap(val.data.Type):
		return true
	case !typeIsPure(val.data.Type):
		return false
	}
	code := val.data.Type.Id
	return code == jntype.Str
}

func isConstNumeric(v string) bool {
	if v == "" {
		return false
	}
	return v[0] == '-' || (v[0] >= '0' && v[0] <= '9')
}

func isConstExpression(v string) bool {
	return isConstNumeric(v) || isstr(v) || ischar(v) || isnil(v) || isbool(v)
}

func checkFloatBit(v models.Data, bit int) bool {
	if bit == 0 {
		return false
	}
	return jnbits.CheckBitFloat(v.Value, bit)
}

func validExprForConst(v value) bool {
	return v.constExpr
}

func okForShifting(v value) bool {
	if !typeIsPure(v.data.Type) ||
		!jntype.IsInteger(v.data.Type.Id) {
		return false
	}
	if !v.constExpr {
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
