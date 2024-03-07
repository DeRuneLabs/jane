package parser

import (
	"github.com/De-Rune/jane/ast"
	"github.com/De-Rune/jane/package/jn"
	"github.com/De-Rune/jane/package/jnbits"
)

func isstr(value string) bool {
	return value[0] == '"' || israwstr(value)
}

func israwstr(value string) bool {
	return value[0] == '`'
}

func isrune(value string) bool {
	return value[0] == '\''
}

func isnil(value string) bool {
	return value == "nil"
}

func isbool(value string) bool {
	return value == "true" || value == "false"
}

func isBoolExpr(val value) bool {
	switch {
	case typeIsNilCompatible(val.ast.Type):
		return true
	case val.ast.Type.Id == jn.Bool && typeIsSingle(val.ast.Type):
		return true
	}
	return false
}

func isForeachIterExpr(val value) bool {
	switch {
	case typeIsArray(val.ast.Type):
		return true
	case !typeIsSingle(val.ast.Type):
		return false
	}
	code := val.ast.Type.Id
	return code == jn.Str
}

func isConstNum(v string) bool {
	if v == "" {
		return false
	}
	return v[0] >= '0' && v[0] <= '9'
}

func checkIntBit(v ast.Value, bit int) bool {
	if bit == 0 {
		return false
	}
	if jn.IsSignedNumericType(v.Type.Id) {
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
