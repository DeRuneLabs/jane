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
	"github.com/DeRuneLabs/jane/ast/models"
	"github.com/DeRuneLabs/jane/package/jnapi"
	"github.com/DeRuneLabs/jane/package/jntype"
)

type foreachChecker struct {
	p       *Parser
	profile *models.IterForeach
	val     value
}

func (fc *foreachChecker) array() {
	fc.checkKeyASize()
	if jnapi.IsIgnoreId(fc.profile.KeyB.Id) {
		return
	}
	componentType := *fc.profile.ExprType.ComponentType
	keyB := &fc.profile.KeyB
	if keyB.Type.Id == jntype.Void {
		keyB.Type = componentType
		return
	}
	fc.p.checkType(componentType, keyB.Type, true, fc.profile.InTok)
}

func (fc *foreachChecker) slice() {
	fc.checkKeyASize()
	if jnapi.IsIgnoreId(fc.profile.KeyB.Id) {
		return
	}
	componentType := *fc.profile.ExprType.ComponentType
	keyB := &fc.profile.KeyB
	if keyB.Type.Id == jntype.Void {
		keyB.Type = componentType
		return
	}
	fc.p.checkType(componentType, keyB.Type, true, fc.profile.InTok)
}

func (fc *foreachChecker) xmap() {
	fc.checkKeyAMapKey()
	fc.checkKeyBMapVal()
}

func (fc *foreachChecker) checkKeyASize() {
	if jnapi.IsIgnoreId(fc.profile.KeyA.Id) {
		return
	}
	keyA := &fc.profile.KeyA
	if keyA.Type.Id == jntype.Void {
		keyA.Type.Id = jntype.UInt
		keyA.Type.Kind = jntype.CppId(keyA.Type.Id)
		return
	}
	var ok bool
	keyA.Type, ok = fc.p.realType(keyA.Type, true)
	if ok {
		if !typeIsPure(keyA.Type) || !jntype.IsNumeric(keyA.Type.Id) {
			fc.p.pusherrtok(keyA.Token, "incompatible_datatype",
				keyA.Type.Kind, jntype.NumericTypeStr)
		}
	}
}

func (fc *foreachChecker) checkKeyAMapKey() {
	if jnapi.IsIgnoreId(fc.profile.KeyA.Id) {
		return
	}
	keyType := fc.val.data.Type.Tag.([]DataType)[0]
	keyA := &fc.profile.KeyA
	if keyA.Type.Id == jntype.Void {
		keyA.Type = keyType
		return
	}
	fc.p.checkType(keyType, keyA.Type, true, fc.profile.InTok)
}

func (fc *foreachChecker) checkKeyBMapVal() {
	if jnapi.IsIgnoreId(fc.profile.KeyB.Id) {
		return
	}
	valType := fc.val.data.Type.Tag.([]DataType)[1]
	keyB := &fc.profile.KeyB
	if keyB.Type.Id == jntype.Void {
		keyB.Type = valType
		return
	}
	fc.p.checkType(valType, keyB.Type, true, fc.profile.InTok)
}

func (fc *foreachChecker) str() {
	fc.checkKeyASize()
	if jnapi.IsIgnoreId(fc.profile.KeyB.Id) {
		return
	}
	runeType := DataType{
		Id:   jntype.U8,
		Kind: jntype.CppId(jntype.U8),
	}
	keyB := &fc.profile.KeyB
	if keyB.Type.Id == jntype.Void {
		keyB.Type = runeType
		return
	}
	fc.p.checkType(runeType, keyB.Type, true, fc.profile.InTok)
}

func (fc *foreachChecker) check() {
	switch {
	case typeIsSlice(fc.val.data.Type):
		fc.slice()
	case typeIsArray(fc.val.data.Type):
		fc.array()
	case typeIsMap(fc.val.data.Type):
		fc.xmap()
	case fc.val.data.Type.Id == jntype.Str:
		fc.str()
	}
}
