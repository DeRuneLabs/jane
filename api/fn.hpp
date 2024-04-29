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

#ifndef __JANE_FN_HPP
#define __JANE_FN_HPP

#include "typedef.hpp"
#include <cstddef>

template <typename _Function_t> struct fn_jnt;

template <typename _Function_t> struct fn_jnt {
  std::function<_Function_t> __buffer;

  fn_jnt<_Function_t>(void) noexcept {}

  fn_jnt<_Function_t>(const std::function<_Function_t> &_Function) noexcept {
    this->__buffer = _Function;
  }
  fn_jnt<_Function_t>(const _Function_t &_Function) noexcept {
    this->__buffer = _Function;
  }

  template <typename... _Arguments_t>
  auto operator()(_Arguments_t... _Arguments) noexcept {
    if (this->__buffer == nil) {
      JANE_ID(panic)(__JANE_ERROR_INVALID_MEMORY);
    }
    return (this->__buffer(_Arguments...));
  }

  inline void operator=(std::nullptr_t) noexcept { this->__buffer = nil; }

  inline void operator=(const std::function<_Function_t> &_Function) noexcept {
    this->__buffer = _Function;
  }

  inline void operator=(const _Function_t &_Function) noexcept {
    this->__buffer = _Function;
  }

  inline bool operator==(std::nullptr_t) const noexcept {
    return (this->__buffer == nil);
  }

  inline bool operator!=(std::nullptr_t) const noexcept {
    return (!this->operator==(nil));
  }
};

#endif // !__JANE_FN_HPP
