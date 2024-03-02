package jnbits

import (
	"strconv"
	"strings"
)

func CheckBitInt(value string, bit int) bool {
	var err error
	if strings.HasPrefix(value, "0x") {
		_, err = strconv.ParseInt(value[2:], 16, bit)
	} else {
		_, err = strconv.ParseInt(value, 10, bit)
	}
	return err == nil
}

func CheckBitUint(value string, bit int) bool {
	var err error
	if strings.HasPrefix(value, "0x") {
		_, err = strconv.ParseUint(value[2:], 16, bit)
	} else {
		_, err = strconv.ParseUint(value, 10, bit)
	}
	return err == nil
}
