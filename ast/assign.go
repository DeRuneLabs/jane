package ast

import "github.com/DeRuneLabs/jane/lexer/tokens"

type AssignInfo struct {
	Left   Toks
	Right  Toks
	Setter Tok
	Ok     bool
	IsExpr bool
}

var SuffixOperators = [...]string{
	0: tokens.DOUBLE_PLUS,
	1: tokens.DOUBLE_MINUS,
}

var AssignOperators = [...]string{
	0:  tokens.EQUAL,
	1:  tokens.PLUS_EQUAL,
	2:  tokens.MINUS_EQUAL,
	3:  tokens.SLASH_EQUAL,
	4:  tokens.STAR_EQUAL,
	5:  tokens.PERCENT_EQUAL,
	6:  tokens.RSHIFT_EQUAL,
	7:  tokens.LSHIFT_EQUAL,
	8:  tokens.VLINE_EQUAL,
	9:  tokens.AMPER_EQUAL,
	10: tokens.CARET_EQUAL,
}

func IsAssign(id uint8) bool {
	switch id {
	case tokens.Id, tokens.Brace, tokens.Operator:
		return true
	}
	return false
}

func IsSuffixOperator(kind string) bool {
	for _, operator := range SuffixOperators {
		if kind == operator {
			return true
		}
	}
	return false
}

func IsAssignOperator(kind string) bool {
	if IsSuffixOperator(kind) {
		return true
	}
	for _, operator := range AssignOperators {
		if kind == operator {
			return true
		}
	}
	return false
}

func CheckAssignToks(toks Toks) bool {
	if len(toks) == 0 || !IsAssign(toks[0].Id) {
		return false
	}
	braceCount := 0
	for _, tok := range toks {
		if tok.Id == tokens.Brace {
			switch tok.Kind {
			case tokens.LBRACE, tokens.LBRACKET, tokens.LPARENTHESES:
				braceCount++
			default:
				braceCount--
			}
		}
		if braceCount < 0 {
			return false
		} else if braceCount > 0 {
			continue
		} else if tok.Id == tokens.Operator && IsAssignOperator(tok.Kind) {
			return true
		}
	}
	return false
}
