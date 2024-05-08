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

package ast

import (
	"github.com/DeRuneLabs/jane/build"
	"github.com/DeRuneLabs/jane/lexer"
)

const nil_type_str = "<nil>"
const void_type_str = "<void>"

const void_t = 0
const i8_t = 1
const i16_t = 2
const i32_t = 3
const i64_t = 4
const u8_t = 5
const u16_t = 6
const u32_t = 7
const u64_t = 8
const bool_t = 9
const str_t = 10
const f32_t = 11
const f64_t = 12
const any_t = 13
const id_t = 14
const fn_t = 15
const nil_t = 16
const uint_t = 17
const int_t = 18
const map_t = 19
const uintptr_t = 20
const enum_t = 21
const struct_t = 22
const trait_t = 23
const slice_t = 24
const array_t = 25
const unsafe_t = 26

var type_map = map[uint8]string{
	void_t:    void_type_str,
	nil_t:     nil_type_str,
	i8_t:      lexer.KND_I8,
	i16_t:     lexer.KND_I16,
	i32_t:     lexer.KND_I32,
	i64_t:     lexer.KND_I64,
	u8_t:      lexer.KND_U8,
	u16_t:     lexer.KND_U16,
	u32_t:     lexer.KND_U32,
	u64_t:     lexer.KND_U64,
	str_t:     lexer.KND_STR,
	bool_t:    lexer.KND_BOOL,
	f32_t:     lexer.KND_F32,
	f64_t:     lexer.KND_F64,
	any_t:     lexer.KND_ANY,
	uint_t:    lexer.KND_UINT,
	int_t:     lexer.KND_INT,
	uintptr_t: lexer.KND_UINTPTR,
	unsafe_t:  lexer.KND_UNSAFE,
}

func cpp_id(t uint8) string {
	if t == void_t || t == unsafe_t {
		return "void"
	}
	id := type_map[t]
	if id == "" {
		return id
	}
	id = build.AsTypeId(id)
	return id
}
