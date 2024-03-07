package jnapi

import "strings"

const RawStrMark = "LR"
const StrMark = "L"

func ToStr(literal string) string {
	var cxx strings.Builder
	cxx.WriteString("str(")
	cxx.WriteString(StrMark)
	cxx.WriteString(literal)
	cxx.WriteByte(')')
	return cxx.String()
}

func ToRawStr(literal string) string {
	var cxx strings.Builder
	cxx.WriteString("str(")
	cxx.WriteString(RawStrMark)
	cxx.WriteString(literal)
	cxx.WriteByte(')')
	return cxx.String()
}

func ToRune(literal string) string {
	var cxx strings.Builder
	cxx.WriteString(StrMark)
	cxx.WriteString(literal)
	return cxx.String()
}
