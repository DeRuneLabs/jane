package parser

import (
	"github.com/De-Rune/jane/ast"
	"github.com/De-Rune/jane/package/jn"
	"github.com/De-Rune/jane/package/jnbits"
)

func IsString(value string) bool {
	return value[0] == '"'
}

func IsRune(value string) bool {
	return value[0] == '\''
}

func IsBoolean(value string) bool {
	return value == "true" || value == "false"
}

func IsNil(value string) bool {
	return value == "nil"
}

func isConstantNumeric(v string) bool {
	if v == "" {
		return false
	}
	return v[0] >= '0' && v[0] <= '9'
}

func checkIntBit(v ast.ValueAST, bit int) bool {
	if bit == 0 {
		return false
	}
	if jn.IsSignedNumericType(v.Type.Code) {
		return jnbits.CheckBitInt(v.Value, bit)
	}
	return jnbits.CheckBitUint(v.Value, bit)
}
