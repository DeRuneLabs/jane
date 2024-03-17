package jnapi

import "strings"

func ToJnAlloc(t, expr string) string {
	var cxx strings.Builder
	cxx.WriteString("jnalloc<")
	cxx.WriteString(t)
	cxx.WriteString(">(")
	cxx.WriteString(expr)
	cxx.WriteByte(')')
	return cxx.String()
}
