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
	return typeIsPure(t)
}

func typeIsSinglePtr(t DataType) bool {
	return t.Id == jntype.Voidptr
}

func typeIsStruct(dt DataType) bool {
	return dt.Id == jntype.Struct
}

func typeIsEnum(dt DataType) bool {
	return dt.Id == jntype.Enum
}

func typeIsGeneric(generics []*GenericType, t DataType) bool {
	if t.Id != jntype.Id {
		return false
	}
	id, _ := t.KindId()
	for _, generic := range generics {
		if id == generic.Id {
			return true
		}
	}
	return false
}

func typeOfArrayComponents(t DataType) DataType {
	t.Kind = t.Kind[2:]
	return t
}

func typeIsExplicitPtr(t DataType) bool {
	if t.Kind == "" {
		return false
	}
	return t.Kind[0] == '*'
}

func typeIsPtr(t DataType) bool {
	return typeIsExplicitPtr(t) || typeIsSinglePtr(t)
}

func typeIsArray(t DataType) bool {
	if t.Kind == "" {
		return false
	}
	return strings.HasPrefix(t.Kind, "[]")
}

func typeIsMap(t DataType) bool {
	if t.Kind == "" || t.Id != jntype.Map {
		return false
	}
	return t.Id == jntype.Map && t.Kind[0] == '[' && !strings.HasPrefix(t.Kind, "[]")
}

func typeIsFunc(t DataType) bool {
	if t.Id != jntype.Func || t.Kind == "" {
		return false
	}
	return t.Kind[0] == '('
}

func typeIsPure(t DataType) bool {
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
	return arrT.Kind == t.Kind
}

func checkMapCompability(mapT, t DataType) bool {
	if t.Id == jntype.Nil {
		return true
	}
	return mapT.Kind == t.Kind
}

func typeIsLvalue(t DataType) bool {
	return typeIsPtr(t) || typeIsArray(t) || typeIsMap(t)
}

func checkPtrCompability(t1, t2 DataType) bool {
	if typeIsPtr(t2) {
		return true
	}
	if typeIsPure(t2) && jntype.IsIntegerType(t2.Id) {
		return true
	}
	return false
}

func typesEquals(t1, t2 DataType) bool {
	return t1.Id == t2.Id && t1.Kind == t2.Kind
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
		return t1.Id == t2.Id && t1.Kind == t2.Kind
	case t1.Id == jntype.Struct, t2.Id == jntype.Struct:
		if t2.Id == jntype.Struct {
			t1, t2 = t2, t1
		}
		return checkStructCompability(t1, t2)
	}
	return jntype.TypesAreCompatible(t1.Id, t2.Id, ignoreany)
}
