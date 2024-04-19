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

#ifndef __JNC_MAP_HPP
#define __JNC_MAP_HPP

#include "jn_util.hpp"
#include "slice.hpp"
#include "typedef.hpp"

template <typename _Key_t, typename _Value_t> class map;

template <typename _Key_t, typename _Value_t>
class map : public std::unordered_map<_Key_t, _Value_t> {
public:
  map<_Key_t, _Value_t>(void) noexcept {}
  map<_Key_t, _Value_t>(const std::nullptr_t) noexcept {}

  map<_Key_t, _Value_t>(
      const std::initializer_list<std::pair<_Key_t, _Value_t>> _Src) noexcept {
    for (const auto _data : _Src) {
      this->insert(_data);
    }
  }

  slice<_Key_t> keys(void) noexcept {
    slice<_Key_t> _keys(this->size());
    uint_jnt _index{0};
    for (const auto &_pair : *this) {
      _keys._alloc[_index++] = _pair.first;
    }
    return (_keys);
  }

  slice<_Value_t> values(void) const noexcept {
    slice<_Value_t> _keys(this->size());
    uint_jnt _index{0};
    for (const auto &_pair : *this) {
      _keys._alloc[_index++] = _pair.second;
    }
    return (_keys);
  }

  inline constexpr bool has(const _Key_t _Key) const noexcept {
    return (this->find(_Key) != this->end());
  }

  inline int_jnt len(void) const noexcept { return (this->size()); }

  inline void del(const _Key_t _Key) noexcept { this->erase(_Key); }

  inline bool operator==(const std::nullptr_t) const noexcept {
    return (this->empty());
  }

  inline bool operator!=(const std::nullptr_t) const noexcept {
    return (!this->operator==(nil));
  }

  friend std::ostream &operator<<(std::ostream &_Stream,
                                  const map<_Key_t, _Value_t> &_Src) noexcept {
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

#endif // !__JNC_MAP_HPP
