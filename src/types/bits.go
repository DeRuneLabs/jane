// Copyright (c) 2024 arfy slowy - DeRuneLabs
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
	"math"
	"strconv"
	"strings"
)

const MAX_INT = 64

type bit_checker = func(v string, bas int, bit int) error

func BitsizeType(t uint8) int {
	switch t {
	case I8, U8:
		return 0b1000
	case I16, U16:
		return 0b00010000
	case I32, U32, F32:
		return 0b00100000
	case I64, U64, F64:
		return 0b01000000
	case UINT, INT:
		return BIT_SIZE
	default:
		return 0
	}
}

func CheckBitFloat(val string, bit int) bool {
	_, err := strconv.ParseFloat(val, bit)
	return err == nil
}

func BitsizeFloat(x float64) uint64 {
	switch {
	case x >= -math.MaxFloat32 && x <= math.MaxFloat32:
		return 32
	default:
		return 64
	}
}

func CheckBitInt(v string, bit int) bool {
	return check_bit(v, bit, func(v string, base int, bit int) error {
		_, err := strconv.ParseInt(v, base, bit)
		return err
	})
}

func check_bit(v string, bit int, checker bit_checker) bool {
	var err error
	switch {
	case v == "":
		return false
	case len(v) == 1:
		return true
	case strings.HasPrefix(v, "0x"):
		err = checker(v[2:], 16, bit)
	case strings.HasPrefix(v, "0b"):
		err = checker(v[2:], 2, bit)
	case v[0] == '0':
		err = checker(v[1:], 8, bit)
	default:
		err = checker(v, 10, bit)
	}
	return err == nil
}

func BitsizeInt(x int64) uint64 {
	switch {
	case x >= math.MinInt8 && x <= math.MaxInt8:
		return 8
	case x >= math.MinInt16 && x <= math.MaxInt16:
		return 16
	case x >= math.MinInt32 && x <= math.MaxInt32:
		return 32
	default:
		return MAX_INT
	}
}

func BitsizeUint(x uint64) uint64 {
	switch {
	case x <= math.MaxUint8:
		return 8
	case x <= math.MaxUint16:
		return 16
	case x <= math.MaxUint32:
		return 32
	default:
		return MAX_INT
	}
}
