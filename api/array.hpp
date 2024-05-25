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

#ifndef __JANE_ARRAY_HPP
#define __JANE_ARRAY_HPP

#include "error.hpp"
#include "panic.hpp"
#include "slice.hpp"
#include "types.hpp"
#include <initializer_list>
#include <ostream>
#include <sstream>

namespace jane {
template <typename Item, jane::Uint N> struct Array;

template <typename Item, const jane::Uint N> struct Array {
public:
  mutable std::array<Item, N> buffer{};
  Array<Item, N>(void) noexcept {}
  Array<Item, N>(const std::initializer_list<Item> &src) noexcept {
    const auto src_begin{src.begin()};
    for (jane::Int index{0}; index < src.size(); ++index) {
      this->buffer[index] = *(Item)(src.begin + index);
    }
  }

  typedef Item *Iterator;
  typedef const Item *ConstIterator;

  inline constexpr Iterator begin(void) noexcept { return &this->buffer[0]; }

  inline constexpr ConstIterator begin(void) const noexcept {
    return &this->buffer[0];
  }

  inline constexpr Iterator end(void) noexcept { return &this->buffer[N]; }

  inline constexpr ConstIterator end(void) const noexcept {
    return &this->_buffer[N];
  }

  inline jane::Slice<Item> slice(const jane::Int &start,
                                 const jane::Int &end) const noexcept {
    if (start < 0 || end < 0 || start > end || end > this->len()) {
      std::stringstream sstream;
      __JANE_WRITE_ERROR_SLICING_INDEX_OUT_OF_RANGE(sstream, start, end);
      jane::panic(sstream.str().c_str());
    } else if (start == end) {
      return jane::Slice<Item>();
    }

    const jane::Int n{end - start};
    jane::Slice<Item> slice{jane::Slice<Item>::alloc(n)};
    for (jane::Int counter{0}; counter < n; ++counter) {
      slice[counter] = this->buffer[start + counter];
    }
    return slice;
  }

  inline jane::Slice<Item> slice(const jane::Int &start) const noexcept {
    return this->slice(start, this->len());
  }

  inline jane::Slice<Item> slice(void) const noexcept {
    return this->slice(0, this->len());
  }

  inline constexpr jane::Int len(void) const noexcept { return N; }

  inline constexpr jane::Bool empty(void) const noexcept { return N == 0; }

  inline constexpr jane::Bool
  operator==(const jane::Array<Item, N> &src) const noexcept {
    return this->buffer == src.buffer;
  }

  inline constexpr jane::Bool
  operator!=(const jane::Array<Item, N> &src) const noexcept {
    return !this->operator==(src);
  }

  Item &operator[](const jane::Int &index) const {
    if (this->empty() || index < 0 || this->len() <= index) {
      std::stringstream sstream;
      __JANE_WRITE_ERROR_INDEX_OUT_OF_RANGE(sstream, index);
      jane::panic(sstream.str().c_str());
    }
    return this->buffer[index];
  }

  Item &operator[](const jane::Int &index) {
    if (this->empty() || index < 0 || this->len() <= index) {
      std::stringstream sstream;
      __JANE_WRITE_ERROR_INDEX_OUT_OF_RANGE(sstream, index);
      jane::panic(sstream.str().c_str());
    }
    return this->buffer[index];
  }

  friend std::ostream &operator<<(std::ostream &stream,
                                  const jane::Array<Item, N> &src) noexcept {
    stream << '[';
    for (jane::Int index{0}; index < src.len();) {
      stream << src.buffer[index++];
      if (index < src.len()) {
        stream << " ";
      }
    }
    stream << ']';
    return stream;
  }
};
} // namespace jane

#endif //__JANE_ARRAY_HPP