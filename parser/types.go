package parser

import (
	"github.com/De-Rune/jane/ast"
	"github.com/De-Rune/jane/package/jn"
)

func typeIsVoidRet(t ast.DataType) bool {
	return t.Id == jn.Void && !t.MultiTyped
}

func typeOfArrayElements(t ast.DataType) ast.DataType {
	t.Val = t.Val[2:]
	return t
}

func typeIsPtr(t ast.DataType) bool {
	if t.Val == "" {
		return false
	}
	return t.Val[0] == '*'
}

func typeIsArray(t ast.DataType) bool {
	if t.Val == "" {
		return false
	}
	return t.Val[0] == '['
}

func typeIsSingle(dt ast.DataType) bool {
	return !typeIsPtr(dt) &&
		!typeIsArray(dt) &&
		dt.Id != jn.Func
}

func typeIsNilCompatible(t ast.DataType) bool {
	return t.Id == jn.Func || typeIsPtr(t)
}

func checkArrayCompatiblity(arrT, t ast.DataType) bool {
	if t.Id == jn.Nil {
		return true
	}
	return arrT.Val == t.Val
}

func typeIsLvalue(t ast.DataType) bool {
	return typeIsPtr(t) || typeIsArray(t)
}

func typesAreCompatible(t1, t2 ast.DataType, ignoreany bool) bool {
	switch {
	case typeIsArray(t1) || typeIsArray(t2):
		if typeIsArray(t2) {
			t1, t2 = t2, t1
		}
		return checkArrayCompatiblity(t1, t2)
	case typeIsNilCompatible(t1) || typeIsNilCompatible(t2):
		return t1.Id == jn.Nil || t2.Id == jn.Nil
	}
	return jn.TypesAreCompatible(t1.Id, t2.Id, ignoreany)
}

func typeIsVariadicable(t ast.DataType) bool {
	return typeIsArray(t)
}

func typeIsMut(t ast.DataType) bool {
	return typeIsPtr(t)
}
