package parser

import (
	"github.com/DeRuneLabs/jane/ast/models"
	"github.com/DeRuneLabs/jane/lexer/tokens"
	"github.com/DeRuneLabs/jane/package/jnbits"
	"github.com/DeRuneLabs/jane/package/jntype"
)

type solver struct {
	p        *Parser
	left     Toks
	leftVal  models.Data
	right    Toks
	rightVal models.Data
	operator Tok
}

func (s *solver) ptr() (v models.Data) {
	v.Tok = s.operator
	if !typesAreCompatible(s.leftVal.Type, s.rightVal.Type, true) {
		s.p.pusherrtok(s.operator, "incompatible_datatype",
			s.rightVal.Type.Kind, s.leftVal.Type.Kind)
		return
	}
	if !typeIsPtr(s.leftVal.Type) {
		s.leftVal, s.rightVal = s.rightVal, s.leftVal
	}
	switch s.operator.Kind {
	case tokens.PLUS, tokens.MINUS:
		v.Type = s.leftVal.Type
	case tokens.EQUALS, tokens.NOT_EQUALS:
		v.Type.Id = jntype.Bool
		v.Type.Kind = tokens.BOOL
	default:
		s.p.pusherrtok(s.operator, "operator_notfor_jntype", s.operator.Kind, "pointer")
	}
	return
}

func (s *solver) enum() (v models.Data) {
	if s.leftVal.Type.Id == jntype.Enum {
		s.leftVal.Type = s.leftVal.Type.Tag.(*Enum).Type
	}
	if s.rightVal.Type.Id == jntype.Enum {
		s.rightVal.Type = s.rightVal.Type.Tag.(*Enum).Type
	}
	return s.solve()
}

func (s *solver) str() (v models.Data) {
	v.Tok = s.operator
	if s.leftVal.Type.Id != s.rightVal.Type.Id {
		s.p.pusherrtok(s.operator, "incompatible_datatype",
			s.leftVal.Type.Kind, s.rightVal.Type.Kind)
		return
	}
	switch s.operator.Kind {
	case tokens.PLUS:
		v.Type.Id = jntype.Str
		v.Type.Kind = tokens.STR
	case tokens.EQUALS, tokens.NOT_EQUALS:
		v.Type.Id = jntype.Bool
		v.Type.Kind = tokens.BOOL
	default:
		s.p.pusherrtok(s.operator, "operator_notfor_jntype",
			s.operator.Kind, tokens.STR)
	}
	return
}

func (s *solver) any() (v models.Data) {
	v.Tok = s.operator
	switch s.operator.Kind {
	case tokens.EQUALS, tokens.NOT_EQUALS:
		v.Type.Id = jntype.Bool
		v.Type.Kind = tokens.BOOL
	default:
		s.p.pusherrtok(s.operator, "operator_notfor_jntype", s.operator.Kind, tokens.ANY)
	}
	return
}

func (s *solver) bool() (v models.Data) {
	v.Tok = s.operator
	if !typesAreCompatible(s.leftVal.Type, s.rightVal.Type, true) {
		s.p.pusherrtok(s.operator, "incompatible_datatype",
			s.rightVal.Type.Kind, s.leftVal.Type.Kind)
		return
	}
	switch s.operator.Kind {
	case tokens.EQUALS, tokens.NOT_EQUALS:
		v.Type.Id = jntype.Bool
		v.Type.Kind = tokens.BOOL
	default:
		s.p.pusherrtok(s.operator, "operator_notfor_jntype",
			s.operator.Kind, tokens.BOOL)
	}
	return
}

func (s *solver) float() (v models.Data) {
	v.Tok = s.operator
	if !jntype.IsNumericType(s.leftVal.Type.Id) ||
		!jntype.IsNumericType(s.rightVal.Type.Id) {
		s.p.pusherrtok(s.operator, "incompatible_datatype",
			s.rightVal.Type.Kind, s.leftVal.Type.Kind)
		return
	}
	switch s.operator.Kind {
	case tokens.EQUALS, tokens.NOT_EQUALS, tokens.LESS, tokens.GREAT,
		tokens.GREAT_EQUAL, tokens.LESS_EQUAL:
		v.Type.Id = jntype.Bool
		v.Type.Kind = tokens.BOOL
	case tokens.PLUS, tokens.MINUS, tokens.STAR, tokens.SLASH:
		v.Type = s.leftVal.Type
		if jntype.TypeGreaterThan(s.rightVal.Type.Id, v.Type.Id) {
			v.Type = s.rightVal.Type
		}
	default:
		s.p.pusherrtok(s.operator, "operator_notfor_float", s.operator.Kind)
	}
	return
}

func (s *solver) signed() (v models.Data) {
	v.Tok = s.operator
	if !jntype.IsNumericType(s.leftVal.Type.Id) ||
		!jntype.IsNumericType(s.rightVal.Type.Id) {
		s.p.pusherrtok(s.operator, "incompatible_datatype",
			s.rightVal.Type.Kind, s.leftVal.Type.Kind)
		return
	}
	switch s.operator.Kind {
	case tokens.EQUALS, tokens.NOT_EQUALS, tokens.LESS,
		tokens.GREAT, tokens.GREAT_EQUAL, tokens.LESS_EQUAL:
		v.Type.Id = jntype.Bool
		v.Type.Kind = tokens.BOOL
	case tokens.PLUS, tokens.MINUS, tokens.STAR, tokens.SLASH,
		tokens.PERCENT, tokens.AMPER, tokens.VLINE, tokens.CARET:
		v.Type = s.leftVal.Type
		if jntype.TypeGreaterThan(s.rightVal.Type.Id, v.Type.Id) {
			v.Type = s.rightVal.Type
		}
	case tokens.RSHIFT, tokens.LSHIFT:
		v.Type = s.leftVal.Type
		if !jntype.IsUnsignedNumericType(s.rightVal.Type.Id) &&
			!checkIntBit(s.rightVal, jnbits.BitsizeType(jntype.U64)) {
			s.p.pusherrtok(s.rightVal.Tok, "bitshift_must_unsigned")
		}
	default:
		s.p.pusherrtok(s.operator, "operator_notfor_int", s.operator.Kind)
	}
	return
}

func (s *solver) unsigned() (v models.Data) {
	v.Tok = s.operator
	if !jntype.IsNumericType(s.leftVal.Type.Id) ||
		!jntype.IsNumericType(s.rightVal.Type.Id) {
		s.p.pusherrtok(s.operator, "incompatible_datatype",
			s.rightVal.Type.Kind, s.leftVal.Type.Kind)
		return
	}
	switch s.operator.Kind {
	case tokens.EQUALS, tokens.NOT_EQUALS, tokens.LESS,
		tokens.GREAT, tokens.GREAT_EQUAL, tokens.LESS_EQUAL:
		v.Type.Id = jntype.Bool
		v.Type.Kind = tokens.BOOL
	case tokens.PLUS, tokens.MINUS, tokens.STAR, tokens.SLASH,
		tokens.PERCENT, tokens.AMPER, tokens.VLINE, tokens.CARET:
		v.Type = s.leftVal.Type
		if jntype.TypeGreaterThan(s.rightVal.Type.Id, v.Type.Id) {
			v.Type = s.rightVal.Type
		}
	case tokens.RSHIFT, tokens.LSHIFT:
		v.Type = s.leftVal.Type
		if jntype.TypeGreaterThan(s.rightVal.Type.Id, v.Type.Id) {
			v.Type = s.rightVal.Type
		}
	default:
		s.p.pusherrtok(s.operator, "operator_notfor_uint", s.operator.Kind)
	}
	return
}

func (s *solver) logical() (v models.Data) {
	v.Tok = s.operator
	v.Type.Id = jntype.Bool
	v.Type.Kind = tokens.BOOL
	if s.leftVal.Type.Id != jntype.Bool || s.rightVal.Type.Id != jntype.Bool {
		s.p.pusherrtok(s.operator, "logical_not_bool")
	}
	return
}

func (s *solver) array() (v models.Data) {
	v.Tok = s.operator
	if !typesAreCompatible(s.leftVal.Type, s.rightVal.Type, true) {
		s.p.pusherrtok(s.operator, "incompatible_datatype",
			s.rightVal.Type.Kind, s.leftVal.Type.Kind)
		return
	}
	switch s.operator.Kind {
	case tokens.EQUALS, tokens.NOT_EQUALS:
		v.Type.Id = jntype.Bool
		v.Type.Kind = tokens.BOOL
	default:
		s.p.pusherrtok(s.operator, "operator_notfor_jntype", s.operator.Kind, "array")
	}
	return
}

func (s *solver) nil() (v models.Data) {
	v.Tok = s.operator
	if !typesAreCompatible(s.leftVal.Type, s.rightVal.Type, false) {
		s.p.pusherrtok(s.operator, "incompatible_datatype",
			s.rightVal.Type.Kind, s.leftVal.Type.Kind)
		return
	}
	switch s.operator.Kind {
	case tokens.NOT_EQUALS, tokens.EQUALS:
		v.Type.Id = jntype.Bool
		v.Type.Kind = tokens.BOOL
	default:
		s.p.pusherrtok(s.operator, "operator_notfor_jntype",
			s.operator.Kind, tokens.NIL)
	}
	return
}

func (s *solver) check() bool {
	switch s.operator.Kind {
	case tokens.PLUS, tokens.MINUS, tokens.STAR, tokens.SLASH, tokens.PERCENT, tokens.RSHIFT,
		tokens.LSHIFT, tokens.AMPER, tokens.VLINE, tokens.CARET, tokens.EQUALS, tokens.NOT_EQUALS,
		tokens.GREAT, tokens.LESS, tokens.GREAT_EQUAL, tokens.LESS_EQUAL:
	case tokens.AND, tokens.OR:
	default:
		s.p.pusherrtok(s.operator, "invalid_operator")
		return false
	}
	return true
}

func (s *solver) solve() (v models.Data) {
	defer func() {
		if v.Type.Id == jntype.Void {
			v.Type.Kind = jntype.VoidTypeStr
		}
	}()
	if !s.check() {
		return
	}
	switch s.operator.Kind {
	case tokens.AND, tokens.OR:
		return s.logical()
	}
	switch {
	case typeIsArray(s.leftVal.Type), typeIsArray(s.rightVal.Type):
		return s.array()
	case typeIsPtr(s.leftVal.Type), typeIsPtr(s.rightVal.Type):
		return s.ptr()
	case s.leftVal.Type.Id == jntype.Enum, s.rightVal.Type.Id == jntype.Enum:
		return s.enum()
	case s.leftVal.Type.Id == jntype.Nil, s.rightVal.Type.Id == jntype.Nil:
		return s.nil()
	case s.leftVal.Type.Id == jntype.Any, s.rightVal.Type.Id == jntype.Any:
		return s.any()
	case s.leftVal.Type.Id == jntype.Bool, s.rightVal.Type.Id == jntype.Bool:
		return s.bool()
	case s.leftVal.Type.Id == jntype.Str, s.rightVal.Type.Id == jntype.Str:
		return s.str()
	case jntype.IsFloatType(s.leftVal.Type.Id),
		jntype.IsFloatType(s.rightVal.Type.Id):
		return s.float()
	case jntype.IsUnsignedNumericType(s.leftVal.Type.Id),
		jntype.IsUnsignedNumericType(s.rightVal.Type.Id):
		return s.unsigned()
	case jntype.IsSignedNumericType(s.leftVal.Type.Id),
		jntype.IsSignedNumericType(s.rightVal.Type.Id):
		return s.signed()
	}
	return
}
