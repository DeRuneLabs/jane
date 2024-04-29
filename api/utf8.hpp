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

#ifndef __JANE_UTF8_HPP
#define __JANE_UTF8_HPP

#include "slice.hpp"
#include "typedef.hpp"
constexpr signed int __JANE_UTF8_RUNE_ERROR{65533};
constexpr signed int __JANE_UTF8_MASKX{63};
constexpr signed int __JANE_UTF8_MASK2{31};
constexpr signed int __JANE_UTF8_MASK3{15};
constexpr signed int __JANE_UTF8_MASK4{7};
constexpr signed int __JANE_UTF8_LOCB{128};
constexpr signed int __JANE_UTF8_HICB{191};
constexpr signed int __JANE_UTF8_XX{241};
constexpr signed int __JANE_UTF8_AS{240};
constexpr signed int __JANE_UTF8_S1{2};
constexpr signed int __JANE_UTF8_S2{19};
constexpr signed int __JANE_UTF8_S3{3};
constexpr signed int __JANE_UTF8_S4{35};
constexpr signed int __JANE_UTF8_S5{52};
constexpr signed int __JANE_UTF8_S6{4};
constexpr signed int __JANE_UTF8_S7{68};
constexpr signed int __JANE_UTF8_RUNE1_MAX{127};
constexpr signed int __JANE_UTF8_RUNE2_MAX{2047};
constexpr signed int __JANE_UTF8_RUNE3_MAX{65535};
constexpr signed int __JANE_UTF8_TX{128};
constexpr signed int __JANE_UTF8_T2{192};
constexpr signed int __JANE_UTF8_T3{224};
constexpr signed int __JANE_UTF8_T4{240};
constexpr signed int __JANE_UTF8_MAX_RUNE{1114111};
constexpr signed int __JANE_UTF8_SURROGATE_MIN{55296};
constexpr signed int __JANE_UTF8_SURROGATE_MAX{57343};

struct __jane_utf8_accept_range;
std::tuple<i32_jnt, int_jnt>
__jane_utf8_decode_rune_str(const char *_S, const int_jnt &_Len) noexcept;
slice_jnt<u8_jnt> __jane_utf8_rune_to_bytes(const i32_jnt &_R) noexcept;

constexpr u8_jnt __jane_utf8_first[256] = {
    __JANE_UTF8_AS, __JANE_UTF8_AS, __JANE_UTF8_AS, __JANE_UTF8_AS,
    __JANE_UTF8_AS, __JANE_UTF8_AS, __JANE_UTF8_AS, __JANE_UTF8_AS,
    __JANE_UTF8_AS, __JANE_UTF8_AS, __JANE_UTF8_AS, __JANE_UTF8_AS,
    __JANE_UTF8_AS, __JANE_UTF8_AS, __JANE_UTF8_AS, __JANE_UTF8_AS,
    __JANE_UTF8_AS, __JANE_UTF8_AS, __JANE_UTF8_AS, __JANE_UTF8_AS,
    __JANE_UTF8_AS, __JANE_UTF8_AS, __JANE_UTF8_AS, __JANE_UTF8_AS,
    __JANE_UTF8_AS, __JANE_UTF8_AS, __JANE_UTF8_AS, __JANE_UTF8_AS,
    __JANE_UTF8_AS, __JANE_UTF8_AS, __JANE_UTF8_AS, __JANE_UTF8_AS,
    __JANE_UTF8_AS, __JANE_UTF8_AS, __JANE_UTF8_AS, __JANE_UTF8_AS,
    __JANE_UTF8_AS, __JANE_UTF8_AS, __JANE_UTF8_AS, __JANE_UTF8_AS,
    __JANE_UTF8_AS, __JANE_UTF8_AS, __JANE_UTF8_AS, __JANE_UTF8_AS,
    __JANE_UTF8_AS, __JANE_UTF8_AS, __JANE_UTF8_AS, __JANE_UTF8_AS,
    __JANE_UTF8_AS, __JANE_UTF8_AS, __JANE_UTF8_AS, __JANE_UTF8_AS,
    __JANE_UTF8_AS, __JANE_UTF8_AS, __JANE_UTF8_AS, __JANE_UTF8_AS,
    __JANE_UTF8_AS, __JANE_UTF8_AS, __JANE_UTF8_AS, __JANE_UTF8_AS,
    __JANE_UTF8_AS, __JANE_UTF8_AS, __JANE_UTF8_AS, __JANE_UTF8_AS,
    __JANE_UTF8_AS, __JANE_UTF8_AS, __JANE_UTF8_AS, __JANE_UTF8_AS,
    __JANE_UTF8_AS, __JANE_UTF8_AS, __JANE_UTF8_AS, __JANE_UTF8_AS,
    __JANE_UTF8_AS, __JANE_UTF8_AS, __JANE_UTF8_AS, __JANE_UTF8_AS,
    __JANE_UTF8_AS, __JANE_UTF8_AS, __JANE_UTF8_AS, __JANE_UTF8_AS,
    __JANE_UTF8_AS, __JANE_UTF8_AS, __JANE_UTF8_AS, __JANE_UTF8_AS,
    __JANE_UTF8_AS, __JANE_UTF8_AS, __JANE_UTF8_AS, __JANE_UTF8_AS,
    __JANE_UTF8_AS, __JANE_UTF8_AS, __JANE_UTF8_AS, __JANE_UTF8_AS,
    __JANE_UTF8_AS, __JANE_UTF8_AS, __JANE_UTF8_AS, __JANE_UTF8_AS,
    __JANE_UTF8_AS, __JANE_UTF8_AS, __JANE_UTF8_AS, __JANE_UTF8_AS,
    __JANE_UTF8_AS, __JANE_UTF8_AS, __JANE_UTF8_AS, __JANE_UTF8_AS,
    __JANE_UTF8_AS, __JANE_UTF8_AS, __JANE_UTF8_AS, __JANE_UTF8_AS,
    __JANE_UTF8_AS, __JANE_UTF8_AS, __JANE_UTF8_AS, __JANE_UTF8_AS,
    __JANE_UTF8_AS, __JANE_UTF8_AS, __JANE_UTF8_AS, __JANE_UTF8_AS,
    __JANE_UTF8_AS, __JANE_UTF8_AS, __JANE_UTF8_AS, __JANE_UTF8_AS,
    __JANE_UTF8_AS, __JANE_UTF8_AS, __JANE_UTF8_AS, __JANE_UTF8_AS,
    __JANE_UTF8_AS, __JANE_UTF8_AS, __JANE_UTF8_AS, __JANE_UTF8_AS,
    __JANE_UTF8_XX, __JANE_UTF8_XX, __JANE_UTF8_XX, __JANE_UTF8_XX,
    __JANE_UTF8_XX, __JANE_UTF8_XX, __JANE_UTF8_XX, __JANE_UTF8_XX,
    __JANE_UTF8_XX, __JANE_UTF8_XX, __JANE_UTF8_XX, __JANE_UTF8_XX,
    __JANE_UTF8_XX, __JANE_UTF8_XX, __JANE_UTF8_XX, __JANE_UTF8_XX,
    __JANE_UTF8_XX, __JANE_UTF8_XX, __JANE_UTF8_XX, __JANE_UTF8_XX,
    __JANE_UTF8_XX, __JANE_UTF8_XX, __JANE_UTF8_XX, __JANE_UTF8_XX,
    __JANE_UTF8_XX, __JANE_UTF8_XX, __JANE_UTF8_XX, __JANE_UTF8_XX,
    __JANE_UTF8_XX, __JANE_UTF8_XX, __JANE_UTF8_XX, __JANE_UTF8_XX,
    __JANE_UTF8_XX, __JANE_UTF8_XX, __JANE_UTF8_XX, __JANE_UTF8_XX,
    __JANE_UTF8_XX, __JANE_UTF8_XX, __JANE_UTF8_XX, __JANE_UTF8_XX,
    __JANE_UTF8_XX, __JANE_UTF8_XX, __JANE_UTF8_XX, __JANE_UTF8_XX,
    __JANE_UTF8_XX, __JANE_UTF8_XX, __JANE_UTF8_XX, __JANE_UTF8_XX,
    __JANE_UTF8_XX, __JANE_UTF8_XX, __JANE_UTF8_XX, __JANE_UTF8_XX,
    __JANE_UTF8_XX, __JANE_UTF8_XX, __JANE_UTF8_XX, __JANE_UTF8_XX,
    __JANE_UTF8_XX, __JANE_UTF8_XX, __JANE_UTF8_XX, __JANE_UTF8_XX,
    __JANE_UTF8_XX, __JANE_UTF8_XX, __JANE_UTF8_XX, __JANE_UTF8_XX,
    __JANE_UTF8_XX, __JANE_UTF8_XX, __JANE_UTF8_S1, __JANE_UTF8_S1,
    __JANE_UTF8_S1, __JANE_UTF8_S1, __JANE_UTF8_S1, __JANE_UTF8_S1,
    __JANE_UTF8_S1, __JANE_UTF8_S1, __JANE_UTF8_S1, __JANE_UTF8_S1,
    __JANE_UTF8_S1, __JANE_UTF8_S1, __JANE_UTF8_S1, __JANE_UTF8_S1,
    __JANE_UTF8_S1, __JANE_UTF8_S1, __JANE_UTF8_S1, __JANE_UTF8_S1,
    __JANE_UTF8_S1, __JANE_UTF8_S1, __JANE_UTF8_S1, __JANE_UTF8_S1,
    __JANE_UTF8_S1, __JANE_UTF8_S1, __JANE_UTF8_S1, __JANE_UTF8_S1,
    __JANE_UTF8_S1, __JANE_UTF8_S1, __JANE_UTF8_S1, __JANE_UTF8_S1,
    __JANE_UTF8_S2, __JANE_UTF8_S3, __JANE_UTF8_S3, __JANE_UTF8_S3,
    __JANE_UTF8_S3, __JANE_UTF8_S3, __JANE_UTF8_S3, __JANE_UTF8_S3,
    __JANE_UTF8_S3, __JANE_UTF8_S3, __JANE_UTF8_S3, __JANE_UTF8_S3,
    __JANE_UTF8_S3, __JANE_UTF8_S4, __JANE_UTF8_S3, __JANE_UTF8_S3,
    __JANE_UTF8_S5, __JANE_UTF8_S6, __JANE_UTF8_S6, __JANE_UTF8_S6,
    __JANE_UTF8_S7, __JANE_UTF8_XX, __JANE_UTF8_XX, __JANE_UTF8_XX,
    __JANE_UTF8_XX, __JANE_UTF8_XX, __JANE_UTF8_XX, __JANE_UTF8_XX,
    __JANE_UTF8_XX, __JANE_UTF8_XX, __JANE_UTF8_XX, __JANE_UTF8_XX,
};

struct __jane_utf8_accept_range {
  const u8_jnt _lo, _hi;
};

constexpr struct __jane_utf8_accept_range __jane_utf8_accept_ranges[16] = {
    {__JANE_UTF8_LOCB, __JANE_UTF8_HICB},
    {0xA0, __JANE_UTF8_HICB},
    {__JANE_UTF8_LOCB, 0x9F},
    {0x90, __JANE_UTF8_HICB},
    {__JANE_UTF8_LOCB, 0x8F},
};

std::tuple<i32_jnt, int_jnt>
__jane_utf8_decode_rune_str(const char *_S, const int_jnt &_Len) noexcept {
  if (_Len < 1) {
    return (std::make_tuple(__JANE_UTF8_RUNE_ERROR, 0));
  }
  const u8_jnt _s0{static_cast<u8_jnt>(_S[0])};
  const u8_jnt _x{__jane_utf8_first[_s0]};
  if (_x >= __JANE_UTF8_AS) {
    const i32_jnt _mask{_x << 31 >> 31};
    return (std::make_tuple((static_cast<i32_jnt>(_S[0]) & ~_mask) |
                                (__JANE_UTF8_RUNE_ERROR & _mask),
                            1));
  }
  const int_jnt _sz{static_cast<int_jnt>(_x & 7)};
  const struct __jane_utf8_accept_range _accept {
    __jane_utf8_accept_ranges[_x >> 4]
  };
  if (_Len < _sz) {
    return (std::make_tuple(__JANE_UTF8_RUNE_ERROR, 1));
  }
  const u8_jnt _s1{static_cast<u8_jnt>(_S[1])};
  if (_s1 < _accept._lo || _accept._hi < _s1) {
    return (std::make_tuple(__JANE_UTF8_RUNE_ERROR, 1));
  }
  if (_sz <= 2) {
    return (
        std::make_tuple((static_cast<i32_jnt>(_s0 & __JANE_UTF8_MASK2) << 6) |
                            static_cast<i32_jnt>(_s1 & __JANE_UTF8_MASKX),
                        2));
  }
  const u8_jnt _s2{static_cast<u8_jnt>(_S[2])};
  if (_s2 < __JANE_UTF8_LOCB || __JANE_UTF8_HICB < _s2) {
    return (std::make_tuple(__JANE_UTF8_RUNE_ERROR, 1));
  }
  if (_sz <= 3) {
    return (std::make_tuple(
        (static_cast<i32_jnt>(_s0 & __JANE_UTF8_MASK3) << 12) |
            (static_cast<i32_jnt>(_s1 & __JANE_UTF8_MASKX) << 6) |
            static_cast<i32_jnt>(_s2 & __JANE_UTF8_MASKX),
        3));
  }
  const u8_jnt _s3{static_cast<u8_jnt>(_S[3])};
  if (_s3 < __JANE_UTF8_LOCB || __JANE_UTF8_HICB < _s3) {
    return std::make_tuple(__JANE_UTF8_RUNE_ERROR, 1);
  }
  return (std::make_tuple(
      (static_cast<i32_jnt>(_s0 & __JANE_UTF8_MASK4) << 18) |
          (static_cast<i32_jnt>(_s1 & __JANE_UTF8_MASKX) << 12) |
          (static_cast<i32_jnt>(_s2 & __JANE_UTF8_MASKX) << 6) |
          static_cast<i32_jnt>(_s3 & __JANE_UTF8_MASKX),
      4));
}

slice_jnt<u8_jnt> __jane_utf8_rune_to_bytes(const i32_jnt &_R) noexcept {
  if (static_cast<u32_jnt>(_R) <= __JANE_UTF8_RUNE1_MAX) {
    return (slice_jnt<u8_jnt>({static_cast<u8_jnt>(_R)}));
  }
  const u32_jnt _i{static_cast<u32_jnt>(_R)};
  if (_i < __JANE_UTF8_RUNE2_MAX) {
    return (slice_jnt<u8_jnt>(
        {static_cast<u8_jnt>(__JANE_UTF8_T2 | static_cast<u8_jnt>(_R >> 6)),
         static_cast<u8_jnt>(__JANE_UTF8_TX |
                             (static_cast<u8_jnt>(_R) & __JANE_UTF8_MASKX))}));
  }
  i32_jnt _r{_R};
  if ((_i > __JANE_UTF8_MAX_RUNE) ||
      (__JANE_UTF8_SURROGATE_MIN <= _i && _i <= __JANE_UTF8_SURROGATE_MAX)) {
    _r = __JANE_UTF8_RUNE_ERROR;
  }
  if (_i <= __JANE_UTF8_RUNE3_MAX) {
    return (slice_jnt<u8_jnt>(
        {static_cast<u8_jnt>(__JANE_UTF8_T3 | static_cast<u8_jnt>(_r >> 12)),
         static_cast<u8_jnt>(__JANE_UTF8_TX | (static_cast<u8_jnt>(_r >> 6) &
                                               __JANE_UTF8_MASKX)),
         static_cast<u8_jnt>(__JANE_UTF8_TX |
                             (static_cast<u8_jnt>(_r) & __JANE_UTF8_MASKX))}));
  }
  return (slice_jnt<u8_jnt>(
      {static_cast<u8_jnt>(__JANE_UTF8_T4 | static_cast<u8_jnt>(_r >> 18)),
       static_cast<u8_jnt>(__JANE_UTF8_TX |
                           (static_cast<u8_jnt>(_r >> 12) & __JANE_UTF8_MASKX)),
       static_cast<u8_jnt>(__JANE_UTF8_TX |
                           (static_cast<u8_jnt>(_r >> 6) & __JANE_UTF8_MASKX)),
       static_cast<u8_jnt>(__JANE_UTF8_TX |
                           (static_cast<u8_jnt>(_r) & __JANE_UTF8_MASKX))}));
}

#endif // !__JANE_UTF8_HPP
