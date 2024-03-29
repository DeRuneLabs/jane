package jntype

import "math"

func MinOfType(id uint8) int64 {
	if !IsInteger(id) {
		return 0
	}
	id = GetRealCode(id)
	switch id {
	case I8:
		return math.MinInt8
	case I16:
		return math.MinInt16
	case I32:
		return math.MinInt32
	case I64:
		return math.MinInt64
	}
	return 0
}

func MaxOfType(id uint8) uint64 {
	if !IsInteger(id) {
		return 0
	}
	id = GetRealCode(id)
	switch id {
	case I8:
		return math.MaxInt8
	case I16:
		return math.MaxInt16
	case I32:
		return math.MaxInt32
	case I64:
		return math.MaxInt64
	case U8:
		return math.MaxUint8
	case U16:
		return math.MaxUint16
	case U32:
		return math.MaxUint32
	case U64:
		return math.MaxUint64
	}
	return 0
}
