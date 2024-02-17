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

func (ast *AST) pushError(err string) {
	message := jane.Errors[err]
	token := ast.Tokens[ast.Position]
	ast.Errors = append(ast.Errors, fmt.Sprintf("%s %s: %d", token.File.Path, message, token.Line))
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
	function.ReturnType.Type = uint8(jane.Void)
	ast.Position += 3
	if ast.Ended() {
		ast.Position--
		ast.pushError("function_body")
		ast.Position = -1
		return
	}
	token := ast.Tokens[ast.Position]
	if token.Type == lexer.Type {
		function.ReturnType.Type = jane.TypeFromName(token.Value)
		function.ReturnType.Value = token.Value
		ast.Position++
		if ast.Ended() {
			ast.Position--
			ast.pushError("function_body")
			ast.Position = -1
			return
		}
		token = ast.Tokens[ast.Position]
	}
	switch token.Type {
	case lexer.Brace:
		if token.Value != "{" {
			ast.pushError("invalid_syntax")
			ast.Position = -1
			return
		}
		ast.Position += 2
	default:
		ast.pushError("invalid_syntax")
		ast.Position = -1
		return
	}
	ast.Tree = append(ast.Tree, Object{
		Token: function.Token,
		Type:  Statement,
		Value: StatementAST{
			Token: function.Token,
			Type:  StatementFunction,
			Value: function,
		},
	})
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
