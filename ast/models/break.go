package models

import "strings"

type Break struct {
	Tok  Tok
	Case *Case
}

func (b Break) String() string {
	if b.Case != nil {
		var cpp strings.Builder
		cpp.WriteString("goto ")
		cpp.WriteString(b.Case.Match.EndLabel())
		return cpp.String()
	}
	return "break;"
}
