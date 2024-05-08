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

package gen

import (
	"strings"

	"github.com/DeRuneLabs/jane/ast"
)

type AnonFuncExpr struct {
	Ast *ast.Fn
}

func (af AnonFuncExpr) String() string {
	var cpp strings.Builder
	t := ast.Type{
		Token: af.Ast.Token,
		Kind:  af.Ast.TypeKind(),
		Tag:   af.Ast,
	}
	cpp.WriteString(t.FnString())
	cpp.WriteString("([=]")
	cpp.WriteString(gen_params(af.Ast.Params))
	cpp.WriteString(" mutable -> ")
	cpp.WriteString(af.Ast.RetType.String())
	cpp.WriteByte(' ')
	vars := af.Ast.RetType.Vars(af.Ast.Block)
	cpp.WriteString(gen_fn_block(vars, af.Ast.Block))
	cpp.WriteByte(')')
	return cpp.String()
}
