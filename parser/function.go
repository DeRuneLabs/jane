package parser

import (
	"github.com/De-Rune/jane/package/io"
	"strings"
)

type Function struct {
	FILE       *io.FILE
	Line       int
	Name       string
	ReturnType uint
}

func (f Function) String() string {
	var sb strings.Builder
	sb.WriteString(cxxTypeNameFromType(f.ReturnType))
	sb.WriteByte(' ')
	sb.WriteString(f.Name)
	sb.WriteString("()")
	sb.WriteString(" {")
	sb.WriteByte('}')
	return sb.String()
}
