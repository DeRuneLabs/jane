package parser

import (
	"github.com/DeRuneLabs/jane/lexer/tokens"
	"github.com/DeRuneLabs/jane/package/jnapi"
	"github.com/DeRuneLabs/jane/package/jntype"
)

type unary struct {
	tok    Tok
	toks   Toks
	model  *exprModel
	parser *Parser
}

func (u *unary) minus() value {
	v := u.parser.evalExprPart(u.toks, u.model)
	if !typeIsPure(v.data.Type) || !jntype.IsNumericType(v.data.Type.Id) {
		u.parser.pusherrtok(u.tok, "invalid_type_unary_operator", '-')
	}
	if isConstNumeric(v.data.Value) {
		v.data.Value = tokens.MINUS + v.data.Value
	}
	return v
}

func (u *unary) plus() value {
	v := u.parser.evalExprPart(u.toks, u.model)
	if !typeIsPure(v.data.Type) || !jntype.IsNumericType(v.data.Type.Id) {
		u.parser.pusherrtok(u.tok, "invalid_type_unary_operator", '+')
	}
	return v
}

func (u *unary) tilde() value {
	v := u.parser.evalExprPart(u.toks, u.model)
	if !typeIsPure(v.data.Type) || !jntype.IsIntegerType(v.data.Type.Id) {
		u.parser.pusherrtok(u.tok, "invalid_type_unary_operator", '~')
	}
	return v
}

func (u *unary) logicalNot() value {
	v := u.parser.evalExprPart(u.toks, u.model)
	if !isBoolExpr(v) {
		u.parser.pusherrtok(u.tok, "invalid_type_unary_operator", '!')
	}
	v.data.Type.Id = jntype.Bool
	v.data.Type.Kind = tokens.BOOL
	return v
}

func (u *unary) star() value {
	v := u.parser.evalExprPart(u.toks, u.model)
	v.lvalue = true
	if !typeIsExplicitPtr(v.data.Type) {
		u.parser.pusherrtok(u.tok, "invalid_type_unary_operator", '*')
	} else {
		v.data.Type.Kind = v.data.Type.Kind[1:]
	}
	return v
}

func (u *unary) amper() value {
	v := u.parser.evalExprPart(u.toks, u.model)
	switch {
	case typeIsFunc(v.data.Type):
		mainNode := &u.model.nodes[u.model.index]
		mainNode.nodes = mainNode.nodes[1:]
		node := &u.model.nodes[u.model.index].nodes[0]
		switch t := (*node).(type) {
		case anonFuncExpr:
			if t.capture == jnapi.LambdaByReference {
				u.parser.pusherrtok(u.tok, "invalid_type_unary_operator", tokens.AMPER)
				break
			}
			t.capture = jnapi.LambdaByReference
			*node = t
		default:
			u.parser.pusherrtok(u.tok, "invalid_type_unary_operator", tokens.AMPER)
		}
	default:
		if !canGetPtr(v) {
			u.parser.pusherrtok(u.tok, "invalid_type_unary_operator", tokens.AMPER)
		}
		v.lvalue = true
		v.data.Type.Kind = tokens.STAR + v.data.Type.Kind
	}
	return v
}
