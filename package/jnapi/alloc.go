package jnapi

import "strings"

func ToJnAlloc(t string) string {
	var cxx strings.Builder
	cxx.WriteString("JNALLOC(")
	cxx.WriteString(t)
	cxx.WriteByte(')')
	return cxx.String()
}
