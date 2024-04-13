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

#ifndef __JNC_DEFER_HPP
#define __JNC_DEFER_HPP

#include "jn_util.hpp"

struct defer;

struct defer {
  typedef std::function<void(void)> _Function_t;
  template <class Callable>
  defer(Callable &&_function) : _function(std::forward<Callable>(_function)) {}
  defer(defer &&_Src) : _function(std::move(_Src._function)) {
    _Src._function = nullptr;
  }
  ~defer() noexcept {
    if (this->_function) {
      this->_function();
    }
  }
  defer(const defer &) = delete;
  void operator=(const defer &) = delete;
  _Function_t _function;
};

#endif // !__JNC_DEFER_HPP
