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

#ifndef __JNC_TRACER_HPP
#define __JNC_TRACER_HPP

#include "jn_util.hpp"
#include "str.hpp"
#include "typedef.hpp"

struct tracer;

struct tracer {
  static constexpr uint_jnt _n{20};

  std::array<str_jnt, _n> _traces;

  void push(const str_jnt &_Src) {
    for (int_jnt _index{_n - 1}; _index > 0; _index--) {
      this->_traces[_index] = this->_traces[_index - 1];
    }
    this->_traces[0] = _Src;
  }

  str_jnt string(void) noexcept {
    str_jnt _traces{};
    for (const str_jnt &_trace : this->_traces) {
      if (_trace.empty()) {
        break;
      }
      _traces += _trace;
      _traces += "\n";
    }
    return _traces;
  }

  void ok(void) noexcept {
    for (int_jnt _index{0}; _index < _n; _index++) {
      this->_traces[_index] = this->_traces[_index + 1];
      if (this->_traces[_index + 1].empty()) {
        break;
      }
    }
  }
};

tracer ___trace{};

#endif // !_JNC_TRACER_HPP
