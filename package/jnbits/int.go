package jnbits

import "strconv"

func CheckBitInt(value string, bit int) bool {
	_, err := strconv.ParseInt(value, 10, bit)
	return err == nil
}

func CheckBitUint(value string, bit int) bool {
	_, err := strconv.ParseUint(value, 10, bit)
	return err == nil
}
