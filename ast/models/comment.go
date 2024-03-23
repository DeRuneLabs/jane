package models

import "strings"

type Comment struct {
	Content string
}

func (c Comment) String() string {
	var cxx strings.Builder
	cxx.WriteString("// ")
	cxx.WriteString(c.Content)
	return cxx.String()
}
