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

#ifndef __JNC_UTF16_HPP
#define __JNC_UTF16_HPP

#include "jn_util.hpp"
#include "slice.hpp"
#include "typedef.hpp"

constexpr signed int __JNC_UTF16_REPLACEMENT_CHAR{65533};
constexpr signed int __JNC_UTF16_SURR1{0xd800};
constexpr signed int __JNC_UTF16_SURR2{0xdc00};
constexpr signed int __JNC_UTF16_SURR3{0xe000};
constexpr signed int __JNC_UTF16_SURR_SELF{0x10000};
constexpr signed int __JNC_UTF16_MAX_RUNE{1114111};

inline i32_jnt __jnc_utf16_decode_rune(const i32_jnt _R1,
                                       const i32_jnt _R2) noexcept;
slice<i32_jnt> __jnc_utf16_decode(const slice<i32_jnt> _S) noexcept;
str_jnt __julec_utf16_to_utf8_str(const wchar_t *_WStr,
                                  const std::size_t _Len) noexcept;

inline i32_jnt __jnc_utf16_decode_rune(const i32_jnt _R1,
                                       const i32_jnt _R2) noexcept {
  if (__JNC_UTF16_SURR1 <= _R1 && _R1 < __JNC_UTF16_SURR2 &&
      __JNC_UTF16_SURR2 <= _R2 && _R2 < __JNC_UTF16_SURR3) {
    return ((_R1 - __JNC_UTF16_SURR1) << 10 |
            (_R2 - __JNC_UTF16_SURR2) + __JNC_UTF16_SURR_SELF);
  }
  return (__JNC_UTF16_REPLACEMENT_CHAR);
}

slice<i32_jnt> __jnc_utf16_decode(const slice<u16_jnt> &_S) noexcept {
  slice<i32_jnt> _a(_S.len());
  int_jnt _n{0};
  for (int_jnt _i{0}; _i < _S.len(); ++_i) {
    u16_jnt _r{_S[_i]};
    if (_r < __JNC_UTF16_SURR1 || __JNC_UTF16_SURR3 <= _r) {
      _a[_n] = static_cast<i32_jnt>(_r);
    } else if (__JNC_UTF16_SURR1 <= _r && _r < __JNC_UTF16_SURR2 &&
               _i + 1 < _S.len() && __JNC_UTF16_SURR2 <= _S[_i + 1] &&
               _S[_i + 1] < __JNC_UTF16_SURR3) {
      _a[_n] = __jnc_utf16_decode_rune(static_cast<i32_jnt>(_r),
                                       static_cast<i32_jnt>(_S[_i + 1]));
      ++_i;
    } else {
      _a[_n] = __JNC_UTF16_REPLACEMENT_CHAR;
    }
    ++_n;
  }
  return (_a.___slice(0, _n));
}

// TODO: implemented this code
// str_jnt __jnc_utf16_to_utf8_str(const wchar_t *_WStr,
//                                 const std::size_t _Len) noexcept {
//   slice<u16_jnt> _code_page(_Len);
//   for (int_jnt _i{0}; _i < _Len; ++_i) {
//     _code_page[_i] = static_cast<u16_jnt>(_WStr[_i]);
//   }
//   return (static_cast<str_jnt>(__jnc_utf16_decode(_code_page)));
// }

std::tuple<i32_jnt, i32_jnt> __jnc_utf16_encode_rune(i32_jnt _R) noexcept {
  if (_R < __JNC_UTF16_SURR_SELF || _R > __JNC_UTF16_MAX_RUNE) {
    return (std::make_tuple(__JNC_UTF16_REPLACEMENT_CHAR,
                            __JNC_UTF16_REPLACEMENT_CHAR));
  }
  _R -= __JNC_UTF16_SURR_SELF;
  return (std::make_tuple(__JNC_UTF16_SURR1 + (_R >> 10) & 0x3ff,
                          __JNC_UTF16_SURR2 + _R & 0x3ff));
}

slice<u16_jnt> encode(const slice<i32_jnt> &_Runes) noexcept {
  int_jnt _n{_Runes.len()};
  for (const i32_jnt _v : _Runes) {
    if (_v >= __JNC_UTF16_SURR_SELF) {
      ++_n;
    }
  }
  slice<u16_jnt> _a{slice<u16_jnt>(_n)};
  _n = 0;
  for (const i32_jnt _v : _Runes) {
    if (0 <= _v && _v < __JNC_UTF16_SURR1 ||
        __JNC_UTF16_SURR3 <= _v && _v < __JNC_UTF16_SURR_SELF) {
      _a[_n] = static_cast<u16_jnt>(_v);
      ++_n;
    } else if (__JNC_UTF16_SURR_SELF <= _v && _v <= __JNC_UTF16_MAX_RUNE) {
      i32_jnt _r1;
      i32_jnt _r2;
      std::tie(_r1, _r2) = __jnc_utf16_encode_rune(_v);
      _a[_n] = static_cast<u16_jnt>(_r1);
      _a[_n + 1] = static_cast<u16_jnt>(_r2);
      _n += 2;
    } else {
      _a[_n] = static_cast<u16_jnt>(__JNC_UTF16_REPLACEMENT_CHAR);
      ++_n;
    }
  }
  return (_a.___slice(0, _n));
}

#endif // !__JNC_UTF16_HPP
