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

type AssignInfo struct {
	Left   []lexer.Token
	Right  []lexer.Token
	Setter lexer.Token
	Ok     bool
}

var POSTFIX_OPS = [...]string{
	lexer.KND_DBL_PLUS,
	lexer.KND_DBL_MINUS,
}

var ASSING_OPS = [...]string{
	lexer.KND_EQ,
	lexer.KND_PLUS_EQ,
	lexer.KND_MINUS_EQ,
	lexer.KND_SOLIDUS_EQ,
	lexer.KND_STAR_EQ,
	lexer.KND_PERCENT_EQ,
	lexer.KND_RSHIFT_EQ,
	lexer.KND_LSHIFT_EQ,
	lexer.KND_VLINE_EQ,
	lexer.KND_AMPER_EQ,
	lexer.KND_CARET_EQ,
}

func IsAssign(id uint8) bool {
	switch id {
	case lexer.ID_IDENT,
		lexer.ID_CPP,
		lexer.ID_LET,
		lexer.ID_DOT,
		lexer.ID_SELF,
		lexer.ID_BRACE,
		lexer.ID_OP:
		return true
	}
	return false
}

func IsPostfixOp(kind string) bool {
	for _, operator := range POSTFIX_OPS {
		if kind == operator {
			return true
		}
	}
	return false
}

func IsAssignOp(kind string) bool {
	if IsPostfixOp(kind) {
		return true
	}
	for _, operator := range ASSING_OPS {
		if kind == operator {
			return true
		}
	}
	return false
}

func CheckAssignTokens(toks []lexer.Token) bool {
	if len(toks) == 0 || !IsAssign(toks[0].Id) {
		return false
	}
	brace_n := 0
	for _, t := range toks {
		if t.Id == lexer.ID_BRACE {
			switch t.Kind {
			case lexer.KND_LBRACE, lexer.KND_LBRACKET, lexer.KND_LPAREN:
				brace_n++
			default:
				brace_n--
			}
		}
		if brace_n < 0 {
			return false
		} else if brace_n > 0 {
			continue
		} else if t.Id == lexer.ID_OP && IsAssignOp(t.Kind) {
			return true
		}
	}
	return false
}
