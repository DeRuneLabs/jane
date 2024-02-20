package parser

import (
	"fmt"
	"strings"

	"github.com/De-Rune/jane/ast"
	"github.com/De-Rune/jane/lexer"
	"github.com/De-Rune/jane/package/jane"
	"github.com/De-Rune/jane/package/jnbits"
)

type CxxParser struct {
	Functions       []*Function
	GlobalVariables []*Variable
	BlockVariables  []*Variable

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
	cp.PFI.Errors = append(cp.PFI.Errors, fmt.Sprintf("%s:%d:%d %s", token.File.Path, token.Line, token.Column, message))
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
	if token := cp.existName(fnAst.Name); token.Type != ast.NA {
		cp.PushErrorToken(fnAst.Token, "exist_name")
		return
	}
	fn := new(Function)
	fn.Token = fnAst.Token
	fn.Name = fnAst.Name
	fn.ReturnType = fnAst.ReturnType.Type
	fn.Block = fnAst.Block
	fn.Params = fnAst.Params
	cp.Functions = append(cp.Functions, fn)
}

func variablesFromParameters(params []ast.ParameterAST) []*Variable {
	var vars []*Variable
	for _, param := range params {
		variable := new(Variable)
		variable.Name = param.Name
		variable.Token = param.Token
		variable.Type = param.Type.Type
	}
	return vars
}

func (cp *CxxParser) checkFunctionReturn(fn *Function) {
	if fn.ReturnType == jane.Void {
		return
	}
	miss := true
	for _, s := range fn.Block.Content {
		if s.Type == ast.StatementReturn {
			value := cp.computeExpression(s.Value.(ast.ReturnAST).Expression)
			if !jane.TypesAreCompatible(value.Type, fn.ReturnType) {
				cp.PushErrorToken(s.Token, "incompatible_type")
			}
			miss = false
		}
	}
	if miss {
		cp.PushErrorToken(fn.Token, "missing_return")
	}
}

func (cp *CxxParser) functionByName(name string) *Function {
	for _, function := range cp.Functions {
		if function.Name == name {
			return function
		}
	}
	return nil
}

func (cp *CxxParser) variableByName(name string) *Variable {
	for _, variable := range cp.BlockVariables {
		if variable.Name == name {
			return variable
		}
	}
	for _, variable := range cp.GlobalVariables {
		if variable.Name == name {
			return variable
		}
	}
	return nil
}

func (cp *CxxParser) existName(name string) lexer.Token {
	fn := cp.functionByName(name)
	if fn != nil {
		return fn.Token
	}
	return lexer.Token{}
}

func (cp *CxxParser) finalCheck() {
	if cp.functionByName(jane.EntryPoint) == nil {
		cp.PushError("no_entry_point")
	}
	for _, fn := range cp.Functions {
		cp.BlockVariables = variablesFromParameters(fn.Params)
		cp.checkFunctionReturn(fn)
	}
}

func (cp *CxxParser) computeProcesses(processes [][]lexer.Token) ast.ValueAST {
	if processes == nil {
		return ast.ValueAST{}
	}
	if len(processes) == 1 {
		value := cp.processValuePart(processes[0])
		return value
	}
	var process arithmeticProcess
	var value ast.ValueAST
	process.cp = cp
	j := cp.nextOperator(processes)
	for j != -1 {
		if j == 0 {
			process.leftVal = value
			process.operator = processes[j][0]
			process.right = processes[j+1]
			process.rightVal = cp.processValuePart(process.right)
			value = process.solve()
			processes = processes[2:]
			j = cp.nextOperator(processes)
			continue
		} else if j == len(processes)-1 {
			process.operator = processes[j][0]
			process.left = processes[j-1]
			process.leftVal = cp.processValuePart(process.left)
			process.rightVal = value
			value = process.solve()
			processes = processes[:j-1]
			j = cp.nextOperator(processes)
			continue
		} else if prev := processes[j-1]; prev[0].Type == lexer.Operator &&
			len(prev) == 1 {
			process.leftVal = value
			process.operator = processes[j][0]
			process.right = processes[j+1]
			process.rightVal = cp.processValuePart(process.right)
			value = process.solve()
			processes = append(processes[:j], processes[j+2:]...)
			j = cp.nextOperator(processes)
			continue
		}
		process.left = processes[j-1]
		process.leftVal = cp.processValuePart(process.left)
		process.operator = processes[j][0]
		process.right = processes[j+1]
		process.rightVal = cp.processValuePart(process.right)
		solvedValue := process.solve()
		if value.Type != ast.NA {
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
		j = cp.nextOperator(processes)
	}
	return value
}

func (cp *CxxParser) computeTokens(tokens []lexer.Token) ast.ValueAST {
	return cp.computeProcesses(new(ast.AST).BuildExpression(tokens).Processes)
}

func (cp *CxxParser) computeExpression(ex ast.ExpressionAST) ast.ValueAST {
	return cp.computeProcesses(ex.Processes)
}

func (cp *CxxParser) nextOperator(tokens [][]lexer.Token) int {
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
			cp.PushErrorToken(part[0], "invalid_operator")
		}
	}
	if high != -1 {
		return high
	} else if mid != -1 {
		return mid
	}
	return low
}

type arithmeticProcess struct {
	cp       *CxxParser
	left     []lexer.Token
	leftVal  ast.ValueAST
	right    []lexer.Token
	rightVal ast.ValueAST
	operator lexer.Token
}

func (p arithmeticProcess) solveString() (value ast.ValueAST) {
	if p.leftVal.Type != p.rightVal.Type {
		p.cp.PushErrorToken(p.operator, "invalid_data_types")
		return
	}
	value.Type = jane.String
	switch p.operator.Value {
	case "+":
		value.Value = p.leftVal.String() + p.rightVal.String()
	default:
		p.cp.PushErrorToken(p.operator, "operator_notfor_strings")
	}
	return
}

func (p arithmeticProcess) solve() (value ast.ValueAST) {
	switch {
	case p.leftVal.Type == jane.Boolean || p.rightVal.Type == jane.Boolean:
		p.cp.PushErrorToken(p.operator, "operator_notfor_booleans")
		return
	case p.leftVal.Type == jane.String || p.rightVal.Type == jane.String:
		return p.solveString()
	}
	if jane.IsSignedNumericType(p.leftVal.Type) != jane.IsSignedNumericType(p.rightVal.Type) {
		p.cp.PushErrorToken(p.operator, "operator_notfor_uint_and_int")
		return
	}
	value.Type = p.leftVal.Type
	if jane.TypeGreaterThan(p.rightVal.Type, value.Type) {
		value.Type = p.rightVal.Type
	}
	return
}

const functionName = 0x0000A

func (cp *CxxParser) processSingleValuePart(token lexer.Token) (result ast.ValueAST) {
	result.Type = ast.NA
	result.Token = token
	switch token.Type {
	case lexer.Value:
		if IsString(token.Value) {
			result.Value = token.Value
			result.Type = jane.String
		} else if IsBoolean(token.Value) {
			result.Value = token.Value
			result.Type = jane.Boolean
		} else {
			if strings.Contains(token.Value, ".") || strings.ContainsAny(token.Value, "eE") {
				result.Type = jane.Float64
			} else {
				result.Type = jane.Int32
				ok := jnbits.CheckBitInt(token.Value, 32)
				if !ok {
					result.Type = jane.Int64
				}
			}
			result.Value = token.Value
		}
	case lexer.Name:
		if cp.functionByName(token.Value) != nil {
			result.Value = token.Value
			result.Type = functionName
		} else if variable := cp.variableByName(token.Value); variable != nil {
			result.Value = token.Value
			result.Type = variable.Type
		} else {
			cp.PushErrorToken(token, "name_not_defined")
		}
	default:
		cp.PushErrorToken(token, "invalid_syntax")
	}
	return
}

func (cp *CxxParser) processValuePart(tokens []lexer.Token) (result ast.ValueAST) {
	if len(tokens) == 1 {
		result = cp.processSingleValuePart(tokens[0])
		if result.Type != ast.NA {
			goto end
		}
	}
	switch token := tokens[len(tokens)-1]; token.Type {
	case lexer.Brace:
		switch token.Value {
		case ")":
			return cp.processParenthesesValuePart(tokens)
		}
	default:
		cp.PushErrorToken(tokens[0], "invalid_syntax")
	}
end:
	return
}

func (cp *CxxParser) processParenthesesValuePart(tokens []lexer.Token) ast.ValueAST {
	var valueTokens []lexer.Token
	j := len(tokens) - 1
	braceCount := 0
	for ; j >= 0; j-- {
		token := tokens[j]
		if token.Type != lexer.Brace {
			continue
		}
		switch token.Value {
		case ")":
			braceCount++
		case "(":
			braceCount--
		}
		if braceCount > 0 {
			continue
		}
		valueTokens = tokens[:j]
		break
	}
	if len(valueTokens) == 0 && braceCount == 0 {
		tk := tokens[0]
		tokens = tokens[1 : len(tokens)-1]
		if len(tokens) == 0 {
			cp.PushErrorToken(tk, "invalid_syntax")
		}
		return cp.computeTokens(tokens)
	}
	value := cp.processValuePart(valueTokens)
	switch value.Type {
	case functionName:
		fn := cp.functionByName(value.Value)
		cp.parseFunctionCallStatement(fn, tokens[len(valueTokens):])
		value.Type = fn.ReturnType
	default:
		cp.PushErrorToken(tokens[len(valueTokens)], "invalid_syntax")
	}
	return value
}

func (cp *CxxParser) parseFunctionCallStatement(fn *Function, tokens []lexer.Token) {
	errToken := tokens[0]
	tokens = cp.getRangeTokens("(", ")", tokens)
	if tokens == nil {
		tokens = make([]lexer.Token, 0)
	}
	if cp.parseArgs(fn, tokens) < len(fn.Params) {
		cp.PushErrorToken(errToken, "argument_missing")
	}
}

func (cp *CxxParser) parseArgs(fn *Function, tokens []lexer.Token) int {
	last := 0
	braceCount := 0
	count := 0
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
		count++
		cp.parseArg(fn, count, tokens[last:index], token)
		last = index + 1
	}
	if last < len(tokens) {
		count++
		if last == 0 {
			cp.parseArg(fn, count, tokens[last:], tokens[last])
		} else {
			cp.parseArg(fn, count, tokens[last:], tokens[last-1])
		}
	}
	return count
}

func (cp *CxxParser) parseArg(fn *Function, count int, tokens []lexer.Token, err lexer.Token) {
	if len(tokens) == 0 {
		cp.PushErrorToken(err, "invalid_syntax")
		return
	}
	if count > len(fn.Params) {
		cp.PushErrorToken(err, "argument_overflow")
		return
	}
	if !jane.TypesAreCompatible(cp.computeTokens(tokens).Type, fn.Params[count-1].Type.Type) {
		cp.PushErrorToken(err, "incompatible_type")
	}
}

func (cp *CxxParser) getRangeTokens(open, close string, tokens []lexer.Token) []lexer.Token {
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
	cp.PushErrorToken(tokens[0], "brace_not_closed")
	return nil
}
