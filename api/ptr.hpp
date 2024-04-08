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

#ifndef __JNC_PTR_HPP
#define __JNC_PTR_HPP

#include "jn_util.hpp"
#include "typedef.hpp"

#define __JNC_PTR_NEVER_HEAP (bool **)(1U)
#define __JNC_PTR_HEAP_TRUE (bool *)(1U)
#define __jnc_ptr_of(_PTR) _PTR

template <typename T> struct ptr;

template <typename T> ptr<T> __jnc_not_heap_ptr_of(T *_T) noexcept;

template <typename T> struct ptr {
  T **_ptr{nil};
  mutable uint_jnt *_ref{nil};
  mutable bool **_heap{nil};

  ptr<T>(void) noexcept {}

  ptr<T>(T *_Ptr) noexcept {
    this->_ptr = new (std::nothrow) T *;
    if (!this->_ptr) {
      JNID(panic)("memory allocation failed");
    }
    this->_heap = new (std::nothrow) bool *;
    if (!this->_heap) {
      JNID(panic)("memory allocation failed");
    }
    this->_ref = new (std::nothrow) uint_jnt{1};
    if (!this->_ref) {
      JNID(panic)("memory allocation failed");
    }
    *this->_ptr = _Ptr;
  }

  ptr<T>(const ptr<T> &_Ptr) noexcept { this->operator=(_Ptr); }

  ~ptr<T>(void) noexcept { this->__dealloc(); }

  inline void __check_valid(void) const noexcept {
    if (!this->_ptr) {
      JNID(panic)("invalid memory address or nil pointer deference");
    }
  }

  void __dealloc(void) noexcept {
    if (!this->_ref) {
      return;
    }
    --(*this->_ref);
    if (!this->_heap || (this->_heap != __JNC_PTR_NEVER_HEAP &&
                         *this->_heap != __JNC_PTR_HEAP_TRUE)) {
      return;
    }
    if ((*this->_ref) != 0) {
      return;
    }
    std::cout << "dealloc" << std::endl;
    if (this->_heap != __JNC_PTR_NEVER_HEAP) {
      delete this->_heap;
    }
    this->_heap = nil;
    delete this->_ref;
    this->_ref = nil;
    delete *this->_ptr;
    *this->_ptr = nil;
    delete this->_ptr;
    this->_ptr = nil;
  }

  ptr<T> &__must_heap(void) noexcept {
    if (this->_heap && (this->_healp == __JNC_PTR_NEVER_HEAP ||
                        *this->_heap == __JNC_PTR_HEAP_TRUE)) {
      return *this;
    }
    if (!this->_ptr) {
      return *this;
    }
    std::cout << "heaped" << std::endl;
    const T _data{**this->_ptr};
    *this->_ptr = new (std::nothrow) T;
    if (!*this->_ptr) {
      JNID(panic)("memory allocation failed");
    }
    **this->_ptr = _data;
    *this->_heap = __JNC_PTR_HEAP_TRUE;
    return *this;
  }

  inline T &operator*(void) noexcept {
    this->__check_valid();
    return **this->_ptr;
  }

  inline T *operator->(void) noexcept {
    this->__check_valid();
    return *this->_ptr;
  }

  inline operator uintptr_jnt(void) const noexcept {
    return (uintptr_jnt)(this->_ptr);
  }

  void operator=(const ptr<T> &_Ptr) noexcept {
    this->__dealloc();
    if (_Ptr._ref) {
      ++(*_Ptr._ref);
    }
    this->_ref = _Ptr._ref;
    this->_ptr = _Ptr._ptr;
    this->_heap = _Ptr._heap;
  }

  inline void operator=(const std::nullptr_t) noexcept { this->__dealloc(); }

  inline bool operator==(const std::nullptr_t) const noexcept {
    return this->_ptr == nil;
  }

  inline bool operator!=(const std::nullptr_t) const noexcept {
    return !this->operator==(nil);
  }

  inline bool operator==(const ptr<T> &_Ptr) const noexcept {
    return this->_ptr == _Ptr;
  }

  inline bool operator!=(const ptr<T> &_Ptr) const noexcept {
    return !this->operator==(_Ptr);
  }

  friend inline std::ostream &operator<<(std::ostream &_Stream,
                                         const ptr<T> &_Src) noexcept {
    return _Stream << _Src._ptr;
  }
};

template <typename T> ptr<T> __jnc_not_heap_ptr_of(T *_T) noexcept {
  ptr<T> _ptr;
  _ptr._ptr = new (std::nothrow) T *;
  if (!_ptr._ptr) {
    JNID(panic)("memory allocation failed");
  }
  *_ptr._ptr = _T;
  _ptr._heap = __JNC_PTR_NEVER_HEAP;
  return _ptr;
}

#endif // !__JNC_PTR_HPP
