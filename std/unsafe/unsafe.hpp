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

#ifndef __JNC_STD_UNSAFE_UNSAFE_HPP
#define __JNC_STD_UNSAFE_UNSAFE_HPP

#include "../../api/ptr.hpp"
#include "../../api/typedef.hpp"

template <typename T>
inline ptr<T> __jnc_uintptr_cast_to_raw(const uintptr_jnt &_Addr) noexcept;

template <typename T>
inline ptr<T> __jnc_uintptr_cast_to_raw(const uintptr_jnt &_Addr) noexcept {
  ptr<T> _ptr;
  _ptr._ptr = (T **)(&_Addr);
  _ptr._heap = __JNC_PTR_NEVER_HEAP;
  return _ptr;
}

#endif // !__JNC_STD_UNSAFE_UNSAFE_HPP
