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
	"github.com/DeRuneLabs/jane/package/jnapi"
	"github.com/DeRuneLabs/jane/package/jntype"
)

type foreachChecker struct {
	p       *Transpiler
	profile *models.IterForeach
	val     value
}

func (fc *foreachChecker) array() {
	fc.checkKeyASize()
	if jnapi.IsIgnoreId(fc.profile.KeyB.Id) {
		return
	}
	componentType := *fc.profile.ExprType.ComponentType
	b := &fc.profile.KeyB
	b.Type = componentType
	val := fc.val
	val.data.Type = componentType
	fc.p.check_valid_init_expr(b.Mutable, val, fc.profile.InToken)
}

func (fc *foreachChecker) slice() {
	fc.checkKeyASize()
	if jnapi.IsIgnoreId(fc.profile.KeyB.Id) {
		return
	}
	componentType := *fc.profile.ExprType.ComponentType
	b := &fc.profile.KeyB
	b.Type = componentType
	val := fc.val
	val.data.Type = componentType
	fc.p.check_valid_init_expr(b.Mutable, val, fc.profile.InToken)
}

func (fc *foreachChecker) hashmap() {
	fc.checkKeyAMapKey()
	fc.checkKeyBMapVal()
}

func (fc *foreachChecker) checkKeyASize() {
	if jnapi.IsIgnoreId(fc.profile.KeyA.Id) {
		return
	}
	a := &fc.profile.KeyA
	a.Type.Id = jntype.Int
	a.Type.Kind = jntype.TypeMap[a.Type.Id]
}

func (fc *foreachChecker) checkKeyAMapKey() {
	if jnapi.IsIgnoreId(fc.profile.KeyA.Id) {
		return
	}
	keyType := fc.val.data.Type.Tag.([]Type)[0]
	a := &fc.profile.KeyA
	a.Type = keyType
	val := fc.val
	val.data.Type = keyType
	fc.p.check_valid_init_expr(a.Mutable, val, fc.profile.InToken)
}

func (fc *foreachChecker) checkKeyBMapVal() {
	if jnapi.IsIgnoreId(fc.profile.KeyB.Id) {
		return
	}
	valType := fc.val.data.Type.Tag.([]Type)[1]
	b := &fc.profile.KeyB
	b.Type = valType
	val := fc.val
	val.data.Type = valType
	fc.p.check_valid_init_expr(b.Mutable, val, fc.profile.InToken)
}

func (fc *foreachChecker) str() {
	fc.checkKeyASize()
	if jnapi.IsIgnoreId(fc.profile.KeyB.Id) {
		return
	}
	runeType := Type{
		Id:   jntype.U8,
		Kind: jntype.TypeMap[jntype.U8],
	}
	b := &fc.profile.KeyB
	b.Type = runeType
}

func (fc *foreachChecker) check() {
	switch {
	case typeIsSlice(fc.val.data.Type):
		fc.slice()
	case typeIsArray(fc.val.data.Type):
		fc.array()
	case typeIsMap(fc.val.data.Type):
		fc.hashmap()
	case fc.val.data.Type.Id == jntype.Str:
		fc.str()
	}
}
