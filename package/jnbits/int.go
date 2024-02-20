package jnbits

import "strconv"

func CheckBitInt(value string, bit int) bool {
	_, err := strconv.ParseInt(value, 10, bit)
	return err == nil
}
