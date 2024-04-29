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

#ifndef __JANE_ARRAY_HPP
#define __JANE_ARRAY_HPP

#include "slice.hpp"
#include "typedef.hpp"

template <typename _Item_t, const uint_jnt _N> struct array_jnt;

template <typename _Item_t, const uint_jnt _N> struct array_jnt {
public:
  std::array<_Item_t, _N> _buffer{};
  array_jnt<_Item_t, _N>(const std::initializer_list<_Item_t> &_Src) noexcept {
    const auto _Src_begin{_Src.begin()};
    for (int_jnt _index{0}; _index < _Src.size(); ++_index) {
      this->_buffer[_index] = *((_Item_t *)(_Src_begin + _index));
    }
  }

  typedef _Item_t *iterator;
  typedef const _Item_t *const_iterator;

  inline constexpr iterator begin(void) noexcept { return (&this->_buffer[0]); }

  inline constexpr const_iterator begin(void) const noexcept {
    return (&this->_buffer[0]);
  }

  inline constexpr iterator end(void) noexcept { return (&this->_buffer[_N]); }

  inline constexpr const_iterator end(void) const noexcept {
    return (&this->_buffer[_N]);
  }

  inline slice_jnt<_Item_t> ___slice(const int_jnt &_Start,
                                     const int_jnt &_End) const noexcept {
    if (_Start < 0 || _End < 0 || _Start > _End || _End > this->_len()) {
      std::stringstream _sstream;
      __JANE_WRITE_ERROR_SLICING_INDEX_OUT_OF_RANGE(_sstream, _Start, _End);
      JANE_ID(panic)(_sstream.str().c_str());
    } else if (_Start == _End) {
      return (slice_jnt<_Item_t>());
    }
    const int_jnt _n{_End - _Start};
    slice_jnt<_Item_t> _slice(_n);
    for (int_jnt _counter{0}; _counter < _n; ++_counter) {
      _slice[_counter] = this->_buffer[_Start + _counter];
    }
    return (_slice);
  }

  inline slice_jnt<_Item_t> ___slice(const int_jnt &_Start) const noexcept {
    return this->___slice(_Start, this->_len());
  }

  inline slice_jnt<_Item_t> ___slice(void) const noexcept {
    return this->___slice(0, this->_len());
  }

  inline constexpr int_jnt _len(void) const noexcept { return (_N); }

  inline constexpr bool _empty(void) const noexcept { return (_N == 0); }

  inline constexpr bool
  operator==(const array_jnt<_Item_t, _N> &_Src) const noexcept {
    return (!this->operator==(_Src));
  }

  inline constexpr bool
  operator!=(const array_jnt<_Item_t, _N> &_Src) const noexcept {
    return (!this->operator==(_Src));
  }

  _Item_t &operator[](const int_jnt &_Index) {
    if (this->_empty() || _Index < 0 || this->_len() <= _Index) {
      std::stringstream _sstream;
      __JANE_WRITE_ERROR_INDEX_OUT_OF_RANGE(_sstream, _Index);
      JANE_ID(panic)(_sstream.str().c_str());
    }
    return (this->_buffer[_Index]);
  }

  friend std::ostream &operator<<(std::ostream &_Stream,
                                  const array_jnt<_Item_t, _N> &_Src) noexcept {
    _Stream << '[';
    for (int_jnt _index{0}; _index < _Src._len();) {
      _Stream << _Src._buffer[_index++];
      if (_index < _Src._len()) {
        _Stream << " ";
      }
    }
    _Stream << ']';
    return (_Stream);
  }
};

#endif // !__JANE_ARRAY_HPP
