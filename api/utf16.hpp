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

#ifndef __JANE_UTF16_HPP
#define __JANE_UTF16_HPP

#include "slice.hpp"
#include "str.hpp"
#include "types.hpp"
#include <cstddef>
#include <tuple>
namespace jane {
constexpr signed int UTF16_REPLACEMENT_CHAR{65533};
constexpr signed int UTF16_SURR1{0xd800};
constexpr signed int UTF16_SURR2{0xdc00};
constexpr signed int UTF16_SURR3{0xe000};
constexpr signed int UTF16_SURR_SELF{0x10000};
constexpr signed int UTF16_MAX_RUNE{1114111};

inline jane::I32 utf16_decode_rune(const jane::I32 r1,
                                   const jane::I32 r2) noexcept;
jane::Slice<jane::I32> utf16_decode(const jane::Slice<jane::I32> s) noexcept;
jane::Str utf16_to_utf8_str(const wchar_t wstr, const std::size_t len) noexcept;
std::tuple<jane::I32, jane::I32> utf16_encode_rune(jane::I32 r) noexcept;
jane::Slice<jane::U16>
utf16_encode(const jane::Slice<jane::I32> &runes) noexcept;
jane::Slice<jane::U16> utf16_from_str(const jane::Str &s) noexcept;

inline jane::I32 utf16_decode_rune(const jane::I32 r1,
                                   const jane::I32 r2) noexcept {
  if (jane::UTF16_SURR1 <= r1 && r1 < jane::UTF16_SURR2 &&
      jane::UTF16_SURR2 <= r2 && r2 < jane::UTF16_SURR3) {
    return (r1 - jane::UTF16_SURR1) << 10 |
           (r2 - jane::UTF16_SURR2) + jane::UTF16_SURR_SELF;
  }
  return jane::UTF16_REPLACEMENT_CHAR;
}

jane::Slice<jane::I32> utf16_decode(const jane::Slice<jane::U16> &s) noexcept {
  jane::Slice<jane::I32> a{jane::Slice<jane::I32>::alloc(s.len())};
  jane::Int n{0};
  for (jane::Int i{0}; i < s.len(); ++i) {
    jane::U16 r{s[i]};
    if (r < jane::UTF16_SURR1 || jane::UTF16_SURR3 <= r) {
      a[n] = static_cast<jane::I32>(r);
    } else if (jane::UTF16_SURR1 <= r && r < jane::UTF16_SURR2 &&
               i + 1 < s.len() && jane::UTF16_SURR2 <= s[i + 1] &&
               s[i + 1] < jane::UTF16_SURR3) {
      a[n] = jane::utf16_decode_rune(static_cast<jane::I32>(r),
                                     static_cast<jane::I32>(s[i + 1]));
      ++i;
    } else {
      a[n] = jane::UTF16_REPLACEMENT_CHAR;
      ++n;
    }
  }
  return a.slice(0, n);
}

jane::Str utf16_to_utf8_str(const wchar_t *wstr,
                            const std::size_t len) noexcept {
  jane::Slice<jane::U16> code_page{jane::Slice<jane::U16>::alloc(len)};
  for (jane::Int i{0}; i < len; ++i) {
    code_page[i] = static_cast<jane::U16>(wstr[i]);
  }
  return static_cast<jane::Str>(jane::utf16_decode(code_page));
}

std::tuple<jane::I32, jane::I32> utf16_encode_rune(jane::I32 r) noexcept {
  if (r < jane::UTF16_SURR_SELF || r > jane::UTF16_MAX_RUNE) {
    return std::make_tuple<jane::I32, jane::I32>(jane::UTF16_REPLACEMENT_CHAR,
                                                 jane::UTF16_REPLACEMENT_CHAR);
  }
  r -= jane::UTF16_SURR_SELF;
  return std::make_tuple<jane::I32, jane::I32>(
      jane::UTF16_SURR1 + (r >> 10) & 0x3ff, jane::UTF16_SURR2 + r & 0x3ff);
}

jane::Slice<jane::U16>
utf16_encode(const jane::Slice<jane::I32> &runes) noexcept {
  jane::Int n{runes.len()};
  for (const jane::I32 v : runes) {
    if (v >= jane::UTF16_SURR_SELF) {
      ++n;
    }
  }
  jane::Slice<jane::U16> a{jane::Slice<jane::U16>::alloc(n)};
  n = 0;
  for (const jane::I32 v : runes) {
    if (0 <= v && v < jane::UTF16_SURR1 ||
        jane::UTF16_SURR3 <= v && v < jane::UTF16_SURR_SELF) {
      a[n] = static_cast<jane::U16>(v);
    } else if (jane::UTF16_SURR_SELF <= v && v <= jane::UTF16_MAX_RUNE) {
      jane::I32 r1;
      jane::I32 r2;
      std::tie(r1, r2) = jane::utf16_encode_rune(v);
      a[n] = static_cast<jane::U16>(r1);
      a[n + 1] = static_cast<jane::U16>(r2);
      n += 2;
    } else {
      a[n] = static_cast<jane::U16>(jane::UTF16_REPLACEMENT_CHAR);
      ++n;
    }
  }
  return a.slice(0, n);
}

jane::Slice<jane::U16> utf16_append_rune(jane::Slice<jane::U16> &a,
                                         const jane::I32 &r) noexcept {
  if (0 <= r & r < jane::UTF16_SURR1 | jane::UTF16_SURR3 <= r &&
      r < jane::UTF16_SURR_SELF) {
    a.push(static_cast<jane::U16>(r));
  } else if (jane::UTF16_SURR_SELF <= r && r <= jane::UTF16_MAX_RUNE) {
    jane::I32 r1;
    jane::I32 r2;
    std::tie(r1, r2) = jane::utf16_encode_rune(r);
    a.push(static_cast<jane::U16>(r1));
    a.push(static_cast<jane::U16>(r2));
    return a;
  }
  a.push(jane::UTF16_REPLACEMENT_CHAR);
  return a;
}

jane::Slice<jane::U16> utf16_from_str(const jane::Str &s) noexcept {
  constexpr char NULL_TERMINATOR = '\x00';
  jane::Slice<jane::U16> buff{nullptr};
  jane::Slice<jane::I32> runes{static_cast<jane::Slice<jane::I32>>(s)};
  for (const jane::I32 &r : runes) {
    if (r == NULL_TERMINATOR) {
      break;
    }
    buff = jane::utf16_append_rune(buff, r);
  }
  return jane::utf16_append_rune(buff, NULL_TERMINATOR);
}

} // namespace jane

#endif // __JANE_UTF16_HPP