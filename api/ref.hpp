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

#ifndef __JNC_REF_HPP
#define __JNC_REF_HPP

#include "atomicity.hpp"
#include "jn_util.hpp"
#include "typedef.hpp"

template <typename T> struct jn_ref;

template <typename T> struct jn_ref {
  T *_alloc{nil};
  mutable uint_jnt *_ref{nil};

  jn_ref<T>(void) noexcept : jn_ref<T>(T()) {}
  jn_ref<T>(std::nullptr_t) noexcept {}

  jn_ref<T>(T *_Ptr, uint_jnt *_Ref) noexcept {
    this->_alloc = _Ptr;
    this->_ref = _Ref;
  }

  jn_ref<T>(T *_Ptr) noexcept {
    this->_ref = (new (std::nothrow) uint_jnt);
    if (!this->_ref) {
      JNC_ID(panic)(__JNC_ERROR_MEMORY_ALLOCATION_FAILED);
    }
    *this->_ref = 1;
    this->_alloc = _Ptr;
  }

  jn_ref<T>(const T &_Instance) noexcept {
    this->_alloc = (new (std::nothrow) T);
    if (!this->_alloc) {
      JNC_ID(panic)(__JNC_ERROR_MEMORY_ALLOCATION_FAILED);
    }
    this->_ref = (new (std::nothrow) uint_jnt);
    if (!this->_ref) {
      JNC_ID(panic)(__JNC_ERROR_MEMORY_ALLOCATION_FAILED);
    }
    *this->_ref = __JNC_REFERENCE_DELTA;
    *this->_alloc = _Instance;
  }

  jn_ref<T>(const jn_ref<T> &_Ref) noexcept { this->operator=(_Ref); }

  ~jn_ref<T>(void) noexcept { this->__drop(); }

  inline int_jnt __drop_ref(void) const noexcept {
    return (__jnc_atomic_add(this->_ref, -__JNC_REFERENCE_DELTA));
  }

  inline int_jnt __add_ref(void) const noexcept {
    return (__jnc_atomic_add(this->_ref, __JNC_REFERENCE_DELTA));
  }

  inline uint_jnt __get_ref_n(void) const noexcept {
    return (__jnc_atomic_load(this->_ref));
  }

  void __drop(void) noexcept {
    if (!this->_ref) {
      return;
    }
    if ((this->__drop_ref()) != __JNC_REFERENCE_DELTA) {
      return;
    }
    delete this->_ref;
    this->_ref = nil;
    delete this->_alloc;
    this->_alloc = nil;
  }

  inline T *operator->(void) noexcept { return (this->_alloc); }

  inline operator T(void) const noexcept { return (*this->_alloc); }

  inline operator T &(void) noexcept { return (*this->_alloc); }

  void operator=(const jn_ref<T> &_Ref) noexcept {
    this->__drop();
    if (_Ref._ref) {
      _Ref.__add_ref();
    }
    this->_ref = _Ref._ref;
    this->_alloc = _Ref._alloc;
  }

  inline void operator=(const T &_Val) const noexcept {
    (*this->_alloc) = (_Val);
  }

  inline bool operator==(const jn_ref<T> &_Ref) const noexcept {
    return ((*this->_alloc) == (*_Ref._alloc));
  }

  inline bool operator!=(const jn_ref<T> &_Ref) const noexcept {
    return (!this->operator==(_Ref));
  }

  friend inline std::ostream &operator<<(std::ostream &_Stream,
                                         const jn_ref<T> &_Ref) noexcept {
    return (_Stream << _Ref.operator T());
  }
};

#endif // !__JNC_REF_HPP
