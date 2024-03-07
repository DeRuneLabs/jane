package jnbits

import "strconv"

func CheckBitFloat(val string, bit int) bool {
	_, err := strconv.ParseFloat(val, bit)
	return err == nil
}
