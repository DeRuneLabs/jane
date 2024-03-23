package models

import "github.com/DeRuneLabs/jane/package/jnapi"

type Defer struct {
	Tok  Tok
	Expr Expr
}

func (d Defer) String() string {
	return jnapi.ToDeferredCall(d.Expr.String())
}

type ConcurrentCall struct {
	Tok  Tok
	Expr Expr
}

func (cc ConcurrentCall) String() string {
	return jnapi.ToConcurrentCall(cc.Expr.String())
}
