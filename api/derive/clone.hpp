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

#ifndef __JANE_DERIVE_CLONE_HPP
#define __JANE_DERIVE_CLONE_HPP

#include "../array.hpp"
#include "../fn.hpp"
#include "../map.hpp"
#include "../str.hpp"
#include "../trait.hpp"
#include "../types.hpp"

namespace jane {
char clone(const char &x) noexcept;
signed char clone(const signed char &x) noexcept;
unsigned char clone(const unsigned char &x) noexcept;
char *clone(char *x) noexcept;
const char *clone(const char *x) noexcept;
jane::Int clone(const jane::Int &x) noexcept;
jane::Uint clone(const jane::Uint &x) noexcept;
jane::Bool clone(const jane::Bool &x) noexcept;
jane::Str clone(const jane::Str &x) noexcept;

template <typename Item>
jane::Slice<Item> clone(const jane::Slice<Item> &s) noexcept;
template <typename Item, const jane::Uint N>
jane::Array<Item, N> clone(const jane::Array<Item, N> &arr) noexcept;
template <typename Key, typename Value>
jane::Map<Key, Value> clone(const jane::Map<Key, Value> &m) noexcept;
template <typename T> jane::Ref<T> clone(const jane::Ref<T> &r) noexcept;
template <typename T> jane::Trait<T> clone(const jane::Trait<T> &t) noexcept;
template <typename T> jane::Fn<T> clone(const jane::Fn<T> &fn) noexcept;
template <typename T> T *clone(T *ptr) noexcept;
template <typename T> const T *clone(const T *ptr) noexcept;
template <typename T> T clone(const T &t) noexcept;

inline char clone(const char &x) noexcept { return x; }
inline signed char clone(const signed char &x) noexcept { return x; }
inline unsigned char clone(const unsigned char &x) noexcept { return x; }
inline char *clone(char *x) noexcept { return x; }
inline const char *clone(const char *x) noexcept { return x; }

inline jane::Int clone(const jane::Int &x) noexcept { return x; }

inline jane::Uint clone(const jane::Uint &x) noexcept { return x; }

inline jane::Bool clone(const jane::Bool &x) noexcept { return x; }

inline jane::Str clone(const jane::Str &x) noexcept { return x; }

template <typename Item>
jane::Slice<Item> clone(const jane::Slice<Item> &s) noexcept {
  jane::Slice<Item> s_clone(s.len());
  for (int i{0}; i < s.len(); ++i) {
    s_clone.operator[](i) = jane::clone(s.operator[](i));
  }
  return s_clone;
}

template <typename Item, const jane::Uint N>
jane::Array<Item, N> clone(const jane::Array<Item, N> &arr) noexcept {
  jane::Array<Item, N> arr_clone{};
  for (int i{0}; i < arr.len(); ++i) {
    arr_clone.operator[](i) = jane::clone(arr.operator[](i));
  }
  return arr_clone;
}

template <typename Key, typename Value>
jane::Map<Key, Value> clone(const jane::Map<Key, Value> &m) noexcept {
  jane::Map<Key, Value> m_clone;
  for (const auto &pair : m) {
    m_clone[jane::clone(pair.first)] = jane::clone(pair.second);
  }
  return m_clone;
}

template <typename T> jane::Ref<T> clone(const jane::Ref<T> &r) noexcept {
  if (!r.real()) {
    return r;
  }
  jane::Ref<T> r_clone{jane::Ref<T>::make(jane::clone(r.operator T()))};
  return r_clone;
}

template <typename T> jane::Fn<T> clone(const jane::Fn<T> &fn) noexcept {
  return fn;
}

template <typename T> T *clone(T *ptr) noexcept { return ptr; }

template <typename T> const T *clone(const T *ptr) noexcept { return ptr; }

template <typename T> T clone(const T &t) noexcept { return t.clone(); }

} // namespace jane

#endif // __JANE_DERIVE_CLONE_HPP