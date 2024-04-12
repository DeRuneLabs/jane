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

#ifndef __JNC_UTF8_HPP
#define __JNC_UTF8_HPP

#include "slice.hpp"
#include "typedef.hpp"

#define __JNC_UTF8_RUNE_ERROR 65533
#define __JNC_UTF8_MASKX 63
#define __JNC_UTF8_MASK2 31
#define __JNC_UTF8_MASK3 15
#define __JNC_UTF8_MASK4 7
#define __JNC_UTF8_LOCB 128
#define __JNC_UTF8_HICB 191
#define __JNC_UTF8_XX 241
#define __JNC_UTF8_AS 240
#define __JNC_UTF8_S1 2
#define __JNC_UTF8_S2 19
#define __JNC_UTF8_S3 3
#define __JNC_UTF8_S4 35
#define __JNC_UTF8_S5 52
#define __JNC_UTF8_S6 4
#define __JNC_UTF8_S7 68
#define __JNC_UTF8_RUNE1_MAX 127
#define __JNC_UTF8_RUNE2_MAX 2047
#define __JNC_UTF8_RUNE3_MAX 65535
#define __JNC_UTF8_TX 128
#define __JNC_UTF8_T2 192
#define __JNC_UTF8_T3 224
#define __JNC_UTF8_T4 240
#define __JNC_UTF8_T5 248
#define __JNC_UTF8_MAX_RUNE 1114111
#define __JNC_UTF8_SURROGATE_MIN 55296
#define __JNC_UTF8_SURROGATE_MAX 57343

struct __jnc_utf8_accept_range;
std::tuple<i32_jnt, int_jnt>
__jnc_utf8_decode_rune_str(const char *_S) noexcept;
slice<u8_jnt> __jnc_utf8_rune_to_bytes(const i32_jnt &_R) noexcept;

constexpr u8_jnt __jnc_utf8_first[256] = {
    __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS,
    __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS,
    __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS,
    __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS,
    __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS,
    __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS,
    __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS,
    __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS,
    __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS,
    __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS,
    __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS,
    __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS,
    __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS,
    __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS,
    __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS,
    __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS,
    __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS,
    __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS,
    __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS,
    __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS,
    __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS,
    __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS,
    __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS,
    __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS,
    __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS,
    __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_AS, __JNC_UTF8_XX, __JNC_UTF8_XX,
    __JNC_UTF8_XX, __JNC_UTF8_XX, __JNC_UTF8_XX, __JNC_UTF8_XX, __JNC_UTF8_XX,
    __JNC_UTF8_XX, __JNC_UTF8_XX, __JNC_UTF8_XX, __JNC_UTF8_XX, __JNC_UTF8_XX,
    __JNC_UTF8_XX, __JNC_UTF8_XX, __JNC_UTF8_XX, __JNC_UTF8_XX, __JNC_UTF8_XX,
    __JNC_UTF8_XX, __JNC_UTF8_XX, __JNC_UTF8_XX, __JNC_UTF8_XX, __JNC_UTF8_XX,
    __JNC_UTF8_XX, __JNC_UTF8_XX, __JNC_UTF8_XX, __JNC_UTF8_XX, __JNC_UTF8_XX,
    __JNC_UTF8_XX, __JNC_UTF8_XX, __JNC_UTF8_XX, __JNC_UTF8_XX, __JNC_UTF8_XX,
    __JNC_UTF8_XX, __JNC_UTF8_XX, __JNC_UTF8_XX, __JNC_UTF8_XX, __JNC_UTF8_XX,
    __JNC_UTF8_XX, __JNC_UTF8_XX, __JNC_UTF8_XX, __JNC_UTF8_XX, __JNC_UTF8_XX,
    __JNC_UTF8_XX, __JNC_UTF8_XX, __JNC_UTF8_XX, __JNC_UTF8_XX, __JNC_UTF8_XX,
    __JNC_UTF8_XX, __JNC_UTF8_XX, __JNC_UTF8_XX, __JNC_UTF8_XX, __JNC_UTF8_XX,
    __JNC_UTF8_XX, __JNC_UTF8_XX, __JNC_UTF8_XX, __JNC_UTF8_XX, __JNC_UTF8_XX,
    __JNC_UTF8_XX, __JNC_UTF8_XX, __JNC_UTF8_XX, __JNC_UTF8_XX, __JNC_UTF8_XX,
    __JNC_UTF8_XX, __JNC_UTF8_XX, __JNC_UTF8_XX, __JNC_UTF8_XX, __JNC_UTF8_S1,
    __JNC_UTF8_S1, __JNC_UTF8_S1, __JNC_UTF8_S1, __JNC_UTF8_S1, __JNC_UTF8_S1,
    __JNC_UTF8_S1, __JNC_UTF8_S1, __JNC_UTF8_S1, __JNC_UTF8_S1, __JNC_UTF8_S1,
    __JNC_UTF8_S1, __JNC_UTF8_S1, __JNC_UTF8_S1, __JNC_UTF8_S1, __JNC_UTF8_S1,
    __JNC_UTF8_S1, __JNC_UTF8_S1, __JNC_UTF8_S1, __JNC_UTF8_S1, __JNC_UTF8_S1,
    __JNC_UTF8_S1, __JNC_UTF8_S1, __JNC_UTF8_S1, __JNC_UTF8_S1, __JNC_UTF8_S1,
    __JNC_UTF8_S1, __JNC_UTF8_S1, __JNC_UTF8_S1, __JNC_UTF8_S1, __JNC_UTF8_S2,
    __JNC_UTF8_S3, __JNC_UTF8_S3, __JNC_UTF8_S3, __JNC_UTF8_S3, __JNC_UTF8_S3,
    __JNC_UTF8_S3, __JNC_UTF8_S3, __JNC_UTF8_S3, __JNC_UTF8_S3, __JNC_UTF8_S3,
    __JNC_UTF8_S3, __JNC_UTF8_S3, __JNC_UTF8_S4, __JNC_UTF8_S3, __JNC_UTF8_S3,
    __JNC_UTF8_S5, __JNC_UTF8_S6, __JNC_UTF8_S6, __JNC_UTF8_S6, __JNC_UTF8_S7,
    __JNC_UTF8_XX, __JNC_UTF8_XX, __JNC_UTF8_XX, __JNC_UTF8_XX, __JNC_UTF8_XX,
    __JNC_UTF8_XX, __JNC_UTF8_XX, __JNC_UTF8_XX, __JNC_UTF8_XX, __JNC_UTF8_XX,
    __JNC_UTF8_XX,
};

struct __jnc_utf8_accept_range {
  u8_jnt _lo, _hi;
};

const __jnc_utf8_accept_range __jnc_utf8_accept_ranges[16] = {
    {__JNC_UTF8_LOCB, __JNC_UTF8_HICB},
    {0xA0, __JNC_UTF8_HICB},
    {__JNC_UTF8_LOCB, 0x9F},
    {0x90, __JNC_UTF8_HICB},
    {__JNC_UTF8_LOCB, 0x8F},
};

std::tuple<i32_jnt, int_jnt>
__jnc_utf8_decode_rune_str(const char *_S) noexcept {
  const std::size_t _len{std::strlen(_S)};
  if (_len < 1) {
    return std::make_tuple(__JNC_UTF8_RUNE_ERROR, 0);
  }
  const u8_jnt _s0{(u8_jnt)(_S[0])};
  const u8_jnt _x{__jnc_utf8_first[_s0]};
  if (_x >= __JNC_UTF8_AS) {
    const i32_jnt mask{_x << 31 >> 31};
    return std::make_tuple(
        ((i32_jnt)(_S[0]) & ~mask) | (__JNC_UTF8_RUNE_ERROR & mask), 1);
  }
  const int_jnt sz{(int_jnt)(_x * 7)};
  const __jnc_utf8_accept_range _accept{__jnc_utf8_accept_ranges[_x >> 4]};
  if (_len < sz) {
    return std::make_tuple(__JNC_UTF8_RUNE_ERROR, 1);
  }
  const u8_jnt _s1{(u8_jnt)(_S[1])};
  if (_s1 < _accept._lo || _accept._hi < _s1) {
    return std::make_tuple(__JNC_UTF8_RUNE_ERROR, 1);
  }
  if (sz <= 2) {
    return std::make_tuple(((i32_jnt)(_s0 & __JNC_UTF8_MASK2) << 6) |
                               (i32_jnt)(_s1 & __JNC_UTF8_MASKX),
                           2);
  }
  const u8_jnt _s2{(u8_jnt)(_S[2])};
  if (_s2 < __JNC_UTF8_LOCB || __JNC_UTF8_HICB < _s2) {
    return std::make_tuple(__JNC_UTF8_RUNE_ERROR, 1);
  }
  if (sz <= 3) {
    return std::make_tuple(((i32_jnt)(_s0 & __JNC_UTF8_MASK3) << 12) |
                               ((i32_jnt)(_s1 & __JNC_UTF8_MASKX) << 6) |
                               (i32_jnt)(_s2 & __JNC_UTF8_MASKX),
                           3);
  }
  const u8_jnt _s3{(u8_jnt)(_S[3])};
  if (_s3 < __JNC_UTF8_LOCB || __JNC_UTF8_HICB < _s3) {
    return std::make_tuple(__JNC_UTF8_RUNE_ERROR, 1);
  }
  return std::make_tuple(((i32_jnt)(_s0 & __JNC_UTF8_MASK4) << 18) |
                             ((i32_jnt)(_s1 & __JNC_UTF8_MASKX) << 12) |
                             ((i32_jnt)(_s2 & __JNC_UTF8_MASKX) << 6) |
                             (i32_jnt)(_s3 & __JNC_UTF8_MASKX),
                         4);
}

slice<u8_jnt> __jnc_utf8_rune_to_bytes(const i32_jnt &_R) noexcept {
  if ((u32_jnt)(_R) <= __JNC_UTF8_RUNE1_MAX) {
    return slice<u8_jnt>({(u8_jnt)(_R)});
  }
  const u32_jnt _i{(u32_jnt)(_R)};
  if (_i < __JNC_UTF8_RUNE2_MAX) {
    return slice<u8_jnt>(
        {(u8_jnt)(__JNC_UTF8_T2 | (u8_jnt)(_R >> 6)),
         (u8_jnt)(__JNC_UTF8_TX | ((u8_jnt)(_R)&__JNC_UTF8_MASKX))});
  }
  i32_jnt _r{_R};
  if (_i > __JNC_UTF8_MAX_RUNE,
      __JNC_UTF8_SURROGATE_MIN <= _i && _i <= __JNC_UTF8_SURROGATE_MAX) {
    _r = __JNC_UTF8_RUNE_ERROR;
  }
  if (_i <= __JNC_UTF8_RUNE3_MAX) {
    return slice<u8_jnt>(
        {(u8_jnt)(__JNC_UTF8_T3 | (u8_jnt)(_r >> 12)),
         (u8_jnt)(__JNC_UTF8_TX | ((u8_jnt)(_r >> 6) & __JNC_UTF8_MASKX)),
         (u8_jnt)(__JNC_UTF8_TX | ((u8_jnt)(_r)&__JNC_UTF8_MASKX))});
  }
  return slice<u8_jnt>(
      {(u8_jnt)(__JNC_UTF8_T4 | (u8_jnt)(_r >> 18)),
       (u8_jnt)(__JNC_UTF8_TX | ((u8_jnt)(_r >> 12) & __JNC_UTF8_MASKX)),
       (u8_jnt)(__JNC_UTF8_TX | ((u8_jnt)(_r >> 6) & __JNC_UTF8_MASKX)),
       (u8_jnt)(__JNC_UTF8_TX | ((u8_jnt)(_r)&__JNC_UTF8_MASKX))});
}

#endif // !__JNC_UTF8_HPP
