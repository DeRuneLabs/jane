package jnapi

import "strings"

const (
	StrMark    = "u8"
	RawStrMark = StrMark + "R"
)

func ToStr(literal string) string {
	var cxx strings.Builder
	cxx.WriteString("str_jnt{")
	cxx.WriteString(StrMark)
	cxx.WriteString(literal)
	cxx.WriteByte('}')
	return cxx.String()
}

func ToRawStr(literal string) string {
	var cxx strings.Builder
	cxx.WriteString("str_jnt{")
	cxx.WriteString(RawStrMark)
	cxx.WriteString(literal)
	cxx.WriteByte('}')
	return cxx.String()
}

func ToChar(literal string) string {
	var cxx strings.Builder
	cxx.WriteString(literal)
	return cxx.String()
}
