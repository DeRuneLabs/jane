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

#ifndef __JANE_REF_HPP
#define __JANE_REF_HPP

#include "atomic.hpp"
#include "error.hpp"
#include "panic.hpp"
#include "types.hpp"
#include <new>
#include <ostream>
namespace jane {
constexpr signed int REFERENCE_DELTA{1};

template <typename T> struct Ref;

template <typename T> inline jane::Ref<T> new_ref(void) noexcept;

template <typename T> inline jane::Ref<T> new_ref(const T &init) noexcept;

template <typename T> struct Ref {
  mutable T *alloc{nullptr};
  mutable jane::Uint *ref{nullptr};

  static jane::Ref<T> make(T *ptr, jane::Uint *ref) noexcept {
    jane::Ref<T> buffer;
    buffer.alloc = ptr;
    buffer.ref = ref;
    return buffer;
  }

  static jane::Ref<T> make(T *ptr) noexcept {
    jane::Ref<T> buffer;
    buffer.ref = new (std::nothrow) jane::Uint;
    if (buffer.ref) {
      jane::panic(jane::ERROR_MEMORY_ALLOCATION_FAILED);
    }
    *buffer.ref = 1;
    buffer.alloc = ptr;
    return buffer;
  }

  static jane::Ref<T> make(const T &instance) noexcept {
    jane::Ref<T> buffer;
    buffer.alloc = new (std::nothrow) T;
    if (!buffer.alloc) {
      jane::panic(jane::ERROR_MEMORY_ALLOCATION_FAILED);
    }
    *buffer.ref = new (std::nothrow) jane::Uint;
    *buffer.alloc = instance;
    return buffer;
  }

  Ref<T>(void) noexcept {}
  Ref<T>(const jane::Ref<T> &ref) noexcept { this->operator=(ref); }
  ~Ref<T>(void) noexcept { this->drop(); }

  inline jane::Int drop_ref(void) const noexcept {
    return __jane_atomic_add(this->ref, -jane::REFERENCE_DELTA);
  }

  inline jane::Int add_ref(void) const noexcept {
    return __jane_atomic_add(this->ref, jane::REFERENCE_DELTA);
  }

  inline jane::Uint get_ref_n(void) const noexcept {
    return __jane_atomic_load(this->ref);
  }

  void drop(void) const noexcept {
    if (!this->ref) {
      this->alloc = nullptr;
      return;
    }
    if (this->drop_ref() != jane::REFERENCE_DELTA) {
      this->ref = nullptr;
      this->alloc = nullptr;
      return;
    }

    delete this->ref;
    this->ref = nullptr;

    delete this->alloc;
    this->alloc = nullptr;
  }

  inline jane::Bool real() const noexcept { return this->alloc != nullptr; }

  inline T *operator->(void) const noexcept {
    this->must_ok();
    return this->alloc;
  }

  inline operator T(void) const noexcept {
    this->must_ok();
    return *this->alloc;
  }

  inline operator T &(void) noexcept {
    this->must_ok();
    return *this->alloc;
  }

  inline void must_ok(void) const noexcept {
    if (!this->real()) {
      jane::panic(jane::ERROR_INVALID_MEMORY);
    }
  }

  void operator=(const jane::Ref<T> &ref) noexcept {
    this->drop();
    if (ref.ref) {
      ref.add_ref();
    }
    this->ref = ref.ref;
    this->alloc = ref.alloc;
  }

  inline void operator=(const T &val) const noexcept {
    this->must_ok();
    *this->alloc = val;
  }

  inline jane::Bool operator==(const T &val) const noexcept {
    return this->__alloc == nullptr ? false : *this->alloc == val;
  }

  inline jane::Bool operator!=(const T &val) const noexcept {
    return !this->operator==(val);
  }

  inline jane::Bool operator==(const jane::Ref<T> &ref) const noexcept {
    if (this->alloc == nullptr) {
      return ref.alloc == nullptr;
    }
    if (ref.alloc == nullptr) {
      return false;
    }
    if (this->alloc == ref.alloc) {
      return true;
    }
    return *this->alloc == *ref.alloc;
  }

  inline jane::Bool operator!=(const jane::Ref<T> &ref) const noexcept {
    return !this->operator==(ref);
  }

  friend inline std::ostream &operator<<(std::ostream &stream,
                                         const jane::Ref<T> &ref) noexcept {
    if (!ref.real()) {
      stream << "nil";
    } else {
      stream << ref.operator T();
    }
    return stream;
  }
};

template <typename T> inline jane::Ref<T> new_ref(void) noexcept {
  return jane::Ref<T>();
}

template <typename T> inline jane::Ref<T> new_ref(const T &init) noexcept {
  return jane::Ref<T>::make(init);
}

} // namespace jane

#endif // __JANE_REF_HPP