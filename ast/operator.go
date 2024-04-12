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

import "github.com/DeRuneLabs/jane/lexer/tokens"

var UnaryOperators = [...]string{
	0: tokens.MINUS,
	1: tokens.PLUS,
	2: tokens.CARET,
	3: tokens.EXCLAMATION,
	4: tokens.STAR,
	5: tokens.AMPER,
}

var SolidOperators = [...]string{
	0:  tokens.PLUS,
	1:  tokens.MINUS,
	2:  tokens.STAR,
	3:  tokens.SOLIDUS,
	4:  tokens.PERCENT,
	5:  tokens.AMPER,
	6:  tokens.VLINE,
	7:  tokens.CARET,
	8:  tokens.LESS,
	9:  tokens.GREAT,
	10: tokens.EXCLAMATION,
}

var ExpressionOperators = [...]string{
	0: tokens.TRIPLE_DOT,
	1: tokens.COLON,
}

func IsUnaryOperator(kind string) bool {
	return existOperator(kind, UnaryOperators[:])
}

func IsSolidOperator(kind string) bool {
	return existOperator(kind, SolidOperators[:])
}

func IsExprOperator(kind string) bool {
	return existOperator(kind, ExpressionOperators[:])
}

func existOperator(kind string, operators []string) bool {
	for _, operator := range operators {
		if kind == operator {
			return true
		}
	}
	return false
}
