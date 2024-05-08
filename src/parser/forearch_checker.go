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

type foreachChecker struct {
	p       *Parser
	profile *ast.IterForeach
	val     value
}

func (fc *foreachChecker) array() {
	fc.check_size_key()
	if lexer.IsIgnoreId(fc.profile.KeyB.Id) {
		return
	}
	componentType := *fc.profile.ExprType.ComponentType
	b := &fc.profile.KeyB
	b.DataType = componentType
	val := fc.val
	val.data.DataType = componentType
	fc.p.check_valid_init_expr(b.Mutable, val, fc.profile.InToken)
}

func (fc *foreachChecker) slice() {
	fc.check_size_key()
	if lexer.IsIgnoreId(fc.profile.KeyB.Id) {
		return
	}
	componentType := *fc.profile.ExprType.ComponentType
	b := &fc.profile.KeyB
	b.DataType = componentType
	val := fc.val
	val.data.DataType = componentType
	fc.p.check_valid_init_expr(b.Mutable, val, fc.profile.InToken)
}

func (fc *foreachChecker) hashmap() {
	fc.check_map_key_a()
	fc.check_map_key_b()
}

func (fc *foreachChecker) check_size_key() {
	if lexer.IsIgnoreId(fc.profile.KeyA.Id) {
		return
	}
	a := &fc.profile.KeyA
	a.DataType.Id = types.INT
	a.DataType.Kind = types.TYPE_MAP[a.DataType.Id]
}

func (fc *foreachChecker) check_map_key_a() {
	if lexer.IsIgnoreId(fc.profile.KeyA.Id) {
		return
	}
	keyType := fc.val.data.DataType.Tag.([]Type)[0]
	a := &fc.profile.KeyA
	a.DataType = keyType
	val := fc.val
	val.data.DataType = keyType
	fc.p.check_valid_init_expr(a.Mutable, val, fc.profile.InToken)
}

func (fc *foreachChecker) check_map_key_b() {
	if lexer.IsIgnoreId(fc.profile.KeyB.Id) {
		return
	}
	valType := fc.val.data.DataType.Tag.([]Type)[1]
	b := &fc.profile.KeyB
	b.DataType = valType
	val := fc.val
	val.data.DataType = valType
	fc.p.check_valid_init_expr(b.Mutable, val, fc.profile.InToken)
}

func (fc *foreachChecker) str() {
	fc.check_size_key()
	if lexer.IsIgnoreId(fc.profile.KeyB.Id) {
		return
	}
	runeType := Type{
		Id:   types.U8,
		Kind: types.TYPE_MAP[types.U8],
	}
	b := &fc.profile.KeyB
	b.DataType = runeType
}

func (fc *foreachChecker) check() {
	switch {
	case types.IsSlice(fc.val.data.DataType):
		fc.slice()
	case types.IsArray(fc.val.data.DataType):
		fc.array()
	case types.IsMap(fc.val.data.DataType):
		fc.hashmap()
	case fc.val.data.DataType.Id == types.STR:
		fc.str()
	}
}
