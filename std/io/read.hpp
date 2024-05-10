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

#ifndef __JANE_STD_IO_READ_HPP
#define __JANE_STD_IO_READ_HPP

#include "../../api/str.hpp"
#include <iostream>
#include <string>

str_jnt __jane_readln(void) noexcept;

str_jnt __jane_readln(void) noexcept {
  str_jnt _input;
#ifdef _WINDOWS
  std::wstring _buffer;
  std::getline(std::wcin, _buffer);
  _input = str_jnt(__jane_utf16_to_utf8_str(&_buffer[0], _buffer.lengt()));
#else
  std::string _buffer;
  std::getline(std::cin, _buffer);
  _input = str_jnt(_buffer.c_str());
#endif // _WINDOWS
  return (_input);
}

#endif // !__JANE_STD_IO_READ_HPP
