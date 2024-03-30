package models

import "strings"

type Comment struct {
	Content string
}

func (c Comment) String() string {
	var cpp strings.Builder
	cpp.WriteString("// ")
	cpp.WriteString(c.Content)
	return cpp.String()
}
