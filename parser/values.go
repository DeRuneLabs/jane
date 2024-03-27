package parser

import (
	"strings"

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

func valIsStructType(v value) bool {
	return v.isType && typeIsStruct(v.data.Type)
}

func valIsEnumType(v value) bool {
	return v.isType && typeIsEnum(v.data.Type)
}

func isBoolExpr(val value) bool {
	switch {
	case typeIsNilCompatible(val.data.Type):
		return true
	case val.data.Type.Id == jntype.Bool && typeIsPure(val.data.Type):
		return true
	}
	return false
}

func isfloat(s string) bool {
	if strings.HasPrefix(s, "0x") {
		return false
	}
	return strings.Contains(s, tokens.DOT) || strings.ContainsAny(s, "eE")
}

func isForeachIterExpr(val value) bool {
	switch {
	case typeIsArray(val.data.Type),
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

func checkIntBit(v models.Data, bit int) bool {
	if bit == 0 {
		return false
	}
	if jntype.IsSignedNumericType(v.Type.Id) {
		return jnbits.CheckBitInt(v.Value, bit)
	}
	return jnbits.CheckBitUInt(v.Value, bit)
}

func checkFloatBit(v models.Data, bit int) bool {
	if bit == 0 {
		return false
	}
	return jnbits.CheckBitFloat(v.Value, bit)
}

func defaultValueOfType(t DataType) string {
	if typeIsNilCompatible(t) {
		return tokens.NIL
	}
	return jntype.DefaultValOfType(t.Id)
}
