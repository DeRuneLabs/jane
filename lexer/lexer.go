package lexer

import (
	"fmt"
	"github.com/slowy07/jane/package/jane"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

func (lexer *Lexer) error(error string) {
	lexer.Errors = append(lexer.Errors, fmt.Sprintf("%s %d:%d %s", lexer.File.Path, lexer.Line, lexer.Column, jane.Errors[error]))
}

func (lexer *Lexer) Tokenize() []Token {
	var tokens []Token
	lexer.Errors = nil
	for lexer.Position < len(lexer.File.Content) {
		token := lexer.Token()
		if token.Type != NA {
			tokens = append(tokens, token)
		}
	}
	return tokens
}

func isKeyword(lexerline, kw string) bool {
	return regexp.MustCompile("^" + kw + `(\s+|$|[[:punct:]])`).MatchString(lexerline)
}

func (lexer *Lexer) lexName(lexerline string) string {
	if lexerline[0] != '_' {
		r, _ := utf8.DecodeRuneInString(lexerline)
		if !unicode.IsLetter(r) {
			return ""
		}
	}
	var sb strings.Builder
	for _, run := range lexerline {
		if run != '_' && ('0' > run || '9' < run) && !unicode.IsLetter(run) {
			break
		}
		sb.WriteRune(run)
		lexer.Position++
	}
	return sb.String()
}

func (lexer *Lexer) resume() string {
	var lexerline string
	runes := lexer.File.Content[lexer.Position:]
	for i, r := range runes {
		if unicode.IsSpace(r) {
			lexer.Column++
			lexer.Position++
			if r == '\n' || r == '\r' {
				lexer.Line++
				lexer.Column = 1
			}
			continue
		}
		lexerline = string(runes[i:])
		break
	}
	return lexerline
}

func (lexer *Lexer) Token() Token {
	tk := Token{
		File: lexer.File,
		Type: NA,
	}
	lexerline := lexer.resume()
	if lexerline == "" {
		return tk
	}

	tk.Column = lexer.Column
	tk.Line = lexer.Line

	switch {
	case lexerline[0] == '(':
		tk.Value = "("
		tk.Type = Brace
		lexer.Position++
	case lexerline[0] == ')':
		tk.Value = ")"
		tk.Type = Brace
		lexer.Position++
	case lexerline[0] == '{':
		tk.Value = "{"
		tk.Type = Brace
		lexer.Position++
	case lexerline[0] == '}':
		tk.Value = "}"
		tk.Type = Brace
		lexer.Position++
	case isKeyword(lexerline, "int16"):
		tk.Value = "int16"
		tk.Type = Type
		lexer.Position += 5
	case isKeyword(lexerline, "int32"):
		tk.Value = "int32"
		tk.Type = Type
		lexer.Position += 5
	default:
		if chk := lexer.lexName(lexerline); chk != "" {
			tk.Value = chk
			tk.Type = Name
			break
		}
		lexer.error("invalid_token")
		lexer.Column++
		lexer.Position++
		return tk
	}
	lexer.Column += len(tk.Value)
	return tk
}
