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

#ifndef __JANE_SLICE_HPP
#define __JANE_SLICE_HPP

#include <cstddef>
#include <initializer_list>
#include <ostream>
#include <sstream>

#include "error.hpp"
#include "panic.hpp"
#include "ref.hpp"
#include "types.hpp"

namespace jane {
template <typename Item> class Slice;

template <typename Item> class Slice {
public:
  jane::Ref<Item> data{};
  Item *_slice{nullptr};
  jane::Uint _len{0};
  jane::Uint _cap{0};

  static jane::Slice<Item> alloc(const jane::Uint &n) noexcept {
    jane::Slice<Item> buffer;
    buffer.alloc_new(n < 0 ? 0 : n);
    return buffer;
  }

  Slice<Item>(void) noexcept {}
  Slice<Item>(const std::nullptr_t) noexcept {}

  Slice<Item>(const jane::Slice<Item> &src) noexcept { this->operator=(src); }

  Slice<Item>(const std::initializer_list<Item> &src) noexcept {
    if (src.size() == 0) {
      return;
    }

    this->alloc_new(src.size());
    const auto src_begin{src.begin()};
    for (jane::Int i{0}; i < this->_len; ++i) {
      this->data.alloc[i] = *reinterpret_cast<const Item *>(src_begin + i);
    }
  }

  ~Slice<Item>(void) noexcept { this->dealloc(); }

  inline void check(void) const noexcept {
    if (this->operator==(nullptr)) {
      jane::panic(jane::ERROR_INVALID_MEMORY);
    }
  }

  void dealloc(void) noexcept {
    this->_len = 0;
    this->_cap = 0;

    if (!this->data.ref) {
      this->data.alloc = nullptr;
      return;
    }

    if (this->data.get_ref_n() != jane::REFERENCE_DELTA) {
      this->data.alloc = nullptr;
      return;
    }

    delete this->data.ref;
    this->data.ref = nullptr;

    delete[] this->data.alloc;
    this->data.alloc = nullptr;
    this->data.ref = nullptr;
    this->_slice = nullptr;
  }

  void alloc_new(const jane::Int n) noexcept {
    this->dealloc();

    Item *alloc{n == 0 ? new (std::nothrow) Item[0]
                       : new (std::nothrow) Item[n]{Item()}};
    if (!alloc) {
      jane::panic(jane::ERROR_MEMORY_ALLOCATION_FAILED);
    }
    this->data = jane::Ref<Item>::make(alloc);
    this->_len = n;
    this->_cap = n;
    this->_slice = &alloc[0];
  }

  typedef Item *Iterator;
  typedef const Item *ConstIterator;

  inline constexpr Iterator begin(void) noexcept { return &this->_slice[0]; }

  inline constexpr ConstIterator begin(void) const noexcept {
    return &this->_slice[0];
  }

  inline constexpr Iterator end(void) noexcept {
    return &this->_slice[this->_len];
  }

  inline constexpr ConstIterator end(void) const noexcept {
    return &this->_slice[this->_len];
  }

  inline Slice<Item> slice(const jane::Int &start,
                           const jane::Int &end) const noexcept {
    this->check();

    if (start < 0 || end < 0 || start > end || end > this->cap()) {
      std::stringstream sstream;
      __JANE_WRITE_ERROR_SLICING_INDEX_OUT_OF_RANGE(sstream, start, end);
      jane::panic(sstream.str().c_str());
    }

    jane::Slice<Item> slice;
    slice.data = this->data;
    slice._slice = &this->_slice[start];
    slice._len = end - start;
    slice._cap = this->_cap - start;
    return slice;
  }

  inline jane::Slice<Item> slice(const jane::Int &start) const noexcept {
    return this->slice(start, this->len());
  }

  inline jane::Slice<Item> slice(void) const noexcept {
    return this->slice(0, this->len());
  }

  inline constexpr jane::Int len(void) const noexcept { return this->_len; }

  inline constexpr jane::Int cap(void) const noexcept { return this->_cap; }

  inline jane::Bool empty(void) const noexcept {
    return !this->_slice || this->_len == 0 || this->_cap == 0;
  }

  void push(const Item &item) noexcept {
    if (this->_len == this->_cap) {
      Item *_new{new (std::nothrow) Item[this->_len + 1]};
      if (!_new) {
        jane::panic(jane::ERROR_MEMORY_ALLOCATION_FAILED);
      }
      for (jane::Int index{0}; index < this->_len; ++index) {
        _new[index] = this->data.alloc[index];
      }
      _new[this->_len] = item;

      delete[] this->data.alloc;
      this->data.alloc = nullptr;
      this->data.alloc = _new;
      this->_slice = this->data.alloc;
    } else {
      this->_slice[this->_len] = item;
    }
    ++this->_len;
  }

  jane::Bool operator==(const jane::Slice<Item> &src) const noexcept {
    if (this->_len != src._len) {
      return false;
    }
    for (jane::Int index{0}; index < this->_len; ++index) {
      if (this->_slice[index] != src._slice[index]) {
        return false;
      }
    }
    return true;
  }

  inline constexpr jane::Bool
  operator!=(const jane::Slice<Item> &src) const noexcept {
    return !this->operator==(src);
  }

  inline constexpr jane::Bool operator==(const std::nullptr_t) const noexcept {
    return !this->_slice;
  }

  inline constexpr jane::Bool operator!=(const std::nullptr_t) const noexcept {
    return !this->operator==(nullptr);
  }

  Item &operator[](const jane::Int &index) const {
    this->check();
    if (this->empty() || index < 0 || this->len() <= index) {
      std::stringstream sstream;
      __JANE_WRITE_ERROR_INDEX_OUT_OF_RANGE(sstream, index);
      jane::panic(sstream.str().c_str());
    }
    return this->_slice[index];
  }

  void operator=(const jane::Slice<Item> &src) noexcept {
    if (this->data.alloc == src.data.alloc) {
      this->_len = src._len;
      this->_cap = src._cap;
      this->data = src.data;
      this->_slice = src._slice;
      return;
    }
    this->dealloc();
    if (src.operator==(nullptr)) {
      return;
    }
    this->_len = src._len;
    this->_cap = src._cap;
    this->data = src.data;
    this->_slice = src._slice;
  }

  void operator=(const std::nullptr_t) noexcept { this->dealloc(); }

  friend std::ostream &operator<<(std::ostream &stream,
                                  const jane::Slice<Item> &src) noexcept {
    if (src.empty()) {
      return stream << "[]";
    }
    stream << '[';
    for (jane::Int index{0}; index < src._len;) {
      stream << src._slice[index++];
      if (index < src._len) {
        stream << ' ';
      }
    }
    stream << ']';
    return stream;
  }
};
} // namespace jane

#endif // __JANE_SLICE_HPP