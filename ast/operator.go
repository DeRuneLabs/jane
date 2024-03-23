package ast

import "github.com/DeRuneLabs/jane/lexer/tokens"

func IsSingleOperator(kind string) bool {
	return kind == tokens.MINUS ||
		kind == tokens.PLUS ||
		kind == tokens.TILDE ||
		kind == tokens.EXCLAMATION ||
		kind == tokens.STAR ||
		kind == tokens.AMPER
}

func IsSolidOperator(kind string) bool {
	return kind == tokens.PLUS ||
		kind == tokens.MINUS ||
		kind == tokens.STAR ||
		kind == tokens.SLASH ||
		kind == tokens.PERCENT ||
		kind == tokens.AMPER ||
		kind == tokens.VLINE ||
		kind == tokens.CARET ||
		kind == tokens.LESS ||
		kind == tokens.GREAT ||
		kind == tokens.TILDE ||
		kind == tokens.EXCLAMATION
}

func IsExprOperator(kind string) bool { return kind == tokens.TRIPLE_DOT }
