package parser

import (
	"fmt"
	"strings"

	"github.com/De-Rune/jane/ast"
	"github.com/De-Rune/jane/lexer"
	"github.com/De-Rune/jane/package/jane"
)

const entryPointStandard = `
  // Entry point standard codes.
#if WIN32
  _setmode(0x1, 0x40000);
#else
  setmode(0x1, 0x40000);
#endif
`

type Function struct {
	Token      lexer.Token
	Name       string
	ReturnType uint8
	Block      ast.BlockAST
}

func (f Function) String() string {
	var sb strings.Builder
	sb.WriteString(cxxTypeNameFromType(f.ReturnType))
	sb.WriteByte(' ')
	sb.WriteString(f.Name)
	sb.WriteString("()")
	sb.WriteString(" {")
	sb.WriteString(getFunctionStandardClose(f.Name))
	for _, s := range f.Block.Content {
		sb.WriteByte('\n')
		sb.WriteString("\t" + fmt.Sprint(s.Value))
		sb.WriteByte(';')
	}
	sb.WriteString("\n}")
	return sb.String()
}

func getFunctionStandardClose(name string) string {
	switch name {
	case jane.EntryPoint:
		return entryPointStandard
	}
	return ""
}
