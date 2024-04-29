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

#ifndef __JANE_TRAIT_HPP
#define __JANE_TRAIT_HPP

#include "ref.hpp"
#include "typedef.hpp"

template <typename T> struct trait_jnt;

template <typename T> struct trait_jnt {
public:
  ref_jnt<T> __data{};
  const char *__type_id{nil};

  trait_jnt<T>(void) noexcept {}
  trait_jnt<T>(std::nullptr_t) noexcept {}

  template <typename TT> trait_jnt<T>(const TT &_Data) noexcept {
    TT *_alloc{new (std::nothrow) TT};
    if (!_alloc) {
      JANE_ID(panic)(__JANE_ERROR_MEMORY_ALLOCATION_FAILED);
    }
    *_alloc = _Data;
    this->__data = ref_jnt<T>::make((T *)(_alloc));
    this->__type_id = typeid(_Data).name();
  }

  template <typename TT> trait_jnt<T>(const ref_jnt<TT> &_Ref) noexcept {
    this->__data = ref_jnt<T>::make(((T *)(_Ref.__alloc)), _Ref.__ref);
    this->__data.__add_ref();
    this->__type_id = typeid(_Ref).name();
  }

  trait_jnt<T>(const trait_jnt<T> &_Src) noexcept { this->operator=(_Src); }

  void __dealloc(void) noexcept { this->__data._drop(); }

  inline void __must_ok(void) noexcept {
    if (this->operator==(nil)) {
      JANE_ID(panic)(__JANE_ERROR_INVALID_MEMORY);
    }
  }

  inline T &_get(void) noexcept {
    this->__must_ok();
    return this->__data;
  }

  ~trait_jnt(void) noexcept {}

  template <typename TT> operator TT(void) noexcept {
    this->__must_ok();
    if (std::strcmp(this->__type_id, typeid(TT).name()) != 0) {
      JANE_ID(panic)(__JANE_ERROR_INCOMPATIBLE_TYPE);
    }
    this->__data.__add_ref();
    return (ref_jnt<TT>((TT *)(this->__data.__alloc), this->__data.__ref));
  }

  inline void operator==(const std::nullptr_t) noexcept { this->__dealloc(); }

  inline void operator=(const trait_jnt<T> &_Src) noexcept {
    this->__dealloc();
    if (_Src == nil) {
      return;
    }
    this->__data = _Src.__data;
    this->__type_id = _Src.__type_id;
  }

  inline bool operator==(const trait_jnt<T> &_Src) const noexcept {
    return (this->__data.__alloc == this->__data.__alloc);
  }

  inline bool operator!=(const trait_jnt<T> &_Src) const noexcept {
    return (!this->operator==(_Src));
  }

  inline bool operator==(std::nullptr_t) const noexcept {
    return (this->__data.__alloc == nil);
  }

  inline bool operator!=(std::nullptr_t) const noexcept {
    return (!this->operator==(nil));
  }

  friend inline std::ostream &operator<<(std::ostream &_Stream,
                                         const trait_jnt<T> &_Src) noexcept {
    return (_Stream << _Src.__data.__alloc);
  }
};

#endif // !__JANE_TRAIT_HPP
