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
	"github.com/DeRuneLabs/jane/package/jnapi"
)

type RetType struct {
	Type        Type
	Identifiers []lexer.Token
}

func (rt RetType) String() string {
	return rt.Type.String()
}

func (rt *RetType) AnyVar() bool {
	for _, tok := range rt.Identifiers {
		if !jnapi.IsIgnoreId(tok.Kind) {
			return true
		}
	}
	return false
}

func (rt *RetType) Vars(owner *Block) []*Var {
	get := func(tok lexer.Token, t Type) *Var {
		v := new(Var)
		v.Token = tok
		if jnapi.IsIgnoreId(tok.Kind) {
			v.Id = jnapi.Ignore
		} else {
			v.Id = tok.Kind
		}
		v.Type = t
		v.Owner = owner
		v.Mutable = true
		return v
	}
	if !rt.Type.MultiTyped {
		if len(rt.Identifiers) > 0 {
			v := get(rt.Identifiers[0], rt.Type)
			if v == nil {
				return nil
			}
			return []*Var{v}
		}
		return nil
	}
	var vars []*Var
	types := rt.Type.Tag.([]Type)
	for i, tok := range rt.Identifiers {
		v := get(tok, types[i])
		if v != nil {
			vars = append(vars, v)
		}
	}
	return vars
}

type Ret struct {
	Token lexer.Token
	Expr  Expr
}

func (r Ret) String() string {
	if r.Expr.Model == nil {
		return "return;"
	}
	var cpp strings.Builder
	cpp.WriteString(r.Expr.String())
	cpp.WriteByte(';')
	return cpp.String()
}
