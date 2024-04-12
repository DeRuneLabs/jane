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

package jnbits

import (
	"math"
	"strconv"
	"strings"
)

const MaxInt = 64

type bitChecker = func(val string, base, bit int) error

func CheckBitInt(val string, bit int) bool {
	return checkBit(val, bit, func(val string, base, bit int) error {
		_, err := strconv.ParseInt(val, base, bit)
		return err
	})
}

func CheckBitUInt(val string, bit int) bool {
	return checkBit(val, bit, func(val string, base, bit int) error {
		_, err := strconv.ParseUint(val, base, bit)
		return err
	})
}

func checkBit(val string, bit int, checker bitChecker) bool {
	var err error
	switch {
	case val == "":
		return false
	case len(val) == 1:
		return true
	case strings.HasPrefix(val, "0x"):
		err = checker(val[2:], 16, bit)
	case strings.HasPrefix(val, "0b"):
		err = checker(val[2:], 2, bit)
	case val[0] == '0':
		err = checker(val[1:], 8, bit)
	default:
		err = checker(val, 10, bit)
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
		return MaxInt
	}
}

func BitsizeUInt(x uint64) uint64 {
	switch {
	case x <= math.MaxUint8:
		return 8
	case x <= math.MaxUint16:
		return 16
	case x <= math.MaxUint32:
		return 32
	default:
		return MaxInt
	}
}
