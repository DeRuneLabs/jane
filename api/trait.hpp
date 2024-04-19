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
#include "ref.hpp"
#include <iterator>

template <typename T> struct trait;

template <typename T> struct trait {
public:
  jn_ref<T> _data{nil};
  const char *type_id{nil};

  trait<T>(void) noexcept {}
  trait<T>(std::nullptr_t) noexcept {}

  template <typename TT> trait<T>(const TT &_Data) noexcept {
    TT *_alloc{new (std::nothrow) TT};
    if (!_alloc) {
      JNC_ID(panic)(__JNC_ERROR_MEMORY_ALLOCATION_FAILED);
    }
    *_alloc = _Data;
    this->_data = jn_ref<T>((T *)(_alloc));
    this->type_id = typeid(_Data).name();
  }

  template <typename TT> trait<T>(const jn_ref<TT> &_Ref) noexcept {
    this->_data = jn_ref<T>(((T *)(_Ref._alloc)), _Ref._ref);
    this->_data.__add_ref();
    this->type_id = typeid(_Ref).name();
  }

  trait<T>(const trait<T> &_Src) noexcept { this->operator=(_Src); }

  void __dealloc(void) noexcept { this->_data.__drop(); }

  inline void __must_ok(void) noexcept {
    if (this->operator==(nil)) {
      JNC_ID(panic)(__JNC_ERROR_INVALID_MEMORY);
    }
  }

  inline T &get(void) noexcept {
    this->__must_ok();
    return this->_data;
  }

  ~trait(void) noexcept { this->__dealloc(); }

  template <typename TT> operator TT(void) noexcept {
    this->__must_ok();
    if (std::strcmp(this->type_id, typeid(TT).name()) != 0) {
      JNC_ID(panic)(__JNC_ERROR_INCOMPATIBLE_TYPE);
    }
    return (*((TT *)(this->_data._alloc)));
  }

  template <typename TT> operator jn_ref<TT>(void) noexcept {
    this->__must_ok();
    if (std::strcmp(this->type_id, typeid(jn_ref<TT>).name()) != 0) {
      JNC_ID(panic)(__JNC_ERROR_INCOMPATIBLE_TYPE);
    }
    this->_data.__add_ref();
    return (jn_ref<TT>((TT *)(this->_data._alloc), this->_data._ref));
  }

  inline void operator=(const std::nullptr_t) noexcept { this->__dealloc(); }

  inline void operator=(const trait<T> &_Src) noexcept {
    this->__dealloc();
    if (_Src == nil) {
      return;
    }
    this->_data = _Src._data;
  }

  inline bool operator==(const trait<T> &_Src) const noexcept {
    return (this->_data._alloc == this->_data._alloc);
  }

  inline bool operator!=(const trait<T> &_Src) const noexcept {
    return (!this->operator==(_Src));
  }

  inline bool operator==(std::nullptr_t) const noexcept {
    return (this->_data._alloc == nil);
  }

  inline bool operator!=(std::nullptr_t) const noexcept {
    return (!this->operator==(nil));
  }

  friend inline std::ostream &operator<<(std::ostream &_Stream,
                                         const trait<T> &_Src) noexcept {
    return (_Stream << _Src._data._alloc);
  }
};

#endif // !__JNC_TRAIT_HPP
