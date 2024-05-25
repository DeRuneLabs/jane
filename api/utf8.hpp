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

#ifndef __JANE_UTF8_HPP
#define __JANE_UTF8_HPP

#include "slice.hpp"
#include "types.hpp"
#include <tuple>

namespace jane {
constexpr signed int UTF8_RUNE_ERROR{65533};
constexpr signed int UTF8_MASKX{63};
constexpr signed int UTF8_MASK2{31};
constexpr signed int UTF8_MASK3{15};
constexpr signed int UTF8_MASK4{7};
constexpr signed int UTF8_LOCB{128};
constexpr signed int UTF8_HICB{191};
constexpr signed int UTF8_XX{241};
constexpr signed int UTF8_AS{240};
constexpr signed int UTF8_S1{2};
constexpr signed int UTF8_S2{19};
constexpr signed int UTF8_S3{3};
constexpr signed int UTF8_S4{35};
constexpr signed int UTF8_S5{52};
constexpr signed int UTF8_S6{4};
constexpr signed int UTF8_S7{68};
constexpr signed int UTF8_RUNE1_MAX{127};
constexpr signed int UTF8_RUNE2_MAX{2047};
constexpr signed int UTF8_RUNE3_MAX{65535};
constexpr signed int UTF8_TX{128};
constexpr signed int UTF8_T2{192};
constexpr signed int UTF8_T3{224};
constexpr signed int UTF8_T4{240};
constexpr signed int UTF8_MAX_RUNE{1114111};
constexpr signed int UTF8_SURROGATE_MIN{55296};
constexpr signed int UTF8_SURROGATE_MAX{57343};

struct UTF8AcceptRange;
std::tuple<jane::I32, jane::Int>
utf8_decode_rune_str(const char *s, const jane::Int &len) noexcept;
jane::Slice<jane::U8> utf8_rune_to_bytes(const jane::I32 &r) noexcept;

constexpr jane::U8 utf8_first[256] = {
    jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS,
    jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS,
    jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS,
    jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS,
    jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS,
    jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS,
    jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS,
    jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS,
    jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS,
    jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS,
    jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS,
    jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS,
    jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS,
    jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS,
    jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS,
    jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS,
    jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS,
    jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS,
    jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS,
    jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS,
    jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS,
    jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS,
    jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS,
    jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS,
    jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS,
    jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_AS, jane::UTF8_XX, jane::UTF8_XX,
    jane::UTF8_XX, jane::UTF8_XX, jane::UTF8_XX, jane::UTF8_XX, jane::UTF8_XX,
    jane::UTF8_XX, jane::UTF8_XX, jane::UTF8_XX, jane::UTF8_XX, jane::UTF8_XX,
    jane::UTF8_XX, jane::UTF8_XX, jane::UTF8_XX, jane::UTF8_XX, jane::UTF8_XX,
    jane::UTF8_XX, jane::UTF8_XX, jane::UTF8_XX, jane::UTF8_XX, jane::UTF8_XX,
    jane::UTF8_XX, jane::UTF8_XX, jane::UTF8_XX, jane::UTF8_XX, jane::UTF8_XX,
    jane::UTF8_XX, jane::UTF8_XX, jane::UTF8_XX, jane::UTF8_XX, jane::UTF8_XX,
    jane::UTF8_XX, jane::UTF8_XX, jane::UTF8_XX, jane::UTF8_XX, jane::UTF8_XX,
    jane::UTF8_XX, jane::UTF8_XX, jane::UTF8_XX, jane::UTF8_XX, jane::UTF8_XX,
    jane::UTF8_XX, jane::UTF8_XX, jane::UTF8_XX, jane::UTF8_XX, jane::UTF8_XX,
    jane::UTF8_XX, jane::UTF8_XX, jane::UTF8_XX, jane::UTF8_XX, jane::UTF8_XX,
    jane::UTF8_XX, jane::UTF8_XX, jane::UTF8_XX, jane::UTF8_XX, jane::UTF8_XX,
    jane::UTF8_XX, jane::UTF8_XX, jane::UTF8_XX, jane::UTF8_XX, jane::UTF8_XX,
    jane::UTF8_XX, jane::UTF8_XX, jane::UTF8_XX, jane::UTF8_XX, jane::UTF8_S1,
    jane::UTF8_S1, jane::UTF8_S1, jane::UTF8_S1, jane::UTF8_S1, jane::UTF8_S1,
    jane::UTF8_S1, jane::UTF8_S1, jane::UTF8_S1, jane::UTF8_S1, jane::UTF8_S1,
    jane::UTF8_S1, jane::UTF8_S1, jane::UTF8_S1, jane::UTF8_S1, jane::UTF8_S1,
    jane::UTF8_S1, jane::UTF8_S1, jane::UTF8_S1, jane::UTF8_S1, jane::UTF8_S1,
    jane::UTF8_S1, jane::UTF8_S1, jane::UTF8_S1, jane::UTF8_S1, jane::UTF8_S1,
    jane::UTF8_S1, jane::UTF8_S1, jane::UTF8_S1, jane::UTF8_S1, jane::UTF8_S2,
    jane::UTF8_S3, jane::UTF8_S3, jane::UTF8_S3, jane::UTF8_S3, jane::UTF8_S3,
    jane::UTF8_S3, jane::UTF8_S3, jane::UTF8_S3, jane::UTF8_S3, jane::UTF8_S3,
    jane::UTF8_S3, jane::UTF8_S3, jane::UTF8_S4, jane::UTF8_S3, jane::UTF8_S3,
    jane::UTF8_S5, jane::UTF8_S6, jane::UTF8_S6, jane::UTF8_S6, jane::UTF8_S7,
    jane::UTF8_XX, jane::UTF8_XX, jane::UTF8_XX, jane::UTF8_XX, jane::UTF8_XX,
    jane::UTF8_XX, jane::UTF8_XX, jane::UTF8_XX, jane::UTF8_XX, jane::UTF8_XX,
    jane::UTF8_XX,
};

struct UTF8AcceptRange {
  const jane::U8 lo, hi;
};

constexpr struct jane::UTF8AcceptRange utf8_accept_ranges[16] = {
    {jane::UTF8_LOCB, jane::UTF8_HICB},
    {0xA0, jane::UTF8_HICB},
    {jane::UTF8_LOCB, 0x9F},
    {0x90, jane::UTF8_HICB},
    {jane::UTF8_LOCB, 0x8F},
};

inline std::tuple<jane::I32, jane::Int>
utf8_decode_rune_st(const char *s, const jane::Int &len) noexcept {
  if (len < 1) {
    return std::make_tuple<jane::I32, jane::Int>(jane::UTF8_RUNE_ERROR, 0);
  }
  const jane::U8 s0{static_cast<jane::U8>(s[0])};
  const jane::U8 x{jane::utf8_first[s0]};
  if (x >= jane::UTF8_AS) {
    const jane::I32 mask{x << 31 >> 31};
    return std::make_tuple((static_cast<jane::I32>(s[0]) & ~mask) |
                               (jane::UTF8_RUNE_ERROR & mask),
                           1);
  }
  const jane::Int sz{static_cast<jane::Int>(x & 7)};
  const struct jane::UTF8AcceptRange accept {
    jane::utf8_accept_ranges[x >> 4]
  };
  if (len < sz) {
    return std::make_tuple<jane::I32, jane::Int>(jane::UTF8_RUNE_ERROR, 1);
  }
  const jane::U8 s1{static_cast<jane::U8>(s[1])};
  if (s1 < accept.lo || accept.hi < s1) {
    return std::make_tuple<jane::I32, jane::Int>(jane::UTF8_RUNE_ERROR, 1);
  }

  if (sz <= 2) {
    return std::make_tuple<jane::I32, jane::Int>(
        (static_cast<jane::I32>(s0 & jane::UTF8_MASK2) << 6) |
            static_cast<jane::I32>(s1 & jane::UTF8_MASKX),
        2);
  }

  const jane::U8 s2{static_cast<jane::U8>(s[2])};
  if (s2 < jane::UTF8_LOCB || jane::UTF8_HICB < s2) {
    return std::make_tuple<jane::I32, jane::Int>(jane::UTF8_RUNE_ERROR, 1);
  }

  if (sz <= 3) {
    return std::make_tuple<jane::I32, jane::Int>(
        (static_cast<jane::I32>(s0 & jane::UTF8_MASK3) << 12) |
            (static_cast<jane::I32>(s1 & jane::UTF8_MASKX) << 6) |
            static_cast<jane::I32>(s2 & jane::UTF8_MASKX),
        3);
  }

  const jane::U8 s3{static_cast<jane::U8>(s[3])};
  if (s3 < jane::UTF8_LOCB || jane::UTF8_HICB < s3) {
    return std::make_tuple<jane::I32, jane::Int>(jane::UTF8_RUNE_ERROR, 1);
  }

  return std::make_tuple(
      (static_cast<jane::I32>(s0 & jane::UTF8_MASK4) << 18) |
          (static_cast<jane::I32>(s1 & jane::UTF8_MASKX) << 12) |
          (static_cast<jane::I32>(s2 & jane::UTF8_MASKX) << 6) |
          static_cast<jane::I32>(s3 & jane::UTF8_MASKX),
      4);
}

inline jane::Slice<jane::U8> utf8_rune_to_bytes(const jane::I32 &r) noexcept {
  if (static_cast<jane::U32>(r) <= jane::UTF8_RUNE1_MAX) {
    return jane::Slice<jane::U8>({static_cast<jane::U8>(r)});
  }

  const jane::U32 i{static_cast<jane::U32>(r)};
  if (i < jane::UTF8_RUNE2_MAX) {
    return jane::Slice<jane::U8>(
        {static_cast<jane::U8>(jane::UTF8_T2 | static_cast<jane::U8>(r >> 6)),
         static_cast<jane::U8>(jane::UTF8_TX |
                               (static_cast<jane::U8>(r) & jane::UTF8_MASKX))});
  }

  jane::I32 _r{r};
  if (i > jane::UTF8_MAX_RUNE ||
      jane::UTF8_SURROGATE_MIN <= i && i <= jane::UTF8_SURROGATE_MAX)
    _r = jane::UTF8_RUNE_ERROR;

  if (i <= jane::UTF8_RUNE3_MAX)
    return jane::Slice<jane::U8>(
        {static_cast<jane::U8>(jane::UTF8_T3 | static_cast<jane::U8>(_r >> 12)),
         static_cast<jane::U8>(jane::UTF8_TX | (static_cast<jane::U8>(_r >> 6) &
                                                jane::UTF8_MASKX)),
         static_cast<jane::U8>(
             jane::UTF8_TX | (static_cast<jane::U8>(_r) & jane::UTF8_MASKX))});

  return jane::Slice<jane::U8>(
      {static_cast<jane::U8>(jane::UTF8_T4 | static_cast<jane::U8>(_r >> 18)),
       static_cast<jane::U8>(jane::UTF8_TX | (static_cast<jane::U8>(_r >> 12) &
                                              jane::UTF8_MASKX)),
       static_cast<jane::U8>(
           jane::UTF8_TX | (static_cast<jane::U8>(_r >> 6) & jane::UTF8_MASKX)),
       static_cast<jane::U8>(jane::UTF8_TX |
                             (static_cast<jane::U8>(_r) & jane::UTF8_MASKX))});
}
} // namespace jane

#endif // __JANE_UTF8_HPP