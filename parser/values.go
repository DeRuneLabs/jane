package parser

import (
	"github.com/DeRuneLabs/jane/ast"
	"github.com/DeRuneLabs/jane/lexer/tokens"
	"github.com/DeRuneLabs/jane/package/jnbits"
	"github.com/DeRuneLabs/jane/package/jntype"
)

func isstr(value string) bool {
	return value[0] == '"' || israwstr(value)
}

func israwstr(value string) bool {
	return value[0] == '`'
}

func ischar(value string) bool {
	return value[0] == '\''
}

func isnil(value string) bool {
	return value == tokens.NIL
}

func isbool(value string) bool {
	return value == tokens.TRUE || value == tokens.FALSE
}

func isBoolExpr(val value) bool {
	switch {
	case typeIsNilCompatible(val.ast.Type):
		return true
	case val.ast.Type.Id == jntype.Bool && typeIsSingle(val.ast.Type):
		return true
	}
	return false
}

func isForeachIterExpr(val value) bool {
	switch {
	case typeIsArray(val.ast.Type),
		typeIsMap(val.ast.Type):
		return true
	case !typeIsSingle(val.ast.Type):
		return false
	}
	code := val.ast.Type.Id
	return code == jntype.Str
}

func isConstNum(v string) bool {
	if v == "" {
		return false
	}
	return v[0] == '-' || (v[0] >= '0' && v[0] <= '9')
}

func isConstExpr(v string) bool {
	return isConstNum(v) || isstr(v) || ischar(v) || isnil(v) || isbool(v)
}

func checkIntBit(v ast.Value, bit int) bool {
	if bit == 0 {
		return false
	}
	if jntype.IsSignedNumericType(v.Type.Id) {
		return jnbits.CheckBitInt(v.Data, bit)
	}
	return jnbits.CheckBitUInt(v.Data, bit)
}

func checkFloatBit(v ast.Value, bit int) bool {
	if bit == 0 {
		return false
	}
	return jnbits.CheckBitFloat(v.Data, bit)
}

func defaultValueOfType(t DataType) string {
	if typeIsNilCompatible(t) {
		return tokens.NIL
	}
	return jntype.DefaultValOfType(t.Id)
}
