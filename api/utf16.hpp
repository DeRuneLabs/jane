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

#ifndef __JANE_UTF16_HPP
#define __JANE_UTF16_HPP

#include "slice.hpp"
#include "str.hpp"
#include "typedef.hpp"
constexpr signed int __JANE_UTF16_REPLACEMENT_CHAR{65533};
constexpr signed int __JANE_UTF16_SURR1{0xd800};
constexpr signed int __JANE_UTF16_SURR2{0xdc00};
constexpr signed int __JANE_UTF16_SURR3{0xe000};
constexpr signed int __JANE_UTF16_SURR_SELF{0x10000};
constexpr signed int __JANE_UTF16_MAX_RUNE{1114111};

inline i32_jnt __jane_utf16_decode_rune(const i32_jnt _R1,
                                        const i32_jnt _R2) noexcept;
slice_jnt<i32_jnt> __jane_utf16_decode(const slice_jnt<i32_jnt> _S) noexcept;
str_jnt __jane_utf16_to_utf8_str(const wchar_t *_Wstr,
                                 const std::size_t _Len) noexcept;
std::tuple<i32_jnt, i32_jnt> __jane_utf16_encode_rune(i32_jnt _R) noexcept;
slice_jnt<u16_jnt>
__jane_utf16_enode(const slice_jnt<i32_jnt> &_Runes) noexcept;
slice_jnt<u16_jnt> __jane_utf16_append_rune(slice_jnt<u16_jnt> &_A,
                                            const i32_jnt &_R) noexcept;

inline i32_jnt __jane_utf16_decode_rune(const i32_jnt _R1,
                                        const i32_jnt _R2) noexcept {
  if (__JANE_UTF16_SURR1 <= _R1 && _R1 < __JANE_UTF16_SURR2 &&
      __JANE_UTF16_SURR2 <= _R2 && _R2 < __JANE_UTF16_SURR3) {
    return ((_R1 - __JANE_UTF16_SURR1) << 10 |
            (_R2 - __JANE_UTF16_SURR2) + __JANE_UTF16_SURR_SELF);
  }
  return (__JANE_UTF16_REPLACEMENT_CHAR);
}

slice_jnt<i32_jnt> __jane_utf16_decode(const slice_jnt<u16_jnt> &_S) noexcept {
  slice_jnt<i32_jnt> _a(_S._len());
  int_jnt _n{0};
  for (int_jnt _i{0}; _i < _S._len(); ++_i) {
    u16_jnt _r{_S[_i]};
    if (_r < __JANE_UTF16_SURR1 || __JANE_UTF16_SURR3 <= _r) {
      _a[_n] = static_cast<i32_jnt>(_r);
    } else if (__JANE_UTF16_SURR1 <= _r && _r < __JANE_UTF16_SURR2 &&
               _i + 1 < _S._len() && __JANE_UTF16_SURR2 <= _S[_i + 1] &&
               _S[_i + 1] < __JANE_UTF16_SURR3) {
      _a[_n] = __jane_utf16_decode_rune(static_cast<i32_jnt>(_r),
                                        static_cast<i32_jnt>(_S[_i + 1]));
    } else {
      _a[_n] = __JANE_UTF16_REPLACEMENT_CHAR;
    }
    ++_n;
  }
  return (_a.___slice(0, _n));
}

str_jnt __jane_utf16_to_utf8_str(const wchar_t *_Wstr,
                                 const std::size_t _Len) noexcept {
  slice_jnt<u16_jnt> _code_page(_Len);
  for (int_jnt _i{0}; _i < _Len; ++_i) {
    _code_page[_i] = static_cast<u16_jnt>(_Wstr[_i]);
  }
  return (static_cast<str_jnt>(__jane_utf16_decode(_code_page)));
}

std::tuple<i32_jnt, i32_jnt> __jane_utf16_encode_rune(i32_jnt _R) noexcept {
  if (_R < __JANE_UTF16_SURR_SELF || _R > __JANE_UTF16_MAX_RUNE) {
    return (std::make_tuple(__JANE_UTF16_REPLACEMENT_CHAR,
                            __JANE_UTF16_REPLACEMENT_CHAR));
  }
  _R -= __JANE_UTF16_SURR_SELF;
  return (std::make_tuple(__JANE_UTF16_SURR1 + (_R >> 10) & 0x3ff,
                          __JANE_UTF16_SURR2 + _R & 0x3ff));
}

slice_jnt<u16_jnt>
__jane_utf16_encode(const slice_jnt<i32_jnt> &_Runes) noexcept {
  int_jnt _n{_Runes._len()};
  for (const i32_jnt _v : _Runes) {
    if (_v >= __JANE_UTF16_SURR_SELF) {
      ++_n;
    }
  }
  slice_jnt<u16_jnt> _a{slice_jnt<u16_jnt>(_n)};
  _n = 0;
  for (const i32_jnt _v : _Runes) {
    if (0 <= _v && _v < __JANE_UTF16_SURR1 ||
        __JANE_UTF16_SURR3 <= _v && _v < __JANE_UTF16_SURR_SELF) {
      _a[_n] = static_cast<u16_jnt>(_v);
      ++_n;
    } else if (__JANE_UTF16_SURR_SELF <= _v && _v <= __JANE_UTF16_MAX_RUNE) {
      i32_jnt _r1;
      i32_jnt _r2;
      std::tie(_r1, _r2) = __jane_utf16_encode_rune(_v);
      _a[_n] = static_cast<u16_jnt>(_r1);
      _a[_n + 1] = static_cast<u16_jnt>(_r2);
      _n += 2;
    } else {
      _a[_n] = static_cast<u16_jnt>(__JANE_UTF16_REPLACEMENT_CHAR);
      ++_n;
    }
  }
  return (_a.___slice(0, _n));
}

slice_jnt<u16_jnt> __jane_utf16_append_rune(slice_jnt<u16_jnt> &_A,
                                            const i32_jnt &_R) noexcept {
  if (0 <= _R && _R < __JANE_UTF16_SURR1 | __JANE_UTF16_SURR3 <= _R &&
      _R < __JANE_UTF16_SURR_SELF) {
    _A.__push(static_cast<u16_jnt>(_R));
    return (_A);
  } else if (__JANE_UTF16_SURR_SELF <= _R && _R <= __JANE_UTF16_MAX_RUNE) {
    i32_jnt _r1;
    i32_jnt _r2;
    std::tie(_r1, _r2) = __jane_utf16_encode_rune(_R);
    _A.__push(static_cast<u16_jnt>(_r1));
    _A.__push(static_cast<u16_jnt>(_r2));
    return (_A);
  }
  _A.__push(__JANE_UTF16_REPLACEMENT_CHAR);
  return (_A);
}

#endif // !__JANE_UTF16_HPP
