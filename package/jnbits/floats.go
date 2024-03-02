package jnbits

import "strconv"

func CheckBitFloat(value string, bit int) bool {
	_, err := strconv.ParseFloat(value, bit)
	return err == nil
}
