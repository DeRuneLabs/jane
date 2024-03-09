package ast

import (
	"os"
	"strings"
	"sync"

	"github.com/DeRuneLabs/jane/lexer"
	"github.com/DeRuneLabs/jane/package/jn"
	"github.com/DeRuneLabs/jane/package/jnapi"
	"github.com/DeRuneLabs/jane/package/jnbits"
	"github.com/DeRuneLabs/jane/package/jnlog"
)

type Builder struct {
	wg  sync.WaitGroup
	pub bool

	Tree   []Obj
	Errs   []jnlog.CompilerLog
	Tokens []lexer.Token
	Pos    int
}

func NewBuilder(tokens []lexer.Token) *Builder {
	b := new(Builder)
	b.Tokens = tokens
	b.Pos = 0
	return b
}

func (b *Builder) pusherr(token lexer.Token, key string, args ...interface{}) {
	b.Errs = append(b.Errs, jnlog.CompilerLog{
		Type:   jnlog.Err,
		Row:    token.Row,
		Column: token.Column,
		Path:   token.File.Path,
		Msg:    jn.GetErr(key, args...),
	})
}

func (ast *Builder) Ended() bool {
	return ast.Pos >= len(ast.Tokens)
}

func (b *Builder) buildNode(tokens []lexer.Token) {
	tok := tokens[0]
	switch tok.Id {
	case lexer.Use:
		b.Use(tokens)
	case lexer.At:
		b.Attribute(tokens)
	case lexer.Id:
		b.Id(tokens)
	case lexer.Const, lexer.Volatile:
		b.GlobalVar(tokens)
	case lexer.Type:
		t := b.Type(tokens)
		t.Pub = b.pub
		b.pub = false
		b.Tree = append(b.Tree, Obj{t.Token, t})
	case lexer.Comment:
		b.Comment(tokens[0])
	case lexer.Preprocessor:
		b.Preprocessor(tokens)
	default:
		b.pusherr(tok, "invalid_syntax")
		return
	}
	if b.pub {
		b.pusherr(tok, "def_not_support_pub")
	}
}

func (b *Builder) Build() {
	for b.Pos != -1 && !b.Ended() {
		toks := b.skipStatement()
		b.pub = toks[0].Id == lexer.Pub
		if b.pub {
			if len(toks) == 1 {
				if b.Ended() {
					b.pusherr(toks[0], "invalid_syntax")
					continue
				}
				toks = b.skipStatement()
			} else {
				toks = toks[1:]
			}
		}
		b.buildNode(toks)
	}
	b.wg.Wait()
}

func (b *Builder) Type(tokens []lexer.Token) (t Type) {
	i := 1
	if i >= len(tokens) {
		b.pusherr(tokens[i-1], "invalid_syntax")
		return
	}
	tok := tokens[i]
	if tok.Id != lexer.Id {
		b.pusherr(tok, "invalid_syntax")
	}
	i++
	if i >= len(tokens) {
		b.pusherr(tokens[i-1], "invalid_syntax")
		return
	}
	destType, _ := b.DataType(tokens[i:], new(int), true)
	tok = tokens[1]
	return Type{
		Pub:   b.pub,
		Token: tok,
		Id:    tok.Kind,
		Type:  destType,
	}
}

func (b *Builder) Comment(token lexer.Token) {
	token.Kind = strings.TrimSpace(token.Kind[2:])
	if strings.HasPrefix(token.Kind, "cxx:") {
		b.Tree = append(b.Tree, Obj{token, CxxEmbed{token.Kind[4:]}})
	} else {
		b.Tree = append(b.Tree, Obj{token, Comment{token.Kind}})
	}
}

func (b *Builder) Preprocessor(tokens []lexer.Token) {
	if len(tokens) == 1 {
		b.pusherr(tokens[0], "invalid_syntax")
		return
	}
	var pp Preprocessor
	tokens = tokens[1:]
	tok := tokens[0]
	if tok.Id != lexer.Id {
		b.pusherr(pp.Token, "invalid_syntax")
		return
	}
	ok := false
	switch tok.Kind {
	case "pragma":
		ok = b.Pragma(&pp, tokens)
	default:
		b.pusherr(tok, "invalid_preprocessor")
		return
	}
	if ok {
		b.Tree = append(b.Tree, Obj{pp.Token, pp})
	}
}

func (b *Builder) Pragma(pp *Preprocessor, tokens []lexer.Token) bool {
	if len(tokens) == 1 {
		b.pusherr(tokens[0], "missing_pragma_directive")
		return false
	}
	tokens = tokens[1:]
	tok := tokens[0]
	if tok.Id != lexer.Id {
		b.pusherr(tok, "invalid_syntax")
		return false
	}
	var d Directive
	ok := false
	switch tok.Kind {
	case "enofi":
		ok = b.pragmaEnofi(&d, tokens)
	default:
		b.pusherr(tok, "invalid_pragma_directive")
	}
	pp.Command = d
	return ok
}

func (b *Builder) pragmaEnofi(d *Directive, tokens []lexer.Token) bool {
	if len(tokens) > 1 {
		b.pusherr(tokens[1], "invalid_syntax")
		return false
	}
	d.Command = EnofiDirective{}
	return true
}

func (b *Builder) Id(tokens []lexer.Token) {
	if len(tokens) == 1 {
		b.pusherr(tokens[0], "invalid_syntax")
		return
	}
	tok := tokens[1]
	switch tok.Id {
	case lexer.Colon:
		b.GlobalVar(tokens)
		return
	case lexer.Brace:
		switch tok.Kind {
		case "(":
			f := b.Func(tokens, false)
			s := Statement{f.Token, f, false}
			b.Tree = append(b.Tree, Obj{f.Token, s})
			return
		}
	}
	b.pusherr(tok, "invalid_syntax")
}

func (b *Builder) Use(tokens []lexer.Token) {
	var use Use
	use.Token = tokens[0]
	if len(tokens) < 2 {
		b.pusherr(use.Token, "missing_use_path")
		return
	}
	use.Path = b.usePath(tokens[1:])
	b.Tree = append(b.Tree, Obj{use.Token, use})
}

func (b *Builder) usePath(tokens []lexer.Token) string {
	var path strings.Builder
	path.WriteString(jn.StdlibPath)
	path.WriteRune(os.PathSeparator)
	for i, tok := range tokens {
		if i%2 != 0 {
			if tok.Id != lexer.Dot {
				b.pusherr(tok, "invalid_syntax")
			}
			path.WriteRune(os.PathSeparator)
			continue
		}
		if tok.Id != lexer.Id {
			b.pusherr(tok, "invalid_syntax")
		}
		path.WriteString(tok.Kind)
	}
	return path.String()
}

func (b *Builder) Attribute(tokens []lexer.Token) {
	var a Attribute
	i := 0
	a.Token = tokens[i]
	i++
	if b.Ended() {
		b.pusherr(tokens[i-1], "invalid_syntax")
		return
	}
	a.Tag = tokens[i]
	if a.Tag.Id != lexer.Id ||
		a.Token.Column+1 != a.Tag.Column {
		b.pusherr(a.Tag, "invalid_syntax")
		return
	}
	b.Tree = append(b.Tree, Obj{a.Token, a})
}

func (b *Builder) Func(tokens []lexer.Token, anonymous bool) (f Func) {
	f.Token = tokens[0]
	i := 0
	f.Pub = b.pub
	b.pub = false
	if anonymous {
		f.Id = "anonymous"
	} else {
		if f.Token.Id != lexer.Id {
			b.pusherr(f.Token, "invalid_syntax")
		}
		f.Id = f.Token.Kind
		i++
	}
	f.RetType.Id = jn.Void
	paramToks := b.getrange(&i, "(", ")", &tokens)
	if len(paramToks) > 0 {
		b.Params(&f, paramToks)
	}
	if i >= len(tokens) {
		b.pusherr(f.Token, "body_not_exist")
		return
	}
	tok := tokens[i]
	t, ok := b.FuncRetDataType(tokens, &i)
	if ok {
		f.RetType = t
		i++
		if i >= len(tokens) {
			b.pusherr(f.Token, "body_not_exist")
			return
		}
		tok = tokens[i]
	}
	if tok.Id != lexer.Brace || tok.Kind != "{" {
		b.pusherr(tok, "invalid_syntax")
		return
	}
	blockToks := b.getrange(&i, "{", "}", &tokens)
	if blockToks == nil {
		b.pusherr(f.Token, "body_not_exist")
		return
	}
	if i < len(tokens) {
		b.pusherr(tokens[i], "invalid_syntax")
	}
	f.Block = b.Block(blockToks)
	return
}

func (b *Builder) GlobalVar(tokens []lexer.Token) {
	if tokens == nil {
		return
	}
	s := b.VarStatement(tokens)
	b.Tree = append(b.Tree, Obj{s.Token, s})
}

func (b *Builder) Params(fn *Func, tokens []lexer.Token) {
	last := 0
	braceCount := 0
	for i, tok := range tokens {
		if tok.Id == lexer.Brace {
			switch tok.Kind {
			case "{", "[", "(":
				braceCount++
			default:
				braceCount--
			}
		}
		if braceCount > 0 || tok.Id != lexer.Comma {
			continue
		}
		b.pushParam(fn, tokens[last:i], tok)
		last = i + 1
	}
	if last < len(tokens) {
		if last == 0 {
			b.pushParam(fn, tokens[last:], tokens[last])
		} else {
			b.pushParam(fn, tokens[last:], tokens[last-1])
		}
	}
	b.wg.Add(1)
	go b.checkParamsAsync(fn)
}

func (b *Builder) checkParamsAsync(f *Func) {
	defer func() { b.wg.Done() }()
	for _, p := range f.Params {
		if p.Type.Token.Id == lexer.NA {
			b.pusherr(p.Token, "missing_type")
		}
	}
}

func (b *Builder) pushParam(f *Func, tokens []lexer.Token, errtok lexer.Token) {
	if len(tokens) == 0 {
		b.pusherr(errtok, "invalid_syntax")
		return
	}
	past := Parameter{Token: tokens[0]}
	for i, tok := range tokens {
		switch tok.Id {
		case lexer.Const:
			if past.Const {
				b.pusherr(tok, "already_constant")
				continue
			}
			past.Const = true
		case lexer.Volatile:
			if past.Volatile {
				b.pusherr(tok, "already_volatile")
				continue
			}
			past.Volatile = true
		case lexer.Operator:
			if tok.Kind != "..." {
				b.pusherr(tok, "invalid_syntax")
				continue
			}
			if past.Variadic {
				b.pusherr(tok, "already_variadic")
				continue
			}
			past.Variadic = true
		case lexer.Id:
			tokens = tokens[i:]
			if !jnapi.IsIgnoreId(tok.Kind) {
				for _, param := range f.Params {
					if param.Id == tok.Kind {
						b.pusherr(tok, "parameter_exist", tok.Kind)
						break
					}
				}
				past.Id = tok.Kind
			}
			if len(tokens) > 1 {
				i := 1
				past.Type, _ = b.DataType(tokens, &i, true)
				i++
				if i < len(tokens) {
					b.pusherr(tokens[i], "invalid_syntax")
				}
				i = len(f.Params) - 1
				for ; i >= 0; i-- {
					param := &f.Params[i]
					if param.Type.Token.Id != lexer.NA {
						break
					}
					param.Type = past.Type
				}
			}
			goto end
		default:
			if t, ok := b.DataType(tokens, &i, true); ok {
				if i+1 == len(tokens) {
					past.Type = t
					goto end
				}
			}
			b.pusherr(tok, "invalid_syntax")
			goto end
		}
	}
end:
	f.Params = append(f.Params, past)
}

func (b *Builder) DataType(tokens []lexer.Token, i *int, err bool) (t DataType, ok bool) {
	first := *i
	var dtv strings.Builder
	for ; *i < len(tokens); *i++ {
		tok := tokens[*i]
		switch tok.Id {
		case lexer.DataType:
			t.Token = tok
			t.Id = jn.TypeFromId(t.Token.Kind)
			dtv.WriteString(t.Token.Kind)
			ok = true
			goto ret
		case lexer.Id:
			t.Token = tok
			t.Id = jn.Id
			dtv.WriteString(t.Token.Kind)
			ok = true
			goto ret
		case lexer.Operator:
			if tok.Kind == "*" {
				dtv.WriteString(tok.Kind)
				break
			}
			if err {
				b.pusherr(tok, "invalid_syntax")
			}
			return
		case lexer.Brace:
			switch tok.Kind {
			case "(":
				t.Token = tok
				t.Id = jn.Func
				val, f := b.FuncDataTypeHead(tokens, i)
				f.RetType, _ = b.FuncRetDataType(tokens, i)
				dtv.WriteString(val)
				t.Tag = f
				ok = true
				goto ret
			case "[":
				*i++
				if *i > len(tokens) {
					if err {
						b.pusherr(tok, "invalid_syntax")
					}
					return
				}
				tok = tokens[*i]
				if tok.Id == lexer.Brace && tok.Kind == "]" {
					dtv.WriteString("[]")
					continue
				}
				*i--
				dt, val := b.MapDataType(tokens, i, err)
				if val == "" {
					if err {
						b.pusherr(tok, "invalid_syntax")
					}
					return
				}
				t = dt
				dtv.WriteString(val)
				ok = true
				goto ret
			}
			return
		default:
			if err {
				b.pusherr(tok, "invalid_syntax")
			}
			return
		}
	}
	if err {
		b.pusherr(tokens[first], "invalid_type")
	}
ret:
	t.Val = dtv.String()
	return
}

func (b *Builder) MapDataType(tokens []lexer.Token, i *int, err bool) (t DataType, _ string) {
	t.Id = jn.Map
	t.Token = tokens[0]
	braceCount := 0
	colon := -1
	start := *i
	var mapToks []lexer.Token
	for ; *i < len(tokens); *i++ {
		tok := tokens[*i]
		if tok.Id == lexer.Brace {
			switch tok.Kind {
			case "(", "[", "{":
				braceCount++
			default:
				braceCount--
			}
		}
		if braceCount == 0 {
			if start+1 > *i {
				return
			}
			mapToks = tokens[start+1 : *i]
			break
		} else if braceCount != 1 {
			continue
		}
		if colon == -1 && tok.Id == lexer.Colon {
			colon = *i - start - 1
		}
	}
	if mapToks == nil || colon == -1 {
		return
	}
	colonTok := tokens[colon]
	if colon == 0 || colon+1 >= len(mapToks) {
		b.pusherr(colonTok, "missing_expr")
		return t, " "
	}
	keyTypeToks := mapToks[:colon]
	valTypeToks := mapToks[colon+1:]
	types := make([]DataType, 2)
	j := 0
	types[0], _ = b.DataType(keyTypeToks, &j, err)
	if j < len(keyTypeToks) && err {
		b.pusherr(keyTypeToks[j], "invalid_syntax")
	}
	j = 0
	types[1], _ = b.DataType(valTypeToks, &j, err)
	if j < len(valTypeToks) && err {
		b.pusherr(valTypeToks[j], "invalid_syntax")
	}
	t.Tag = types
	var val strings.Builder
	val.WriteByte('[')
	val.WriteString(types[0].Val)
	val.WriteByte(':')
	val.WriteString(types[1].Val)
	val.WriteByte(']')
	return t, val.String()
}

func (b *Builder) FuncDataTypeHead(tokens []lexer.Token, i *int) (string, Func) {
	var f Func
	var typeVal strings.Builder
	typeVal.WriteByte('(')
	brace := 1
	firstIndex := *i
	for *i++; *i < len(tokens); *i++ {
		tok := tokens[*i]
		typeVal.WriteString(tok.Kind)
		switch tok.Id {
		case lexer.Brace:
			switch tok.Kind {
			case "{", "[", "(":
				brace++
			default:
				brace--
			}
		}
		if brace == 0 {
			b.Params(&f, tokens[firstIndex+1:*i])
			*i++
			return typeVal.String(), f
		}
	}
	b.pusherr(tokens[firstIndex], "invalid_type")
	return "", f
}

func (b *Builder) pushTypeToTypes(types *[]DataType, tokens []lexer.Token, errTok lexer.Token) {
	if len(tokens) == 0 {
		b.pusherr(errTok, "missing_expr")
		return
	}
	currentDt, _ := b.DataType(tokens, new(int), false)
	*types = append(*types, currentDt)
}

func (b *Builder) FuncRetDataType(tokens []lexer.Token, i *int) (t DataType, ok bool) {
	if *i >= len(tokens) {
		return
	}
	tok := tokens[*i]
	start := *i
	if tok.Id == lexer.Brace && tok.Kind == "[" {
		t.Val += tok.Kind
		*i++
		if *i >= len(tokens) {
			*i--
			goto end
		}
		tok = tokens[*i]
		if tok.Id == lexer.Brace && tok.Kind == "]" {
			*i--
			goto end
		}
		var types []DataType
		braceCount := 1
		last := *i
		for ; *i < len(tokens); *i++ {
			tok := tokens[*i]
			t.Val += tok.Kind
			if tok.Id == lexer.Brace {
				switch tok.Kind {
				case "(", "[", "{":
					braceCount++
				default:
					braceCount--
				}
			}
			if braceCount == 0 {
				if tok.Id == lexer.Colon {
					*i = start
					goto end
				}
				b.pushTypeToTypes(&types, tokens[last:*i], tokens[last-1])
				break
			} else if braceCount > 1 {
				continue
			}
			if tok.Id != lexer.Comma {
				continue
			}
			b.pushTypeToTypes(&types, tokens[last:*i], tokens[*i-1])
			last = *i + 1
		}
		if len(types) > 1 {
			t.MultiTyped = true
			t.Tag = types
		} else {
			t = types[0]
		}
		ok = true
		return
	}
end:
	return b.DataType(tokens, i, false)
}

func IsSingleOperator(kind string) bool {
	return kind == "-" ||
		kind == "+" ||
		kind == "~" ||
		kind == "!" ||
		kind == "*" ||
		kind == "&"
}

func (b *Builder) pushStatementToBlock(bs *blockStatement) {
	if len(bs.tokens) == 0 {
		return
	}
	lastTok := bs.tokens[len(bs.tokens)-1]
	if lastTok.Id == lexer.SemiColon {
		if len(bs.tokens) == 1 {
			return
		}
		bs.tokens = bs.tokens[:len(bs.tokens)-1]
	}
	s := b.Statement(bs)
	if s.Val == nil {
		return
	}
	s.WithTerminator = bs.withTerminator
	bs.block.Tree = append(bs.block.Tree, s)
}

func IsStatement(current, prev lexer.Token) (ok bool, withTerminator bool) {
	ok = current.Id == lexer.SemiColon || prev.Row < current.Row
	withTerminator = current.Id == lexer.SemiColon
	return
}

func nextStatementPos(tokens []lexer.Token, start int) (int, bool) {
	braceCount := 0
	i := start
	for ; i < len(tokens); i++ {
		var isStatement, withTerminator bool
		tok := tokens[i]
		if tok.Id == lexer.Brace {
			switch tok.Kind {
			case "{", "[", "(":
				braceCount++
				continue
			default:
				braceCount--
				if braceCount == 0 {
					if i+1 < len(tokens) {
						isStatement, withTerminator = IsStatement(tokens[i+1], tok)
						if isStatement {
							i++
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
		if i > start {
			isStatement, withTerminator = IsStatement(tok, tokens[i-1])
		} else {
			isStatement, withTerminator = IsStatement(tok, tok)
		}
		if !isStatement {
			continue
		}
	ret:
		if withTerminator {
			i++
		}
		return i, withTerminator
	}
	return i, false
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
		if b.Pos == -1 {
			return
		}
		i, withTerminator := nextStatementPos(tokens, 0)
		statementToks := tokens[:i]
		bs := new(blockStatement)
		bs.block = &block
		bs.blockTokens = &tokens
		bs.tokens = statementToks
		bs.withTerminator = withTerminator
		b.pushStatementToBlock(bs)
	next:
		if len(bs.nextTokens) > 0 {
			bs.tokens = bs.nextTokens
			bs.nextTokens = nil
			b.pushStatementToBlock(bs)
			goto next
		}
		if i >= len(tokens) {
			break
		}
		tokens = tokens[i:]
	}
	return
}

func (b *Builder) Statement(bs *blockStatement) (s Statement) {
	s, ok := b.AssignStatement(bs.tokens, false)
	if ok {
		return s
	}
	tok := bs.tokens[0]
	switch tok.Id {
	case lexer.Id:
		s, ok := b.IdStatement(bs.tokens)
		if ok {
			return s
		}
	case lexer.Const, lexer.Volatile:
		return b.VarStatement(bs.tokens)
	case lexer.Ret:
		return b.RetStatement(bs.tokens)
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
	case lexer.Type:
		t := b.Type(bs.tokens)
		s.Token = t.Token
		s.Val = t
		return
	case lexer.Operator:
		if tok.Kind == "<" {
			return b.RetStatement(bs.tokens)
		}
	case lexer.Comment:
		return b.CommentStatement(bs.tokens[0])
	}
	return b.ExprStatement(bs.tokens)
}

type assignInfo struct {
	selectorTokens []lexer.Token
	exprTokens     []lexer.Token
	setter         lexer.Token
	ok             bool
	isExpr         bool
}

func (b *Builder) assignInfo(tokens []lexer.Token) (info assignInfo) {
	info.ok = true
	braceCount := 0
	for i, tok := range tokens {
		if tok.Id == lexer.Brace {
			switch tok.Kind {
			case "(", "[", "{":
				braceCount++
			default:
				braceCount--
			}
		}
		if braceCount > 0 {
			continue
		}
		if tok.Id == lexer.Operator &&
			tok.Kind[len(tok.Kind)-1] == '=' {
			info.selectorTokens = tokens[:i]
			if info.selectorTokens == nil {
				b.pusherr(tok, "invalid_syntax")
				info.ok = false
			}
			info.setter = tok
			if i+1 >= len(tokens) {
				// b.pusherr(tok, "missing_expr")
				info.ok = false
			} else {
				info.exprTokens = tokens[i+1:]
			}
			return
		}
	}
	return
}

func (b *Builder) pushAssignSelector(
	selectors *[]AssignSelector,
	last, current int,
	info assignInfo,
) {
	var selector AssignSelector
	selector.Expr.Tokens = info.selectorTokens[last:current]
	if last-current == 0 {
		b.pusherr(info.selectorTokens[current-1], "missing_expr")
		return
	}
	if selector.Expr.Tokens[0].Id == lexer.Id &&
		current-last > 1 &&
		selector.Expr.Tokens[1].Id == lexer.Colon {
		if info.isExpr {
			b.pusherr(selector.Expr.Tokens[0], "notallow_declares")
		}
		selector.Var.New = true
		selector.Var.IdToken = selector.Expr.Tokens[0]
		selector.Var.Id = selector.Var.IdToken.Kind
		selector.Var.SetterToken = info.setter
		if current-last > 2 {
			selector.Var.Type, _ = b.DataType(selector.Expr.Tokens[2:], new(int), false)
		}
	} else {
		if selector.Expr.Tokens[0].Id == lexer.Id {
			selector.Var.IdToken = selector.Expr.Tokens[0]
			selector.Var.Id = selector.Var.IdToken.Kind
		}
		selector.Expr = b.Expr(selector.Expr.Tokens)
	}
	*selectors = append(*selectors, selector)
}

func (b *Builder) assignSelectors(info assignInfo) []AssignSelector {
	var selectors []AssignSelector
	braceCount := 0
	lastIndex := 0
	for i, tok := range info.selectorTokens {
		if tok.Id == lexer.Brace {
			switch tok.Kind {
			case "(", "[", "{":
				braceCount++
			default:
				braceCount--
			}
		}
		if braceCount > 0 {
			continue
		} else if tok.Id != lexer.Comma {
			continue
		}
		b.pushAssignSelector(&selectors, lastIndex, i, info)
		lastIndex = i + 1
	}
	if lastIndex < len(info.selectorTokens) {
		b.pushAssignSelector(&selectors, lastIndex, len(info.selectorTokens), info)
	}
	return selectors
}

func (b *Builder) pushAssignExpr(exps *[]Expr, last, current int, info assignInfo) {
	toks := info.exprTokens[last:current]
	if toks == nil {
		b.pusherr(info.exprTokens[current-1], "missing_expr")
		return
	}
	*exps = append(*exps, b.Expr(toks))
}

func (b *Builder) assignExprs(info assignInfo) []Expr {
	var exprs []Expr
	braceCount := 0
	lastIndex := 0
	for i, tok := range info.exprTokens {
		if tok.Id == lexer.Brace {
			switch tok.Kind {
			case "(", "[", "{":
				braceCount++
			default:
				braceCount--
			}
		}
		if braceCount > 0 {
			continue
		} else if tok.Id != lexer.Comma {
			continue
		}
		b.pushAssignExpr(&exprs, lastIndex, i, info)
		lastIndex = i + 1
	}
	if lastIndex < len(info.exprTokens) {
		b.pushAssignExpr(&exprs, lastIndex, len(info.exprTokens), info)
	}
	return exprs
}

func isAssignTok(id uint8) bool {
	return id == lexer.Id ||
		id == lexer.Brace ||
		id == lexer.Operator
}

func isAssignOperator(kind string) bool {
	return kind == "=" ||
		kind == "+=" ||
		kind == "-=" ||
		kind == "/=" ||
		kind == "*=" ||
		kind == "%=" ||
		kind == ">>=" ||
		kind == "<<=" ||
		kind == "|=" ||
		kind == "&=" ||
		kind == "^="
}

func checkAssignToks(tokens []lexer.Token) bool {
	if !isAssignTok(tokens[0].Id) {
		return false
	}
	braceCount := 0
	for _, tok := range tokens {
		if tok.Id == lexer.Brace {
			switch tok.Kind {
			case "{", "[", "(":
				braceCount++
			default:
				braceCount--
			}
		}
		if braceCount < 0 {
			return false
		} else if braceCount > 0 {
			continue
		}
		if tok.Id == lexer.Operator &&
			isAssignOperator(tok.Kind) {
			return true
		}
	}
	return false
}

func (b *Builder) AssignStatement(tokens []lexer.Token, isExpr bool) (s Statement, _ bool) {
	assign, ok := b.AssignExpr(tokens, isExpr)
	if !ok {
		return
	}
	s.Token = tokens[0]
	s.Val = assign
	return s, true
}

func (b *Builder) AssignExpr(tokens []lexer.Token, isExpr bool) (assign Assign, ok bool) {
	if !checkAssignToks(tokens) {
		return
	}
	info := b.assignInfo(tokens)
	if !info.ok {
		return
	}
	ok = true
	info.isExpr = isExpr
	assign.IsExpr = isExpr
	assign.Setter = info.setter
	assign.SelectExprs = b.assignSelectors(info)
	if isExpr && len(assign.SelectExprs) > 1 {
		b.pusherr(assign.Setter, "notallow_multiple_assign")
	}
	assign.ValueExprs = b.assignExprs(info)
	return
}

func (b *Builder) IdStatement(tokens []lexer.Token) (s Statement, _ bool) {
	tok := tokens[1]
	switch tok.Id {
	case lexer.Colon:
		return b.VarStatement(tokens), true
	}
	return
}

func (b *Builder) ExprStatement(tokens []lexer.Token) Statement {
	block := ExprStatement{b.Expr(tokens)}
	return Statement{tokens[0], block, false}
}

func (b *Builder) Args(tokens []lexer.Token) []Arg {
	var args []Arg
	last := 0
	braceCount := 0
	for i, tok := range tokens {
		if tok.Id == lexer.Brace {
			switch tok.Kind {
			case "{", "[", "(":
				braceCount++
			default:
				braceCount--
			}
		}
		if braceCount > 0 || tok.Id != lexer.Comma {
			continue
		}
		b.pushArg(&args, tokens[last:i], tok)
		last = i + 1
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

func (b *Builder) pushArg(args *[]Arg, tokens []lexer.Token, err lexer.Token) {
	if len(tokens) == 0 {
		b.pusherr(err, "invalid_syntax")
		return
	}
	var arg Arg
	arg.Token = tokens[0]
	arg.Expr = b.Expr(tokens)
	*args = append(*args, arg)
}

func (b *Builder) VarStatement(tokens []lexer.Token) (s Statement) {
	var vast Var
	vast.Pub = b.pub
	b.pub = false
	i := 0
	vast.DefToken = tokens[i]
	for ; i < len(tokens); i++ {
		tok := tokens[i]
		if tok.Id == lexer.Id {
			break
		}
		switch tok.Id {
		case lexer.Const:
			if vast.Const {
				b.pusherr(tok, "invalid_constant")
				break
			}
			vast.Const = true
		case lexer.Volatile:
			if vast.Volatile {
				b.pusherr(tok, "invalid_volatile")
				break
			}
			vast.Volatile = true
		default:
			b.pusherr(tok, "invalid_syntax")
		}
	}
	if i >= len(tokens) {
		return
	}
	vast.IdToken = tokens[i]
	if vast.IdToken.Id != lexer.Id {
		b.pusherr(vast.IdToken, "invalid_syntax")
	}
	vast.Id = vast.IdToken.Kind
	vast.Type = DataType{Id: jn.Void}

	i++
	if vast.DefToken.File != nil {
		if tokens[i].Id != lexer.Colon {
			b.pusherr(tokens[i], "invalid_syntax")
			return
		}
		i++
	} else {
		i++
	}
	if i < len(tokens) {
		tok := tokens[i]
		t, ok := b.DataType(tokens, &i, false)
		if ok {
			vast.Type = t
			i++
			if i >= len(tokens) {
				goto ret
			}
			tok = tokens[i]
		}
		if tok.Id == lexer.Operator {
			if tok.Kind != "=" {
				b.pusherr(tok, "invalid_syntax")
				return
			}
			valueToks := tokens[i+1:]
			if len(valueToks) == 0 {
				b.pusherr(tok, "missing_expr")
				return
			}
			vast.Val = b.Expr(valueToks)
			vast.SetterToken = tok
		}
	}
ret:
	return Statement{vast.IdToken, vast, false}
}

func (b *Builder) CommentStatement(token lexer.Token) (s Statement) {
	s.Token = token
	token.Kind = strings.TrimSpace(token.Kind[2:])
	if strings.HasPrefix(token.Kind, "cxx:") {
		s.Val = CxxEmbed{token.Kind[4:]}
	} else {
		s.Val = Comment{token.Kind}
	}
	return
}

func (b *Builder) RetStatement(tokens []lexer.Token) Statement {
	var returnModel Ret
	returnModel.Token = tokens[0]
	if len(tokens) > 1 {
		returnModel.Expr = b.Expr(tokens[1:])
	}
	return Statement{returnModel.Token, returnModel, false}
}

func (b *Builder) FreeStatement(tokens []lexer.Token) Statement {
	var free Free
	free.Token = tokens[0]
	tokens = tokens[1:]
	if len(tokens) == 0 {
		b.pusherr(free.Token, "missing_expr")
	} else {
		free.Expr = b.Expr(tokens)
	}
	return Statement{free.Token, free, false}
}

func blockExprToks(tokens []lexer.Token) (expr []lexer.Token) {
	braceCount := 0
	for i, tok := range tokens {
		if tok.Id == lexer.Brace {
			switch tok.Kind {
			case "{":
				if braceCount > 0 {
					braceCount++
					break
				}
				return tokens[:i]
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

func (b *Builder) pushVarsToksPart(
	vars *[][]lexer.Token,
	tokens []lexer.Token,
	errTok lexer.Token,
) {
	if len(tokens) == 0 {
		b.pusherr(errTok, "missing_expr")
	}
	*vars = append(*vars, tokens)
}

func (b *Builder) getForeachVarsToks(tokens []lexer.Token) [][]lexer.Token {
	var vars [][]lexer.Token
	braceCount := 0
	last := 0
	for i, tok := range tokens {
		if tok.Id == lexer.Brace {
			switch tok.Kind {
			case "(", "[", "{":
				braceCount++
			default:
				braceCount--
			}
		}
		if braceCount > 0 {
			continue
		}
		if tok.Id == lexer.Comma {
			part := tokens[last:i]
			b.pushVarsToksPart(&vars, part, tok)
			last = i + 1
		}
	}
	if last < len(tokens) {
		part := tokens[last:]
		b.pushVarsToksPart(&vars, part, tokens[last])
	}
	return vars
}

func (b *Builder) getForeachIterVars(varsTokens [][]lexer.Token) []Var {
	var vars []Var
	for _, toks := range varsTokens {
		var vast Var
		vast.IdToken = toks[0]
		if vast.IdToken.Id != lexer.Id {
			b.pusherr(vast.IdToken, "invalid_syntax")
			vars = append(vars, vast)
			continue
		}
		vast.Id = vast.IdToken.Kind
		if len(toks) == 1 {
			vars = append(vars, vast)
			continue
		}
		if colon := toks[1]; colon.Id != lexer.Colon {
			b.pusherr(colon, "invalid_syntax")
			vars = append(vars, vast)
			continue
		}
		vast.New = true
		i := new(int)
		*i = 2
		if *i >= len(toks) {
			vars = append(vars, vast)
			continue
		}
		vast.Type, _ = b.DataType(toks, i, true)
		if *i < len(toks)-1 {
			b.pusherr(toks[*i], "invalid_syntax")
		}
		vars = append(vars, vast)
	}
	return vars
}

func (b *Builder) getForeachIterProfile(
	varTokens, exprTokens []lexer.Token,
	inToken lexer.Token,
) ForeachProfile {
	var profile ForeachProfile
	profile.InToken = inToken
	profile.Expr = b.Expr(exprTokens)
	if len(varTokens) == 0 {
		profile.KeyA.Id = jnapi.Ignore
		profile.KeyB.Id = jnapi.Ignore
	} else {
		varsToks := b.getForeachVarsToks(varTokens)
		if len(varsToks) == 0 {
			return profile
		}
		if len(varsToks) > 2 {
			b.pusherr(inToken, "much_foreach_vars")
		}
		vars := b.getForeachIterVars(varsToks)
		profile.KeyA = vars[0]
		if len(vars) > 1 {
			profile.KeyB = vars[1]
		} else {
			profile.KeyB.Id = jnapi.Ignore
		}
	}
	return profile
}

func (b *Builder) getIterProfile(tokens []lexer.Token) IterProfile {
	braceCount := 0
	for i, tok := range tokens {
		if tok.Id == lexer.Brace {
			switch tok.Kind {
			case "(", "[", "{":
				braceCount++
			default:
				braceCount--
			}
		}
		if braceCount != 0 {
			continue
		}
		if tok.Id == lexer.In {
			varToks := tokens[:i]
			exprToks := tokens[i+1:]
			return b.getForeachIterProfile(varToks, exprToks, tok)
		}
	}
	return b.getWhileIterProfile(tokens)
}

func (b *Builder) IterExpr(tokens []lexer.Token) (s Statement) {
	var iter Iter
	iter.Token = tokens[0]
	tokens = tokens[1:]
	if len(tokens) == 0 {
		b.pusherr(iter.Token, "body_not_exist")
		return
	}
	exprToks := blockExprToks(tokens)
	if len(exprToks) > 0 {
		iter.Profile = b.getIterProfile(exprToks)
	}
	i := new(int)
	*i = len(exprToks)
	blockToks := b.getrange(i, "{", "}", &tokens)
	if blockToks == nil {
		b.pusherr(iter.Token, "body_not_exist")
		return
	}
	if *i < len(tokens) {
		b.pusherr(tokens[*i], "invalid_syntax")
	}
	iter.Block = b.Block(blockToks)
	return Statement{iter.Token, iter, false}
}

func (b *Builder) IfExpr(bs *blockStatement) (s Statement) {
	var ifast If
	ifast.Token = bs.tokens[0]
	bs.tokens = bs.tokens[1:]
	exprToks := blockExprToks(bs.tokens)
	if len(exprToks) == 0 {
		b.pusherr(ifast.Token, "missing_expr")
	}
	i := new(int)
	*i = len(exprToks)
	blockToks := b.getrange(i, "{", "}", &bs.tokens)
	if blockToks == nil {
		b.pusherr(ifast.Token, "body_not_exist")
		return
	}
	if *i < len(bs.tokens) {
		if bs.tokens[*i].Id == lexer.Else {
			bs.nextTokens = bs.tokens[*i:]
		} else {
			b.pusherr(bs.tokens[*i], "invalid_syntax")
		}
	}
	ifast.Expr = b.Expr(exprToks)
	ifast.Block = b.Block(blockToks)
	return Statement{ifast.Token, ifast, false}
}

func (b *Builder) ElseIfExpr(bs *blockStatement) (s Statement) {
	var elif ElseIf
	elif.Token = bs.tokens[1]
	bs.tokens = bs.tokens[2:]
	exprToks := blockExprToks(bs.tokens)
	if len(exprToks) == 0 {
		b.pusherr(elif.Token, "missing_expr")
	}
	i := new(int)
	*i = len(exprToks)
	blockToks := b.getrange(i, "{", "}", &bs.tokens)
	if blockToks == nil {
		b.pusherr(elif.Token, "body_not_exist")
		return
	}
	if *i < len(bs.tokens) {
		if bs.tokens[*i].Id == lexer.Else {
			bs.nextTokens = bs.tokens[*i:]
		} else {
			b.pusherr(bs.tokens[*i], "invalid_syntax")
		}
	}
	elif.Expr = b.Expr(exprToks)
	elif.Block = b.Block(blockToks)
	return Statement{elif.Token, elif, false}
}

func (b *Builder) ElseBlock(bs *blockStatement) (s Statement) {
	if len(bs.tokens) > 1 && bs.tokens[1].Id == lexer.If {
		return b.ElseIfExpr(bs)
	}
	var elseast Else
	elseast.Token = bs.tokens[0]
	bs.tokens = bs.tokens[1:]
	i := new(int)
	blockToks := b.getrange(i, "{", "}", &bs.tokens)
	if blockToks == nil {
		if *i < len(bs.tokens) {
			b.pusherr(elseast.Token, "else_have_expr")
		} else {
			b.pusherr(elseast.Token, "body_not_exist")
		}
		return
	}
	if *i < len(bs.tokens) {
		b.pusherr(bs.tokens[*i], "invalid_syntax")
	}
	elseast.Block = b.Block(blockToks)
	return Statement{elseast.Token, elseast, false}
}

func (b *Builder) BreakStatement(tokens []lexer.Token) Statement {
	var breakAST Break
	breakAST.Token = tokens[0]
	if len(tokens) > 1 {
		b.pusherr(tokens[1], "invalid_syntax")
	}
	return Statement{breakAST.Token, breakAST, false}
}

func (b *Builder) ContinueStatement(tokens []lexer.Token) Statement {
	var continueAST Continue
	continueAST.Token = tokens[0]
	if len(tokens) > 1 {
		b.pusherr(tokens[1], "invalid_syntax")
	}
	return Statement{continueAST.Token, continueAST, false}
}

func (b *Builder) Expr(tokens []lexer.Token) (e Expr) {
	e.Processes = b.getExprProcesses(tokens)
	e.Tokens = tokens
	return
}

func isOverflowOperator(kind string) bool {
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

func isExprOperator(kind string) bool {
	return kind == "..."
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
	for i := 0; i < len(tokens); i++ {
		tok := tokens[i]
		switch tok.Id {
		case lexer.Operator:
			if newKeyword ||
				isExprOperator(tok.Kind) ||
				isAssignOperator(tok.Kind) {
				part = append(part, tok)
				continue
			}
			if !operator {
				if IsSingleOperator(tok.Kind) && !singleOperatored {
					part = append(part, tok)
					singleOperatored = true
					continue
				}
				if braceCount == 0 && isOverflowOperator(tok.Kind) {
					b.pusherr(tok, "operator_overflow")
				}
			}
			singleOperatored = false
			operator = false
			value = true
			if braceCount > 0 {
				part = append(part, tok)
				continue
			}
			processes = append(processes, part)
			processes = append(processes, []lexer.Token{tok})
			part = []lexer.Token{}
			continue
		case lexer.Brace:
			switch tok.Kind {
			case "(", "[", "{":
				if tok.Kind == "[" {
					oldIndex := i
					_, ok := b.DataType(tokens, &i, false)
					if ok {
						part = append(part, tokens[oldIndex:i+1]...)
						continue
					}
					i = oldIndex
				}
				singleOperatored = false
				braceCount++
			default:
				braceCount--
			}
		case lexer.New:
			newKeyword = true
		case lexer.Id:
			if braceCount == 0 {
				newKeyword = false
			}
		}
		if i > 0 && braceCount == 0 {
			lt := tokens[i-1]
			if (lt.Id == lexer.Id || lt.Id == lexer.Value) &&
				(tok.Id == lexer.Id || tok.Id == lexer.Value) {
				b.pusherr(tok, "invalid_syntax")
				pushedError = true
			}
		}
		b.checkExprTok(tok)
		part = append(part, tok)
		operator = requireOperatorForProcess(tok, i, len(tokens))
		value = false
	}
	if len(part) > 0 {
		processes = append(processes, part)
	}
	if value {
		b.pusherr(processes[len(processes)-1][0], "operator_overflow")
		pushedError = true
	}
	if pushedError {
		return nil
	}
	return processes
}

func requireOperatorForProcess(token lexer.Token, index, len int) bool {
	switch token.Id {
	case lexer.Comma:
		return false
	case lexer.Brace:
		if token.Kind == "(" ||
			token.Kind == "{" {
			return false
		}
	}
	return index < len-1
}

func (b *Builder) checkExprTok(token lexer.Token) {
	if token.Kind[0] >= '0' && token.Kind[0] <= '9' {
		var result bool
		if strings.Contains(token.Kind, ".") ||
			strings.ContainsAny(token.Kind, "eE") {
			result = jnbits.CheckBitFloat(token.Kind, 64)
		} else {
			result = jnbits.CheckBitInt(token.Kind, 64)
			if !result {
				result = jnbits.CheckBitUInt(token.Kind, 64)
			}
		}
		if !result {
			b.pusherr(token, "invalid_numeric_range")
		}
	}
}

func (b *Builder) getrange(i *int, open, close string, toks *[]lexer.Token) []lexer.Token {
	rang := getrange(i, open, close, *toks)
	if rang != nil {
		return rang
	}
	if b.Ended() {
		return nil
	}
	*i = 0
	*toks = b.skipStatement()
	rang = getrange(i, open, close, *toks)
	return rang
}

func getrange(i *int, open, close string, tokens []lexer.Token) []lexer.Token {
	if *i >= len(tokens) {
		return nil
	}
	tok := tokens[*i]
	if tok.Id == lexer.Brace && tok.Kind == open {
		*i++
		braceCount := 1
		start := *i
		for ; braceCount > 0 && *i < len(tokens); *i++ {
			tok := tokens[*i]
			if tok.Id != lexer.Brace {
				continue
			}
			if tok.Kind == open {
				braceCount++
			} else if tok.Kind == close {
				braceCount--
			}
		}
		return tokens[start : *i-1]
	}
	return nil
}

func (b *Builder) skipStatement() []lexer.Token {
	start := b.Pos
	b.Pos, _ = nextStatementPos(b.Tokens, start)
	toks := b.Tokens[start:b.Pos]
	if toks[len(toks)-1].Id == lexer.SemiColon {
		if len(toks) == 1 {
			return b.skipStatement()
		}
		toks = toks[:len(toks)-1]
	}
	return toks
}
