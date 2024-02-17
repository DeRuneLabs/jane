package parser

import (
	"fmt"
	"strings"

	"github.com/De-Rune/jane/ast"
	"github.com/De-Rune/jane/lexer"
	"github.com/De-Rune/jane/package/jane"
)

type CxxParser struct {
	Functions []*Function

	Tokens []lexer.Token
	PFI    *ParseFileInfo
}

func NewParser(tokens []lexer.Token, PFI *ParseFileInfo) *CxxParser {
	parser := new(CxxParser)
	parser.Tokens = tokens
	parser.PFI = PFI
	return parser
}

func (cp *CxxParser) PushErrorToken(token lexer.Token, err string) {
	message := jane.Errors[err]
	cp.PFI.Errors = append(cp.PFI.Errors, fmt.Sprintf("%s:%d %s", token.File.Path, token.Line, message))
}

func (cp *CxxParser) PushError(err string) {
	cp.PFI.Errors = append(cp.PFI.Errors, jane.Errors[err])
}

func (cp CxxParser) String() string {
	var sb strings.Builder
	for _, function := range cp.Functions {
		sb.WriteString(function.String())
	}
	return sb.String()
}

func (cp *CxxParser) Parse() {
	astModel := ast.New(cp.Tokens)
	astModel.Build()
	if astModel.Errors != nil {
		cp.PFI.Errors = append(cp.PFI.Errors, astModel.Errors...)
		return
	}
	for _, model := range astModel.Tree {
		switch model.Type {
		case ast.Statement:
			cp.ParseStatement(model.Value.(ast.StatementAST))
		default:
			cp.PushErrorToken(model.Token, "invalid_syntax")
		}
	}
}

func (cp *CxxParser) ParseStatement(s ast.StatementAST) {
	switch s.Type {
	case ast.StatementFunction:
		cp.ParseFunction(s.Value.(ast.FunctionAST))
	default:
		cp.PushErrorToken(s.Token, "invalid_syntax")
	}
}

func (cp *CxxParser) ParseFunction(f ast.FunctionAST) {
	if function := cp.functionByName(f.Name); function != nil {
		cp.PushErrorToken(f.Token, "exist_name")
		return
	}
	function := new(Function)
	function.Name = f.Name
	function.Line = f.Token.Line
	function.FILE = f.Token.File
	function.ReturnType = f.ReturnType.Type
	cp.Functions = append(cp.Functions, function)
}

func (cp *CxxParser) functionByName(name string) *Function {
	for _, function := range cp.Functions {
		if function.Name == name {
			return function
		}
	}
	return nil
}

func (cp *CxxParser) finalChek() {
	if cp.functionByName(jane.EntryPoint) == nil {
		cp.PushError("no_entry_point")
	}
}
