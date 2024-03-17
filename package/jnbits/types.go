package jnbits

import "github.com/DeRuneLabs/jane/package/jntype"

func BitsizeType(t uint8) int {
	switch t {
	case jntype.I8, jntype.U8:
		return 0b1000
	case jntype.I16, jntype.U16:
		return 0b00010000
	case jntype.I32, jntype.U32, jntype.F32:
		return 0b00100000
	case jntype.I64, jntype.U64, jntype.F64:
		return 0b01000000
	case jntype.UInt, jntype.Int:
		return jntype.BitSize
	default:
		return 0
	}
}
