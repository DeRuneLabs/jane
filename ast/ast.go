package ast

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/De-Rune/jane/lexer"
	"github.com/De-Rune/jane/package/jn"
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
	message := jn.Errors[err]
	ast.Errors = append(
		ast.Errors,
		fmt.Sprintf("%s:%d:%d %s", token.File.Path, token.Row, token.Column, message),
	)
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
		switch firstToken.Id {
		case lexer.At:
			ast.BuildAttribute()
		case lexer.Name:
			ast.BuildName()
		case lexer.Const:
			ast.BuildGlobalVariable()
		case lexer.Type:
			ast.BuildType()
		default:
			ast.PushError("invalid_syntax")
			ast.Position++
		}
	}
}

func (ast *AST) BuildType() {
	position := 1
	tokens := ast.skipStatement()
	if position >= len(tokens) {
		ast.PushErrorToken(tokens[position-1], "invalid_syntax")
		return
	}
	token := tokens[position]
	if token.Id != lexer.Name {
		ast.PushErrorToken(token, "invalid_syntax")
	}
	position++
	if position >= len(tokens) {
		ast.PushErrorToken(tokens[position-1], "invalid_syntax")
		return
	}
	destinationType, _ := ast.BuildDataType(tokens[position:], new(int), true)
	ast.Tree = append(ast.Tree, Object{
		Token: tokens[1],
		Value: TypeAST{
			Token: tokens[1],
			Name:  tokens[1].Kind,
			Type:  destinationType,
		},
	})
}

func (ast *AST) BuildName() {
	ast.Position++
	if ast.Ended() {
		ast.PushErrorToken(ast.Tokens[ast.Position-1], "invalid_syntax")
		return
	}
	token := ast.Tokens[ast.Position]
	ast.Position--
	switch token.Id {
	case lexer.Colon:
		ast.BuildGlobalVariable()
	case lexer.Brace:
		switch token.Kind {
		case "(":
			funAST := ast.BuildFunction(false)
			ast.Tree = append(ast.Tree, Object{
				Token: funAST.Token,
				Value: StatementAST{
					Token: funAST.Token,
					Value: funAST,
				},
			})
			return
		}
	}
	ast.Position++
	ast.PushErrorToken(token, "invalid_syntax")
}

func (ast *AST) BuildAttribute() {
	var attribute AttributeAST
	attribute.Token = ast.Tokens[ast.Position]
	ast.Position++
	if ast.Ended() {
		ast.PushErrorToken(ast.Tokens[ast.Position-1], "invalid_syntax")
		return
	}
	attribute.Tag = ast.Tokens[ast.Position]
	if attribute.Tag.Id != lexer.Name || attribute.Token.Column+1 != attribute.Tag.Column {
		ast.PushErrorToken(attribute.Tag, "invalid_syntax")
		ast.Position = -1
		return
	}
	ast.Tree = append(ast.Tree, Object{
		Token: attribute.Token,
		Value: attribute,
	})
	ast.Position++
}

func (ast *AST) BuildFunction(anonyomus bool) (funAST FunctionAST) {
	funAST.Token = ast.Tokens[ast.Position]
	if anonyomus {
		funAST.Name = "anonyomus"
	} else {
		if funAST.Token.Id != lexer.Name {
			ast.PushErrorToken(funAST.Token, "invalid_syntax")
		}
		funAST.Name = funAST.Token.Kind
		ast.Position++
		if ast.Ended() {
			ast.Position--
			ast.PushError("function_body_not_exist")
			ast.Position = -1
			return
		}
	}
	funAST.ReturnType.Code = jn.Void
	tokens := ast.getRange("(", ")")
	if tokens == nil {
		return
	} else if len(tokens) > 0 {
		ast.BuildParameters(&funAST, tokens)
	}
	if ast.Ended() {
		ast.Position--
		ast.PushError("function_body_not_exist")
		ast.Position = -1
		return
	}
	token := ast.Tokens[ast.Position]
	t, ok := ast.BuildFunctionReturnDataType(ast.Tokens, &ast.Position)
	if ok {
		funAST.ReturnType = t
		ast.Position++
		if ast.Ended() {
			ast.Position--
			ast.PushError("function_body_not_exist")
			ast.Position = -1
			return
		}
		token = ast.Tokens[ast.Position]
	}
	if token.Id != lexer.Brace || token.Kind != "{" {
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
	funAST.Block = ast.BuildBlock(blockTokens)
	return
}

func (ast *AST) BuildGlobalVariable() {
	statementTokens := ast.skipStatement()
	if statementTokens == nil {
		return
	}
	statement := ast.BuildVariableStatement(statementTokens)
	ast.Tree = append(ast.Tree, Object{
		Token: statement.Token,
		Value: statement,
	})
}

func (ast *AST) BuildParameters(fn *FunctionAST, tokens []lexer.Token) {
	last := 0
	braceCount := 0
	for index, token := range tokens {
		if token.Id == lexer.Brace {
			switch token.Kind {
			case "{", "[", "(":
				braceCount++
			default:
				braceCount--
			}
		}
		if braceCount > 0 || token.Id != lexer.Comma {
			continue
		}
		ast.pushParameter(fn, tokens[last:index], token)
		last = index + 1
	}
	if last < len(tokens) {
		if last == 0 {
			ast.pushParameter(fn, tokens[last:], tokens[last])
		} else {
			ast.pushParameter(fn, tokens[last:], tokens[last-1])
		}
	}
}

func (ast *AST) pushParameter(fn *FunctionAST, tokens []lexer.Token, err lexer.Token) {
	if len(tokens) == 0 {
		ast.PushErrorToken(err, "invalid_syntax")
		return
	}
	paramAST := ParameterAST{
		Token: tokens[0],
	}
	for index, token := range tokens {
		switch token.Id {
		case lexer.Const:
			if paramAST.Const {
				ast.PushErrorToken(token, "already_constant")
				continue
			}
			paramAST.Const = true
		case lexer.Name:
			tokens = tokens[index:]
			if len(tokens) < 2 {
				ast.PushErrorToken(paramAST.Token, "missing_type")
				return
			}
			if !jn.IsIgnoreName(token.Kind) {
				for _, param := range fn.Params {
					if param.Name == token.Kind {
						ast.PushErrorToken(token, "parameter_exist")
						break
					}
				}
				paramAST.Name = token.Kind
			}
			index := 1
			paramAST.Type, _ = ast.BuildDataType(tokens, &index, true)
			if index+1 < len(tokens) {
				ast.PushErrorToken(tokens[index+1], "invalid_syntax")
			}
			goto end
		default:
			if t, ok := ast.BuildDataType(tokens, &index, true); ok {
				if index+1 == len(tokens) {
					paramAST.Type = t
					goto end
				}
			}
			ast.PushErrorToken(token, "invalid_syntax")
			goto end
		}
	}
end:
	if paramAST.Type.Code == jn.Void {
		ast.PushErrorToken(paramAST.Token, "invalid_syntax")
	}
	fn.Params = append(fn.Params, paramAST)
}

func (ast *AST) BuildDataType(
	tokens []lexer.Token,
	index *int,
	err bool,
) (dt DataTypeAST, ok bool) {
	first := *index
	for ; *index < len(tokens); *index++ {
		token := tokens[*index]
		switch token.Id {
		case lexer.DataType:
			BuildDataType(token, &dt)
			return dt, true
		case lexer.Name:
			buildNameType(token, &dt)
			return dt, true
		case lexer.Operator:
			if token.Kind == "*" {
				dt.Value += token.Kind
				break
			}
			if err {
				ast.PushErrorToken(token, "invalid_syntax")
			}
			return dt, false
		case lexer.Brace:
			switch token.Kind {
			case "(":
				ast.buildFunctionDataType(token, tokens, index, &dt)
				return dt, true
			case "[":
				*index++
				if *index > len(tokens) {
					if err {
						ast.PushErrorToken(token, "invalid_syntax")
					}
					return dt, false
				}
				token = tokens[*index]
				if token.Id != lexer.Brace || token.Kind != "]" {
					if err {
						ast.PushErrorToken(token, "invalid_syntax")
					}
					return dt, false
				}
				dt.Value += "[]"
				continue
			}
			return dt, false
		default:
			if err {
				ast.PushErrorToken(token, "invalid_syntax")
			}
			return dt, false
		}
	}
	if err {
		ast.PushErrorToken(tokens[first], "invalid_type")
	}
	return dt, false
}

func BuildDataType(token lexer.Token, dt *DataTypeAST) {
	dt.Token = token
	dt.Code = jn.TypeFromName(dt.Token.Kind)
	dt.Value += dt.Token.Kind
}

func buildNameType(token lexer.Token, dt *DataTypeAST) {
	dt.Token = token
	dt.Code = jn.Name
	dt.Value += dt.Token.Kind
}

func (ast *AST) buildFunctionDataType(
	token lexer.Token,
	tokens []lexer.Token,
	index *int,
	dt *DataTypeAST,
) {
	dt.Token = token
	dt.Code = jn.Function
	value, funAST := ast.buildFunctionDataTypeHead(tokens, index)
	funAST.ReturnType, _ = ast.BuildFunctionReturnDataType(tokens, index)
	dt.Value += value
	dt.Tag = funAST
}

func (ast *AST) buildFunctionDataTypeHead(tokens []lexer.Token, index *int) (string, FunctionAST) {
	var funAST FunctionAST
	var typeValue strings.Builder
	typeValue.WriteByte('{')
	brace := 1
	firstIndex := *index
	for *index++; *index < len(tokens); *index++ {
		token := tokens[*index]
		typeValue.WriteString(token.Kind)
		switch token.Id {
		case lexer.Brace:
			switch token.Kind {
			case "{", "[", "(":
				brace++
			default:
				brace--
			}
		}
		if brace == 0 {
			ast.BuildParameters(&funAST, tokens[firstIndex+1:*index])
			return typeValue.String(), funAST
		}
	}
	ast.PushErrorToken(tokens[firstIndex], "invalid_type")
	return "", funAST
}

func (ast *AST) pushTypeToTypes(types *[]DataTypeAST, tokens []lexer.Token, errToken lexer.Token) {
	if len(tokens) == 0 {
		ast.PushErrorToken(errToken, "missing_value")
		return
	}
	currentDt, _ := ast.BuildDataType(tokens, new(int), false)
	*types = append(*types, currentDt)
}

func (ast *AST) BuildFunctionReturnDataType(
	tokens []lexer.Token,
	index *int,
) (dt DataTypeAST, ok bool) {
	if *index >= len(tokens) {
		goto end
	}
	if tokens[*index].Id == lexer.Brace && tokens[*index].Kind == "[" {
		*index++
		if *index >= len(tokens) {
			*index--
			goto end
		}
		if tokens[*index].Id == lexer.Brace && tokens[*index].Kind == "]" {
			*index--
			goto end
		}
		var types []DataTypeAST
		braceCount := 1
		last := *index
		for ; *index < len(tokens); *index++ {
			token := tokens[*index]
			if token.Id == lexer.Brace {
				switch token.Kind {
				case "(", "[", "{":
					braceCount++
				default:
					braceCount--
				}
			}
			if braceCount == 0 {
				ast.pushTypeToTypes(&types, tokens[last:*index], tokens[last-1])
				break
			} else if braceCount > 1 {
				continue
			}
			if token.Id != lexer.Comma {
				continue
			}
			ast.pushTypeToTypes(&types, tokens[last:*index], tokens[*index-1])
			last = *index + 1
		}
		if len(types) > 1 {
			dt.MultiTyped = true
			dt.Tag = types
		} else {
			dt = types[0]
		}
		ok = true
		return
	}
end:
	return ast.BuildDataType(tokens, index, false)
}

func IsSingleOperator(operator string) bool {
	return operator == "-" ||
		operator == "+" ||
		operator == "~" ||
		operator == "!" ||
		operator == "*" ||
		operator == "&"
}

func IsStatement(current, prev lexer.Token) (yes bool, withSemicolon bool) {
	yes = current.Id == lexer.SemiColon || prev.Row < current.Row
	if yes {
		withSemicolon = current.Id == lexer.SemiColon
	}
	return
}

func (ast *AST) pushStatementToBlock(b *BlockAST, tokens []lexer.Token) {
	if len(tokens) == 0 {
		return
	}
	if tokens[len(tokens)-1].Id == lexer.SemiColon {
		if len(tokens) == 1 {
			return
		}
		tokens = tokens[:len(tokens)-1]
	}
	b.Statements = append(b.Statements, ast.BuildStatement(tokens))
}

func nextStatementPos(tokens []lexer.Token, start int) int {
	braceCount := 0
	index := start
	for ; index < len(tokens); index++ {
		token := tokens[index]
		if token.Id == lexer.Brace {
			switch token.Kind {
			case "{", "[", "(":
				braceCount++
				continue
			default:
				braceCount--
				continue
			}
		}
		if braceCount > 0 {
			continue
		}
		var isStatement, withSemicolon bool
		if index > start {
			isStatement, withSemicolon = IsStatement(token, tokens[index-1])
		} else {
			isStatement, withSemicolon = IsStatement(token, token)
		}
		if !isStatement {
			continue
		}
		if withSemicolon {
			index++
		}
		return index
	}
	return index
}

func (ast *AST) BuildBlock(tokens []lexer.Token) (b BlockAST) {
	var index, start int
	for {
		if ast.Position == -1 {
			return
		}
		index = nextStatementPos(tokens, index)
		ast.pushStatementToBlock(&b, tokens[start:index])
		if index >= len(tokens) {
			break
		}
		start = index
	}
	return
}

func (ast *AST) BuildStatement(tokens []lexer.Token) (s StatementAST) {
	s, ok := ast.BuildVariableSetStatement(tokens)
	if ok {
		return s
	}
	firstToken := tokens[0]
	switch firstToken.Id {
	case lexer.Name:
		return ast.BuildNameStatement(tokens)
	case lexer.Const:
		return ast.BuildVariableStatement(tokens)
	case lexer.Return:
		return ast.BuildReturnStatement(tokens)
	case lexer.Brace:
		if firstToken.Kind == "(" {
			return ast.BuildExpressionStatement(tokens)
		}
	case lexer.Operator:
		if firstToken.Kind == "<" {
			return ast.BuildReturnStatement(tokens)
		}
	}
	ast.PushErrorToken(firstToken, "invalid_syntax")
	return
}

func checkVariableSetStatementTokens(tokens []lexer.Token) bool {
	braceCount := 0
	for _, token := range tokens {
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
		if token.Id == lexer.Operator && token.Kind[len(token.Kind)-1] == '=' {
			return true
		}
	}
	return false
}

type varsetInfo struct {
	selectorTokens   []lexer.Token
	expressionTokens []lexer.Token
	setter           lexer.Token
	ok               bool
	JustDeclare      bool
}

func (ast *AST) variableSetInfo(tokens []lexer.Token) (info varsetInfo) {
	info.ok = true
	braceCount := 0
	for index, token := range tokens {
		if token.Id == lexer.Brace {
			switch token.Kind {
			case "(", "[", "{":
				braceCount++
			default:
				braceCount--
			}
		}
		if braceCount > 0 {
			continue
		}
		if token.Id == lexer.Operator && token.Kind[len(token.Kind)-1] == '=' {
			info.selectorTokens = tokens[:index]
			if info.selectorTokens == nil {
				ast.PushErrorToken(token, "invalid_syntax")
				info.ok = false
			}
			info.setter = token
			if index+1 >= len(tokens) {
				ast.PushErrorToken(token, "missing_value")
				info.ok = false
			} else {
				info.expressionTokens = tokens[index+1:]
			}
			return
		}
	}
	info.JustDeclare = true
	info.selectorTokens = tokens
	return
}

func (ast *AST) pushVarsetSelector(
	selectors *[]VarsetSelector,
	last, current int,
	info varsetInfo,
) {
	var selector VarsetSelector
	selector.Expression.Tokens = info.selectorTokens[last:current]
	if last-current == 0 {
		ast.PushErrorToken(info.selectorTokens[current-1], "missing_value")
		return
	}
	if selector.Expression.Tokens[0].Id == lexer.Name && current-last > 1 &&
		selector.Expression.Tokens[1].Id == lexer.Colon {
		selector.NewVariable = true
		selector.Variable.NameToken = selector.Expression.Tokens[0]
		selector.Variable.Name = selector.Variable.NameToken.Kind
		selector.Variable.SetterToken = info.setter
		if current-last > 2 {
			selector.Variable.Type, _ = ast.BuildDataType(
				selector.Expression.Tokens[2:],
				new(int),
				false,
			)
		}
	} else {
		if selector.Expression.Tokens[0].Id == lexer.Name {
			selector.Variable.NameToken = selector.Expression.Tokens[0]
			selector.Variable.Name = selector.Variable.NameToken.Kind
		}
		selector.Expression = ast.BuildExpression(selector.Expression.Tokens)
	}
	*selectors = append(*selectors, selector)
}

func (ast *AST) varsetSelectors(info varsetInfo) []VarsetSelector {
	var selectors []VarsetSelector
	braceCount := 0
	lastIndex := 0
	for index, token := range info.selectorTokens {
		if token.Id == lexer.Brace {
			switch token.Kind {
			case "(", "[", "{":
				braceCount++
			default:
				braceCount--
			}
		}
		if braceCount > 0 {
			continue
		} else if token.Id != lexer.Comma {
			continue
		}
		ast.pushVarsetSelector(&selectors, lastIndex, index, info)
		lastIndex = index + 1
	}
	if lastIndex < len(info.selectorTokens) {
		ast.pushVarsetSelector(&selectors, lastIndex, len(info.selectorTokens), info)
	}
	return selectors
}

func (ast *AST) pushVarsetExpression(exps *[]ExpressionAST, last, current int, info varsetInfo) {
	tokens := info.expressionTokens[last:current]
	if tokens == nil {
		ast.PushErrorToken(info.expressionTokens[current-1], "missing_value")
		return
	}
	*exps = append(*exps, ast.BuildExpression(tokens))
}

func (ast *AST) varsetExpression(info varsetInfo) []ExpressionAST {
	var expressions []ExpressionAST
	braceCount := 0
	lastIndex := 0
	for index, token := range info.expressionTokens {
		if token.Id == lexer.Brace {
			switch token.Kind {
			case "(", "[", "{":
				braceCount++
			default:
				braceCount--
			}
		}
		if braceCount > 0 {
			continue
		} else if token.Id != lexer.Comma {
			continue
		}
		ast.pushVarsetExpression(&expressions, lastIndex, index, info)
		lastIndex = index + 1
	}
	if lastIndex < len(info.expressionTokens) {
		ast.pushVarsetExpression(&expressions, lastIndex, len(info.expressionTokens), info)
	}
	return expressions
}

func (ast *AST) BuildVariableSetStatement(tokens []lexer.Token) (s StatementAST, _ bool) {
	if !checkVariableSetStatementTokens(tokens) {
		return
	}
	info := ast.variableSetInfo(tokens)
	if !info.ok {
		return
	}
	var varAST VariableSetAST
	varAST.Setter = info.setter
	varAST.JustDeclare = info.JustDeclare
	varAST.SelectExpressions = ast.varsetSelectors(info)
	if !info.JustDeclare {
		varAST.ValueExpressions = ast.varsetExpression(info)
	}
	s.Token = tokens[0]
	s.Value = varAST
	return s, true
}

func (ast *AST) BuildNameStatement(tokens []lexer.Token) (s StatementAST) {
	if len(tokens) == 1 {
		ast.PushErrorToken(tokens[0], "invalid_syntax")
		return
	}
	switch tokens[1].Id {
	case lexer.Colon:
		return ast.BuildVariableStatement(tokens)
	case lexer.Brace:
		switch tokens[1].Kind {
		case "(":
			return ast.BuildFunctionCallStatement(tokens)
		}
	}
	ast.PushErrorToken(tokens[0], "invalid_syntax")
	return
}

func (ast *AST) BuildFunctionCallStatement(tokens []lexer.Token) StatementAST {
	return ast.BuildExpressionStatement(tokens)
}

func (ast *AST) BuildExpressionStatement(tokens []lexer.Token) StatementAST {
	return StatementAST{
		Token: tokens[0],
		Value: BlockExpressionAST{
			Expression: ast.BuildExpression(tokens),
		},
	}
}

func (ast *AST) BuildArgs(tokens []lexer.Token) []ArgAST {
	var args []ArgAST
	last := 0
	braceCount := 0
	for index, token := range tokens {
		if token.Id == lexer.Brace {
			switch token.Kind {
			case "{", "[", "(":
				braceCount++
			default:
				braceCount--
			}
		}
		if braceCount > 0 || token.Id != lexer.Comma {
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

func (ast *AST) BuildVariableStatement(tokens []lexer.Token) (s StatementAST) {
	var varAST VariableAST
	position := 0
	if tokens[position].Id != lexer.Name {
		varAST.DefineToken = tokens[position]
		position++
	}
	varAST.NameToken = tokens[position]
	if varAST.NameToken.Id != lexer.Name {
		ast.PushErrorToken(varAST.NameToken, "invalid_syntax")
	}
	varAST.Name = varAST.NameToken.Kind
	varAST.Type = DataTypeAST{Code: jn.Void}

	position++
	if varAST.DefineToken.File != nil {
		if tokens[position].Id != lexer.Colon {
			ast.PushErrorToken(tokens[position], "invalid_syntax")
			return
		}
		position++
	} else {
		position++
	}
	if position < len(tokens) {
		token := tokens[position]
		t, ok := ast.BuildDataType(tokens, &position, false)
		if ok {
			varAST.Type = t
			position++
			if position >= len(tokens) {
				goto ret
			}
			token = tokens[position]
		}
		if token.Id == lexer.Operator {
			if token.Kind != "=" {
				ast.PushErrorToken(token, "invalid_syntax")
				return
			}
			valueTokens := tokens[position+1:]
			if len(valueTokens) == 0 {
				ast.PushErrorToken(token, "missing_value")
				return
			}
			varAST.Value = ast.BuildExpression(valueTokens)
			varAST.SetterToken = token
		}
	}
ret:
	return StatementAST{
		Token: varAST.NameToken,
		Value: varAST,
	}
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
	for index := 0; index < len(tokens); index++ {
		token := tokens[index]
		switch token.Id {
		case lexer.Operator:
			if !operator {
				if IsSingleOperator(token.Kind) && !singleOperatored {
					part = append(part, token)
					singleOperatored = true
					continue
				}
				if braceCount == 0 {
					ast.PushErrorToken(token, "operator_overflow")
				}
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
			switch token.Kind {
			case "(", "[", "{":
				if token.Kind == "[" {
					oldIndex := index
					if _, ok := ast.BuildDataType(tokens, &index, false); ok {
						part = append(part, tokens[oldIndex:index+1]...)
						continue
					}
					index = oldIndex
				}
				singleOperatored = false
				braceCount++
			default:
				braceCount--
			}
		}
		if index > 0 && braceCount == 0 {
			lt := tokens[index-1]
			if (lt.Id == lexer.Name || lt.Id == lexer.Value) &&
				(token.Id == lexer.Name || token.Id == lexer.Value) {
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
	switch token.Id {
	case lexer.Comma:
		return false
	case lexer.Brace:
		if token.Kind == "(" || token.Kind == "{" {
			return false
		}
	}
	return index < tokensLen-1
}

func (ast *AST) checkExpressionToken(token lexer.Token) {
	if token.Kind[0] >= '0' && token.Kind[0] <= '9' {
		var result bool
		if strings.IndexByte(token.Kind, '.') != -1 {
			_, result = new(big.Float).SetString(token.Kind)
		} else {
			result = jnbits.CheckBitInt(token.Kind, 64)
		}
		if !result {
			ast.PushErrorToken(token, "invalid_numeric_range")
		}
	}
}

func (ast *AST) getRange(open, close string) []lexer.Token {
	token := ast.Tokens[ast.Position]
	if token.Id == lexer.Brace && token.Kind == open {
		ast.Position++
		braceCount := 1
		start := ast.Position
		for ; braceCount > 0 && !ast.Ended(); ast.Position++ {
			token := ast.Tokens[ast.Position]
			if token.Id != lexer.Brace {
				continue
			}
			if token.Kind == open {
				braceCount++
			} else if token.Kind == close {
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

func (ast *AST) skipStatement() []lexer.Token {
	start := ast.Position
	ast.Position = nextStatementPos(ast.Tokens, start)
	return ast.Tokens[start:ast.Position]
}
