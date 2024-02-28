package lexer

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/De-Rune/jane/package/jn"
)

func (lex *Lexer) pushError(err string) {
	lex.Errors = append(lex.Errors, fmt.Sprintf("%s %d:%d %s", lex.File.Path, lex.Line, lex.Column, jn.Errors[err]))
}

func (lex *Lexer) Tokenize() []Token {
	var tokens []Token
	lex.Errors = nil
	for lex.Position < len(lex.File.Content) {
		token := lex.Token()
		if token.Id != NA {
			tokens = append(tokens, token)
		}
	}
	return tokens
}

func isKeyword(lexerline, kw string) bool {
	if !strings.HasPrefix(lexerline, kw) {
		return false
	}
	lexerline = lexerline[len(kw):]
	switch {
	case lexerline == "", unicode.IsSpace(rune(lexerline[0])), unicode.IsPunct(rune(lexerline[0])):
		return true
	}
	return false
}

func (lex *Lexer) lexName(lexerline string) string {
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
		lex.Position++
	}
	return sb.String()
}

func (lex *Lexer) resume() string {
	var lexerline string
	runes := lex.File.Content[lex.Position:]
	for i, r := range runes {
		if unicode.IsSpace(r) {
			lex.Column++
			lex.Position++
			if r == '\n' {
				lex.NewLine()
			}
			continue
		}
		lexerline = string(runes[i:])
		break
	}
	return lexerline
}

func (lex *Lexer) lexLineComment() {
	lex.Position += 2
	for ; lex.Position < len(lex.File.Content); lex.Position++ {
		if lex.File.Content[lex.Position] == '\n' {
			lex.Position++
			lex.NewLine()
			return
		}
	}
}

func (lex *Lexer) lexBlockComment() {
	lex.Position += 2
	for ; lex.Position < len(lex.File.Content); lex.Position++ {
		run := lex.File.Content[lex.Position]
		if run == '\n' {
			lex.NewLine()
			continue
		}
		lex.Column += len(string(run))
		if strings.HasPrefix(string(lex.File.Content[lex.Position:]), "*/") {
			lex.Column += 2
			lex.Position += 2
			return
		}
	}
	lex.pushError("missing_block_comment")
}

var numericRegexp = *regexp.MustCompile(`^((0x[[:xdigit:]]+)|(\d+((\.\d+)?((e|E)(\-|\+|)\d+)?|(\.\d+))))`)

func (lex *Lexer) lexNumeric(content string) string {
	value := numericRegexp.FindString(content)
	lex.Position += len(value)
	return value
}

var escapeSequenceRegexp = regexp.MustCompile(`^\\([\\'"abfnrtv]|U.{8}|u.{4}|x..|[0-7]{1,3})`)

func (lex *Lexer) getEscapeSequence(content string) string {
	seq := escapeSequenceRegexp.FindString(content)
	if seq != "" {
		lex.Position += len(seq)
		return seq
	}
	lex.Position++
	lex.pushError("invalid_escape_sequence")
	return seq
}

func (lex *Lexer) getRune(content string) string {
	if content[0] == '\\' {
		return lex.getEscapeSequence(content)
	}
	run, _ := utf8.DecodeRuneInString(content)
	lex.Position++
	return string(run)
}

func (lex *Lexer) lexRune(content string) string {
	var sb strings.Builder
	sb.WriteByte('\'')
	lex.Column++
	content = content[1:]
	count := 0
	for index := 0; index < len(content); index++ {
		if content[index] == '\n' {
			lex.pushError("missing_rune_end")
			lex.Position++
			lex.NewLine()
			return ""
		}
		run := lex.getRune(content[index:])
		sb.WriteString(run)
		length := len(run)
		lex.Column += length
		if run == "'" {
			lex.Position++
			break
		}
		if length > 1 {
			index += length - 1
		}
		count++
	}
	if count == 0 {
		lex.pushError("rune_empty")
	} else if count > 1 {
		lex.pushError("rune_overflow")
	}
	return sb.String()
}

func (lex *Lexer) lexString(content string) string {
	var sb strings.Builder
	sb.WriteByte('"')
	lex.Column++
	content = content[1:]
	for index, run := range content {
		if run == '\n' {
			lex.pushError("missing_string_end")
			lex.Position++
			lex.NewLine()
			return ""
		}
		run := lex.getRune(content[index:])
		sb.WriteString(run)
		length := len(run)
		lex.Column += length
		if run == `"` {
			lex.Position++
			break
		}
		if length > 1 {
			index += length - 1
		}
	}
	return sb.String()
}

func (lex *Lexer) NewLine() {
	lex.Line++
	lex.Column = 1
}

func (lex *Lexer) lexPunct(content, kind string, id uint8, token *Token) bool {
	if !strings.HasPrefix(content, kind) {
		return false
	}
	token.Kind = kind
	token.Id = id
	lex.Position += len([]rune(kind))
	return true
}

func (lex *Lexer) lexKeyword(content, kind string, id uint8, token *Token) bool {
	if !isKeyword(content, kind) {
		return false
	}
	token.Kind = kind
	token.Id = id
	lex.Position += len([]rune(kind))
	return true
}

func (lex *Lexer) Token() Token {
	token := Token{
		File: lex.File,
		Id:   NA,
	}
	content := lex.resume()
	if content == "" {
		return token
	}
	token.Column = lex.Column
	token.Row = lex.Line

	switch {
	case content[0] == '\'':
		token.Kind = lex.lexRune(content)
		token.Id = Value
		return token
	case content[0] == '"':
		token.Kind = lex.lexString(content)
		token.Id = Value
		return token
	case strings.HasPrefix(content, "//"):
		lex.lexLineComment()
		return token
	case strings.HasPrefix(content, "/*"):
		lex.lexBlockComment()
		return token
	case
		lex.lexPunct(content, ":", Colon, &token),
		lex.lexPunct(content, ";", SemiColon, &token),
		lex.lexPunct(content, ",", Comma, &token),
		lex.lexPunct(content, "@", At, &token),
		lex.lexPunct(content, "(", Brace, &token),
		lex.lexPunct(content, ")", Brace, &token),
		lex.lexPunct(content, "{", Brace, &token),
		lex.lexPunct(content, "}", Brace, &token),
		lex.lexPunct(content, "[", Brace, &token),
		lex.lexPunct(content, "]", Brace, &token),
		lex.lexPunct(content, "+=", Operator, &token),
		lex.lexPunct(content, "-=", Operator, &token),
		lex.lexPunct(content, "*=", Operator, &token),
		lex.lexPunct(content, "/=", Operator, &token),
		lex.lexPunct(content, "%=", Operator, &token),
		lex.lexPunct(content, "<<=", Operator, &token),
		lex.lexPunct(content, ">>=", Operator, &token),
		lex.lexPunct(content, "^=", Operator, &token),
		lex.lexPunct(content, "&=", Operator, &token),
		lex.lexPunct(content, "|=", Operator, &token),
		lex.lexPunct(content, "==", Operator, &token),
		lex.lexPunct(content, "!=", Operator, &token),
		lex.lexPunct(content, ">=", Operator, &token),
		lex.lexPunct(content, "<=", Operator, &token),
		lex.lexPunct(content, "&&", Operator, &token),
		lex.lexPunct(content, "||", Operator, &token),
		lex.lexPunct(content, "<<", Operator, &token),
		lex.lexPunct(content, ">>", Operator, &token),
		lex.lexPunct(content, "+", Operator, &token),
		lex.lexPunct(content, "-", Operator, &token),
		lex.lexPunct(content, "*", Operator, &token),
		lex.lexPunct(content, "/", Operator, &token),
		lex.lexPunct(content, "%", Operator, &token),
		lex.lexPunct(content, "~", Operator, &token),
		lex.lexPunct(content, "&", Operator, &token),
		lex.lexPunct(content, "|", Operator, &token),
		lex.lexPunct(content, "^", Operator, &token),
		lex.lexPunct(content, "!", Operator, &token),
		lex.lexPunct(content, "<", Operator, &token),
		lex.lexPunct(content, ">", Operator, &token),
		lex.lexPunct(content, "=", Operator, &token),
		lex.lexKeyword(content, "i8", DataType, &token),
		lex.lexKeyword(content, "i16", DataType, &token),
		lex.lexKeyword(content, "i32", DataType, &token),
		lex.lexKeyword(content, "i64", DataType, &token),
		lex.lexKeyword(content, "u8", DataType, &token),
		lex.lexKeyword(content, "u16", DataType, &token),
		lex.lexKeyword(content, "u32", DataType, &token),
		lex.lexKeyword(content, "u64", DataType, &token),
		lex.lexKeyword(content, "f32", DataType, &token),
		lex.lexKeyword(content, "f64", DataType, &token),
		lex.lexKeyword(content, "bool", DataType, &token),
		lex.lexKeyword(content, "rune", DataType, &token),
		lex.lexKeyword(content, "str", DataType, &token),
		lex.lexKeyword(content, "true", Value, &token),
		lex.lexKeyword(content, "false", Value, &token),
		lex.lexKeyword(content, "nil", Value, &token),
		lex.lexKeyword(content, "const", Const, &token),
		lex.lexKeyword(content, "ret", Return, &token),
		lex.lexKeyword(content, "type", Type, &token),
		lex.lexKeyword(content, "new", New, &token):
	default:
		l := lex.lexName(content)
		if l != "" {
			token.Kind = "_" + l
			token.Id = Name
			break
		}
		l = lex.lexNumeric(content)
		if l != "" {
			token.Kind = l
			token.Id = Value
			break
		}
		lex.pushError("invalid_token")
		lex.Column++
		lex.Position++
		return token
	}
	lex.Column += len(token.Kind)
	return token
}
