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
	"strings"

	"github.com/DeRuneLabs/jane/ast"
	"github.com/DeRuneLabs/jane/lexer"
	"github.com/DeRuneLabs/jane/types"
)

type type_builder struct {
	r      *builder
	t      *ast.Type
	tokens []lexer.Token
	i      *int
	err    bool
	first  int
	kind   string
	ok     bool
}

func (tb *type_builder) dt(tok lexer.Token) {
	tb.t.Token = tok
	tb.t.Id = types.TypeFromId(tb.t.Token.Kind)
	tb.kind += tb.t.Token.Kind
	tb.ok = true
}

func (tb *type_builder) unsafe_kw(tok lexer.Token) {
	tb.t.Id = types.UNSAFE
	tb.t.Token = tok
	tb.kind += tok.Kind
	tb.ok = true
}

func (tb *type_builder) op(tok lexer.Token) (imret bool) {
	switch tok.Kind {
	case lexer.KND_STAR, lexer.KND_AMPER, lexer.KND_DBL_AMPER:
		tb.kind += tok.Kind
	default:
		if tb.err {
			tb.r.pusherr(tok, "invalid_syntax")
		}
		return true
	}
	return false
}

func (tb *type_builder) function(tok lexer.Token) {
	tb.t.Token = tok
	tb.t.Id = types.FN
	f, proto_ok := tb.r.fn_prototype(tb.tokens, tb.i, false, true)
	if !proto_ok {
		tb.r.pusherr(tok, "invalid_type")
		return
	}
	*tb.i--
	tb.t.Tag = &f
	tb.kind += f.TypeKind()
	tb.ok = true
}

func (tb *type_builder) ident(tok lexer.Token) {
	tb.kind += tok.Kind
	if *tb.i+1 < len(tb.tokens) && tb.tokens[*tb.i+1].Id == lexer.ID_DBLCOLON {
		return
	}
	tb.t.Id = types.ID
	tb.t.Token = tok
	tb.ident_end()
	tb.ok = true
}

func (tb *type_builder) ident_end() {
	if *tb.i+1 >= len(tb.tokens) {
		return
	}
	*tb.i++
	tok := tb.tokens[*tb.i]
	if tok.Id != lexer.ID_BRACE || tok.Kind != lexer.KND_LBRACKET {
		*tb.i--
		return
	}
	tb.kind += "["
	var genericsStr strings.Builder
	parts := tb.ident_generics()
	generics := make([]ast.Type, len(parts))
	for i, part := range parts {
		index := 0
		t, _ := tb.r.DataType(part, &index, true)
		if index+1 < len(part) {
			tb.r.pusherr(part[index+1], "invalid_syntax")
		}
		genericsStr.WriteString(t.String())
		genericsStr.WriteByte(',')
		generics[i] = t
	}
	tb.kind += genericsStr.String()[:genericsStr.Len()-1] + "]"
	tb.t.Tag = generics
}

func (tb *type_builder) ident_generics() [][]lexer.Token {
	first := *tb.i
	brace_n := 0
	for ; *tb.i < len(tb.tokens); *tb.i++ {
		tok := tb.tokens[*tb.i]
		if tok.Id == lexer.ID_BRACE {
			switch tok.Kind {
			case lexer.KND_LBRACKET:
				brace_n++
			case lexer.KND_RBRACKET:
				brace_n--
			}
		}
		if brace_n == 0 {
			break
		}
	}
	tokens := tb.tokens[first+1 : *tb.i]
	parts, errs := ast.Parts(tokens, lexer.ID_COMMA, true)
	tb.r.Errors = append(tb.r.Errors, errs...)
	return parts
}

func (tb *type_builder) cpp_kw(tok lexer.Token) (imret bool) {
	if *tb.i+1 >= len(tb.tokens) {
		if tb.err {
			tb.r.pusherr(tok, "invalid_syntax")
		}
		return true
	}
	*tb.i++
	if tb.tokens[*tb.i].Id != lexer.ID_DOT {
		if tb.err {
			tb.r.pusherr(tb.tokens[*tb.i], "invalid_syntax")
		}
	}
	if *tb.i+1 >= len(tb.tokens) {
		if tb.err {
			tb.r.pusherr(tok, "invalid_syntax")
		}
		return true
	}
	*tb.i++
	if tb.tokens[*tb.i].Id != lexer.ID_IDENT {
		if tb.err {
			tb.r.pusherr(tb.tokens[*tb.i], "invalid_syntax")
		}
	}
	tb.t.CppLinked = true
	tb.t.Id = types.ID
	tb.t.Token = tb.tokens[*tb.i]
	tb.kind += tb.t.Token.Kind
	tb.ident_end()
	tb.ok = true
	return false
}

func (tb *type_builder) enumerable(tok lexer.Token) (imret bool) {
	*tb.i++
	if *tb.i >= len(tb.tokens) {
		if tb.err {
			tb.r.pusherr(tok, "invalid_syntax")
		}
		return
	}
	tok = tb.tokens[*tb.i]
	if tok.Id == lexer.ID_BRACE && tok.Kind == lexer.KND_RBRACKET {
		tb.kind += lexer.PREFIX_SLICE
		tb.t.ComponentType = new(ast.Type)
		tb.t.Id = types.SLICE
		tb.t.Token = tok
		*tb.i++
		if *tb.i >= len(tb.tokens) {
			if tb.err {
				tb.r.pusherr(tok, "invalid_syntax")
			}
			return
		}
		*tb.t.ComponentType, tb.ok = tb.r.DataType(tb.tokens, tb.i, tb.err)
		tb.kind += tb.t.ComponentType.Kind
		return false
	}
	*tb.i--
	tb.ok = tb.map_or_array()
	if tb.t.Id == types.VOID {
		if tb.err {
			tb.r.pusherr(tok, "invalid_syntax")
		}
		return true
	}
	tb.t.Token = tok
	return false
}

func (tb *type_builder) array() (ok bool) {
	defer func() { tb.t.Original = *tb.t }()
	if *tb.i+1 >= len(tb.tokens) {
		return
	}
	tb.t.Id = types.ARRAY
	*tb.i++
	exprI := *tb.i
	tb.t.ComponentType = new(ast.Type)
	ok = tb.r.datatype(tb.t.ComponentType, tb.tokens, tb.i, tb.err)
	if !ok {
		return
	}
	_, exprToks := ast.RangeLast(tb.tokens[:exprI])
	exprToks = exprToks[1 : len(exprToks)-1]
	tok := exprToks[0]
	if len(exprToks) == 1 && tok.Id == lexer.ID_OP && tok.Kind == lexer.KND_TRIPLE_DOT {
		tb.t.Size.AutoSized = true
		tb.t.Size.Expr.Tokens = exprToks
	} else {
		tb.t.Size.Expr = tb.r.Expr(exprToks)
	}
	tb.kind = tb.kind + lexer.PREFIX_ARRAY + tb.t.ComponentType.Kind
	return
}

func (tb *type_builder) map_or_array() (ok bool) {
	ok = tb.map_t()
	if !ok {
		ok = tb.array()
	}
	return
}

func (tb *type_builder) map_t() (ok bool) {
	typeToks, colon := ast.SplitColon(tb.tokens, tb.i)
	if typeToks == nil || colon == -1 {
		return
	}
	defer func() { tb.t.Original = *tb.t }()
	tb.t.Id = types.MAP
	tb.t.Token = tb.tokens[0]
	colonTok := tb.tokens[colon]
	if colon == 0 || colon+1 >= len(typeToks) {
		if tb.err {
			tb.r.pusherr(colonTok, "missing_expr")
		}
		return
	}
	keyTypeToks := typeToks[:colon]
	valueTypeToks := typeToks[colon+1:]
	types := make([]ast.Type, 2)
	j := 0
	types[0], _ = tb.r.DataType(keyTypeToks, &j, tb.err)
	j = 0
	types[1], _ = tb.r.DataType(valueTypeToks, &j, tb.err)
	tb.t.Tag = types
	tb.kind = tb.kind + tb.t.MapKind()
	ok = true
	return
}

func (tb *type_builder) step() (imret bool) {
	tok := tb.tokens[*tb.i]
	switch tok.Id {
	case lexer.ID_DT:
		tb.dt(tok)
		return
	case lexer.ID_IDENT:
		tb.ident(tok)
		return
	case lexer.ID_CPP:
		imret = tb.cpp_kw(tok)
		return
	case lexer.ID_DBLCOLON:
		tb.kind += tok.Kind
		return
	case lexer.ID_UNSAFE:
		if *tb.i+1 >= len(tb.tokens) || tb.tokens[*tb.i+1].Id != lexer.ID_FN {
			tb.unsafe_kw(tok)
			return
		}
		fallthrough
	case lexer.ID_FN:
		tb.function(tok)
		return
	case lexer.ID_OP:
		imret = tb.op(tok)
		return
	case lexer.ID_BRACE:
		switch tok.Kind {
		case lexer.KND_LBRACKET:
			imret = tb.enumerable(tok)
			return
		}
		imret = true
		return
	default:
		if tb.err {
			tb.r.pusherr(tok, "invalid_syntax")
		}
		imret = true
		return
	}
}

func (tb *type_builder) build() bool {
	defer func() { tb.t.Original = *tb.t }()
	tb.first = *tb.i
	for ; *tb.i < len(tb.tokens); *tb.i++ {
		imret := tb.step()
		if tb.ok {
			break
		} else if imret {
			return tb.ok
		}
	}
	tb.t.Kind = tb.kind
	return tb.ok
}
