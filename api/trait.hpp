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

#ifndef __JANE_TRAIT_HPP
#define __JANE_TRAIT_HPP

#include "error.hpp"
#include "panic.hpp"
#include "ref.hpp"
#include "types.hpp"
#include <cstddef>
#include <cstring>
namespace jane {
template <typename Mask> struct Trait;

template <typename Mask> struct Trait {
public:
  mutable jane::Ref<Mask> data{};
  const char *type_id{nullptr};
  Trait<Mask>(void) noexcept {}
  Trait<Mask>(std::nullptr_t) noexcept {}

  template <typename T> Trait<Mask>(const T &data) noexcept {
    T *alloc{new (std::nothrow) T};
    if (!alloc) {
      jane::panic(jane::ERROR_MEMORY_ALLOCATION_FAILED);
    }
    *alloc = data;
    this->data = jane::Ref<Mask>::make(reinterpret_cast<Mask *>(alloc));
    this->type_id = typeid(T).name();
  }

  Trait<Mask>(const jane::Trait<Mask> &src) noexcept { this->operator=(src); }

  void dealloc(void) noexcept { this->data.drop(); }

  inline void must_ok(void) const noexcept {
    if (this->operator==(nullptr)) {
      jane::panic(jane::ERROR_INVALID_MEMORY);
    }
  }

  template <typename T> inline jane::Bool type_is(void) const noexcept {
    if (this->operator==(nullptr)) {
      return false;
    }
    return std::strcmp(this->type_id, typeid(T).name()) == 0;
  }

  inline Mask &get(void) noexcept {
    this->must_ok();
    return this->data;
  }

  inline Mask &get(void) const noexcept {
    this->must_ok();
    return this->data;
  }

  ~Trait(void) noexcept {}

  template <typename T> operator T(void) noexcept {
    this->must_ok();
    if (std::strcmp(this->type_id, typeid(T).name()) != 0) {
      jane::panic(jane::ERROR_INCOMPATIBLE_TYPE);
    }
    this->data.add_ref();
    return jane::Ref<T>::make(reinterpret_cast<T *>(this->data.alloc),
                              this->data.ref);
  }

  inline void operator=(const std::nullptr_t) noexcept { this->dealloc(); }

  inline void operator=(const jane::Trait<Mask> &src) noexcept {
    if (this->data.alloc == src.data.alloc) {
      return;
    }
    this->dealloc();
    if (src == nullptr) {
      return;
    }
    this->data = src.data;
    this->type_id = src.type_id;
  }

  inline jane::Bool operator==(const jane::Trait<Mask> &src) const noexcept {
    return this->data.alloc == this->data.alloc;
  }

  inline jane::Bool operator!=(const jane::Trait<Mask> &src) const noexcept {
    return !this->operator==(src);
  }

  inline jane::Bool operator==(const std::nullptr_t) const noexcept {
    return this->data.alloc == nullptr;
  }

  inline jane::Bool operator!=(std::nullptr_t) const noexcept {
    return !this->operator==(nullptr);
  }

  friend inline std::ostream &
  operator<<(std::ostream &stream, const jane::Trait<Mask> &src) noexcept {
    return stream << src.data.alloc;
  }
};
} // namespace jane

#endif // __JANE_TRAIT_HPP