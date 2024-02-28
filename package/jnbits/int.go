package jnbits

import (
	"strconv"
	"strings"
)

func CheckBitInt(value string, bit int) bool {
	value = strings.TrimPrefix(value, "0x")
	_, err := strconv.ParseInt(value, 16, bit)
	return err == nil
}

func CheckBitUint(value string, bit int) bool {
	value = strings.TrimPrefix(value, "0x")
	_, err := strconv.ParseUint(value, 16, bit)
	return err == nil
}
