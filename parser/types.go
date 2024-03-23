package parser

import (
	"strings"

	"github.com/DeRuneLabs/jane/lexer/tokens"
	"github.com/DeRuneLabs/jane/package/jntype"
)

func typeIsVoid(t DataType) bool {
	return t.Id == jntype.Void && !t.MultiTyped
}

func typeIsVariadicable(t DataType) bool {
	return typeIsArray(t)
}

func typeIsMut(t DataType) bool {
	return typeIsPtr(t)
}

func typeIsAllowForConst(t DataType) bool {
	return typeIsSingle(t)
}

func typeIsSinglePtr(t DataType) bool {
	return t.Id == jntype.Voidptr
}

func typeIsGeneric(generics []*GenericType, t DataType) bool {
	if t.Id != jntype.Id {
		return false
	}
	id, _ := t.GetValId()
	for _, generic := range generics {
		if id == generic.Id {
			return true
		}
	}
	return false
}

func typeOfArrayComponents(t DataType) DataType {
	t.Val = t.Val[2:]
	return t
}

func typeIsExplicitPtr(t DataType) bool {
	if t.Val == "" {
		return false
	}
	return t.Val[0] == '*'
}

func typeIsPtr(t DataType) bool {
	return typeIsExplicitPtr(t) || typeIsSinglePtr(t)
}

func typeIsArray(t DataType) bool {
	if t.Val == "" {
		return false
	}
	return strings.HasPrefix(t.Val, "[]")
}

func typeIsMap(t DataType) bool {
	if t.Val == "" {
		return false
	}
	return t.Id == jntype.Map && t.Val[0] == '[' && !strings.HasPrefix(t.Val, "[]")
}

func typeIsFunc(t DataType) bool {
	if t.Id != jntype.Func || t.Val == "" {
		return false
	}
	return t.Val[0] == '('
}

func typeIsSingle(t DataType) bool {
	return !typeIsPtr(t) &&
		!typeIsArray(t) &&
		!typeIsMap(t) &&
		!typeIsFunc(t)
}

func subIdAccessorOfType(t DataType) string {
	if typeIsPtr(t) {
		return "->"
	}
	return tokens.DOT
}

func typeIsNilCompatible(t DataType) bool {
	return typeIsFunc(t) || typeIsPtr(t) || typeIsArray(t) || typeIsMap(t)
}

func checkArrayCompatiblity(arrT, t DataType) bool {
	if t.Id == jntype.Nil {
		return true
	}
	return arrT.Val == t.Val
}

func checkMapCompability(mapT, t DataType) bool {
	if t.Id == jntype.Nil {
		return true
	}
	return mapT.Val == t.Val
}

func typeIsLvalue(t DataType) bool {
	return typeIsPtr(t) || typeIsArray(t) || typeIsMap(t)
}

func checkPtrCompability(t1, t2 DataType) bool {
	if typeIsPtr(t2) {
		return true
	}
	if typeIsSingle(t2) && jntype.IsIntegerType(t2.Id) {
		return true
	}
	return false
}

func typesEquals(t1, t2 DataType) bool {
	return t1.Id == t2.Id && t1.Val == t2.Val
}

func checkStructCompability(t1, t2 DataType) bool {
	s1, s2 := t1.Tag.(*jnstruct), t2.Tag.(*jnstruct)
	switch {
	case s1.Ast.Id != s2.Ast.Id,
		s1.Ast.Tok.File != s2.Ast.Tok.File:
		return false
	}
	if len(s1.Ast.Generics) == 0 {
		return true
	}
	n1, n2 := len(s1.generics), len(s2.generics)
	if n1 > 0 || n2 > 0 {
		if n1 != n2 {
			return false
		}
		for i, g1 := range s1.generics {
			g2 := s2.generics[i]
			if !typesEquals(g1, g2) {
				return false
			}
		}
	}
	return true
}

func typesAreCompatible(t1, t2 DataType, ignoreany bool) bool {
	switch {
	case typeIsPtr(t1), typeIsPtr(t2):
		if typeIsPtr(t2) {
			t1, t2 = t2, t1
		}
		return checkPtrCompability(t1, t2)
	case typeIsArray(t1), typeIsArray(t2):
		if typeIsArray(t2) {
			t1, t2 = t2, t1
		}
		return checkArrayCompatiblity(t1, t2)
	case typeIsMap(t1), typeIsMap(t2):
		if typeIsMap(t2) {
			t1, t2 = t2, t1
		}
		return checkMapCompability(t1, t2)
	case typeIsNilCompatible(t1), typeIsNilCompatible(t2):
		return t1.Id == jntype.Nil || t2.Id == jntype.Nil
	case t1.Id == jntype.Enum, t2.Id == jntype.Enum:
		return t1.Id == t2.Id && t1.Val == t2.Val
	case t1.Id == jntype.Struct, t2.Id == jntype.Struct:
		if t2.Id == jntype.Struct {
			t1, t2 = t2, t1
		}
		return checkStructCompability(t1, t2)
	}
	return jntype.TypesAreCompatible(t1.Id, t2.Id, ignoreany)
}
