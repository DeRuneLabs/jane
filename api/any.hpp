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

#ifndef __JNC_ANY_HPP
#define __JNC_ANY_HPP

#include "jn_util.hpp"
#include "ref.hpp"
#include "slice.hpp"

struct any_jnt;

struct any_jnt {
  jn_ref<void *> _data{nil};
  const char *_type_id{nil};

  template <typename T> any_jnt(const T &_Expr) noexcept {
    this->operator=(_Expr);
  }

  any_jnt(const any_jnt &_Src) noexcept { this->operator=(_Src); }

  ~any_jnt(void) noexcept { this->__dealloc(); }

  inline void __dealloc(void) noexcept {
    this->_type_id = nil;
    if (!this->_data._ref) {
      this->_data._alloc = nil;
      return;
    }
    if ((this->_data.__get_ref_n()) != __JNC_REFERENCE_DELTA) {
      return;
    }
    delete this->_data._ref;
    this->_data._ref = nil;
    std::free(*this->_data._alloc);
    this->_data._alloc = nil;
  }

  template <typename T> inline bool __type_is(void) const noexcept {
    if (std::is_same<T, std::nullptr_t>::value) {
      return (false);
    }
    if (this->operator==(nil)) {
      return (false);
    }
    return std::strcmp(this->_type_id, typeid(T).name()) == 0;
  }

  template <typename T> void operator=(const T &_Expr) noexcept {
    this->__dealloc();
    T *_alloc{new (std::nothrow) T};
    if (!_alloc) {
      JNC_ID(panic)(__JNC_ERROR_MEMORY_ALLOCATION_FAILED);
    }
    void **_main_alloc{new (std::nothrow) void *};
    if (!_main_alloc) {
      JNC_ID(panic)(__JNC_ERROR_MEMORY_ALLOCATION_FAILED);
    }
    *_alloc = _Expr;
    *_main_alloc = ((void *)(_alloc));
    this->_data = jn_ref<void *>(_main_alloc);
    this->_type_id = typeid(_Expr).name();
  }

  void operator=(const any_jnt &_Src) noexcept {
    if (_Src.operator==(nil)) {
      this->operator=(nil);
      return;
    }
    this->__dealloc();
    this->_data = _Src._data;
    this->_type_id = _Src._type_id;
  }

  inline void operator=(const std::nullptr_t) noexcept { this->__dealloc(); }

  template <typename T> operator T(void) const noexcept {
    if (this->operator==(nil)) {
      JNC_ID(panic)(__JNC_ERROR_INVALID_MEMORY);
    }
    if (!this->__type_is<T>()) {
      JNC_ID(panic)(__JNC_ERROR_INCOMPATIBLE_TYPE);
    }
    return (*((T *)(*this->_data._alloc)));
  }

  template <typename T> inline bool operator==(const T &_Expr) const noexcept {
    return (this->__type_is<T>() && this->operator T() == _Expr);
  }

  template <typename T>
  inline constexpr bool operator!=(const T &_Expr) const noexcept {
    return (!this->operator==(_Expr));
  }

  inline bool operator==(const any_jnt &_Any) const noexcept {
    if (this->operator==(nil) && _Any.operator==(nil)) {
      return (true);
    }
    return (std::strcmp(this->_type_id, _Any._type_id) == 0);
  }

  inline bool operator!=(const any_jnt &_Any) const noexcept {
    return (!this->operator==(_Any));
  }

  inline bool operator==(std::nullptr_t) const noexcept {
    return (!this->_data._alloc);
  }

  inline bool operator!=(std::nullptr_t) const noexcept {
    return (!this->operator==(nil));
  }

  friend std::ostream &operator<<(std::ostream &_Stream,
                                  const any_jnt &_Src) noexcept {
    if (_Src.operator!=(nil)) {
      _Stream << "<any>";
    } else {
      _Stream << 0;
    }
    return (_Stream);
  }
};

#endif // !__JNC_ANY_HPP
