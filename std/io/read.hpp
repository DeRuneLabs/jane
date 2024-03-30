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

#ifndef __JNC_STD_IO_READ_HPP
#define __JNC_STD_IO_READ_HPP

#include "../../api/str.hpp"

str_jnt __jnc_read() noexcept;
str_jnt __jnc_readln() noexcept;

#ifdef _WINDOWS

inline std::string
__jnc_std_io_encode_utf8(const std::wstring &_WStr) noexcept {
  std::wstring_convert<std::codecvt_utf8<wchar_t>, wchar_t> conv{};
  return conv.to_bytes(_WStr);
}
#endif // _WINDOWS

str_jnt __jnc_read() noexcept {
#ifdef _WINDOWS
  std::wstring buffer{};
  std::wcin >> buffer;
  return __jnc_std_io_encode_utf8(buffer).c_str();
#else
  std::string buffer{};
  std::cin >> buffer;
  return buffer.c_str();
#endif // _WINDOWS
}

str_jnt __jnc_readln() noexcept {
#ifdef _WINDOWS
  std::wstring buffer{};
  std::getline(std::wcin, buffer);
  return __jnc_std_io_encode_utf8(buffer).c_str();
#else
  std::string buffer{};
  std::getline(std::cin, buffer);
  return buffer.c_str();
#endif // _WINDOWS
}

#endif // !__JNC_STD_IO_READ_HPP
