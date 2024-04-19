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

package transpiler

import (
	"github.com/DeRuneLabs/jane/ast/models"
	"github.com/DeRuneLabs/jane/lexer"
)

func make_slice(p *Transpiler, m *exprModel, t models.Type, args *models.Args, errtok lexer.Token) (v value) {
	v.data.Type = t
	v.data.Value = " "
	if len(args.Src) < 2 {
		p.pusherrtok(errtok, "missing_expr_for", "len")
		return
	} else if len(args.Src) > 2 {
		p.pusherrtok(errtok, "argument_overflow")
	}
	len_expr := args.Src[1].Expr
	len_v, len_expr_model := p.evalExpr(len_expr, nil)
	err_key := check_value_for_indexing(len_v)
	if err_key != "" {
		p.pusherrtok(errtok, err_key)
	} else if typeIsRef(*t.ComponentType) {
		p.pusherrtok(errtok, "reference_not_initialized")
	}
	m.nodes[m.index].nodes[0] = nil
	m.appendSubNode(exprNode{t.String()})
	m.appendSubNode(exprNode{"("})
	m.appendSubNode(len_expr_model)
	m.appendSubNode(exprNode{")"})
	return
}
