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

#ifndef __JANE_BUILTIN_HPP
#define __JANE_BUILTIN_HPP

#include "ref.hpp"
#include "slice.hpp"
#include "str.hpp"
#include "typedef.hpp"

typedef u8_jnt(JANE_ID(byte));
typedef i32_jnt(JANE_ID(rune));

template <typename _Obj_t> str_jnt __jane_to_str(const _Obj_t &_Obj) noexcept;
slice_jnt<u16_jnt> __jane_utf16_from_str(const str_jnt &_Str) noexcept;

template <typename _Obj_t>
inline void JANE_ID(print)(const _Obj_t &_Obj) noexcept;
template <typename _Obj_t>
inline void JANE_ID(println)(const _Obj_t &_Obj) noexcept;
struct JANE_ID(Error);
template <typename _Item_t>
int_jnt JANE_ID(copy)(const slice_jnt<_Item_t> &_Dest,
                      const slice_jnt<_Item_t> &Components) noexcept;

template <typename T> inline ref_jnt<T> JANE_ID(new)(void) noexcept;
template <typename T> inline ref_jnt<T> JANE_ID(new)(const T &_Expr) noexcept;
template <typename T> inline void JANE_ID(drop)(T &_Obj) noexcept;
template <typename T> inline bool JANE_ID(real)(T &_Obj) noexcept;

template <typename _Obj_t>
inline void JANE_ID(print)(const _Obj_t &_Obj) noexcept {
#ifdef _WINDOWS
  const str_jnt _str{__jane_to_str<_Obj_t>(_Obj)};
  const slice_jnt<u16_jnt> _utf16_str{__jane_utf16_from_str(_str)};
  HANDLE _handle{GetStdHandle(STD_OUTPUT_HANDLE)};
  WriteConsoleW(_handle, &_utf16_str[0], _utf16_str._len(), nullptr, nullptr);
#else
  std::cout << _Obj;
#endif // DEBUG
}

template <typename _Obj_t>
inline void JANE_ID(println)(const _Obj_t &_Obj) noexcept {
  JANE_ID(print)(_Obj);
  std::cout << std::endl;
}

struct JANE_ID(Error) {
  virtual str_jnt _error(void) { return {}; }
  virtual ~JANE_ID(Error)(void) noexcept {}

  bool operator==(const JANE_ID(Error) & _Src) { return false; }
  bool operator!=(const JANE_ID(Error) & _Src) {
    return !this->operator==(_Src);
  }
};

template <typename _Item_t>
int_jnt JANE_ID(copy)(const slice_jnt<_Item_t> &_Dest,
                      const slice_jnt<_Item_t> &_Src) noexcept {
  if (_Dest._empty() || _Src._empty()) {
    return 0;
  }
  int_jnt _len = (_Dest._len() > _Src._len())   ? _Src._len()
                 : (_Src._len() > _Dest._len()) ? _Dest._len()
                                                : _Src._len();
  for (int_jnt _index{0}; _index < _len; ++_index) {
    _Dest.__slice[_index] = _Src.__slice[_index];
  }
  return (_len);
}

template <typename _Item_t>
slice_jnt<_Item_t>
JANE_ID(append)(const slice_jnt<_Item_t> &_Src,
                const slice_jnt<_Item_t> &_Components) noexcept {
  const int_jnt _N{_Src._len() + _Components._len()};
  slice_jnt<_Item_t> _buffer(_N);
  JANE_ID(copy)<_Item_t>(_buffer, _Src);
  for (int_jnt _index{0}; _index < _Components._len(); ++_index) {
    _buffer[_Src._len() + _index] = _Components.__slice[_index];
  }
  return (_buffer);
}

template <typename T> inline ref_jnt<T> JANE_ID(new)(void) noexcept {
  return (ref_jnt<T>());
}

template <typename T> inline ref_jnt<T> JANE_ID(new)(const T &_Expr) noexcept {
  return (ref_jnt<T>::make(_Expr));
}

template <typename T> inline void JANE_ID(drop)(T &_Obj) noexcept {
  _Obj._drop();
}

template <typename T> inline bool JANE_ID(real)(T &_Obj) noexcept {
  return (_Obj._real());
}

#endif // !__JANE_BUILTIN_HPP
