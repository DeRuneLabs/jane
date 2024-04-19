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

#ifndef __JNC_UTIL_HPP
#define __JNC_UTIL_HPP

#if defined(WIN32) || defined(_WIN32) || defined(__WIN32__) || defined(__NT__)
#define _WINDOWS
#elif defined(__linux__) || defined(linux) || defined(__linux)
#define _LINUX
#elif defined(__APPLE__) || defined(__MACH__)
#define _DARWIN
#endif

#if defined(_LINUX) || defined(_DARWIN)
#define _UNIX
#endif // defined(_LINUX) || defined(_DARWIN)

#if defined(__amd64) || defined(__amd64__) || defined(__x86_64) ||             \
    defined(__x86_64__) || defined(_M_AMD64)
#define _AMD64
#elif defined(__aarch64__)
`#define _ARM64
#elif defined(i386) || defined(__thumb__) || defined(_M_ARM) || defined(__arm)
#define _ARM
#elif defined(__aarch64__)
#define _ARM64
#elif defined(i386) || defined(__i386) || defined(__i386__) ||                 \
    defined(_X86_) || defined(__I86__) || defined(__386)
#define _I386
#endif

#if defined(_AMD64) || defined(_ARM64)
#define _64BIT
#else
#define _32BIT
#endif

#include <cstddef>
#include <cstring>
#include <functional>
#include <iostream>
#include <iterator>
#include <sstream>
#include <string>
#include <thread>
#include <typeinfo>
#include <unordered_map>
#ifdef _WINDOWS
#include <codecvt>
#include <fcntl.h>
#include <windows.h>
#endif // _WINDOWS

constexpr const char *__JNC_ERROR_INVALID_MEMORY{
    "invalid memory address or nil pointer deference"};
constexpr const char *__JNC_ERROR_INCOMPATIBLE_TYPE{"incompatible type"};
constexpr const char *__JNC_ERROR_MEMORY_ALLOCATION_FAILED{
    "memory allocation failed"};
constexpr const char *__JNC_ERROR_INDEX_OUT_OF_RANGE{"index out of range"};
constexpr signed int __JNC_EXIT_PANIC{2};
constexpr std::nullptr_t nil{nullptr};

// writing error slicing
#define __JNC_WRITE_ERROR_SLICING_INDEX_OUT_OF_RANGE(_STREAM, _START, _LEN)    \
  (_STREAM << __JNC_ERROR_INDEX_OUT_OF_RANGE << '[' << _START << ':' << _LEN   \
           << ']')
#define __JNC_WRITE_ERROR_INDEX_OUT_OF_RANGE(_STREAM, _INDEX)                  \
  (_STREAM << __JNC_ERROR_INDEX_OUT_OF_RANGE << '[' << _INDEX << ']')

#define __JNC_CCONCAT(_A, _B) _A##_B
#define __JNC_CONCAT(_A, _B) __JNC_CCONCAT(_A, _B)

// identifier jane C
#define __JNC_IDENTIFIER_PREFIX _
#define JNC_ID(_IDENTIFIER) __JNC_CONCAT(__JNC_IDENTIFIER_PREFIX, _IDENTIFIER)
#define __JNC_CO(_EXPR)                                                        \
  (std::thread{[&](void) mutable -> void { _EXPR; }}.detach())
#define __JNC_DEFER(_EXPR)                                                     \
  defer __JNC_CONCAT(JNC_DEFER_, __LINE__) {                                   \
    [&](void) -> void { _EXPR; }                                               \
  }

constexpr signed int __JNC_REFERENCE_DELTA{1};

// builtin str type
class str_jnt;

template <typename _Obj_t> void JNC_ID(panic)(const _Obj_t &_Expr);
inline std::ostream &operator<<(std::ostream &_Stream,
                                const signed char _I8) noexcept;
inline std::ostream &operator<<(std::ostream &_Stream,
                                const unsigned char _U8) noexcept;

#endif // !JNC_UTIL_HPP
