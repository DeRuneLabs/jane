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

#ifndef __JNC_TRAIT_HPP
#define __JNC_TRAIT_HPP

#include "jn_util.hpp"
#include "typedef.hpp"

template <typename T> struct trait;

template <typename T> struct trait {
public:
  T *_data{nil};
  mutable uint_jnt *_ref{nil};

  trait<T>(void) noexcept {}
  trait<T>(std::nullptr_t) noexcept {}

  template <typename TT> trait<T>(const TT &_Data) noexcept {
    TT *_alloc = new (std::nothrow) TT{_Data};
    if (!_alloc) {
      JNC_ID(panic)(__JNC_ERROR_MEMORY_ALLOCATION_FAILED);
    }
    this->_data = (T *)(_alloc);
    this->_ref = new (std::nothrow) uint_jnt{1};
    if (!this->_ref) {
      JNC_ID(panic)(__JNC_ERROR_MEMORY_ALLOCATION_FAILED);
    }
  }

  trait<T>(const trait<T> &_Src) noexcept { this->operator=(_Src); }

  void __dealloc(void) noexcept {
    if (!this->_ref) {
      return;
    }
    (*this->_ref)--;
    if (*this->_ref != 0) {
      return;
    }
    delete this->_ref;
    this->_ref = nil;
    delete this->_data;
    this->_data = nil;
  }

  T &get(void) noexcept {
    if (this->operator==(nil)) {
      JNC_ID(panic)(__JNC_ERROR_INVALID_MEMORY);
    }
    return *this->_data;
  }

  ~trait(void) noexcept { this->__dealloc(); }

  inline void operator=(const std::nullptr_t) noexcept { this->__dealloc(); }

  inline void operator=(const trait<T> &_Src) noexcept {
    this->__dealloc();
    if (_Src == nil) {
      return;
    }
    (*_Src._ref)++;
    this->_data = _Src._data;
    this->_ref = _Src._ref;
  }

  inline bool operator==(const trait<T> &_Src) const noexcept {
    return this->_data == this->_data;
  }

  inline bool operator!=(const trait<T> &_Src) const noexcept {
    return !this->operator==(_Src);
  }

  inline bool operator==(std::nullptr_t) const noexcept { return !this->_data; }

  inline bool operator!=(std::nullptr_t) const noexcept {
    return !this->operator==(nil);
  }

  friend inline std::ostream &operator<<(std::ostream &_Stream,
                                         const trait<T> &_Src) noexcept {
    return _Stream << _Src._data;
  }
};

#endif // !__JNC_TRAIT_HPP
