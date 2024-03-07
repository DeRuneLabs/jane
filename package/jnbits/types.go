package jnbits

import (
	"strconv"

	"github.com/De-Rune/jane/package/jn"
)

func BitsizeType(t uint8) int {
	switch t {
	case jn.I8, jn.U8:
		return 8
	case jn.I16, jn.U16:
		return 16
	case jn.I32, jn.U32, jn.F32:
		return 32
	case jn.I64, jn.U64, jn.F64:
		return 64
	case jn.Size, jn.SSize:
		return strconv.IntSize
	default:
		return 0
	}
}
