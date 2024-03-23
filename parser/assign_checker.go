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
	if typeIsSingle(ac.t) && isConstNum(ac.v.ast.Data) {
		switch {
		case jntype.IsSignedIntegerType(ac.t.Id):
			if jnbits.CheckBitInt(ac.v.ast.Data, jnbits.BitsizeType(ac.t.Id)) {
				return
			}
		case jntype.IsFloatType(ac.t.Id):
			if checkFloatBit(ac.v.ast, jnbits.BitsizeType(ac.t.Id)) {
				return
			}
		case jntype.IsUnsignedNumericType(ac.t.Id):
			if jnbits.CheckBitUInt(ac.v.ast.Data, jnbits.BitsizeType(ac.t.Id)) {
				return
			}
		}
	}
	ac.p.wg.Add(1)
	go ac.p.checkTypeAsync(ac.t, ac.v.ast.Type, ac.ignoreAny, ac.errtok)
}
