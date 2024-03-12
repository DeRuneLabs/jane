package jnapi

import "strings"

const CxxIgnore = "std::ignore"

func ToDefer(expr string) string {
	var cxx strings.Builder
	cxx.WriteString("DEFER(")
	cxx.WriteString(expr)
	cxx.WriteString(");")
	return cxx.String()
}
