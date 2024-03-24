package jnapi

import (
	"encoding/hex"
	"strconv"
	"strings"
	"unicode/utf8"
)

const RawStrMark = "R"

func ToStr(bytes []byte) string {
	var cxx strings.Builder
	cxx.WriteString("str_jnt{\"")
	cxx.WriteString(bytesToStr(bytes))
	cxx.WriteString("\"}")
	return cxx.String()
}

func ToRawStr(bytes []byte) string {
	var cxx strings.Builder
	cxx.WriteString("str_jnt{")
	cxx.WriteString(RawStrMark)
	cxx.WriteString("\"(")
	cxx.WriteString(bytesToStr(bytes))
	cxx.WriteString(")\"}")
	return cxx.String()
}

func ToChar(b byte) string {
	return "'" + string(b) + "'"
}

func ToRune(bytes []byte) string {
	if len(bytes) == 0 {
		return ""
	} else if bytes[0] == '\\' {
		if len(bytes) > 1 && (bytes[1] == 'u' || bytes[1] == 'U') {
			bytes = bytes[2:]
			i, _ := strconv.ParseInt(string(bytes), 16, 32)
			return "0x" + strconv.FormatInt(i, 16)
		}
	}
	r, _ := utf8.DecodeRune(bytes)
	return "0x" + strconv.FormatInt(int64(r), 16)
}

func btoa(b byte) string {
	if b <= 127 {
		return string(b)
	}
	return "\\x" + hex.EncodeToString([]byte{b})
}

func bytesToStr(bytes []byte) string {
	var str strings.Builder
	for _, b := range bytes {
		str.WriteString(btoa(b))
	}
	return str.String()
}
