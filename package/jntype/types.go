package jntype

import (
	"strconv"

	"github.com/DeRuneLabs/jane/package/jnapi"
)

var (
	IntCode  uint8
	UIntCode uint8
	BitSize  int
)

const (
	NumericTypeStr = "<numeric>"
	NilTypeStr     = "<nil>"
	VoidTypeStr    = "<void>"
)

func GetRealCode(t uint8) uint8 {
	switch t {
	case Int:
		t = IntCode
	case UInt, UIntptr:
		t = UIntCode
	}
	return t
}

func I16GreaterThan(t uint8) bool {
	t = GetRealCode(t)
	return t == U8
}

func I32GreaterThan(t uint8) bool {
	t = GetRealCode(t)
	return t == I8 || t == I16
}

func I64GreaterThan(t uint8) bool {
	t = GetRealCode(t)
	return t == I8 || t == I16 || t == I32
}

func U16GreaterThan(t uint8) bool {
	t = GetRealCode(t)
	return t == U8
}

func U32GreaterThan(t uint8) bool {
	t = GetRealCode(t)
	return t == U8 || t == U16
}

func U64GreaterThan(t uint8) bool {
	t = GetRealCode(t)
	return t == U8 || t == U16 || t == U32
}

func F32GreaterThan(t uint8) bool {
	return t != Any && t != F64
}

func F64GreaterThan(t uint8) bool {
	return t != Any
}

func TypeGreaterThan(t1, t2 uint8) bool {
	t1 = GetRealCode(t1)
	switch t1 {
	case I16:
		return I16GreaterThan(t2)
	case I32:
		return I32GreaterThan(t2)
	case I64:
		return I64GreaterThan(t2)
	case U16:
		return U16GreaterThan(t2)
	case U32:
		return U32GreaterThan(t2)
	case U64:
		return U64GreaterThan(t2)
	case F32:
		return F32GreaterThan(t2)
	case F64:
		return F64GreaterThan(t2)
	case Enum, Any:
		return true
	}
	return false
}

func I8CompatibleWith(t uint8) bool {
	t = GetRealCode(t)
	return t == I8
}

func I16CompatibleWith(t uint8) bool {
	t = GetRealCode(t)
	return t == I8 || t == I16 || t == U8
}

func I32CompatibleWith(t uint8) bool {
	t = GetRealCode(t)
	return t == I8 || t == I16 || t == I32 || t == U8 || t == U16
}

func I64CompatibleWith(t uint8) bool {
	t = GetRealCode(t)
	switch t {
	case I8, I16, I32, I64, U8, U16, U32:
		return true
	default:
		return false
	}
}

func U8CompatibleWith(t uint8) bool {
	t = GetRealCode(t)
	return t == U8
}

func U16CompatibleWith(t uint8) bool {
	t = GetRealCode(t)
	return t == U8 || t == U16
}

func U32CompatibleWith(t uint8) bool {
	t = GetRealCode(t)
	return t == U8 || t == U16 || t == U32
}

func U64CompatibleWith(t uint8) bool {
	t = GetRealCode(t)
	return t == U8 || t == U16 || t == U32 || t == U64
}

func F32CompatibleWith(t uint8) bool {
	t = GetRealCode(t)
	switch t {
	case F32, I8, I16, I32, I64, U8, U16, U32, U64:
		return true
	default:
		return false
	}
}

func F64CompatibleWith(t uint8) bool {
	t = GetRealCode(t)
	switch t {
	case F64, F32, I8, I16, I32, I64, U8, U16, U32, U64:
		return true
	default:
		return false
	}
}

func TypesAreCompatible(t1, t2 uint8, ignoreany bool) bool {
	t1 = GetRealCode(t1)
	switch t1 {
	case Any:
		return !ignoreany
	case I8:
		return I8CompatibleWith(t2)
	case I16:
		return I16CompatibleWith(t2)
	case I32:
		return I32CompatibleWith(t2)
	case I64:
		return I64CompatibleWith(t2)
	case U8:
		return U8CompatibleWith(t2)
	case U16:
		return U16CompatibleWith(t2)
	case U32:
		return U32CompatibleWith(t2)
	case U64:
		return U64CompatibleWith(t2)
	case Bool:
		return t2 == Bool
	case Str:
		return t2 == Str
	case F32:
		return F32CompatibleWith(t2)
	case F64:
		return F64CompatibleWith(t2)
	case Nil:
		return t2 == Nil
	}
	return false
}

func IsInteger(t uint8) bool {
	return IsSignedInteger(t) || IsUnsignedInteger(t)
}

func IsNumericType(t uint8) bool {
	return IsInteger(t) || IsFloat(t)
}

func IsFloat(t uint8) bool {
	return t == F32 || t == F64
}

func IsNumeric(t uint8) bool {
	return IsInteger(t) || IsFloat(t)
}

func IsSignedNumeric(t uint8) bool {
	return IsSignedInteger(t) || IsFloat(t)
}

func IsSignedInteger(t uint8) bool {
	t = GetRealCode(t)
	switch t {
	case I8, I16, I32, I64, Int:
		return true
	default:
		return false
	}
}

func IsUnsignedInteger(t uint8) bool {
	t = GetRealCode(t)
	switch t {
	case U8, U16, U32, U64, UInt, UIntptr:
		return true
	default:
		return false
	}
}

func TypeFromId(id string) uint8 {
	for t, tid := range TypeMap {
		if id == tid {
			return t
		}
	}
	return 0
}

func CxxId(t uint8) string {
	if t == Void {
		return "void"
	}
	id := TypeMap[t]
	if id == "" {
		return id
	}
	id = jnapi.AsTypeId(id)
	return id
}

func DefaultValOfType(t uint8) string {
	t = GetRealCode(t)
	if IsNumericType(t) || t == Enum {
		return "0"
	}
	switch t {
	case Bool:
		return "false"
	case Str:
		return `""`
	}
	return "nil"
}

func IntFromBits(bits uint64) uint8 {
	switch bits {
	case 8:
		return I8
	case 16:
		return I16
	case 32:
		return I32
	default:
		return I64
	}
}

func UIntFromBits(bits uint64) uint8 {
	switch bits {
	case 8:
		return U8
	case 16:
		return U16
	case 32:
		return U32
	default:
		return U64
	}
}

func FloatFromBits(bits uint64) uint8 {
	switch bits {
	case 32:
		return F32
	default:
		return F64
	}
}

func init() {
	BitSize = strconv.IntSize
	switch BitSize {
	case 8:
		IntCode = I8
		UIntCode = U8
	case 16:
		IntCode = I16
		UIntCode = U16
	case 32:
		IntCode = I32
		UIntCode = U32
	case 64:
		IntCode = I64
		UIntCode = U64
	}
}
