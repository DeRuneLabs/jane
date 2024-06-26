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

#ifndef __JANE_DEFER_HPP
#define __JANE_DEFER_HPP

#include <functional>
#define __JANE_CCONCAT(A, B) A##B
#define __JANE_CONCAT(A, B) __JANE_CCONCAT(A, B)

#define __JANE_DEFER(BLOCK)                                                    \
  jane::DeferBase __JANE_CONCAT(__deffered__, __LINE__) { [=] BLOCK }

namespace jane {
struct DeferBase;
struct DeferBase {
public:
  std::function<void(void)> scope;
  DeferBase(const std::function<void(void)> &fn) noexcept { this->scope = fn; }
  ~DeferBase(void) noexcept { this->scope(); }
};
} // namespace jane

#endif // __JANE_DEFER_HPP