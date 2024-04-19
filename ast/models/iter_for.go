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

import "strings"

type IterFor struct {
	Once      Statement
	Condition Expr
	Next      Statement
}

func (f IterFor) String(i *Iter) string {
	var cpp strings.Builder
	var indent string
	if f.Once.Data != nil {
		cpp.WriteString("{\n")
		AddIndent()
		indent = IndentString()
		cpp.WriteString(indent)
		cpp.WriteString(f.Once.String())
		cpp.WriteByte('\n')
		cpp.WriteString(indent)
	} else {
		indent = IndentString()
	}
	begin := i.BeginLabel()
	cpp.WriteString(begin)
	cpp.WriteString(":;\n")
	cpp.WriteString(indent)
	cpp.WriteString(i.Block.String())
	cpp.WriteByte('\n')
	cpp.WriteString(indent)
	cpp.WriteString(i.NextLabel())
	cpp.WriteString(":;\n")
	cpp.WriteString(indent)
	if f.Next.Data != nil {
		cpp.WriteString(f.Next.String())
		cpp.WriteByte('\n')
		cpp.WriteString(indent)
	}
	condition := f.Condition.String()
	if condition != "" {
		cpp.WriteString("if (")
		cpp.WriteString(f.Condition.String())
		cpp.WriteString(") { goto ")
		cpp.WriteString(begin)
		cpp.WriteString("; }\n")
	} else {
		cpp.WriteString("goto ")
		cpp.WriteString(begin)
		cpp.WriteString(";\n")
	}
	cpp.WriteString(indent)
	cpp.WriteString(i.EndLabel())
	cpp.WriteString(":;")
	if f.Once.Data != nil {
		cpp.WriteByte('\n')
		DoneIndent()
		cpp.WriteString(IndentString())
		cpp.WriteByte('}')
	}
	return cpp.String()
}
