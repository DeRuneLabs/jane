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

package transpiler

import (
	"math"
	"strconv"

	"github.com/DeRuneLabs/jane/lexer"
	"github.com/DeRuneLabs/jane/package/jnbits"
	"github.com/DeRuneLabs/jane/package/jntype"
)

func floatAssignable(dt uint8, v value) bool {
	switch t := v.expr.(type) {
	case float64:
		v.data.Value = strconv.FormatFloat(t, 'e', -1, 64)
	case int64:
		v.data.Value = strconv.FormatFloat(float64(t), 'e', -1, 64)
	case uint64:
		v.data.Value = strconv.FormatFloat(float64(t), 'e', -1, 64)
	}
	return checkFloatBit(v.data, jnbits.BitsizeType(dt))
}

func signedAssignable(dt uint8, v value) bool {
	minn := jntype.MinOfType(dt)
	maxn := int64(jntype.MaxOfType(dt))
	switch t := v.expr.(type) {
	case float64:
		i, frac := math.Modf(t)
		if frac != 0 {
			return false
		}
		return i >= float64(minn) && i <= float64(maxn)
	case uint64:
		if t <= uint64(maxn) {
			return true
		}
	case int64:
		return t >= minn && t <= maxn
	}
	return false
}

func unsignedAssignable(dt uint8, v value) bool {
	maxnum := jntype.MaxOfType(dt)
	switch t := v.expr.(type) {
	case float64:
		if t < 0 {
			return false
		}
		i, frac := math.Modf(t)
		if frac != 0 {
			return false
		}
		return i <= float64(maxnum)
	case uint64:
		if t <= maxnum {
			return true
		}
	case int64:
		if t < 0 {
			return false
		}
		return uint64(t) <= maxnum
	}
	return false
}

func integerAssignable(dt uint8, v value) bool {
	switch {
	case jntype.IsSignedInteger(dt):
		return signedAssignable(dt, v)
	case jntype.IsUnsignedInteger(dt):
		return unsignedAssignable(dt, v)
	}
	return false
}

type assign_checker struct {
	t                *Transpiler
	expr_t           Type
	v                value
	ignoreAny        bool
	not_allow_assign bool
	errtok           lexer.Token
}

func (ac *assign_checker) has_error() bool {
	return ac.t.eval.has_error || ac.v.data.Value == ""
}

func (ac *assign_checker) check_validity() (valid bool) {
	valid = true
	if typeIsFunc(ac.v.data.Type) {
		f := ac.v.data.Type.Tag.(*Func)
		if f.Receiver != nil {
			ac.t.pusherrtok(ac.errtok, "method_as_anonymous_fn")
			valid = false
		} else if len(f.Generics) > 0 {
			ac.t.pusherrtok(ac.errtok, "genericed_fn_as_anonymous_fn")
			valid = false
		}
	}
	return
}

func (ac *assign_checker) check_constant() (ok bool) {
	if !ac.v.constExpr || !typeIsPure(ac.expr_t) ||
		!typeIsPure(ac.v.data.Type) || !jntype.IsNumeric(ac.v.data.Type.Id) {
		return
	}
	ok = true
	switch {
	case jntype.IsFloat(ac.expr_t.Id):
		if !floatAssignable(ac.expr_t.Id, ac.v) {
			ac.t.pusherrtok(ac.errtok, "overflow_limits")
			ok = false
		}
	case jntype.IsInteger(ac.expr_t.Id):
		if !integerAssignable(ac.expr_t.Id, ac.v) {
			ac.t.pusherrtok(ac.errtok, "overflow_limits")
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
	} else if ac.check_constant() {
		return
	}
	ac.t.checkType(ac.expr_t, ac.v.data.Type, ac.ignoreAny, !ac.not_allow_assign, ac.errtok)
}
