package jnapi

import "strings"

func ToJnAlloc(t string) string {
	var cxx strings.Builder
	cxx.WriteString("jnalloc<")
	cxx.WriteString(t)
	cxx.WriteString(">()")
	return cxx.String()
}
