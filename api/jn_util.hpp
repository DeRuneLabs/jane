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

#define JN_EXIT_PANIC 2
#define _CONCAT(_A, _B) _A##_B
#define CONCAT(_A, _B) _CONCAT(_A, _B)
#define JNID(_Identifier) CONCAT(_, _Identifier)
#define nil nullptr
#define CO(_Expr) std::thread{[&](void) mutable -> void { _Expr; }}.detach()

void JNID(panic)(const char *_Message);

#endif // !__JNC_UTIL_LIBS_HPP
