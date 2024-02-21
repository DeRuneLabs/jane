package ast

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/De-Rune/jane/lexer"
	"github.com/De-Rune/jane/package/jane"
	"github.com/De-Rune/jane/package/jnbits"
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
	ast.Errors = append(ast.Errors, fmt.Sprintf("%s:%d:%d %s", token.File.Path, token.Line, token.Column, message))
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
		case lexer.Brace:
			ast.BuildBrace()
		case lexer.Function:
			ast.BuildFunction()
		default:
			ast.PushError("invalid_syntax")
			ast.Position++
		}
	}
}

func (ast *AST) BuildBrace() {
	token := ast.Tokens[ast.Position]
	switch token.Value {
	case "[":
		ast.BuildTag()
	default:
		ast.PushErrorToken(token, "invalid_syntax")
	}
}

func (ast *AST) BuildTag() {
	var tag AttributeAST
	ast.Position++
	if ast.Ended() {
		ast.PushErrorToken(ast.Tokens[ast.Position-1], "invalid_syntax")
		return
	}
	ast.Position++
	if ast.Ended() {
		ast.PushErrorToken(ast.Tokens[ast.Position-1], "invalid_syntax")
		return
	}
	tag.Token = ast.Tokens[ast.Position]
	if tag.Token.Type != lexer.Brace || tag.Token.Value != "]" {
		ast.PushErrorToken(tag.Token, "invalid_syntax")
		ast.Position = -1
		return
	}
	tag.Token = ast.Tokens[ast.Position-1]
	tag.Value = tag.Token.Value
	ast.Tree = append(ast.Tree, Object{
		Token: tag.Token,
		Type:  Attribute,
		Value: tag,
	})
	ast.Position++
}

func (ast *AST) BuildFunction() {
	ast.Position++
	var funAst FunctionAST
	funAst.Token = ast.Tokens[ast.Position]
	funAst.Name = funAst.Token.Value
	funAst.ReturnType.Type = jane.Void
	ast.Position++
	tokens := ast.getRange("(", ")")
	if tokens == nil {
		return
	} else if len(tokens) > 0 {
		ast.BuildParameters(&funAst, tokens)
	}
	if ast.Ended() {
		ast.Position++
		ast.PushError("function_body_not_exist")
		ast.Position = -1
		return
	}
	token := ast.Tokens[ast.Position]
	if token.Type == lexer.Type {
		funAst.ReturnType = ast.BuildType(token)
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
	funAst.Block = ast.BuildBlock(blockTokens)
	ast.Tree = append(ast.Tree, Object{
		Token: funAst.Token,
		Type:  Statement,
		Value: StatementAST{
			Token: funAst.Token,
			Type:  StatementFunction,
			Value: funAst,
		},
	})
}

func (ast *AST) BuildParameters(function *FunctionAST, tokens []lexer.Token) {
	last := 0
	braceCount := 0
	for index, token := range tokens {
		if token.Type == lexer.Brace {
			switch token.Value {
			case "{", "[", "(":
				braceCount++
			default:
				braceCount--
			}
		}
		if braceCount > 0 || token.Type != lexer.Comma {
			continue
		}
		ast.pushParameter(function, tokens[last:index], token)
		last = index + 1
	}
	if last < len(tokens) {
		if last == 0 {
			ast.pushParameter(function, tokens[last:], tokens[last])
		} else {
			ast.pushParameter(function, tokens[last:], tokens[last-1])
		}
	}
}

func (ast *AST) pushParameter(function *FunctionAST, tokens []lexer.Token, err lexer.Token) {
	if len(tokens) == 0 {
		ast.PushErrorToken(err, "invalid_syntax")
		return
	}
	nameToken := tokens[0]
	if nameToken.Type != lexer.Name {
		ast.PushErrorToken(nameToken, "invalid_syntax")
	}
	if len(tokens) < 2 {
		ast.PushErrorToken(nameToken, "type_missing")
	}
	for _, param := range function.Params {
		if param.Name == nameToken.Value {
			ast.PushErrorToken(nameToken, "parameter_exist")
			break
		}
	}
	function.Params = append(function.Params, ParameterAST{
		Token: nameToken,
		Name:  nameToken.Value,
		Type:  ast.BuildType(tokens[1]),
	})
}

func IsStatement(token lexer.Token) bool {
	return token.Type == lexer.SemiColon
}

func (ast *AST) BuildType(token lexer.Token) (t TypeAST) {
	if token.Type != lexer.Type {
		ast.PushErrorToken(token, "invalid_type")
		return
	}
	t.Token = token
	t.Type = jane.TypeFromName(token.Value)
	t.Value = token.Value
	return t
}

func IsSingleOperator(operator string) bool {
	return operator == "-"
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
		if braceCount > 0 || !IsStatement(token) {
			continue
		}
		if index-oldStatementPoint == 0 {
			continue
		}
		b.Content = append(b.Content, ast.BuildStatement(tokens[oldStatementPoint:index]))
		if ast.Position == -1 {
			break
		}
		oldStatementPoint = index + 1
	}
	if oldStatementPoint < len(tokens) {
		ast.PushErrorToken(tokens[len(tokens)-1], "missing_semicolon")
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

func (ast *AST) BuildNameStatement(tokens []lexer.Token) (s StatementAST) {
	if len(tokens) == 1 {
		ast.PushErrorToken(tokens[0], "invalid_syntax")
		return
	}
	switch tokens[1].Type {
	case lexer.Brace:
		switch tokens[1].Value {
		case "(":
			return ast.BuildFunctionCallStatement(tokens)
		}
	}
	ast.PushErrorToken(tokens[0], "invalid_syntax")
	return
}

func (ast *AST) BuildFunctionCallStatement(tokens []lexer.Token) StatementAST {
	var fnCall FunctionCallAST
	fnCall.Expression = ast.BuildExpression(tokens)
	fnCall.Token = tokens[0]
	fnCall.Name = fnCall.Token.Value
	tokens = tokens[1:]
	args := ast.getRangeTokens("(", ")", tokens)
	if args == nil {
		ast.Position = -1
		return StatementAST{}
	} else if len(args) != len(tokens)-2 {
		ast.PushErrorToken(tokens[len(tokens)-2], "invalid_syntax")
		ast.Position = -1
		return StatementAST{}
	}
	fnCall.Args = ast.BuildArgs(args)
	return StatementAST{
		Token: fnCall.Token,
		Value: fnCall,
		Type:  StatementFunctionCall,
	}
}

func (ast *AST) BuildArgs(tokens []lexer.Token) []ArgAST {
	var args []ArgAST
	last := 0
	braceCount := 0
	for index, token := range tokens {
		if token.Type == lexer.Brace {
			switch token.Value {
			case "{", "[", "(":
				braceCount++
			default:
				braceCount--
			}
		}
		if braceCount > 0 || token.Type != lexer.Comma {
			continue
		}
		ast.pushArg(&args, tokens[last:index], token)
		last = index + 1
	}
	if last < len(tokens) {
		if last == 0 {
			ast.pushArg(&args, tokens[last:], tokens[last])
		} else {
			ast.pushArg(&args, tokens[last:], tokens[last-1])
		}
	}
	return args
}

func (ast *AST) pushArg(args *[]ArgAST, tokens []lexer.Token, err lexer.Token) {
	if len(tokens) == 0 {
		ast.PushErrorToken(err, "invalid_syntax")
		return
	}
	var arg ArgAST
	arg.Token = tokens[0]
	arg.Tokens = tokens
	arg.Expression = ast.BuildExpression(arg.Tokens)
	*args = append(*args, arg)
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
	e.Processes = ast.getExpressionProcesses(tokens)
	e.Tokens = tokens
	return
}

func (ast *AST) getExpressionProcesses(tokens []lexer.Token) [][]lexer.Token {
	var processes [][]lexer.Token
	var part []lexer.Token
	operator := false
	value := false
	braceCount := 0
	pushedError := false
	singleOperatored := false
	for index, token := range tokens {
		switch token.Type {
		case lexer.Operator:
			if !operator {
				if IsSingleOperator(token.Value) && !singleOperatored {
					part = append(part, token)
					singleOperatored = true
					continue
				}
				ast.PushErrorToken(token, "operator_overflow")
			}
			singleOperatored = false
			operator = false
			value = true
			if braceCount > 0 {
				part = append(part, token)
				continue
			}
			processes = append(processes, part)
			processes = append(processes, []lexer.Token{token})
			part = []lexer.Token{}
			continue
		case lexer.Brace:
			switch token.Value {
			case "(", "[", "{":
				singleOperatored = false
				braceCount++
			default:
				braceCount--
			}
		}
		if index > 0 {
			lt := tokens[index+1]
			if (lt.Type == lexer.Name || lt.Type == lexer.Value) && (token.Type == lexer.Name || token.Type == lexer.Value) {
				ast.PushErrorToken(token, "invalid_syntax")
				pushedError = true
			}
		}
		ast.checkExpressionToken(token)
		part = append(part, token)
		operator = requireOperatorForProcess(token, index, len(tokens))
		value = false
	}
	if len(part) > 0 {
		processes = append(processes, part)
	}
	if value {
		ast.PushErrorToken(processes[len(processes)-1][0], "operator_overflow")
		pushedError = true
	}
	if pushedError {
		return nil
	}
	return processes
}

func requireOperatorForProcess(token lexer.Token, index, tokensLen int) bool {
	switch token.Type {
	case lexer.Brace:
		if token.Value == "[" || token.Value == "(" || token.Value == "{" {
			return false
		}
	}
	return index < tokensLen-1
}

func (ast *AST) checkExpressionToken(token lexer.Token) {
	if token.Value[0] >= '0' && token.Value[0] <= '9' {
		var result bool
		if strings.IndexByte(token.Value, '.') != -1 {
			_, result = new(big.Float).SetString(token.Value)
		} else {
			result = jnbits.CheckBitInt(token.Value, 64)
		}
		if !result {
			ast.PushErrorToken(token, "invalid_numeric_range")
		}
	}
}

func (ast *AST) getRangeTokens(open, close string, tokens []lexer.Token) []lexer.Token {
	braceCount := 0
	start := 1
	for index, token := range tokens {
		if token.Type != lexer.Brace {
			continue
		}
		if token.Value == open {
			braceCount++
		} else if token.Value == close {
			braceCount--
		}
		if braceCount > 0 {
			continue
		}
		return tokens[start:index]
	}
	ast.PushErrorToken(tokens[0], "brace_not_closed")
	return nil
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
