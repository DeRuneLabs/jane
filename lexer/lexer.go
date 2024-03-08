package lexer

import (
	"regexp"
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"

	"github.com/DeRuneLabs/jane/package/jn"
	"github.com/DeRuneLabs/jane/package/jnio"
	"github.com/DeRuneLabs/jane/package/jnlog"
)

type Lexer struct {
	wg               sync.WaitGroup
	firstTokenOfLine bool

	File   *jnio.File
	Pos    int
	Column int
	Row    int
	Logs   []jnlog.CompilerLog
	braces []Token
}

func NewLexer(f *jnio.File) *Lexer {
	lex := new(Lexer)
	lex.File = f
	lex.Pos = 0
	lex.Row = -1
	lex.NewLine()
	return lex
}

func (lex *Lexer) pusherr(key string, args ...interface{}) {
	lex.Logs = append(lex.Logs, jnlog.CompilerLog{
		Type:   jnlog.Err,
		Row:    lex.Row,
		Column: lex.Column,
		Path:   lex.File.Path,
		Msg:    jn.GetErr(key, args...),
	})
}

func (lex *Lexer) pusherrtok(token Token, err string) {
	lex.Logs = append(lex.Logs, jnlog.CompilerLog{
		Type:   jnlog.Err,
		Row:    token.Row,
		Column: token.Column,
		Path:   lex.File.Path,
		Msg:    jn.Errs[err],
	})
}

func (lex *Lexer) Lexer() []Token {
	var toks []Token
	lex.Logs = nil
	lex.NewLine()
	for lex.Pos < len(lex.File.Text) {
		tok := lex.Token()
		if tok.Id != NA {
			toks = append(toks, tok)
		}
	}
	lex.wg.Add(1)
	go lex.checkRangesAsync()
	lex.wg.Wait()
	return toks
}

func (lex *Lexer) checkRangesAsync() {
	defer func() { lex.wg.Done() }()
	for _, token := range lex.braces {
		switch token.Kind {
		case "(":
			lex.pusherrtok(token, "wait_close_parentheses")
		case "{":
			lex.pusherrtok(token, "wait_close_brace")
		case "[":
			lex.pusherrtok(token, "wait_close_bracket")
		}
	}
}

func iskw(lexerline, kw string) bool {
	if !strings.HasPrefix(lexerline, kw) {
		return false
	}
	lexerline = lexerline[len(kw):]
	return lexerline == "" ||
		unicode.IsSpace(rune(lexerline[0])) ||
		unicode.IsPunct(rune(lexerline[0]))
}

func (lex *Lexer) id(lexerline string) string {
	if lexerline[0] != '_' {
		r, _ := utf8.DecodeRuneInString(lexerline)
		if !unicode.IsLetter(r) {
			return ""
		}
	}
	var sb strings.Builder
	for _, r := range lexerline {
		if r != '_' &&
			('0' > r || '9' < r) &&
			!unicode.IsLetter(r) {
			break
		}
		sb.WriteRune(r)
		lex.Pos++
	}
	return sb.String()
}

func (lex *Lexer) resume() string {
	var lexerline string
	runes := lex.File.Text[lex.Pos:]
	for i, r := range runes {
		if unicode.IsSpace(r) {
			lex.Pos++
			if r == '\n' {
				lex.NewLine()
			} else {
				lex.Column++
			}
			continue
		}
		lexerline = string(runes[i:])
		break
	}
	return lexerline
}

func (lex *Lexer) lncomment(token *Token) {
	start := lex.Pos
	lex.Pos += 2
	for ; lex.Pos < len(lex.File.Text); lex.Pos++ {
		if lex.File.Text[lex.Pos] == '\n' {
			if lex.firstTokenOfLine {
				token.Id = Comment
				token.Kind = string(lex.File.Text[start:lex.Pos])
			}
			return
		}
	}
	if lex.firstTokenOfLine {
		token.Id = Comment
		token.Kind = string(lex.File.Text[start:])
	}
}

func (lex *Lexer) rangecomment() {
	lex.Pos += 2
	for ; lex.Pos < len(lex.File.Text); lex.Pos++ {
		run := lex.File.Text[lex.Pos]
		if run == '\n' {
			lex.NewLine()
			continue
		}
		lex.Column += len(string(run))
		if strings.HasPrefix(string(lex.File.Text[lex.Pos:]), "*/") {
			lex.Column += 2
			lex.Pos += 2
			return
		}
	}
	lex.pusherr("missing_block_comment")
}

var numRegexp = *regexp.MustCompile(`^((0x[[:xdigit:]]+)|(\d+((\.\d+)?((e|E)(\-|\+|)\d+)?|(\.\d+))))`)

func (lex *Lexer) num(txt string) string {
	val := numRegexp.FindString(txt)
	lex.Pos += len(val)
	return val
}

var escSeqRegexp = regexp.MustCompile(`^\\([\\'"abfnrtv]|U.{8}|u.{4}|x..|[0-7]{1,3})`)

func (lex *Lexer) escseq(txt string) string {
	seq := escSeqRegexp.FindString(txt)
	if seq != "" {
		lex.Pos += len([]rune(seq))
		return seq
	}
	lex.Pos++
	lex.pusherr("invalid_escape_sequence")
	return seq
}

func (lex *Lexer) getrune(txt string, raw bool) string {
	if !raw && txt[0] == '\\' {
		return lex.escseq(txt)
	}
	run, _ := utf8.DecodeRuneInString(txt)
	lex.Pos++
	return string(run)
}

func (lex *Lexer) rune(txt string) string {
	var sb strings.Builder
	sb.WriteByte('\'')
	lex.Column++
	txt = txt[1:]
	count := 0
	for i := 0; i < len(txt); i++ {
		if txt[i] == '\n' {
			lex.pusherr("missing_rune_end")
			lex.Pos++
			lex.NewLine()
			return ""
		}
		run := lex.getrune(txt[i:], false)
		sb.WriteString(run)
		length := len(run)
		lex.Column += length
		if run == "'" {
			lex.Pos++
			break
		}
		if length > 1 {
			i += length - 1
		}
		count++
	}
	if count == 0 {
		lex.pusherr("rune_empty")
	} else if count > 1 {
		lex.pusherr("rune_overflow")
	}
	return sb.String()
}

func (lex *Lexer) str(txt string) string {
	var sb strings.Builder
	mark := txt[0]
	raw := mark == '`'
	sb.WriteByte(mark)
	lex.Column++
	txt = txt[1:]
	for i := 0; i < len(txt); i++ {
		ch := txt[i]
		if ch == '\n' {
			defer lex.NewLine()
			if !raw {
				lex.pusherr("missing_string_end")
				lex.Pos++
				return ""
			}
		}
		run := lex.getrune(txt[i:], raw)
		sb.WriteString(run)
		length := len(run)
		lex.Column += length
		if ch == mark {
			lex.Pos++
			break
		}
		if length > 1 {
			i += length - 1
		}
	}
	return sb.String()
}

func (lex *Lexer) NewLine() {
	lex.firstTokenOfLine = true
	lex.Row++
	lex.Column = 1
}

func (lex *Lexer) punct(txt, kind string, id uint8, token *Token) bool {
	if !strings.HasPrefix(txt, kind) {
		return false
	}
	token.Kind = kind
	token.Id = id
	lex.Pos += len([]rune(kind))
	return true
}

func (lex *Lexer) kw(txt, kind string, id uint8, tok *Token) bool {
	if !iskw(txt, kind) {
		return false
	}
	tok.Kind = kind
	tok.Id = id
	lex.Pos += len([]rune(kind))
	return true
}

func (lex *Lexer) Token() Token {
	defer func() { lex.firstTokenOfLine = false }()

	tok := Token{File: lex.File, Id: NA}

	txt := lex.resume()
	if txt == "" {
		return tok
	}

	tok.Column = lex.Column
	tok.Row = lex.Row

	switch {
	case txt[0] == '\'':
		tok.Kind = lex.rune(txt)
		tok.Id = Value
		return tok
	case txt[0] == '"', txt[0] == '`':
		tok.Kind = lex.str(txt)
		tok.Id = Value
		return tok
	case strings.HasPrefix(txt, "//"):
		lex.lncomment(&tok)
		goto ret
	case strings.HasPrefix(txt, "/*"):
		lex.rangecomment()
		return tok
	case lex.punct(txt, "(", Brace, &tok):
		lex.braces = append(lex.braces, tok)
	case lex.punct(txt, ")", Brace, &tok):
		len := len(lex.braces)
		if len == 0 {
			lex.pusherrtok(tok, "extra_closed_parentheses")
			break
		} else if lex.braces[len-1].Kind != "(" {
			lex.wg.Add(1)
			go lex.pushWrongOrderCloseErrAsync(tok)
		}
		lex.rmrange(len-1, tok.Kind)
	case lex.punct(txt, "{", Brace, &tok):
		lex.braces = append(lex.braces, tok)
	case lex.punct(txt, "}", Brace, &tok):
		len := len(lex.braces)
		if len == 0 {
			lex.pusherrtok(tok, "extra_closed_braces")
			break
		} else if lex.braces[len-1].Kind != "{" {
			lex.wg.Add(1)
			go lex.pushWrongOrderCloseErrAsync(tok)
		}
		lex.rmrange(len-1, tok.Kind)
	case lex.punct(txt, "[", Brace, &tok):
		lex.braces = append(lex.braces, tok)
	case lex.punct(txt, "]", Brace, &tok):
		len := len(lex.braces)
		if len == 0 {
			lex.pusherrtok(tok, "extra_closed_brackets")
			break
		} else if lex.braces[len-1].Kind != "[" {
			lex.wg.Add(1)
			go lex.pushWrongOrderCloseErrAsync(tok)
		}
		lex.rmrange(len-1, tok.Kind)
	case
		lex.firstTokenOfLine && lex.punct(txt, "#", Preprocessor, &tok),
		lex.punct(txt, ":", Colon, &tok),
		lex.punct(txt, ";", SemiColon, &tok),
		lex.punct(txt, ",", Comma, &tok),
		lex.punct(txt, "@", At, &tok),
		lex.punct(txt, "...", Operator, &tok),
		lex.punct(txt, ".", Dot, &tok),
		lex.punct(txt, "+=", Operator, &tok),
		lex.punct(txt, "-=", Operator, &tok),
		lex.punct(txt, "*=", Operator, &tok),
		lex.punct(txt, "/=", Operator, &tok),
		lex.punct(txt, "%=", Operator, &tok),
		lex.punct(txt, "<<=", Operator, &tok),
		lex.punct(txt, ">>=", Operator, &tok),
		lex.punct(txt, "^=", Operator, &tok),
		lex.punct(txt, "&=", Operator, &tok),
		lex.punct(txt, "|=", Operator, &tok),
		lex.punct(txt, "==", Operator, &tok),
		lex.punct(txt, "!=", Operator, &tok),
		lex.punct(txt, ">=", Operator, &tok),
		lex.punct(txt, "<=", Operator, &tok),
		lex.punct(txt, "&&", Operator, &tok),
		lex.punct(txt, "||", Operator, &tok),
		lex.punct(txt, "<<", Operator, &tok),
		lex.punct(txt, ">>", Operator, &tok),
		lex.punct(txt, "+", Operator, &tok),
		lex.punct(txt, "-", Operator, &tok),
		lex.punct(txt, "*", Operator, &tok),
		lex.punct(txt, "/", Operator, &tok),
		lex.punct(txt, "%", Operator, &tok),
		lex.punct(txt, "~", Operator, &tok),
		lex.punct(txt, "&", Operator, &tok),
		lex.punct(txt, "|", Operator, &tok),
		lex.punct(txt, "^", Operator, &tok),
		lex.punct(txt, "!", Operator, &tok),
		lex.punct(txt, "<", Operator, &tok),
		lex.punct(txt, ">", Operator, &tok),
		lex.punct(txt, "=", Operator, &tok),
		lex.kw(txt, "i8", DataType, &tok),
		lex.kw(txt, "i16", DataType, &tok),
		lex.kw(txt, "i32", DataType, &tok),
		lex.kw(txt, "i64", DataType, &tok),
		lex.kw(txt, "u8", DataType, &tok),
		lex.kw(txt, "u16", DataType, &tok),
		lex.kw(txt, "u32", DataType, &tok),
		lex.kw(txt, "u64", DataType, &tok),
		lex.kw(txt, "f32", DataType, &tok),
		lex.kw(txt, "f64", DataType, &tok),
		lex.kw(txt, "byte", DataType, &tok),
		lex.kw(txt, "sbyte", DataType, &tok),
		lex.kw(txt, "size", DataType, &tok),
		lex.kw(txt, "ssize", DataType, &tok),
		lex.kw(txt, "bool", DataType, &tok),
		lex.kw(txt, "rune", DataType, &tok),
		lex.kw(txt, "str", DataType, &tok),
		lex.kw(txt, "true", Value, &tok),
		lex.kw(txt, "false", Value, &tok),
		lex.kw(txt, "nil", Value, &tok),
		lex.kw(txt, "const", Const, &tok),
		lex.kw(txt, "ret", Ret, &tok),
		lex.kw(txt, "type", Type, &tok),
		lex.kw(txt, "new", New, &tok),
		lex.kw(txt, "free", Free, &tok),
		lex.kw(txt, "iter", Iter, &tok),
		lex.kw(txt, "break", Break, &tok),
		lex.kw(txt, "continue", Continue, &tok),
		lex.kw(txt, "in", In, &tok),
		lex.kw(txt, "if", If, &tok),
		lex.kw(txt, "else", Else, &tok),
		lex.kw(txt, "volatile", Volatile, &tok),
		lex.kw(txt, "use", Use, &tok),
		lex.kw(txt, "pub", Pub, &tok):
	default:
		l := lex.id(txt)
		if l != "" {
			tok.Kind = l
			tok.Id = Id
			break
		}
		l = lex.num(txt)
		if l != "" {
			tok.Kind = l
			tok.Id = Value
			break
		}
		r, sz := utf8.DecodeRuneInString(txt)
		lex.pusherr("invalid_token", r)
		lex.Column += sz
		lex.Pos++
		return tok
	}
	lex.Column += len(tok.Kind)
ret:
	return tok
}

func (lex *Lexer) rmrange(i int, kind string) {
	var close string
	switch kind {
	case ")":
		close = "("
	case "}":
		close = "{"
	case "]":
		close = "["
	}
	for ; i >= 0; i-- {
		tok := lex.braces[i]
		if tok.Kind != close {
			continue
		}
		lex.braces = append(lex.braces[:i], lex.braces[i+1:]...)
		break
	}
}

func (lex *Lexer) pushWrongOrderCloseErrAsync(tok Token) {
	defer func() { lex.wg.Done() }()
	var msg string
	switch lex.braces[len(lex.braces)-1].Kind {
	case "(":
		msg = "expected_parentheses_close"
	case "{":
		msg = "expected_brace_close"
	case "[":
		msg = "expected_bracket_close"
	}
	lex.pusherrtok(tok, msg)
}
