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

#ifndef __JANE_STD_SYSCALL_UNIX_HPP
#define __JANE_STD_SYSCALL_UNIX_HPP

#include <limits.h>
#include <unistd.h>
#include <cstddef>

#include "../../api/str.hpp"
#include "../../api/typedef.hpp"

str_jnt __jane_str_from_byte_ptr(const char *_Ptr) noexcept;
str_jnt __jane_str_from_byte_ptr(const JANE_ID(byte) *_Ptr) noexcept;
str_jnt __jane_stat(const char *_Path, struct stat *_Stat) noexcept;

str_jnt __jane_from_byte_ptr(const char *_Ptr) noexcept {
  return (__jane_str_from_byte_ptr((const JANE_ID(byte)*)(_Ptr)));
}

str_jnt __jane_str_from_byte_ptr(const JANE_ID(byte) *_Ptr) noexcept {
  return (str_jnt(_Ptr));
}

int_jnt __jane_stat(const char *_Path, struct stat *_Stat) noexcept {
  return (stat(_Path, _Stat));
}

#endif // !__JANE_STD_SYSCALL_UNIX_HPP
