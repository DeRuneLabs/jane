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
	"strconv"
	"strings"
)

type Fallthrough struct {
	Tok  Tok
	Case *Case
}

func (f Fallthrough) String() string {
	var cpp strings.Builder
	cpp.WriteString("goto ")
	cpp.WriteString(f.Case.Next.BeginLabel())
	cpp.WriteByte(';')
	return cpp.String()
}

type Case struct {
	Tok   Tok
	Exprs []Expr
	Block *Block
	Match *Match
	Next  *Case
}

func (c *Case) BeginLabel() string {
	var cpp strings.Builder
	cpp.WriteString("case_begin_")
	cpp.WriteString(strconv.FormatInt(int64(c.Tok.Row), 10))
	cpp.WriteString(strconv.FormatInt(int64(c.Tok.Column), 10))
	return cpp.String()
}

func (c *Case) EndLabel() string {
	var cpp strings.Builder
	cpp.WriteString("case_end_")
	cpp.WriteString(strconv.FormatInt(int64(c.Tok.Row), 10))
	cpp.WriteString(strconv.FormatInt(int64(c.Tok.Column), 10))
	return cpp.String()
}

func (c *Case) String(matchExpr string) string {
	endlabel := c.EndLabel()
	var cpp strings.Builder
	if len(c.Exprs) > 0 {
		cpp.WriteString("if (!(")
		for i, expr := range c.Exprs {
			cpp.WriteString(expr.String())
			if matchExpr != "" {
				cpp.WriteString(" == ")
				cpp.WriteString(matchExpr)
			}
			if i+1 < len(c.Exprs) {
				cpp.WriteString(" || ")
			}
		}
		cpp.WriteString(")) { goto ")
		cpp.WriteString(endlabel)
		cpp.WriteString("; }\n")
	}
	if len(c.Block.Tree) > 0 {
		cpp.WriteString(IndentString())
		cpp.WriteString(c.BeginLabel())
		cpp.WriteString(":;\n")
		cpp.WriteString(IndentString())
		cpp.WriteString(c.Block.String())
		cpp.WriteByte('\n')
		cpp.WriteString(IndentString())
		cpp.WriteString("goto ")

		cpp.WriteString(c.Match.EndLabel())
		cpp.WriteString(";")
		cpp.WriteByte('\n')

	}
	cpp.WriteString(IndentString())
	cpp.WriteString(endlabel)
	cpp.WriteString(":;")
	return cpp.String()
}

type Match struct {
	Tok      Tok
	Expr     Expr
	ExprType DataType
	Default  *Case
	Cases    []Case
}

func (m *Match) MatchExprString() string {
	if len(m.Cases) == 0 {
		if m.Default != nil {
			return m.Default.String("")
		}
		return ""
	}
	var cpp strings.Builder
	cpp.WriteString("{\n")
	AddIndent()
	cpp.WriteString(IndentString())
	cpp.WriteString(m.ExprType.String())
	cpp.WriteString(" expr{")
	cpp.WriteString(m.Expr.String())
	cpp.WriteString("};\n")
	cpp.WriteString(IndentString())
	if len(m.Cases) > 0 {
		cpp.WriteString(m.Cases[0].String("expr"))
		for _, c := range m.Cases[1:] {
			cpp.WriteByte('\n')
			cpp.WriteString(IndentString())
			cpp.WriteString(c.String("expr"))
		}
	}
	if m.Default != nil {
		cpp.WriteString(m.Default.String(""))
	}
	cpp.WriteByte('\n')
	DoneIndent()
	cpp.WriteString(IndentString())
	cpp.WriteByte('}')
	return cpp.String()
}

func (m *Match) MatchBoolString() string {
	var cpp strings.Builder
	if len(m.Cases) > 0 {
		cpp.WriteString(m.Cases[0].String(""))
		for _, c := range m.Cases[1:] {
			cpp.WriteByte('\n')
			cpp.WriteString(IndentString())
			cpp.WriteString(c.String(""))
		}
	}
	if m.Default != nil {
		cpp.WriteByte('\n')
		cpp.WriteString(m.Default.String(""))
		cpp.WriteByte('\n')
	}
	return cpp.String()
}

func (m *Match) EndLabel() string {
	var cpp strings.Builder
	cpp.WriteString("match_end_")

	cpp.WriteString(strconv.FormatInt(int64(m.Tok.Row), 10))
	cpp.WriteString(strconv.FormatInt(int64(m.Tok.Column), 10))
	return cpp.String()
}

func (m Match) String() string {
	var cpp strings.Builder
	if m.Expr.Model != nil {
		cpp.WriteString(m.MatchExprString())
	} else {
		cpp.WriteString(m.MatchBoolString())
	}
	cpp.WriteByte('\n')
	cpp.WriteString(IndentString())
	cpp.WriteString(m.EndLabel())
	cpp.WriteString(":;")
	return cpp.String()
}
