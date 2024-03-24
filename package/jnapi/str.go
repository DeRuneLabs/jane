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
	return strconv.Itoa(int(b))
}

func ToRune(bytes []byte) string {
	r, _ := utf8.DecodeRune(bytes)
	return strconv.FormatInt(int64(r), 10)
}

func bytesToStr(bytes []byte) string {
	var str strings.Builder
	for _, b := range bytes {
		if b <= 127 {
			str.WriteByte(b)
		} else {
			str.WriteString("\\x")
			str.WriteString(hex.EncodeToString([]byte{b}))
		}
	}
	return str.String()
}
