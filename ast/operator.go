package ast

import "github.com/DeRuneLabs/jane/lexer/tokens"

var UnaryOperators = [...]string{
	0: tokens.MINUS,
	1: tokens.PLUS,
	2: tokens.TILDE,
	3: tokens.EXCLAMATION,
	4: tokens.STAR,
	5: tokens.AMPER,
}

var SolidOperators = [...]string{
	0:  tokens.PLUS,
	1:  tokens.MINUS,
	2:  tokens.STAR,
	3:  tokens.SOLIDUS,
	4:  tokens.PERCENT,
	5:  tokens.AMPER,
	6:  tokens.VLINE,
	7:  tokens.CARET,
	8:  tokens.LESS,
	9:  tokens.GREAT,
	10: tokens.EXCLAMATION,
}

var ExpressionOperators = [...]string{
	0: tokens.TRIPLE_DOT,
	1: tokens.COLON,
}

func IsUnaryOperator(kind string) bool {
	return existOperator(kind, UnaryOperators[:])
}

func IsSolidOperator(kind string) bool {
	return existOperator(kind, SolidOperators[:])
}

func IsExprOperator(kind string) bool {
	return existOperator(kind, ExpressionOperators[:])
}

func existOperator(kind string, operators []string) bool {
	for _, operator := range operators {
		if kind == operator {
			return true
		}
	}
	return false
}
