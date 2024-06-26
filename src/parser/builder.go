// Copyright (c) 2024 arfy slowy - DeRuneLabs
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package parser

import (
	"os"
	"strings"

	"github.com/DeRuneLabs/jane"
	"github.com/DeRuneLabs/jane/ast"
	"github.com/DeRuneLabs/jane/build"
	"github.com/DeRuneLabs/jane/lexer"
	"github.com/DeRuneLabs/jane/types"
)

func compilerErr(t lexer.Token, key string, args ...any) build.Log {
	return build.Log{
		Type:   build.ERR,
		Row:    t.Row,
		Column: t.Column,
		Path:   t.File.Path(),
		Text:   build.Errorf(key, args...),
	}
}

type block_st struct {
	pos        int
	block      *ast.Block
	srcToks    *[]lexer.Token
	toks       []lexer.Token
	nextToks   []lexer.Token
	terminated bool
}

type builder struct {
	public bool
	Tree   []ast.Node
	Errors []build.Log
	Tokens []lexer.Token
	Pos    int
}

func new_builder(t []lexer.Token) *builder {
	b := new(builder)
	b.Tokens = t
	b.Pos = 0
	return b
}

func (b *builder) pusherr(t lexer.Token, key string, args ...any) {
	b.Errors = append(b.Errors, compilerErr(t, key, args...))
}

func (b *builder) Ended() bool {
	return b.Pos >= len(b.Tokens)
}

func (b *builder) buildNode(toks []lexer.Token) {
	t := toks[0]
	switch t.Id {
	case lexer.ID_USE:
		b.Use(toks)
	case lexer.ID_FN, lexer.ID_UNSAFE:
		s := ast.St{Token: t}
		s.Data = b.Func(toks, false, false, false)
		b.Tree = append(b.Tree, ast.Node{Token: s.Token, Data: s})
	case lexer.ID_CONST, lexer.ID_LET, lexer.ID_MUT:
		b.GlobalVar(toks)
	case lexer.ID_TYPE:
		b.Tree = append(b.Tree, b.GlobalTypeAlias(toks))
	case lexer.ID_ENUM:
		b.Enum(toks)
	case lexer.ID_STRUCT:
		b.Structure(toks)
	case lexer.ID_TRAIT:
		b.Trait(toks)
	case lexer.ID_IMPL:
		b.Impl(toks)
	case lexer.ID_CPP:
		b.CppLink(toks)
	case lexer.ID_COMMENT:
		b.Tree = append(b.Tree, b.Comment(toks[0]))
	default:
		b.pusherr(t, "invalid_syntax")
		return
	}
	if b.public {
		b.pusherr(t, "def_not_support_pub")
	}
}

func (b *builder) Build() {
	for b.Pos != -1 && !b.Ended() {
		toks := b.next_builder_st()
		b.public = toks[0].Id == lexer.ID_PUB
		if b.public {
			if len(toks) == 1 {
				if b.Ended() {
					b.pusherr(toks[0], "invalid_syntax")
					continue
				}
				toks = b.next_builder_st()
			} else {
				toks = toks[1:]
			}
		}
		b.buildNode(toks)
	}
}

func (b *builder) TypeAlias(toks []lexer.Token) (t ast.TypeAlias) {
	i := 1
	if i >= len(toks) {
		b.pusherr(toks[i-1], "invalid_syntax")
		return
	}
	t.Token = toks[1]
	t.Id = t.Token.Kind
	token := toks[i]
	if token.Id != lexer.ID_IDENT {
		b.pusherr(token, "invalid_syntax")
	}
	i++
	if i >= len(toks) {
		b.pusherr(toks[i-1], "invalid_syntax")
		return
	}
	token = toks[i]
	if token.Id != lexer.ID_COLON {
		b.pusherr(toks[i-1], "invalid_syntax")
		return
	}
	i++
	if i >= len(toks) {
		b.pusherr(toks[i-1], "missing_type")
		return
	}
	destType, ok := b.DataType(toks, &i, true)
	t.TargetType = destType
	if ok && i+1 < len(toks) {
		b.pusherr(toks[i+1], "invalid_syntax")
	}
	return
}

func (b *builder) buildEnumItemExpr(i *int, toks []lexer.Token) ast.Expr {
	brace_n := 0
	exprStart := *i
	for ; *i < len(toks); *i++ {
		t := toks[*i]
		if t.Id == lexer.ID_BRACE {
			switch t.Kind {
			case lexer.KND_LBRACE, lexer.KND_LBRACKET, lexer.KND_LPAREN:
				brace_n++
				continue
			default:
				brace_n--
			}
		}
		if brace_n > 0 {
			continue
		}
		if t.Id == lexer.ID_COMMA || *i+1 >= len(toks) {
			var exprToks []lexer.Token
			if t.Id == lexer.ID_COMMA {
				exprToks = toks[exprStart:*i]
			} else {
				exprToks = toks[exprStart:]
			}
			return b.Expr(exprToks)
		}
	}
	return ast.Expr{}
}

func (b *builder) buildEnumItems(toks []lexer.Token) []*ast.EnumItem {
	items := make([]*ast.EnumItem, 0)
	for i := 0; i < len(toks); i++ {
		t := toks[i]
		if t.Id == lexer.ID_COMMENT {
			continue
		}
		item := new(ast.EnumItem)
		item.Token = t
		if item.Token.Id != lexer.ID_IDENT {
			b.pusherr(item.Token, "invalid_syntax")
		}
		item.Id = item.Token.Kind
		if i+1 >= len(toks) || toks[i+1].Id == lexer.ID_COMMA {
			if i+1 < len(toks) {
				i++
			}
			items = append(items, item)
			continue
		}
		i++
		t = toks[i]
		if t.Id != lexer.ID_OP && t.Kind != lexer.KND_EQ {
			b.pusherr(toks[0], "invalid_syntax")
		}
		i++
		if i >= len(toks) || toks[i].Id == lexer.ID_COMMA {
			b.pusherr(toks[0], "missing_expr")
			continue
		}
		item.Expr = b.buildEnumItemExpr(&i, toks)
		items = append(items, item)
	}
	return items
}

func (b *builder) Enum(toks []lexer.Token) {
	if len(toks) < 2 || len(toks) < 3 {
		b.pusherr(toks[0], "invalid_syntax")
		return
	}
	e := &ast.Enum{}
	e.Token = toks[1]
	if e.Token.Id != lexer.ID_IDENT {
		b.pusherr(e.Token, "invalid_syntax")
	}
	e.Id = e.Token.Kind
	i := 2
	if toks[i].Id == lexer.ID_COLON {
		i++
		if i >= len(toks) {
			b.pusherr(toks[i-1], "invalid_syntax")
			return
		}
		e.DataType, _ = b.DataType(toks, &i, true)
		i++
		if i >= len(toks) {
			b.stop()
			b.pusherr(e.Token, "body_not_exist")
			return
		}
	} else {
		e.DataType = ast.Type{Id: types.U32, Kind: types.TYPE_MAP[types.U32]}
	}
	itemToks := b.getrange(&i, lexer.KND_LBRACE, lexer.KND_RBRACE, &toks)
	if itemToks == nil {
		b.stop()
		b.pusherr(e.Token, "body_not_exist")
		return
	} else if i < len(toks) {
		b.pusherr(toks[i], "invalid_syntax")
	}
	e.Pub = b.public
	b.public = false
	e.Items = b.buildEnumItems(itemToks)
	b.Tree = append(b.Tree, ast.Node{Token: e.Token, Data: e})
}

func (b *builder) Comment(t lexer.Token) ast.Node {
	t.Kind = strings.TrimSpace(t.Kind[2:])
	return ast.Node{
		Token: t,
		Data: ast.Comment{
			Token:   t,
			Content: t.Kind,
		},
	}
}

func (b *builder) structFields(toks []lexer.Token, cpp_linked bool) []*ast.Var {
	var fields []*ast.Var
	i := 0
	for i < len(toks) {
		var_tokens := b.skip_st(&i, &toks)
		if var_tokens[0].Id == lexer.ID_COMMENT {
			continue
		}
		is_pub := var_tokens[0].Id == lexer.ID_PUB
		if is_pub {
			if len(var_tokens) == 1 {
				b.pusherr(var_tokens[0], "invalid_syntax")
				continue
			}
			var_tokens = var_tokens[1:]
		}
		is_mut := var_tokens[0].Id == lexer.ID_MUT
		if is_mut {
			if len(var_tokens) == 1 {
				b.pusherr(var_tokens[0], "invalid_syntax")
				continue
			}
			var_tokens = var_tokens[1:]
		}
		v := b.Var(var_tokens, false, false)
		v.Public = is_pub
		v.Mutable = is_mut
		v.IsField = true
		v.CppLinked = cpp_linked
		fields = append(fields, &v)
	}
	return fields
}

func (b *builder) build_struct(toks []lexer.Token, cpp_linked bool) ast.Struct {
	var s ast.Struct
	s.Pub = b.public
	b.public = false
	if len(toks) < 3 {
		b.pusherr(toks[0], "invalid_syntax")
		return s
	}

	i := 1
	s.Token = toks[i]
	if s.Token.Id != lexer.ID_IDENT {
		b.pusherr(s.Token, "invalid_syntax")
	}
	i++
	if i >= len(toks) {
		b.pusherr(toks[i], "invalid_syntax")
		return s
	}
	s.Id = s.Token.Kind

	generics_toks := ast.Range(&i, lexer.KND_LBRACKET, lexer.KND_RBRACKET, toks)
	if generics_toks != nil {
		s.Generics = b.Generics(generics_toks)
	}
	if i >= len(toks) {
		b.pusherr(toks[i], "invalid_syntax")
		return s
	}

	body_toks := b.getrange(&i, lexer.KND_LBRACE, lexer.KND_RBRACE, &toks)
	if body_toks == nil {
		b.stop()
		b.pusherr(s.Token, "body_not_exist")
		return s
	}
	if i < len(toks) {
		b.pusherr(toks[i], "invalid_syntax")
	}
	s.Fields = b.structFields(body_toks, cpp_linked)
	return s
}

func (b *builder) Structure(toks []lexer.Token) {
	s := b.build_struct(toks, false)
	b.Tree = append(b.Tree, ast.Node{Token: s.Token, Data: s})
}

func (b *builder) traitFns(toks []lexer.Token, trait_id string) []*ast.Fn {
	var fns []*ast.Fn
	i := 0
	for i < len(toks) {
		fnToks := b.skip_st(&i, &toks)
		f := b.Func(fnToks, true, false, true)
		b.setup_receiver(&f, trait_id)
		f.Public = true
		fns = append(fns, &f)
	}
	return fns
}

func (b *builder) Trait(toks []lexer.Token) {
	var t ast.Trait
	t.Pub = b.public
	b.public = false
	if len(toks) < 3 {
		b.pusherr(toks[0], "invalid_syntax")
		return
	}
	t.Token = toks[1]
	if t.Token.Id != lexer.ID_IDENT {
		b.pusherr(t.Token, "invalid_syntax")
	}
	t.Id = t.Token.Kind
	i := 2
	bodyToks := b.getrange(&i, lexer.KND_LBRACE, lexer.KND_RBRACE, &toks)
	if bodyToks == nil {
		b.stop()
		b.pusherr(t.Token, "body_not_exist")
		return
	}
	if i < len(toks) {
		b.pusherr(toks[i], "invalid_syntax")
	}
	t.Funcs = b.traitFns(bodyToks, t.Id)
	b.Tree = append(b.Tree, ast.Node{Token: t.Token, Data: t})
}

func (b *builder) implTraitFns(impl *ast.Impl, toks []lexer.Token) {
	pos, btoks := b.Pos, make([]lexer.Token, len(b.Tokens))
	copy(btoks, b.Tokens)
	b.Pos = 0
	b.Tokens = toks
	for b.Pos != -1 && !b.Ended() {
		fnToks := b.next_builder_st()
		tok := fnToks[0]
		switch tok.Id {
		case lexer.ID_COMMENT:
			impl.Tree = append(impl.Tree, b.Comment(tok))
			continue
		case lexer.ID_FN, lexer.ID_UNSAFE:
			f := b.get_method(fnToks)
			f.Public = true
			b.setup_receiver(f, impl.Target.Kind)
			impl.Tree = append(impl.Tree, ast.Node{Token: f.Token, Data: f})
		default:
			b.pusherr(tok, "invalid_syntax")
			continue
		}
	}
	b.Pos, b.Tokens = pos, btoks
}

func (b *builder) implStruct(impl *ast.Impl, toks []lexer.Token) {
	pos, btoks := b.Pos, make([]lexer.Token, len(b.Tokens))
	copy(btoks, b.Tokens)
	b.Pos = 0
	b.Tokens = toks
	for b.Pos != -1 && !b.Ended() {
		fnToks := b.next_builder_st()
		tok := fnToks[0]
		pub := false
		switch tok.Id {
		case lexer.ID_COMMENT:
			impl.Tree = append(impl.Tree, b.Comment(tok))
			continue
		}
		if tok.Id == lexer.ID_PUB {
			pub = true
			if len(fnToks) == 1 {
				b.pusherr(fnToks[0], "invalid_syntax")
				continue
			}
			fnToks = fnToks[1:]
			if len(fnToks) > 0 {
				tok = fnToks[0]
			}
		}
		switch tok.Id {
		case lexer.ID_FN, lexer.ID_UNSAFE:
			f := b.get_method(fnToks)
			f.Public = pub
			b.setup_receiver(f, impl.Base.Kind)
			impl.Tree = append(impl.Tree, ast.Node{Token: f.Token, Data: f})
		default:
			b.pusherr(tok, "invalid_syntax")
			continue
		}
	}
	b.Pos, b.Tokens = pos, btoks
}

func (b *builder) get_method(toks []lexer.Token) *ast.Fn {
	tok := toks[0]
	if tok.Id == lexer.ID_UNSAFE {
		toks = toks[1:]
		if len(toks) == 0 || toks[0].Id != lexer.ID_FN {
			b.pusherr(tok, "invalid_syntax")
			return nil
		}
	} else if toks[0].Id != lexer.ID_FN {
		b.pusherr(tok, "invalid_syntax")
		return nil
	}
	f := new(ast.Fn)
	*f = b.Func(toks, true, false, false)
	f.IsUnsafe = tok.Id == lexer.ID_UNSAFE
	if f.Block != nil {
		f.Block.IsUnsafe = f.IsUnsafe
	}
	return f
}

func (b *builder) implFuncs(impl *ast.Impl, toks []lexer.Token) {
	if impl.Target.Id != types.VOID {
		b.implTraitFns(impl, toks)
		return
	}
	b.implStruct(impl, toks)
}

func (b *builder) Impl(toks []lexer.Token) {
	tok := toks[0]
	if len(toks) < 2 {
		b.pusherr(tok, "invalid_syntax")
		return
	}
	tok = toks[1]
	if tok.Id != lexer.ID_IDENT {
		b.pusherr(tok, "invalid_syntax")
		return
	}
	var impl ast.Impl
	if len(toks) < 3 {
		b.pusherr(tok, "invalid_syntax")
		return
	}
	impl.Base = tok
	tok = toks[2]
	if tok.Id != lexer.ID_ITER {
		if tok.Id == lexer.ID_BRACE && tok.Kind == lexer.KND_LBRACE {
			toks = toks[2:]
			goto body
		}
		b.pusherr(tok, "invalid_syntax")
		return
	}
	if len(toks) < 4 {
		b.pusherr(tok, "invalid_syntax")
		return
	}
	tok = toks[3]
	if tok.Id != lexer.ID_IDENT {
		b.pusherr(tok, "invalid_syntax")
		return
	}
	{
		i := 0
		impl.Target, _ = b.DataType(toks[3:4], &i, true)
		toks = toks[4:]
	}
body:
	i := 0
	bodyToks := b.getrange(&i, lexer.KND_LBRACE, lexer.KND_RBRACE, &toks)
	if bodyToks == nil {
		b.stop()
		b.pusherr(impl.Base, "body_not_exist")
		return
	}
	if i < len(toks) {
		b.pusherr(toks[i], "invalid_syntax")
	}
	b.implFuncs(&impl, bodyToks)
	b.Tree = append(b.Tree, ast.Node{Token: impl.Base, Data: impl})
}

func (b *builder) link_fn(toks []lexer.Token) {
	tok := toks[0]
	bpub := b.public
	b.public = false

	var link ast.CppLinkFn
	link.Token = tok
	link.Link = new(ast.Fn)
	*link.Link = b.Func(toks[1:], false, false, true)
	b.Tree = append(b.Tree, ast.Node{Token: tok, Data: link})

	b.public = bpub
}

func (b *builder) link_var(toks []lexer.Token) {
	tok := toks[0]
	bpub := b.public
	b.public = false

	var link ast.CppLinkVar
	link.Token = tok
	link.Link = new(ast.Var)
	*link.Link = b.Var(toks[1:], true, false)
	b.Tree = append(b.Tree, ast.Node{Token: tok, Data: link})

	b.public = bpub
}

func (b *builder) link_struct(toks []lexer.Token) {
	tok := toks[0]
	bpub := b.public
	b.public = false

	var link ast.CppLinkStruct
	link.Token = tok
	link.Link = b.build_struct(toks[1:], true)
	b.Tree = append(b.Tree, ast.Node{Token: tok, Data: link})

	b.public = bpub
}

func (b *builder) link_type_alias(toks []lexer.Token) {
	tok := toks[0]
	bpub := b.public
	b.public = false

	var link ast.CppLinkAlias
	link.Token = tok
	link.Link = b.TypeAlias(toks[1:])
	b.Tree = append(b.Tree, ast.Node{Token: tok, Data: link})

	b.public = bpub
}

func (b *builder) CppLink(toks []lexer.Token) {
	tok := toks[0]
	if len(toks) == 1 {
		b.pusherr(tok, "invalid_syntax")
		return
	}
	tok = toks[1]
	switch tok.Id {
	case lexer.ID_FN, lexer.ID_UNSAFE:
		b.link_fn(toks)
	case lexer.ID_LET:
		b.link_var(toks)
	case lexer.ID_STRUCT:
		b.link_struct(toks)
	case lexer.ID_TYPE:
		b.link_type_alias(toks)
	default:
		b.pusherr(tok, "invalid_syntax")
	}
}

func tokstoa(toks []lexer.Token) string {
	var str strings.Builder
	for _, tok := range toks {
		str.WriteString(tok.Kind)
	}
	return str.String()
}

func (b *builder) Use(toks []lexer.Token) {
	var use ast.UseDecl
	use.Token = toks[0]
	if len(toks) < 2 {
		b.pusherr(use.Token, "missing_use_path")
		return
	}
	toks = toks[1:]
	b.buildUseDecl(&use, toks)
	b.Tree = append(b.Tree, ast.Node{Token: use.Token, Data: use})
}

func (b *builder) getSelectors(toks []lexer.Token) []lexer.Token {
	i := 0
	toks = b.getrange(&i, lexer.KND_LBRACE, lexer.KND_RBRACE, &toks)
	parts, errs := ast.Parts(toks, lexer.ID_COMMA, true)
	if len(errs) > 0 {
		b.Errors = append(b.Errors, errs...)
		return nil
	}
	selectors := make([]lexer.Token, len(parts))
	for i, part := range parts {
		if len(part) > 1 {
			b.pusherr(part[1], "invalid_syntax")
		}
		tok := part[0]
		if tok.Id != lexer.ID_IDENT && tok.Id != lexer.ID_SELF {
			b.pusherr(tok, "invalid_syntax")
			continue
		}
		selectors[i] = tok
	}
	return selectors
}

func (b *builder) buildUseCppDecl(use *ast.UseDecl, toks []lexer.Token) {
	if len(toks) > 2 {
		b.pusherr(toks[2], "invalid_syntax")
	}
	tok := toks[1]
	if tok.Id != lexer.ID_LITERAL || (tok.Kind[0] != '`' && tok.Kind[0] != '"') {
		b.pusherr(tok, "invalid_expr")
		return
	}
	use.Cpp = true
	use.Path = tok.Kind[1 : len(tok.Kind)-1]
}

func (b *builder) buildUseDecl(use *ast.UseDecl, toks []lexer.Token) {
	var path strings.Builder
	path.WriteString(jane.STDLIB_PATH)
	path.WriteRune(os.PathSeparator)
	tok := toks[0]
	isStd := false
	if tok.Id == lexer.ID_CPP {
		b.buildUseCppDecl(use, toks)
		return
	}
	if tok.Id != lexer.ID_IDENT || tok.Kind != "std" {
		b.pusherr(toks[0], "invalid_syntax")
	}
	isStd = true
	if len(toks) < 3 {
		b.pusherr(tok, "invalid_syntax")
		return
	}
	toks = toks[2:]
	tok = toks[len(toks)-1]
	switch tok.Id {
	case lexer.ID_DBLCOLON:
		b.pusherr(tok, "invalid_syntax")
		return
	case lexer.ID_BRACE:
		if tok.Kind != lexer.KND_RBRACE {
			b.pusherr(tok, "invalid_syntax")
			return
		}
		var selectors []lexer.Token
		toks, selectors = ast.RangeLast(toks)
		use.Selectors = b.getSelectors(selectors)
		if len(toks) == 0 {
			b.pusherr(tok, "invalid_syntax")
			return
		}
		tok = toks[len(toks)-1]
		if tok.Id != lexer.ID_DBLCOLON {
			b.pusherr(tok, "invalid_syntax")
			return
		}
		toks = toks[:len(toks)-1]
		if len(toks) == 0 {
			b.pusherr(tok, "invalid_syntax")
			return
		}
	case lexer.ID_OP:
		if tok.Kind != lexer.KND_STAR {
			b.pusherr(tok, "invalid_syntax")
			return
		}
		toks = toks[:len(toks)-1]
		if len(toks) == 0 {
			b.pusherr(tok, "invalid_syntax")
			return
		}
		tok = toks[len(toks)-1]
		if tok.Id != lexer.ID_DBLCOLON {
			b.pusherr(tok, "invalid_syntax")
			return
		}
		toks = toks[:len(toks)-1]
		if len(toks) == 0 {
			b.pusherr(tok, "invalid_syntax")
			return
		}
		use.FullUse = true
	}
	for i, tok := range toks {
		if i%2 != 0 {
			if tok.Id != lexer.ID_DBLCOLON {
				b.pusherr(tok, "invalid_syntax")
			}
			path.WriteRune(os.PathSeparator)
			continue
		}
		if tok.Id != lexer.ID_IDENT {
			b.pusherr(tok, "invalid_syntax")
		}
		path.WriteString(tok.Kind)
	}
	use.LinkString = tokstoa(toks)
	if isStd {
		use.LinkString = "std::" + use.LinkString
	}
	use.Path = path.String()
}

func (b *builder) setup_receiver(f *ast.Fn, owner_id string) {
	if len(f.Params) == 0 {
		b.pusherr(f.Token, "missing_receiver")
		return
	}
	param := f.Params[0]
	if param.Id != lexer.KND_SELF {
		b.pusherr(f.Token, "missing_receiver")
		return
	}
	f.Receiver = new(ast.Var)
	f.Receiver.DataType = ast.Type{
		Id:   types.STRUCT,
		Kind: owner_id,
	}
	f.Receiver.Mutable = param.Mutable
	if param.DataType.Kind != "" && param.DataType.Kind[0] == '&' {
		f.Receiver.DataType.Kind = lexer.KND_AMPER + f.Receiver.DataType.Kind
	}
	f.Params = f.Params[1:]
}

func (b *builder) fn_prototype(toks []lexer.Token, i *int, method, anon bool) (f ast.Fn, ok bool) {
	ok = true
	f.Token = toks[*i]
	if f.Token.Id == lexer.ID_UNSAFE {
		f.IsUnsafe = true
		*i++
		if *i >= len(toks) {
			b.pusherr(f.Token, "invalid_syntax")
			ok = false
			return
		}
		f.Token = toks[*i]
	}
	*i++
	if *i >= len(toks) {
		b.pusherr(f.Token, "invalid_syntax")
		ok = false
		return
	}
	f.Public = b.public
	b.public = false
	if anon {
		f.Id = lexer.ANONYMOUS_ID
	} else {
		tok := toks[*i]
		if tok.Id != lexer.ID_IDENT {
			b.pusherr(tok, "invalid_syntax")
			ok = false
		}
		f.Id = tok.Kind
		*i++
	}

	f.RetType.DataType.Id = types.VOID
	f.RetType.DataType.Kind = types.TYPE_MAP[f.RetType.DataType.Id]
	if *i >= len(toks) {
		b.pusherr(f.Token, "invalid_syntax")
		return
	}

	generics_toks := ast.Range(i, lexer.KND_LBRACKET, lexer.KND_RBRACKET, toks)
	if generics_toks != nil {
		f.Generics = b.Generics(generics_toks)
		if len(f.Generics) > 0 {
			f.Combines = new([][]ast.Type)
		}
	}

	if toks[*i].Kind != lexer.KND_LPAREN {
		b.pusherr(toks[*i], "missing_function_parentheses")
		return
	}
	params_toks := b.getrange(i, lexer.KND_LPAREN, lexer.KND_RPARENT, &toks)
	if len(params_toks) > 0 {
		f.Params = b.Params(params_toks, method, false)
	}

	t, ret_ok := b.FnRetDataType(toks, i)
	if ret_ok {
		f.RetType = t
		*i++
	}
	return
}

func (b *builder) stop() {
	b.Pos = -1
}

func (b *builder) Func(toks []lexer.Token, method, anon, prototype bool) (f ast.Fn) {
	var ok bool
	i := 0
	f, ok = b.fn_prototype(toks, &i, method, anon)
	if prototype {
		if i+1 < len(toks) {
			b.pusherr(toks[i+1], "invalid_syntax")
		}
		return
	} else if !ok {
		return
	}
	if i >= len(toks) {
		b.stop()
		b.pusherr(f.Token, "body_not_exist")
		return
	}
	block_toks := b.getrange(&i, lexer.KND_LBRACE, lexer.KND_RBRACE, &toks)
	if block_toks != nil {
		f.Block = b.Block(block_toks)
		f.Block.IsUnsafe = f.IsUnsafe
		if i < len(toks) {
			b.pusherr(toks[i], "invalid_syntax")
		}
	} else {
		b.stop()
		b.pusherr(f.Token, "body_not_exist")
		b.Tokens = append(toks, b.Tokens...)
	}
	return
}

func (b *builder) generic(toks []lexer.Token) *ast.GenericType {
	if len(toks) > 1 {
		b.pusherr(toks[1], "invalid_syntax")
	}
	gt := new(ast.GenericType)
	gt.Token = toks[0]
	if gt.Token.Id != lexer.ID_IDENT {
		b.pusherr(gt.Token, "invalid_syntax")
	}
	gt.Id = gt.Token.Kind
	return gt
}

func (b *builder) Generics(toks []lexer.Token) []*ast.GenericType {
	tok := toks[0]
	if len(toks) == 0 {
		b.pusherr(tok, "missing_expr")
		return nil
	}
	parts, errs := ast.Parts(toks, lexer.ID_COMMA, true)
	b.Errors = append(b.Errors, errs...)
	generics := make([]*ast.GenericType, len(parts))
	for i, part := range parts {
		if len(parts) == 0 {
			continue
		}
		generics[i] = b.generic(part)
	}
	return generics
}

func (b *builder) GlobalTypeAlias(toks []lexer.Token) ast.Node {
	t := b.TypeAlias(toks)
	t.Pub = b.public
	b.public = false
	return ast.Node{Token: t.Token, Data: t}
}

func (b *builder) GlobalVar(toks []lexer.Token) {
	if toks == nil {
		return
	}
	bs := block_st{toks: toks}
	s := b.VarSt(&bs, true)
	b.Tree = append(b.Tree, ast.Node{
		Token: s.Token,
		Data:  s,
	})
}

func (b *builder) build_self(toks []lexer.Token) (model ast.Param) {
	if len(toks) == 0 {
		return
	}
	i := 0
	if toks[i].Id == lexer.ID_MUT {
		model.Mutable = true
		i++
		if i >= len(toks) {
			b.pusherr(toks[i-1], "invalid_syntax")
			return
		}
	}
	if toks[i].Kind == lexer.KND_AMPER {
		model.DataType.Kind = "&"
		i++
		if i >= len(toks) {
			b.pusherr(toks[i-1], "invalid_syntax")
			return
		}
	}
	if toks[i].Id == lexer.ID_SELF {
		model.Id = lexer.KND_SELF
		model.Token = toks[i]
		i++
		if i < len(toks) {
			b.pusherr(toks[i+1], "invalid_syntax")
		}
	}
	return
}

func (b *builder) Params(toks []lexer.Token, method, mustPure bool) []ast.Param {
	parts, errs := ast.Parts(toks, lexer.ID_COMMA, true)
	b.Errors = append(b.Errors, errs...)
	if len(parts) == 0 {
		return nil
	}
	var params []ast.Param
	if method && len(parts) > 0 {
		param := b.build_self(parts[0])
		if param.Id == lexer.KND_SELF {
			params = append(params, param)
			parts = parts[1:]
		}
	}
	for _, part := range parts {
		b.pushParam(&params, part, mustPure)
	}
	b.checkParams(&params)
	return params
}

func (b *builder) checkParams(params *[]ast.Param) {
	for i := range *params {
		param := &(*params)[i]
		if param.Id == lexer.KND_SELF || param.DataType.Token.Id != lexer.ID_NA {
			continue
		}
		if param.Token.Id == lexer.ID_NA {
			b.pusherr(param.Token, "missing_type")
		} else {
			param.DataType.Token = param.Token
			param.DataType.Id = types.ID
			param.DataType.Kind = param.DataType.Token.Kind
			param.DataType.Original = param.DataType
			param.Id = lexer.ANONYMOUS_ID
			param.Token = lexer.Token{}
		}
	}
}

func (b *builder) paramTypeBegin(param *ast.Param, i *int, toks []lexer.Token) {
	for ; *i < len(toks); *i++ {
		tok := toks[*i]
		switch tok.Id {
		case lexer.ID_OP:
			switch tok.Kind {
			case lexer.KND_TRIPLE_DOT:
				if param.Variadic {
					b.pusherr(tok, "already_variadic")
					continue
				}
				param.Variadic = true
			default:
				return
			}
		default:
			return
		}
	}
}

func (b *builder) paramBodyId(param *ast.Param, tok lexer.Token) {
	if lexer.IsIgnoreId(tok.Kind) {
		param.Id = lexer.ANONYMOUS_ID
		return
	}
	param.Id = tok.Kind
}

func (b *builder) paramBody(param *ast.Param, i *int, toks []lexer.Token, mustPure bool) {
	b.paramBodyId(param, toks[*i])
	tok := toks[*i]
	toks = toks[*i+1:]
	if len(toks) == 0 {
		return
	} else if len(toks) < 2 {
		b.pusherr(tok, "missing_type")
		return
	}
	tok = toks[*i]
	if tok.Id != lexer.ID_COLON {
		b.pusherr(tok, "invalid_syntax")
		return
	}
	toks = toks[*i+1:]
	b.paramType(param, toks, mustPure)
}

func (b *builder) paramType(param *ast.Param, toks []lexer.Token, mustPure bool) {
	i := 0
	if !mustPure {
		b.paramTypeBegin(param, &i, toks)
		if i >= len(toks) {
			return
		}
	}
	param.DataType, _ = b.DataType(toks, &i, true)
	i++
	if i < len(toks) {
		b.pusherr(toks[i], "invalid_syntax")
	}
}

func (b *builder) pushParam(params *[]ast.Param, toks []lexer.Token, mustPure bool) {
	var param ast.Param
	param.Token = toks[0]
	if param.Token.Id == lexer.ID_MUT {
		param.Mutable = true
		if len(toks) == 1 {
			b.pusherr(toks[0], "invalid_syntax")
			return
		}
		toks = toks[1:]
		param.Token = toks[0]
	}
	if param.Token.Id != lexer.ID_IDENT {
		param.Id = lexer.ANONYMOUS_ID
		b.paramType(&param, toks, mustPure)
	} else {
		i := 0
		b.paramBody(&param, &i, toks, mustPure)
	}
	*params = append(*params, param)
}

func (b *builder) datatype(t *ast.Type, toks []lexer.Token, i *int, err bool) (ok bool) {
	tb := type_builder{
		r:      b,
		t:      t,
		tokens: toks,
		i:      i,
		err:    err,
	}
	return tb.build()
}

func (b *builder) DataType(toks []lexer.Token, i *int, err bool) (t ast.Type, ok bool) {
	tok := toks[*i]
	ok = b.datatype(&t, toks, i, err)
	if err && t.Token.Id == lexer.ID_NA {
		b.pusherr(tok, "invalid_type")
	}
	return
}

func (b *builder) fnMultiTypeRet(toks []lexer.Token, i *int) (t ast.RetType, ok bool) {
	tok := toks[*i]
	t.DataType.Kind += tok.Kind
	*i++
	if *i >= len(toks) {
		*i--
		t.DataType, ok = b.DataType(toks, i, false)
		return
	}
	tok = toks[*i]
	*i--
	rang := ast.Range(i, lexer.KND_LPAREN, lexer.KND_RPARENT, toks)
	params := b.Params(rang, false, true)
	types := make([]ast.Type, len(params))
	for i, param := range params {
		types[i] = param.DataType
		if param.Id != lexer.ANONYMOUS_ID {
			param.Token.Kind = param.Id
		} else {
			param.Token.Kind = lexer.IGNORE_ID
		}
		t.Identifiers = append(t.Identifiers, param.Token)
	}
	if len(types) > 1 {
		t.DataType.MultiTyped = true
		t.DataType.Tag = types
	} else {
		t.DataType = types[0]
	}
	*i--
	ok = true
	return
}

func (b *builder) FnRetDataType(toks []lexer.Token, i *int) (t ast.RetType, ok bool) {
	t.DataType.Id = types.VOID
	t.DataType.Kind = types.TYPE_MAP[t.DataType.Id]
	if *i >= len(toks) {
		return
	}
	tok := toks[*i]
	switch tok.Id {
	case lexer.ID_BRACE:
		if tok.Kind == lexer.KND_LBRACE {
			return
		}
	case lexer.ID_OP:
		if tok.Kind == lexer.KND_EQ {
			return
		}
	case lexer.ID_COLON:
		if *i+1 >= len(toks) {
			b.pusherr(tok, "missing_type")
			return
		}
		*i++
		tok = toks[*i]
		if tok.Id == lexer.ID_BRACE {
			switch tok.Kind {
			case lexer.KND_LPAREN:
				return b.fnMultiTypeRet(toks, i)
			case lexer.KND_LBRACE:
				return
			}
		}
		t.DataType, ok = b.DataType(toks, i, true)
		return
	}
	*i++
	b.pusherr(tok, "invalid_syntax")
	return
}

func (b *builder) pushStToBlock(bs *block_st) {
	if len(bs.toks) == 0 {
		return
	}
	lastTok := bs.toks[len(bs.toks)-1]
	if lastTok.Id == lexer.ID_SEMICOLON {
		if len(bs.toks) == 1 {
			return
		}
		bs.toks = bs.toks[:len(bs.toks)-1]
	}
	s := b.St(bs)
	if s.Data == nil {
		return
	}
	s.WithTerminator = bs.terminated
	bs.block.Tree = append(bs.block.Tree, s)
}

func get_next_st(toks []lexer.Token) []lexer.Token {
	pos, terminated := ast.NextStPos(toks, 0)
	if terminated {
		return toks[:pos-1]
	}
	return toks[:pos]
}

func set_to_next_st(bs *block_st) {
	if bs.nextToks != nil {
		bs.toks = bs.nextToks
		bs.nextToks = nil
		return
	}
	*bs.srcToks = (*bs.srcToks)[bs.pos:]
	bs.pos, bs.terminated = ast.NextStPos(*bs.srcToks, 0)
	if bs.terminated {
		bs.toks = (*bs.srcToks)[:bs.pos-1]
	} else {
		bs.toks = (*bs.srcToks)[:bs.pos]
	}
}

func blockStFinished(bs *block_st) bool {
	return bs.nextToks == nil && bs.pos >= len(*bs.srcToks)
}

func block_from_tree(tree []ast.St) *ast.Block {
	block := new(ast.Block)
	block.Tree = tree
	return block
}

func (b *builder) Block(toks []lexer.Token) (block *ast.Block) {
	block = new(ast.Block)
	var bs block_st
	bs.block = block
	bs.srcToks = &toks
	for {
		set_to_next_st(&bs)
		b.pushStToBlock(&bs)
		if blockStFinished(&bs) {
			break
		}
	}
	return
}

func (b *builder) St(bs *block_st) (s ast.St) {
	tok := bs.toks[0]
	if tok.Id == lexer.ID_IDENT {
		s, ok := b.IdSt(bs)
		if ok {
			return s
		}
	}
	s, ok := b.AssignSt(bs.toks)
	if ok {
		return s
	}
	switch tok.Id {
	case lexer.ID_CONST, lexer.ID_LET, lexer.ID_MUT:
		return b.VarSt(bs, true)
	case lexer.ID_RET:
		return b.RetSt(bs.toks)
	case lexer.ID_ITER:
		return b.IterExpr(bs)
	case lexer.ID_BREAK:
		return b.BreakSt(bs.toks)
	case lexer.ID_CONTINUE:
		return b.ContinueSt(bs.toks)
	case lexer.ID_IF:
		return b.Conditional(bs)
	case lexer.ID_COMMENT:
		return b.CommentSt(bs.toks[0])
	case lexer.ID_CO:
		return b.ConcurrentCallSt(bs.toks)
	case lexer.ID_GOTO:
		return b.GotoSt(bs.toks)
	case lexer.ID_FALL:
		return b.Fallthrough(bs.toks)
	case lexer.ID_TYPE:
		t := b.TypeAlias(bs.toks)
		s.Token = t.Token
		s.Data = t
		return
	case lexer.ID_MATCH:
		return b.MatchCase(bs.toks)
	case lexer.ID_UNSAFE, lexer.ID_DEFER:
		return b.blockSt(bs.toks)
	case lexer.ID_BRACE:
		if tok.Kind == lexer.KND_LBRACE {
			return b.blockSt(bs.toks)
		}
	}
	if ast.IsFnCall(bs.toks) != nil {
		return b.ExprSt(bs)
	}
	b.pusherr(tok, "invalid_syntax")
	return
}

func (b *builder) blockSt(toks []lexer.Token) ast.St {
	is_unsafe := false
	is_deferred := false
	tok := toks[0]
	if tok.Id == lexer.ID_UNSAFE {
		is_unsafe = true
		toks = toks[1:]
		if len(toks) == 0 {
			b.pusherr(tok, "invalid_syntax")
			return ast.St{}
		}
		tok = toks[0]
		if tok.Id == lexer.ID_DEFER {
			is_deferred = true
			toks = toks[1:]
			if len(toks) == 0 {
				b.pusherr(tok, "invalid_syntax")
				return ast.St{}
			}
		}
	} else if tok.Id == lexer.ID_DEFER {
		is_deferred = true
		toks = toks[1:]
		if len(toks) == 0 {
			b.pusherr(tok, "invalid_syntax")
			return ast.St{}
		}
	}

	i := 0
	toks = ast.Range(&i, lexer.KND_LBRACE, lexer.KND_RBRACE, toks)
	if len(toks) == 0 {
		b.pusherr(tok, "invalid_syntax")
		return ast.St{}
	} else if i < len(toks) {
		b.pusherr(toks[i], "invalid_syntax")
	}
	block := b.Block(toks)
	block.IsUnsafe = is_unsafe
	block.Deferred = is_deferred
	return ast.St{Token: tok, Data: block}
}

func (b *builder) assignInfo(toks []lexer.Token) (info ast.AssignInfo) {
	info.Ok = true
	brace_n := 0
	for i, tok := range toks {
		if tok.Id == lexer.ID_BRACE {
			switch tok.Kind {
			case lexer.KND_LBRACE, lexer.KND_LBRACKET, lexer.KND_LPAREN:
				brace_n++
			default:
				brace_n--
			}
		}
		if brace_n > 0 {
			continue
		} else if tok.Id != lexer.ID_OP {
			continue
		} else if !ast.IsAssignOp(tok.Kind) {
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
			info.Ok = ast.IsPostfixOp(info.Setter.Kind)
			break
		}
		info.Right = toks[i+1:]
		if ast.IsPostfixOp(info.Setter.Kind) {
			if info.Right != nil {
				b.pusherr(info.Right[0], "invalid_syntax")
				info.Right = nil
			}
		}
		break
	}
	return
}

func (b *builder) build_assign_left(toks []lexer.Token) (l ast.AssignLeft) {
	l.Expr.Tokens = toks
	if l.Expr.Tokens[0].Id == lexer.ID_IDENT {
		l.Var.Token = l.Expr.Tokens[0]
		l.Var.Id = l.Var.Token.Kind
	}
	l.Expr = b.Expr(l.Expr.Tokens)
	return
}

func (b *builder) assignLefts(parts [][]lexer.Token) []ast.AssignLeft {
	var lefts []ast.AssignLeft
	for _, p := range parts {
		l := b.build_assign_left(p)
		lefts = append(lefts, l)
	}
	return lefts
}

func (b *builder) assignExprs(toks []lexer.Token) []ast.Expr {
	parts, errs := ast.Parts(toks, lexer.ID_COMMA, true)
	if len(errs) > 0 {
		b.Errors = append(b.Errors, errs...)
		return nil
	}
	exprs := make([]ast.Expr, len(parts))
	for i, p := range parts {
		exprs[i] = b.Expr(p)
	}
	return exprs
}

func (b *builder) AssignSt(toks []lexer.Token) (s ast.St, _ bool) {
	assign, ok := b.AssignExpr(toks)
	if !ok {
		return
	}
	s.Token = toks[0]
	s.Data = assign
	return s, true
}

func (b *builder) AssignExpr(toks []lexer.Token) (assign ast.Assign, ok bool) {
	if !ast.CheckAssignTokens(toks) {
		return
	}
	switch toks[0].Id {
	case lexer.ID_LET:
		return b.letDeclAssign(toks)
	default:
		return b.plainAssign(toks)
	}
}

func (b *builder) letDeclAssign(toks []lexer.Token) (assign ast.Assign, ok bool) {
	if len(toks) < 1 {
		return
	}
	toks = toks[1:]
	tok := toks[0]
	if tok.Id != lexer.ID_BRACE || tok.Kind != lexer.KND_LPAREN {
		return
	}
	ok = true
	var i int
	rang := ast.Range(&i, lexer.KND_LPAREN, lexer.KND_RPARENT, toks)
	if rang == nil {
		b.pusherr(tok, "invalid_syntax")
		return
	} else if i+1 < len(toks) {
		assign.Setter = toks[i]
		i++
		assign.Right = b.assignExprs(toks[i:])
	}
	parts, errs := ast.Parts(rang, lexer.ID_COMMA, true)
	if len(errs) > 0 {
		b.Errors = append(b.Errors, errs...)
		return
	}
	for _, p := range parts {
		mutable := false
		tok := p[0]
		if tok.Id == lexer.ID_MUT {
			mutable = true
			p = p[1:]
			if len(p) != 1 {
				b.pusherr(tok, "invalid_syntax")
				continue
			}
		}
		if p[0].Id != lexer.ID_IDENT && p[0].Id != lexer.ID_BRACE && p[0].Kind != lexer.KND_LPAREN {
			b.pusherr(tok, "invalid_syntax")
			continue
		}
		l := b.build_assign_left(p)
		l.Var.Mutable = mutable
		l.Var.New = l.Var.Id != "" && !lexer.IsIgnoreId(l.Var.Id)
		l.Var.SetterTok = assign.Setter
		assign.Left = append(assign.Left, l)
	}
	return
}

func (b *builder) plainAssign(toks []lexer.Token) (assign ast.Assign, ok bool) {
	info := b.assignInfo(toks)
	if !info.Ok {
		return
	}
	ok = true
	assign.Setter = info.Setter
	parts, errs := ast.Parts(info.Left, lexer.ID_COMMA, true)
	if len(errs) > 0 {
		b.Errors = append(b.Errors, errs...)
		return
	}
	assign.Left = b.assignLefts(parts)
	if info.Right != nil {
		assign.Right = b.assignExprs(info.Right)
	}
	return
}

func (b *builder) IdSt(bs *block_st) (s ast.St, ok bool) {
	if len(bs.toks) == 1 {
		return
	}
	tok := bs.toks[1]
	switch tok.Id {
	case lexer.ID_COLON:
		return b.LabelSt(bs), true
	}
	return
}

func (b *builder) LabelSt(bs *block_st) ast.St {
	var l ast.Label
	l.Token = bs.toks[0]
	l.Label = l.Token.Kind
	if len(bs.toks) > 2 {
		bs.nextToks = bs.toks[2:]
	}
	return ast.St{
		Token: l.Token,
		Data:  l,
	}
}

func (b *builder) ExprSt(bs *block_st) ast.St {
	st := ast.ExprSt{Expr: b.Expr(bs.toks)}
	return ast.St{Token: bs.toks[0], Data: st}
}

func (b *builder) Args(toks []lexer.Token, targeting bool) *ast.Args {
	args := new(ast.Args)
	last := 0
	brace_n := 0
	for i, tok := range toks {
		if tok.Id == lexer.ID_BRACE {
			switch tok.Kind {
			case lexer.KND_LBRACE, lexer.KND_LBRACKET, lexer.KND_LPAREN:
				brace_n++
			default:
				brace_n--
			}
		}
		if brace_n > 0 || tok.Id != lexer.ID_COMMA {
			continue
		}
		b.pushArg(args, targeting, toks[last:i], tok)
		last = i + 1
	}
	if last < len(toks) {
		if last == 0 {
			if len(toks) > 0 {
				b.pushArg(args, targeting, toks[last:], toks[last])
			}
		} else {
			b.pushArg(args, targeting, toks[last:], toks[last-1])
		}
	}
	return args
}

func (b *builder) pushArg(args *ast.Args, targeting bool, toks []lexer.Token, err lexer.Token) {
	if len(toks) == 0 {
		b.pusherr(err, "invalid_syntax")
		return
	}
	var arg ast.Arg
	arg.Token = toks[0]
	if targeting && arg.Token.Id == lexer.ID_IDENT {
		if len(toks) > 1 {
			tok := toks[1]
			if tok.Id == lexer.ID_COLON {
				args.Targeted = true
				arg.TargetId = arg.Token.Kind
				toks = toks[2:]
			}
		}
	}
	arg.Expr = b.Expr(toks)
	args.Src = append(args.Src, arg)
}

func (b *builder) varBegin(v *ast.Var, i *int, toks []lexer.Token) {
	tok := toks[*i]
	switch tok.Id {
	case lexer.ID_LET:
		*i++
		if toks[*i].Id == lexer.ID_MUT {
			v.Mutable = true
			*i++
		}
	case lexer.ID_CONST:
		*i++
		if v.Constant {
			b.pusherr(tok, "already_const")
			break
		}
		v.Constant = true
		if !v.Mutable {
			break
		}
		fallthrough
	default:
		b.pusherr(tok, "invalid_syntax")
		return
	}
	if *i >= len(toks) {
		b.pusherr(tok, "invalid_syntax")
	}
}

func (b *builder) varTypeNExpr(v *ast.Var, toks []lexer.Token, i int, expr bool) {
	tok := toks[i]
	if tok.Id == lexer.ID_COLON {
		if i >= len(toks) ||
			(toks[i].Id == lexer.ID_OP && toks[i].Kind == lexer.KND_EQ) {
			b.pusherr(tok, "missing_type")
			return
		}
		t, ok := b.DataType(toks, &i, false)
		if ok {
			v.DataType = t
			i++
			if i >= len(toks) {
				return
			}
			tok = toks[i]
		}
	}
	if expr && tok.Id == lexer.ID_OP {
		if tok.Kind != lexer.KND_EQ {
			b.pusherr(tok, "invalid_syntax")
			return
		}
		valueToks := toks[i+1:]
		if len(valueToks) == 0 {
			b.pusherr(tok, "missing_expr")
			return
		}
		v.Expr = b.Expr(valueToks)
		v.SetterTok = tok
	} else {
		b.pusherr(tok, "invalid_syntax")
	}
}

func (b *builder) Var(toks []lexer.Token, begin, expr bool) (v ast.Var) {
	v.Public = b.public
	b.public = false
	i := 0
	v.Token = toks[i]
	if begin {
		b.varBegin(&v, &i, toks)
		if i >= len(toks) {
			return
		}
	}
	v.Token = toks[i]
	if v.Token.Id != lexer.ID_IDENT {
		b.pusherr(v.Token, "invalid_syntax")
		return
	}
	v.Id = v.Token.Kind
	v.DataType.Id = types.VOID
	v.DataType.Kind = types.TYPE_MAP[v.DataType.Id]
	if i >= len(toks) {
		return
	}
	i++
	if i < len(toks) {
		b.varTypeNExpr(&v, toks, i, expr)
	} else if !expr {
		b.pusherr(v.Token, "missing_type")
	}
	return
}

func (b *builder) VarSt(bs *block_st, expr bool) ast.St {
	v := b.Var(bs.toks, true, expr)
	v.Owner = bs.block
	return ast.St{Token: v.Token, Data: v}
}

func (b *builder) CommentSt(tok lexer.Token) (s ast.St) {
	s.Token = tok
	tok.Kind = strings.TrimSpace(tok.Kind[2:])
	s.Data = ast.Comment{Content: tok.Kind}
	return
}

func (b *builder) ConcurrentCallSt(toks []lexer.Token) (s ast.St) {
	var cc ast.ConcurrentCall
	cc.Token = toks[0]
	toks = toks[1:]
	if len(toks) == 0 {
		b.pusherr(cc.Token, "missing_expr")
		return
	}
	if ast.IsFnCall(toks) == nil {
		b.pusherr(cc.Token, "expr_not_func_call")
	}
	cc.Expr = b.Expr(toks)
	s.Token = cc.Token
	s.Data = cc
	return
}

func (b *builder) Fallthrough(toks []lexer.Token) (s ast.St) {
	s.Token = toks[0]
	if len(toks) > 1 {
		b.pusherr(toks[1], "invalid_syntax")
	}
	s.Data = ast.Fall{
		Token: s.Token,
	}
	return
}

func (b *builder) GotoSt(toks []lexer.Token) (s ast.St) {
	s.Token = toks[0]
	if len(toks) == 1 {
		b.pusherr(s.Token, "missing_goto_label")
		return
	} else if len(toks) > 2 {
		b.pusherr(toks[2], "invalid_syntax")
	}
	idTok := toks[1]
	if idTok.Id != lexer.ID_IDENT {
		b.pusherr(idTok, "invalid_syntax")
		return
	}
	var gt ast.Goto
	gt.Token = s.Token
	gt.Label = idTok.Kind
	s.Data = gt
	return
}

func (b *builder) RetSt(toks []lexer.Token) ast.St {
	var ret ast.Ret
	ret.Token = toks[0]
	if len(toks) > 1 {
		ret.Expr = b.Expr(toks[1:])
	}
	return ast.St{
		Token: ret.Token,
		Data:  ret,
	}
}

func (b *builder) getWhileIterProfile(toks []lexer.Token) ast.IterWhile {
	return ast.IterWhile{
		Expr: b.Expr(toks),
	}
}

func (b *builder) getForeachVarsToks(toks []lexer.Token) [][]lexer.Token {
	vars, errs := ast.Parts(toks, lexer.ID_COMMA, true)
	b.Errors = append(b.Errors, errs...)
	return vars
}

func (b *builder) getVarProfile(toks []lexer.Token) (v ast.Var) {
	if len(toks) == 0 {
		return
	}
	v.Token = toks[0]
	if v.Token.Id == lexer.ID_MUT {
		v.Mutable = true
		if len(toks) == 1 {
			b.pusherr(v.Token, "invalid_syntax")
		}
		v.Token = toks[1]
	} else if len(toks) > 1 {
		b.pusherr(toks[1], "invalid_syntax")
	}
	if v.Token.Id != lexer.ID_IDENT {
		b.pusherr(v.Token, "invalid_syntax")
		return
	}
	v.Id = v.Token.Kind
	v.New = true
	return
}

func (b *builder) getForeachIterVars(varsToks [][]lexer.Token) []ast.Var {
	var vars []ast.Var
	for _, toks := range varsToks {
		vars = append(vars, b.getVarProfile(toks))
	}
	return vars
}

func (b *builder) setup_foreach_explicit_vars(f *ast.IterForeach, toks []lexer.Token) {
	i := 0
	rang := ast.Range(&i, lexer.KND_LPAREN, lexer.KND_RPARENT, toks)
	if i < len(toks) {
		b.pusherr(f.InToken, "invalid_syntax")
	}
	b.setup_foreach_plain_vars(f, rang)
}

func (b *builder) setup_foreach_plain_vars(f *ast.IterForeach, toks []lexer.Token) {
	varsToks := b.getForeachVarsToks(toks)
	if len(varsToks) == 0 {
		return
	}
	if len(varsToks) > 2 {
		b.pusherr(f.InToken, "much_foreach_vars")
	}
	vars := b.getForeachIterVars(varsToks)
	f.KeyA = vars[0]
	if len(vars) > 1 {
		f.KeyB = vars[1]
	} else {
		f.KeyB.Id = lexer.IGNORE_ID
	}
}

func (b *builder) setup_foreach_vars(f *ast.IterForeach, toks []lexer.Token) {
	if toks[0].Id == lexer.ID_BRACE {
		if toks[0].Kind != lexer.KND_LPAREN {
			b.pusherr(toks[0], "invalid_syntax")
			return
		}
		b.setup_foreach_explicit_vars(f, toks)
		return
	}
	b.setup_foreach_plain_vars(f, toks)
}

func (b *builder) getForeachIterProfile(
	varToks, exprToks []lexer.Token,
	inTok lexer.Token,
) ast.IterForeach {
	var foreach ast.IterForeach
	foreach.InToken = inTok
	if len(exprToks) == 0 {
		b.pusherr(inTok, "missing_expr")
		return foreach
	}
	foreach.Expr = b.Expr(exprToks)
	if len(varToks) == 0 {
		foreach.KeyA.Id = lexer.IGNORE_ID
		foreach.KeyB.Id = lexer.IGNORE_ID
	} else {
		b.setup_foreach_vars(&foreach, varToks)
	}
	return foreach
}

func (b *builder) getIterProfile(toks []lexer.Token, _ lexer.Token) any {
	brace_n := 0
	for i, tok := range toks {
		if tok.Id == lexer.ID_BRACE {
			switch tok.Kind {
			case lexer.KND_LBRACE, lexer.KND_LBRACKET, lexer.KND_LPAREN:
				brace_n++
				continue
			default:
				brace_n--
			}
		}
		if brace_n != 0 {
			continue
		}
		switch tok.Id {
		case lexer.ID_IN:
			varToks := toks[:i]
			exprToks := toks[i+1:]
			return b.getForeachIterProfile(varToks, exprToks, tok)
		}
	}
	return b.getWhileIterProfile(toks)
}

func (b *builder) next_st(toks []lexer.Token) ast.St {
	s := b.St(&block_st{toks: toks})
	switch s.Data.(type) {
	case ast.ExprSt, ast.Assign, ast.Var:
	default:
		b.pusherr(toks[0], "invalid_syntax")
	}
	return s
}

func (b *builder) getWhileNextIterProfile(bs *block_st) (s ast.St) {
	var iter ast.Iter
	iter.Token = bs.toks[0]
	bs.toks = bs.toks[1:]
	profile := ast.IterWhile{}
	if len(bs.toks) > 0 {
		profile.Expr = b.Expr(bs.toks)
	}
	if blockStFinished(bs) {
		b.pusherr(iter.Token, "invalid_syntax")
		return
	}
	set_to_next_st(bs)
	st_toks := ast.GetBlockExpr(bs.toks)
	if len(st_toks) > 0 {
		profile.Next = b.next_st(st_toks)
	}
	i := len(st_toks)
	blockToks := b.getrange(&i, lexer.KND_LBRACE, lexer.KND_RBRACE, &bs.toks)
	if blockToks == nil {
		b.stop()
		b.pusherr(iter.Token, "body_not_exist")
		return
	}
	if i < len(bs.toks) {
		b.pusherr(bs.toks[i], "invalid_syntax")
	}
	iter.Block = b.Block(blockToks)
	iter.Profile = profile
	return ast.St{Token: iter.Token, Data: iter}
}

func (b *builder) commonIterProfile(toks []lexer.Token) (s ast.St) {
	var iter ast.Iter
	iter.Token = toks[0]
	toks = toks[1:]
	if len(toks) == 0 {
		b.stop()
		b.pusherr(iter.Token, "body_not_exist")
		return
	}
	exprToks := ast.GetBlockExpr(toks)
	if len(exprToks) > 0 {
		iter.Profile = b.getIterProfile(exprToks, iter.Token)
	}
	i := len(exprToks)
	blockToks := b.getrange(&i, lexer.KND_LBRACE, lexer.KND_RBRACE, &toks)
	if blockToks == nil {
		b.stop()
		b.pusherr(iter.Token, "body_not_exist")
		return
	}
	if i < len(toks) {
		b.pusherr(toks[i], "invalid_syntax")
	}
	iter.Block = b.Block(blockToks)
	return ast.St{Token: iter.Token, Data: iter}
}

func (b *builder) IterExpr(bs *block_st) ast.St {
	if bs.terminated {
		return b.getWhileNextIterProfile(bs)
	}
	return b.commonIterProfile(bs.toks)
}

func (b *builder) caseexprs(toks *[]lexer.Token, type_match bool) []ast.Expr {
	var exprs []ast.Expr
	push_expr := func(toks []lexer.Token, _ lexer.Token) {
		if len(toks) > 0 {
			if type_match {
				i := 0
				t, ok := b.DataType(toks, &i, true)
				if ok {
					exprs = append(exprs, ast.Expr{
						Tokens: toks,
						Op:     t,
					})
				}
				i++
				if i < len(toks) {
					b.pusherr(toks[i], "invalid_syntax")
				}
				return
			}
			exprs = append(exprs, b.Expr(toks))
		}
	}
	brace_n := 0
	j := 0
	var i int
	var tok lexer.Token
	for i, tok = range *toks {
		if tok.Id == lexer.ID_BRACE {
			switch tok.Kind {
			case lexer.KND_LPAREN, lexer.KND_LBRACE, lexer.KND_LBRACKET:
				brace_n++
			default:
				brace_n--
			}
			continue
		} else if brace_n != 0 {
			continue
		}
		switch {
		case tok.Id == lexer.ID_OP && tok.Kind == lexer.KND_VLINE:
			push_expr((*toks)[j:i], tok)
			j = i + 1
		case tok.Id == lexer.ID_COLON:
			push_expr((*toks)[j:i], tok)
			*toks = (*toks)[i+1:]
			return exprs
		}
	}
	b.pusherr((*toks)[0], "invalid_syntax")
	*toks = nil
	return nil
}

func (b *builder) caseblock(toks *[]lexer.Token) *ast.Block {
	n := 0
	for {
		next := get_next_st((*toks)[n:])
		if len(next) == 0 {
			break
		}
		tok := next[0]
		if tok.Id != lexer.ID_OP || tok.Kind != lexer.KND_VLINE {
			n += len(next)
			continue
		}
		block := b.Block((*toks)[:n])
		*toks = (*toks)[n:]
		return block
	}
	block := b.Block(*toks)
	*toks = nil
	return block
}

func (b *builder) getcase(toks *[]lexer.Token, type_match bool) (ast.Case, bool) {
	var c ast.Case
	c.Token = (*toks)[0]
	*toks = (*toks)[1:]
	c.Exprs = b.caseexprs(toks, type_match)
	c.Block = b.caseblock(toks)
	is_default := len(c.Exprs) == 0
	return c, is_default
}

func (b *builder) cases(toks []lexer.Token, type_match bool) ([]ast.Case, *ast.Case) {
	var cases []ast.Case
	var def *ast.Case
	for len(toks) > 0 {
		tok := toks[0]
		if tok.Id != lexer.ID_OP || tok.Kind != lexer.KND_VLINE {
			b.pusherr(tok, "invalid_syntax")
			break
		}
		c, is_default := b.getcase(&toks, type_match)
		if is_default {
			c.Token = tok
			if def == nil {
				def = new(ast.Case)
				*def = c
			} else {
				b.pusherr(tok, "invalid_syntax")
			}
		} else {
			cases = append(cases, c)
		}
	}
	return cases, def
}

func (b *builder) MatchCase(toks []lexer.Token) (s ast.St) {
	m := new(ast.Match)
	m.Token = toks[0]
	s.Token = m.Token
	toks = toks[1:]

	if len(toks) > 0 && toks[0].Id == lexer.ID_TYPE {
		m.TypeMatch = true
		toks = toks[1:]
	}

	exprToks := ast.GetBlockExpr(toks)
	if len(exprToks) > 0 {
		m.Expr = b.Expr(exprToks)
	} else if m.TypeMatch {
		b.pusherr(m.Token, "missing_expr")
	}

	i := len(exprToks)
	block_toks := b.getrange(&i, lexer.KND_LBRACE, lexer.KND_RBRACE, &toks)
	if block_toks == nil {
		b.stop()
		b.pusherr(m.Token, "body_not_exist")
		return
	}

	m.Cases, m.Default = b.cases(block_toks, m.TypeMatch)

	for i := range m.Cases {
		c := &m.Cases[i]
		c.Match = m
		if i > 0 {
			m.Cases[i-1].Next = c
		}
	}
	if m.Default != nil {
		if len(m.Cases) > 0 {
			m.Cases[len(m.Cases)-1].Next = m.Default
		}
		m.Default.Match = m
	}

	s.Data = m
	return
}

func (b *builder) if_expr(bs *block_st) *ast.If {
	model := new(ast.If)
	model.Token = bs.toks[0]
	bs.toks = bs.toks[1:]
	exprToks := ast.GetBlockExpr(bs.toks)
	i := 0
	if len(exprToks) == 0 {
		b.pusherr(model.Token, "missing_expr")
	} else {
		i = len(exprToks)
	}
	blockToks := b.getrange(&i, lexer.KND_LBRACE, lexer.KND_RBRACE, &bs.toks)
	if blockToks == nil {
		b.stop()
		b.pusherr(model.Token, "body_not_exist")
		return nil
	}
	if i < len(bs.toks) {
		if bs.toks[i].Id == lexer.ID_ELSE {
			bs.nextToks = bs.toks[i:]
		} else {
			b.pusherr(bs.toks[i], "invalid_syntax")
		}
	}
	model.Expr = b.Expr(exprToks)
	model.Block = b.Block(blockToks)
	return model
}

func (b *builder) conditional_default(bs *block_st) *ast.Else {
	model := new(ast.Else)
	model.Token = bs.toks[0]
	bs.toks = bs.toks[1:]
	i := 0
	blockToks := b.getrange(&i, lexer.KND_LBRACE, lexer.KND_RBRACE, &bs.toks)
	if blockToks == nil {
		if i < len(bs.toks) {
			b.pusherr(model.Token, "else_have_expr")
		} else {
			b.stop()
			b.pusherr(model.Token, "body_not_exist")
		}
		return nil
	}
	if i < len(bs.toks) {
		b.pusherr(bs.toks[i], "invalid_syntax")
	}
	model.Block = b.Block(blockToks)
	return model
}

func (b *builder) Conditional(bs *block_st) (s ast.St) {
	s.Token = bs.toks[0]
	var c ast.Conditional
	c.If = b.if_expr(bs)
	if c.If == nil {
		return
	}
node:
	if bs.terminated || blockStFinished(bs) {
		goto end
	}
	set_to_next_st(bs)
	if bs.toks[0].Id == lexer.ID_ELSE {
		if len(bs.toks) > 1 && bs.toks[1].Id == lexer.ID_IF {
			bs.toks = bs.toks[1:]
			elif := b.if_expr(bs)
			c.Elifs = append(c.Elifs, elif)
			goto node
		}
		c.Default = b.conditional_default(bs)
	} else {
		bs.nextToks = bs.toks
	}
end:
	s.Data = c
	return
}

func (b *builder) BreakSt(toks []lexer.Token) ast.St {
	var breakAST ast.Break
	breakAST.Token = toks[0]
	if len(toks) > 1 {
		if toks[1].Id != lexer.ID_IDENT {
			b.pusherr(toks[1], "invalid_syntax")
		} else {
			breakAST.LabelToken = toks[1]
			if len(toks) > 2 {
				b.pusherr(toks[1], "invalid_syntax")
			}
		}
	}
	return ast.St{
		Token: breakAST.Token,
		Data:  breakAST,
	}
}

func (b *builder) ContinueSt(toks []lexer.Token) ast.St {
	var continueAST ast.Continue
	continueAST.Token = toks[0]
	if len(toks) > 1 {
		if toks[1].Id != lexer.ID_IDENT {
			b.pusherr(toks[1], "invalid_syntax")
		} else {
			continueAST.LoopLabel = toks[1]
			if len(toks) > 2 {
				b.pusherr(toks[1], "invalid_syntax")
			}
		}
	}
	return ast.St{Token: continueAST.Token, Data: continueAST}
}

func (b *builder) Expr(toks []lexer.Token) (e ast.Expr) {
	e.Op = b.build_expr_op(toks)
	e.Tokens = toks
	return
}

func (b *builder) build_binop_expr(toks []lexer.Token) any {
	i := b.find_lowest_precedenced_operator(toks)
	if i != -1 {
		return b.build_binop(toks)
	}
	return ast.BinopExpr{Tokens: toks}
}

func (b *builder) build_binop(toks []lexer.Token) ast.Binop {
	op := ast.Binop{}
	i := b.find_lowest_precedenced_operator(toks)
	op.L = b.build_binop_expr(toks[:i])
	op.R = b.build_binop_expr(toks[i+1:])
	op.Op = toks[i]
	return op
}

func eliminate_comments(toks []lexer.Token) []lexer.Token {
	cutted := []lexer.Token{}
	for _, token := range toks {
		if token.Id != lexer.ID_COMMENT {
			cutted = append(cutted, token)
		}
	}
	return cutted
}

func (b *builder) build_expr_op(toks []lexer.Token) any {
	toks = eliminate_comments(toks)
	i := b.find_lowest_precedenced_operator(toks)
	if i == -1 {
		return b.build_binop_expr(toks)
	}
	return b.build_binop(toks)
}

func (b *builder) find_lowest_precedenced_operator(toks []lexer.Token) int {
	prec := precedencer{}
	brace_n := 0
	for i, tok := range toks {
		switch {
		case tok.Id == lexer.ID_BRACE:
			switch tok.Kind {
			case lexer.KND_LBRACE, lexer.KND_LPAREN, lexer.KND_LBRACKET:
				brace_n++
			default:
				brace_n--
			}
			continue
		case i == 0:
			continue
		case tok.Id != lexer.ID_OP:
			continue
		case brace_n > 0:
			continue
		}
		if toks[i-1].Id == lexer.ID_OP {
			continue
		}
		p := tok.Prec()
		if p != -1 {
			prec.set(p, i)
		}
	}
	data := prec.get_lower()
	if data == nil {
		return -1
	}
	return data.(int)
}

func (b *builder) getrange(i *int, open, close string, toks *[]lexer.Token) []lexer.Token {
	rang := ast.Range(i, open, close, *toks)
	if rang != nil {
		return rang
	}
	if b.Ended() {
		return nil
	}
	*i = 0
	*toks = b.next_builder_st()
	rang = ast.Range(i, open, close, *toks)
	return rang
}

func (b *builder) skip_st(i *int, toks *[]lexer.Token) []lexer.Token {
	start := *i
	*i, _ = ast.NextStPos(*toks, start)
	stoks := (*toks)[start:*i]
	if stoks[len(stoks)-1].Id == lexer.ID_SEMICOLON {
		if len(stoks) == 1 {
			return b.skip_st(i, toks)
		}
		stoks = stoks[:len(stoks)-1]
	}
	return stoks
}

func (b *builder) next_builder_st() []lexer.Token {
	return b.skip_st(&b.Pos, &b.Tokens)
}
