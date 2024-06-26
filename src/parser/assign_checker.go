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

package parser

import (
	"math"
	"strconv"

	"github.com/DeRuneLabs/jane/lexer"
	"github.com/DeRuneLabs/jane/types"
)

func float_assignable(dt uint8, v value) bool {
	switch t := v.expr.(type) {
	case float64:
		v.data.Value = strconv.FormatFloat(t, 'e', -1, 64)
	case int64:
		v.data.Value = strconv.FormatFloat(float64(t), 'e', -1, 64)
	case uint64:
		v.data.Value = strconv.FormatFloat(float64(t), 'e', -1, 64)
	}
	return check_float_bit(v.data, types.BitsizeType(dt))
}

func signed_assignable(dt uint8, v value) bool {
	min := types.MinOfType(dt)
	max := int64(types.MaxOfType(dt))
	switch t := v.expr.(type) {
	case float64:
		i, frac := math.Modf(t)
		if frac != 0 {
			return false
		}
		return i >= float64(min) && i <= float64(max)
	case uint64:
		if t <= uint64(max) {
			return true
		}
	case int64:
		return t >= min && t <= max
	}
	return false
}

func unsigned_assignable(dt uint8, v value) bool {
	max := types.MaxOfType(dt)
	switch t := v.expr.(type) {
	case float64:
		if t < 0 {
			return false
		}
		i, frac := math.Modf(t)
		if frac != 0 {
			return false
		}
		return i <= float64(max)
	case uint64:
		if t <= max {
			return true
		}
	case int64:
		if t < 0 {
			return false
		}
		return uint64(t) <= max
	}
	return false
}

func int_assignable(dt uint8, v value) bool {
	switch {
	case types.IsSignedInteger(dt):
		return signed_assignable(dt, v)
	case types.IsUnsignedInteger(dt):
		return unsigned_assignable(dt, v)
	}
	return false
}

type assign_checker struct {
	p                *Parser
	t                Type
	v                value
	ignoreAny        bool
	not_allow_assign bool
	errtok           lexer.Token
}

func (ac *assign_checker) has_error() bool {
	return ac.p.eval.has_error || ac.v.data.Value == ""
}

func (ac *assign_checker) check_validity() (valid bool) {
	valid = true
	if types.IsFn(ac.v.data.DataType) {
		f := ac.v.data.DataType.Tag.(*Fn)
		if f.Receiver != nil {
			ac.p.pusherrtok(ac.errtok, "method_as_anonymous_fn")
			valid = false
		} else if len(f.Generics) > 0 {
			ac.p.pusherrtok(ac.errtok, "genericed_fn_as_anonymous_fn")
			valid = false
		}
	}
	return
}

func (ac *assign_checker) check_const() (ok bool) {
	if !ac.v.constant || !types.IsPure(ac.t) ||
		!types.IsPure(ac.v.data.DataType) || !types.IsNumeric(ac.v.data.DataType.Id) {
		return
	}
	ok = true
	switch {
	case types.IsFloat(ac.t.Id):
		if !float_assignable(ac.t.Id, ac.v) {
			ac.p.pusherrtok(ac.errtok, "overflow_limits")
			ok = false
		}
	case types.IsInteger(ac.t.Id):
		if !int_assignable(ac.t.Id, ac.v) {
			ac.p.pusherrtok(ac.errtok, "overflow_limits")
			ok = false
		}
	default:
		ok = false
	}
	return
}

func (ac assign_checker) check() {
	if ac.has_error() {
		return
	} else if !ac.check_validity() {
		return
	} else if ac.check_const() {
		return
	}
	ac.p.check_type(ac.t, ac.v.data.DataType, ac.ignoreAny, !ac.not_allow_assign, ac.errtok)
}
