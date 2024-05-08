// copyright (c) 2024 arfy slowy - derunelabs
//
// permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "software"), to deal
// in the software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the software, and to permit persons to whom the software is
// furnished to do so, subject to the following conditions:
//
// the above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the software.
//
// the software is provided "as is", without warranty of any kind, express or
// implied, including but not limited to the warranties of merchantability,
// fitness for a particular purpose and noninfringement. in no event shall the
// authors or copyright holders be liable for any claim, damages or other
// liability, whether in an action of contract, tort or otherwise, arising from,
// out of or in connection with the software or the use or other dealings in the
// software.

package types

import (
	"github.com/DeRuneLabs/jane/ast"
	"github.com/DeRuneLabs/jane/build"
	"github.com/DeRuneLabs/jane/lexer"
)

func I8CompatibleWith(t uint8) bool {
	t = GetRealCode(t)
	return t == I8
}

func I16CompatibleWith(t uint8) bool {
	t = GetRealCode(t)
	return t == I8 || t == I16 || t == U8
}

func I32CompatibleWith(t uint8) bool {
	t = GetRealCode(t)
	return t == I8 || t == I16 || t == I32 || t == U8 || t == U16
}

func I64CompatibleWith(t uint8) bool {
	t = GetRealCode(t)
	switch t {
	case I8, I16, I32, I64, U8, U16, U32:
		return true
	default:
		return false
	}
}

func U8CompatibleWith(t uint8) bool {
	t = GetRealCode(t)
	return t == U8
}

func U16CompatibleWith(t uint8) bool {
	t = GetRealCode(t)
	return t == U8 || t == U16
}

func U32CompatibleWith(t uint8) bool {
	t = GetRealCode(t)
	return t == U8 || t == U16 || t == U32
}

func U64CompatibleWith(t uint8) bool {
	t = GetRealCode(t)
	return t == U8 || t == U16 || t == U32 || t == U64
}

func F32CompatibleWith(t uint8) bool {
	t = GetRealCode(t)
	switch t {
	case F32, I8, I16, I32, I64, U8, U16, U32, U64:
		return true
	default:
		return false
	}
}

func F64CompatibleWith(t uint8) bool {
	t = GetRealCode(t)
	switch t {
	case F64, F32, I8, I16, I32, I64, U8, U16, U32, U64:
		return true
	default:
		return false
	}
}

func TypesAreCompatible(t1, t2 uint8, ignoreany bool) bool {
	t1 = GetRealCode(t1)
	switch t1 {
	case ANY:
		return !ignoreany
	case I8:
		return I8CompatibleWith(t2)
	case I16:
		return I16CompatibleWith(t2)
	case I32:
		return I32CompatibleWith(t2)
	case I64:
		return I64CompatibleWith(t2)
	case U8:
		return U8CompatibleWith(t2)
	case U16:
		return U16CompatibleWith(t2)
	case U32:
		return U32CompatibleWith(t2)
	case U64:
		return U64CompatibleWith(t2)
	case BOOL:
		return t2 == BOOL
	case STR:
		return t2 == STR
	case F32:
		return F32CompatibleWith(t2)
	case F64:
		return F64CompatibleWith(t2)
	case NIL:
		return t2 == NIL
	}
	return false
}

type Checker struct {
	ErrTok      lexer.Token
	L           Type
	R           Type
	ErrorLogged bool
	IgnoreAny   bool
	AllowAssign bool
	Errors      []build.Log
}

func (c *Checker) push_err_tok(tok lexer.Token, key string, args ...any) {
	c.Errors = append(c.Errors, build.Log{
		Type:   build.ERR,
		Row:    tok.Row,
		Column: tok.Column,
		Path:   tok.File.Path(),
		Text:   build.Errorf(key, args...),
	})
}

func (c *Checker) check_ref() bool {
	if c.L.Kind == c.R.Kind {
		return true
	} else if !c.AllowAssign {
		return false
	}
	c.L = Elem(c.L)
	return c.Check()
}

func (c *Checker) check_ptr() bool {
	if c.R.Id == NIL {
		return true
	} else if IsUnsafePtr(c.L) {
		return true
	}
	return c.L.Kind == c.R.Kind
}

func trait_has_reference_receiver(t *ast.Trait) bool {
	for _, f := range t.Defines.Fns {
		if IsRef(f.Receiver.DataType) {
			return true
		}
	}
	return false
}

func (c *Checker) check_trait() bool {
	if c.R.Id == NIL {
		return true
	}
	t := c.L.Tag.(*ast.Trait)
	lm := c.L.Modifiers()
	ref := false
	switch {
	case IsRef(c.R):
		ref = true
		c.R = Elem(c.R)
		if !IsStruct(c.R) {
			break
		}
		fallthrough
	case IsStruct(c.R):
		if lm != "" {
			return false
		}
		rm := c.R.Modifiers()
		if rm != "" {
			return false
		}
		s := c.R.Tag.(*ast.Struct)
		if !s.HasTrait(t) {
			return false
		}
		if trait_has_reference_receiver(t) && !ref {
			c.ErrorLogged = true
			c.push_err_tok(c.ErrTok, "trait_has_reference_parametered_function")
			return false
		}
		return true
	case IsTrait(c.R):
		return t == c.R.Tag.(*ast.Trait) && lm == c.R.Modifiers()
	}
	return false
}

func (c *Checker) check_struct() bool {
	if c.R.Tag == nil {
		return false
	}
	s1, s2 := c.L.Tag.(*ast.Struct), c.R.Tag.(*ast.Struct)
	switch {
	case s1.Id != s2.Id,
		s1.Token.File != s2.Token.File:
		return false
	}
	if len(s1.Generics) == 0 {
		return true
	}
	n1, n2 := len(s1.GetGenerics()), len(s2.GetGenerics())
	if n1 != n2 {
		return false
	}
	for i, g1 := range s1.GetGenerics() {
		g2 := s2.GetGenerics()[i]
		if !Equals(g1, g2) {
			return false
		}
	}
	return true
}

func (c *Checker) check_slice() bool {
	if c.R.Id == NIL {
		return true
	}
	return c.L.Kind == c.R.Kind
}

func (c *Checker) check_array() bool {
	if !IsArray(c.R) {
		return false
	}
	return c.L.Size.N == c.R.Size.N
}

func (c *Checker) check_map() bool {
	if c.R.Id == NIL {
		return true
	}
	return c.L.Kind == c.R.Kind
}

func (c *Checker) Check() bool {
	switch {
	case IsTrait(c.L), IsTrait(c.R):
		if IsTrait(c.R) {
			c.L, c.R = c.R, c.L
		}
		return c.check_trait()
	case IsRef(c.L), IsRef(c.R):
		if IsRef(c.R) {
			c.L, c.R = c.R, c.L
		}
		return c.check_ref()
	case IsPtr(c.L), IsPtr(c.R):
		if !IsPtr(c.L) && IsPtr(c.R) {
			c.L, c.R = c.R, c.L
		}
		return c.check_ptr()
	case IsSlice(c.L), IsSlice(c.R):
		if IsSlice(c.R) {
			c.L, c.R = c.R, c.L
		}
		return c.check_slice()
	case IsArray(c.L), IsArray(c.R):
		if IsArray(c.R) {
			c.L, c.R = c.R, c.L
		}
		return c.check_array()
	case IsMap(c.L), IsMap(c.R):
		if IsMap(c.R) {
			c.L, c.R = c.R, c.L
		}
		return c.check_map()
	case IsNilCompatible(c.L):
		return c.R.Id == NIL
	case IsNilCompatible(c.R):
		return c.L.Id == NIL
	case IsEnum(c.L), IsEnum(c.R):
		return c.L.Id == c.R.Id && c.L.Kind == c.R.Kind
	case IsStruct(c.L), IsStruct(c.R):
		if c.R.Id == STRUCT {
			c.L, c.R = c.R, c.L
		}
		return c.check_struct()
	}
	return TypesAreCompatible(c.L.Id, c.R.Id, c.IgnoreAny)
}
