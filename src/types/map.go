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

package types

import "github.com/DeRuneLabs/jane/lexer"

const (
	VOID    uint8 = 0
	I8      uint8 = 1
	I16     uint8 = 2
	I32     uint8 = 3
	I64     uint8 = 4
	U8      uint8 = 5
	U16     uint8 = 6
	U32     uint8 = 7
	U64     uint8 = 8
	BOOL    uint8 = 9
	STR     uint8 = 10
	F32     uint8 = 11
	F64     uint8 = 12
	ANY     uint8 = 13
	ID      uint8 = 14
	FN      uint8 = 15
	NIL     uint8 = 16
	UINT    uint8 = 17
	INT     uint8 = 18
	MAP     uint8 = 19
	UINTPTR uint8 = 20
	ENUM    uint8 = 21
	STRUCT  uint8 = 22
	TRAIT   uint8 = 23
	SLICE   uint8 = 24
	ARRAY   uint8 = 25
	UNSAFE  uint8 = 26
)

var TYPE_MAP = map[uint8]string{
	VOID:    VOID_TYPE_STR,
	NIL:     NIL_TYPE_STR,
	I8:      lexer.KND_I8,
	I16:     lexer.KND_I16,
	I32:     lexer.KND_I32,
	I64:     lexer.KND_I64,
	U8:      lexer.KND_U8,
	U16:     lexer.KND_U16,
	U32:     lexer.KND_U32,
	U64:     lexer.KND_U64,
	STR:     lexer.KND_STR,
	BOOL:    lexer.KND_BOOL,
	F32:     lexer.KND_F32,
	F64:     lexer.KND_F64,
	ANY:     lexer.KND_ANY,
	UINT:    lexer.KND_UINT,
	INT:     lexer.KND_INT,
	UINTPTR: lexer.KND_UINTPTR,
	UNSAFE:  lexer.KND_UNSAFE,
}
