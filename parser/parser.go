package parser

import (
	"fmt"
	"strings"

	"github.com/De-Rune/jane/ast"
	"github.com/De-Rune/jane/lexer"
	"github.com/De-Rune/jane/package/jn"
	"github.com/De-Rune/jane/package/jnbits"
)

type Parser struct {
	attributes []ast.AttributeAST

	Functions              []*function
	GlobalVariables        []ast.VariableAST
	Types                  []ast.TypeAST
	WaitingGlobalVariables []ast.VariableAST
	BlockVariables         []ast.VariableAST
	Tokens                 []lexer.Token
	PFI                    *ParseFileInfo
}

func NewParser(tokens []lexer.Token, PFI *ParseFileInfo) *Parser {
	parser := new(Parser)
	parser.Tokens = tokens
	parser.PFI = PFI
	return parser
}

func (p *Parser) PushErrorToken(token lexer.Token, err string) {
	message := jn.Errors[err]
	p.PFI.Errors = append(
		p.PFI.Errors,
		fmt.Sprintf("%s:%d:%d %s", token.File.Path, token.Row, token.Column, message),
	)
}

func (p *Parser) AppendErrors(errors ...string) {
	p.PFI.Errors = append(p.PFI.Errors, errors...)
}

func (p *Parser) PushError(err string) {
	p.PFI.Errors = append(p.PFI.Errors, jn.Errors[err])
}

func (p Parser) String() string {
	return p.Cxx()
}

func (p *Parser) CxxTypes() string {
	var cxx strings.Builder
	cxx.WriteString("#pragma region TYPES\n")
	for _, t := range p.Types {
		cxx.WriteString(t.String())
		cxx.WriteByte('\n')
	}
	cxx.WriteString("#pragma endregion TYPES\n")
	return cxx.String()
}

func (p *Parser) CxxPrototypes() string {
	var cxx strings.Builder
	cxx.WriteString("#pragma endregion PROTOTYPES\n")
	for _, fun := range p.Functions {
		cxx.WriteString(fun.Prototype())
		cxx.WriteByte('\n')
	}
	cxx.WriteString("#pragma endregion PROTOTYPES")
	return cxx.String()
}

func (p *Parser) CxxGlobalVariables() string {
	var cxx strings.Builder
	cxx.WriteString("#pragma region GLOBAL_VARIABLES\n")
	for _, va := range p.GlobalVariables {
		cxx.WriteString(va.String())
		cxx.WriteByte('\n')
	}
	cxx.WriteString("#pragma endregion GLOBAL_VARIABLES")
	return cxx.String()
}

func (p *Parser) CxxFunctions() string {
	var cxx strings.Builder
	cxx.WriteString("#pragma region FUNCTIONS")
	cxx.WriteString("\n\n")
	for _, fun := range p.Functions {
		cxx.WriteString(fun.String())
		cxx.WriteString("\n\n")
	}
	cxx.WriteString("#pragma endregion FUNCTIONS")
	return cxx.String()
}

func (p *Parser) Cxx() string {
	var cxx strings.Builder
	cxx.WriteString(p.CxxTypes() + "\n\n")
	cxx.WriteString(p.CxxPrototypes() + "\n\n")
	cxx.WriteString(p.CxxGlobalVariables() + "\n\n")
	cxx.WriteString(p.CxxFunctions())
	return cxx.String()
}

func (p *Parser) Parse() {
	astModel := ast.New(p.Tokens)
	astModel.Build()
	if astModel.Errors != nil {
		p.PFI.Errors = append(p.PFI.Errors, astModel.Errors...)
		return
	}
	for _, model := range astModel.Tree {
		switch t := model.Value.(type) {
		case ast.AttributeAST:
			p.PushAttribute(t)
		case ast.StatementAST:
			p.ParseStatement(t)
		case ast.TypeAST:
			p.ParseType(t)
		default:
			p.PushErrorToken(model.Token, "invalid_syntax")
		}
	}
	p.finalCheck()
}

func (p *Parser) ParseType(t ast.TypeAST) {
	if p.existName(t.Name).Id != lexer.NA {
		p.PushErrorToken(t.Token, "exist_name")
		return
	} else if jn.IsIgnoreName(t.Name) {
		p.PushErrorToken(t.Token, "ignore_name_identifier")
		return
	}
	p.Types = append(p.Types, t)
}

func (p *Parser) PushAttribute(attribute ast.AttributeAST) {
	switch attribute.Tag.Kind {
	case "_inline":
	default:
		p.PushErrorToken(attribute.Tag, "undefined_tag")
	}
	for _, attr := range p.attributes {
		if attr.Tag.Kind == attribute.Tag.Kind {
			p.PushErrorToken(attribute.Tag, "attribute_repeat")
			return
		}
	}
	p.attributes = append(p.attributes, attribute)
}

func (p *Parser) ParseStatement(s ast.StatementAST) {
	switch t := s.Value.(type) {
	case ast.FunctionAST:
		p.ParseFunction(t)
	case ast.VariableAST:
		p.ParseGlobalVariable(t)
	default:
		p.PushErrorToken(s.Token, "invalid_syntax")
	}
}

func (p *Parser) ParseFunction(funAST ast.FunctionAST) {
	if p.existName(funAST.Name).Id != lexer.NA {
		p.PushErrorToken(funAST.Token, "exist_name")
	} else if jn.IsIgnoreName(funAST.Name) {
		p.PushErrorToken(funAST.Token, "ignore_name_identifier")
	}
	fun := new(function)
	fun.Ast = funAST
	fun.Attributes = p.attributes
	p.attributes = nil
	p.checkFunctionAttributes(fun.Attributes)
	p.Functions = append(p.Functions, fun)
}

func (p *Parser) ParseGlobalVariable(varAST ast.VariableAST) {
	if p.existName(varAST.Name).Id != lexer.NA {
		p.PushErrorToken(varAST.NameToken, "exist_name")
		return
	}
	p.WaitingGlobalVariables = append(p.WaitingGlobalVariables, varAST)
}

func (p *Parser) ParseWaitingGlobalVariables() {
	for _, varAST := range p.WaitingGlobalVariables {
		p.GlobalVariables = append(p.GlobalVariables, p.ParseVariable(varAST))
	}
}

func (p *Parser) ParseVariable(varAST ast.VariableAST) ast.VariableAST {
	if jn.IsIgnoreName(varAST.Name) {
		p.PushErrorToken(varAST.NameToken, "ignore_name_identifier")
	}
	var dt ast.DataTypeAST
	switch t := varAST.Tag.(type) {
	case ast.DataTypeAST:
		dt = t
	default:
		if varAST.SetterToken.Id != lexer.NA {
			var val value
			val, varAST.Value.Model = p.computeExpression(varAST.Value)
			dt = val.ast.Type
		}
	}
	if varAST.Type.Code != jn.Void {
		if varAST.SetterToken.Id != lexer.NA {
			p.checkType(varAST.Type, dt, false, varAST.NameToken)
		} else {
			var valueToken lexer.Token
			valueToken.Id = lexer.Value
			dt, ok := p.readyType(varAST.Type)
			if ok {
				valueToken.Kind = p.defaultValueOfType(dt)
				valueTokens := []lexer.Token{valueToken}
				varAST.Value = ast.ExpressionAST{
					Tokens:    valueTokens,
					Processes: [][]lexer.Token{valueTokens},
				}
			}
		}
	} else {
		if varAST.SetterToken.Id == lexer.NA {
			p.PushErrorToken(varAST.NameToken, "missing_autotype_value")
		} else {
			varAST.Type = dt
			p.checkValidityForAutoType(varAST.Type, varAST.SetterToken)
		}
	}
	if varAST.DefineToken.Kind == "const" {
		if varAST.SetterToken.Id == lexer.NA {
			p.PushErrorToken(varAST.NameToken, "missing_const_value")
		} else if !checkValidityConstantDataType(varAST.Type) {
			p.PushErrorToken(varAST.NameToken, "invalid_const_data_type")
		}
	}
	return varAST
}

func (p *Parser) checkFunctionAttributes(attributes []ast.AttributeAST) {
	for _, attribute := range attributes {
		switch attribute.Tag.Kind {
		case "_inline":
		default:
			p.PushErrorToken(attribute.Token, "invalid_attribute")
		}
	}
}

func variablesFromParameters(params []ast.ParameterAST) []ast.VariableAST {
	var vars []ast.VariableAST
	for _, param := range params {
		var variable ast.VariableAST
		variable.Name = param.Name
		variable.NameToken = param.Token
		variable.Type = param.Type
		if param.Const {
			variable.DefineToken.Id = lexer.Const
		}
		vars = append(vars, variable)
	}
	return vars
}

func (p *Parser) typeByName(name string) *ast.TypeAST {
	for _, t := range p.Types {
		if t.Name == name {
			return &t
		}
	}
	return nil
}

func (p *Parser) FunctionByName(name string) *function {
	for _, fun := range builtinFunctions {
		if fun.Ast.Name == name {
			return fun
		}
	}
	for _, fun := range p.Functions {
		if fun.Ast.Name == name {
			return fun
		}
	}
	return nil
}

func (p *Parser) variableByName(name string) *ast.VariableAST {
	for _, variable := range p.BlockVariables {
		if variable.Name == name {
			return &variable
		}
	}
	for _, variable := range p.GlobalVariables {
		if variable.Name == name {
			return &variable
		}
	}
	return nil
}

func (p *Parser) existName(name string) lexer.Token {
	t := p.typeByName(name)
	if t != nil {
		return t.Token
	}
	fun := p.FunctionByName(name)
	if fun != nil {
		return fun.Ast.Token
	}
	variable := p.variableByName(name)
	if variable != nil {
		return variable.NameToken
	}
	for _, varAST := range p.WaitingGlobalVariables {
		if varAST.Name == name {
			return varAST.NameToken
		}
	}
	return lexer.Token{}
}

func (p *Parser) finalCheck() {
	if p.FunctionByName("_"+jn.EntryPoint) == nil {
		p.PushError("no_entry_point")
	}
	p.checkTypes()
	p.ParseWaitingGlobalVariables()
	p.WaitingGlobalVariables = nil
	p.checkFunctions()
}

func (p *Parser) checkTypes() {
	for _, t := range p.Types {
		_, ok := p.readyType(t.Type)
		if !ok {
			p.PushErrorToken(t.Token, "invalid_type_source")
		}
	}
}

func (p *Parser) checkFunctions() {
	for _, fun := range p.Functions {
		p.BlockVariables = variablesFromParameters(fun.Ast.Params)
		p.checkFunctionSpecialCases(fun)
		p.checkFunction(&fun.Ast)
	}
}

type value struct {
	ast      ast.ValueAST
	constant bool
}

func (p *Parser) computeProcesses(processes [][]lexer.Token) (v value, e expressionModel) {
	if processes == nil {
		return
	}
	builder := newExpBuilder()
	if len(processes) == 1 {
		builder.setIndex(0)
		v = p.processValPart(processes[0], builder)
		e = builder.build()
		return
	}
	process := arithmeticProcess{p: p}
	j := p.nextOperator(processes)
	boolean := false
	for j != -1 {
		if !boolean {
			boolean = v.ast.Type.Code == jn.Bool
		}
		if boolean {
			v.ast.Type.Code = jn.Bool
		}
		if j == 0 {
			process.leftVal = v.ast
			process.operator = processes[j][0]
			builder.setIndex(j + 1)
			builder.appendNode(tokenExpNode{process.operator})
			process.right = processes[j+1]
			builder.setIndex(j + 1)
			process.rightVal = p.processValPart(process.right, builder).ast
			v.ast = process.solve()
			processes = processes[2:]
			goto end
		} else if j == len(processes)-1 {
			process.operator = processes[j][0]
			process.left = processes[j-1]
			builder.setIndex(j - 1)
			process.leftVal = p.processValPart(process.left, builder).ast
			process.rightVal = v.ast
			builder.setIndex(j)
			builder.appendNode(tokenExpNode{process.operator})
			v.ast = process.solve()
			processes = processes[:j-1]
			goto end
		} else if prev := processes[j-1]; prev[0].Id == lexer.Operator &&
			len(prev) == 1 {
			process.leftVal = v.ast
			process.operator = processes[j][0]
			builder.setIndex(j)
			builder.appendNode(tokenExpNode{process.operator})
			process.right = processes[j+1]
			builder.setIndex(j + 1)
			process.rightVal = p.processValPart(process.right, builder).ast
			v.ast = process.solve()
			processes = append(processes[:j], processes[j+2:]...)
			goto end
		}
		process.left = processes[j-1]
		builder.setIndex(j - 1)
		process.leftVal = p.processValPart(process.left, builder).ast
		process.operator = processes[j][0]
		builder.setIndex(j)
		builder.appendNode(tokenExpNode{process.operator})
		process.right = processes[j+1]
		builder.setIndex(j + 1)
		process.rightVal = p.processValPart(process.right, builder).ast
		{
			solvedValue := process.solve()
			if v.ast.Type.Code != jn.Void {
				process.operator.Kind = "+"
				process.leftVal = v.ast
				process.right = processes[j+1]
				process.rightVal = solvedValue
				v.ast = process.solve()
			} else {
				v.ast = solvedValue
			}
		}
		processes = append(processes[:j-1], processes[j+2:]...)
		if len(processes) == 1 {
			break
		}
	end:
		j = p.nextOperator(processes)
	}
	e = builder.build()
	return
}

func (p *Parser) computeTokens(tokens []lexer.Token) (value, expressionModel) {
	return p.computeProcesses(new(ast.AST).BuildExpression(tokens).Processes)
}

func (p *Parser) computeExpression(ex ast.ExpressionAST) (value, expressionModel) {
	processes := make([][]lexer.Token, len(ex.Processes))
	copy(processes, ex.Processes)
	return p.computeProcesses(processes)
}

func (p *Parser) nextOperator(tokens [][]lexer.Token) int {
	precedence5 := -1
	precedence4 := -1
	precedence3 := -1
	precedence2 := -1
	precedence1 := -1
	for index, part := range tokens {
		if len(part) != 1 {
			continue
		} else if part[0].Id != lexer.Operator {
			continue
		}
		switch part[0].Kind {
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
			p.PushErrorToken(part[0], "invalid_operator")
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
	p        *Parser
	left     []lexer.Token
	leftVal  ast.ValueAST
	right    []lexer.Token
	rightVal ast.ValueAST
	operator lexer.Token
}

func (ap arithmeticProcess) solvePointer() (v ast.ValueAST) {
	if ap.leftVal.Type.Value != ap.rightVal.Type.Value {
		ap.p.PushErrorToken(ap.operator, "incompatible_type")
		return
	}
	switch ap.operator.Kind {
	case "!=", "==":
		v.Type.Code = jn.Bool
	default:
		ap.p.PushErrorToken(ap.operator, "operator_notfor_pointer")
	}
	return
}

func (ap arithmeticProcess) solveString() (v ast.ValueAST) {
	if ap.leftVal.Type.Code != ap.rightVal.Type.Code {
		ap.p.PushErrorToken(ap.operator, "incompatible_type")
		return
	}
	switch ap.operator.Kind {
	case "+":
		v.Type.Code = jn.Str
	case "==", "!=":
		v.Type.Code = jn.Bool
	default:
		ap.p.PushErrorToken(ap.operator, "operator_notfor_string")
	}
	return
}

func (ap arithmeticProcess) solveAny() (v ast.ValueAST) {
	switch ap.operator.Kind {
	case "!=", "==":
		v.Type.Code = jn.Bool
	default:
		ap.p.PushErrorToken(ap.operator, "operator_notfor_any")
	}
	return
}

func (ap arithmeticProcess) solveBool() (v ast.ValueAST) {
	if !typesAreCompatible(ap.leftVal.Type, ap.rightVal.Type, true) {
		ap.p.PushErrorToken(ap.operator, "incompatible_type")
	}
	switch ap.operator.Kind {
	case "!=", "==":
		v.Type.Code = jn.Bool
	default:
		ap.p.PushErrorToken(ap.operator, "operator_notfor_bool")
	}
	return
}

func (ap arithmeticProcess) solveFloat() (v ast.ValueAST) {
	if !typesAreCompatible(ap.leftVal.Type, ap.rightVal.Type, true) {
		if !isConstantNumeric(ap.leftVal.Value) && !isConstantNumeric(ap.rightVal.Value) {
			ap.p.PushErrorToken(ap.operator, "incompatible_type")
			return
		}
	}
	switch ap.operator.Kind {
	case "!=", "==", "<", ">", ">=", "<=":
		v.Type.Code = jn.Bool
	case "+", "-", "*", "/":
		v.Type.Code = jn.F32
		if ap.leftVal.Type.Code == jn.F64 || ap.rightVal.Type.Code == jn.F64 {
			v.Type.Code = jn.F64
		}
	default:
		ap.p.PushErrorToken(ap.operator, "operator_notfor_float")
	}
	return
}

func (ap arithmeticProcess) solveSigned() (v ast.ValueAST) {
	if !typesAreCompatible(ap.leftVal.Type, ap.rightVal.Type, true) {
		if !isConstantNumeric(ap.leftVal.Value) && !isConstantNumeric(ap.rightVal.Value) {
			ap.p.PushErrorToken(ap.operator, "incompatible_type")
			return
		}
	}
	switch ap.operator.Kind {
	case "!=", "==", "<", ">", ">=", "<=":
		v.Type.Code = jn.Bool
	case "+", "-", "*", "/", "%", "&", "|", "^":
		v.Type = ap.leftVal.Type
		if jn.TypeGreaterThan(ap.rightVal.Type.Code, v.Type.Code) {
			v.Type = ap.rightVal.Type
		}
	case ">>", "<<":
		v.Type = ap.leftVal.Type
		if !jn.IsUnsignedNumericType(ap.rightVal.Type.Code) &&
			!checkIntBit(ap.rightVal, jnbits.BitsizeOfType(jn.U64)) {
			ap.p.PushErrorToken(ap.rightVal.Token, "bitshift_must_unsigned")
		}
	default:
		ap.p.PushErrorToken(ap.operator, "operator_notfor_int")
	}
	return
}

func (ap arithmeticProcess) solveUnsigned() (v ast.ValueAST) {
	if !typesAreCompatible(ap.leftVal.Type, ap.rightVal.Type, true) {
		if !isConstantNumeric(ap.leftVal.Value) && !isConstantNumeric(ap.rightVal.Value) {
			ap.p.PushErrorToken(ap.operator, "incompatible_type")
			return
		}
		return
	}
	switch ap.operator.Kind {
	case "!=", "==", "<", ">", ">=", "<=":
		v.Type.Code = jn.Bool
	case "+", "-", "*", "/", "%", "&", "|", "^":
		v.Type = ap.leftVal.Type
		if jn.TypeGreaterThan(ap.rightVal.Type.Code, v.Type.Code) {
			v.Type = ap.rightVal.Type
		}
	default:
		ap.p.PushErrorToken(ap.operator, "operator_notfor_uint")
	}
	return
}

func (ap arithmeticProcess) solveLogical() (v ast.ValueAST) {
	v.Type.Code = jn.Bool
	if ap.leftVal.Type.Code != jn.Bool {
		ap.p.PushErrorToken(ap.leftVal.Token, "logical_not_bool")
	}
	if ap.rightVal.Type.Code != jn.Bool {
		ap.p.PushErrorToken(ap.rightVal.Token, "logical_not_bool")
	}
	return
}

func (ap arithmeticProcess) solveRune() (v ast.ValueAST) {
	if !typesAreCompatible(ap.leftVal.Type, ap.rightVal.Type, true) {
		ap.p.PushErrorToken(ap.operator, "incompatible_type")
		return
	}
	switch ap.operator.Kind {
	case "!=", "==", ">", "<", ">=", "<=":
		v.Type.Code = jn.Bool
	case "+", "-", "*", "/", "^", "&", "%", "|":
		v.Type.Code = jn.Rune
	default:
		ap.p.PushErrorToken(ap.operator, "operator_notfor_rune")
	}
	return
}

func (ap arithmeticProcess) solveArray() (v ast.ValueAST) {
	if !typesAreCompatible(ap.leftVal.Type, ap.rightVal.Type, true) {
		ap.p.PushErrorToken(ap.operator, "incompatible_type")
		return
	}
	switch ap.operator.Kind {
	case "!=", "==":
		v.Type.Code = jn.Bool
	default:
		ap.p.PushErrorToken(ap.operator, "operator_notfor_array")
	}
	return
}

func (ap arithmeticProcess) solveNil() (v ast.ValueAST) {
	if !typesAreCompatible(ap.leftVal.Type, ap.rightVal.Type, false) {
		ap.p.PushErrorToken(ap.operator, "incompatible_type")
		return
	}
	switch ap.operator.Kind {
	case "!=", "==":
		v.Type.Code = jn.Bool
	default:
		ap.p.PushErrorToken(ap.operator, "operator_notfor_nil")
	}
	return
}

func (ap arithmeticProcess) solve() (v ast.ValueAST) {
	switch ap.operator.Kind {
	case "+", "-", "*", "/", "%", ">>",
		"<<", "&", "|", "^", "==", "!=",
		">=", "<=", ">", "<":
	case "&&", "||":
		return ap.solveLogical()
	default:
		ap.p.PushErrorToken(ap.operator, "invalid_operator")
	}
	switch {
	case typeIsArray(ap.leftVal.Type) || typeIsArray(ap.rightVal.Type):
		return ap.solveArray()
	case typeIsPointer(ap.leftVal.Type) || typeIsPointer(ap.rightVal.Type):
		return ap.solvePointer()
	case ap.leftVal.Type.Code == jn.Nil || ap.rightVal.Type.Code == jn.Nil:
		return ap.solveNil()
	case ap.leftVal.Type.Code == jn.Rune || ap.rightVal.Type.Code == jn.Rune:
		return ap.solveRune()
	case ap.leftVal.Type.Code == jn.Any || ap.rightVal.Type.Code == jn.Any:
		return ap.solveAny()
	case ap.leftVal.Type.Code == jn.Bool || ap.rightVal.Type.Code == jn.Bool:
		return ap.solveBool()
	case ap.leftVal.Type.Code == jn.Str || ap.rightVal.Type.Code == jn.Str:
		return ap.solveString()
	case jn.IsFloatType(ap.leftVal.Type.Code) ||
		jn.IsFloatType(ap.rightVal.Type.Code):
		return ap.solveFloat()
	case jn.IsSignedNumericType(ap.leftVal.Type.Code) ||
		jn.IsSignedNumericType(ap.rightVal.Type.Code):
		return ap.solveSigned()
	case jn.IsUnsignedNumericType(ap.leftVal.Type.Code) ||
		jn.IsUnsignedNumericType(ap.rightVal.Type.Code):
		return ap.solveUnsigned()
	}
	return
}

type singleValueProcessor struct {
	token   lexer.Token
	builder *expressionModelBuilder
	parser  *Parser
}

func (p *singleValueProcessor) string() value {
	var v value
	v.ast.Value = p.token.Kind
	v.ast.Type.Code = jn.Str
	v.ast.Type.Value = "str"
	p.builder.appendNode(strExpNode{p.token})
	return v
}

func (p *singleValueProcessor) rune() value {
	var v value
	v.ast.Value = p.token.Kind
	v.ast.Type.Code = jn.Rune
	v.ast.Type.Value = "rune"
	p.builder.appendNode(runeExpNode{p.token})
	return v
}

func (p *singleValueProcessor) boolean() value {
	var v value
	v.ast.Value = p.token.Kind
	v.ast.Type.Code = jn.Bool
	v.ast.Type.Value = "bool"
	p.builder.appendNode(tokenExpNode{p.token})
	return v
}

func (p *singleValueProcessor) nil() value {
	var v value
	v.ast.Value = p.token.Kind
	v.ast.Type.Code = jn.Nil
	p.builder.appendNode(tokenExpNode{p.token})
	return v
}

func (p *singleValueProcessor) numeric() value {
	var v value
	if strings.Contains(p.token.Kind, ".") ||
		strings.ContainsAny(p.token.Kind, "eE") {
		v.ast.Type.Code = jn.F64
		v.ast.Type.Value = "f64"
	} else {
		v.ast.Type.Code = jn.I32
		v.ast.Type.Value = "i32"
		ok := jnbits.CheckBitInt(p.token.Kind, 32)
		if !ok {
			v.ast.Type.Code = jn.I64
			v.ast.Type.Value = "i4"
		}
	}
	v.ast.Value = p.token.Kind
	p.builder.appendNode(tokenExpNode{p.token})
	return v
}

func (p *singleValueProcessor) name() (v value, ok bool) {
	if variable := p.parser.variableByName(p.token.Kind); variable != nil {
		v.ast.Value = p.token.Kind
		v.ast.Type = variable.Type
		v.constant = variable.DefineToken.Id == lexer.Const
		v.ast.Token = variable.NameToken
		p.builder.appendNode(tokenExpNode{p.token})
		ok = true
	} else if fun := p.parser.FunctionByName(p.token.Kind); fun != nil {
		v.ast.Value = p.token.Kind
		v.ast.Type.Code = jn.Function
		v.ast.Type.Tag = fun.Ast
		v.ast.Type.Value = fun.Ast.DataTypeString()
		v.ast.Token = fun.Ast.Token
		p.builder.appendNode(tokenExpNode{p.token})
		ok = true
	} else {
		p.parser.PushErrorToken(p.token, "name_not_defined")
	}
	return
}

func (p *Parser) processSingleValPart(
	token lexer.Token,
	builder *expressionModelBuilder,
) (v value, ok bool) {
	processor := singleValueProcessor{
		token:   token,
		builder: builder,
		parser:  p,
	}
	v.ast.Type.Code = jn.Void
	v.ast.Token = token
	switch token.Id {
	case lexer.Value:
		ok = true
		switch {
		case IsString(token.Kind):
			v = processor.string()
		case IsRune(token.Kind):
			v = processor.rune()
		case IsBoolean(token.Kind):
			v = processor.boolean()
		case IsNil(token.Kind):
			v = processor.nil()
		default:
			v = processor.numeric()
		}
	case lexer.Name:
		v, ok = processor.name()
	default:
		p.PushErrorToken(token, "invalid_syntax")
	}
	return
}

type singleOperatorProcessor struct {
	token   lexer.Token
	tokens  []lexer.Token
	builder *expressionModelBuilder
	parser  *Parser
}

func (p *singleOperatorProcessor) unary() value {
	v := p.parser.processValPart(p.tokens, p.builder)
	if !typeIsSingle(v.ast.Type) {
		p.parser.PushErrorToken(p.token, "invalid_data_unary")
	} else if !jn.IsNumericType(v.ast.Type.Code) {
		p.parser.PushErrorToken(p.token, "invalid_data_unary")
	}
	return v
}

func (p *singleOperatorProcessor) plus() value {
	v := p.parser.processValPart(p.tokens, p.builder)
	if !typeIsSingle(v.ast.Type) {
		p.parser.PushErrorToken(p.token, "invalid_data_plus")
	} else if !jn.IsNumericType(v.ast.Type.Code) {
		p.parser.PushErrorToken(p.token, "invalid_data_plus")
	}
	return v
}

func (p *singleOperatorProcessor) tilde() value {
	v := p.parser.processValPart(p.tokens, p.builder)
	if !typeIsSingle(v.ast.Type) {
		p.parser.PushErrorToken(p.token, "invalid_data_tilde")
	} else if !jn.IsIntegerType(v.ast.Type.Code) {
		p.parser.PushErrorToken(p.token, "invalid_data_tilde")
	}
	return v
}

func (p *singleOperatorProcessor) logicalNot() value {
	v := p.parser.processValPart(p.tokens, p.builder)
	if !typeIsSingle(v.ast.Type) {
		p.parser.PushErrorToken(p.token, "invalid_data_logical_not")
	} else if v.ast.Type.Code != jn.Bool {
		p.parser.PushErrorToken(p.token, "invalid_data_logical_not")
	}
	return v
}

func (p *singleOperatorProcessor) star() value {
	v := p.parser.processValPart(p.tokens, p.builder)
	if !typeIsPointer(v.ast.Type) {
		p.parser.PushErrorToken(p.token, "invalid_data_star")
	} else {
		v.ast.Type.Value = v.ast.Type.Value[1:]
	}
	return v
}

func (p *singleOperatorProcessor) amper() value {
	v := p.parser.processValPart(p.tokens, p.builder)
	if !canGetPointer(v) {
		p.parser.PushErrorToken(p.token, "invalid_data_amper")
	}
	v.ast.Type.Value = "*" + v.ast.Type.Value
	return v
}

func (p *Parser) processSingleOperatorPart(
	tokens []lexer.Token,
	builder *expressionModelBuilder,
) value {
	var v value
	processor := singleOperatorProcessor{
		token:   tokens[0],
		tokens:  tokens[1:],
		builder: builder,
		parser:  p,
	}
	builder.appendNode(tokenExpNode{processor.token})
	if processor.tokens == nil {
		p.PushErrorToken(processor.token, "invalid_syntax")
		return v
	}
	switch processor.token.Kind {
	case "-":
		v = processor.unary()
	case "+":
		v = processor.plus()
	case "~":
		v = processor.tilde()
	case "!":
		v = processor.logicalNot()
	case "*":
		v = processor.star()
	case "&":
		v = processor.amper()
	default:
		p.PushErrorToken(processor.token, "invalid_syntax")
	}
	v.ast.Token = processor.token
	return v
}

func canGetPointer(v value) bool {
	if v.ast.Type.Code == jn.Function {
		return false
	}
	return v.ast.Token.Id == lexer.Name
}

func (p *Parser) computeNewHeapAllocation(
	tokens []lexer.Token,
	builder *expressionModelBuilder,
) (v value) {
	if len(tokens) == 1 {
		p.PushErrorToken(tokens[0], "invalid_syntax_keyword_new")
		return
	}
	v.ast.Token = tokens[0]
	tokens = tokens[1:]
	astb := new(ast.AST)
	index := new(int)
	dt, ok := astb.BuildDataType(tokens, index, true)
	builder.appendNode(newHeapAllocationExpModel{dt})
	dt.Value = "*" + dt.Value
	v.ast.Type = dt
	if !ok {
		p.PushErrorToken(tokens[0], "fail_build_heap_allocation_type")
		return
	}
	if *index < len(tokens)-1 {
		p.PushErrorToken(tokens[*index+1], "invalid_syntax")
	}
	return
}

func (p *Parser) processValPart(tokens []lexer.Token, builder *expressionModelBuilder) (v value) {
	if len(tokens) == 1 {
		value, ok := p.processSingleValPart(tokens[0], builder)
		if ok {
			v = value
			return
		}
	}
	firstTok := tokens[0]
	switch firstTok.Id {
	case lexer.Operator:
		return p.processSingleOperatorPart(tokens, builder)
	case lexer.New:
		return p.computeNewHeapAllocation(tokens, builder)
	}
	switch token := tokens[len(tokens)-1]; token.Id {
	case lexer.Brace:
		switch token.Kind {
		case ")":
			return p.processParenthesesValPart(tokens, builder)
		case "}":
			return p.processBraceValPart(tokens, builder)
		case "]":
			return p.processBracketValPart(tokens, builder)
		}
	default:
		p.PushErrorToken(tokens[0], "invalid_syntax")
	}
	return
}

func (p *Parser) processParenthesesValPart(
	tokens []lexer.Token,
	builder *expressionModelBuilder,
) (v value) {
	var valueTokens []lexer.Token
	j := len(tokens) - 1
	braceCount := 0
	for ; j >= 0; j-- {
		token := tokens[j]
		if token.Id != lexer.Brace {
			continue
		}
		switch token.Kind {
		case ")", "}", "]":
			braceCount++
		case "(", "{", "[":
			braceCount--
		}
		if braceCount > 0 {
			continue
		}
		valueTokens = tokens[:j]
		break
	}
	if len(valueTokens) == 0 && braceCount == 0 {
		builder.appendNode(tokenExpNode{lexer.Token{Kind: "("}})
		defer builder.appendNode(tokenExpNode{lexer.Token{Kind: ")"}})

		tk := tokens[0]
		tokens = tokens[1 : len(tokens)-1]
		if len(tokens) == 0 {
			p.PushErrorToken(tk, "invalid_syntax")
		}
		value, model := p.computeTokens(tokens)
		v = value
		builder.appendNode(model)
		return
	}
	v = p.processValPart(valueTokens, builder)

	builder.appendNode(tokenExpNode{lexer.Token{Kind: "("}})
	defer builder.appendNode(tokenExpNode{lexer.Token{Kind: ")"}})

	switch v.ast.Type.Code {
	case jn.Function:
		fun := v.ast.Type.Tag.(ast.FunctionAST)
		p.parseFunctionCallStatement(fun, tokens[len(valueTokens):], builder)
		v.ast.Type = fun.ReturnType
	default:
		p.PushErrorToken(tokens[len(valueTokens)], "invalid_syntax")
	}
	return
}

func (p *Parser) processBraceValPart(
	tokens []lexer.Token,
	builder *expressionModelBuilder,
) (v value) {
	var valueTokens []lexer.Token
	j := len(tokens) - 1
	braceCount := 0
	for ; j >= 0; j-- {
		token := tokens[j]
		if token.Id != lexer.Brace {
			continue
		}
		switch token.Kind {
		case "}", "]", ")":
			braceCount++
		case "{", "(", "[":
			braceCount--
		}
		if braceCount > 0 {
			continue
		}
		valueTokens = tokens[:j]
		break
	}
	valTokensLen := len(valueTokens)
	if valTokensLen == 0 || braceCount > 0 {
		p.PushErrorToken(tokens[0], "invalid_syntax")
		return
	}
	switch valueTokens[0].Id {
	case lexer.Brace:
		switch valueTokens[0].Kind {
		case "[":
			ast := ast.New(nil)
			dt, ok := ast.BuildDataType(valueTokens, new(int), true)
			if !ok {
				p.AppendErrors(ast.Errors...)
				return
			}
			valueTokens = tokens[len(valueTokens):]
			var model expressionNode
			v, model = p.buildArray(p.buildEnumerableParts(valueTokens), dt, valueTokens[0])
			builder.appendNode(model)
			return
		case "(":
			astBuilder := ast.New(tokens)
			funAST := astBuilder.BuildFunction(true)
			if len(astBuilder.Errors) > 0 {
				p.AppendErrors(astBuilder.Errors...)
				return
			}
			p.checkAnonymousFunction(&funAST)
			v.ast.Type.Tag = funAST
			v.ast.Type.Code = jn.Function
			v.ast.Type.Value = funAST.DataTypeString()
			builder.appendNode(anonymousFunctionExp{funAST})
			return
		default:
			p.PushErrorToken(valueTokens[0], "invalid_syntax")
		}
	default:
		p.PushErrorToken(valueTokens[0], "invalid_syntax")
	}
	return
}

func (p *Parser) processBracketValPart(
	tokens []lexer.Token,
	builder *expressionModelBuilder,
) (v value) {
	var valueTokens []lexer.Token
	j := len(tokens) - 1
	braceCount := 0
	for ; j >= 0; j-- {
		token := tokens[j]
		if token.Id != lexer.Brace {
			continue
		}
		switch token.Kind {
		case "}", "]", ")":
			braceCount++
		case "{", "(", "[":
			braceCount--
		}
		if braceCount > 0 {
			continue
		}
		valueTokens = tokens[:j]
		break
	}
	valTokensLen := len(valueTokens)
	if valTokensLen == 0 || braceCount > 0 {
		p.PushErrorToken(tokens[0], "invalid_syntax")
		return
	}
	var model expressionNode
	v, model = p.computeTokens(valueTokens)
	builder.appendNode(model)
	tokens = tokens[len(valueTokens)+1 : len(tokens)-1]
	builder.appendNode(tokenExpNode{lexer.Token{Kind: "["}})
	selectv, model := p.computeTokens(tokens)
	builder.appendNode(model)
	builder.appendNode(tokenExpNode{lexer.Token{Kind: "]"}})
	return p.processEnumerableSelect(v, selectv, tokens[0])
}

func (p *Parser) processEnumerableSelect(enumv, selectv value, err lexer.Token) (v value) {
	switch {
	case typeIsArray(enumv.ast.Type):
		return p.processArraySelect(enumv, selectv, err)
	case typeIsSingle(enumv.ast.Type):
		return p.processStringSelect(enumv, selectv, err)
	}
	p.PushErrorToken(err, "not_enumerable")
	return
}

func (p *Parser) processArraySelect(arrv, selectv value, err lexer.Token) value {
	arrv.ast.Type = typeOfArrayElements(arrv.ast.Type)
	if !typeIsSingle(selectv.ast.Type) || !jn.IsIntegerType(selectv.ast.Type.Code) {
		p.PushErrorToken(err, "notint_array_select")
	}
	return arrv
}

func (p *Parser) processStringSelect(strv, selectv value, err lexer.Token) value {
	strv.ast.Type.Code = jn.Rune
	if !typeIsSingle(selectv.ast.Type) || !jn.IsIntegerType(selectv.ast.Type.Code) {
		p.PushErrorToken(err, "notint_string_select")
	}
	return strv
}

type enumPart struct {
	tokens []lexer.Token
}

func (p *Parser) buildEnumerableParts(tokens []lexer.Token) []enumPart {
	tokens = tokens[1 : len(tokens)-1]
	braceCount := 0
	lastcComma := -1
	var parts []enumPart
	for index, token := range tokens {
		if token.Id == lexer.Brace {
			switch token.Kind {
			case "{", "[", "(":
				braceCount++
			default:
				braceCount--
			}
		}
		if braceCount > 0 {
			continue
		}
		if token.Id == lexer.Comma {
			if index-lastcComma-1 == 0 {
				p.PushErrorToken(token, "missing_expression")
				lastcComma = index
				continue
			}
			parts = append(parts, enumPart{tokens[lastcComma+1 : index]})
			lastcComma = index
		}
	}
	if lastcComma+1 < len(tokens) {
		parts = append(parts, enumPart{tokens[lastcComma+1:]})
	}
	return parts
}

func (p *Parser) buildArray(
	parts []enumPart,
	dt ast.DataTypeAST,
	err lexer.Token,
) (value, expressionNode) {
	var v value
	v.ast.Type = dt
	model := arrayExp{dataType: dt}
	elementType := typeOfArrayElements(dt)
	for _, part := range parts {
		partValue, expModel := p.computeTokens(part.tokens)
		model.expressions = append(model.expressions, expModel)
		p.checkType(elementType, partValue.ast.Type, false, part.tokens[0])
	}
	return v, model
}

func (p *Parser) checkAnonymousFunction(fun *ast.FunctionAST) {
	globalVariables := p.GlobalVariables
	blockVariables := p.BlockVariables
	p.GlobalVariables = append(blockVariables, p.GlobalVariables...)
	p.BlockVariables = variablesFromParameters(fun.Params)
	p.checkFunction(fun)
	p.GlobalVariables = globalVariables
	p.BlockVariables = blockVariables
}

func (p *Parser) parseFunctionCallStatement(
	fun ast.FunctionAST,
	tokens []lexer.Token,
	builder *expressionModelBuilder,
) {
	errToken := tokens[0]
	tokens, _ = p.getRangeTokens("(", ")", tokens)
	if tokens == nil {
		tokens = make([]lexer.Token, 0)
	}
	ast := new(ast.AST)
	args := ast.BuildArgs(tokens)
	if len(ast.Errors) > 0 {
		p.AppendErrors(ast.Errors...)
	}
	p.parseArgs(fun, args, errToken, builder)
	if builder != nil {
		builder.appendNode(argsExp{args})
	}
}

func (p *Parser) parseArgs(
	fun ast.FunctionAST,
	args []ast.ArgAST,
	errToken lexer.Token,
	builder *expressionModelBuilder,
) {
	if len(args) < len(fun.Params) {
		p.PushErrorToken(errToken, "missing_argument")
	}
	for index, arg := range args {
		p.parseArg(fun, index, &arg)
		args[index] = arg
	}
}

func (p *Parser) parseArg(fun ast.FunctionAST, index int, arg *ast.ArgAST) {
	if index >= len(fun.Params) {
		p.PushErrorToken(arg.Token, "argument_overflow")
		return
	}
	value, model := p.computeExpression(arg.Expression)
	arg.Expression.Model = model
	param := fun.Params[index]
	p.checkType(param.Type, value.ast.Type, false, arg.Token)
}

func (p *Parser) getRangeTokens(
	open, close string,
	tokens []lexer.Token,
) (_ []lexer.Token, ok bool) {
	braceCount := 0
	start := 1
	if tokens[0].Id != lexer.Brace {
		return nil, false
	}
	for index, token := range tokens {
		if token.Id != lexer.Brace {
			continue
		}
		if token.Kind == open {
			braceCount++
		} else if token.Kind == close {
			braceCount--
		}
		if braceCount > 0 {
			continue
		}
		return tokens[start:index], true
	}
	p.PushErrorToken(tokens[0], "brace_not_closed")
	return nil, false
}

func (p *Parser) checkFunctionSpecialCases(fun *function) {
	switch fun.Ast.Name {
	case "_" + jn.EntryPoint:
		p.checkEntryPointSpecialCases(fun)
	}
}

func (p *Parser) checkEntryPointSpecialCases(fun *function) {
	if len(fun.Ast.Params) > 0 {
		p.PushErrorToken(fun.Ast.Token, "entrypoint_have_parameters")
	}
	if fun.Ast.ReturnType.Code != jn.Void {
		p.PushErrorToken(fun.Ast.ReturnType.Token, "entrypoint_have_return")
	}
	if fun.Attributes != nil {
		p.PushErrorToken(fun.Ast.Token, "entrypoint_have_attributes")
	}
}

func (p *Parser) checkBlock(b *ast.BlockAST) {
	for index := 0; index < len(b.Statements); index++ {
		model := &b.Statements[index]
		switch t := model.Value.(type) {
		case ast.BlockExpressionAST:
			_, t.Expression.Model = p.computeExpression(t.Expression)
			model.Value = t
		case ast.VariableAST:
			p.checkVariableStatement(&t, false)
			model.Value = t
		case ast.VariableSetAST:
			p.checkVarsetStatement(&t)
			model.Value = t
		case ast.FreeAST:
			p.checkFreeStatement(&t)
			model.Value = t
		case ast.ReturnAST:
		default:
			p.PushErrorToken(model.Token, "invalid_syntax")
		}
	}
}

func (p *Parser) checkParameters(params []ast.ParameterAST) {
	for _, param := range params {
		if !param.Const {
			continue
		}
		if !checkValidityConstantDataType(param.Type) {
			p.PushErrorToken(param.Type.Token, "invalid_const_data_type")
		}
	}
}

type returnChecker struct {
	p        *Parser
	retAST   *ast.ReturnAST
	fun      ast.FunctionAST
	expModel multiReturnExpModel
	values   []value
}

func (rc *returnChecker) pushValue(last, current int, errTk lexer.Token) {
	if current-last == 0 {
		rc.p.PushErrorToken(errTk, "missing_value")
		return
	}
	tokens := rc.retAST.Expression.Tokens[last:current]
	value, model := rc.p.computeTokens(tokens)
	rc.expModel.models = append(rc.expModel.models, model)
	rc.values = append(rc.values, value)
}

func (rc *returnChecker) checkValues() {
	braceCount := 0
	last := 0
	for index, token := range rc.retAST.Expression.Tokens {
		if token.Id == lexer.Brace {
			switch token.Kind {
			case "(", "{", "[":
			default:
				braceCount--
			}
		}
		if braceCount > 0 || token.Id != lexer.Comma {
			continue
		}
		rc.pushValue(last, index, token)
		last = index + 1
	}
	length := len(rc.retAST.Expression.Tokens)
	if last < length {
		if last == 0 {
			rc.pushValue(0, length, rc.retAST.Token)
		} else {
			rc.pushValue(last, length, rc.retAST.Expression.Tokens[last-1])
		}
	}
	if !typeIsVoidReturn(rc.fun.ReturnType) {
		rc.checkValueTypes()
	}
}

func (rc *returnChecker) checkValueTypes() {
	valLength := len(rc.values)
	if !rc.fun.ReturnType.MultiTyped {
		rc.retAST.Expression.Model = rc.expModel.models[0]
		if valLength > 1 {
			rc.p.PushErrorToken(rc.retAST.Token, "overflow_return")
		}
		rc.p.checkType(rc.fun.ReturnType, rc.values[0].ast.Type, true, rc.retAST.Token)
		return
	}
	rc.retAST.Expression.Model = rc.expModel
	types := rc.fun.ReturnType.Tag.([]ast.DataTypeAST)
	if valLength == 1 {
		rc.p.PushErrorToken(rc.retAST.Token, "missing_multi_return")
	} else if valLength > len(types) {
		rc.p.PushErrorToken(rc.retAST.Token, "overflow_return")
	}
	for index, t := range types {
		if index >= valLength {
			break
		}
		rc.p.checkType(t, rc.values[index].ast.Type, true, rc.retAST.Token)
	}
}

func (rc *returnChecker) check() {
	exprTokensLen := len(rc.retAST.Expression.Tokens)
	if exprTokensLen == 0 && !typeIsVoidReturn(rc.fun.ReturnType) {
		rc.p.PushErrorToken(rc.retAST.Token, "require_return_value")
		return
	}
	if exprTokensLen > 0 && typeIsVoidReturn(rc.fun.ReturnType) {
		rc.p.PushErrorToken(rc.retAST.Token, "void_function_return_value")
	}
	rc.checkValues()
}

func (p *Parser) checkReturns(fun ast.FunctionAST) {
	missed := true
	for index, s := range fun.Block.Statements {
		switch t := s.Value.(type) {
		case ast.ReturnAST:
			rc := returnChecker{p: p, retAST: &t, fun: fun}
			rc.check()
			fun.Block.Statements[index].Value = t
			missed = false
		}
	}
	if missed && !typeIsVoidReturn(fun.ReturnType) {
		p.PushErrorToken(fun.Token, "missing_return")
	}
}

func (p *Parser) checkFunction(fun *ast.FunctionAST) {
	p.checkBlock(&fun.Block)
	p.checkReturns(*fun)
	p.checkParameters(fun.Params)
}

func (p *Parser) checkVariableStatement(varAST *ast.VariableAST, noParse bool) {
	for _, t := range p.Types {
		if varAST.Name == t.Name {
			p.PushErrorToken(varAST.NameToken, "exist_name")
			break
		}
	}
	for _, variable := range p.BlockVariables {
		if varAST.Name == variable.Name {
			p.PushErrorToken(varAST.NameToken, "exist_name")
			break
		}
	}
	if !noParse {
		*varAST = p.ParseVariable(*varAST)
	}
	p.BlockVariables = append(p.BlockVariables, *varAST)
}

func (p *Parser) checkVarsetOperation(selected value, err lexer.Token) bool {
	if selected.constant {
		p.PushErrorToken(err, "const_value_update")
		return false
	}
	switch selected.ast.Type.Tag.(type) {
	case ast.FunctionAST:
		if p.FunctionByName(selected.ast.Token.Kind) != nil {
			p.PushErrorToken(err, "type_not_support_value_update")
			return false
		}
	}
	return true
}

func (p *Parser) checkOneVarset(vsAST *ast.VariableSetAST) {
	selected, _ := p.computeExpression(vsAST.SelectExpressions[0].Expression)
	if !p.checkVarsetOperation(selected, vsAST.Setter) {
		return
	}
	value, model := p.computeExpression(vsAST.ValueExpressions[0])
	vsAST.ValueExpressions[0] = model.ExpressionAST()
	if vsAST.Setter.Kind != "=" {
		vsAST.Setter.Kind = vsAST.Setter.Kind[:len(vsAST.Setter.Kind)-1]
		value.ast = arithmeticProcess{
			p:        p,
			left:     vsAST.SelectExpressions[0].Expression.Tokens,
			leftVal:  selected.ast,
			right:    vsAST.ValueExpressions[0].Tokens,
			rightVal: value.ast,
			operator: vsAST.Setter,
		}.solve()
		vsAST.Setter.Kind += "="
	}
	p.checkType(selected.ast.Type, value.ast.Type, false, vsAST.Setter)
}

func (p *Parser) parseVarsetSelections(vsAST *ast.VariableSetAST) {
	for index, selector := range vsAST.SelectExpressions {
		p.checkVariableStatement(&selector.Variable, false)
		vsAST.SelectExpressions[index] = selector
	}
}

func (p *Parser) getVarsetTypes(vsAST *ast.VariableSetAST) []ast.DataTypeAST {
	values := make([]ast.DataTypeAST, len(vsAST.ValueExpressions))
	for index, expression := range vsAST.ValueExpressions {
		val, model := p.computeExpression(expression)
		vsAST.ValueExpressions[index].Model = model
		values[index] = val.ast.Type
	}
	return values
}

func (p *Parser) processFuncMultiVarset(vsAST *ast.VariableSetAST, funcVal value) {
	types := funcVal.ast.Type.Tag.([]ast.DataTypeAST)
	if len(types) != len(vsAST.SelectExpressions) {
		p.PushErrorToken(vsAST.Setter, "missing_multiassign_identifiers")
		return
	}
	p.processMultiVarset(vsAST, types)
}

func (p *Parser) processMultiVarset(vsAST *ast.VariableSetAST, types []ast.DataTypeAST) {
	for index := range vsAST.SelectExpressions {
		selector := &vsAST.SelectExpressions[index]
		selector.Ignore = jn.IsIgnoreName(selector.Variable.Name)
		dt := types[index]
		if !selector.NewVariable {
			if selector.Ignore {
				continue
			}
			selected, _ := p.computeExpression(selector.Expression)
			if !p.checkVarsetOperation(selected, vsAST.Setter) {
				return
			}
			p.checkType(selected.ast.Type, dt, false, vsAST.Setter)
			continue
		}
		selector.Variable.Tag = dt
		p.checkVariableStatement(&selector.Variable, false)
	}
}

func (p *Parser) checkVarsetStatement(vsAST *ast.VariableSetAST) {
	selectLength := len(vsAST.SelectExpressions)
	valueLength := len(vsAST.ValueExpressions)
	if vsAST.JustDeclare {
		p.parseVarsetSelections(vsAST)
		return
	} else if selectLength == 1 && !vsAST.SelectExpressions[0].NewVariable {
		p.checkOneVarset(vsAST)
		return
	} else if vsAST.Setter.Kind != "=" {
		p.PushErrorToken(vsAST.Setter, "invalid_syntax")
		return
	}
	if valueLength == 1 {
		firstVal, _ := p.computeExpression(vsAST.ValueExpressions[0])
		if firstVal.ast.Type.MultiTyped {
			vsAST.MultipleReturn = true
			p.processFuncMultiVarset(vsAST, firstVal)
			return
		}
	}
	switch {
	case selectLength > valueLength:
		p.PushErrorToken(vsAST.Setter, "overflow_multiassign_identifiers")
		return
	case selectLength < valueLength:
		p.PushErrorToken(vsAST.Setter, "missing_multiassign_identifiers")
		return
	}
	p.processMultiVarset(vsAST, p.getVarsetTypes(vsAST))
}

func (p *Parser) checkFreeStatement(freeAST *ast.FreeAST) {
	val, model := p.computeExpression(freeAST.Expression)
	freeAST.Expression.Model = model
	if !typeIsPointer(val.ast.Type) {
		p.PushErrorToken(freeAST.Token, "free_nonpointer")
	}
}
