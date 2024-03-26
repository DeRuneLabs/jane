package parser

import (
	"github.com/DeRuneLabs/jane/lexer/tokens"
	"github.com/DeRuneLabs/jane/package/jnapi"
	"github.com/DeRuneLabs/jane/package/jnbits"
	"github.com/DeRuneLabs/jane/package/jntype"
)

type valueEvaluator struct {
	tok   Tok
	model *exprModel
	p     *Parser
}

func (ve *valueEvaluator) str() value {
	var v value
	v.data.Value = ve.tok.Kind
	v.data.Type.Id = jntype.Str
	v.data.Type.Kind = tokens.STR
	content := []byte(ve.tok.Kind[1 : len(ve.tok.Kind)-1])
	if israwstr(ve.tok.Kind) {
		ve.model.appendSubNode(exprNode{jnapi.ToRawStr(content)})
	} else {
		ve.model.appendSubNode(exprNode{jnapi.ToStr(content)})
	}
	return v
}

func toCharLiteral(kind string) (string, bool) {
	kind = kind[1 : len(kind)-1]
	isByte := false
	switch {
	case len(kind) == 1 && kind[0] <= 255:
		isByte = true
	case kind[0] == '\\' && kind[1] == 'x':
		isByte = true
	case kind[0] == '\\' && kind[1] >= '0' && kind[1] <= '7':
		isByte = true
	}
	return kind, isByte
}

func (ve *valueEvaluator) char() value {
	var v value
	v.data.Value = ve.tok.Kind
	content, isByte := toCharLiteral(ve.tok.Kind)
	if isByte {
		v.data.Type.Id = jntype.U8
		v.data.Type.Kind = tokens.U8
		content = jnapi.ToChar(content[0])
	} else {
		v.data.Type.Id = jntype.I32
		v.data.Type.Kind = tokens.I32
		content = jnapi.ToRune([]byte(content))
	}
	ve.model.appendSubNode(exprNode{content})
	return v
}

func (ve *valueEvaluator) bool() value {
	var v value
	v.data.Value = ve.tok.Kind
	v.data.Type.Id = jntype.Bool
	v.data.Type.Kind = tokens.BOOL
	ve.model.appendSubNode(exprNode{ve.tok.Kind})
	return v
}

func (ve *valueEvaluator) nil() value {
	var v value
	v.data.Value = ve.tok.Kind
	v.data.Type.Id = jntype.Nil
	v.data.Type.Kind = jntype.NilTypeStr
	ve.model.appendSubNode(exprNode{ve.tok.Kind})
	return v
}

func (ve *valueEvaluator) float() value {
	var v value
	v.data.Value = ve.tok.Kind
	v.data.Type.Id = jntype.F64
	v.data.Type.Kind = tokens.F64
	return v
}

func (ve *valueEvaluator) integer() value {
	var v value
	v.data.Value = ve.tok.Kind
	intbit := jnbits.BitsizeType(jntype.Int)
	switch {
	case jnbits.CheckBitInt(ve.tok.Kind, intbit):
		v.data.Type.Id = jntype.Int
		v.data.Type.Kind = tokens.INT
	case intbit < jnbits.MaxInt && jnbits.CheckBitInt(ve.tok.Kind, jnbits.MaxInt):
		v.data.Type.Id = jntype.I64
		v.data.Type.Kind = tokens.I64
	default:
		v.data.Type.Id = jntype.U64
		v.data.Type.Kind = tokens.U64
	}
	return v
}

func (ve *valueEvaluator) numeric() value {
	var v value
	if isfloat(ve.tok.Kind) {
		v = ve.float()
	} else {
		v = ve.integer()
	}
	cxxId := jntype.CxxTypeIdFromType(v.data.Type.Id)
	node := exprNode{cxxId + "{" + ve.tok.Kind + "}"}
	ve.model.appendSubNode(node)
	return v
}

func (ve *valueEvaluator) varId(id string, variable *Var) (v value) {
	variable.Used = true
	v.data.Value = id
	v.data.Type = variable.Type
	v.constant = variable.Const
	v.volatile = variable.Volatile
	v.data.Tok = variable.IdTok
	v.lvalue = true
	if variable.IdTok.Id == tokens.NA {
		ve.model.appendSubNode(exprNode{jnapi.OutId(id, nil)})
	} else {
		ve.model.appendSubNode(exprNode{jnapi.OutId(id, variable.IdTok.File)})
	}
	return
}

func (ve *valueEvaluator) funcId(id string, f *function) (v value) {
	f.used = true
	v.data.Value = id
	v.data.Type.Id = jntype.Func
	v.data.Type.Tag = f.Ast
	v.data.Type.Kind = f.Ast.DataTypeString()
	v.data.Tok = f.Ast.Tok
	ve.model.appendSubNode(exprNode{f.outId()})
	return
}

func (ve *valueEvaluator) enumId(id string, e *Enum) (v value) {
	e.Used = true
	v.data.Value = id
	v.data.Type.Id = jntype.Enum
	v.data.Type.Tag = e
	v.data.Type.Kind = e.Id
	v.data.Tok = e.Tok
	v.constant = true
	v.isType = true
	if e.Tok.Id == tokens.NA {
		ve.model.appendSubNode(exprNode{jnapi.OutId(id, nil)})
	} else {
		ve.model.appendSubNode(exprNode{jnapi.OutId(id, e.Tok.File)})
	}
	return
}

func (ve *valueEvaluator) structId(id string, s *jnstruct) (v value) {
	s.Used = true
	v.data.Value = id
	v.data.Type.Id = jntype.Struct
	v.data.Type.Tag = s
	v.data.Type.Kind = s.Ast.Id
	v.data.Type.Tok = s.Ast.Tok
	v.data.Tok = s.Ast.Tok
	v.isType = true
	if s.Ast.Tok.Id == tokens.NA {
		ve.model.appendSubNode(exprNode{jnapi.OutId(id, nil)})
	} else {
		ve.model.appendSubNode(exprNode{jnapi.OutId(id, s.Ast.Tok.File)})
	}
	return
}

func (ve *valueEvaluator) typeId(id string, t *Type) (_ value, _ bool) {
	dt, ok := ve.p.realType(t.Type, true)
	if !ok {
		return
	}
	if typeIsStruct(dt) {
		return ve.structId(id, dt.Tag.(*jnstruct)), true
	}
	return
}

func (ve *valueEvaluator) id() (_ value, ok bool) {
	id := ve.tok.Kind
	if v, _ := ve.p.varById(id); v != nil {
		return ve.varId(id, v), true
	} else if f, _, _ := ve.p.FuncById(id); f != nil {
		return ve.funcId(id, f), true
	} else if e, _, _ := ve.p.enumById(id); e != nil {
		return ve.enumId(id, e), true
	} else if s, _, _ := ve.p.structById(id); s != nil {
		return ve.structId(id, s), true
	} else if t, _, _ := ve.p.typeById(id); t != nil {
		return ve.typeId(id, t)
	} else {
		ve.p.pusherrtok(ve.tok, "id_noexist", id)
	}
	return
}
