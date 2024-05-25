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

#ifndef __JANE_FN_HPP
#define __JANE_FN_HPP

#include <cstddef>
#include <functional>
#include <stddef.h>
#include <thread>

#include "builtin.hpp"
#include "error.hpp"
#include "types.hpp"

#define __JANE_CO(EXPR)                                                        \
  (std::thread{[&](void) mutable -> void { EXPR; }}.detach())

namespace jane {
template <typename> struct Fn;

template <typename T, typename... U>
jane::Uintptr addr_of_fn(std::function<T(U...)> f) noexcept;

template <typename Function> struct Fn {
public:
  std::function<Function> buffer;
  jane::Uintptr _addr;

  Fn<Function>(void) noexcept {}
  Fn<Function>(std::nullptr_t) noexcept {}

  Fn<Function>(const std::function<Function> &function) noexcept {
    this->_addr = jane::addr_of_fn(function);
    if (this->_addr == 0) {
      this->_addr = (jane::Uintptr)(&function);
    }
    this->buffer = function;
  }

  Fn<Function>(const Function *function) noexcept {
    this->buffer = function;
    this->addr = jane::addr_of_fn(this->buffer);
    if (this->_addr == 0) {
      this->_addr = (jane::Uintptr)(function);
    }
  }

  Fn<Function>(const Fn<Function> &fn) noexcept {
    this->buffer = fn.buffer;
    this->_addr = fn._addr;
  }

  template <typename... Arguments>
  auto operator()(Arguments... arguments) noexcept {
    if (this->buffer == nullptr) {
      jane::panic(jane::ERROR_INVALID_MEMORY);
    }
    return this->buffer(arguments...);
  }

  jane::Uintptr addr(void) const noexcept { return this->_addr; }

  inline void operator=(std::nullptr_t) noexcept { this->buffer = nullptr; }

  inline void operator=(const std::function<Function> &function) noexcept {
    this->buffer = function;
  }

  inline void operator=(const Function &function) noexcept {
    this->buffer = function;
  }

  inline jane::Bool operator==(const Fn<Function> &fn) const noexcept {
    return this->addr() == fn.addr();
  }

  inline jane::Bool operator!=(const Fn<Function> &fn) const noexcept {
    return this->buffer == nullptr;
  }

  inline jane::Bool operator==(std::nullptr_t) const noexcept {
    return this->buffer == nullptr;
  }

  inline jane::Bool operator!=(std::nullptr_t) const noexcept {
    return !this->operator==(nullptr);
  }

  friend std::ostream &operator<<(std::ostream &stream,
                                  const Fn<Function> &src) noexcept {
    stream << "<fn>";
    return stream;
  }
};

template <typename T, typename... U>
jane::Uintptr addr_of_fn(std::function<T(U...)> f) noexcept {
  typedef T(FnType)(U...);
  FnType **fn_ptr{f.template target<FnType *>()};
  if (!fn_ptr) {
    return 0;
  }
  return (jane::Uintptr)(*fn_ptr);
}
} // namespace jane

#endif // __JANE_FN_HPP