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

import (
	"strconv"
	"strings"
)

type bit_checker = func(v string, base int, bit int) bool

const BIT_SIZE = 32 << (^uint(0) >> 63)

var (
	SYS_INT  string
	SYS_UINT string
)

func check_bit(v string, bit int, checker bit_checker) bool {
	switch {
	case v == "":
		return false
	case len(v) == 1:
		return true
	case strings.HasPrefix(v, "0x"):
		return checker(v[2:], 0b00010000, bit)
	case strings.HasPrefix(v, "0b"):
		return checker(v[2:], 0b10, bit)
	case v[0] == '0':
		return checker(v[1:], 0b1000, bit)
	default:
		return checker(v, 0b1010, bit)
	}
}

func Real_kind_of(kind string) string {
	switch kind {
	case TypeKind_INT:
		return SYS_INT
	case TypeKind_UINT, TypeKind_UINTPTR:
		return SYS_UINT
	default:
		return kind
	}
}

func Bitsize_of(k string) int {
	switch k {
	case TypeKind_I8, TypeKind_U8:
		return 0b1000
	case TypeKind_I16, TypeKind_U16:
		return 0b00010000
	case TypeKind_I32, TypeKind_U32, TypeKind_F32:
		return 0b00100000
	case TypeKind_I64, TypeKind_U64, TypeKind_F64:
		return 0b01000000
	case TypeKind_UINT, TypeKind_INT:
		return BIT_SIZE
	default:
		return -1
	}
}

func Int_from_bits(bits uint64) string {
	switch bits {
	case 0b1000:
		return TypeKind_I8
	case 0b00010000:
		return TypeKind_I16
	case 0b00100000:
		return TypeKind_I32
	case 0b01000000:
		return TypeKind_I64

	default:
		return ""
	}
}

func Uint_from_bits(bits uint64) string {
	switch bits {
	case 0b1000:
		return TypeKind_U8
	case 0b00010000:
		return TypeKind_U16
	case 0b00100000:
		return TypeKind_U32
	case 0b01000000:
		return TypeKind_U64
	default:
		return ""
	}
}

func Float_from_bits(bits uint64) string {
	switch bits {
	case 0b00100000:
		return TypeKind_F32
	case 0b01000000:
		return TypeKind_F64
	default:
		return ""
	}
}

func Check_bit_int(v string, bit int) bool {
	return check_bit(v, bit, func(v string, base int, bit int) bool {
		_, err := strconv.ParseInt(v, base, bit)
		return err == nil
	})
}

func Check_bit_uint(v string, bit int) bool {
	return check_bit(v, bit, func(v string, base int, bit int) bool {
		_, err := strconv.ParseUint(v, base, bit)
		return err == nil
	})
}

func Check_bit_float(val string, bit int) bool {
	_, err := strconv.ParseFloat(val, bit)
	return err == nil
}

func Bitsize_of_float(x float64) uint64 {
	switch {
	case MIN_F32 <= x && x <= MAX_F32:
		return 0b00100000
	default:
		return 0b01000000
	}
}

func Bitsize_of_int(x int64) uint64 {
	switch {
	case MIN_I8 <= x && x <= MAX_I8:
		return 0b1000
	case MIN_I16 <= x && x <= MAX_I16:
		return 0b00010000
	case MIN_I32 <= x && x <= MAX_I32:
		return 0b00100000
	default:
		return 0b01000000
	}
}

func Bitsize_of_uint(x uint64) uint64 {
	switch {
	case x <= MAX_U8:
		return 0b1000
	case x <= MAX_U16:
		return 0b00010000
	case x <= MAX_U32:
		return 0b00100000
	default:
		return 0b01000000
	}
}

func init() {
	switch BIT_SIZE {
	case 0b00100000:
		SYS_INT = TypeKind_I32
		SYS_UINT = TypeKind_U32
	case 0b01000000:
		SYS_INT = TypeKind_I64
		SYS_UINT = TypeKind_U64
	}
}
