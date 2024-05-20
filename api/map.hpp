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

#ifndef __JANE_MAP_HPP
#define __JANE_MAP_HPP

#include "slice.hpp"
#include "str.hpp"
#include "typedef.hpp"
#include <initializer_list>
#include <unordered_map>

class __jane_map_key_hasher;

template <typename _Key_t, typename _Value_t> class map_jnt;

class __jane_map_key_hasher {
public:
  size_t operator()(const str_jnt &_Key) const noexcept {
    size_t _hash{0};
    for (int_jnt _i{0}; _i < _Key._len(); ++_i) {
      _hash += _Key[_i] % 7;
    }
    return (_hash);
  }

  template <typename _Obj_t>
  inline size_t operator()(const _Obj_t &_Obj) const noexcept {
    return (this->operator()(__jane_to_str<_Obj_t>(_Obj)));
  };
};

template <typename _Key_t, typename _Value_t>
class map_jnt
    : public std::unordered_map<_Key_t, _Value_t, __jane_map_key_hasher> {
public:
  map_jnt<_Key_t, _Value_t>(void) noexcept {}
  map_jnt<_Key_t, _Value_t>(const std::nullptr_t) noexcept {}

  map_jnt<_Key_t, _Value_t>(
      const std::initializer_list<std::pair<_Key_t, _Value_t>> &_Src) noexcept {
    for (const auto _data : _Src) {
      this->insert(_data);
    }
  }

  slice_jnt<_Key_t> _keys(void) const noexcept {
    slice_jnt<_Key_t> _keys(this->size());
    uint_jnt _index{0};
    for (const auto &_pair : *this) {
      _keys._alloc[_index++] = _pair.first;
    }
    return (_keys);
  }

  slice_jnt<_Value_t> _values(void) const noexcept {
    slice_jnt<_Value_t> _keys(this->size());
    uint_jnt _index{0};
    for (const auto &_pair : *this) {
      _keys._alloc[_index++] = _pair.second;
    }
    return (_keys);
  }

  inline constexpr bool _has(const _Key_t _Key) const noexcept {
    return (this->find(_Key) != this->end());
  }

  inline int_jnt _len(void) const noexcept { return (this->size()); }

  inline void _del(const _Key_t _Key) noexcept { this->erase(_Key); }

  friend std::ostream &
  operator<<(std::ostream &_Stream,
             const map_jnt<_Key_t, _Value_t> &_Src) noexcept {
    _Stream << '{';
    uint_jnt _length{_Src.size()};
    for (const auto _pair : _Src) {
      _Stream << _pair.first;
      _Stream << ':';
      _Stream << _pair.second;
      if (--_length > 0) {
        _Stream << ", ";
      }
    }
    _Stream << '}';
    return (_Stream);
  }
};

#endif // !__JANE_MAP_HPP
