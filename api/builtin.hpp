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

#ifndef __JNC_BUILTIN_HPP
#define __JNC_BUILTIN_HPP

#include "jn_util.hpp"
#include "ptr.hpp"
#include "str.hpp"
#include "trait.hpp"
#include "typedef.hpp"

typedef u8_jnt JNC_ID(byte);
typedef i32_jnt JNC_ID(rune);

struct JNC_ID(Error);

template <typename _Item_t>
int_jnt JNC_ID(copy)(const slice<_Item_t> &_Dest,
                     const slice<_Item_t> &_Src) noexcept;

template <typename _Item_t>
slice<_Item_t> JNC_ID(append)(const slice<_Item_t> &_Src,
                              const slice<_Item_t> &_Components) noexcept;

template <typename T> ptr<T> JNC_ID(new)(void) noexcept;

struct JNC_ID(Error) {
  virtual str_jnt error(void) = 0;
};

template <typename _Item_t>
inline slice<_Item_t> JNC_ID(make)(const int_jnt &_N) noexcept {
  return _N < 0 ? nil : slice<_Item_t>(_N);
}

template <typename _Item_t>
int_jnt JNC_ID(copy)(const slice<_Item_t> &_Dest,
                     const slice<_Item_t> &_Src) noexcept {
  if (_Dest.empty() || _Src.empty()) {
    return 0;
  }
  int_jnt _len = _Dest.len() > _Src.len()   ? _len = _Src.len()
                 : _Src.len() > _Dest.len() ? _len = _Dest.len()
                                            : _len = _Src.len();
  for (int_jnt _index{0}; _index < _len; ++_index) {
    _Dest._slice[_index] = _Src._slice[_index];
  }
  return _len;
}

template <typename _Item_t>
slice<_Item_t> JNC_ID(append)(const slice<_Item_t> &_Src,
                              const slice<_Item_t> &_Components) noexcept {
  const int_jnt _N{_Src.len() + _Components.len()};
  slice<_Item_t> _buffer{JNC_ID(make) < _Item_t > (_N)};
  JNC_ID(copy)<_Item_t>(_buffer, _Src);
  for (int_jnt _index{0}; _index < _Components.len(); ++_index) {
    _buffer[_Src.len() + _index] = _Components._slice[_index];
  }
  return _buffer;
}

template <typename T> ptr<T> JNC_ID(new)(void) noexcept {
  ptr<T> _ptr;
  _ptr._heap = new (std::nothrow) bool *{__JNC_PTR_HEAP_TRUE};
  if (!_ptr._heap) {
    JNC_ID(panic)(__JNC_ERROR_MEMORY_ALLOCATION_FAILED);
  }
  _ptr._ptr = new (std::nothrow) T *;
  if (!_ptr._ptr) {
    JNC_ID(panic)(__JNC_ERROR_MEMORY_ALLOCATION_FAILED);
  }
  *_ptr._ptr = new (std::nothrow) T;
  if (!*_ptr._ptr) {
    JNC_ID(panic)(__JNC_ERROR_MEMORY_ALLOCATION_FAILED);
  }
  _ptr._ref = new (std::nothrow) uint_jnt{1};
  if (!_ptr._ref) {
    JNC_ID(panic)(__JNC_ERROR_MEMORY_ALLOCATION_FAILED);
  }
  return _ptr;
}

#endif // !__JNC_BUILTIN_HPP
