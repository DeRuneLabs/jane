// Copyright (c) 2024 - DeRuneLabs
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
	"github.com/DeRuneLabs/jane/lexer"
	"github.com/DeRuneLabs/jane/lexer/tokens"
)

func IsFuncCall(toks []lexer.Token) []lexer.Token {
	switch toks[0].Id {
	case tokens.Brace, tokens.Id, tokens.DataType:
	default:
		tok := toks[len(toks)-1]
		if tok.Id != tokens.Brace && tok.Kind != tokens.RPARENTHESES {
			return nil
		}
	}
	tok := toks[len(toks)-1]
	if tok.Id != tokens.Brace || tok.Kind != tokens.RPARENTHESES {
		return nil
	}
	brace_n := 0
	for i := len(toks) - 1; i >= 1; i-- {
		tok := toks[i]
		if tok.Id == tokens.Brace {
			switch tok.Kind {
			case tokens.RPARENTHESES:
				brace_n++
			case tokens.LPARENTHESES:
				brace_n--
			}
			if brace_n == 0 {
				return toks[:i]
			}
		}
	}
	return nil
}

func RequireOperatorToProcess(tok lexer.Token, index, len int) bool {
	switch tok.Id {
	case tokens.Comma:
		return false
	case tokens.Brace:
		if tok.Kind == tokens.LPARENTHESES ||
			tok.Kind == tokens.LBRACE {
			return false
		}
	}
	return index < len-1
}

func BlockExpr(toks []lexer.Token) (expr []lexer.Token) {
	brace_n := 0
	for i, tok := range toks {
		if tok.Id == tokens.Brace {
			switch tok.Kind {
			case tokens.LBRACE:
				if brace_n > 0 {
					brace_n++
					break
				}
				return toks[:i]
			case tokens.LBRACKET, tokens.LPARENTHESES:
				brace_n++
			default:
				brace_n--
			}
		}
	}
	return nil
}
