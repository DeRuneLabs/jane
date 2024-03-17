package jnapi

import "strings"

const CxxIgnore = "std::ignore"

func ToDeferredCall(expr string) string {
	var cxx strings.Builder
	cxx.WriteString("DEFER(")
	cxx.WriteString(expr)
	cxx.WriteString(");")
	return cxx.String()
}

func ToConcurrentCall(expr string) string {
	var cxx strings.Builder
	cxx.WriteString("CO(")
	cxx.WriteString(expr)
	cxx.WriteString(");")
	return cxx.String()
}
