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

func IsString(value string) bool {
	return value[0] == '"'
}

func IsBoolean(value string) bool {
	return value == "true" || value == "false"
}

func CheckBitInt(value string, bit int) bool {
	_, err := strconv.ParseInt(value, 10, bit)
	return err == nil
}

func IsSingleOperator(operator string) bool {
	return operator == "-" ||
		operator == "!" ||
		operator == "*" ||
		operator == "&"
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
	return ast.processExpression(tokens)
}

func (ast *AST) processSingleValuePart(token lexer.Token) (result ValueAST) {
	result.Type = NA
	result.Token = token
	switch token.Type {
	case lexer.Value:
		if IsString(token.Value) {
			result.Value = token.Value[1 : len(token.Value)-1]
			result.Type = jane.String
		} else if IsBoolean(token.Value) {
			result.Value = token.Value
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
		result.Value = token.Value
	}
	return
}

type arithmeticProcess struct {
	ast      *AST
	left     []lexer.Token
	leftVal  ValueAST
	right    []lexer.Token
	rightVal ValueAST
	operator lexer.Token
}

func (p arithmeticProcess) solveString() (value ValueAST) {
	if p.leftVal != p.rightVal {
		p.ast.PushErrorToken(p.operator, "invalid_data_types")
		return
	}
	value.Type = jane.String
	switch p.operator.Value {
	case "+":
		value.Value = p.leftVal.String() + p.rightVal.String()
	default:
		p.ast.PushErrorToken(p.operator, "oeprator_notfor_string")
	}
	return
}

func (p arithmeticProcess) solve() (value ValueAST) {
	switch {
	case p.leftVal.Type == jane.Boolean || p.rightVal.Type == jane.Boolean:
		p.ast.PushErrorToken(p.operator, "operator_notfor_booleans")
		return
	case p.leftVal.Type == jane.String || p.rightVal.Type == jane.String:
		return p.solveString()
	}
	if jane.IsSignedNumericType(p.leftVal.Type) != jane.IsSignedNumericType(p.rightVal.Type) {
		p.ast.PushErrorToken(p.operator, "operator_notfor_uint_and_int")
		return
	}
	value.Type = p.leftVal.Type
	if jane.TypeGreaterThan(p.rightVal.Type, value.Type) {
		value.Type = p.rightVal.Type
	}
	return
}

func (ast *AST) processValuePart(tokens []lexer.Token) (result ValueAST) {
	if len(tokens) == 1 {
		result = ast.processSingleValuePart(tokens[0])
		if result.Type != NA {
			goto end
		}
		switch token := tokens[len(tokens)-1]; token.Type {
		default:
			ast.PushErrorToken(tokens[0], "invalid_syntax")
		}
	}
end:
	return
}

func (ast *AST) processExpression(tokens []lexer.Token) ExpressionAST {
	processes := ast.getExpressionProcesses(tokens)
	if len(processes) == 1 {
		value := ast.processValuePart(processes[0])
		return ExpressionAST{
			Content: []ExpressionNode{{
				Content: value,
				Type:    ExpressionNodeValue,
			}},
			Type: value.Type,
		}
	}
	result := buildExpressionByProcesses(processes)
	var process arithmeticProcess
	var value ValueAST
	process.ast = ast
	j := ast.nextOperator(processes)
	for j != -1 {
		if j == 0 {
			process.leftVal = value
			process.operator = processes[j][0]
			process.right = processes[j+1]
			process.rightVal = ast.processValuePart(process.right)
			value = process.solve()
			processes = processes[2:]
			j = ast.nextOperator(processes)
			continue
		} else if j == len(processes)-1 {
			process.operator = processes[j][0]
			process.left = processes[j-1]
			process.leftVal = ast.processValuePart(process.left)
			process.rightVal = value
			value = process.solve()
			processes = processes[:j-1]
			j = ast.nextOperator(processes)
			continue
		} else if prev := processes[j-1]; prev[0].Type == lexer.Operator && len(prev) == 1 {
			process.leftVal = value
			process.operator = processes[j][0]
			process.right = processes[j+1]
			process.rightVal = ast.processValuePart(process.right)
			value = process.solve()
			processes = append(processes[:j], processes[j+2:]...)
			j = ast.nextOperator(processes)
			continue
		}
		process.left = processes[j-1]
		process.leftVal = ast.processValuePart(process.left)
		process.operator = processes[j][0]
		process.right = processes[j+1]
		process.rightVal = ast.processValuePart(process.right)
		solvedValue := process.solve()
		if value.Type != NA {
			process.operator.Value = "+"
			process.right = processes[j+1]
			process.leftVal = value
			process.rightVal = solvedValue
			value = process.solve()
		} else {
			value = solvedValue
		}
		processes = append(processes[:j-1], processes[j+2:]...)
		if len(processes) == 1 {
			break
		}
		j = ast.nextOperator(processes)
	}
	result.Type = value.Type
	return result
}

func buildExpressionByProcesses(processes [][]lexer.Token) ExpressionAST {
	var result ExpressionAST
	for _, part := range processes {
		for _, token := range part {
			switch token.Type {
			case lexer.Operator:
				result.Content = append(result.Content, ExpressionNode{
					Content: OperatorAST{
						Token: token,
						Value: token.Value,
					},
					Type: ExpressionNodeOperator,
				})
			case lexer.Value:
				result.Content = append(result.Content, ExpressionNode{
					Content: ValueAST{
						Token: token,
						Value: token.Value,
					},
					Type: ExpressionNodeValue,
				})
			case lexer.Brace:
				result.Content = append(result.Content, ExpressionNode{
					Content: BraceAST{
						Token: token,
						Value: token.Value,
					},
					Type: ExpressionNodeBrace,
				})
			}
		}
	}
	return result
}

func (ast *AST) nextOperator(tokens [][]lexer.Token) int {
	high, mid, low := -1, -1, -1
	for index, part := range tokens {
		if len(part) != 1 {
			continue
		} else if part[0].Type != lexer.Operator {
			continue
		}
		switch part[0].Value {
		case "<<", ">>":
			return index
		case "&", "&^", "%":
			if high == -1 {
				high = index
			}
		case "*", "/", "\\", "|":
			if mid == -1 {
				mid = index
			}
		case "+", "-":
			if low == -1 {
				low = index
			}
		default:
			ast.PushErrorToken(part[0], "invalid_operator")
		}
	}
	if high != -1 {
		return high
	} else if mid != -1 {
		return mid
	}
	return low
}

func (ast *AST) getExpressionProcesses(tokens []lexer.Token) [][]lexer.Token {
	var processes [][]lexer.Token
	var part []lexer.Token
	operator := false
	braceCount := 0
	pushedError := false
	for index, token := range tokens {
		switch token.Type {
		case lexer.Operator:
			if !operator {
				if IsSingleOperator(token.Value) {
					part = append(part, token)
					continue
				}
				ast.PushErrorToken(token, "operator_overflow")
			}
			operator = false
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
		operator = requireOperatorForProcess(token, index, len(tokens))
	}
	if len(part) != 0 {
		processes = append(processes, part)
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
