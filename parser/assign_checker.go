package parser

import (
	"strconv"

	"github.com/DeRuneLabs/jane/package/jnbits"
	"github.com/DeRuneLabs/jane/package/jntype"
)

func floatAssignable(dt DataType, v value) bool {
	switch t := v.expr.(type) {
	case float64:
		v.data.Value = strconv.FormatFloat(t, 'e', -1, 64)
	case int64:
		v.data.Value = strconv.FormatFloat(float64(t), 'e', -1, 64)
	case uint64:
		v.data.Value = strconv.FormatFloat(float64(t), 'e', -1, 64)
	}
	return checkFloatBit(v.data, jnbits.BitsizeType(dt.Id))
}

func signedAssignable(dt DataType, v value) bool {
	min := jntype.MinOfType(dt.Id)
	max := int64(jntype.MaxOfType(dt.Id))
	switch t := v.expr.(type) {
	case float64:
	case uint64:
		if t <= uint64(max) {
			return true
		}
	case int64:
		return t >= min && t <= max
	}
	return false
}

func unsignedAssignable(dt DataType, v value) bool {
	max := jntype.MaxOfType(dt.Id)
	switch t := v.expr.(type) {
	case float64:
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

func integerAssignable(dt DataType, v value) bool {
	switch {
	case jntype.IsSignedInteger(dt.Id):
		return signedAssignable(dt, v)
	case jntype.IsUnsignedInteger(dt.Id):
		return unsignedAssignable(dt, v)
	}
	return false
}

type assignChecker struct {
	p         *Parser
	t         DataType
	v         value
	ignoreAny bool
	errtok    Tok
}

func (ac assignChecker) checkAssignType() {
	if ac.p.eval.hasError || ac.v.data.Value == "" {
		return
	}
	if typeIsPure(ac.t) && ac.v.constExpr && typeIsPure(ac.v.data.Type) {
		switch {
		case jntype.IsFloat(ac.t.Id):
			if !floatAssignable(ac.t, ac.v) {
				ac.p.pusherrtok(ac.errtok, "overflow_limits")
			}
			return
		case jntype.IsInteger(ac.t.Id) && jntype.IsInteger(ac.v.data.Type.Id):
			if !integerAssignable(ac.t, ac.v) {
				ac.p.pusherrtok(ac.errtok, "overflow_limits")
			}
			return
		}
	}
	ac.p.checkType(ac.t, ac.v.data.Type, ac.ignoreAny, ac.errtok)
}
