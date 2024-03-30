package jnapi

import "strings"

var JNCHeader = ""

const (
	CppIgnore = "std::ignore"
	CppSelf   = "this"
)

func ToDeferredCall(expr string) string {
	var cpp strings.Builder
	cpp.WriteString("DEFER(")
	cpp.WriteString(expr)
	cpp.WriteString(");")
	return cpp.String()
}

func ToConcurrentCall(expr string) string {
	var cpp strings.Builder
	cpp.WriteString("CO(")
	cpp.WriteString(expr)
	cpp.WriteString(");")
	return cpp.String()
}
