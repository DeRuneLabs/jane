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
	"github.com/DeRuneLabs/jane/ast/models"
	"github.com/DeRuneLabs/jane/lexer/tokens"
)

type blockStatement struct {
	pos            int
	block          *models.Block
	srcToks        *Toks
	toks           Toks
	nextToks       Toks
	withTerminator bool
}

func IsStatement(current, prev Tok) (ok bool, withTerminator bool) {
	ok = current.Id == tokens.SemiColon || prev.Row < current.Row
	withTerminator = current.Id == tokens.SemiColon
	return
}

func NextStatementPos(toks Toks, start int) (int, bool) {
	braceCount := 0
	i := start
	for ; i < len(toks); i++ {
		var isStatement, withTerminator bool
		tok := toks[i]
		if tok.Id == tokens.Brace {
			switch tok.Kind {
			case tokens.LBRACE, tokens.LBRACKET, tokens.LPARENTHESES:
				if braceCount == 0 && i > start {
					isStatement, withTerminator = IsStatement(tok, toks[i-1])
					if isStatement {
						goto ret
					}
				}
				braceCount++
				continue
			default:
				braceCount--
				if braceCount == 0 && i+1 < len(toks) {
					isStatement, withTerminator = IsStatement(toks[i+1], tok)
					if isStatement {
						i++
						goto ret
					}
				}
				continue
			}
		}
		if braceCount != 0 {
			continue
		} else if i > start {
			isStatement, withTerminator = IsStatement(tok, toks[i-1])
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
