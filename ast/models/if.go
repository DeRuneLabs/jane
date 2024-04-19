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

package models

import (
	"strings"

	"github.com/DeRuneLabs/jane/lexer"
)

type If struct {
	Token lexer.Token
	Expr  Expr
	Block *Block
}

func (ifast If) String() string {
	var cpp strings.Builder
	cpp.WriteString("if (")
	cpp.WriteString(ifast.Expr.String())
	cpp.WriteString(") ")
	cpp.WriteString(ifast.Block.String())
	return cpp.String()
}

type ElseIf struct {
	Token lexer.Token
	Expr  Expr
	Block *Block
}

func (elif ElseIf) String() string {
	var cpp strings.Builder
	cpp.WriteString("else if (")
	cpp.WriteString(elif.Expr.String())
	cpp.WriteString(") ")
	cpp.WriteString(elif.Block.String())
	return cpp.String()
}

type Else struct {
	Token lexer.Token
	Block *Block
}

func (elseast Else) String() string {
	var cpp strings.Builder
	cpp.WriteString("else ")
	cpp.WriteString(elseast.Block.String())
	return cpp.String()
}
