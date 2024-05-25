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

#ifndef __JANE_ERROR_HPP
#define __JANE_ERROR_HPP

#define __JANE_WRITE_ERROR_SLICING_INDEX_OUT_OF_RANGE(STREAM, START, LEN)      \
  (STREAM << jane::ERROR_INDEX_OUT_OF_RANGE << '[' << START << ':' << LEN      \
          << ']')

#define __JANE_WRITE_ERROR_INDEX_OUT_OF_RANGE(STREAM, INDEX)                   \
  (STREAM << jane::ERROR_INDEX_OUT_OF_RANGE << '[' << INDEX << ']')

namespace jane {
constexpr const char *ERROR_INVALID_MEMORY{
    "invalid memory address or nil pointer deference"};
constexpr const char *ERROR_INCOMPATIBLE_TYPE{"incompatible type"};
constexpr const char *ERROR_MEMORY_ALLOCATION_FAILED{
    "memory allocation failed"};
constexpr const char *ERROR_INDEX_OUT_OF_RANGE{"index out of range"};
constexpr const char *ERROR_DIVIDE_BY_ZERO{"divide by zero"};
constexpr signed int EXIT_PANIC{2};
} // namespace jane

#endif // __JANE_ERROR_HPP