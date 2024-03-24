package models

import (
	"strings"

	"github.com/DeRuneLabs/jane/lexer/tokens"
	"github.com/DeRuneLabs/jane/package/jnapi"
)

type IterForeach struct {
	KeyA     Var
	KeyB     Var
	InTok    Tok
	Expr     Expr
	ExprType DataType
}

func (f IterForeach) String(iter Iter) string {
	if !jnapi.IsIgnoreId(f.KeyA.Id) {
		return f.ForeachString(iter)
	}
	return f.IterationString(iter)
}

func (f *IterForeach) ClassicString(iter Iter) string {
	var cxx strings.Builder
	cxx.WriteString("foreach<")
	cxx.WriteString(f.ExprType.String())
	cxx.WriteByte(',')
	cxx.WriteString(f.KeyA.Type.String())
	if !jnapi.IsIgnoreId(f.KeyB.Id) {
		cxx.WriteByte(',')
		cxx.WriteString(f.KeyB.Type.String())
	}
	cxx.WriteString(">(")
	cxx.WriteString(f.Expr.String())
	cxx.WriteString(", [&](")
	cxx.WriteString(f.KeyA.Type.String())
	cxx.WriteByte(' ')
	cxx.WriteString(jnapi.OutId(f.KeyA.Id, f.KeyA.IdTok.File))
	if !jnapi.IsIgnoreId(f.KeyB.Id) {
		cxx.WriteByte(',')
		cxx.WriteString(f.KeyB.Type.String())
		cxx.WriteByte(' ')
		cxx.WriteString(jnapi.OutId(f.KeyB.Id, f.KeyB.IdTok.File))
	}
	cxx.WriteString(") -> void ")
	cxx.WriteString(iter.Block.String())
	cxx.WriteString(");")
	return cxx.String()
}

func (f *IterForeach) MapString(iter Iter) string {
	var cxx strings.Builder
	cxx.WriteString("foreach<")
	types := f.ExprType.Tag.([]DataType)
	cxx.WriteString(types[0].String())
	cxx.WriteByte(',')
	cxx.WriteString(types[1].String())
	cxx.WriteString(">(")
	cxx.WriteString(f.Expr.String())
	cxx.WriteString(", [&](")
	cxx.WriteString(f.KeyA.Type.String())
	cxx.WriteByte(' ')
	cxx.WriteString(jnapi.OutId(f.KeyA.Id, f.KeyA.IdTok.File))
	if !jnapi.IsIgnoreId(f.KeyB.Id) {
		cxx.WriteByte(',')
		cxx.WriteString(f.KeyB.Type.String())
		cxx.WriteByte(' ')
		cxx.WriteString(jnapi.OutId(f.KeyB.Id, f.KeyB.IdTok.File))
	}
	cxx.WriteString(") -> void ")
	cxx.WriteString(iter.Block.String())
	cxx.WriteString(");")
	return cxx.String()
}

func (f *IterForeach) ForeachString(iter Iter) string {
	switch {
	case f.ExprType.Kind == tokens.STR,
		strings.HasPrefix(f.ExprType.Kind, "[]"):
		return f.ClassicString(iter)
	case f.ExprType.Kind[0] == '[':
		return f.MapString(iter)
	}
	return ""
}

func (f IterForeach) IterationString(iter Iter) string {
	var cxx strings.Builder
	cxx.WriteString("for (auto ")
	cxx.WriteString(jnapi.OutId(f.KeyB.Id, f.KeyB.IdTok.File))
	cxx.WriteString(" : ")
	cxx.WriteString(f.Expr.String())
	cxx.WriteString(") ")
	cxx.WriteString(iter.Block.String())
	return cxx.String()
}
