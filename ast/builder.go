package ast

import (
	"os"
	"strings"
	"sync"

	"github.com/DeRuneLabs/jane/ast/models"
	"github.com/DeRuneLabs/jane/lexer"
	"github.com/DeRuneLabs/jane/lexer/tokens"
	"github.com/DeRuneLabs/jane/package/jn"
	"github.com/DeRuneLabs/jane/package/jnapi"
	"github.com/DeRuneLabs/jane/package/jnbits"
	"github.com/DeRuneLabs/jane/package/jnlog"
	"github.com/DeRuneLabs/jane/package/jntype"
)

type (
	Tok  = lexer.Tok
	Toks = []Tok
)

type Builder struct {
	wg     sync.WaitGroup
	pub    bool
	Tree   []models.Object
	Errors []jnlog.CompilerLog
	Toks   Toks
	Pos    int
}

func NewBuilder(toks Toks) *Builder {
	b := new(Builder)
	b.Toks = toks
	b.Pos = 0
	return b
}

func compilerErr(tok Tok, key string, args ...any) jnlog.CompilerLog {
	return jnlog.CompilerLog{
		Type:    jnlog.Error,
		Row:     tok.Row,
		Column:  tok.Column,
		Path:    tok.File.Path(),
		Message: jn.GetError(key, args...),
	}
}

func (b *Builder) pusherr(tok Tok, key string, args ...any) {
	b.Errors = append(b.Errors, compilerErr(tok, key, args...))
}

func (ast *Builder) Ended() bool {
	return ast.Pos >= len(ast.Toks)
}

func (b *Builder) buildNode(toks Toks) {
	tok := toks[0]
	switch tok.Id {
	case tokens.Use:
		b.Use(toks)
	case tokens.At:
		b.Attribute(toks)
	case tokens.Id:
		b.Id(toks)
	case tokens.Const, tokens.Volatile:
		b.GlobalVar(toks)
	case tokens.Type:
		b.Tree = append(b.Tree, b.TypeOrGenerics(toks))
	case tokens.Enum:
		b.Enum(toks)
	case tokens.Struct:
		b.Struct(toks)
	case tokens.Comment:
		b.Comment(toks[0])
	case tokens.Preprocessor:
		b.Preprocessor(toks)
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
		toks := b.nextBuilderStatement()
		b.pub = toks[0].Id == tokens.Pub
		if b.pub {
			if len(toks) == 1 {
				if b.Ended() {
					b.pusherr(toks[0], "invalid_syntax")
					continue
				}
				toks = b.nextBuilderStatement()
			} else {
				toks = toks[1:]
			}
		}
		b.buildNode(toks)
	}
	b.Wait()
}

func (b *Builder) Wait() {
	b.wg.Wait()
}

func (b *Builder) Type(toks Toks) (t models.Type) {
	i := 1
	if i >= len(toks) {
		b.pusherr(toks[i-1], "invalid_syntax")
		return
	}
	tok := toks[i]
	if tok.Id != tokens.Id {
		b.pusherr(tok, "invalid_syntax")
	}
	i++
	if i >= len(toks) {
		b.pusherr(toks[i-1], "invalid_syntax")
		return
	}
	destType, _ := b.DataType(toks[i:], new(int), true)
	tok = toks[1]
	return models.Type{
		Tok:  tok,
		Id:   tok.Kind,
		Type: destType,
	}
}

func (b *Builder) buildEnumItemExpr(i *int, toks Toks) models.Expr {
	braceCount := 0
	exprStart := *i
	for ; *i < len(toks); *i++ {
		tok := toks[*i]
		if tok.Id == tokens.Brace {
			switch tok.Kind {
			case tokens.LBRACE, tokens.LBRACKET, tokens.LPARENTHESES:
				braceCount++
				continue
			default:
				braceCount--
			}
		}
		if braceCount > 0 {
			continue
		}
		if tok.Id == tokens.Comma || *i+1 >= len(toks) {
			var exprToks Toks
			if tok.Id == tokens.Comma {
				exprToks = toks[exprStart:*i]
			} else {
				exprToks = toks[exprStart:]
			}
			return b.Expr(exprToks)
		}
	}
	return models.Expr{}
}

func (b *Builder) buildEnumItems(toks Toks) []*models.EnumItem {
	items := make([]*models.EnumItem, 0)
	for i := 0; i < len(toks); i++ {
		tok := toks[i]
		item := new(models.EnumItem)
		item.Tok = tok
		if item.Tok.Id != tokens.Id {
			b.pusherr(item.Tok, "invalid_syntax")
		}
		item.Id = item.Tok.Kind
		if i+1 >= len(toks) || toks[i+1].Id == tokens.Comma {
			if i+1 < len(toks) {
				i++
			}
			items = append(items, item)
			continue
		}
		i++
		tok = toks[i]
		if tok.Id != tokens.Operator && tok.Kind != tokens.EQUAL {
			b.pusherr(toks[0], "invalid_syntax")
		}
		i++
		if i >= len(toks) || toks[i].Id == tokens.Comma {
			b.pusherr(toks[0], "missing_expr")
			continue
		}
		item.Expr = b.buildEnumItemExpr(&i, toks)
		items = append(items, item)
	}
	return items
}

func (b *Builder) Enum(toks Toks) {
	var enum models.Enum
	if len(toks) < 2 || len(toks) < 3 {
		b.pusherr(toks[0], "invalid_syntax")
		return
	}
	enum.Tok = toks[1]
	if enum.Tok.Id != tokens.Id {
		b.pusherr(enum.Tok, "invalid_syntax")
	}
	enum.Id = enum.Tok.Kind
	i := 2
	if toks[i].Id == tokens.Colon {
		i++
		if i >= len(toks) {
			b.pusherr(toks[i-1], "invalid_syntax")
			return
		}
		enum.Type, _ = b.DataType(toks, &i, true)
		i++
		if i >= len(toks) {
			b.pusherr(enum.Tok, "body_not_exist")
			return
		}
	} else {
		enum.Type = models.DataType{Id: jntype.U32, Kind: tokens.U32}
	}
	itemToks := b.getrange(&i, tokens.LBRACE, tokens.RBRACE, &toks)
	if itemToks == nil {
		b.pusherr(enum.Tok, "body_not_exist")
		return
	} else if i < len(toks) {
		b.pusherr(toks[i], "invalid_syntax")
	}
	enum.Pub = b.pub
	b.pub = false
	enum.Items = b.buildEnumItems(itemToks)
	b.Tree = append(b.Tree, models.Object{
		Tok:   enum.Tok,
		Value: enum,
	})
}

func (b *Builder) Comment(tok Tok) {
	tok.Kind = strings.TrimSpace(tok.Kind[2:])
	if strings.HasPrefix(tok.Kind, "cxx:") {
		b.Tree = append(b.Tree, models.Object{
			Tok: tok,
			Value: models.CxxEmbed{
				Tok:     tok,
				Content: tok.Kind[4:],
			},
		})
		return
	}
	b.Tree = append(b.Tree, models.Object{
		Tok: tok,
		Value: models.Comment{
			Content: tok.Kind,
		},
	})
}

func (b *Builder) Preprocessor(toks Toks) {
	if len(toks) == 1 {
		b.pusherr(toks[0], "invalid_syntax")
		return
	}
	var pp models.Preprocessor
	toks = toks[1:]
	tok := toks[0]
	if tok.Id != tokens.Id {
		b.pusherr(pp.Tok, "invalid_syntax")
		return
	}
	ok := false
	switch tok.Kind {
	case jn.PreprocessorDirective:
		ok = b.PreprocessorDirective(&pp, toks)
	default:
		b.pusherr(tok, "invalid_preprocessor")
		return
	}
	if ok {
		b.Tree = append(b.Tree, models.Object{
			Tok:   pp.Tok,
			Value: pp,
		})
	}
}

func (b *Builder) PreprocessorDirective(pp *models.Preprocessor, toks Toks) bool {
	if len(toks) == 1 {
		b.pusherr(toks[0], "missing_pragma_directive")
		return false
	}
	toks = toks[1:]
	tok := toks[0]
	if tok.Id != tokens.Id {
		b.pusherr(tok, "invalid_syntax")
		return false
	}
	var d models.Directive
	ok := false
	switch tok.Kind {
	case jn.PreprocessorDirectiveEnofi:
		ok = b.directiveEnofi(&d, toks)
	default:
		b.pusherr(tok, "invalid_pragma_directive")
	}
	pp.Command = d
	return ok
}

func (b *Builder) directiveEnofi(d *models.Directive, toks Toks) bool {
	if len(toks) > 1 {
		b.pusherr(toks[1], "invalid_syntax")
		return false
	}
	d.Command = models.DirectiveEnofi{}
	return true
}

func (b *Builder) Id(toks Toks) {
	if len(toks) == 1 {
		b.pusherr(toks[0], "invalid_syntax")
		return
	}
	tok := toks[1]
	switch tok.Id {
	case tokens.Colon:
		b.GlobalVar(toks)
		return
	case tokens.DoubleColon:
		b.Namespace(toks)
		return
	case tokens.Brace:
		switch tok.Kind {
		case tokens.LBRACE:
			b.Namespace(toks)
			return
		case tokens.LPARENTHESES:
			f := b.Func(toks, false)
			s := models.Statement{
				Tok: f.Tok,
				Val: f,
			}
			b.Tree = append(b.Tree, models.Object{
				Tok:   f.Tok,
				Value: s,
			})
			return
		}
	}
	b.pusherr(tok, "invalid_syntax")
}

func (b *Builder) nsIds(toks Toks, i *int) []string {
	var ids []string
	for ; *i < len(toks); *i++ {
		tok := toks[*i]
		if (*i+1)%2 != 0 {
			if tok.Id != tokens.Id {
				b.pusherr(tok, "invalid_syntax")
				continue
			}
			ids = append(ids, tok.Kind)
			continue
		}
		switch tok.Id {
		case tokens.DoubleColon:
			continue
		default:
			goto ret
		}
	}
ret:
	return ids
}

func (b *Builder) Namespace(toks Toks) {
	var ns models.Namespace
	ns.Tok = toks[0]
	i := new(int)
	ns.Ids = b.nsIds(toks, i)
	treeToks := b.getrange(i, tokens.LBRACE, tokens.RBRACE, &toks)
	if treeToks == nil {
		b.pusherr(ns.Tok, "body_not_exist")
		return
	}
	if *i < len(toks) {
		b.pusherr(toks[*i], "invalid_syntax")
	}
	tree := b.Tree
	b.Tree = nil
	btoks := b.Toks
	pos := b.Pos
	b.Toks = treeToks
	b.Pos = 0
	b.Build()
	b.Toks = btoks
	b.Pos = pos
	ns.Tree = b.Tree
	b.Tree = tree
	b.Tree = append(b.Tree, models.Object{
		Tok:   ns.Tok,
		Value: ns,
	})
}

func (b *Builder) structFields(toks Toks) []*models.Var {
	fields := make([]*models.Var, 0)
	i := new(int)
	for *i < len(toks) {
		varToks := b.skipStatement(i, &toks)
		pub := varToks[0].Id == tokens.Pub
		if pub {
			if len(varToks) == 1 {
				b.pusherr(varToks[0], "invalid_syntax")
				continue
			}
			varToks = varToks[1:]
		}
		vast := b.Var(varToks)
		vast.Pub = pub
		fields = append(fields, &vast)
	}
	return fields
}

func (b *Builder) Struct(toks Toks) {
	var s models.Struct
	s.Pub = b.pub
	b.pub = false
	if len(toks) < 3 {
		b.pusherr(toks[0], "invalid_syntax")
		return
	}
	s.Tok = toks[1]
	if s.Tok.Id != tokens.Id {
		b.pusherr(s.Tok, "invalid_syntax")
	}
	s.Id = s.Tok.Kind
	i := 2
	bodyToks := b.getrange(&i, tokens.LBRACE, tokens.RBRACE, &toks)
	if bodyToks == nil {
		b.pusherr(s.Tok, "body_not_exist")
		return
	}
	if i < len(toks) {
		b.pusherr(toks[i], "invalid_syntax")
	}
	s.Fields = b.structFields(bodyToks)
	b.Tree = append(b.Tree, models.Object{
		Tok:   s.Tok,
		Value: s,
	})
}

func tokstoa(toks Toks) string {
	var str strings.Builder
	for _, tok := range toks {
		str.WriteString(tok.Kind)
	}
	return str.String()
}

func (b *Builder) Use(toks Toks) {
	var use models.Use
	use.Tok = toks[0]
	if len(toks) < 2 {
		b.pusherr(use.Tok, "missing_use_path")
		return
	}
	toks = toks[1:]
	use.LinkString = tokstoa(toks)
	use.Path = b.usePath(toks)
	b.Tree = append(b.Tree, models.Object{
		Tok:   use.Tok,
		Value: use,
	})
}

func (b *Builder) usePath(toks Toks) string {
	var path strings.Builder
	path.WriteString(jn.StdlibPath)
	path.WriteRune(os.PathSeparator)
	for i, tok := range toks {
		if i%2 != 0 {
			if tok.Id != tokens.Dot {
				b.pusherr(tok, "invalid_syntax")
			}
			path.WriteRune(os.PathSeparator)
			continue
		}
		if tok.Id != tokens.Id {
			b.pusherr(tok, "invalid_syntax")
		}
		path.WriteString(tok.Kind)
	}
	return path.String()
}

func (b *Builder) Attribute(toks Toks) {
	var a models.Attribute
	i := 0
	a.Tok = toks[i]
	i++
	if b.Ended() {
		b.pusherr(toks[i-1], "invalid_syntax")
		return
	}
	a.Tag = toks[i]
	if a.Tag.Id != tokens.Id || a.Tok.Column+1 != a.Tag.Column {
		b.pusherr(a.Tag, "invalid_syntax")
		return
	}
	b.Tree = append(b.Tree, models.Object{
		Tok:   a.Tok,
		Value: a,
	})
}

func (b *Builder) Func(toks Toks, anon bool) (f models.Func) {
	f.Tok = toks[0]
	i := 0
	f.Pub = b.pub
	b.pub = false
	if anon {
		f.Id = jn.Anonymous
	} else {
		if f.Tok.Id != tokens.Id {
			b.pusherr(f.Tok, "invalid_syntax")
		}
		f.Id = f.Tok.Kind
		i++
	}
	f.RetType.Type.Id = jntype.Void
	f.RetType.Type.Kind = jntype.VoidTypeStr
	paramToks := b.getrange(&i, tokens.LPARENTHESES, tokens.RPARENTHESES, &toks)
	if len(paramToks) > 0 {
		b.Params(&f, paramToks)
	}
	if i >= len(toks) {
		if b.Ended() {
			b.pusherr(f.Tok, "body_not_exist")
			return
		}
		i = 0
		toks = b.nextBuilderStatement()
	}
	tok := toks[i]
	t, ok := b.FuncRetDataType(toks, &i)
	if ok {
		f.RetType = t
		i++
		if i >= len(toks) {
			if b.Ended() {
				b.pusherr(f.Tok, "body_not_exist")
				return
			}
			i = 0
			toks = b.nextBuilderStatement()
		}
		tok = toks[i]
	}
	if tok.Id != tokens.Brace || tok.Kind != tokens.LBRACE {
		b.pusherr(tok, "invalid_syntax")
		return
	}
	blockToks := b.getrange(&i, tokens.LBRACE, tokens.RBRACE, &toks)
	if blockToks == nil {
		b.pusherr(f.Tok, "body_not_exist")
		return
	} else if i < len(toks) {
		b.pusherr(toks[i], "invalid_syntax")
	}
	f.Block = b.Block(blockToks)
	return
}

func (b *Builder) generic(toks Toks) models.GenericType {
	if len(toks) > 1 {
		b.pusherr(toks[1], "invalid_syntax")
	}
	var gt models.GenericType
	gt.Tok = toks[0]
	if gt.Tok.Id != tokens.Id {
		b.pusherr(gt.Tok, "invalid_syntax")
	}
	gt.Id = gt.Tok.Kind
	return gt
}

func (b *Builder) Generics(toks Toks) []models.GenericType {
	tok := toks[0]
	i := 1
	genericsToks := Range(&i, tokens.LBRACKET, tokens.RBRACKET, toks)
	if len(genericsToks) == 0 {
		b.pusherr(tok, "missing_expr")
		return make([]models.GenericType, 0)
	} else if i < len(toks) {
		b.pusherr(toks[i], "invalid_syntax")
	}
	parts, errs := Parts(genericsToks, tokens.Comma)
	b.Errors = append(b.Errors, errs...)
	generics := make([]models.GenericType, len(parts))
	for i, part := range parts {
		if len(parts) == 0 {
			continue
		}
		generics[i] = b.generic(part)
	}
	return generics
}

func (b *Builder) TypeOrGenerics(toks Toks) models.Object {
	if len(toks) > 1 {
		tok := toks[1]
		if tok.Id == tokens.Brace && tok.Kind == tokens.LBRACKET {
			generics := b.Generics(toks)
			return models.Object{
				Tok:   tok,
				Value: generics,
			}
		}
	}
	t := b.Type(toks)
	t.Pub = b.pub
	b.pub = false
	return models.Object{
		Tok:   t.Tok,
		Value: t,
	}
}

func (b *Builder) GlobalVar(toks Toks) {
	if toks == nil {
		return
	}
	s := b.VarStatement(toks)
	b.Tree = append(b.Tree, models.Object{
		Tok:   s.Tok,
		Value: s,
	})
}

func (b *Builder) Params(f *models.Func, toks Toks) {
	parts, errs := Parts(toks, tokens.Comma)
	b.Errors = append(b.Errors, errs...)
	for _, part := range parts {
		b.pushParam(f, part)
	}
	b.wg.Add(1)
	go b.checkParamsAsync(f)
}

func (b *Builder) checkParamsAsync(f *models.Func) {
	defer func() { b.wg.Done() }()
	for i := range f.Params {
		p := &f.Params[i]
		if p.Type.Tok.Id == tokens.NA {
			if p.Tok.Id == tokens.NA {
				b.pusherr(p.Tok, "missing_type")
			} else {
				p.Type.Tok = p.Tok
				p.Type.Id = jntype.Id
				p.Type.Kind = p.Type.Tok.Kind
				p.Type.Original = p.Type
				p.Id = jn.Anonymous
				p.Tok = lexer.Tok{}
			}
		}
	}
}

func (b *Builder) paramBegin(p *models.Param, i *int, toks Toks) {
	for ; *i < len(toks); *i++ {
		tok := toks[*i]
		switch tok.Id {
		case tokens.Const:
			if p.Const {
				b.pusherr(tok, "already_constant")
				continue
			}
			p.Const = true
		case tokens.Volatile:
			if p.Volatile {
				b.pusherr(tok, "already_volatile")
				continue
			}
			p.Volatile = true
		case tokens.Operator:
			switch tok.Kind {
			case tokens.TRIPLE_DOT:
				if p.Variadic {
					b.pusherr(tok, "already_variadic")
					continue
				}
				p.Variadic = true
			case tokens.AMPER:
				if p.Reference {
					b.pusherr(tok, "already_reference")
					continue
				}
				p.Reference = true
			default:
				b.pusherr(tok, "invalid_syntax")
			}
		default:
			return
		}
	}
}

func (b *Builder) paramBodyId(f *models.Func, p *models.Param, tok Tok) {
	if jnapi.IsIgnoreId(tok.Kind) {
		p.Id = jn.Anonymous
		return
	}
	for _, param := range f.Params {
		if param.Id == tok.Kind {
			b.pusherr(tok, "parameter_exist", tok.Kind)
			break
		}
	}
	p.Id = tok.Kind
}

type exprNode struct {
	expr string
}

func (en exprNode) String() string {
	return en.expr
}

func (b *Builder) paramBodyDefaultExpr(p *models.Param, toks *Toks) {
	braceCount := 0
	for i, tok := range *toks {
		if tok.Id == tokens.Brace {
			switch tok.Kind {
			case tokens.LBRACE, tokens.LBRACKET, tokens.LPARENTHESES:
				braceCount++
				continue
			default:
				braceCount--
			}
		}
		if braceCount > 0 {
			continue
		}
		exprToks := (*toks)[1:i]
		*toks = (*toks)[i+1:]
		if len(exprToks) > 0 {
			p.Default = b.Expr(exprToks)
		} else {
			p.Default.Model = exprNode{jnapi.DefaultExpr}
		}
		break
	}
}

func (b *Builder) paramBodyDataType(f *models.Func, p *models.Param, toks Toks) {
	i := 0
	p.Type, _ = b.DataType(toks, &i, true)
	i++
	if i < len(toks) {
		b.pusherr(toks[i], "invalid_syntax")
	}
	i = len(f.Params) - 1
	for ; i >= 0; i-- {
		param := &f.Params[i]
		if param.Type.Tok.Id != tokens.NA {
			break
		}
		param.Type = p.Type
	}
}

func (b *Builder) paramBody(f *models.Func, p *models.Param, i *int, toks Toks) {
	b.paramBodyId(f, p, toks[*i])
	toks = toks[*i+1:]
	if len(toks) == 0 {
		return
	}
	if tok := toks[0]; tok.Id == tokens.Brace && tok.Kind == tokens.LBRACE {
		b.paramBodyDefaultExpr(p, &toks)
	}
	if len(toks) > 0 {
		b.paramBodyDataType(f, p, toks)
	}
}

func (b *Builder) pushParam(f *models.Func, toks Toks) {
	var p models.Param
	i := 0
	b.paramBegin(&p, &i, toks)
	if i >= len(toks) {
		return
	}
	tok := toks[i]
	p.Tok = tok
	if tok.Id != tokens.Id {
		if t, ok := b.DataType(toks, &i, true); ok {
			if i+1 == len(toks) {
				p.Type = t
				goto end
			}
		}
		b.pusherr(tok, "invalid_syntax")
		goto end
	}
	b.paramBody(f, &p, &i, toks)
end:
	f.Params = append(f.Params, p)
}

func (b *Builder) idGenericsParts(toks Toks, i *int) []Toks {
	first := *i
	braceCount := 0
	for ; *i < len(toks); *i++ {
		tok := toks[*i]
		if tok.Id == tokens.Brace {
			switch tok.Kind {
			case tokens.LBRACKET:
				braceCount++
			case tokens.RBRACKET:
				braceCount--
			}
		}
		if braceCount == 0 {
			break
		}
	}
	toks = toks[first+1 : *i]
	parts, errs := Parts(toks, tokens.Comma)
	b.Errors = append(b.Errors, errs...)
	return parts
}

func (b *Builder) idDataTypePartEnd(t *models.DataType, dtv *strings.Builder, toks Toks, i *int) {
	if *i+1 >= len(toks) {
		return
	}
	*i++
	tok := toks[*i]
	if tok.Id != tokens.Brace || tok.Kind != tokens.LBRACKET {
		*i--
		return
	}
	dtv.WriteByte('[')
	var genericsStr strings.Builder
	parts := b.idGenericsParts(toks, i)
	generics := make([]models.DataType, len(parts))
	for i, part := range parts {
		index := 0
		t, _ := b.DataType(part, &index, true)
		if index+1 < len(part) {
			b.pusherr(part[index+1], "invalid_syntax")
		}
		genericsStr.WriteString(t.String())
		genericsStr.WriteByte(',')
		generics[i] = t
	}
	dtv.WriteString(genericsStr.String()[:genericsStr.Len()-1])
	dtv.WriteByte(']')
	t.Tag = generics
}

func (b *Builder) DataType(toks Toks, i *int, err bool) (t models.DataType, ok bool) {
	defer func() { t.Original = t }()
	first := *i
	var dtv strings.Builder
	for ; *i < len(toks); *i++ {
		tok := toks[*i]
		switch tok.Id {
		case tokens.DataType:
			t.Tok = tok
			t.Id = jntype.TypeFromId(t.Tok.Kind)
			dtv.WriteString(t.Tok.Kind)
			ok = true
			goto ret
		case tokens.Id:
			t.Id = jntype.Id
			t.Tok = tok
			dtv.WriteString(t.Tok.Kind)
			b.idDataTypePartEnd(&t, &dtv, toks, i)
			ok = true
			goto ret
		case tokens.Operator:
			if tok.Kind == tokens.STAR {
				dtv.WriteString(tok.Kind)
				break
			}
			if err {
				b.pusherr(tok, "invalid_syntax")
			}
			return
		case tokens.Brace:
			switch tok.Kind {
			case tokens.LPARENTHESES:
				t.Tok = tok
				t.Id = jntype.Func
				f := b.FuncDataTypeHead(toks, i)
				*i++
				f.RetType, ok = b.FuncRetDataType(toks, i)
				if !ok {
					*i--
				}
				t.Tag = &f
				dtv.WriteString(f.DataTypeString())
				ok = true
				goto ret
			case tokens.LBRACKET:
				*i++
				if *i > len(toks) {
					if err {
						b.pusherr(tok, "invalid_syntax")
					}
					return
				}
				tok = toks[*i]
				if tok.Id == tokens.Brace && tok.Kind == tokens.RBRACKET {
					dtv.WriteString("[]")
					continue
				}
				*i--
				dt, val := b.MapDataType(toks, i, err)
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
		b.pusherr(toks[first], "invalid_type")
	}
ret:
	t.Kind = dtv.String()
	return
}

func (b *Builder) MapDataType(toks Toks, i *int, err bool) (t models.DataType, _ string) {
	defer func() { t.Original = t }()
	t.Id = jntype.Map
	t.Tok = toks[0]
	typeToks, colon := MapDataTypeInfo(toks, i)
	if typeToks == nil || colon == -1 {
		return
	}
	colonTok := toks[colon]
	if colon == 0 || colon+1 >= len(typeToks) {
		b.pusherr(colonTok, "missing_expr")
		return t, " "
	}
	keyTypeToks := typeToks[:colon]
	valTypeToks := typeToks[colon+1:]
	types := make([]models.DataType, 2)
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
	val.WriteString(types[0].Kind)
	val.WriteByte(':')
	val.WriteString(types[1].Kind)
	val.WriteByte(']')
	return t, val.String()
}

func (b *Builder) FuncDataTypeHead(toks Toks, i *int) models.Func {
	var f models.Func
	brace := 1
	firstIndex := *i
	for *i++; *i < len(toks); *i++ {
		tok := toks[*i]
		switch tok.Id {
		case tokens.Brace:
			switch tok.Kind {
			case tokens.LBRACE, tokens.LBRACKET, tokens.LPARENTHESES:
				brace++
			default:
				brace--
			}
		}
		if brace == 0 {
			b.Params(&f, toks[firstIndex+1:*i])
			return f
		}
	}
	b.pusherr(toks[firstIndex], "invalid_type")
	return f
}

func (b *Builder) pushTypeToTypes(ids *Toks, types *[]models.DataType, toks Toks, errTok Tok) {
	if len(toks) == 0 {
		b.pusherr(errTok, "missing_expr")
		return
	}
	tok := toks[0]
	if tok.Id == tokens.Id && len(toks) > 1 {
		*ids = append(*ids, tok)
		toks = toks[1:]
	} else {
		*ids = append(*ids, Tok{Kind: jnapi.Ignore})
	}
	index := new(int)
	currentDt, ok := b.DataType(toks, index, true)
	if !ok {
		return
	} else if *index < len(toks)-1 {
		b.pusherr(toks[*index], "invalid_syntax")
	}
	*types = append(*types, currentDt)
}

func (b *Builder) funcMultiTypeRet(toks Toks, i *int) (t models.RetType, ok bool) {
	defer func() { t.Type.Original = t.Type }()
	start := *i
	tok := toks[*i]
	t.Type.Kind += tok.Kind
	*i++
	if *i >= len(toks) {
		*i--
		t.Type, ok = b.DataType(toks, i, false)
		return
	}
	tok = toks[*i]
	if tok.Id == tokens.Brace && tok.Kind == tokens.RBRACKET {
		*i--
		t.Type, ok = b.DataType(toks, i, false)
		return
	}
	var types []models.DataType
	braceCount := 1
	last := *i
	for ; *i < len(toks); *i++ {
		tok := toks[*i]
		t.Type.Kind += tok.Kind
		if tok.Id == tokens.Brace {
			switch tok.Kind {
			case tokens.LBRACE, tokens.LBRACKET, tokens.LPARENTHESES:
				braceCount++
			default:
				braceCount--
			}
		}
		if braceCount == 0 {
			if tok.Id == tokens.Colon {
				*i = start
				t.Type, ok = b.DataType(toks, i, false)
				return
			}
			b.pushTypeToTypes(&t.Identifiers, &types, toks[last:*i], toks[last-1])
			break
		} else if braceCount > 1 {
			continue
		}
		switch tok.Id {
		case tokens.Comma:
		case tokens.Colon:
			*i = start
			t.Type, ok = b.DataType(toks, i, false)
			return
		default:
			continue
		}
		b.pushTypeToTypes(&t.Identifiers, &types, toks[last:*i], toks[*i-1])
		last = *i + 1
	}
	if len(types) > 1 {
		t.Type.MultiTyped = true
		t.Type.Tag = types
	} else {
		t.Type = types[0]
	}
	ok = true
	return
}

func (b *Builder) FuncRetDataType(toks Toks, i *int) (t models.RetType, ok bool) {
	defer func() { t.Type.Original = t.Type }()
	t.Type.Id = jntype.Void
	t.Type.Kind = jntype.VoidTypeStr
	if *i >= len(toks) {
		return
	}
	tok := toks[*i]
	if tok.Id == tokens.Brace && tok.Kind == tokens.LBRACKET {
		return b.funcMultiTypeRet(toks, i)
	}
	t.Type, ok = b.DataType(toks, i, false)
	return
}

func (b *Builder) pushStatementToBlock(bs *blockStatement) {
	if len(bs.toks) == 0 {
		return
	}
	lastTok := bs.toks[len(bs.toks)-1]
	if lastTok.Id == tokens.SemiColon {
		if len(bs.toks) == 1 {
			return
		}
		bs.toks = bs.toks[:len(bs.toks)-1]
	}
	s := b.Statement(bs)
	if s.Val == nil {
		return
	}
	s.WithTerminator = bs.withTerminator
	bs.block.Tree = append(bs.block.Tree, s)
}

func (b *Builder) Block(toks Toks) (block models.Block) {
	bs := new(blockStatement)
	bs.block = &block
	bs.srcToks = &toks
	for {
		bs.pos, bs.withTerminator = NextStatementPos(toks, 0)
		statementToks := toks[:bs.pos]
		bs.blockToks = &toks
		bs.toks = statementToks
		b.pushStatementToBlock(bs)
	next:
		if len(bs.nextToks) > 0 {
			bs.toks = bs.nextToks
			bs.nextToks = nil
			b.pushStatementToBlock(bs)
			goto next
		}
		if bs.pos >= len(toks) {
			break
		}
		toks = toks[bs.pos:]
	}
	return
}

func (b *Builder) Statement(bs *blockStatement) (s models.Statement) {
	s, ok := b.AssignStatement(bs.toks, false)
	if ok {
		return s
	}
	tok := bs.toks[0]
	switch tok.Id {
	case tokens.Id:
		s, ok := b.IdStatement(bs.toks)
		if ok {
			return s
		}
	case tokens.Const, tokens.Volatile:
		return b.VarStatement(bs.toks)
	case tokens.Ret:
		return b.RetStatement(bs.toks)
	case tokens.Iter:
		return b.IterExpr(bs.toks)
	case tokens.Break:
		return b.BreakStatement(bs.toks)
	case tokens.Continue:
		return b.ContinueStatement(bs.toks)
	case tokens.If:
		return b.IfExpr(bs)
	case tokens.Else:
		return b.ElseBlock(bs)
	case tokens.Comment:
		return b.CommentStatement(bs.toks[0])
	case tokens.Defer:
		return b.DeferStatement(bs.toks)
	case tokens.Co:
		return b.ConcurrentCallStatement(bs.toks)
	case tokens.Goto:
		return b.GotoStatement(bs.toks)
	case tokens.Try:
		return b.TryBlock(bs)
	case tokens.Catch:
		return b.CatchBlock(bs)
	case tokens.Type:
		t := b.Type(bs.toks)
		s.Tok = t.Tok
		s.Val = t
		return
	case tokens.Match:
		return b.MatchCase(bs.toks)
	case tokens.Brace:
		if tok.Kind == tokens.LBRACE {
			return b.blockStatement(bs.toks)
		}
	}
	if IsFuncCall(bs.toks) != nil {
		return b.ExprStatement(bs.toks)
	}
	tok = Tok{
		File:   tok.File,
		Id:     tokens.Ret,
		Kind:   tokens.RET,
		Row:    tok.Row,
		Column: tok.Column,
	}
	bs.toks = append([]Tok{tok}, bs.toks...)
	return b.RetStatement(bs.toks)
}

func (b *Builder) blockStatement(toks Toks) models.Statement {
	i := new(int)
	tok := toks[0]
	toks = Range(i, tokens.LBRACE, tokens.RBRACE, toks)
	if *i < len(toks) {
		b.pusherr(toks[*i], "invalid_syntax")
	}
	block := b.Block(toks)
	return models.Statement{Tok: tok, Val: block}
}

func (b *Builder) assignInfo(toks Toks) (info AssignInfo) {
	info.Ok = true
	braceCount := 0
	for i, tok := range toks {
		if tok.Id == tokens.Brace {
			switch tok.Kind {
			case tokens.LBRACE, tokens.LBRACKET, tokens.LPARENTHESES:
				braceCount++
			default:
				braceCount--
			}
		}
		if braceCount > 0 {
			continue
		} else if tok.Id != tokens.Operator {
			continue
		} else if !IsAssignOperator(tok.Kind) {
			continue
		}
		info.Left = toks[:i]
		if info.Left == nil {
			b.pusherr(tok, "invalid_syntax")
			info.Ok = false
		}
		info.Setter = tok
		if i+1 >= len(toks) {
			info.Right = nil
			info.Ok = IsSuffixOperator(info.Setter.Kind)
			break
		}
		info.Right = toks[i+1:]
		if IsSuffixOperator(info.Setter.Kind) {
			if info.Right != nil {
				b.pusherr(info.Right[0], "invalid_syntax")
				info.Right = nil
			}
		}
		break
	}
	return
}

func (b *Builder) pushAssignSelector(
	selectors *[]models.AssignLeft,
	last, current int,
	info AssignInfo,
) {
	var selector models.AssignLeft
	selector.Expr.Toks = info.Left[last:current]
	if last-current == 0 {
		b.pusherr(info.Left[current-1], "missing_expr")
		return
	}
	if selector.Expr.Toks[0].Id == tokens.Id &&
		current-last > 1 &&
		selector.Expr.Toks[1].Id == tokens.Colon {
		if info.IsExpr {
			b.pusherr(selector.Expr.Toks[0], "notallow_declares")
		}
		selector.Var.New = true
		selector.Var.IdTok = selector.Expr.Toks[0]
		selector.Var.Id = selector.Var.IdTok.Kind
		selector.Var.SetterTok = info.Setter
		if current-last > 2 {
			selector.Var.Type, _ = b.DataType(selector.Expr.Toks[2:], new(int), false)
		}
	} else {
		if selector.Expr.Toks[0].Id == tokens.Id {
			selector.Var.IdTok = selector.Expr.Toks[0]
			selector.Var.Id = selector.Var.IdTok.Kind
		}
		selector.Expr = b.Expr(selector.Expr.Toks)
	}
	*selectors = append(*selectors, selector)
}

func (b *Builder) assignSelectors(info AssignInfo) []models.AssignLeft {
	var selectors []models.AssignLeft
	braceCount := 0
	lastIndex := 0
	for i, tok := range info.Left {
		if tok.Id == tokens.Brace {
			switch tok.Kind {
			case tokens.LBRACE, tokens.LBRACKET, tokens.LPARENTHESES:
				braceCount++
			default:
				braceCount--
			}
		}
		if braceCount > 0 {
			continue
		} else if tok.Id != tokens.Comma {
			continue
		}
		b.pushAssignSelector(&selectors, lastIndex, i, info)
		lastIndex = i + 1
	}
	if lastIndex < len(info.Left) {
		b.pushAssignSelector(&selectors, lastIndex, len(info.Left), info)
	}
	return selectors
}

func (b *Builder) pushAssignExpr(exps *[]models.Expr, last, current int, info AssignInfo) {
	toks := info.Right[last:current]
	if toks == nil {
		b.pusherr(info.Right[current-1], "missing_expr")
		return
	}
	*exps = append(*exps, b.Expr(toks))
}

func (b *Builder) assignExprs(info AssignInfo) []models.Expr {
	var exprs []models.Expr
	braceCount := 0
	lastIndex := 0
	for i, tok := range info.Right {
		if tok.Id == tokens.Brace {
			switch tok.Kind {
			case tokens.LBRACE, tokens.LBRACKET, tokens.LPARENTHESES:
				braceCount++
			default:
				braceCount--
			}
		}
		if braceCount > 0 {
			continue
		} else if tok.Id != tokens.Comma {
			continue
		}
		b.pushAssignExpr(&exprs, lastIndex, i, info)
		lastIndex = i + 1
	}
	if lastIndex < len(info.Right) {
		b.pushAssignExpr(&exprs, lastIndex, len(info.Right), info)
	}
	return exprs
}

func (b *Builder) AssignStatement(toks Toks, isExpr bool) (s models.Statement, _ bool) {
	assign, ok := b.AssignExpr(toks, isExpr)
	if !ok {
		return
	}
	s.Tok = toks[0]
	s.Val = assign
	return s, true
}

func (b *Builder) AssignExpr(toks Toks, isExpr bool) (assign models.Assign, ok bool) {
	if !CheckAssignToks(toks) {
		return
	}
	info := b.assignInfo(toks)
	if !info.Ok {
		return
	}
	ok = true
	info.IsExpr = isExpr
	assign.IsExpr = isExpr
	assign.Setter = info.Setter
	assign.Left = b.assignSelectors(info)
	if isExpr && len(assign.Left) > 1 {
		b.pusherr(assign.Setter, "notallow_multiple_assign")
	}
	if info.Right != nil {
		assign.Right = b.assignExprs(info)
	}
	return
}

func (b *Builder) IdStatement(toks Toks) (s models.Statement, _ bool) {
	if len(toks) == 1 {
		return
	}
	tok := toks[1]
	switch tok.Id {
	case tokens.Colon:
		if len(toks) == 2 {
			return b.LabelStatement(toks[0]), true
		}
		return b.VarStatement(toks), true
	}
	return
}

func (b *Builder) LabelStatement(tok Tok) models.Statement {
	var l models.Label
	l.Tok = tok
	l.Label = tok.Kind
	return models.Statement{Tok: tok, Val: l}
}

func (b *Builder) ExprStatement(toks Toks) models.Statement {
	block := models.ExprStatement{
		Expr: b.Expr(toks),
	}
	return models.Statement{
		Tok: toks[0],
		Val: block,
	}
}

func (b *Builder) Args(toks Toks) *models.Args {
	args := new(models.Args)
	last := 0
	braceCount := 0
	for i, tok := range toks {
		if tok.Id == tokens.Brace {
			switch tok.Kind {
			case tokens.LBRACE, tokens.LBRACKET, tokens.LPARENTHESES:
				braceCount++
			default:
				braceCount--
			}
		}
		if braceCount > 0 || tok.Id != tokens.Comma {
			continue
		}
		b.pushArg(args, toks[last:i], tok)
		last = i + 1
	}
	if last < len(toks) {
		if last == 0 {
			b.pushArg(args, toks[last:], toks[last])
		} else {
			b.pushArg(args, toks[last:], toks[last-1])
		}
	}
	return args
}

func (b *Builder) pushArg(args *models.Args, toks Toks, err Tok) {
	if len(toks) == 0 {
		b.pusherr(err, "invalid_syntax")
		return
	}
	var arg models.Arg
	arg.Tok = toks[0]
	if arg.Tok.Id == tokens.Id {
		if len(toks) > 1 {
			tok := toks[1]
			if tok.Id == tokens.Operator && tok.Kind == tokens.EQUAL {
				args.Targeted = true
				arg.TargetId = arg.Tok.Kind
				toks = toks[2:]
			}
		}
	}
	arg.Expr = b.Expr(toks)
	args.Src = append(args.Src, arg)
}

func (b *Builder) varBegin(v *models.Var, i *int, toks Toks) {
	for ; *i < len(toks); *i++ {
		tok := toks[*i]
		if tok.Id == tokens.Id {
			break
		}
		switch tok.Id {
		case tokens.Const:
			if v.Const {
				b.pusherr(tok, "already_constant")
				break
			}
			v.Const = true
		case tokens.Volatile:
			if v.Volatile {
				b.pusherr(tok, "already_volatile")
				break
			}
			v.Volatile = true
		default:
			b.pusherr(tok, "invalid_syntax")
		}
	}
}

func (b *Builder) varTypeNExpr(v *models.Var, toks Toks, i int) {
	tok := toks[i]
	t, ok := b.DataType(toks, &i, false)
	if ok {
		v.Type = t
		i++
		if i >= len(toks) {
			return
		}
		tok = toks[i]
	}
	if tok.Id == tokens.Operator {
		if tok.Kind != tokens.EQUAL {
			b.pusherr(tok, "invalid_syntax")
			return
		}
		valueToks := toks[i+1:]
		if len(valueToks) == 0 {
			b.pusherr(tok, "missing_expr")
			return
		}
		v.Val = b.Expr(valueToks)
		v.SetterTok = tok
	} else {
		b.pusherr(tok, "invalid_syntax")
	}
}

func (b *Builder) Var(toks Toks) (v models.Var) {
	v.Pub = b.pub
	b.pub = false
	i := 0
	v.DefTok = toks[i]
	b.varBegin(&v, &i, toks)
	if i >= len(toks) {
		return
	}
	v.IdTok = toks[i]
	if v.IdTok.Id != tokens.Id {
		b.pusherr(v.IdTok, "invalid_syntax")
	}
	v.Id = v.IdTok.Kind
	v.Type.Id = jntype.Void
	v.Type.Kind = jntype.VoidTypeStr
	v.Type.Kind = jntype.VoidTypeStr
	i++
	if i >= len(toks) {
		b.pusherr(toks[i-1], "invalid_syntax")
		return
	}
	if v.DefTok.File != nil {
		if toks[i].Id != tokens.Colon {
			b.pusherr(toks[i], "invalid_syntax")
			return
		}
		i++
	}
	if i < len(toks) {
		b.varTypeNExpr(&v, toks, i)
	}
	return
}

func (b *Builder) VarStatement(toks Toks) models.Statement {
	vast := b.Var(toks)
	return models.Statement{
		Tok: vast.IdTok,
		Val: vast,
	}
}

func (b *Builder) CommentStatement(tok Tok) (s models.Statement) {
	s.Tok = tok
	tok.Kind = strings.TrimSpace(tok.Kind[2:])
	if strings.HasPrefix(tok.Kind, "cxx:") {
		s.Val = models.CxxEmbed{
			Tok:     tok,
			Content: tok.Kind[4:],
		}
	} else {
		s.Val = models.Comment{
			Content: tok.Kind,
		}
	}
	return
}

func (b *Builder) DeferStatement(toks Toks) (s models.Statement) {
	var d models.Defer
	d.Tok = toks[0]
	toks = toks[1:]
	if len(toks) == 0 {
		b.pusherr(d.Tok, "missing_expr")
		return
	}
	if IsFuncCall(toks) == nil {
		b.pusherr(d.Tok, "expr_not_func_call")
	}
	d.Expr = b.Expr(toks)
	s.Tok = d.Tok
	s.Val = d
	return
}

func (b *Builder) ConcurrentCallStatement(toks Toks) (s models.Statement) {
	var cc models.ConcurrentCall
	cc.Tok = toks[0]
	toks = toks[1:]
	if len(toks) == 0 {
		b.pusherr(cc.Tok, "missing_expr")
		return
	}
	if IsFuncCall(toks) == nil {
		b.pusherr(cc.Tok, "expr_not_func_call")
	}
	cc.Expr = b.Expr(toks)
	s.Tok = cc.Tok
	s.Val = cc
	return
}

func (b *Builder) GotoStatement(toks Toks) (s models.Statement) {
	s.Tok = toks[0]
	if len(toks) == 1 {
		b.pusherr(s.Tok, "missing_goto_label")
		return
	} else if len(toks) > 2 {
		b.pusherr(toks[2], "invalid_syntax")
	}
	idTok := toks[1]
	if idTok.Id != tokens.Id {
		b.pusherr(idTok, "invalid_syntax")
		return
	}
	var gt models.Goto
	gt.Tok = s.Tok
	gt.Label = idTok.Kind
	s.Val = gt
	return
}

func (b *Builder) RetStatement(toks Toks) models.Statement {
	var ret models.Ret
	ret.Tok = toks[0]
	if len(toks) > 1 {
		ret.Expr = b.Expr(toks[1:])
	}
	return models.Statement{
		Tok: ret.Tok,
		Val: ret,
	}
}

func (b *Builder) getWhileIterProfile(toks Toks) models.IterWhile {
	return models.IterWhile{
		Expr: b.Expr(toks),
	}
}

func (b *Builder) getForeachVarsToks(toks Toks) []Toks {
	vars, errs := Parts(toks, tokens.Comma)
	b.Errors = append(b.Errors, errs...)
	return vars
}

func (b *Builder) getVarProfile(toks Toks) (vast models.Var) {
	if len(toks) == 0 {
		return
	}
	vast.IdTok = toks[0]
	if vast.IdTok.Id != tokens.Id {
		b.pusherr(vast.IdTok, "invalid_syntax")
		return
	}
	vast.Id = vast.IdTok.Kind
	if len(toks) == 1 {
		return
	}
	if colon := toks[1]; colon.Id != tokens.Colon {
		b.pusherr(colon, "invalid_syntax")
		return
	}
	vast.New = true
	i := new(int)
	*i = 2
	if *i >= len(toks) {
		return
	}
	vast.Type, _ = b.DataType(toks, i, true)
	if *i < len(toks)-1 {
		b.pusherr(toks[*i], "invalid_syntax")
	}
	return
}

func (b *Builder) getForeachIterVars(varsToks []Toks) []models.Var {
	var vars []models.Var
	for _, toks := range varsToks {
		vars = append(vars, b.getVarProfile(toks))
	}
	return vars
}

func (b *Builder) getForeachIterProfile(varToks, exprToks Toks, inTok Tok) models.IterForeach {
	var foreach models.IterForeach
	foreach.InTok = inTok
	foreach.Expr = b.Expr(exprToks)
	if len(varToks) == 0 {
		foreach.KeyA.Id = jnapi.Ignore
		foreach.KeyB.Id = jnapi.Ignore
	} else {
		varsToks := b.getForeachVarsToks(varToks)
		if len(varsToks) == 0 {
			return foreach
		}
		if len(varsToks) > 2 {
			b.pusherr(inTok, "much_foreach_vars")
		}
		vars := b.getForeachIterVars(varsToks)
		foreach.KeyA = vars[0]
		if len(vars) > 1 {
			foreach.KeyB = vars[1]
		} else {
			foreach.KeyB.Id = jnapi.Ignore
		}
	}
	return foreach
}

func (b *Builder) getIterProfile(toks Toks) models.IterProfile {
	braceCount := 0
	for i, tok := range toks {
		if tok.Id == tokens.Brace {
			switch tok.Kind {
			case tokens.LBRACE, tokens.LBRACKET, tokens.LPARENTHESES:
				braceCount++
			default:
				braceCount--
			}
		}
		if braceCount != 0 {
			continue
		}
		if tok.Id == tokens.In {
			varToks := toks[:i]
			exprToks := toks[i+1:]
			return b.getForeachIterProfile(varToks, exprToks, tok)
		}
	}
	return b.getWhileIterProfile(toks)
}

func (b *Builder) IterExpr(toks Toks) (s models.Statement) {
	var iter models.Iter
	iter.Tok = toks[0]
	toks = toks[1:]
	if len(toks) == 0 {
		b.pusherr(iter.Tok, "body_not_exist")
		return
	}
	exprToks := BlockExpr(toks)
	if len(exprToks) > 0 {
		iter.Profile = b.getIterProfile(exprToks)
	}
	i := new(int)
	*i = len(exprToks)
	blockToks := b.getrange(i, tokens.LBRACE, tokens.RBRACE, &toks)
	if blockToks == nil {
		b.pusherr(iter.Tok, "body_not_exist")
		return
	}
	if *i < len(toks) {
		b.pusherr(toks[*i], "invalid_syntax")
	}
	iter.Block = b.Block(blockToks)
	return models.Statement{
		Tok: iter.Tok,
		Val: iter,
	}
}

func (b *Builder) TryBlock(bs *blockStatement) (s models.Statement) {
	var try models.Try
	try.Tok = bs.toks[0]
	bs.toks = bs.toks[1:]
	i := new(int)
	blockToks := b.getrange(i, tokens.LBRACE, tokens.RBRACE, &bs.toks)
	if blockToks == nil {
		b.pusherr(try.Tok, "body_not_exist")
		return
	}
	if *i < len(bs.toks) {
		if bs.toks[*i].Id == tokens.Catch {
			bs.nextToks = bs.toks[*i:]
		} else {
			b.pusherr(bs.toks[*i], "invalid_syntax")
		}
	}
	try.Block = b.Block(blockToks)
	return models.Statement{
		Tok: try.Tok,
		Val: try,
	}
}

func (b *Builder) CatchBlock(bs *blockStatement) (s models.Statement) {
	var catch models.Catch
	catch.Tok = bs.toks[0]
	bs.toks = bs.toks[1:]
	varToks := BlockExpr(bs.toks)
	i := new(int)
	*i = len(varToks)
	blockToks := b.getrange(i, tokens.LBRACE, tokens.RBRACE, &bs.toks)
	if blockToks == nil {
		b.pusherr(catch.Tok, "body_not_exist")
		return
	}
	if *i < len(bs.toks) {
		if bs.toks[*i].Id == tokens.Catch {
			bs.nextToks = bs.toks[*i:]
		} else {
			b.pusherr(bs.toks[*i], "invalid_syntax")
		}
	}
	if len(varToks) > 0 {
		catch.Var = b.getVarProfile(varToks)
	}
	catch.Block = b.Block(blockToks)
	return models.Statement{
		Tok: catch.Tok,
		Val: catch,
	}
}

func (b *Builder) caseexprs(toks *Toks, caseIsDefault bool) []models.Expr {
	var exprs []models.Expr
	pushExpr := func(toks Toks, tok Tok) {
		if caseIsDefault {
			if len(toks) > 0 {
				b.pusherr(tok, "invalid_syntax")
			}
			return
		}
		if len(toks) > 0 {
			exprs = append(exprs, b.Expr(toks))
			return
		}
		b.pusherr(tok, "missing_expr")
	}
	braceCount := 0
	j := 0
	var i int
	var tok Tok
	for i, tok = range *toks {
		if tok.Id == tokens.Brace {
			switch tok.Kind {
			case tokens.LPARENTHESES, tokens.LBRACE, tokens.LBRACKET:
				braceCount++
			default:
				braceCount--
			}
			continue
		} else if braceCount != 0 {
			continue
		}
		switch tok.Id {
		case tokens.Comma:
			pushExpr((*toks)[j:i], tok)
			j = i + 1
		case tokens.Colon:
			pushExpr((*toks)[j:i], tok)
			*toks = (*toks)[i+1:]
			return exprs
		}
	}
	b.pusherr((*toks)[0], "invalid_syntax")
	*toks = nil
	return nil
}

func (b *Builder) caseblock(toks *Toks) models.Block {
	braceCount := 0
	for i, tok := range *toks {
		if tok.Id == tokens.Brace {
			switch tok.Kind {
			case tokens.LPARENTHESES, tokens.LBRACE, tokens.LBRACKET:
				braceCount++
			default:
				braceCount--
			}
			continue
		} else if braceCount != 0 {
			continue
		}
		switch tok.Id {
		case tokens.Case, tokens.Default:
			blockToks := (*toks)[:i]
			*toks = (*toks)[i:]
			return b.Block(blockToks)
		}
	}
	block := b.Block(*toks)
	*toks = nil
	return block
}

func (b *Builder) getcase(toks *Toks) models.Case {
	var c models.Case
	tok := (*toks)[0]
	*toks = (*toks)[1:]
	c.Exprs = b.caseexprs(toks, tok.Id == tokens.Default)
	c.Block = b.caseblock(toks)
	return c
}

func (b *Builder) cases(toks Toks) ([]models.Case, *models.Case) {
	var cases []models.Case
	var def *models.Case
	for len(toks) > 0 {
		tok := toks[0]
		switch tok.Id {
		case tokens.Case:
			cases = append(cases, b.getcase(&toks))
		case tokens.Default:
			c := b.getcase(&toks)
			if def == nil {
				def = new(models.Case)
				*def = c
				break
			}
			fallthrough
		default:
			b.pusherr(tok, "invalid_syntax")
		}
	}
	return cases, def
}

func (b *Builder) MatchCase(toks Toks) (s models.Statement) {
	var match models.Match
	match.Tok = toks[0]
	s.Tok = match.Tok
	toks = toks[1:]
	exprToks := BlockExpr(toks)
	if len(exprToks) > 0 {
		match.Expr = b.Expr(exprToks)
	}
	i := new(int)
	*i = len(exprToks)
	blockToks := b.getrange(i, tokens.LBRACE, tokens.RBRACE, &toks)
	if blockToks == nil {
		b.pusherr(match.Tok, "body_not_exist")
		return
	}
	match.Cases, match.Default = b.cases(blockToks)
	s.Val = match
	return
}

func (b *Builder) IfExpr(bs *blockStatement) (s models.Statement) {
	var ifast models.If
	ifast.Tok = bs.toks[0]
	bs.toks = bs.toks[1:]
	exprToks := BlockExpr(bs.toks)
	i := new(int)
	if len(exprToks) == 0 {
		if len(bs.toks) == 0 || bs.pos >= len(*bs.srcToks) {
			b.pusherr(ifast.Tok, "missing_expr")
			return
		}
		exprToks = bs.toks
		*bs.srcToks = (*bs.srcToks)[bs.pos:]
		bs.pos, bs.withTerminator = NextStatementPos(*bs.srcToks, 0)
		bs.toks = (*bs.srcToks)[:bs.pos]
	} else {
		*i = len(exprToks)
	}
	blockToks := b.getrange(i, tokens.LBRACE, tokens.RBRACE, &bs.toks)
	if blockToks == nil {
		b.pusherr(ifast.Tok, "body_not_exist")
		return
	}
	if *i < len(bs.toks) {
		if bs.toks[*i].Id == tokens.Else {
			bs.nextToks = bs.toks[*i:]
		} else {
			b.pusherr(bs.toks[*i], "invalid_syntax")
		}
	}
	ifast.Expr = b.Expr(exprToks)
	ifast.Block = b.Block(blockToks)
	return models.Statement{
		Tok: ifast.Tok,
		Val: ifast,
	}
}

func (b *Builder) ElseIfExpr(bs *blockStatement) (s models.Statement) {
	var elif models.ElseIf
	elif.Tok = bs.toks[1]
	bs.toks = bs.toks[2:]
	exprToks := BlockExpr(bs.toks)
	i := new(int)
	if len(exprToks) == 0 {
		if len(bs.toks) == 0 || bs.pos >= len(*bs.srcToks) {
			b.pusherr(elif.Tok, "missing_expr")
			return
		}
		exprToks = bs.toks
		*bs.srcToks = (*bs.srcToks)[bs.pos:]
		bs.pos, bs.withTerminator = NextStatementPos(*bs.srcToks, 0)
		bs.toks = (*bs.srcToks)[:bs.pos]
	} else {
		*i = len(exprToks)
	}
	blockToks := b.getrange(i, tokens.LBRACE, tokens.RBRACE, &bs.toks)
	if blockToks == nil {
		b.pusherr(elif.Tok, "body_not_exist")
		return
	}
	if *i < len(bs.toks) {
		if bs.toks[*i].Id == tokens.Else {
			bs.nextToks = bs.toks[*i:]
		} else {
			b.pusherr(bs.toks[*i], "invalid_syntax")
		}
	}
	elif.Expr = b.Expr(exprToks)
	elif.Block = b.Block(blockToks)
	return models.Statement{
		Tok: elif.Tok,
		Val: elif,
	}
}

func (b *Builder) ElseBlock(bs *blockStatement) (s models.Statement) {
	if len(bs.toks) > 1 && bs.toks[1].Id == tokens.If {
		return b.ElseIfExpr(bs)
	}
	var elseast models.Else
	elseast.Tok = bs.toks[0]
	bs.toks = bs.toks[1:]
	i := new(int)
	blockToks := b.getrange(i, tokens.LBRACE, tokens.RBRACE, &bs.toks)
	if blockToks == nil {
		if *i < len(bs.toks) {
			b.pusherr(elseast.Tok, "else_have_expr")
		} else {
			b.pusherr(elseast.Tok, "body_not_exist")
		}
		return
	}
	if *i < len(bs.toks) {
		b.pusherr(bs.toks[*i], "invalid_syntax")
	}
	elseast.Block = b.Block(blockToks)
	return models.Statement{
		Tok: elseast.Tok,
		Val: elseast,
	}
}

func (b *Builder) BreakStatement(toks Toks) models.Statement {
	var breakAST models.Break
	breakAST.Tok = toks[0]
	if len(toks) > 1 {
		b.pusherr(toks[1], "invalid_syntax")
	}
	return models.Statement{
		Tok: breakAST.Tok,
		Val: breakAST,
	}
}

func (b *Builder) ContinueStatement(toks Toks) models.Statement {
	var continueAST models.Continue
	continueAST.Tok = toks[0]
	if len(toks) > 1 {
		b.pusherr(toks[1], "invalid_syntax")
	}
	return models.Statement{
		Tok: continueAST.Tok,
		Val: continueAST,
	}
}

func (b *Builder) Expr(toks Toks) (e models.Expr) {
	e.Processes = b.exprProcesses(toks)
	e.Toks = toks
	return
}

type exprProcessInfo struct {
	processes        []Toks
	part             Toks
	operator         bool
	value            bool
	singleOperatored bool
	pushedError      bool
	braceCount       int
	toks             Toks
	i                int
}

func (b *Builder) exprOperatorPart(info *exprProcessInfo, tok Tok) {
	if IsExpressionOperator(tok.Kind) ||
		IsAssignOperator(tok.Kind) {
		info.part = append(info.part, tok)
		return
	}
	if !info.operator {
		if IsUnaryOperator(tok.Kind) && !info.singleOperatored {
			info.part = append(info.part, tok)
			info.singleOperatored = true
			return
		}
		if info.braceCount == 0 && IsSolidOperator(tok.Kind) {
			b.pusherr(tok, "operator_overflow")
		}
	}
	info.singleOperatored = false
	info.operator = false
	info.value = true
	if info.braceCount > 0 {
		info.part = append(info.part, tok)
		return
	}
	info.processes = append(info.processes, info.part)
	info.processes = append(info.processes, Toks{tok})
	info.part = Toks{}
}

func (b *Builder) exprValuePart(info *exprProcessInfo, tok Tok) {
	if info.i > 0 && info.braceCount == 0 {
		lt := info.toks[info.i-1]
		if (lt.Id == tokens.Id || lt.Id == tokens.Value) &&
			(tok.Id == tokens.Id || tok.Id == tokens.Value) {
			b.pusherr(tok, "invalid_syntax")
			info.pushedError = true
		}
	}
	b.checkExprTok(tok)
	info.part = append(info.part, tok)
	info.operator = RequireOperatorToProcess(tok, info.i, len(info.toks))
	info.value = false
}

func (b *Builder) exprBracePart(info *exprProcessInfo, tok Tok) bool {
	switch tok.Kind {
	case tokens.LBRACE, tokens.LBRACKET, tokens.LPARENTHESES:
		if tok.Kind == tokens.LBRACKET {
			oldIndex := info.i
			_, ok := b.DataType(info.toks, &info.i, false)
			if ok {
				info.part = append(info.part, info.toks[oldIndex:info.i+1]...)
				return true
			}
			info.i = oldIndex
		}
		info.singleOperatored = false
		info.braceCount++
	default:
		info.braceCount--
	}
	return false
}

func (b *Builder) exprProcesses(toks Toks) []Toks {
	var info exprProcessInfo
	info.toks = toks
	for ; info.i < len(info.toks); info.i++ {
		tok := info.toks[info.i]
		switch tok.Id {
		case tokens.Operator:
			b.exprOperatorPart(&info, tok)
			continue
		case tokens.Brace:
			skipStep := b.exprBracePart(&info, tok)
			if skipStep {
				continue
			}
		}
		b.exprValuePart(&info, tok)
	}
	if len(info.part) > 0 {
		info.processes = append(info.processes, info.part)
	}
	if info.value {
		b.pusherr(info.processes[len(info.processes)-1][0], "operator_overflow")
		info.pushedError = true
	}
	if info.pushedError {
		return nil
	}
	return info.processes
}

func (b *Builder) checkExprTok(tok Tok) {
	if lexer.NumRegexp.MatchString(tok.Kind) {
		var result bool
		if strings.Contains(tok.Kind, tokens.DOT) ||
			(!strings.HasPrefix(tok.Kind, "0x") && strings.ContainsAny(tok.Kind, "eE")) {
			result = jnbits.CheckBitFloat(tok.Kind, 64)
		} else {
			result = jnbits.CheckBitInt(tok.Kind, jnbits.MaxInt)
			if !result {
				result = jnbits.CheckBitUInt(tok.Kind, jnbits.MaxInt)
			}
		}
		if !result {
			b.pusherr(tok, "invalid_numeric_range")
		}
	}
}

func (b *Builder) getrange(i *int, open, close string, toks *Toks) Toks {
	rang := Range(i, open, close, *toks)
	if rang != nil {
		return rang
	}
	if b.Ended() {
		return nil
	}
	*i = 0
	*toks = b.nextBuilderStatement()
	rang = Range(i, open, close, *toks)
	return rang
}

func (b *Builder) skipStatement(i *int, toks *Toks) Toks {
	start := *i
	*i, _ = NextStatementPos(*toks, start)
	stoks := (*toks)[start:*i]
	if stoks[len(stoks)-1].Id == tokens.SemiColon {
		if len(stoks) == 1 {
			return b.skipStatement(i, toks)
		}
		stoks = stoks[:len(stoks)-1]
	}
	return stoks
}

func (b *Builder) nextBuilderStatement() Toks {
	return b.skipStatement(&b.Pos, &b.Toks)
}
