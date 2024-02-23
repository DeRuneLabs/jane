package parser

import (
	"github.com/De-Rune/jane/ast"
	"github.com/De-Rune/jane/lexer"
	"github.com/De-Rune/jane/package/jane"
	"strings"
)

type variable struct {
	Name  string
	Token lexer.Token
	Value ast.ExpressionAST
	Type  uint8
}

func (v variable) String() string {
	var sb strings.Builder
	sb.WriteString(jane.CxxTypeNameFromType(v.Type))
	sb.WriteByte(' ')
	sb.WriteString(v.Name)
	sb.WriteByte('=')
	sb.WriteString(v.Value.String())
	sb.WriteByte(';')
	return sb.String()
}
