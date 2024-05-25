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

#ifndef __JANE_ANY_HPP
#define __JANE_ANY_HPP

#include <cstddef>
#include <cstdlib>
#include <cstring>
#include <ostream>
#include <stddef.h>

#include "builtin.hpp"
#include "error.hpp"
#include "ref.hpp"
#include "str.hpp"
#include "types.hpp"
namespace jane {
class Any;

class Any {
private:
  template <typename T> struct DynamicType {
  public:
    static const char *type_id(void) noexcept { return typeid(T).name(); }

    static void dealloc(void *alloc) noexcept {
      delete reinterpret_cast<T *>(alloc);
    }

    static jane::Bool eq(void *alloc, void *other) noexcept {
      T *l{reinterpret_cast<T *>(alloc)};
      T *r{reinterpret_cast<T *>(other)};
      return *l == *r;
    }

    static const jane::Str to_str(const void *alloc) noexcept {
      const T *v{reinterpret_cast<T *>(alloc)};
      return jane::to_str(*v);
    }
  };

  struct Type {
  public:
    const char *(*type_id)(void) noexcept;
    void (*dealloc)(void *alloc) noexcept;
    jane::Bool (*eq)(void *alloc, void *other) noexcept;
    const jane::Str (*to_str)(const void *alloc) noexcept;
  };

  template <typename T> static jane::Any::Type *new_type(void) noexcept {
    using t = typename std::decay<DynamicType<T>>::type;
    static jane::Any::Type table = {
        t::type_id,
        t::dealloc,
        t::eq,
        t::to_str,
    };
    return &table;
  }

public:
  jane::Ref<void *> data{};
  jane::Any::Type *type{nullptr};
  Any(void) noexcept {}
  template <typename T> Any(const T &expr) noexcept { this->operator=(expr); }

  Any(const jane::Any &src) noexcept { this->operator=(src); }

  ~Any(void) noexcept { this->dealloc(); }

  void dealloc(void) noexcept {
    if (!this->data.ref) {
      this->type = nullptr;
      this->data.alloc = nullptr;
      return;
    }

    if ((this->data.get_ref_n()) != jane::REFERENCE_DELTA) {
      return;
    }
    this->type->dealloc(*this->data.alloc);
    *this->data.alloc = nullptr;
    this->type = nullptr;

    delete this->data.ref;
    this->data.ref = nullptr;
    std::free(this->data.alloc);
    this->data.alloc = nullptr;
  }

  template <typename T> inline jane::Bool type_is(void) const noexcept {
    if (std::is_same<typename std::decay<T>::type, std::nullptr_t>::value) {
      return false;
    }
    if (this->operator==(nullptr)) {
      return false;
    }
    return std::strcmp(this->type->type_id(), typeid(T).name()) == 0;
  }

  template <typename T> void operator=(const T &expr) noexcept {
    this->dealloc();

    T *alloc{new (std::nothrow) void *};
    if (!alloc) {
      jane::panic(jane::ERROR_MEMORY_ALLOCATION_FAILED);
    }

    void **main_alloc{new (std::nothrow) void *};
    if (!main_alloc) {
      jane::panic(jane::ERROR_MEMORY_ALLOCATION_FAILED);
    }

    *alloc = expr;
    *main_alloc = reinterpret_cast<void *>(alloc);
    this->data = jane::Ref<void *>::make(main_alloc);
    this->type = jane::Any::new_type<T>();
  }

  void operator=(const jane::Any &src) noexcept {
    if (this->data.alloc == src.data.alloc) {
      return;
    }
    if (src.operator==(nullptr)) {
      this->operator=(nullptr);
    }

    this->dealloc();
    this->data = src.data;
    this->type = src.type;
  }

  inline void operator=(const std::nullptr_t) noexcept { this->dealloc(); }

  template <typename T> operator T(void) const noexcept {
    if (this->operator==(nullptr)) {
      jane::panic(jane::ERROR_INVALID_MEMORY);
    }

    if (!this->type_is<T>()) {
      jane::panic(jane::ERROR_INCOMPATIBLE_TYPE);
    }

    return *reinterpret_cast<T *>(*this->data.alloc);
  }

  template <typename T>
  inline jane::Bool operator==(const T &_Expr) const noexcept {
    return (this->type_is<T>() && this->operator T() == _Expr);
  }

  template <typename T>
  inline jane::Bool operator!=(const T &_Expr) const noexcept {
    return (!this->operator==(_Expr));
  }

  inline jane::Bool operator==(const jane::Any &other) const noexcept {
    if (this->data.alloc == other.data.alloc) {
      return true;
    }

    if (this->operator==(nullptr) && other.operator==(nullptr)) {
      return true;
    }

    if (std::strcmp(this->type->type_id(), other.type->type_id()) != 0) {
      return false;
    }

    return this->type->eq(*this->data.alloc, *other.data.alloc);
  }

  inline jane::Bool operator!=(const jane::Any &other) const noexcept {
    return !this->operator==(other);
  }

  inline jane::Bool operator==(std::nullptr_t) const noexcept {
    return !this->data.alloc;
  }

  inline jane::Bool operator!=(std::nullptr_t) const noexcept {
    return !this->operator==(nullptr);
  }

  friend std::ostream &operator<<(std::ostream &stream,
                                  const jane::Any &src) noexcept {
    if (src.operator!=(nullptr)) {
      stream << src.type->to_str(*src.data.alloc);
    } else {
      stream << 0;
    }
    return stream;
  }
};
} // namespace jane

#endif //__JANE_ANY_HPP