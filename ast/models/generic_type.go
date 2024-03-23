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
	var cxx strings.Builder
	cxx.WriteString("typename ")
	cxx.WriteString(jnapi.OutId(gt.Id, gt.Tok.File))
	return cxx.String()
}
