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
	elementType := fc.profile.ExprType
	elementType.Kind = elementType.Kind[2:]
	keyB := &fc.profile.KeyB
	if keyB.Type.Id == jntype.Void {
		keyB.Type = elementType
		return
	}
	fc.p.wg.Add(1)
	go fc.p.checkTypeAsync(elementType, keyB.Type, true, fc.profile.InTok)
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
		keyA.Type.Kind = jntype.CxxTypeIdFromType(keyA.Type.Id)
		return
	}
	var ok bool
	keyA.Type, ok = fc.p.realType(keyA.Type, true)
	if ok {
		if !typeIsPure(keyA.Type) || !jntype.IsNumericType(keyA.Type.Id) {
			fc.p.pusherrtok(keyA.IdTok, "incompatible_datatype",
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
	fc.p.wg.Add(1)
	go fc.p.checkTypeAsync(keyType, keyA.Type, true, fc.profile.InTok)
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
	fc.p.wg.Add(1)
	go fc.p.checkTypeAsync(valType, keyB.Type, true, fc.profile.InTok)
}

func (fc *foreachChecker) str() {
	fc.checkKeyASize()
	if jnapi.IsIgnoreId(fc.profile.KeyB.Id) {
		return
	}
	runeType := DataType{
		Id:   jntype.U8,
		Kind: jntype.CxxTypeIdFromType(jntype.U8),
	}
	keyB := &fc.profile.KeyB
	if keyB.Type.Id == jntype.Void {
		keyB.Type = runeType
		return
	}
	fc.p.wg.Add(1)
	go fc.p.checkTypeAsync(runeType, keyB.Type, true, fc.profile.InTok)
}

func (fc *foreachChecker) check() {
	switch {
	case typeIsArray(fc.val.data.Type):
		fc.array()
	case typeIsMap(fc.val.data.Type):
		fc.xmap()
	case fc.val.data.Type.Id == jntype.Str:
		fc.str()
	}
}
