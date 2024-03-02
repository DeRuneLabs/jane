package ast

import (
	"fmt"
	"math/big"
	"strings"
	"sync"

	"github.com/De-Rune/jane/lexer"
	"github.com/De-Rune/jane/package/jn"
	"github.com/De-Rune/jane/package/jnbits"
)

type Builder struct {
	wg sync.WaitGroup

	Tree     []Object
	Errors   []string
	Tokens   []lexer.Token
	Position int
}

func New(tokens []lexer.Token) *Builder {
	ast := new(Builder)
	ast.Tokens = tokens
	ast.Position = 0
	return ast
}

func NewBuilder(tokens []lexer.Token) *Builder {
	ast := new(Builder)
	ast.Tokens = tokens
	ast.Position = 0
	return ast
}

func (b *Builder) PushError(token lexer.Token, err string) {
	message := jn.Errors[err]
	b.Errors = append(
		b.Errors,
		fmt.Sprintf("%s:%d:%d %s", token.File.Path, token.Row, token.Column, message),
	)
}

func (ast *Builder) Ended() bool {
	return ast.Position >= len(ast.Tokens)
}

func (b *Builder) Build() {
	for b.Position != -1 && !b.Ended() {
		tokens := b.skipStatement()
		token := tokens[0]
		switch token.Id {
		case lexer.At:
			b.Attribute(tokens)
		case lexer.Name:
			b.Name(tokens)
		case lexer.Const:
			b.GlobalVariable(tokens)
		case lexer.Type:
			b.Type(tokens)
		default:
			b.PushError(token, "invalid_syntax")
		}
	}
	b.wg.Wait()
}

func (b *Builder) Type(tokens []lexer.Token) {
	position := 1
	if position >= len(tokens) {
		b.PushError(tokens[position-1], "invalid_syntax")
		return
	}
	token := tokens[position]
	if token.Id != lexer.Name {
		b.PushError(token, "invalid_syntax")
	}
	position++
	if position >= len(tokens) {
		b.PushError(tokens[position-1], "invalid_syntax")
		return
	}
	destinationType, _ := b.DataType(tokens[position:], new(int), true)
	b.Tree = append(b.Tree, Object{
		Token: tokens[1],
		Value: TypeAST{
			Token: tokens[1],
			Name:  tokens[1].Kind,
			Type:  destinationType,
		},
	})
}

func (b *Builder) Name(tokens []lexer.Token) {
	if len(tokens) == 1 {
		b.PushError(tokens[0], "invalid_syntax")
		return
	}
	token := tokens[1]
	switch token.Id {
	case lexer.Colon:
		b.GlobalVariable(tokens)
		return
	case lexer.Brace:
		switch token.Kind {
		case "(":
			funAST := b.Function(tokens, false)
			statement := StatementAST{funAST.Token, funAST, false}
			b.Tree = append(b.Tree, Object{funAST.Token, statement})
			return
		}
	}
	b.PushError(token, "invalid_syntax")
}

func (b *Builder) Attribute(tokens []lexer.Token) {
	var attribute AttributeAST
	index := 0
	attribute.Token = tokens[index]
	index++
	if b.Ended() {
		b.PushError(tokens[index-1], "invalid_syntax")
		return
	}
	attribute.Tag = tokens[index]
	if attribute.Tag.Id != lexer.Name || attribute.Token.Column+1 != attribute.Tag.Column {
		b.PushError(attribute.Tag, "invalid_syntax")
		return
	}
	b.Tree = append(b.Tree, Object{attribute.Token, attribute})
}

func (b *Builder) Function(tokens []lexer.Token, anonymous bool) (funAST FunctionAST) {
	funAST.Token = tokens[0]
	index := 0
	if anonymous {
		funAST.Name = "anonymous"
	} else {
		if funAST.Token.Id != lexer.Name {
			b.PushError(funAST.Token, "invalid_syntax")
		}
		funAST.Name = funAST.Token.Kind
		index++
	}
	funAST.ReturnType.Code = jn.Void
	paramTokens := getRange(&index, "(", ")", tokens)
	if len(paramTokens) > 0 {
		b.Parameters(&funAST, paramTokens)
	}
	if index >= len(tokens) {
		b.PushError(funAST.Token, "body_not_exist")
		return
	}
	token := tokens[index]
	t, ok := b.FunctionReturnDataType(tokens, &index)
	if ok {
		funAST.ReturnType = t
		index++
		if index >= len(tokens) {
			b.PushError(funAST.Token, "body_not_exist")
			return
		}
		token = tokens[index]
	}
	if token.Id != lexer.Brace || token.Kind != "{" {
		b.PushError(token, "invalid_syntax")
		return
	}
	blockTokens := getRange(&index, "{", "}", tokens)
	if blockTokens == nil {
		b.PushError(funAST.Token, "body_not_exist")
		return
	}
	if index < len(tokens) {
		b.PushError(tokens[index], "invalid_syntax")
	}
	funAST.Block = b.Block(blockTokens)
	return
}

func (b *Builder) GlobalVariable(tokens []lexer.Token) {
	if tokens == nil {
		return
	}
	statement := b.VariableStatement(tokens)
	b.Tree = append(b.Tree, Object{statement.Token, statement})
}

func (b *Builder) Parameters(fn *FunctionAST, tokens []lexer.Token) {
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
		b.pushParameter(fn, tokens[last:index], token)
		last = index + 1
	}
	if last < len(tokens) {
		if last == 0 {
			b.pushParameter(fn, tokens[last:], tokens[last])
		} else {
			b.pushParameter(fn, tokens[last:], tokens[last-1])
		}
	}
	b.wg.Add(1)
	go b.checkParamsAsync(fn)
}

func (b *Builder) checkParamsAsync(fn *FunctionAST) {
	defer func() { b.wg.Done() }()
	for _, param := range fn.Params {
		if param.Type.Token.Id == lexer.NA {
			b.PushError(param.Token, "missing_type")
		}
	}
}

func (b *Builder) pushParameter(fn *FunctionAST, tokens []lexer.Token, err lexer.Token) {
	if len(tokens) == 0 {
		b.PushError(err, "invalid_syntax")
		return
	}
	paramAST := ParameterAST{
		Token: tokens[0],
	}
	for index, token := range tokens {
		switch token.Id {
		case lexer.Const:
			if paramAST.Const {
				b.PushError(token, "already_constant")
				continue
			}
			paramAST.Const = true
		case lexer.Operator:
			if token.Kind != "..." {
				b.PushError(token, "invalid_syntax")
				continue
			}
			if paramAST.Variadic {
				b.PushError(token, "already_variadic")
				continue
			}
			paramAST.Variadic = true
		case lexer.Name:
			tokens = tokens[index:]
			if !jn.IsIgnoreName(token.Kind) {
				for _, param := range fn.Params {
					if param.Name == token.Kind {
						b.PushError(token, "parameter_exist")
						break
					}
				}
				paramAST.Name = token.Kind
			}
			if len(tokens) > 1 {
				index := 1
				paramAST.Type, _ = b.DataType(tokens, &index, true)
				index++
				if index < len(tokens) {
					b.PushError(tokens[index], "invalid_syntax")
				}
				index = len(fn.Params) - 1
				for ; index >= 0; index-- {
					param := &fn.Params[index]
					if param.Type.Token.Id != lexer.NA {
						break
					}
					param.Type = paramAST.Type
				}
			}
			goto end
		default:
			if t, ok := b.DataType(tokens, &index, true); ok {
				if index+1 == len(tokens) {
					paramAST.Type = t
					goto end
				}
			}
			b.PushError(token, "invalid_syntax")
			goto end
		}
	}
end:
	fn.Params = append(fn.Params, paramAST)
}

func (b *Builder) DataType(tokens []lexer.Token, index *int, err bool) (dt DataTypeAST, ok bool) {
	first := *index
	for ; *index < len(tokens); *index++ {
		token := tokens[*index]
		switch token.Id {
		case lexer.DataType:
			dataType(token, &dt)
			return dt, true
		case lexer.Name:
			nameType(token, &dt)
			return dt, true
		case lexer.Operator:
			if token.Kind == "*" {
				dt.Value += token.Kind
				break
			}
			if err {
				b.PushError(token, "invalid_syntax")
			}
			return dt, false
		case lexer.Brace:
			switch token.Kind {
			case "(":
				b.functionDataType(token, tokens, index, &dt)
				return dt, true
			case "[":
				*index++
				if *index > len(tokens) {
					if err {
						b.PushError(token, "invalid_syntax")
					}
					return dt, false
				}
				token = tokens[*index]
				if token.Id != lexer.Brace || token.Kind != "]" {
					if err {
						b.PushError(token, "invalid_syntax")
					}
					return dt, false
				}
				dt.Value += "[]"
				continue
			}
			return dt, false
		default:
			if err {
				b.PushError(token, "invalid_syntax")
			}
			return dt, false
		}
	}
	if err {
		b.PushError(tokens[first], "invalid_type")
	}
	return dt, false
}

func dataType(token lexer.Token, dt *DataTypeAST) {
	dt.Token = token
	dt.Code = jn.TypeFromName(dt.Token.Kind)
	dt.Value += dt.Token.Kind
}

func nameType(token lexer.Token, dt *DataTypeAST) {
	dt.Token = token
	dt.Code = jn.Name
	dt.Value += dt.Token.Kind
}

func (b *Builder) functionDataType(
	token lexer.Token,
	tokens []lexer.Token,
	index *int,
	dt *DataTypeAST,
) {
	dt.Token = token
	dt.Code = jn.Function
	value, fun := b.FunctionDataTypeHead(tokens, index)
	fun.ReturnType, _ = b.FunctionReturnDataType(tokens, index)
	dt.Value += value
	dt.Tag = fun
}

func (b *Builder) FunctionDataTypeHead(tokens []lexer.Token, index *int) (string, FunctionAST) {
	var funAST FunctionAST
	var typeValue strings.Builder
	typeValue.WriteByte('(')
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
			b.Parameters(&funAST, tokens[firstIndex+1:*index])
			*index++
			return typeValue.String(), funAST
		}
	}
	b.PushError(tokens[firstIndex], "invalid_type")
	return "", funAST
}

func (b *Builder) pushTypeToTypes(
	types *[]DataTypeAST,
	tokens []lexer.Token,
	errToken lexer.Token,
) {
	if len(tokens) == 0 {
		b.PushError(errToken, "missing_value")
		return
	}
	currentDt, _ := b.DataType(tokens, new(int), false)
	*types = append(*types, currentDt)
}

func (b *Builder) FunctionReturnDataType(
	tokens []lexer.Token,
	index *int,
) (dt DataTypeAST, ok bool) {
	if *index >= len(tokens) {
		return
	}
	token := tokens[*index]
	// NOTE: multi typeindex
	if token.Id == lexer.Brace && token.Kind == "[" {
		*index++
		if *index >= len(tokens) {
			*index--
			goto end
		}
		if token.Id == lexer.Brace && token.Kind == "]" {
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
				b.pushTypeToTypes(&types, tokens[last:*index], tokens[last-1])
				break
			} else if braceCount > 1 {
				continue
			}
			if token.Id != lexer.Comma {
				continue
			}
			b.pushTypeToTypes(&types, tokens[last:*index], tokens[*index-1])
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
	return b.DataType(tokens, index, false)
}

func IsSingleOperator(operator string) bool {
	return operator == "-" ||
		operator == "+" ||
		operator == "~" ||
		operator == "!" ||
		operator == "*" ||
		operator == "&"
}

func (b *Builder) pushStatementToBlock(bs *blockStatement) {
	if len(bs.tokens) == 0 {
		return
	}
	lastToken := bs.tokens[len(bs.tokens)-1]
	if lastToken.Id == lexer.SemiColon {
		if len(bs.tokens) == 1 {
			return
		}
		bs.tokens = bs.tokens[:len(bs.tokens)-1]
	}
	statement := b.Statement(bs)
	statement.WithTerminator = bs.withTerminator
	bs.block.Statements = append(bs.block.Statements, statement)
}

func IsStatement(current, prev lexer.Token) (ok bool, withTerminator bool) {
	ok = current.Id == lexer.SemiColon || prev.Row < current.Row
	withTerminator = current.Id == lexer.SemiColon
	return
}

func nextStatementPos(tokens []lexer.Token, start int) (int, bool) {
	braceCount := 0
	index := start
	for ; index < len(tokens); index++ {
		var isStatement, withTerminator bool
		token := tokens[index]
		if token.Id == lexer.Brace {
			switch token.Kind {
			case "{", "[", "(":
				braceCount++
				continue
			default:
				braceCount--
				if braceCount == 0 {
					if index+1 < len(tokens) {
						isStatement, withTerminator = IsStatement(tokens[index+1], token)
						if isStatement {
							index++
							goto ret
						}
					}
				}
				continue
			}
		}
		if braceCount != 0 {
			continue
		}
		if index > start {
			isStatement, withTerminator = IsStatement(token, tokens[index-1])
		} else {
			isStatement, withTerminator = IsStatement(token, token)
		}
		if !isStatement {
			continue
		}
	ret:
		if withTerminator {
			index++
		}
		return index, withTerminator
	}
	return index, false
}

type blockStatement struct {
	block          *BlockAST
	blockTokens    *[]lexer.Token
	tokens         []lexer.Token
	nextTokens     []lexer.Token
	withTerminator bool
}

func (b *Builder) Block(tokens []lexer.Token) (block BlockAST) {
	for {
		if b.Position == -1 {
			return
		}
		index, withTerminator := nextStatementPos(tokens, 0)
		statementTokens := tokens[:index]
		bs := new(blockStatement)
		bs.block = &block
		bs.blockTokens = &tokens
		bs.tokens = statementTokens
		bs.withTerminator = withTerminator
		b.pushStatementToBlock(bs)
	next:
		if len(bs.nextTokens) > 0 {
			bs.tokens = bs.nextTokens
			bs.nextTokens = nil
			b.pushStatementToBlock(bs)
			goto next
		}
		if index >= len(tokens) {
			break
		}
		tokens = tokens[index:]
	}
	return
}

func (b *Builder) Statement(bs *blockStatement) (s StatementAST) {
	s, ok := b.VariableSetStatement(bs.tokens)
	if ok {
		return s
	}
	token := bs.tokens[0]
	switch token.Id {
	case lexer.Name:
		return b.NameStatement(bs.tokens)
	case lexer.Const:
		return b.VariableStatement(bs.tokens)
	case lexer.Return:
		return b.ReturnStatement(bs.tokens)
	case lexer.Free:
		return b.FreeStatement(bs.tokens)
	case lexer.Iter:
		return b.IterExpr(bs.tokens)
	case lexer.Break:
		return b.BreakStatement(bs.tokens)
	case lexer.Continue:
		return b.ContinueStatement(bs.tokens)
	case lexer.If:
		return b.IfExpr(bs)
	case lexer.Else:
		return b.ElseBlock(bs)
	case lexer.Operator:
		if token.Kind == "<" {
			return b.ReturnStatement(bs.tokens)
		}
	}
	return b.ExprStatement(bs.tokens)
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
	selectorTokens []lexer.Token
	exprTokens     []lexer.Token
	setter         lexer.Token
	ok             bool
	justDeclare    bool
}

func (b *Builder) variableSetInfo(tokens []lexer.Token) (info varsetInfo) {
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
				b.PushError(token, "invalid_syntax")
				info.ok = false
			}
			info.setter = token
			if index+1 >= len(tokens) {
				b.PushError(token, "missing_value")
				info.ok = false
			} else {
				info.exprTokens = tokens[index+1:]
			}
			return
		}
	}
	info.justDeclare = true
	info.selectorTokens = tokens
	return
}

func (b *Builder) pushVarsetSelector(
	selectors *[]VarsetSelector,
	last, current int,
	info varsetInfo,
) {
	var selector VarsetSelector
	selector.Expr.Tokens = info.selectorTokens[last:current]
	if last-current == 0 {
		b.PushError(info.selectorTokens[current-1], "missing_value")
		return
	}
	if selector.Expr.Tokens[0].Id == lexer.Name && current-last > 1 &&
		selector.Expr.Tokens[1].Id == lexer.Colon {
		selector.NewVariable = true
		selector.Variable.NameToken = selector.Expr.Tokens[0]
		selector.Variable.Name = selector.Variable.NameToken.Kind
		selector.Variable.SetterToken = info.setter
		if current-last > 2 {
			selector.Variable.Type, _ = b.DataType(selector.Expr.Tokens[2:], new(int), false)
		}
	} else {
		if selector.Expr.Tokens[0].Id == lexer.Name {
			selector.Variable.NameToken = selector.Expr.Tokens[0]
			selector.Variable.Name = selector.Variable.NameToken.Kind
		}
		selector.Expr = b.Expr(selector.Expr.Tokens)
	}
	*selectors = append(*selectors, selector)
}

func (b *Builder) varsetSelectors(info varsetInfo) []VarsetSelector {
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
		b.pushVarsetSelector(&selectors, lastIndex, index, info)
		lastIndex = index + 1
	}
	if lastIndex < len(info.selectorTokens) {
		b.pushVarsetSelector(&selectors, lastIndex, len(info.selectorTokens), info)
	}
	return selectors
}

func (b *Builder) pushVarsetExpr(exps *[]ExprAST, last, current int, info varsetInfo) {
	tokens := info.exprTokens[last:current]
	if tokens == nil {
		b.PushError(info.exprTokens[current-1], "missing_value")
		return
	}
	*exps = append(*exps, b.Expr(tokens))
}

func (b *Builder) varsetExprs(info varsetInfo) []ExprAST {
	var exprs []ExprAST
	braceCount := 0
	lastIndex := 0
	for index, token := range info.exprTokens {
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
		b.pushVarsetExpr(&exprs, lastIndex, index, info)
		lastIndex = index + 1
	}
	if lastIndex < len(info.exprTokens) {
		b.pushVarsetExpr(&exprs, lastIndex, len(info.exprTokens), info)
	}
	return exprs
}

func (b *Builder) VariableSetStatement(tokens []lexer.Token) (s StatementAST, _ bool) {
	if !checkVariableSetStatementTokens(tokens) {
		return
	}
	info := b.variableSetInfo(tokens)
	if !info.ok {
		return
	}
	var varAST VariableSetAST
	varAST.Setter = info.setter
	varAST.JustDeclare = info.justDeclare
	varAST.SelectExprs = b.varsetSelectors(info)
	if !info.justDeclare {
		varAST.ValueExprs = b.varsetExprs(info)
	}
	s.Token = tokens[0]
	s.Value = varAST
	return s, true
}

func (b *Builder) NameStatement(tokens []lexer.Token) (s StatementAST) {
	if len(tokens) == 1 {
		b.PushError(tokens[0], "invalid_syntax")
		return
	}
	switch tokens[1].Id {
	case lexer.Colon:
		return b.VariableStatement(tokens)
	case lexer.Brace:
		switch tokens[1].Kind {
		case "(":
			return b.FunctionCallStatement(tokens)
		}
	}
	b.PushError(tokens[0], "invalid_syntax")
	return
}

func (b *Builder) FunctionCallStatement(tokens []lexer.Token) StatementAST {
	return b.ExprStatement(tokens)
}

func (b *Builder) ExprStatement(tokens []lexer.Token) StatementAST {
	block := ExprStatementAST{b.Expr(tokens)}
	return StatementAST{tokens[0], block, false}
}

func (b *Builder) Args(tokens []lexer.Token) []ArgAST {
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
		b.pushArg(&args, tokens[last:index], token)
		last = index + 1
	}
	if last < len(tokens) {
		if last == 0 {
			b.pushArg(&args, tokens[last:], tokens[last])
		} else {
			b.pushArg(&args, tokens[last:], tokens[last-1])
		}
	}
	return args
}

func (b *Builder) pushArg(args *[]ArgAST, tokens []lexer.Token, err lexer.Token) {
	if len(tokens) == 0 {
		b.PushError(err, "invalid_syntax")
		return
	}
	var arg ArgAST
	arg.Token = tokens[0]
	arg.Expr = b.Expr(tokens)
	*args = append(*args, arg)
}

func (b *Builder) VariableStatement(tokens []lexer.Token) (s StatementAST) {
	var varAST VariableAST
	position := 0
	if tokens[position].Id != lexer.Name {
		varAST.DefineToken = tokens[position]
		position++
	}
	varAST.NameToken = tokens[position]
	if varAST.NameToken.Id != lexer.Name {
		b.PushError(varAST.NameToken, "invalid_syntax")
	}
	varAST.Name = varAST.NameToken.Kind
	varAST.Type = DataTypeAST{Code: jn.Void}
	position++
	if varAST.DefineToken.File != nil {
		if tokens[position].Id != lexer.Colon {
			b.PushError(tokens[position], "invalid_syntax")
			return
		}
		position++
	} else {
		position++
	}
	if position < len(tokens) {
		token := tokens[position]
		t, ok := b.DataType(tokens, &position, false)
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
				b.PushError(token, "invalid_syntax")
				return
			}
			valueTokens := tokens[position+1:]
			if len(valueTokens) == 0 {
				b.PushError(token, "missing_value")
				return
			}
			varAST.Value = b.Expr(valueTokens)
			varAST.SetterToken = token
		}
	}
ret:
	return StatementAST{varAST.NameToken, varAST, false}
}

func (b *Builder) ReturnStatement(tokens []lexer.Token) StatementAST {
	var returnModel ReturnAST
	returnModel.Token = tokens[0]
	if len(tokens) > 1 {
		returnModel.Expr = b.Expr(tokens[1:])
	}
	return StatementAST{returnModel.Token, returnModel, false}
}

func (b *Builder) FreeStatement(tokens []lexer.Token) StatementAST {
	var freeAST FreeAST
	freeAST.Token = tokens[0]
	tokens = tokens[1:]
	if len(tokens) == 0 {
		b.PushError(freeAST.Token, "missing_expression")
	} else {
		freeAST.Expr = b.Expr(tokens)
	}
	return StatementAST{freeAST.Token, freeAST, false}
}

func blockExprTokens(tokens []lexer.Token) (expr []lexer.Token) {
	braceCount := 0
	for index, token := range tokens {
		if token.Id == lexer.Brace {
			switch token.Kind {
			case "{":
				if braceCount > 0 {
					braceCount++
					break
				}
				return tokens[:index]
			case "(", "[":
				braceCount++
			default:
				braceCount--
			}
		}
	}
	return nil
}

func (b *Builder) getWhileIterProfile(tokens []lexer.Token) WhileProfile {
	return WhileProfile{b.Expr(tokens)}
}

func (b *Builder) pushVarsTokensPart(
	vars *[][]lexer.Token,
	part []lexer.Token,
	errTok lexer.Token,
) {
	if len(part) == 0 {
		b.PushError(errTok, "missing_value")
	}
	*vars = append(*vars, part)
}

func (b *Builder) getForeachVarsTokens(tokens []lexer.Token) [][]lexer.Token {
	var vars [][]lexer.Token
	braceCount := 0
	last := 0
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
		if token.Id == lexer.Comma {
			part := tokens[last:index]
			b.pushVarsTokensPart(&vars, part, token)
			last = index + 1
		}
	}
	if last < len(tokens) {
		part := tokens[last:]
		b.pushVarsTokensPart(&vars, part, tokens[last])
	}
	return vars
}

func (b *Builder) getForeachIterVars(varsTokens [][]lexer.Token) []VariableAST {
	var vars []VariableAST
	for _, tokens := range varsTokens {
		var vast VariableAST
		vast.NameToken = tokens[0]
		if vast.NameToken.Id != lexer.Name {
			b.PushError(vast.NameToken, "invalid_syntax")
			vars = append(vars, vast)
			continue
		}
		vast.Name = vast.NameToken.Kind
		if len(tokens) == 1 {
			vars = append(vars, vast)
			continue
		}
		if colon := tokens[1]; colon.Id != lexer.Colon {
			b.PushError(colon, "invalid_syntax")
			vars = append(vars, vast)
			continue
		}
		vast.New = true
		index := new(int)
		*index = 2
		if *index >= len(tokens) {
			vars = append(vars, vast)
			continue
		}
		vast.Type, _ = b.DataType(tokens, index, true)
		if *index < len(tokens)-1 {
			b.PushError(tokens[*index], "invalid_syntax")
		}
		vars = append(vars, vast)
	}
	return vars
}

func (b *Builder) getForeachIterProfile(
	varTokens, exprTokens []lexer.Token,
	inTok lexer.Token,
) ForeachProfile {
	var profile ForeachProfile
	profile.InToken = inTok
	profile.Expr = b.Expr(exprTokens)
	if len(varTokens) == 0 {
		profile.KeyA.Name = "__"
		profile.KeyB.Name = "__"
	} else {
		varsTokens := b.getForeachVarsTokens(varTokens)
		if len(varsTokens) == 0 {
			return profile
		}
		if len(varsTokens) > 2 {
			b.PushError(inTok, "much_foreach_vars")
		}
		vars := b.getForeachIterVars(varsTokens)
		profile.KeyA = vars[0]
		if len(vars) > 1 {
			profile.KeyB = vars[1]
		} else {
			profile.KeyB.Name = "__"
		}
	}
	return profile
}

func (b *Builder) getIterProfile(tokens []lexer.Token) IterProfile {
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
		if braceCount != 0 {
			continue
		}
		if token.Id == lexer.In {
			varTokens := tokens[:index]
			exprTokens := tokens[index+1:]
			return b.getForeachIterProfile(varTokens, exprTokens, token)
		}
	}
	return b.getWhileIterProfile(tokens)
}

func (b *Builder) IterExpr(tokens []lexer.Token) (s StatementAST) {
	var iter IterAST
	iter.Token = tokens[0]
	tokens = tokens[1:]
	if len(tokens) == 0 {
		b.PushError(iter.Token, "body_not_exist")
		return
	}
	exprTokens := blockExprTokens(tokens)
	if len(exprTokens) > 0 {
		iter.Profile = b.getIterProfile(exprTokens)
	}
	index := new(int)
	*index = len(exprTokens)
	blockTokens := getRange(index, "{", "}", tokens)
	if blockTokens == nil {
		b.PushError(iter.Token, "body_not_exist")
		return
	}
	if *index < len(tokens) {
		b.PushError(tokens[*index], "invalid_syntax")
	}
	iter.Block = b.Block(blockTokens)
	return StatementAST{iter.Token, iter, false}
}

func (b *Builder) IfExpr(bs *blockStatement) (s StatementAST) {
	var ifast IfAST
	ifast.Token = bs.tokens[0]
	bs.tokens = bs.tokens[1:]
	exprTokens := blockExprTokens(bs.tokens)
	if len(exprTokens) == 0 {
		b.PushError(ifast.Token, "missing_expression")
	}
	index := new(int)
	*index = len(exprTokens)
	blockTokens := getRange(index, "{", "}", bs.tokens)
	if blockTokens == nil {
		b.PushError(ifast.Token, "body_not_exist")
		return
	}
	if *index < len(bs.tokens) {
		if bs.tokens[*index].Id == lexer.Else {
			bs.nextTokens = bs.tokens[*index:]
		} else {
			b.PushError(bs.tokens[*index], "invalid_syntax")
		}
	}
	ifast.Expr = b.Expr(exprTokens)
	ifast.Block = b.Block(blockTokens)
	return StatementAST{ifast.Token, ifast, false}
}

func (b *Builder) ElseIfExpr(bs *blockStatement) (s StatementAST) {
	var elif ElseIfAST
	elif.Token = bs.tokens[1]
	bs.tokens = bs.tokens[2:]
	exprTokens := blockExprTokens(bs.tokens)
	if len(exprTokens) == 0 {
		b.PushError(elif.Token, "missing_expression")
	}
	index := new(int)
	*index = len(exprTokens)
	blockTokens := getRange(index, "{", "}", bs.tokens)
	if blockTokens == nil {
		b.PushError(elif.Token, "body_not_exist")
		return
	}
	if *index < len(bs.tokens) {
		if bs.tokens[*index].Id == lexer.Else {
			bs.nextTokens = bs.tokens[*index:]
		} else {
			b.PushError(bs.tokens[*index], "invalid_syntax")
		}
	}
	elif.Expr = b.Expr(exprTokens)
	elif.Block = b.Block(blockTokens)
	return StatementAST{elif.Token, elif, false}
}

func (b *Builder) ElseBlock(bs *blockStatement) (s StatementAST) {
	if len(bs.tokens) > 1 && bs.tokens[1].Id == lexer.If {
		return b.ElseIfExpr(bs)
	}
	var elseast ElseAST
	elseast.Token = bs.tokens[0]
	bs.tokens = bs.tokens[1:]
	index := new(int)
	blockTokens := getRange(index, "{", "}", bs.tokens)
	if blockTokens == nil {
		if *index < len(bs.tokens) {
			b.PushError(elseast.Token, "else_have_expr")
		} else {
			b.PushError(elseast.Token, "body_not_exist")
		}
		return
	}
	if *index < len(bs.tokens) {
		b.PushError(bs.tokens[*index], "invalid_syntax")
	}
	elseast.Block = b.Block(blockTokens)
	return StatementAST{elseast.Token, elseast, false}
}

func (b *Builder) BreakStatement(tokens []lexer.Token) StatementAST {
	var breakAST BreakAST
	breakAST.Token = tokens[0]
	if len(tokens) > 1 {
		b.PushError(tokens[1], "invalid_syntax")
	}
	return StatementAST{breakAST.Token, breakAST, false}
}

func (b *Builder) ContinueStatement(tokens []lexer.Token) StatementAST {
	var ContinueAST ContinueAST
	ContinueAST.Token = tokens[0]
	if len(tokens) > 1 {
		b.PushError(tokens[1], "invalid_syntax")
	}
	return StatementAST{ContinueAST.Token, ContinueAST, false}
}

func (b *Builder) Expr(tokens []lexer.Token) (e ExprAST) {
	e.Processes = b.getExprProcesses(tokens)
	e.Tokens = tokens
	return
}

func (b *Builder) isOverflowOperator(kind string) bool {
	return kind == "+" ||
		kind == "-" ||
		kind == "*" ||
		kind == "/" ||
		kind == "%" ||
		kind == "&" ||
		kind == "|" ||
		kind == "^" ||
		kind == "<" ||
		kind == ">" ||
		kind == "~" ||
		kind == "!"
}

func (b *Builder) getExprProcesses(tokens []lexer.Token) [][]lexer.Token {
	var processes [][]lexer.Token
	var part []lexer.Token
	operator := false
	value := false
	braceCount := 0
	pushedError := false
	singleOperatored := false
	newKeyword := false
	for index := 0; index < len(tokens); index++ {
		token := tokens[index]
		switch token.Id {
		case lexer.Operator:
			if newKeyword || !b.isOverflowOperator(token.Kind) {
				part = append(part, token)
				continue
			}
			if !operator {
				if IsSingleOperator(token.Kind) && !singleOperatored {
					part = append(part, token)
					singleOperatored = true
					continue
				}
				if braceCount == 0 {
					b.PushError(token, "operator_overflow")
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
					_, ok := b.DataType(tokens, &index, false)
					if ok {
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
		case lexer.New:
			newKeyword = true
		case lexer.Name:
			if braceCount == 0 {
				newKeyword = false
			}
		}
		if index > 0 && braceCount == 0 {
			lt := tokens[index-1]
			if (lt.Id == lexer.Name || lt.Id == lexer.Value) &&
				(token.Id == lexer.Name || token.Id == lexer.Value) {
				b.PushError(token, "invalid_syntax")
				pushedError = true
			}
		}
		b.checkExprToken(token)
		part = append(part, token)
		operator = requireOperatorForProcess(token, index, len(tokens))
		value = false
	}
	if len(part) > 0 {
		processes = append(processes, part)
	}
	if value {
		b.PushError(processes[len(processes)-1][0], "operator_overflow")
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

func (b *Builder) checkExprToken(token lexer.Token) {
	if token.Kind[0] >= '0' && token.Kind[0] <= '9' {
		var result bool
		if strings.IndexByte(token.Kind, '.') != -1 {
			_, result = new(big.Float).SetString(token.Kind)
		} else {
			result = jnbits.CheckBitInt(token.Kind, 64)
		}
		if !result {
			b.PushError(token, "invalid_numeric_range")
		}
	}
}

func getRange(index *int, open, close string, tokens []lexer.Token) []lexer.Token {
	if *index >= len(tokens) {
		return nil
	}
	token := tokens[*index]
	if token.Id == lexer.Brace && token.Kind == open {
		*index++
		braceCount := 1
		start := *index
		for ; braceCount > 0 && *index < len(tokens); *index++ {
			token := tokens[*index]
			if token.Id != lexer.Brace {
				continue
			}
			if token.Kind == open {
				braceCount++
			} else if token.Kind == close {
				braceCount--
			}
		}
		return tokens[start : *index-1]
	}
	return nil
}

func (b *Builder) skipStatement() []lexer.Token {
	start := b.Position
	b.Position, _ = nextStatementPos(b.Tokens, start)
	tokens := b.Tokens[start:b.Position]
	if tokens[len(tokens)-1].Id == lexer.SemiColon {
		if len(tokens) == 1 {
			return b.skipStatement()
		}
		tokens = tokens[:len(tokens)-1]
	}
	return tokens
}
