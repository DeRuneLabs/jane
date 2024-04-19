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

type AssignInfo struct {
	Left   []lexer.Token
	Right  []lexer.Token
	Setter lexer.Token
	Ok     bool
}

var PostfixOperators = [...]string{
	0: tokens.DOUBLE_PLUS,
	1: tokens.DOUBLE_MINUS,
}

var AssignOperators = [...]string{
	0:  tokens.EQUAL,
	1:  tokens.PLUS_EQUAL,
	2:  tokens.MINUS_EQUAL,
	3:  tokens.SLASH_EQUAL,
	4:  tokens.STAR_EQUAL,
	5:  tokens.PERCENT_EQUAL,
	6:  tokens.RSHIFT_EQUAL,
	7:  tokens.LSHIFT_EQUAL,
	8:  tokens.VLINE_EQUAL,
	9:  tokens.AMPER_EQUAL,
	10: tokens.CARET_EQUAL,
}

func IsAssign(id uint8) bool {
	switch id {
	case tokens.Id,
		tokens.Cpp,
		tokens.Let,
		tokens.Dot,
		tokens.Self,
		tokens.Brace,
		tokens.Operator:
		return true
	}
	return false
}

func IsPostfixOperator(kind string) bool {
	for _, operator := range PostfixOperators {
		if kind == operator {
			return true
		}
	}
	return false
}

func IsAssignOperator(kind string) bool {
	if IsPostfixOperator(kind) {
		return true
	}
	for _, operator := range AssignOperators {
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
		if t.Id == tokens.Brace {
			switch t.Kind {
			case tokens.LBRACE, tokens.LBRACKET, tokens.LPARENTHESES:
				brace_n++
			default:
				brace_n--
			}
		}
		if brace_n < 0 {
			return false
		} else if brace_n > 0 {
			continue
		} else if t.Id == tokens.Operator && IsAssignOperator(t.Kind) {
			return true
		}
	}
	return false
}
