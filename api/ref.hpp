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

#ifndef __JANE_REF_HPP
#define __JANE_REF_HPP

#include "atomicity.hpp"
#include "typedef.hpp"
constexpr signed int __JANE_REFERENCE_DELTA{1};

template <typename T> struct ref_jnt;

template <typename T> struct ref_jnt {
  mutable T *__alloc{nil};
  mutable uint_jnt *__ref{nil};

  static ref_jnt<T> make(T *_Ptr, uint_jnt *_Ref) noexcept {
    ref_jnt<T> _buffer;
    _buffer.__alloc = _Ptr;
    _buffer.__ref = _Ref;
    return (_buffer);
  }

  static ref_jnt<T> make(T *_Ptr) noexcept {
    ref_jnt<T> _buffer;
    _buffer.__ref = (new (std::nothrow) uint_jnt);
    if (!_buffer.__ref) {
      JANE_ID(panic)(__JANE_ERROR_MEMORY_ALLOCATION_FAILED);
    }
    *_buffer.__ref = 1;
    _buffer.__alloc = _Ptr;
    return (_buffer);
  }

  static ref_jnt<T> make(const T &_Instance) noexcept {
    ref_jnt<T> _buffer;
    _buffer.__alloc = (new (std::nothrow) T);
    if (!_buffer.__alloc) {
      JANE_ID(panic)(__JANE_ERROR_MEMORY_ALLOCATION_FAILED);
    }
    _buffer.__ref = (new (std::nothrow) uint_jnt);
    if (!_buffer.__ref) {
      JANE_ID(panic)(__JANE_ERROR_MEMORY_ALLOCATION_FAILED);
    }
    *_buffer.__ref = __JANE_REFERENCE_DELTA;
    *_buffer.__aloc = _Instance;
    return (_buffer);
  }

  ref_jnt<T>(void) noexcept {}

  ref_jnt<T>(const ref_jnt<T> &_Ref) noexcept { this->operator=(_Ref); }

  ~ref_jnt<T>(void) noexcept { this->_drop(); }

  inline int_jnt __drop_ref(void) const noexcept {
    return (__jane_atomic_add(this->__ref, -__JANE_REFERENCE_DELTA));
  }

  inline int_jnt __add_ref(void) const noexcept {
    return (__jane_atomic_add(this->__ref, __JANE_REFERENCE_DELTA));
  }

  inline uint_jnt __get_ref_n(void) const noexcept {
    return (__jane_atomic_load(this->__ref));
  }

  void _drop(void) const noexcept {
    if (!this->__ref) {
      this->__alloc = nil;
      return;
    }
    if ((this->__drop_ref()) != __JANE_REFERENCE_DELTA) {
      this->__ref = nil;
      this->__alloc = nil;
      return;
    }
    delete this->__ref;
    this->__ref = nil;
    delete this->__alloc;
    this->__alloc = nil;
  }

  inline bool _real() const noexcept { return (this->__alloc != nil); }

  inline T *operator->(void) noexcept {
    this->__must_ok();
    return (*this->__alloc);
  }

  inline operator T(void) noexcept {
    this->__must_ok();
    return (*this->__alloc);
  }

  inline operator T &(void) noexcept {
    this->__must_ok();
    return (*this->__alloc);
  }

  inline void __must_ok(void) const noexcept {
    if (!this->_real()) {
      JANE_ID(panic)(__JANE_ERROR_INVALID_MEMORY);
    }
  }

  void operator=(const ref_jnt<T> &_Ref) noexcept {
    this->_drop();
    if (_Ref.__ref) {
      _Ref.__add_ref();
    }
    this->__ref = _Ref.__ref;
    this->__alloc = _Ref.__alloc;
  }

  inline bool operator==(const T &_Val) const noexcept {
    return (this->__alloc == nil ? false : *this->__alloc == _Val);
  }

  inline bool operator!=(const T &_Val) const noexcept {
    return (!this->operator==(_Val));
  }

  inline bool operator==(const ref_jnt<T> &_Ref) const noexcept {
    if (this->__alloc == nil) {
      return _Ref.__alloc == nil;
    }
    if (_Ref.__alloc == nil) {
      return false;
    }
    return ((*this->__alloc) == (*_Ref.__alloc));
  }

  inline bool operator!=(const ref_jnt<T> &_Ref) const noexcept {
    return (!this->operator==(_Ref));
  }

  friend inline std::ostream &operator<<(std::ostream &_Stream,
                                         const ref_jnt<T> &_Ref) noexcept {
    if (!_Ref._real()) {
      _Stream << "nil";
    } else {
      _Stream << _Ref.operator T();
    }
    return (_Stream);
  }
};

#endif // !__JANE_REF_HPP
