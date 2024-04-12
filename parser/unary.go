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

package parser

import (
	"github.com/DeRuneLabs/jane/lexer/tokens"
	"github.com/DeRuneLabs/jane/package/jntype"
)

type unary struct {
	tok   Tok
	toks  Toks
	model *exprModel
	p     *Parser
}

func (u *unary) minus() value {
	v := u.p.eval.process(u.toks, u.model)
	if !typeIsPure(v.data.Type) || !jntype.IsNumeric(v.data.Type.Id) {
		u.p.eval.pusherrtok(u.tok, "invalid_type_unary_operator", tokens.MINUS)
	}
	if v.constExpr {
		v.data.Value = tokens.MINUS + v.data.Value
		switch t := v.expr.(type) {
		case float64:
			v.expr = -t
		case int64:
			v.expr = -t
		case uint64:
			v.expr = -t
		}
		v.model = numericModel(v)
	}
	return v
}

func (u *unary) plus() value {
	v := u.p.eval.process(u.toks, u.model)
	if !typeIsPure(v.data.Type) || !jntype.IsNumeric(v.data.Type.Id) {
		u.p.eval.pusherrtok(u.tok, "invalid_type_unary_operator", tokens.PLUS)
	}
	if v.constExpr {
		switch t := v.expr.(type) {
		case float64:
			v.expr = +t
		case int64:
			v.expr = +t
		case uint64:
			v.expr = +t
		}
		v.model = numericModel(v)
	}
	return v
}

func (u *unary) caret() value {
	v := u.p.eval.process(u.toks, u.model)
	if !typeIsPure(v.data.Type) || !jntype.IsInteger(v.data.Type.Id) {
		u.p.eval.pusherrtok(u.tok, "invalid_type_unary_operator", tokens.CARET)
	}
	if v.constExpr {
		switch t := v.expr.(type) {
		case int64:
			v.expr = ^t
		case uint64:
			v.expr = ^t
		}
		v.model = numericModel(v)
	}
	return v
}

func (u *unary) logicalNot() value {
	v := u.p.eval.process(u.toks, u.model)
	if !isBoolExpr(v) {
		u.p.eval.pusherrtok(u.tok, "invalid_type_unary_operator", tokens.EXCLAMATION)
	} else if v.constExpr {
		v.expr = !v.expr.(bool)
		v.model = boolModel(v)
	}
	v.data.Type.Id = jntype.Bool
	v.data.Type.Kind = tokens.BOOL
	return v
}

func (u *unary) star() value {
	v := u.p.eval.process(u.toks, u.model)
	v.constExpr = false
	v.lvalue = true
	if !typeIsExplicitPtr(v.data.Type) {
		u.p.eval.pusherrtok(u.tok, "invalid_type_unary_operator", tokens.STAR)
	} else {
		v.data.Type.Kind = v.data.Type.Kind[1:]
	}
	return v
}

func (u *unary) amper() value {
	v := u.p.eval.process(u.toks, u.model)
	v.constExpr = false
	if !canGetPtr(v) {
		u.p.eval.pusherrtok(u.tok, "invalid_type_unary_operator", tokens.AMPER)
	}
	v.lvalue = true
	v.data.Type.Kind = tokens.STAR + v.data.Type.Kind
	nodes := &u.model.nodes[u.model.index].nodes
	var func_call string
	if v.heapMust {
		func_call = "__jnc_ptr_of"
	} else {
		func_call = "__jnc_not_heap_ptr_of"
	}
	expr := []iExpr{
		exprNode{func_call},
		exprNode{tokens.LPARENTHESES},
	}
	*nodes = append(expr, *nodes...)
	*nodes = append(*nodes, exprNode{tokens.RPARENTHESES})
	return v
}
