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

#ifndef __JNC_STR_HPP
#define __JNC_STR_HPP

#include "jn_util.hpp"
#include "typedef.hpp"

class str_jnt {
public:
  std::basic_string<u8_jnt> _buffer{};

  str_jnt(void) noexcept {}
  str_jnt(const char *_Src) noexcept {
    if (!_Src) {
      return;
    }
    this->_buffer =
        std::basic_string<u8_jnt>(&_Src[0], &_Src[std::strlen(_Src)]);
  }

  str_jnt(const std::initializer_list<u8_jnt> &_Src) noexcept {
    this->_buffer = _Src;
  }
};

#endif // !__JNC_STR_HPP
