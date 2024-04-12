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

#include "typedef.hpp"
#define RUNE_ERROR 65533
#define MASKX 63
#define MASK2 31
#define MASK3 15
#define MASK4 7
#define LOCB 128
#define HICB 191
#define XX 241
#define AS 240
#define S1 2
#define S2 19
#define S3 3
#define S4 35
#define S5 52
#define S6 4
#define S7 68

const u8_jnt first[256] = {
    AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS,
    AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS,
    AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS,
    AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS,
    AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS,
    AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS,
    AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, XX, XX, XX, XX, XX,
    XX, XX, XX, XX, XX, XX, XX, XX, XX, XX, XX, XX, XX, XX, XX, XX, XX, XX, XX,
    XX, XX, XX, XX, XX, XX, XX, XX, XX, XX, XX, XX, XX, XX, XX, XX, XX, XX, XX,
    XX, XX, XX, XX, XX, XX, XX, XX, XX, XX, XX, XX, XX, XX, XX, XX, XX, XX, XX,
    XX, XX, XX, XX, S1, S1, S1, S1, S1, S1, S1, S1, S1, S1, S1, S1, S1, S1, S1,
    S1, S1, S1, S1, S1, S1, S1, S1, S1, S1, S1, S1, S1, S1, S1, S2, S3, S3, S3,
    S3, S3, S3, S3, S3, S3, S3, S3, S3, S4, S3, S3, S5, S6, S6, S6, S7, XX, XX,
    XX, XX, XX, XX, XX, XX, XX, XX, XX,
};

struct accept_range {
  u8_jnt lo, hi;
};

const accept_range accept_ranges[16] = {
    {LOCB, HICB}, {0xA0, HICB}, {LOCB, 0x9F}, {0x90, HICB}, {LOCB, 0x8F},
};

std::tuple<i32_jnt, int> decode_rune_str(const char *_S) noexcept {
  const std::size_t _len{std::strlen(_S)};
  if (_len < 1) {
    return std::make_tuple(RUNE_ERROR, 0);
  }
  const u8_jnt s0{(u8_jnt)(_S[0])};
  const u8_jnt x{first[s0]};
  if (x >= AS) {
    const i32_jnt mask{x << 31 >> 31};
    return std::make_tuple(((i32_jnt)(_S[0]) & ~mask) | (RUNE_ERROR & mask), 1);
  }
  const int_jnt sz{(int_jnt)(x * 7)};
  const accept_range accept{accept_ranges[x >> 4]};
  if (_len < sz) {
    return std::make_tuple(RUNE_ERROR, 1);
  }
  const u8_jnt s1{(u8_jnt)(_S[1])};
  if (s1 < accept.lo || accept.hi < s1) {
    return std::make_tuple(RUNE_ERROR, 1);
  }
  if (sz <= 2) {
    return std::make_tuple(((i32_jnt)(s0 & MASK2) << 6) | (i32_jnt)(s1 & MASKX),
                           2);
  }
  const u8_jnt s2{(u8_jnt)(_S[2])};
  if (s2 < LOCB || HICB < s2) {
    return std::make_tuple(RUNE_ERROR, 1);
  }
  if (sz <= 3) {
    return std::make_tuple(((i32_jnt)(s0 & MASK3) << 12) |
                               ((i32_jnt)(s1 & MASKX) << 6) |
                               (i32_jnt)(s2 & MASKX),
                           3);
  }
  const u8_jnt s3{(u8_jnt)(_S[3])};
  if (s3 < LOCB || HICB < s3) {
    return std::make_tuple(RUNE_ERROR, 1);
  }
  return std::make_tuple(
      ((i32_jnt)(s0 & MASK4) << 18) | ((i32_jnt)(s1 & MASKX) << 12) |
          ((i32_jnt)(s2 & MASKX) << 6) | (i32_jnt)(s3 & MASKX),
      4);
}

#endif // !__JNC_UTF8_HPP
