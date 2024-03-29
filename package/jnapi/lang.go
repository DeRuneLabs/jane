package jnapi

import "strings"

var JNCHeader = ""

const (
	CxxIgnore = "std::ignore"
	CxxSelf   = "this"
)

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
