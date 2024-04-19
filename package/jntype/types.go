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

package jntype

import (
	"strconv"

	"github.com/DeRuneLabs/jane/package/jnapi"
)

var (
	// integer type code of current platform architecture,
	// equivalent to "int", but specific bit-size integer type code
	IntCode uint8
	// integer type of current platform architecture
	// equivalent to "uint", but specific bit-size integer type code
	UIntCode uint8
	// bit size of architecture
	BitSize int
)

const (
	NumericTypeStr = "<numeric>"
	NilTypeStr     = "<nil>"
	VoidTypeStr    = "<void>"
)

// return real type code of code,
// if type is "int" or "uint", set to bit-specific type code
func GetRealCode(t uint8) uint8 {
	switch t {
	case Int:
		t = IntCode
	case UInt, UIntptr:
		t = UIntCode
	}
	return t
}

// report i16 is greater or not data-type than specified type
func I16GreaterThan(t uint8) bool {
	t = GetRealCode(t)
	return t == U8
}

// report i32 is greater or not data-type than specified type
func I32GreaterThan(t uint8) bool {
	t = GetRealCode(t)
	return t == I8 || t == I16
}

// report i64 is greater or not data-type than specified type
func I64GreaterThan(t uint8) bool {
	t = GetRealCode(t)
	return t == I8 || t == I16 || t == I32
}

// report u16 is greater or not data-type than specified type
func U16GreaterThan(t uint8) bool {
	t = GetRealCode(t)
	return t == U8
}

// report u32 is greater or not data-type than specified type
func U32GreaterThan(t uint8) bool {
	t = GetRealCode(t)
	return t == U8 || t == U16
}

// report u64 is greater or not data-type than specified type
func U64GreaterThan(t uint8) bool {
	t = GetRealCode(t)
	return t == U8 || t == U16 || t == U32
}

// report f32 is greater or not data-type than specified type
func F32GreaterThan(t uint8) bool {
	return t != Any && t != F64
}

// report f64 is greater or not data-type than specified type
func F64GreaterThan(t uint8) bool {
	return t != Any
}

// report type one is greater than type two or not
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

// report i8 is compatible or not with data-type specified type
func I8CompatibleWith(t uint8) bool {
	t = GetRealCode(t)
	return t == I8
}

// report i16 is compatible or not with data-type specified type
func I16CompatibleWith(t uint8) bool {
	t = GetRealCode(t)
	return t == I8 || t == I16 || t == U8
}

// report i32 is compatible or not with data-type specified type
func I32CompatibleWith(t uint8) bool {
	t = GetRealCode(t)
	return t == I8 || t == I16 || t == I32 || t == U8 || t == U16
}

// report i64 is compatible or not with data-type specified type
func I64CompatibleWith(t uint8) bool {
	t = GetRealCode(t)
	switch t {
	case I8, I16, I32, I64, U8, U16, U32:
		return true
	default:
		return false
	}
}

// report u8 is compatible or not with data-type specified type
func U8CompatibleWith(t uint8) bool {
	t = GetRealCode(t)
	return t == U8
}

// report u16 is compatible or not with data-type specified type
func U16CompatibleWith(t uint8) bool {
	t = GetRealCode(t)
	return t == U8 || t == U16
}

// report u32 is compatible or not with data-type specified type
func U32CompatibleWith(t uint8) bool {
	t = GetRealCode(t)
	return t == U8 || t == U16 || t == U32
}

// report u64 is compatible or not with data-type specified type
func U64CompatibleWith(t uint8) bool {
	t = GetRealCode(t)
	return t == U8 || t == U16 || t == U32 || t == U64
}

// report f32 is compatible or not with data-type specified type
func F32CompatibleWith(t uint8) bool {
	t = GetRealCode(t)
	switch t {
	case F32, I8, I16, I32, I64, U8, U16, U32, U64:
		return true
	default:
		return false
	}
}

// report f64 is compatible or not with data-type specified type
func F64CompatibleWith(t uint8) bool {
	t = GetRealCode(t)
	switch t {
	case F64, F32, I8, I16, I32, I64, U8, U16, U32, U64:
		return true
	default:
		return false
	}
}

// report type one and type two is compatible or not
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

// report type is signed/unsigned integer or not
func IsInteger(t uint8) bool {
	return IsSignedInteger(t) || IsUnsignedInteger(t)
}

// report type is numeric or not
func IsNumericType(t uint8) bool {
	return IsInteger(t) || IsFloat(t)
}

// report type is float or not
func IsFloat(t uint8) bool {
	return t == F32 || t == F64
}

// report type is numeric or not
func IsNumeric(t uint8) bool {
	return IsInteger(t) || IsFloat(t)
}

// report type is float or not
func IsSignedNumeric(t uint8) bool {
	return IsSignedInteger(t) || IsFloat(t)
}

// report type is signed itneger or not
func IsSignedInteger(t uint8) bool {
	t = GetRealCode(t)
	switch t {
	case I8, I16, I32, I64, Int:
		return true
	default:
		return false
	}
}

// report type is unsigned integer or not
func IsUnsignedInteger(t uint8) bool {
	t = GetRealCode(t)
	switch t {
	case U8, U16, U32, U64, UInt, UIntptr:
		return true
	default:
		return false
	}
}

// return type id of specified type code
func TypeFromId(id string) uint8 {
	for t, tid := range TypeMap {
		if id == tid {
			return t
		}
	}
	return 0
}

// return cpp output indentifier of data-type
func CppId(t uint8) string {
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

// return default value of specified type
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

// return type code by bits
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

// return type code by bits
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

// return type code by bits
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
