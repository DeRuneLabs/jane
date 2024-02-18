package ast

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"

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

func (ast *AST) PushErrorToken(token lexer.Token, err string) {
	message := jane.Errors[err]
	ast.Errors = append(ast.Errors, fmt.Sprintf("%s:%d %s", token.File.Path, token.Line, message))
}

func (ast *AST) PushError(err string) {
	ast.PushErrorToken(ast.Tokens[ast.Position], err)
}

func (ast *AST) Ended() bool {
	return ast.Position >= len(ast.Tokens)
}

func (ast *AST) Build() {
	for ast.Position != -1 && !ast.Ended() {
		firstToken := ast.Tokens[ast.Position]
		switch firstToken.Type {
		case lexer.Name:
			ast.processName()
		default:
			ast.PushError("invalid_syntax")
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
		ast.PushError("parameters_not_supported")
	}
	if ast.Ended() {
		ast.Position--
		ast.PushError("function_body_not_exist")
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
			ast.PushError("function_body_not_exist")
			ast.Position = -1
			return
		}
		token = ast.Tokens[ast.Position]
	}
	if token.Type != lexer.Brace || token.Value != "{" {
		ast.PushError("invalid_syntax")
		ast.Position = -1
		return
	}
	blockTokens := ast.getRange("{", "}")
	if blockTokens == nil {
		ast.PushError("function_body_not_exist")
		ast.Position = -1
		return
	}
	function.Block = ast.BuildBlock(blockTokens)
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

func IsStatement(before, current lexer.Token) bool {
	return current.Type == lexer.SemiColon || before.Line < current.Line
}

func (ast *AST) BuildBlock(tokens []lexer.Token) (b BlockAST) {
	braceCount := 0
	oldStatementPoint := 0
	for index, token := range tokens {
		if token.Type == lexer.Brace {
			if token.Value == "{" {
				braceCount++
			} else {
				braceCount--
			}
		}
		if braceCount < 0 {
			continue
		}
		for index < len(tokens)-1 {
			if index == 0 && !IsStatement(token, token) {
				continue
			} else if index > 0 && !IsStatement(tokens[index-1], token) {
				continue
			}
		}
		if token.Type != lexer.SemiColon {
			index++
		}
		if index-oldStatementPoint == 0 {
			continue
		}
		b.Content = append(b.Content, ast.BuildStatement(tokens[oldStatementPoint:index]))
		oldStatementPoint = index + 1
	}
	return
}
func (ast *AST) BuildStatement(tokens []lexer.Token) (s StatementAST) {
	firstToken := tokens[0]
	switch firstToken.Type {
	case lexer.Return:
		return ast.BuildReturnStatement(tokens)
	default:
		ast.PushErrorToken(firstToken, "invalid_syntax")
	}
	return
}

func (ast *AST) BuildReturnStatement(tokens []lexer.Token) StatementAST {
	var returnModel ReturnAST
	returnModel.Token = tokens[0]
	if len(tokens) > 1 {
		returnModel.Expression = ast.BuildExpression(tokens[1:])
	}
	return StatementAST{
		Token: returnModel.Token,
		Type:  StatementReturn,
		Value: returnModel,
	}
}

func (ast *AST) BuildExpression(tokens []lexer.Token) (e ExpressionAST) {
	processes := ast.getExpressionProcesses(tokens)
	if len(processes) == 1 {
		value := ast.processExpression(tokens)
		e.Content = append(e.Content, ExpressionNode{
			Content: value,
			Type:    ExpressionNodeValue,
		})
		e.Type = value.Type
		return
	}
	return
}

func IsString(value string) bool {
	return value[0] == '"'
}

func IsBoolean(value string) bool {
	return value == "true" || value == "false"
}

func (ast *AST) processSingleValuePart(token lexer.Token) (result ValueAST) {
	result.Type = NA
	result.Token = token
	switch token.Type {
	case lexer.Value:
		if IsString(token.Value) {
			result.Data = token.Value[1 : len(token.Value)-1]
			result.Type = jane.String
		} else if IsBoolean(token.Value) {
			result.Data = token.Value
			result.Type = jane.Boolean
		}
		if strings.Contains(token.Value, ".") || strings.ContainsAny(token.Value, "eE") {
			result.Type = jane.Float64
		} else {
			result.Type = jane.Float32
			ok := CheckBitInt(token.Value, 32)
			if !ok {
				result.Type = jane.Int64
			}
		}
		result.Data = token.Value
	}
	return
}

func (ast *AST) processExpression(tokens []lexer.Token) (result ValueAST) {
	if len(tokens) == 1 {
		result = ast.processSingleValuePart(tokens[0])
		if result.Type != NA {
			goto end
		}
	}
	ast.PushErrorToken(tokens[0], "invalid_syntax")
end:
	return
}

func (ast *AST) getExpressionProcesses(tokens []lexer.Token) [][]lexer.Token {
	var processes [][]lexer.Token
	var part []lexer.Token
	braceCount := 0
	pushedError := false
	for index, token := range tokens {
		switch token.Type {
		case lexer.Brace:
			switch token.Value {
			case "(", "[", "{":
				braceCount++
			default:
				braceCount--
			}
		}
		if index > 0 {
			lt := tokens[index-1]
			if (lt.Type == lexer.Name || lt.Type == lexer.Value) && (token.Type == lexer.Name || token.Type == lexer.Value) {
				ast.PushErrorToken(token, "invalid_syntax")
				pushedError = true
			}
		}
		ast.checkExpressionToken(token)
		part = append(part, token)
	}
	if len(part) != 0 {
		processes = append(processes, part)
	}
	if pushedError {
		return nil
	}
	return processes
}

func CheckBitInt(value string, bit int) bool {
	_, err := strconv.ParseInt(value, 10, bit)
	return err == nil
}

func (ast *AST) checkExpressionToken(token lexer.Token) {
	if token.Value[0] >= '0' && token.Value[0] <= '9' {
		var result bool
		if strings.IndexByte(token.Value, '.') != -1 {
			_, result = new(big.Float).SetString(token.Value)
		} else {
			result = CheckBitInt(token.Value, 64)
		}
		if !result {
			ast.PushErrorToken(token, "invalid_numeric_range")
		}
	}
}

func (ast *AST) processName() {
	ast.Position++
	if ast.Ended() {
		ast.Position--
		ast.PushError("invalid_syntax")
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
			ast.PushError("invalid_syntax")
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
			ast.PushError("brace_not_closed")
			ast.Position = -1
			return nil
		}
		return ast.Tokens[start : ast.Position-1]
	}
	return nil
}
