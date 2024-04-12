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
#include "str.hpp"
#include "trait.hpp"
#include "typedef.hpp"
#include "ptr.hpp"

typedef u8_jnt JNID(byte);
typedef i32_jnt JNID(rune);

// declaration
template <typename _Obj_t> inline void JNID(print)(const _Obj_t _Obj) noexcept;

template <typename _Obj_t>
inline void JNID(println)(const _Obj_t _Obj) noexcept;

struct JNID(Error);
inline void JNID(panic)(trait<JNID(Error)> _Error);

template <typename _Item_t>
int_jnt JNID(copy)(const slice<_Item_t> &_Dest,
                   const slice<_Item_t> &_Src) noexcept;
template <typename _Item_t>
slice<_Item_t> JNID(append)(const slice<_Item_t> &_Src,
                            const slice<_Item_t> &_Components) noexcept;

template <typename T> ptr<T> JNID(new)(void) noexcept;

// definition
template <typename _Obj_t> inline void JNID(print)(const _Obj_t _Obj) noexcept {
  std::cout << _Obj;
}

template <typename _Obj_t>
inline void JNID(println)(const _Obj_t _Obj) noexcept {
  JNID(print)<_Obj_t>(_Obj);
  std::cout << std::endl;
}

struct JNID(Error) {
  virtual str_jnt error(void) = 0;
};

inline void JNID(panic)(trait<JNID(Error)> _Error) { throw _Error; }

template <typename _Item_t>
inline slice<_Item_t> JNID(make)(const int_jnt &_N) noexcept {
  return _N < 0 ? nil : slice<_Item_t>(_N);
}

template <typename _Item_t>
int_jnt JNID(copy)(const slice<_Item_t> &_Dest,
                   const slice<_Item_t> &_Src) noexcept {
  if (_Dest.empty() || _Src.empty()) {
    return 0;
  }
  int_jnt _len;
  if (_Dest.len() > _Src.len()) {
    _len = _Src.len();
  } else if (_Src.len() > _Dest.len()) {
    _len = _Dest.len();
  } else {
    _len = _Src.len();
  }
  for (int_jnt _index{0}; _index < _len; ++_index) {
    _Dest._buffer[_index] = _Src._buffer[_index];
  }
  return _len;
}

template <typename _Item_t>
slice<_Item_t> JNID(append)(const slice<_Item_t> &_Src,
                            const slice<_Item_t> &_Components) noexcept {
  const int_jnt _N{_Src.len() + _Components.len()};
  slice<_Item_t> _buffer{JNID(make) < _Item_t > (_N)};
  JNID(copy)<_Item_t>(_buffer, _Src);
  for (int_jnt _index{0}; _index < _Components.len(); ++_index) {
    _buffer[_Src.len() + _index] = _Components._buffer[_index];
  }
  return _buffer;
}

template <typename T>
ptr<T> JNID(new)(void) noexcept {
  ptr<T> _ptr;
  _ptr._heap = new(std::nothrow) bool* {__JNC_PTR_HEAP_TRUE};
  if (!_ptr._heap) {
    JNID(panic)("memory allocation failed");
  }
  *_ptr._ptr = new(std::nothrow) T;
  if (!*_ptr._ptr) {
    JNID(panic)("memory allocation failed");
  }
  _ptr._ref = new(std::nothrow) uint_jnt{1};
  if (!_ptr._ref) {
    JNID(panic)("memory allocation failed");
  }
  return _ptr;
}

#endif // !__JNC_BUILTIN_HPP
