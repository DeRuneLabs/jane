package models

import (
	"strings"

	"github.com/DeRuneLabs/jane/lexer/tokens"
	"github.com/DeRuneLabs/jane/package/jn"
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
	var cpp strings.Builder
	cpp.WriteString("foreach<")
	cpp.WriteString(f.ExprType.String())
	cpp.WriteByte(',')
	cpp.WriteString(f.KeyA.Type.String())
	if !jnapi.IsIgnoreId(f.KeyB.Id) {
		cpp.WriteByte(',')
		cpp.WriteString(f.KeyB.Type.String())
	}
	cpp.WriteString(">(")
	cpp.WriteString(f.Expr.String())
	cpp.WriteString(", [&](")
	cpp.WriteString(f.KeyA.Type.String())
	cpp.WriteByte(' ')
	cpp.WriteString(jnapi.OutId(f.KeyA.Id, f.KeyA.IdTok.File))
	if !jnapi.IsIgnoreId(f.KeyB.Id) {
		cpp.WriteByte(',')
		cpp.WriteString(f.KeyB.Type.String())
		cpp.WriteByte(' ')
		cpp.WriteString(jnapi.OutId(f.KeyB.Id, f.KeyB.IdTok.File))
	}
	cpp.WriteString(") -> void ")
	cpp.WriteString(iter.Block.String())
	cpp.WriteString(");")
	return cpp.String()
}

func (f *IterForeach) MapString(iter Iter) string {
	var cpp strings.Builder
	cpp.WriteString("foreach<")
	types := f.ExprType.Tag.([]DataType)
	cpp.WriteString(types[0].String())
	cpp.WriteByte(',')
	cpp.WriteString(types[1].String())
	cpp.WriteString(">(")
	cpp.WriteString(f.Expr.String())
	cpp.WriteString(", [&](")
	cpp.WriteString(f.KeyA.Type.String())
	cpp.WriteByte(' ')
	cpp.WriteString(jnapi.OutId(f.KeyA.Id, f.KeyA.IdTok.File))
	if !jnapi.IsIgnoreId(f.KeyB.Id) {
		cpp.WriteByte(',')
		cpp.WriteString(f.KeyB.Type.String())
		cpp.WriteByte(' ')
		cpp.WriteString(jnapi.OutId(f.KeyB.Id, f.KeyB.IdTok.File))
	}
	cpp.WriteString(") -> void ")
	cpp.WriteString(iter.Block.String())
	cpp.WriteString(");")
	return cpp.String()
}

func (f *IterForeach) ForeachString(iter Iter) string {
	switch {
	case f.ExprType.Kind == tokens.STR,
		strings.HasPrefix(f.ExprType.Kind, jn.Prefix_Slice),
		strings.HasPrefix(f.ExprType.Kind, jn.Prefix_Array):
		return f.ClassicString(iter)
	case f.ExprType.Kind[0] == '[':
		return f.MapString(iter)
	}
	return ""
}

func (f IterForeach) IterationString(iter Iter) string {
	var cpp strings.Builder
	cpp.WriteString("for (auto ")
	cpp.WriteString(jnapi.OutId(f.KeyB.Id, f.KeyB.IdTok.File))
	cpp.WriteString(" : ")
	cpp.WriteString(f.Expr.String())
	cpp.WriteString(") ")
	cpp.WriteString(iter.Block.String())
	return cpp.String()
}
