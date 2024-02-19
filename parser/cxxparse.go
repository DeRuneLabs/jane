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
		sb.WriteString("\n\n")
	}
	return sb.String()
}

func (cp *CxxParser) Parse() {
	astModel := ast.New(cp.Tokens)
	if astModel.Errors != nil {
		cp.PFI.Errors = append(cp.PFI.Errors, astModel.Errors...)
	}
	for _, model := range astModel.Tree {
		switch model.Type {
		case ast.Statement:
			cp.ParseStatement(model.Value.(ast.StatementAST))
		default:
			cp.PushErrorToken(model.Token, "invalid_syntax")
		}
	}
	cp.finalCheck()
}

func (cp *CxxParser) ParseStatement(s ast.StatementAST) {
	switch s.Type {
	case ast.StatementFunction:
		cp.ParseFunction(s.Value.(ast.FunctionAST))
	default:
		cp.PushErrorToken(s.Token, "invalid_syntax")
	}
}

func (cp *CxxParser) ParseFunction(fnAst ast.FunctionAST) {
	if function := cp.functionByBName(fnAst.Name); function != nil {
		cp.PushErrorToken(fnAst.Token, "exist_name")
		return
	}
	fn := new(Function)
	fn.Token = fnAst.Token
	fn.Name = fnAst.Name
	fn.ReturnType = fnAst.ReturnType.Type
	fn.Block = fnAst.Block
	cp.checkFunctionReturn(fn)
	cp.Functions = append(cp.Functions, fn)
}

func (cp *CxxParser) checkFunctionReturn(fn *Function) {
	if fn.ReturnType == jane.Void {
		return
	}
	miss := true
	for _, s := range fn.Block.Content {
		if s.Type == ast.StatementReturn {
			if !jane.TypesAreCompatible(
				s.Value.(ast.ReturnAST).Expression.Type, fn.ReturnType) {
				cp.PushErrorToken(s.Token, "incompatible_type")
			}
			miss = false
		}
	}
	if miss {
		cp.PushErrorToken(fn.Token, "missing_return")
	}
}

func (cp *CxxParser) functionByBName(name string) *Function {
	for _, function := range cp.Functions {
		if function.Name == name {
			return function
		}
	}
	return nil
}

func (cp *CxxParser) finalCheck() {
	if cp.functionByBName(jane.EntryPoint) == nil {
		cp.PushError("no_entry_point")
	}
}
