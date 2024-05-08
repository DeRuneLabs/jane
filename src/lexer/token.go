// Copyright (c) 2024 arfy slowy - DeRuneLabs
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package lexer

import (
	"strings"
	"unicode/utf8"
)

const (
	ID_NA        = 0
	ID_DT        = 1
	ID_IDENT     = 2
	ID_BRACE     = 3
	ID_RET       = 4
	ID_SEMICOLON = 5
	ID_LITERAL   = 6
	ID_OP        = 7
	ID_COMMA     = 8
	ID_CONST     = 9
	ID_TYPE      = 10
	ID_COLON     = 11
	ID_ITER      = 12
	ID_BREAK     = 13
	ID_CONTINUE  = 14
	ID_IN        = 15
	ID_IF        = 16
	ID_ELSE      = 17
	ID_COMMENT   = 18
	ID_USE       = 19
	ID_DOT       = 20
	ID_PUB       = 21
	ID_GOTO      = 22
	ID_DBLCOLON  = 23
	ID_ENUM      = 24
	ID_STRUCT    = 25
	ID_CO        = 26
	ID_MATCH     = 27
	ID_SELF      = 28
	ID_TRAIT     = 29
	ID_IMPL      = 30
	ID_CPP       = 31
	ID_FALL      = 32
	ID_FN        = 33
	ID_LET       = 34
	ID_UNSAFE    = 35
	ID_MUT       = 36
	ID_DEFER     = 37
)

const (
	KND_DBLCOLON     = "::"
	KND_COLON        = ":"
	KND_SEMICOLON    = ";"
	KND_COMMA        = ","
	KND_TRIPLE_DOT   = "..."
	KND_DOT          = "."
	KND_PLUS_EQ      = "+="
	KND_MINUS_EQ     = "-="
	KND_STAR_EQ      = "*="
	KND_SOLIDUS_EQ   = "/="
	KND_PERCENT_EQ   = "%="
	KND_LSHIFT_EQ    = "<<="
	KND_RSHIFT_EQ    = ">>="
	KND_CARET_EQ     = "^="
	KND_AMPER_EQ     = "&="
	KND_VLINE_EQ     = "|="
	KND_EQS          = "=="
	KND_NOT_EQ       = "!="
	KND_GREAT_EQ     = ">="
	KND_LESS_EQ      = "<="
	KND_DBL_AMPER    = "&&"
	KND_DBL_VLINE    = "||"
	KND_LSHIFT       = "<<"
	KND_RSHIFT       = ">>"
	KND_DBL_PLUS     = "++"
	KND_DBL_MINUS    = "--"
	KND_PLUS         = "+"
	KND_MINUS        = "-"
	KND_STAR         = "*"
	KND_SOLIDUS      = "/"
	KND_PERCENT      = "%"
	KND_AMPER        = "&"
	KND_VLINE        = "|"
	KND_CARET        = "^"
	KND_EXCL         = "!"
	KND_LT           = "<"
	KND_GT           = ">"
	KND_EQ           = "="
	KND_LN_COMMENT   = "//"
	KND_RNG_LCOMMENT = "/*"
	KND_RNG_RCOMMENT = "*/"
	KND_LPAREN       = "("
	KND_RPARENT      = ")"
	KND_LBRACKET     = "["
	KND_RBRACKET     = "]"
	KND_LBRACE       = "{"
	KND_RBRACE       = "}"
	KND_I8           = "i8"
	KND_I16          = "i16"
	KND_I32          = "i32"
	KND_I64          = "i64"
	KND_U8           = "u8"
	KND_U16          = "u16"
	KND_U32          = "u32"
	KND_U64          = "u64"
	KND_F32          = "f32"
	KND_F64          = "f64"
	KND_UINT         = "uint"
	KND_INT          = "int"
	KND_UINTPTR      = "uintptr"
	KND_BOOL         = "bool"
	KND_STR          = "str"
	KND_ANY          = "any"
	KND_TRUE         = "true"
	KND_FALSE        = "false"
	KND_NIL          = "nil"
	KND_CONST        = "const"
	KND_RET          = "ret"
	KND_TYPE         = "type"
	KND_ITER         = "for"
	KND_BREAK        = "break"
	KND_CONTINUE     = "continue"
	KND_IN           = "in"
	KND_IF           = "if"
	KND_ELSE         = "else"
	KND_USE          = "use"
	KND_PUB          = "pub"
	KND_GOTO         = "goto"
	KND_ENUM         = "enum"
	KND_STRUCT       = "struct"
	KND_CO           = "co"
	KND_MATCH        = "match"
	KND_SELF         = "self"
	KND_TRAIT        = "trait"
	KND_IMPL         = "impl"
	KND_CPP          = "cpp"
	KND_FALL         = "fall"
	KND_FN           = "fn"
	KND_LET          = "let"
	KND_UNSAFE       = "unsafe"
	KND_MUT          = "mut"
	KND_DEFER        = "defer"
)

const (
	IGNORE_ID    = "_"
	ANONYMOUS_ID = "<anonymous>"
)

const (
	COMMENT_PRAGMA_SEP    = ":"
	PRAGMA_COMMENT_PREFIX = "jane" + COMMENT_PRAGMA_SEP
)

const (
	MARK_ARRAY   = "..."
	PREFIX_SLICE = "[]"
	PREFIX_ARRAY = "[" + MARK_ARRAY + "]"
)

var PUNCTS = [...]rune{
	'!',
	'#',
	'$',
	',',
	'.',
	'\'',
	'"',
	':',
	';',
	'<',
	'>',
	'=',
	'?',
	'-',
	'+',
	'*',
	'(',
	')',
	'[',
	']',
	'{',
	'}',
	'%',
	'&',
	'/',
	'\\',
	'@',
	'^',
	'_',
	'`',
	'|',
	'~',
	'¦',
}

var SPACES = [...]rune{
	' ',
	'\t',
	'\v',
	'\r',
	'\n',
}

type Token struct {
	File   *File
	Row    int
	Column int
	Kind   string
	Id     uint8
}

func (t *Token) Prec() int {
	if t.Id != ID_OP {
		return -1
	}
	switch t.Kind {
	case KND_STAR, KND_PERCENT, KND_SOLIDUS, KND_RSHIFT, KND_LSHIFT, KND_AMPER:
		return 5
	case KND_PLUS, KND_MINUS, KND_VLINE, KND_CARET:
		return 4
	case KND_EQS, KND_NOT_EQ, KND_LT, KND_LESS_EQ, KND_GT, KND_GREAT_EQ:
		return 3
	case KND_DBL_AMPER:
		return 2
	case KND_DBL_VLINE:
		return 1
	default:
		return -1
	}
}

func IsStr(k string) bool {
	return k != "" && (k[0] == '"' || IsRawStr(k))
}

func IsRawStr(k string) bool {
	return k != "" && k[0] == '`'
}

func IsRune(k string) bool {
	return k != "" && k[0] == '\''
}

func IsNil(k string) bool {
	return k == KND_NIL
}

func IsBool(k string) bool {
	return k == KND_TRUE || k == KND_FALSE
}

func contains_any(s string, bytes string) bool {
	for _, b := range bytes {
		i := strings.Index(s, string(b))
		if i >= 0 {
			return true
		}
	}
	return false
}

func IsFloat(k string) bool {
	if strings.HasPrefix(k, "0x") {
		return contains_any(k, ".pP")
	}
	return contains_any(k, ".eE")
}

func IsNum(k string) bool {
	if k == "" {
		return false
	}
	return k[0] == '-' || (k[0] >= '0' && k[0] <= '9')
}

func IsLiteral(k string) bool {
	return IsNum(k) || IsStr(k) || IsRune(k) || IsNil(k) || IsBool(k)
}

func IsIgnoreId(id string) bool {
	return id == IGNORE_ID
}

func IsAnonymousId(id string) bool {
	return id == ANONYMOUS_ID
}

func rune_exist(r rune, runes []rune) bool {
	for _, cr := range runes {
		if r == cr {
			return true
		}
	}
	return false
}

func IsPunct(r rune) bool {
	return rune_exist(r, PUNCTS[:])
}

func IsSpace(r rune) bool {
	return rune_exist(r, SPACES[:])
}

func IsLetter(r rune) bool {
	return ('a' <= r && r <= 'z') || ('A' <= r && r <= 'Z')
}

func IsIdentifierRune(s string) bool {
	if s == "" {
		return false
	}
	if s[0] != '_' {
		r, _ := utf8.DecodeRuneInString(s)
		if !IsLetter(r) {
			return false
		}
	}
	return true
}

func IsDecimal(b byte) bool {
	return '0' <= b && b <= '9'
}

func IsBinary(b byte) bool {
	return b == '0' || b == '1'
}

func IsOctal(b byte) bool {
	return '0' <= b && b <= '7'
}

func IsHex(b byte) bool {
	switch {
	case '0' <= b && b <= '9':
		return true
	case 'a' <= b && b <= 'f':
		return true
	case 'A' <= b && b <= 'F':
		return true
	default:
		return false
	}
}
