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
	"github.com/DeRuneLabs/jane/lexer"
	"github.com/DeRuneLabs/jane/lexer/tokens"
	"github.com/DeRuneLabs/jane/package/jntype"
)

type unary struct {
	token lexer.Token
	toks  []lexer.Token
	model *exprModel
	t     *Transpiler
}

func (u *unary) minus() value {
	v := u.t.eval.process(u.toks, u.model)
	if !typeIsPure(v.data.Type) || !jntype.IsNumeric(v.data.Type.Id) {
		u.t.eval.pusherrtok(u.token, "invalid_expr_unary_operator", tokens.MINUS)
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
	v := u.t.eval.process(u.toks, u.model)
	if !typeIsPure(v.data.Type) || !jntype.IsNumeric(v.data.Type.Id) {
		u.t.eval.pusherrtok(u.token, "invalid_expr_unary_operator", tokens.PLUS)
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
	v := u.t.eval.process(u.toks, u.model)
	if !typeIsPure(v.data.Type) || !jntype.IsInteger(v.data.Type.Id) {
		u.t.eval.pusherrtok(u.token, "invalid_expr_unary_operator", tokens.CARET)
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
	v := u.t.eval.process(u.toks, u.model)
	if !isBoolExpr(v) {
		u.t.eval.pusherrtok(u.token, "invalid_expr_unary_operator", tokens.EXCLAMATION)
	} else if v.constExpr {
		v.expr = !v.expr.(bool)
		v.model = boolModel(v)
	}
	v.data.Type.Id = jntype.Bool
	v.data.Type.Kind = tokens.BOOL
	return v
}

func (u *unary) star() value {
	if !u.t.eval.unsafe_allowed() {
		u.t.pusherrtok(u.token, "unsafe_behavior_at_out_of_unsafe_scope")
	}
	v := u.t.eval.process(u.toks, u.model)
	v.constExpr = false
	v.lvalue = true
	switch {
	case !typeIsExplicitPtr(v.data.Type):
		u.t.eval.pusherrtok(u.token, "invalid_expr_unary_operator", tokens.STAR)
		goto end
	}
	v.data.Type.Kind = v.data.Type.Kind[1:]
end:
	v.data.Value = " "
	return v
}

func (u *unary) amper() value {
	v := u.t.eval.process(u.toks, u.model)
	v.constExpr = false
	v.lvalue = true
	nodes := &u.model.nodes[u.model.index].nodes
	switch {
	case valIsStructIns(v):
		s := v.data.Type.Tag.(*structure)
		if s.Ast.Id != v.data.Value {
			break
		}
		var alloc_model exprNode
		alloc_model.value = "__jnc_new_structure<"
		alloc_model.value += s.OutId()
		alloc_model.value += ">(new( std::nothrow ) "
		(*nodes)[0] = alloc_model
		last := &(*nodes)[len(*nodes)-1]
		*last = exprNode{(*last).String() + ")"}
		v.data.Type.Kind = tokens.AMPER + v.data.Type.Kind
		v.mutable = true
		return v
	case typeIsRef(v.data.Type):
		model := exprNode{(*nodes)[1].String() + "._alloc"}
		*nodes = nil
		*nodes = make([]iExpr, 1)
		(*nodes)[0] = model
		v.data.Type.Kind = tokens.STAR + un_ptr_or_ref_type(v.data.Type).Kind
		return v
	case !canGetPtr(v):
		u.t.eval.pusherrtok(u.token, "invalid_expr_unary_operator", tokens.AMPER)
		return v
	}
	v.data.Type.Kind = tokens.STAR + v.data.Type.Kind
	return v
}
