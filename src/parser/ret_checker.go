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

package parser

import (
	"github.com/DeRuneLabs/jane/ast"
	"github.com/DeRuneLabs/jane/lexer"
	"github.com/DeRuneLabs/jane/types"
)

type ret_checker struct {
	p         *Parser
	ret_ast   *ast.Ret
	f         *Fn
	exp_model retExpr
	values    []value
}

func (rc *ret_checker) push_val(last int, current int, errTok lexer.Token) {
	if current-last == 0 {
		rc.p.pusherrtok(errTok, "missing_expr")
		return
	}
	toks := rc.ret_ast.Expr.Tokens[last:current]
	var prefix Type
	i := len(rc.values)
	if rc.f.RetType.DataType.MultiTyped {
		types := rc.f.RetType.DataType.Tag.([]Type)
		if i < len(types) {
			prefix = types[i]
		}
	} else if i == 0 {
		prefix = rc.f.RetType.DataType
	}
	v, model := rc.p.evalToks(toks, &prefix)
	rc.exp_model.models = append(rc.exp_model.models, model)
	rc.values = append(rc.values, v)
}

func (rc *ret_checker) check_expressions() {
	brace_n := 0
	last := 0
	for i, tok := range rc.ret_ast.Expr.Tokens {
		if tok.Id == lexer.ID_BRACE {
			switch tok.Kind {
			case lexer.KND_LBRACE, lexer.KND_LBRACKET, lexer.KND_LPAREN:
				brace_n++
			default:
				brace_n--
			}
		}
		if brace_n > 0 || tok.Id != lexer.ID_COMMA {
			continue
		}
		rc.push_val(last, i, tok)
		last = i + 1
	}
	n := len(rc.ret_ast.Expr.Tokens)
	if last < n {
		if last == 0 {
			rc.push_val(0, n, rc.ret_ast.Token)
		} else {
			rc.push_val(last, n, rc.ret_ast.Expr.Tokens[last-1])
		}
	}
	if !types.IsVoid(rc.f.RetType.DataType) {
		rc.check_type_safety()
		rc.ret_ast.Expr.Model = rc.exp_model
	}
}

func (rc *ret_checker) check_for_ret_expr(v value) {
	if rc.p.unsafe_allowed() || !lexer.IsIdentifierRune(v.data.Value) {
		return
	}
	if !v.mutable && types.IsMut(v.data.DataType) {
		rc.p.pusherrtok(rc.ret_ast.Token, "ret_with_mut_typed_non_mut")
		return
	}
}

func (rc *ret_checker) single() {
	if len(rc.values) > 1 {
		rc.p.pusherrtok(rc.ret_ast.Token, "overflow_return")
	}
	v := rc.values[0]
	rc.check_for_ret_expr(v)
	assign_checker{
		p:      rc.p,
		t:      rc.f.RetType.DataType,
		v:      v,
		errtok: rc.ret_ast.Token,
	}.check()
}

func (rc *ret_checker) multi() {
	types := rc.f.RetType.DataType.Tag.([]Type)
	n := len(rc.values)
	if n == 1 {
		rc.check_multi_ret_as_mutli_ret()
		return
	} else if n > len(types) {
		rc.p.pusherrtok(rc.ret_ast.Token, "overflow_return")
	}
	for i, t := range types {
		if i >= n {
			break
		}
		v := rc.values[i]
		rc.check_for_ret_expr(v)
		assign_checker{
			p:      rc.p,
			t:      t,
			v:      v,
			errtok: rc.ret_ast.Token,
		}.check()
	}
}

func (rc *ret_checker) check_type_safety() {
	if !rc.f.RetType.DataType.MultiTyped {
		rc.single()
		return
	}
	rc.multi()
}

func (rc *ret_checker) check_multi_ret_as_mutli_ret() {
	v := rc.values[0]
	if !v.data.DataType.MultiTyped {
		rc.p.pusherrtok(rc.ret_ast.Token, "missing_multi_return")
		return
	}
	val_types := v.data.DataType.Tag.([]Type)
	ret_types := rc.f.RetType.DataType.Tag.([]Type)
	if len(val_types) < len(ret_types) {
		rc.p.pusherrtok(rc.ret_ast.Token, "missing_multi_return")
		return
	} else if len(val_types) < len(ret_types) {
		rc.p.pusherrtok(rc.ret_ast.Token, "overflow_return")
		return
	}
	for i, rt := range ret_types {
		vt := val_types[i]
		v := value{data: ast.Data{DataType: vt}}
		v.data.Value = " "
		assign_checker{
			p:                rc.p,
			t:                rt,
			v:                v,
			ignoreAny:        false,
			not_allow_assign: false,
			errtok:           rc.ret_ast.Token,
		}.check()
	}
}

func (rc *ret_checker) rets_vars() {
	if !rc.f.RetType.DataType.MultiTyped {
		for _, v := range rc.f.RetType.Identifiers {
			if !lexer.IsIgnoreId(v.Kind) {
				model := new(expr_model)
				model.index = 0
				model.nodes = make([]expr_build_node, 1)
				val, _ := rc.p.eval.single(v, model)
				rc.exp_model.models = append(rc.exp_model.models, model)
				rc.values = append(rc.values, val)
				break
			}
		}
		rc.ret_ast.Expr.Model = rc.exp_model
		return
	}
	types := rc.f.RetType.DataType.Tag.([]Type)
	for i, v := range rc.f.RetType.Identifiers {
		if lexer.IsIgnoreId(v.Kind) {
			node := exprNode{}
			node.value = types[i].String()
			node.value += types[i].InitValue()
			rc.exp_model.models = append(rc.exp_model.models, node)
			continue
		}
		model := new(expr_model)
		model.index = 0
		model.nodes = make([]expr_build_node, 1)
		val, _ := rc.p.eval.single(v, model)
		rc.exp_model.models = append(rc.exp_model.models, model)
		rc.values = append(rc.values, val)
	}
	rc.ret_ast.Expr.Model = rc.exp_model
}

func (rc *ret_checker) check() {
	n := len(rc.ret_ast.Expr.Tokens)
	if n == 0 && !types.IsVoid(rc.f.RetType.DataType) {
		if !rc.f.RetType.AnyVar() {
			rc.p.pusherrtok(rc.ret_ast.Token, "require_return_value")
		}
		rc.rets_vars()
		return
	}
	if n > 0 && types.IsVoid(rc.f.RetType.DataType) {
		rc.p.pusherrtok(rc.ret_ast.Token, "void_function_return_value")
	}
	rc.exp_model.vars = rc.f.RetType.Vars(rc.p.nodeBlock)
	rc.check_expressions()
}
