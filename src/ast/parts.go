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

package ast

import (
	"github.com/DeRuneLabs/jane/build"
	"github.com/DeRuneLabs/jane/lexer"
)

func compiler_err(t lexer.Token, key string, args ...any) build.Log {
	return build.Log{
		Type:   build.ERR,
		Row:    t.Row,
		Column: t.Column,
		Path:   t.File.Path(),
		Text:   build.Errorf(key, args...),
	}
}

func Range(i *int, open string, close string, toks []lexer.Token) []lexer.Token {
	if *i >= len(toks) {
		return nil
	}
	tok := toks[*i]
	if tok.Id != lexer.ID_BRACE || tok.Kind != open {
		return nil
	}
	*i++
	n := 1
	start := *i
	for ; n != 0 && *i < len(toks); *i++ {
		tok := toks[*i]
		if tok.Id != lexer.ID_BRACE {
			continue
		}
		switch tok.Kind {
		case open:
			n++
		case close:
			n--
		}
	}
	return toks[start : *i-1]
}

func RangeLast(toks []lexer.Token) (cutted, cut []lexer.Token) {
	if len(toks) == 0 {
		return toks, nil
	} else if toks[len(toks)-1].Id != lexer.ID_BRACE {
		return toks, nil
	}
	brace_n := 0
	for i := len(toks) - 1; i >= 0; i-- {
		tok := toks[i]
		if tok.Id == lexer.ID_BRACE {
			switch tok.Kind {
			case lexer.KND_RBRACE, lexer.KND_RBRACKET, lexer.KND_RPARENT:
				brace_n++
				continue
			default:
				brace_n--
			}
		}
		if brace_n == 0 {
			return toks[:i], toks[i:]
		}
	}
	return toks, nil
}

func Parts(toks []lexer.Token, id uint8, exprMust bool) ([][]lexer.Token, []build.Log) {
	if len(toks) == 0 {
		return nil, nil
	}
	var parts [][]lexer.Token
	var errs []build.Log
	brace_n := 0
	last := 0
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
		if brace_n > 0 {
			continue
		}
		if tok.Id == id {
			if exprMust && i-last <= 0 {
				errs = append(errs, compiler_err(tok, "missing_expr"))
			}
			parts = append(parts, toks[last:i])
			last = i + 1
		}
	}
	if last < len(toks) {
		parts = append(parts, toks[last:])
	} else if !exprMust {
		parts = append(parts, []lexer.Token{})
	}
	return parts, errs
}

func SplitColon(toks []lexer.Token, i *int) (rangeToks []lexer.Token, colon int) {
	rangeToks = nil
	colon = -1
	brace_n := 0
	start := *i
	for ; *i < len(toks); *i++ {
		tok := toks[*i]
		if tok.Id == lexer.ID_BRACE {
			switch tok.Kind {
			case lexer.KND_LBRACE, lexer.KND_LBRACKET, lexer.KND_LPAREN:
				brace_n++
				continue
			default:
				brace_n--
			}
		}
		if brace_n == 0 {
			if start+1 > *i {
				return
			}
			rangeToks = toks[start+1 : *i]
			break
		} else if brace_n != 1 {
			continue
		}
		if colon == -1 && tok.Id == lexer.ID_COLON {
			colon = *i - start - 1
		}
	}
	return
}
