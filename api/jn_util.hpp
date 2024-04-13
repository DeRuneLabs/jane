#ifndef __JNC_UTIL_LIBS_HPP
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

#define __JNC_UTIL_LIBS_HPP

#if defined(WIN32) || defined(_WIN32) || defined(__WIN32__) || defined(__NT__)
#ifndef _WINDOWS
#define _WINDOWS
#endif // !_WINDOWS
#endif // defined(WIN32) || defined(_WIN32) || defined(__WIN32__) ||
       // defined(__NT__)

#include <cstddef>
#include <cstdint>
#include <functional>
#include <ostream>
#include <sstream>
#include <type_traits>

#include <any>
#include <cstring>
#include <functional>
#include <iostream>
#include <map>
#include <sstream>
#include <string>
#include <thread>
#include <typeinfo>
#include <valarray>
#include <vector>
#ifdef _WINDOWS
#include <codecvt>
#include <fcntl.h>
#include <windows.h>
#endif // _WINDOWS

// PTR HEAP define
#define __JNC_PTR_NEVER_HEAP ((bool**)(0x0000001))
#define __JNC_PTR_HEAP_TRUE ((bool*)(0x0000001))

// jn ptr
#define __jnc_ptr(_PTR) (_PTR)

// error message cpp api
#define __JNC_ERROR_INVALID_MEMORY                                             \
  ("invalid memory address or nil pointer deference")
#define __JNC_ERROR_INCOMPATIBLE_TYPE ("incompatible type")
#define __JNC_ERROR_MEMORY_ALLOCATION_FAILED ("memory allocation failed")
#define __JNC_ERROR_INDEX_OUT_OF_RANGE ("index of out range")
#define __JNC_WRITE_ERROR_SLICING_INDEX_OUT_OF_RANGE(_STREAM, _START, _LEN)    \
  (_STREAM << __JNC_ERROR_INDEX_OUT_OF_RANGE << '[' << _START << ':' << _LEN   \
           << ']')

#define __JNC_WRITE_ERROR_INDEX_OUT_OF_RANGE(_STREAM, _INDEX)                  \
  (_STREAM << __JNC_ERROR_INDEX_OUT_OF_RANGE << '[' << _INDEX << ']')
#define __JNC_EXIT_PANIC (2)
#define __JNC_CCONCAT(_A, _B) _A##_B
#define __JNC_CONCAT(_A, _B) __JNC_CCONCAT(_A, _B)
#define __JNC_IDENTIFIER_PREFIX _
#define JNC_ID(_IDENTIFIER) __JNC_CONCAT(__JNC_IDENTIFIER_PREFIX, _IDENTIFIER)

#define nil (nullptr)
#define CO(_Expr) (std::thread{[&](void) mutable -> void { _EXPR; }}).detach()

template <typename _Obj_t> void JNC_ID(panic)(const _Obj_t &_Expr);

// print standard cpp
#define _print(_EXPR) (std::cout << _EXPR)
// println standard cpp
template <typename _Obj_t>
inline void JNC_ID(println)(const _Obj_t _Obj) noexcept {
  JNC_ID(print)(_Obj);
  std::cout << std::endl;
}

#endif // !__JNC_UTIL_LIBS_HPP
