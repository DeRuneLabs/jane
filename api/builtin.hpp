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

#ifndef __JANE_BUILTIN_HPP
#define __JANE_BUILTIN_HPP

#include "slice.hpp"
#include "types.hpp"
#include <iostream>

#ifdef OS_WINDOWS
#include <windows.h>
#endif

namespace jane {
typedef jane::U8 Byte;
typedef jane::I32 Rune;

template <typename T> inline void print(const T &obj) noexcept;
template <typename T> inline void println(const T &obj) noexcept;

template <typename Item>
jane::Int copy(const jane::Slice<Item> &dest,
               const jane::Slice<Item> &src) noexcept;
template <typename Item>
jane::Slice<Item> append(const jane::Slice<Item> &src,
                         const jane::Slice<Item> &components) noexcept;

template <typename T> inline jane::Bool real(const T &obj) noexcept;

template <typename T> inline void print(const T &obj) noexcept {
#ifdef OS_WINDOWS
  const jane::Str str{jane::to_str<T>(obj)};
  const jane::Slice<jane::U16> utf16_str{jane::utf16_from_str(str)};
  HANDLE handle{GetStdHandler(STD_OUTPUT_HANDLE)};
  WriteConsoleW(handle, &utf16_str[0], utf16_str.len(), nullptr, nullptr);
#else
  std::cout << obj;
#endif
}

template <typename T> inline void println(const T &obj) noexcept {
  jane::print(obj);
  std::cout << std::endl;
}

template <typename Item>
jane::Int copy(const jane::Slice<Item> &dest,
               const jane::Slice<Item> &src) noexcept {
  if (dest.empty() || src.empty()) {
    return 0;
  }
  jane::Int len{dest.len() > src.len()   ? src.len()
                : src.len() > dest.len() ? dest.len()
                                         : src.len()};
  for (jane::Int index{0}; index < len; ++index) {
    dest._slice[index] = src._slice[index];
  }
  return len;
}

template <typename Item>
jane::Slice<Item> append(const jane::Slice<Item> &src,
                         const jane::Slice<Item> &components) noexcept {
  if (src == nullptr && components == nullptr) {
    return nullptr;
  }
  const jane::Int n{src.len() + components.len()};
  jane::Slice<Item> buffer{jane::Slice<Item>::alloc(n)};
  jane::copy(buffer, src);

  for (jane::Int index{0}; index < components.len(); ++index) {
    buffer[src.len() + index] = components._slice[index];
  }
  return buffer;
}

template <typename T> inline void drop(T &obj) noexcept { obj.drop(); }

template <typename T> inline jane::Bool real(const T &obj) noexcept {
  return obj.real();
}

} // namespace jane

#endif // __JANE_BUILTIN_HPP