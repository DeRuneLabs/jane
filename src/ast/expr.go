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

import "github.com/DeRuneLabs/jane/lexer"

func IsFnCall(toks []lexer.Token) []lexer.Token {
	switch toks[0].Id {
	case lexer.ID_BRACE, lexer.ID_IDENT, lexer.ID_DT:
	default:
		tok := toks[len(toks)-1]
		if tok.Id != lexer.ID_BRACE && tok.Kind != lexer.KND_RPARENT {
			return nil
		}
	}
	tok := toks[len(toks)-1]
	if tok.Id != lexer.ID_BRACE || tok.Kind != lexer.KND_RPARENT {
		return nil
	}
	brace_n := 0
	for i := len(toks) - 1; i >= 1; i-- {
		tok := toks[i]
		if tok.Id == lexer.ID_BRACE {
			switch tok.Kind {
			case lexer.KND_RPARENT:
				brace_n++
			case lexer.KND_LPAREN:
				brace_n--
			}
			if brace_n == 0 {
				return toks[:i]
			}
		}
	}
	return nil
}

func GetBlockExpr(toks []lexer.Token) (expr []lexer.Token) {
	brace_n := 0
	for i, tok := range toks {
		if tok.Id == lexer.ID_BRACE {
			switch tok.Kind {
			case lexer.KND_LBRACE:
				if brace_n > 0 {
					brace_n++
					break
				}
				return toks[:i]
			case lexer.KND_LBRACKET, lexer.KND_LPAREN:
				brace_n++
			default:
				brace_n--
			}
		}
	}
	return nil
}
