package ast

import (
	"github.com/DeRuneLabs/jane/ast/models"
	"github.com/DeRuneLabs/jane/lexer/tokens"
)

type blockStatement struct {
	pos            int
	block          *models.Block
	srcToks        *Toks
	blockToks      *Toks
	toks           Toks
	nextToks       Toks
	withTerminator bool
}

func IsStatement(current, prev Tok) (ok bool, withTerminator bool) {
	ok = current.Id == tokens.SemiColon || prev.Row < current.Row
	withTerminator = current.Id == tokens.SemiColon
	return
}

func NextStatementPos(toks Toks, start int) (int, bool) {
	braceCount := 0
	i := start
	for ; i < len(toks); i++ {
		var isStatement, withTerminator bool
		tok := toks[i]
		if tok.Id == tokens.Brace {
			switch tok.Kind {
			case tokens.LBRACE, tokens.LBRACKET, tokens.LPARENTHESES:
				if braceCount == 0 && i > start {
					isStatement, withTerminator = IsStatement(tok, toks[i-1])
					if isStatement {
						goto ret
					}
				}
				braceCount++
				continue
			default:
				braceCount--
				if braceCount == 0 && i+1 < len(toks) {
					isStatement, withTerminator = IsStatement(toks[i+1], tok)
					if isStatement {
						i++
						goto ret
					}
				}
				continue
			}
		}
		if braceCount != 0 {
			continue
		} else if i > start {
			isStatement, withTerminator = IsStatement(tok, toks[i-1])
		} else {
			isStatement, withTerminator = IsStatement(tok, tok)
		}
		if !isStatement {
			continue
		}
	ret:
		if withTerminator {
			i++
		}
		return i, withTerminator
	}
	return i, false
}
