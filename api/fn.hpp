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

#ifndef __JNC_FN_HPP
#define __JNC_FN_HPP

#include "jn_util.hpp"
#include "slice.hpp"
#include "str.hpp"

template <typename _Fn_t> struct fn;

template <typename _Fn_t> struct fn {
  _Fn_t _buffer;

  fn<_Fn_t>(void) noexcept {}

  fn<_Fn_t>(const _Fn_t &_Fn) noexcept { this->_buffer = _Fn; }

  template <typename... _Args_t> auto operator()(_Args_t... _Args) noexcept {
    if (this->_buffer == nil) {
      JNC_ID(panic)(__JNC_ERROR_INVALID_MEMORY);
    }
    return (this->_buffer(_Args...));
  }

  inline void operator=(std::nullptr_t) noexcept { this->_buffer = nil; }

  inline void operator=(const _Fn_t &_Func) noexcept { this->_buffer = _Func; }

  inline bool operator==(std::nullptr_t) const noexcept {
    return (this->_buffer == nil);
  }

  inline bool operator!=(std::nullptr_t) const noexcept {
    return (!this->operator==(nil));
  }
};

#endif // !__JNC_FN_HPP
