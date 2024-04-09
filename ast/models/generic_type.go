package models

import (
	"strings"

	"github.com/DeRuneLabs/jane/package/jnapi"
)

type GenericType struct {
	Tok Tok
	Id  string
}

func (gt GenericType) String() string {
	var cpp strings.Builder
	cpp.WriteString("typename ")
	cpp.WriteString(jnapi.AsId(gt.Id))
	return cpp.String()
}
