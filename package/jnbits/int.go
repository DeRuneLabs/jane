package jnbits

import (
	"strconv"
	"strings"
)

func CheckBitInt(val string, bit int) bool {
	var err error
	if strings.HasPrefix(val, "0x") {
		_, err = strconv.ParseInt(val[2:], 16, bit)
	} else {
		_, err = strconv.ParseInt(val, 10, bit)
	}
	return err == nil
}

func CheckBitUInt(val string, bit int) bool {
	var err error
	if strings.HasPrefix(val, "0x") {
		_, err = strconv.ParseInt(val[2:], 16, bit)
	} else {
		_, err = strconv.ParseInt(val, 10, bit)
	}
	return err == nil
}
