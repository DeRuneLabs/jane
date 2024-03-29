package jnbits

import (
	"math"
	"strconv"
)

func CheckBitFloat(val string, bit int) bool {
	_, err := strconv.ParseFloat(val, bit)
	return err == nil
}

func BitsizeFloat(x float64) uint64 {
	switch {
	case x >= -math.MaxFloat32 && x <= math.MaxFloat32:
		return 32
	default:
		return 64
	}
}
