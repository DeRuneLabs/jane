package parser

import (
	"github.com/De-Rune/jane/ast"
	"github.com/De-Rune/jane/package/jn"
)

func typeIsVoidReturn(t ast.DataTypeAST) bool {
	return t.Code == jn.Void && !t.MultiTyped
}

func typeOfArrayElements(t ast.DataTypeAST) ast.DataTypeAST {
	t.Value = t.Value[2:]
	return t
}

func typeIsPointer(t ast.DataTypeAST) bool {
	if t.Value == "" {
		return false
	}
	return t.Value[0] == '*'
}

func typeIsArray(t ast.DataTypeAST) bool {
	if t.Value == "" {
		return false
	}
	return t.Value[0] == '['
}

func typeIsSingle(dt ast.DataTypeAST) bool {
	return !typeIsPointer(dt) && !typeIsArray(dt) && dt.Code != jn.Function
}

func checkValidityConstantDataType(dt ast.DataTypeAST) bool {
	return typeIsSingle(dt)
}

func typeIsNilCompatible(t ast.DataTypeAST) bool {
	return t.Code == jn.Function || typeIsPointer(t)
}

func checkArrayCompatiblity(arrT, t ast.DataTypeAST) bool {
	if t.Code == jn.Nil {
		return true
	}
	return arrT.Value == t.Value
}

func typeIsLvalue(t ast.DataTypeAST) bool {
	return typeIsPointer(t) || typeIsArray(t)
}

func typesAreCompatible(t1, t2 ast.DataTypeAST, ignoreany bool) bool {
	switch {
	case typeIsArray(t1) || typeIsArray(t2):
		if typeIsArray(t2) {
			t1, t2 = t2, t1
		}
		return checkArrayCompatiblity(t1, t2)
	case typeIsNilCompatible(t1) || typeIsNilCompatible(t2):
		return t1.Code == jn.Nil || t2.Code == jn.Nil
	}
	return jn.TypesAreCompatible(t1.Code, t2.Code, ignoreany)
}

func typeIsVariadicable(t ast.DataTypeAST) bool {
	return typeIsArray(t)
}
