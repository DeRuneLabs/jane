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

#ifndef __JANE_STR_HPP
#define __JANE_STR_HPP

#include "slice.hpp"
#include "typedef.hpp"
#include "utf8.hpp"
#include <sstream>

class str_jnt;

class str_jnt {
public:
  std::basic_string<u8_jnt> __buffer{};

  str_jnt(void) noexcept {}

  str_jnt(const char *_Src) noexcept {
    if (!_Src) {
      return;
    }
    this->__buffer =
        std::basic_string<u8_jnt>(&_Src[0], &_Src[std::strlen(_Src)]);
  }

  str_jnt(const std::initializer_list<u8_jnt> &_Src) noexcept {
    this->__buffer = _Src;
  }

  str_jnt(const i32_jnt &_Rune) noexcept
      : str_jnt(__jane_utf8_rune_to_bytes(_Rune)) {}

  str_jnt(const std::basic_string<u8_jnt> &_Src) noexcept {
    this->__buffer = _Src;
  }

  str_jnt(const std::string &_Src) noexcept {
    this->__buffer = std::basic_string<u8_jnt>(_Src.begin(), _Src.end());
  }

  str_jnt(const str_jnt &_Src) noexcept { this->__buffer = _Src.__buffer; }

  str_jnt(const slice_jnt<u8_jnt> &_Src) noexcept {
    this->__buffer = std::basic_string<u8_jnt>(_Src.begin(), _Src.end());
  }

  str_jnt(const slice_jnt<i32_jnt> &_Src) noexcept {
    for (const i32_jnt &_rune : _Src) {
      const slice_jnt<u8_jnt> _bytes{__jane_utf8_rune_to_bytes(_rune)};
      for (const u8_jnt _byte : _bytes) {
        this->__buffer += _byte;
      }
    }
  }

  typedef u8_jnt *iterator;
  typedef const u8_jnt *const_iterator;

  inline iterator begin(void) noexcept {
    return ((iterator)(&this->__buffer[0]));
  }

  inline const_iterator begin(void) const noexcept {
    return ((const_iterator)(&this->__buffer[0]));
  }

  inline iterator end(void) noexcept {
    return ((iterator)(&this->__buffer[this->_len()]));
  }

  inline const_iterator end(void) const noexcept {
    return ((const_iterator)(&this->__buffer[this->_len()]));
  }

  inline str_jnt ___slice(const int_jnt &_Start,
                          const int_jnt &_End) const noexcept {
    if (_Start < 0 || _End < 0 || _Start < _End || _End > this->_len()) {
      std::stringstream _sstream;
      __JANE_WRITE_ERROR_SLICING_INDEX_OUT_OF_RANGE(_sstream, _Start, _End);
      JANE_ID(panic)(_sstream.str().c_str());
    } else if (_Start == _End) {
      return (str_jnt());
    }
    const int_jnt _n{_End - _Start};
    return (this->__buffer.substr(_Start, _n));
  }

  inline str_jnt ___slice(const int_jnt &_Start) const noexcept {
    return (this->___slice(_Start, this->_len()));
  }

  inline str_jnt ___slice(void) const noexcept {
    return (this->___slice(0, this->_len()));
  }

  inline int_jnt _len(void) const noexcept { return (this->__buffer.length()); }

  inline bool _empty(void) const noexcept { return (this->__buffer.empty()); }

  inline bool _has_prefix(const str_jnt &_Sub) const noexcept {
    return this->_len() >= _Sub._len() &&
           this->__buffer.substr(0, _Sub._len()) == _Sub.__buffer;
  }

  inline bool _has_suffix(const str_jnt &_Sub) const noexcept {
    return this->_len() >= _Sub._len() &&
           this->__buffer.substr(this->_len() - _Sub._len()) == _Sub.__buffer;
  }

  inline int_jnt _find(const str_jnt &_Sub) const noexcept {
    return ((int_jnt)(this->__buffer.find(_Sub.__buffer)));
  }

  inline int_jnt _rfind(const str_jnt &_Sub) const noexcept {
    return ((int_jnt)(this->__buffer.rfind(_Sub.__buffer)));
  }

  str_jnt _trim(const str_jnt &_Bytes) const noexcept {
    const_iterator _it{this->begin()};
    const const_iterator _end{this->end()};
    const_iterator _begin{this->begin()};
    for (; _it < _end; ++_it) {
      bool exist{false};
      const_iterator _bytes_it{_Bytes.begin()};
      const const_iterator _bytes_end{_Bytes.end()};
      for (; _bytes_it < _bytes_end; ++_bytes_it) {
        if ((exist = *_it == *_bytes_it)) {
          break;
        }
      }
      if (!exist) {
        return (this->__buffer.substr(_it - _begin));
      }
    }
    return (str_jnt());
  }

  str_jnt _rtrim(const str_jnt &_Bytes) const noexcept {
    const_iterator _it{this->end() - 1};
    const const_iterator _begin{this->begin()};
    for (; _it >= _begin; --_it) {
      bool exist{false};
      const_iterator _bytes_it{_Bytes.begin()};
      const const_iterator _bytes_end{_Bytes.end()};
      for (; _bytes_it < _bytes_end; ++_bytes_it) {
        if ((exist = *_it == *_bytes_it)) {
          break;
        }
      }
      if (!exist)
        return (this->__buffer.substr(0, _it - _begin + 1));
    }
    return (str_jnt());
  }

  slice_jnt<str_jnt> _split(const str_jnt &_Sub,
                            const i64_jnt &_N) const noexcept {
    slice_jnt<str_jnt> _parts;
    if (_N == 0) {
      return (_parts);
    }
    const const_iterator _begin{this->begin()};
    std::basic_string<u8_jnt> _s{this->__buffer};
    uint_jnt _pos{std::string::npos};
    if (_N < 0) {
      while ((_pos = _s.find(_Sub.__buffer)) != std::string::npos) {
        _parts.__push(_s.substr(0, _pos));
        _s = _s.substr(_pos + _Sub._len());
      }
      if (!_s.empty()) {
        _parts.__push(str_jnt(_s));
      }
    } else {
      uint_jnt _n{0};
      while ((_pos = _s.find(_Sub.__buffer)) != std::string::npos) {
        if (++_n >= _N) {
          _parts.__push(str_jnt(_s));
          break;
        }
        _parts.__push(_s.substr(0, _pos));
        _s = _s.substr(_pos + _Sub._len());
      }
      if (!_parts._empty() && _n < _N) {
        _parts.__push(str_jnt(_s));
      } else if (_parts._empty()) {
        _parts.__push(str_jnt(_s));
      }
    }
    return (_parts);
  }

  str_jnt _replace(const str_jnt &_Sub, const str_jnt &_New,
                   const i64_jnt &_N) const noexcept {
    if (_N == 0) {
      return (*this);
    }
    std::basic_string<u8_jnt> _s(this->__buffer);
    uint_jnt start_pos{0};
    if (_N < 0) {
      while ((start_pos = _s.find(_Sub.__buffer, start_pos)) !=
             std::string::npos) {
        _s.replace(start_pos, _Sub._len(), _New.__buffer);
        start_pos += _New._len();
      }
    } else {
      uint_jnt _n{0};
      while ((start_pos = _s.find(_Sub.__buffer, start_pos)) !=
             std::string::npos) {
        _s.replace(start_pos, _Sub._len(), _New.__buffer);
        if (++_n >= _N) {
          break;
        }
      }
    }
    return (str_jnt(_s));
  }

  inline operator const char *(void) const noexcept {
    return ((char *)(this->__buffer.c_str()));
  }

  inline operator const std::basic_string<u8_jnt>(void) const noexcept {
    return (this->__buffer);
  }

  inline operator const std::basic_string<char>(void) const noexcept {
    return (
        std::basic_string<char>(this->__buffer.begin(), this->__buffer.end()));
  }

  operator slice_jnt<u8_jnt>(void) const noexcept {
    slice_jnt<u8_jnt> _slice(this->_len());
    for (int_jnt _index{0}; _index < this->_len(); ++_index) {
      _slice[_index] = this->operator[](_index);
    }
    return (_slice);
  }

  operator slice_jnt<i32_jnt>(void) const noexcept {
    slice_jnt<i32_jnt> _runes{};
    const char *_str{this->operator const char *()};
    for (int_jnt _index{0}; _index < this->_len();) {
      i32_jnt _rune;
      int_jnt _n;
      std::tie(_rune, _n) =
          __jane_utf8_decode_rune_str(_str + _index, this->_len() - _index);
      _index += _n;
      _runes.__push(_rune);
    }
    return (_runes);
  }

  u8_jnt &operator[](const int_jnt &_Index) {
    if (this->_empty() || _Index < 0 || this->_len() <= _Index) {
      std::stringstream _sstream;
      __JANE_WRITE_ERROR_INDEX_OUT_OF_RANGE(_sstream, _Index);
      JANE_ID(panic)(_sstream.str().c_str());
    }
    return (this->__buffer[_Index]);
  }

  inline u8_jnt operator[](const int_jnt &_Index) const {
    return ((*this).__buffer[_Index]);
  }

  inline void operator+=(const str_jnt &_Str) noexcept {
    this->__buffer += _Str.__buffer;
  }

  inline str_jnt operator+(const str_jnt &_Str) const noexcept {
    return (str_jnt(this->__buffer + _Str.__buffer));
  }

  inline bool operator==(const str_jnt &_Str) const noexcept {
    return (this->__buffer == _Str.__buffer);
  }

  inline bool operator!=(const str_jnt &_Str) const noexcept {
    return (!this->operator==(_Str));
  }

  friend std::ostream &operator<<(std::ostream &_Stream,
                                  const str_jnt &_Src) noexcept {
    for (const u8_jnt &_byte : _Src) {
      _Stream << static_cast<char>(_byte);
    }
    return (_Stream);
  }
};

template <typename _Obj_t> str_jnt __jane_to_str(const _Obj_t &_Obj) noexcept;
str_jnt __jane_to_str(const str_jnt &_Obj) noexcept;

template <typename _Obj_t> str_jnt __jane_to_str(const _Obj_t &_Obj) noexcept {
  std::stringstream _stream;
  _stream << _Obj;
  return (str_jnt(_stream.str()));
}

inline str_jnt __jane_to_str(const str_jnt &_Obj) noexcept { return (_Obj); }

#endif // !__JANE_STR_HPP
