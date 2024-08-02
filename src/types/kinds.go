// MIT License
//
// # Copyright (c) 2024 arfy slowy - DeRuneLabs
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
package types

// Kind of signed 8-bit integer.
const TypeKind_I8 = "i8"

// Kind of signed 16-bit integer.
const TypeKind_I16 = "i16"

// Kind of signed 32-bit integer.
const TypeKind_I32 = "i32"

// Kind of signed 64-bit integer.
const TypeKind_I64 = "i64"

// Kind of unsigned 8-bit integer.
const TypeKind_U8 = "u8"

// Kind of unsigned 16-bit integer.
const TypeKind_U16 = "u16"

// Kind of unsigned 32-bit integer.
const TypeKind_U32 = "u32"

// Kind of unsigned 64-bit integer.
const TypeKind_U64 = "u64"

// Kind of 32-bit floating-point.
const TypeKind_F32 = "f32"

// Kind of 64-bit floating-point.
const TypeKind_F64 = "f64"

// Kind of system specific bit-size unsigned integer.
const TypeKind_UINT = "uint"

// Kind of system specific bit-size signed integer.
const TypeKind_INT = "int"

// Kind of system specific bit-size unsigned integer.
const TypeKind_UINTPTR = "uintptr"

// Kind of boolean.
const TypeKind_BOOL = "bool"

// Kind of string.
const TypeKind_STR = "str"

// Kind of any type.
const TypeKind_ANY = "any"

func Is_sig_int(k string) bool {
	k = Real_kind_of(k)
	switch k {
	case TypeKind_I8, TypeKind_I16, TypeKind_I32, TypeKind_I64:
		return true
	default:
		return false
	}
}

func Is_unsig_int(k string) bool {
	k = Real_kind_of(k)
	switch k {
	case TypeKind_U8, TypeKind_U16, TypeKind_U32, TypeKind_U64:
		return true
	default:
		return false
	}
}

func Is_int(k string) bool {
	return Is_sig_int(k) || Is_unsig_int(k)
}

func Is_float(k string) bool {
	return k == TypeKind_F32 || k == TypeKind_F64
}

func Is_num(k string) bool {
	return Is_int(k) || Is_float(k)
}

func Is_sig_num(k string) bool {
	return Is_sig_int(k) || Is_float(k)
}
