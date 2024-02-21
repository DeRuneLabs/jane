package jnbits

import "github.com/De-Rune/jane/package/jane"

func BitsizeOfType(t uint8) int {
	switch t {
	case jane.Int8, jane.UInt8:
		return 8
	case jane.Int16, jane.UInt16:
		return 16
	case jane.Int32, jane.UInt32:
		return 32
	case jane.Int64, jane.UInt64:
		return 64
	}
	return 0
}
