package ast

import (
	"fmt"

	"github.com/De-Rune/jane/lexer"
	"github.com/De-Rune/jane/package/jane"
)

type AST struct {
	Tree     []Object
	Errors   []string
	Tokens   []lexer.Token
	Position int
}

func New(tokens []lexer.Token) *AST {
	ast := new(AST)
	ast.Tokens = tokens
	ast.Position = 0
	return ast
}

func (ast *AST) pushErrorToken(token lexer.Token, err string) {
	message := jane.Errors[err]
	ast.Errors = append(ast.Errors, fmt.Sprintf("%s:%d %s", token.File.Path, token.Line, message))
}

func (ast *AST) pushError(err string) {
	ast.pushErrorToken(ast.Tokens[ast.Position], err)
}

func (ast *AST) Ended() bool {
	return ast.Position >= len(ast.Tokens)
}

func (ast *AST) Build() {
	for ast.Position != -1 && !ast.Ended() {
		firstToken := ast.Tokens[ast.Position]
		switch firstToken.Type {
		case lexer.Name:
			ast.ProcessName()
		default:
			ast.pushError("invalid_syntax")
		}
	}
}

func (ast *AST) BuildFunction() {
	var function FunctionAST
	function.Token = ast.Tokens[ast.Position]
	function.Name = function.Token.Value
	function.ReturnType.Type = jane.Void
	ast.Position++
	parameters := ast.getRange("(", ")")
	if parameters == nil {
		return
	} else if len(parameters) > 0 {
		ast.pushError("parameters_no_supported")
	}
	if ast.Ended() {
		ast.Position--
		ast.pushError("function_body_not_exist")
	}
}

func (ast *AST) ProcessName() {
	ast.Position++
	if ast.Ended() {
		ast.Position--
		ast.pushError("invalid_syntax")
		return
	}
	ast.Position--
	secondToken := ast.Tokens[ast.Position+1]
	switch secondToken.Type {
	case lexer.Brace:
		switch secondToken.Value {
		case "(":
			ast.BuildFunction()
		default:
			ast.pushError("invalid_syntax")
		}
	}
}

func (ast *AST) getRange(open, close string) []lexer.Token {
	token := ast.Tokens[ast.Position]
	if token.Type == lexer.Brace && token.Value == open {
		ast.Position++
		braceCount := 1
		start := ast.Position
		for ; braceCount > 0 && !ast.Ended(); ast.Position++ {
			token := ast.Tokens[ast.Position]
			if token.Type != lexer.Brace {
				continue
			}
			if token.Value == open {
				braceCount++
			} else if token.Value == close {
				braceCount--
			}
		}
		if braceCount > 0 {
			ast.Position--
			ast.pushError("brace_not_closed")
			ast.Position = -1
			return nil
		}
		return ast.Tokens[start : ast.Position-1]
	}
	return nil
}
