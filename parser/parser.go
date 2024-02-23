package parser

import (
	"fmt"
	"strings"

	"github.com/De-Rune/jane/ast"
	"github.com/De-Rune/jane/lexer"
	"github.com/De-Rune/jane/package/jane"
	"github.com/De-Rune/jane/package/jnbits"
)

type Parser struct {
	attributes      []ast.AttributeAST
	Functions       []*function
	GlobalVariables []*variable
	BlockVariables  []*variable
	Tokens          []lexer.Token
	PFI             *ParseFileInfo
}

func NewParser(tokens []lexer.Token, PFI *ParseFileInfo) *Parser {
	parser := new(Parser)
	parser.Tokens = tokens
	parser.PFI = PFI
	return parser
}

func (cp *Parser) PushErrorToken(token lexer.Token, err string) {
	message := jane.Errors[err]
	cp.PFI.Errors = append(cp.PFI.Errors, fmt.Sprintf("%s:%d:%d %s", token.File.Path, token.Line, token.Column, message))
}

func (cp *Parser) AppendErrors(errors ...string) {
	cp.PFI.Errors = append(cp.PFI.Errors, errors...)
}

func (cp *Parser) PushError(err string) {
	cp.PFI.Errors = append(cp.PFI.Errors, jane.Errors[err])
}

func (cp Parser) String() string {
	return cp.Cxx()
}

func (p *Parser) Cxx() string {
	var sb strings.Builder
	sb.WriteString("#pragma region JANE_GLOBAL_VARIABLES")
	sb.WriteByte('\n')
	for _, va := range p.GlobalVariables {
		sb.WriteString(va.String())
		sb.WriteByte('\n')
	}
	sb.WriteString("#pragma endregion JANE_GLOBAL_VARIABLES")
	sb.WriteString("\n\n")
	sb.WriteString("#pragma region JANE_FUNCTIONS")
	sb.WriteByte('\n')
	for _, fun := range p.Functions {
		sb.WriteString(fun.String())
		sb.WriteString("\n\n")
	}
	sb.WriteString("#pragma endregion JANE_FUNCTIONS")
	return sb.String()
}

func (cp *Parser) Parse() {
	astModel := ast.New(cp.Tokens)
	astModel.Build()
	if astModel.Errors != nil {
		cp.PFI.Errors = append(cp.PFI.Errors, astModel.Errors...)
		return
	}
	for _, model := range astModel.Tree {
		switch model.Type {
		case ast.Attribute:
			cp.PushAtrribute(model.Value.(ast.AttributeAST))
		case ast.Statement:
			cp.ParseStatement(model.Value.(ast.StatementAST))
		default:
			cp.PushErrorToken(model.Token, "invalid_syntax")
		}
	}
	cp.finalCheck()
}

func (cp *Parser) PushAtrribute(t ast.AttributeAST) {
	switch t.Token.Type {
	case lexer.Inline:
	default:
		cp.PushErrorToken(t.Token, "invalid_syntax")
	}
	cp.attributes = append(cp.attributes, t)
}

func (cp *Parser) ParseStatement(s ast.StatementAST) {
	switch s.Type {
	case ast.StatementFunction:
		cp.ParseFunction(s.Value.(ast.FunctionAST))
	default:
		cp.PushErrorToken(s.Token, "invalid_syntax")
	}
}

func (cp *Parser) ParseFunction(funAst ast.FunctionAST) {
	if token := cp.existName(funAst.Name); token.Type != ast.NA {
		cp.PushErrorToken(funAst.Token, "exist_name")
		return
	}
	fun := new(function)
	fun.Token = funAst.Token
	fun.Name = funAst.Name
	fun.ReturnType = funAst.ReturnType
	fun.Block = funAst.Block
	fun.Params = funAst.Params
	fun.Attributes = cp.attributes
	cp.attributes = nil
	cp.checkFunctionAttributes(fun.Attributes)
	cp.Functions = append(cp.Functions, fun)
}

func (cp *Parser) checkFunctionAttributes(tags []ast.AttributeAST) {
	for _, tag := range tags {
		switch tag.Token.Type {
		case lexer.Inline:
		default:
			cp.PushErrorToken(tag.Token, "invalid_tag")
		}
	}
}

func variablesFromParameters(params []ast.ParameterAST) []*variable {
	var vars []*variable
	for _, param := range params {
		variable := new(variable)
		variable.Name = param.Name
		variable.Token = param.Token
		variable.Type = param.Type.Type
	}
	return vars
}

func (cp *Parser) checkFunctionReturn(fun *function) {
	miss := true
	for _, s := range fun.Block.Content {
		if s.Type == ast.StatementReturn {
			retAST := s.Value.(ast.ReturnAST)
			if len(retAST.Expression.Tokens) == 0 {
				if fun.ReturnType.Type != jane.Void {
					cp.PushErrorToken(retAST.Token, "require_return_value")
				}
			} else {
				if fun.ReturnType.Type == jane.Void {
					cp.PushErrorToken(retAST.Token, "void_function_return_value")
				} else {
					value := cp.computeExpression(retAST.Expression)
					if !jane.TypesAreCompatible(value.Type, fun.ReturnType.Type, true) {
						cp.PushErrorToken(retAST.Token, "incompatible_type")
					}
				}
			}
			miss = false
		}
	}
	if miss && fun.ReturnType.Type != jane.Void {
		cp.PushErrorToken(fun.Token, "missing_return")
	}
}

func (cp *Parser) functionByName(name string) *function {
	for _, fun := range builtinFunctions {
		if fun.Name == name {
			return fun
		}
	}
	for _, fun := range cp.Functions {
		if fun.Name == name {
			return fun
		}
	}
	return nil
}

func (cp *Parser) variableByName(name string) *variable {
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

func (cp *Parser) existName(name string) lexer.Token {
	fun := cp.functionByName(name)
	if fun != nil {
		return fun.Token
	}
	return lexer.Token{}
}

func (cp *Parser) finalCheck() {
	if cp.functionByName(jane.EntryPoint) == nil {
		cp.PushError("no_entry_point")
	}
	for _, fun := range cp.Functions {
		cp.BlockVariables = variablesFromParameters(fun.Params)
		cp.checkFunction(fun)
		cp.checkBlock(fun.Block)
		cp.checkFunctionReturn(fun)
	}
}

func (cp *Parser) computeProcesses(processes [][]lexer.Token) ast.ValueAST {
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
	boolean := false
	for j != -1 {
		if !boolean {
			boolean = value.Type == jane.Bool
		}
		if boolean {
			value.Type = jane.Bool
		}
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
			process.leftVal = value
			process.right = processes[j+1]
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

func (cp *Parser) computeTokens(tokens []lexer.Token) ast.ValueAST {
	return cp.computeProcesses(new(ast.AST).BuildExpression(tokens).Processes)
}

func (cp *Parser) computeExpression(ex ast.ExpressionAST) ast.ValueAST {
	processes := make([][]lexer.Token, len(ex.Processes))
	copy(processes, ex.Processes)
	return cp.computeProcesses(ex.Processes)
}

func (cp *Parser) nextOperator(tokens [][]lexer.Token) int {
	precedence5 := -1
	precedence4 := -1
	precedence3 := -1
	precedence2 := -1
	precedence1 := -1
	for index, part := range tokens {
		if len(part) != 1 {
			continue
		} else if part[0].Type != lexer.Operator {
			continue
		}
		switch part[0].Value {
		case "*", "/", "%", "<<", ">>", "&":
			precedence5 = index
		case "+", "-", "|", "^":
			precedence4 = index
		case "==", "!=", "<", "<=", ">", ">=":
			precedence3 = index
		case "&&":
			precedence2 = index
		case "||":
			precedence1 = index
		default:
			cp.PushErrorToken(part[0], "invalid_operator")
		}
	}
	if precedence5 != -1 {
		return precedence5
	} else if precedence4 != -1 {
		return precedence4
	} else if precedence3 != -1 {
		return precedence3
	} else if precedence2 != -1 {
		return precedence2
	}
	return precedence1
}

type arithmeticProcess struct {
	cp       *Parser
	left     []lexer.Token
	leftVal  ast.ValueAST
	right    []lexer.Token
	rightVal ast.ValueAST
	operator lexer.Token
}

func (p arithmeticProcess) solveString() (value ast.ValueAST) {
	if p.leftVal.Type != p.rightVal.Type {
		p.cp.PushErrorToken(p.operator, "invalid_datatype")
		return
	}
	value.Type = jane.Str
	switch p.operator.Value {
	case "+":
		value.Type = jane.Str
	case "==", "!=":
		value.Type = jane.Bool
	default:
		p.cp.PushErrorToken(p.operator, "operator_notfor_strings")
	}
	return
}

func (p arithmeticProcess) solveAny() (value ast.ValueAST) {
	switch p.operator.Value {
	case "!=", "==":
		value.Type = jane.Bool
	default:
		p.cp.PushErrorToken(p.operator, "operator_notfor_any")
	}
	return
}

func (p arithmeticProcess) solveBool() (value ast.ValueAST) {
	if !jane.TypesAreCompatible(p.leftVal.Type, p.rightVal.Type, true) {
		p.cp.PushErrorToken(p.operator, "incompatible_type")
		return
	}
	switch p.operator.Value {
	case "!=", "==":
		value.Type = jane.Bool
	default:
		p.cp.PushErrorToken(p.operator, "operator_notfor_bool")
	}
	return
}

func (p arithmeticProcess) solveFloat() (value ast.ValueAST) {
	if !jane.TypesAreCompatible(p.leftVal.Type, p.rightVal.Type, true) {
		if !IsConstantNumeric(p.leftVal.Value) && !IsConstantNumeric(p.rightVal.Value) {
			p.cp.PushErrorToken(p.operator, "incompatible_type")
			return
		}
	}
	switch p.operator.Value {
	case "!=", "==", "<", ">", ">=", "<=":
		value.Type = jane.Bool
	case "+", "-", "*", "/":
		value.Type = jane.Float32
		if p.leftVal.Type == jane.Float64 || p.rightVal.Type == jane.Float64 {
			value.Type = jane.Float64
		}
	default:
		p.cp.PushErrorToken(p.operator, "operator_notfor_float")
	}
	return
}

func (p arithmeticProcess) solveSigned() (value ast.ValueAST) {
	if !jane.TypesAreCompatible(p.leftVal.Type, p.rightVal.Type, true) {
		if !IsConstantNumeric(p.leftVal.Value) && !IsConstantNumeric(p.rightVal.Value) {
			p.cp.PushErrorToken(p.operator, "incompatible_type")
			return
		}
	}
	switch p.operator.Value {
	case "!=", "==", "<", ">", ">=", "<=":
		value.Type = jane.Bool
	case "+", "-", "*", "/", "%", "&", "|", "^":
		value.Type = p.leftVal.Type
		if jane.TypeGreaterThan(p.rightVal.Type, value.Type) {
			value.Type = p.rightVal.Type
		}
	case ">>", "<<":
		value.Type = p.leftVal.Type
		if !jane.IsUnsignedNumericType(p.rightVal.Type) && !checkIntBit(p.rightVal, jnbits.BitsizeOfType(jane.UInt64)) {
			p.cp.PushErrorToken(p.rightVal.Token, "bitshift_must_unsigned")
		}
	default:
		p.cp.PushErrorToken(p.operator, "operator_notfor_int")
	}
	return
}

func (p arithmeticProcess) solveUnsigned() (value ast.ValueAST) {
	if !jane.TypesAreCompatible(p.leftVal.Type, p.rightVal.Type, true) {
		if !IsConstantNumeric(p.leftVal.Value) && !IsConstantNumeric(p.rightVal.Value) {
			p.cp.PushErrorToken(p.operator, "incompatible_type")
			return
		}
		return
	}
	switch p.operator.Value {
	case "!=", "==", "<", ">", ">=", "<=":
		value.Type = jane.Bool
	case "+", "-", "*", "/", "%", "&", "|", "^":
		value.Type = p.leftVal.Type
		if jane.TypeGreaterThan(p.rightVal.Type, value.Type) {
			value.Type = p.rightVal.Type
		}
	default:
		p.cp.PushErrorToken(p.operator, "operator_notfor_uint")
	}
	return
}

func (p arithmeticProcess) solveLogical() (value ast.ValueAST) {
	value.Type = jane.Bool
	if p.leftVal.Type != jane.Bool {
		p.cp.PushErrorToken(p.leftVal.Token, "logical_not_bool")
	}
	if p.rightVal.Type != jane.Bool {
		p.cp.PushErrorToken(p.rightVal.Token, "logical_not_bool")
	}
	return
}

func (p arithmeticProcess) solve() (value ast.ValueAST) {
	switch p.operator.Value {
	case "+", "-", "*", "/", "%", ">>",
		"<<", "&", "|", "^", "==", "!=",
		">=", "<=", ">", "<":
	case "&&", "||":
		return p.solveLogical()
	default:
		p.cp.PushErrorToken(p.operator, "invalid_operator")
	}
	switch {
	case p.leftVal.Type == jane.Any || p.rightVal.Type == jane.Any:
		return p.solveAny()
	case p.leftVal.Type == jane.Bool || p.rightVal.Type == jane.Bool:
		return p.solveBool()
	case p.leftVal.Type == jane.Str || p.rightVal.Type == jane.Str:
		return p.solveString()
	case jane.IsFloatType(p.leftVal.Type) || jane.IsFloatType(p.rightVal.Type):
		return p.solveFloat()
	case jane.IsSignedNumericType(p.leftVal.Type) || jane.IsSignedNumericType(p.rightVal.Type):
		return p.solveSigned()
	case jane.IsUnsignedNumericType(p.leftVal.Type) || jane.IsUnsignedNumericType(p.rightVal.Type):
		return p.solveUnsigned()
	}
	return
}

const functionName = 0x0000A

func (cp *Parser) processSingleValuePart(token lexer.Token) (result ast.ValueAST) {
	result.Type = ast.NA
	result.Token = token
	switch token.Type {
	case lexer.Value:
		if IsString(token.Value) {
			result.Value = "L" + token.Value
			result.Type = jane.Str
		} else if IsBoolean(token.Value) {
			result.Value = token.Value
			result.Type = jane.Bool
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

func (cp *Parser) processSingleOperatorPart(tokens []lexer.Token) ast.ValueAST {
	var result ast.ValueAST
	token := tokens[0]
	if len(tokens) == 0 {
		cp.PushErrorToken(token, "invalid_syntax")
		return result
	}
	switch token.Value {
	case "-":
		result = cp.processValuePart(tokens)
		if !jane.IsNumericType(result.Type) {
			cp.PushErrorToken(token, "invalid_data_unary")
		}
	case "+":
		result = cp.processValuePart(tokens)
		if !jane.IsNumericType(result.Type) {
			cp.PushErrorToken(token, "invalid_data_plus")
		}
	case "~":
		result = cp.processValuePart(tokens)
		if !jane.IsIntegerType(result.Type) {
			cp.PushErrorToken(token, "invalid_data_tilde")
		}
	case "!":
		result = cp.processValuePart(tokens)
		if result.Type != jane.Bool {
			cp.PushErrorToken(token, "invalid_data_logical_not")
		}
	default:
		cp.PushErrorToken(token, "invalid_syntax")
	}
	return result
}

func (cp *Parser) processValuePart(tokens []lexer.Token) (result ast.ValueAST) {
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

func (cp *Parser) processParenthesesValuePart(tokens []lexer.Token) ast.ValueAST {
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
		fun := cp.functionByName(value.Value)
		cp.parseFunctionCallStatement(fun, tokens[len(valueTokens):])
		value.Type = fun.ReturnType.Type
	default:
		cp.PushErrorToken(tokens[len(valueTokens)], "invalid_syntax")
	}
	return value
}

func (cp *Parser) parseFunctionCallStatement(fun *function, tokens []lexer.Token) {
	errToken := tokens[0]
	tokens = cp.getRangeTokens("(", ")", tokens)
	if tokens == nil {
		tokens = make([]lexer.Token, 0)
	}
	ast := new(ast.AST)
	args := ast.BuildArgs(tokens)
	if len(ast.Errors) > 0 {
		cp.AppendErrors(ast.Errors...)
	}
	cp.parseArgs(fun, args, errToken)
}

func (cp *Parser) parseArgs(fun *function, args []ast.ArgAST, errToken lexer.Token) {
	if len(args) < len(fun.Params) {
		cp.PushErrorToken(errToken, "argument_missing")
	}
	for index, arg := range args {
		cp.parseArg(fun, index, arg)
	}
}

func (cp *Parser) parseArg(fun *function, index int, arg ast.ArgAST) {
	if index >= len(fun.Params) {
		cp.PushErrorToken(arg.Token, "argument_overflow")
		return
	}
	value := cp.computeExpression(arg.Expression)
	param := fun.Params[index]
	if !jane.TypesAreCompatible(value.Type, param.Type.Type, false) {
		value.Type = param.Type.Type
		if !checkIntBit(value, jnbits.BitsizeOfType(param.Type.Type)) {
			cp.PushErrorToken(arg.Token, "incompatible_type")
		}
	}
}

func (cp *Parser) getRangeTokens(open, close string, tokens []lexer.Token) []lexer.Token {
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

func (cp *Parser) checkFunction(fun *function) {
	switch fun.Name {
	case jane.EntryPoint:
		if len(fun.Params) > 0 {
			cp.PushErrorToken(fun.Token, "entrypoint_have_parameters")
		}
		if fun.ReturnType.Type != jane.Void {
			cp.PushErrorToken(fun.ReturnType.Token, "entrypoint_have_return")
		}
	}
}

func (cp *Parser) checkBlock(b ast.BlockAST) {
	for _, model := range b.Content {
		switch model.Type {
		case ast.StatementFunctionCall:
			cp.checkFunctionCallStatement(model.Value.(ast.FunctionCallAST))
		case ast.StatementReturn:
		default:
			cp.PushErrorToken(model.Token, "invalid_syntax")
		}
	}
}

func (cp *Parser) checkFunctionCallStatement(cs ast.FunctionCallAST) {
	fun := cp.functionByName(cs.Name)
	if fun == nil {
		cp.PushErrorToken(cs.Token, "name_not_defined")
		return
	}
	cp.parseArgs(fun, cs.Args, cs.Token)
}

func IsConstantNumeric(v string) bool {
	if v == "" {
		return false
	}
	return v[0] >= '0' && v[0] <= '9'
}

func checkIntBit(v ast.ValueAST, bit int) bool {
	if bit == 0 {
		return false
	}
	if jane.IsSignedNumericType(v.Type) {
		return jnbits.CheckBitInt(v.Value, bit)
	}
	return jnbits.CheckBitUint(v.Value, bit)
}
