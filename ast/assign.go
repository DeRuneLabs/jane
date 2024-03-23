package ast

import "github.com/DeRuneLabs/jane/lexer/tokens"

type AssignInfo struct {
	Left   Toks
	Right  Toks
	Setter Tok
	Ok     bool
	IsExpr bool
}

func IsAssign(id uint8) bool {
	return id == tokens.Id ||
		id == tokens.Brace ||
		id == tokens.Operator
}

func IsAssignOperator(kind string) bool {
	return kind == tokens.EQUAL ||
		kind == tokens.PLUS_EQUAL ||
		kind == tokens.MINUS_EQUAL ||
		kind == tokens.SLASH_EQUAL ||
		kind == tokens.STAR_EQUAL ||
		kind == tokens.PERCENT_EQUAL ||
		kind == tokens.RSHIFT_EQUAL ||
		kind == tokens.LSHIFT_EQUAL ||
		kind == tokens.VLINE_EQUAL ||
		kind == tokens.AMPER_EQUAL ||
		kind == tokens.CARET_EQUAL
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
		}
		if tok.Id == tokens.Operator &&
			IsAssignOperator(tok.Kind) {
			return true
		}
	}
	return false
}
