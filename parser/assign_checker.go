package parser

import (
	"github.com/DeRuneLabs/jane/package/jnbits"
	"github.com/DeRuneLabs/jane/package/jntype"
)

type assignChecker struct {
	p         *Parser
	constant  bool
	t         DataType
	v         value
	ignoreAny bool
	errtok    Tok
}

func (ac assignChecker) checkAssignTypeAsync() {
	defer func() { ac.p.wg.Done() }()
	ac.p.checkAssignConst(ac.constant, ac.t, ac.v, ac.errtok)
	if typeIsPure(ac.t) && isConstNumeric(ac.v.data.Value) {
		switch {
		case jntype.IsSignedIntegerType(ac.t.Id):
			if jnbits.CheckBitInt(ac.v.data.Value, jnbits.BitsizeType(ac.t.Id)) {
				return
			}
		case jntype.IsFloatType(ac.t.Id):
			if checkFloatBit(ac.v.data, jnbits.BitsizeType(ac.t.Id)) {
				return
			}
		case jntype.IsUnsignedNumericType(ac.t.Id):
			if jnbits.CheckBitUInt(ac.v.data.Value, jnbits.BitsizeType(ac.t.Id)) {
				return
			}
		}
	}
	ac.p.wg.Add(1)
	go ac.p.checkTypeAsync(ac.t, ac.v.data.Type, ac.ignoreAny, ac.errtok)
}
