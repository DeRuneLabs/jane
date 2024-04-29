// Copyright (c) 2024 - DeRuneLabs
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

#ifndef __JANE_TYPEDEF_HPP
#define __JANE_TYPEDEF_HPP

#include <cstring>
#include <new>
#include <ostream>
#include <stdlib.h>
#include <sstream>
#include <functional>
#include <iostream>

#if defined(_32BIT)
typedef unsigned long int(uint_jnt);
typedef signed long int(int_jnt);
typedef unsigned long int(uintptr_jnt);
#else
typedef unsigned long long int(uint_jnt);
typedef signed long long int(int_jnt);
typedef unsigned long long int(uintptr_jnt);
#endif // defined(_32BIT)

typedef signed char(i8_jnt);
typedef signed short int(i16_jnt);
typedef signed long int(i32_jnt);
typedef signed long long int(i64_jnt);
typedef unsigned char(u8_jnt);
typedef unsigned short int(u16_jnt);
typedef unsigned long int(u32_jnt);
typedef unsigned long long int(u64_jnt);
typedef float(f32_jnt);
typedef double(f64_jnt);
typedef bool(bool_jnt);

constexpr const char *__JANE_ERROR_INVALID_MEMORY{
    "invalid memory address or nil pointer deference"};
constexpr const char *__JANE_ERROR_INCOMPATIBLE_TYPE{"incompatible type"};
constexpr const char *__JANE_ERROR_MEMORY_ALLOCATION_FAILED{
    "memory allocation failed"};
constexpr const char *__JANE_ERROR_INDEX_OUT_OF_RANGE{"index out of range"};
constexpr const char *__JANE_ERROR_DIVIDE_BY_ZERO{"divide by zero"};
constexpr signed int __JANE_EXIT_PANIC{2};
constexpr std::nullptr_t nil{nullptr};

#define __JANE_WRITE_ERROR_SLICING_INDEX_OUT_OF_RANGE(_STREAM, _START, _LEN)   \
  (_STREAM << __JANE_ERROR_INDEX_OUT_OF_RANGE << '[' << _START << ':' << _LEN  \
           << ']')

#define __JANE_WRITE_ERROR_INDEX_OUT_OF_RANGE(_STREAM, _INDEX)                 \
  (_STREAM << __JANE_ERROR_INDEX_OUT_OF_RANGE << '[' << _INDEX << ']')

#define __JANE_CCONCAT(_A, _B) _A##_B
#define __JANE_CONCAT(_A, _B) __JANE_CCONCAT(_A, _B)
#define __JANE_IDENTIFIER_PREFIX _
#define JANE_ID(_IDENTIFIER)                                                   \
  __JANE_CONCAT(__JANE_IDENTIFIER_PREFIX, _IDENTIFIER)
#define __JANE_CO(_EXPR)                                                       \
  (std::thread{[&](void) mutable -> void { _EXPR; }}.detach())

template <typename _Obj_t> void JANE_ID(panic)(const _Obj_t &_Expr);
inline std::ostream &operator<<(std::ostream &_Stream,
                                const signed char _I8) noexcept;
inline std::ostream &operator<<(std::ostream &_Stream,
                                const unsigned char _U8) noexcept;

#endif // !__JANE_TYPEDEF_HPP
