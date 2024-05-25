// Copyright (c) 2024 arfy slowy - DeRuneLabs
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

#ifndef __JANE_TYPES_HPP
#define __JANE_TYPES_HPP

#include "platform.hpp"
#include <stddef.h>

namespace jane {
#ifdef ARCH_32BIT
typedef unsigned long long int Uint;
typedef signed long int Int;
typedef unsigned long int Uintptr;
#else
typedef unsigned long long int Uint;
typedef signed long long int Int;
typedef unsigned long long int Uintptr;
#endif

typedef signed char I8;
typedef signed short int I16;
typedef signed long int I32;
typedef signed long long int I64;
typedef unsigned char U8;
typedef unsigned short int U16;
typedef unsigned long int U32;
typedef unsigned long long int U64;
typedef float F32;
typedef double F64;
typedef bool Bool;

constexpr decltype(nullptr) nil{nullptr};

constexpr jane::F32 MAX_F32{0x1p127 * (1 + (1 - 0x1p-23))};
constexpr jane::F32 MIN_F32{-0x1p127 * (1 + (1 - 0x1p-23))};
constexpr jane::F64 MAX_F64{0x1p1023 * (1 + (1 - 0x1p-52))};
constexpr jane::F64 MIN_F64{-0x1p1023 * (1 + (1 - 0x1p-52))};
constexpr jane::I64 MAX_I64{9223372036854775807LL};
constexpr jane::I64 MIN_I64{-9223372036854775807 - 1};
constexpr jane::U64 MAX_U64{18446744073709551615LLU};
} // namespace jane

#endif // __JANE_TYPES_HPP