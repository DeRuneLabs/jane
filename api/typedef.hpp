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

#ifndef __JNC_TYPEDEF_HPP
#define __JNC_TYPEDEF_HPP

#include "jn_util.hpp"

typedef std::size_t(uint_jnt);
typedef std::make_signed<uint_jnt>::type(int_jnt);
typedef signed char(i8_jnt);
typedef signed short(i16_jnt);
typedef signed long(i32_jnt);
typedef signed long long(i64_jnt);
typedef unsigned char(u8_jnt);
typedef unsigned short(u16_jnt);
typedef unsigned long(u32_jnt);
typedef unsigned long long(u64_jnt);
typedef float(f32_jnt);
typedef double(f64_jnt);
typedef bool(bool_jnt);
typedef std::uintptr_t(uintptr_jnt);

#endif // !__JNC_TYPEDEF_HPP