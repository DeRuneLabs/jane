// Copyright (c) 2024 - DeRuneLabs
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

package jntype

import "github.com/DeRuneLabs/jane/lexer/tokens"

// builtin data type constant
const (
	Void    uint8 = 0
	I8      uint8 = 1
	I16     uint8 = 2
	I32     uint8 = 3
	I64     uint8 = 4
	U8      uint8 = 5
	U16     uint8 = 6
	U32     uint8 = 7
	U64     uint8 = 8
	Bool    uint8 = 9
	Str     uint8 = 10
	F32     uint8 = 11
	F64     uint8 = 12
	Any     uint8 = 13
	Id      uint8 = 14
	Fn      uint8 = 15
	Nil     uint8 = 16
	UInt    uint8 = 17
	Int     uint8 = 18
	Map     uint8 = 19
	UIntptr uint8 = 20
	Enum    uint8 = 21
	Struct  uint8 = 22
	Trait   uint8 = 23
	Slice   uint8 = 24
	Array   uint8 = 25
	Unsafe  uint8 = 26
)

var TypeMap = map[uint8]string{
	Void:    VoidTypeStr,
	Nil:     NilTypeStr,
	I8:      tokens.I8,
	I16:     tokens.I16,
	I32:     tokens.I32,
	I64:     tokens.I64,
	U8:      tokens.U8,
	U16:     tokens.U16,
	U32:     tokens.U32,
	U64:     tokens.U64,
	Str:     tokens.STR,
	Bool:    tokens.BOOL,
	F32:     tokens.F32,
	F64:     tokens.F64,
	Any:     tokens.ANY,
	UInt:    tokens.UINT,
	Int:     tokens.INT,
	UIntptr: tokens.UINTPTR,
	Unsafe:  tokens.UNSAFE,
}
