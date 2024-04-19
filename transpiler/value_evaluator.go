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
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/DeRuneLabs/jane/ast/models"
	"github.com/DeRuneLabs/jane/lexer"
	"github.com/DeRuneLabs/jane/lexer/tokens"
	"github.com/DeRuneLabs/jane/package/jnapi"
	"github.com/DeRuneLabs/jane/package/jntype"
)

type valueEvaluator struct {
	token lexer.Token
	model *exprModel
	t     *Transpiler
}

func strModel(v value) iExpr {
	content := v.expr.(string)
	if israwstr(content) {
		return exprNode{jnapi.ToRawStr([]byte(content))}
	}
	return exprNode{jnapi.ToStr([]byte(content))}
}

func boolModel(v value) iExpr {
	if v.expr.(bool) {
		return exprNode{tokens.TRUE}
	}
	return exprNode{tokens.FALSE}
}

func getModel(v value) iExpr {
	switch v.expr.(type) {
	case string:
		return strModel(v)
	case bool:
		return boolModel(v)
	default:
		return numericModel(v)
	}
}

func numericModel(v value) iExpr {
	switch t := v.expr.(type) {
	case uint64:
		fmt := strconv.FormatUint(t, 10)
		return exprNode{fmt + "LLU"}
	case int64:
		fmt := strconv.FormatInt(t, 10)
		return exprNode{fmt + "LL"}
	case float64:
		switch {
		case normalize(&v):
			return numericModel(v)
		case v.data.Type.Id == jntype.F32:
			return exprNode{fmt.Sprint(t) + "f"}
		case v.data.Type.Id == jntype.F64:
			return exprNode{fmt.Sprint(t)}
		}
	}
	return nil
}

func (ve *valueEvaluator) str() value {
	var v value
	v.constExpr = true
	v.data.Value = ve.token.Kind
	v.data.Type.Id = jntype.Str
	v.data.Type.Kind = jntype.TypeMap[v.data.Type.Id]
	content := ve.token.Kind[1 : len(ve.token.Kind)-1]
	v.expr = content
	v.model = strModel(v)
	ve.model.appendSubNode(v.model)
	return v
}

func toCharLiteral(kind string) (string, bool) {
	kind = kind[1 : len(kind)-1]
	isByte := false
	switch {
	case len(kind) == 1 && kind[0] <= 255:
		isByte = true
	case kind[0] == '\\' && kind[1] == 'x':
		isByte = true
	case kind[0] == '\\' && kind[1] >= '0' && kind[1] <= '7':
		isByte = true
	}
	return kind, isByte
}

func (ve *valueEvaluator) char() value {
	var v value
	v.constExpr = true
	v.data.Value = ve.token.Kind
	content, isByte := toCharLiteral(ve.token.Kind)
	if isByte {
		v.data.Type.Id = jntype.U8
	} else { // rune
		v.data.Type.Id = jntype.I32
	}
	content = jnapi.ToRune([]byte(content))
	v.data.Type.Kind = jntype.TypeMap[v.data.Type.Id]
	v.expr, _ = strconv.ParseInt(content[2:], 16, 64)
	v.model = exprNode{content}
	ve.model.appendSubNode(v.model)
	return v
}

func (ve *valueEvaluator) bool() value {
	var v value
	v.constExpr = true
	v.data.Value = ve.token.Kind
	v.data.Type.Id = jntype.Bool
	v.data.Type.Kind = jntype.TypeMap[v.data.Type.Id]
	v.expr = ve.token.Kind == tokens.TRUE
	v.model = boolModel(v)
	ve.model.appendSubNode(v.model)
	return v
}

func (ve *valueEvaluator) nil() value {
	var v value
	v.constExpr = true
	v.data.Value = ve.token.Kind
	v.data.Type.Id = jntype.Nil
	v.data.Type.Kind = jntype.TypeMap[v.data.Type.Id]
	v.expr = nil
	v.model = exprNode{ve.token.Kind}
	ve.model.appendSubNode(v.model)
	return v
}

func normalize(v *value) (normalized bool) {
	switch {
	case !v.constExpr:
		return
	case integerAssignable(jntype.U64, *v):
		v.data.Type.Id = jntype.U64
		v.data.Type.Kind = jntype.TypeMap[v.data.Type.Id]
		v.expr = tonumu(v.expr)
		bitize(v)
		return true
	case integerAssignable(jntype.I64, *v):
		v.data.Type.Id = jntype.I64
		v.data.Type.Kind = jntype.TypeMap[v.data.Type.Id]
		v.expr = tonums(v.expr)
		bitize(v)
		return true
	}
	return
}

func (ve *valueEvaluator) float() value {
	var v value
	v.data.Value = ve.token.Kind
	v.data.Type.Id = jntype.F64
	v.data.Type.Kind = jntype.TypeMap[v.data.Type.Id]
	v.expr, _ = strconv.ParseFloat(v.data.Value, 64)
	return v
}

func (ve *valueEvaluator) integer() value {
	var v value
	v.data.Value = ve.token.Kind
	var bigint big.Int
	switch {
	case strings.HasPrefix(ve.token.Kind, "0x"):
		_, _ = bigint.SetString(ve.token.Kind[2:], 16)
	case strings.HasPrefix(ve.token.Kind, "0b"):
		_, _ = bigint.SetString(ve.token.Kind[2:], 2)
	case ve.token.Kind[0] == '0':
		_, _ = bigint.SetString(ve.token.Kind[1:], 8)
	default:
		_, _ = bigint.SetString(ve.token.Kind, 10)
	}
	if bigint.IsInt64() {
		v.expr = bigint.Int64()
	} else {
		v.expr = bigint.Uint64()
	}
	bitize(&v)
	return v
}

func (ve *valueEvaluator) numeric() value {
	var v value
	if isfloat(ve.token.Kind) {
		v = ve.float()
	} else {
		v = ve.integer()
	}
	v.constExpr = true
	v.model = numericModel(v)
	ve.model.appendSubNode(v.model)
	return v
}

func make_value_from_var(v *Var) (val value) {
	val.data.Value = v.Id
	val.data.Type = v.Type
	val.constExpr = v.Const
	val.data.Token = v.Token
	val.lvalue = !val.constExpr
	val.mutable = v.Mutable
	if val.constExpr {
		val.expr = v.ExprTag
		val.model = v.Expr.Model
	}
	return
}

func (ve *valueEvaluator) varId(id string, variable *Var, global bool) (v value) {
	variable.Used = true
	v = make_value_from_var(variable)
	if v.constExpr {
		ve.model.appendSubNode(v.model)
	} else {
		if variable.Id == tokens.SELF && !typeIsRef(variable.Type) {
			ve.model.appendSubNode(exprNode{"(*this)"})
		} else {
			ve.model.appendSubNode(exprNode{variable.OutId()})
		}
		ve.t.eval.has_error = ve.t.eval.has_error || typeIsVoid(v.data.Type)
	}
	return
}

func make_value_from_fn(f *models.Fn) (v value) {
	v.data.Value = f.Id
	v.data.Type.Id = jntype.Fn
	v.data.Type.Tag = f
	v.data.Type.Kind = f.DataTypeString()
	v.data.Token = f.Token
	return
}

func (ve *valueEvaluator) funcId(id string, f *Fn) (v value) {
	f.used = true
	v = make_value_from_fn(f.Ast)
	ve.model.appendSubNode(exprNode{f.outId()})
	return
}

func (ve *valueEvaluator) enumId(id string, e *Enum) (v value) {
	e.Used = true
	v.data.Value = id
	v.data.Type.Id = jntype.Enum
	v.data.Type.Kind = e.Id
	v.data.Type.Tag = e
	v.data.Token = e.Token
	v.constExpr = true
	v.is_type = true
	if e.Token.Id == tokens.NA {
		ve.model.appendSubNode(exprNode{jnapi.OutId(id, nil)})
	} else {
		ve.model.appendSubNode(exprNode{jnapi.OutId(id, e.Token.File)})
	}
	return
}

func make_value_from_struct(s *structure) (v value) {
	v.data.Value = s.Ast.Id
	v.data.Type.Id = jntype.Struct
	v.data.Type.Tag = s
	v.data.Type.Kind = s.Ast.Id
	v.data.Type.Token = s.Ast.Token
	v.data.Token = s.Ast.Token
	v.is_type = true
	return
}

func (ve *valueEvaluator) structId(id string, s *structure) (v value) {
	s.Used = true
	v = make_value_from_struct(s)
	if s.Ast.Token.Id == tokens.NA {
		ve.model.appendSubNode(exprNode{jnapi.OutId(id, nil)})
	} else {
		ve.model.appendSubNode(exprNode{jnapi.OutId(id, s.Ast.Token.File)})
	}
	return
}

func (ve *valueEvaluator) typeId(id string, t *TypeAlias) (_ value, _ bool) {
	dt, ok := ve.t.realType(t.Type, true)
	if !ok {
		return
	}
	if typeIsStruct(dt) {
		return ve.structId(id, dt.Tag.(*structure)), true
	}
	return
}

func (ve *valueEvaluator) id() (_ value, ok bool) {
	id := ve.token.Kind

	v, _ := ve.t.blockVarById(id)
	if v != nil {
		return ve.varId(id, v, false), true
	} else {
		v, _, _ := ve.t.globalById(id)
		if v != nil {
			return ve.varId(id, v, true), true
		}
	}

	f, _, _ := ve.t.FuncById(id)
	if f != nil {
		return ve.funcId(id, f), true
	}

	e, _, _ := ve.t.enumById(id)
	if e != nil {
		return ve.enumId(id, e), true
	}

	s, _, _ := ve.t.structById(id)
	if s != nil {
		return ve.structId(id, s), true
	}

	t, _, _ := ve.t.typeById(id)
	if t != nil {
		return ve.typeId(id, t)
	}

	ve.t.eval.pusherrtok(ve.token, "id_not_exist", id)
	return
}
