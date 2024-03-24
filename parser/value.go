package parser

import "github.com/DeRuneLabs/jane/ast/models"

type value struct {
	data     models.Data
	constant bool
	volatile bool
	lvalue   bool
	variadic bool
	isType   bool
}
