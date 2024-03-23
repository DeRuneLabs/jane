package parser

import (
	"github.com/DeRuneLabs/jane/ast/models"
	"github.com/DeRuneLabs/jane/lexer/tokens"
	"github.com/DeRuneLabs/jane/package/jnbits"
	"github.com/DeRuneLabs/jane/package/jntype"
)

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

func checkIntBit(v models.Value, bit int) bool {
	if bit == 0 {
		return false
	}
	if jntype.IsSignedNumericType(v.Type.Id) {
		return jnbits.CheckBitInt(v.Data, bit)
	}
	return jnbits.CheckBitUInt(v.Data, bit)
}

func checkFloatBit(v models.Value, bit int) bool {
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
