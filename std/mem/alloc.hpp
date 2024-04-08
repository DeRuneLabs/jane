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

#ifndef __JNC_STD_MEM_ALLOC_HPP
#define __JNC_STD_MEM_ALLOC_HPP

#include "../../api/ptr.hpp"

template <typename T> ptr<T> __jnc_new_heap_ptr(void) noexcept;

template <typename T> ptr<T> __jnc_new_heap_ptr(void) noexcept {
  ptr<T> _ptr;
  _ptr._heap = new (std::nothrow) bool *{__JNC_PTR_HEAP_TRUE};
  if (!_ptr._heap) {
    JNID(panic)("memory allocation failed");
  }
  *_ptr._ptr = new (std::nothrow) T;
  if (!*_ptr._ptr) {
    JNID(panic)("memory allocation failed");
  }
  _ptr._ref = new (std::nothrow) uintptr_jnt{1};
  if (!_ptr._ref) {
    JNID(panic)("memory allocation failed");
  }
  return _ptr;
}

#endif // !__JNC_STD_MEM_ALLOC_HPP
