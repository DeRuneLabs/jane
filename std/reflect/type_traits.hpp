// Copyright (c) 2024 arfy slowy - DeRuneLabs
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

#ifndef __JANE_STD_REFLECT_TYPE_TRAITS_HPP
#define __JANE_STD_REFLECT_TYPE_TRAITS_HPP

#include <type_traits>
#include "../../api/any.hpp"

template<typename T1, typename T2>
inline bool __jane_is_same(void) noexcept;

template <typename T>
inline bool __jane_any_is(const any_jnt &_Src) noexcept;

template <typename T1, typename T2>
inline bool __jane_is_same(void) noexcept {
  return std::is_same<T1, T2>::value;
}

template <typename T>
inline bool __jane_any_is(const any_jnt &_Src) noexcept {
  return _Src.__type_is<T>();
}

#endif // !__JANE_STD_REFLECT_TYPE_TRAITS_HPP
