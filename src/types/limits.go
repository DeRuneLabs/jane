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

// Maximum positive value of 32-bit floating-points.
const MAX_F32 = 0x1p127 * (1 + (1 - 0x1p-23))

// Maximum negative value of 32-bit floating-points.
const MIN_F32 = -0x1p127 * (1 + (1 - 0x1p-23))

// Maximum positive value of 64-bit floating-points.
const MAX_F64 = 0x1p1023 * (1 + (1 - 0x1p-52))

// Maximum negative value of 64-bit floating-points.
const MIN_F64 = -0x1p1023 * (1 + (1 - 0x1p-52))

// Maximum positive value of 8-bit signed integers.
const MAX_I8 = 127

// Maximum negative value of 8-bit signed integers.
const MIN_I8 = -128

// Maximum positive value of 16-bit signed integers.
const MAX_I16 = 32767

// Maximum negative value of 16-bit signed integers.
const MIN_I16 = -32768

// Maximum positive value of 32-bit signed integers.
const MAX_I32 = 2147483647

// Maximum negative value of 32-bit signed integers.
const MIN_I32 = -2147483648

// Maximum positive value of 64-bit signed integers.
const MAX_I64 = 9223372036854775807

// Maximum negative value of 64-bit signed integers.
const MIN_I64 = -9223372036854775808

// Maximum value of 8-bit unsigned integers.
const MAX_U8 = 255

// Maximum value of 16-bit unsigned integers.
const MAX_U16 = 65535

// Maximum value of 32-bit unsigned integers.
const MAX_U32 = 4294967295

// Maximum value of 64-bit unsigned integers.
const MAX_U64 = 18446744073709551615

func Min_of(k string) float64 {
	k = Real_kind_of(k)
	switch k {
	case TypeKind_I8:
		return MIN_I8

	case TypeKind_I16:
		return MIN_I16

	case TypeKind_I32:
		return MIN_I32

	case TypeKind_I64:
		return MIN_I64

	case TypeKind_F32:
		return MIN_F32

	case TypeKind_F64:
		return MIN_F64

	default:
		return 0
	}
}

func Max_of(k string) float64 {
	k = Real_kind_of(k)
	switch k {
	case TypeKind_I8:
		return MAX_I8
	case TypeKind_I16:
		return MAX_I16
	case TypeKind_I32:
		return MAX_I32
	case TypeKind_I64:
		return MAX_I64
	case TypeKind_U8:
		return MAX_U8
	case TypeKind_U16:
		return MAX_U16
	case TypeKind_U32:
		return MAX_U32
	case TypeKind_U64:
		return MAX_U64
	case TypeKind_F32:
		return MAX_F32
	case TypeKind_F64:
		return MAX_F64
	default:
		return 0
	}
}
