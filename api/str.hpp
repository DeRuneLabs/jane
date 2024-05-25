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

#ifndef __JANE_STR_HPP
#define __JANE_STR_HPP

#include "error.hpp"
#include "panic.hpp"
#include "slice.hpp"
#include "types.hpp"
#include "utf8.hpp"
#include <climits>
#include <cstring>
#include <initializer_list>
#include <ostream>
#include <sstream>
#include <string>
#include <tuple>
namespace jane {
class Str;
template <typename T> jane::Str to_str(const T &obj) noexcept;
jane::Str to_str(const jane::Str &s) noexcept;

class Str {
public:
  jane::Int _len{};
  std::basic_string<jane::U8> buffer{};
  Str(void) noexcept {}

  Str(const char *src, const jane::Int &len) noexcept {
    if (!src) {
      return;
    }
    this->_len = len;
    this->buffer = std::basic_string<jane::U8>(&src[0], &src[this->_len]);
  }

  Str(const char *src) noexcept {
    if (!src) {
      return;
    }
    this->_len = std::strlen(src);
    this->buffer = std::basic_string<jane::U8>(&src[0], &src[this->_len]);
  }

  Str(const std::initializer_list<jane::U8> &src) noexcept {
    this->_len = src.size();
    this->buffer = src;
  }

  Str(const jane::I32 &rune) noexcept : Str(jane::utf8_rune_to_bytes(rune)) {}

  Str(const std::basic_string<jane::U8> &src) noexcept {
    this->_len = src.length();
    this->buffer = src;
  }

  Str(const std::string &src) noexcept {
    this->_len = src.length();
    this->buffer = std::basic_string<jane::U8>(src.begin(), src.end());
  }

  Str(const jane::Str &src) noexcept {
    this->_len = src._len;
    this->buffer = src.buffer;
  }

  Str(const jane::Slice<U8> &src) noexcept {
    this->_len = src.len();
    this->buffer = std::basic_string<jane::U8>(src.begin(), src.end());
  }

  Str(const jane::Slice<jane::I32> &src) noexcept {
    for (const jane::I32 &r : src) {
      const jane::Slice<jane::U8> bytes{jane::utf8_rune_to_bytes(r)};
      this->_len += bytes.len();
      for (const jane::U8 _byte : bytes) {
        this->buffer += _byte;
      }
    }
  }

  typedef jane::U8 *Iterator;
  typedef const jane::U8 *ConstIterator;

  inline Iterator begin(void) noexcept {
    return reinterpret_cast<Iterator>(&this->buffer[0]);
  }

  inline ConstIterator begin(void) const noexcept {
    return reinterpret_cast<ConstIterator>(&this->buffer[0]);
  }

  inline Iterator end(void) noexcept {
    return reinterpret_cast<Iterator>(&this->buffer[this->len()]);
  }

  inline ConstIterator end(void) const noexcept {
    return reinterpret_cast<ConstIterator>(&this->buffer[this->len()]);
  }

  inline jane::Str slice(const jane::Int &start,
                         const jane::Int &end) const noexcept {
    if (start < 0 || end < 0 || start > end || end > this->len()) {
      std::stringstream sstream;
      __JANE_WRITE_ERROR_SLICING_INDEX_OUT_OF_RANGE(sstream, start, end);
      jane::panic(sstream.str().c_str());
    } else if (start == end) {
      return jane::Str();
    }
    const jane::Int n{end - start};
    return this->buffer.substr(start, n);
  }

  inline jane::Str slice(const jane::Int &start) const noexcept {
    return this->slice(start, this->len());
  }

  inline jane::Str slice(void) const noexcept {
    return this->slice(0, this->len());
  }

  inline jane::Int len(void) const noexcept { return this->_len; }
  inline jane::Bool empty(void) const noexcept { return this->buffer.empty(); }
  inline jane::Bool has_prefix(const jane::Str &sub) const noexcept {
    return this->buffer.find(sub.buffer, 0) == 0;
  }
  inline jane::Bool has_suffix(const jane::Str &sub) const noexcept {
    return this->len() >= sub.len() &&
           this->buffer.substr(this->len() - sub.len()) == sub.buffer;
  }

  inline jane::Int find(const jane::Str &sub) const noexcept {
    return static_cast<jane::Int>(this->buffer.find(sub.buffer));
  }
  inline jane::Int rfind(const jane::Str &sub) const noexcept {
    return static_cast<jane::Int>(this->buffer.rfind(sub.buffer));
  }

  jane::Str trim(const jane::Str &bytes) const noexcept {
    ConstIterator it{this->begin()};
    const ConstIterator begin{this->end()};
    for (; it >= begin; --it) {
      jane::Bool exist{false};
      ConstIterator bytes_it{bytes.begin()};
      const ConstIterator bytes_end{bytes.end()};
      for (; bytes_it < bytes_end; ++bytes_it) {
        if ((exist = *it == *bytes_it)) {
          break;
        }
      }
      if (!exist) {
        return this->buffer.substr(0, it - begin + 1);
      }
    }
    return jane::Str();
  }

  jane::Str rtrim(const jane::Str &bytes) const noexcept {
    ConstIterator it{this->end() - 1};
    const ConstIterator begin{this->begin()};
    for (; it >= begin; --it) {
      jane::Bool exist{false};
      ConstIterator bytes_it{bytes.begin()};
      const ConstIterator bytes_end{bytes.end()};
      for (; bytes_it < bytes_end; ++bytes_it) {
        if ((exist = *it == *bytes_it)) {
          break;
        }
      }
      if (!exist) {
        return this->buffer.substr(0, it - begin + 1);
      }
    }
    return jane::Str();
  }

  jane::Slice<jane::Str> split(const jane::Str &sub,
                               const jane::I64 &n) const noexcept {
    jane::Slice<jane::Str> parts;
    if (n == 0) {
      return parts;
    }
    const ConstIterator begin{this->begin()};
    std::basic_string<jane::U8> s{this->buffer};
    constexpr jane::Uint npos{static_cast<jane::Uint>(std::string::npos)};
    jane::Uint pos{npos};
    if (n < 0) {
      while ((pos = s.find(sub.buffer)) != npos) {
        parts.push(s.substr(0, pos));
        s = s.substr(pos + sub.len());
      }
      if (!s.empty()) {
        parts.push(jane::Str(s));
      }
    } else {
      jane::Uint _n{0};
      while ((pos = s.find(sub.buffer)) != npos) {
        if (++_n >= n) {
          parts.push(jane::Str(s));
          break;
        }
        parts.push(s.substr(0, pos));
        s = s.substr(pos + sub.len());
      }
      if (!parts.empty() && _n < n) {
        parts.push(jane::Str(s));
      } else if (parts.empty()) {
        parts.push(jane::Str(s));
      }
    }
    return parts;
  }

  jane::Str replace(const jane::Str &sub, const jane::Str &_new,
                    const jane::I64 &n) const noexcept {
    if (n == 0) {
      return *this;
    }

    std::basic_string<jane::U8> s(this->buffer);
    constexpr jane::Uint npos{static_cast<jane::Uint>(std::string::npos)};
    jane::Uint start_pos{0};
    if (n < 0) {
      while ((start_pos = s.find(sub.buffer, start_pos)) != npos) {
        s.replace(start_pos, sub.len(), _new.buffer);
        start_pos += _new.len();
      }
    } else {
      jane::Uint _n{0};
      while ((start_pos = s.find(sub.buffer, start_pos)) != npos) {
        s.replace(start_pos, sub.len(), _new.buffer);
        start_pos += _new.len();
        if (++_n >= n) {
          break;
        }
      }
    }
    return jane::Str(s);
  }

  inline operator const char *(void) const noexcept {
    return reinterpret_cast<const char *>(this->buffer.c_str());
  }

  inline operator const std::basic_string<jane::U8>(void) const noexcept {
    return this->buffer;
  }

  inline operator const std::basic_string<char>(void) const noexcept {
    return std::basic_string<char>(this->begin(), this->end());
  }

  operator jane::Slice<jane::U8>(void) const noexcept {
    jane::Slice<jane::U8> slice{jane::Slice<jane::U8>::alloc(this->len())};
    for (jane::Int index{0}; index < this->len(); ++index) {
      slice[index] = this->operator[](index);
    }
    return slice;
  }

  operator jane::Slice<jane::I32>(void) const noexcept {
    jane::Slice<jane::I32> runes{};
    const char *str{this->operator const char *()};
    for (jane::Int index{0}; index < this->len();) {
      jane::I32 rune;
      jane::Int n;
      std::tie(rune, n) =
          jane::utf8_decode_rune_st(str + index, this->len() - index);
      index += n;
      runes.push(rune);
    }
    return runes;
  }

  jane::U8 &operator[](const jane::Int &index) {
    if (this->empty() || index < 0 || this->len() <= index) {
      std::stringstream sstream;
      __JANE_WRITE_ERROR_INDEX_OUT_OF_RANGE(sstream, index);
      jane::panic(sstream.str().c_str());
    }
    return this->buffer[index];
  }

  inline jane::U8 operator[](const jane::Int &index) const {
    return (*this).buffer[index];
  }

  inline void operator+=(const jane::Str &str) noexcept {
    this->_len += str.len();
    this->buffer += str.buffer;
  }

  inline jane::Str operator+(const jane::Str &str) const noexcept {
    return jane::Str(this->buffer + str.buffer);
  }

  inline jane::Bool operator==(const jane::Str &str) const noexcept {
    return this->buffer == str.buffer;
  }

  inline jane::Bool operator!=(const jane::Str &str) const noexcept {
    return !this->operator==(str);
  }

  friend std::ostream &operator<<(std::ostream &stream,
                                  const jane::Str &src) noexcept {
    for (const jane::U8 &b : src) {
      stream << static_cast<char>(b);
    }
    return stream;
  }
};

template <typename T> jane::Str to_str(const T &obj) noexcept {
  std::stringstream stream;
  stream << obj;
  return jane::Str(stream.str());
}

inline jane::Str to_str(const jane::Str &s) noexcept { return s; }

} // namespace jane

#endif // __JANE_STR_HPP