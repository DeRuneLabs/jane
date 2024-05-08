// copyright (c) 2024 arfy slowy - derunelabs
//
// permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "software"), to deal
// in the software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the software, and to permit persons to whom the software is
// furnished to do so, subject to the following conditions:
//
// the above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the software.
//
// the software is provided "as is", without warranty of any kind, express or
// implied, including but not limited to the warranties of merchantability,
// fitness for a particular purpose and noninfringement. in no event shall the
// authors or copyright holders be liable for any claim, damages or other
// liability, whether in an action of contract, tort or otherwise, arising from,
// out of or in connection with the software or the use or other dealings in the
// software.

package types

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
