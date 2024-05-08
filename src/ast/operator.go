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

var UNARY_OPS = [...]string{
	lexer.KND_MINUS,
	lexer.KND_PLUS,
	lexer.KND_CARET,
	lexer.KND_EXCL,
	lexer.KND_STAR,
	lexer.KND_AMPER,
}

var STRONG_OPS = [...]string{
	lexer.KND_PLUS,
	lexer.KND_MINUS,
	lexer.KND_STAR,
	lexer.KND_SOLIDUS,
	lexer.KND_PERCENT,
	lexer.KND_AMPER,
	lexer.KND_VLINE,
	lexer.KND_CARET,
	lexer.KND_LT,
	lexer.KND_GT,
	lexer.KND_EXCL,
	lexer.KND_DBL_AMPER,
	lexer.KND_DBL_VLINE,
}

var WEAK_OPS = [...]string{
	lexer.KND_TRIPLE_DOT,
	lexer.KND_COLON,
}

func IsUnaryOp(kind string) bool {
	return exist_op(kind, UNARY_OPS[:])
}

func IsStrongOp(kind string) bool {
	return exist_op(kind, STRONG_OPS[:])
}

func IsExprOp(kind string) bool {
	return exist_op(kind, WEAK_OPS[:])
}

func exist_op(kind string, operators []string) bool {
	for _, operator := range operators {
		if kind == operator {
			return true
		}
	}
	return false
}
