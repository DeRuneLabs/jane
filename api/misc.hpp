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

#ifndef __JANE_MISC_HPP
#define __JANE_MISC_HPP

#include "error.hpp"
#include "panic.hpp"
#include "ref.hpp"
#include "types.hpp"
namespace jane {
template <typename T, typename Denominator>
auto div(const T &x, const Denominator &denominator) noexcept;

template <typename T> jane::Ref<T> new_struct(T *ptr);

template <typename T, typename Denominator>
auto div(const T &x, const Denominator &denominator) noexcept {
  if (denominator == 0) {
    jane::panic(jane::ERROR_DIVIDE_BY_ZERO);
  }
  return (x / denominator);
}

template <typename T> jane::Ref<T> new_struct(T *ptr) {
  if (!ptr) {
    jane::panic(jane::ERROR_MEMORY_ALLOCATION_FAILED);
  }
  ptr->self.ref = new (std::nothrow) jane::Uint;
  if (!ptr->self.ref) {
    jane::panic(jane::ERROR_MEMORY_ALLOCATION_FAILED);
  }
  *ptr->self.ref = 0;
  return ptr->self;
}
} // namespace jane

#endif //__JANE_MISC_HPP