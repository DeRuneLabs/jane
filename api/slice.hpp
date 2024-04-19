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

#ifndef __JNC_SLICE_HPP
#define __JNC_SLICE_HPP

#include "jn_util.hpp"
#include "ref.hpp"
#include "typedef.hpp"

template <typename _Item_t> class slice;

template <typename _Item_t> class slice {
public:
  jn_ref<_Item_t> _data{nil};
  _Item_t *_slice{nil};
  int_jnt _n{0};

  slice<_Item_t>(void) noexcept {}
  slice<_Item_t>(const std::nullptr_t) noexcept {}

  slice<_Item_t>(const uint_jnt &_N) noexcept {
    this->__alloc_new(_N < 0 ? 0 : _N);
  }

  slice<_Item_t>(const slice<_Item_t> &_Src) noexcept { this->operator=(_Src); }

  slice<_Item_t>(const std::initializer_list<_Item_t> &_Src) noexcept {
    this->__alloc_new(_Src.size());
    const auto _Src_begin{_Src.begin()};
    for (int_jnt _i{0}; _i < this->_n; ++_i) {
      this->_data._alloc[_i] = *(_Item_t)(_Src_begin + _i);
    }
  }

  ~slice<_Item_t>(void) noexcept { this->__dealloc(); }

  inline void __check(void) const noexcept {
    if (this->operator==(nil)) {
      JNC_ID(panic)(__JNC_ERROR_INVALID_MEMORY);
    }
  }

  void __dealloc(void) noexcept {
    if (!this->_data._ref) {
      return;
    }
    if ((this->_data.__get_ref_n()) != __JNC_REFERENCE_DELTA) {
      return;
    }
    delete this->_data._ref;
    this->_data._ref = nil;
    delete[] this->_data._alloc;
    this->_data._alloc = nil;
    this->_data._ref = nil;
    this->_slice = nil;
    this->_n = 0;
  }

  void __alloc_new(const int_jnt _N) noexcept {
    this->__dealloc();
    _Item_t *_alloc{new (std::nothrow) _Item_t[_N]{_Item_t()}};
    if (!_alloc) {
      JNC_ID(panic)(__JNC_ERROR_MEMORY_ALLOCATION_FAILED);
    }
    this->_data = jn_ref<_Item_t>(_alloc);
    this->_n = _N;
    this->_slice = &_alloc[0];
  }

  typedef _Item_t *iterator;
  typedef const _Item_t *const_iterator;

  inline constexpr iterator begin(void) noexcept { return &this->_slice[0]; }

  inline constexpr const_iterator begin(void) const noexcept {
    return &this->_slice[0];
  }

  inline constexpr iterator end(void) noexcept {
    return &this->_slice[this->_n];
  }

  inline constexpr const_iterator end(void) const noexcept {
    return &this->_slice[this->_n];
  }

  inline slice<_Item_t> ___slice(const int_jnt &_Start,
                                 const int_jnt &_End) const noexcept {
    this->__check();
    if (_Start < 0 || _End < 0 || _Start > _End) {
      std::stringstream _sstream;
      __JNC_WRITE_ERROR_SLICING_INDEX_OUT_OF_RANGE(_sstream, _Start, _End);
      JNC_ID(panic)(_sstream.str().c_str());
    } else if (_Start == _End) {
      return slice<_Item_t>();
    }
    slice<_Item_t> _slice;
    _slice._data = this->_data;
    _slice._slice = &this->_slice[_Start];
    _slice._n = _End - _Start;
    return _slice;
  }

  inline slice<_Item_t> ___slice(const int_jnt &_Start) const noexcept {
    return this->___slice(_Start, this->len());
  }

  inline slice<_Item_t> ___slice(void) const noexcept {
    return this->___slice(0, this->len());
  }

  inline constexpr int_jnt len(void) const noexcept {
    return !this->_slice || this->_n == 0;
  }

  inline bool empty(void) const noexcept {
    return !this->_slice || this->_n == 0;
  }

  void __push(const _Item_t &_Item) noexcept {
    _Item_t *_new = new (std::nothrow) _Item_t[this->_n + 1];
    if (!_new) {
      JNC_ID(panic)(__JNC_ERROR_MEMORY_ALLOCATION_FAILED);
    }
    for (int_jnt _index{0}; _index < this->_n; ++_index) {
      _new[_index] = this->_data._alloc[_index];
    }
    _new[this->_n] = _Item;
    delete[] this->_data._alloc;
    this->_data._alloc = nil;
    this->_data._alloc = _new;
    this->_slice = this->_data._alloc;
    ++this->_n;
  }

  bool operator==(const slice<_Item_t> &_Src) const noexcept {
    if (this->_n != _Src._n) {
      return false;
    }
    for (int_jnt _index{0}; _index < this->_n; ++_index) {
      if (this->_slice[_index] != _Src._slice[_index]) {
        return false;
      }
    }
    return true;
  }

  inline constexpr bool operator!=(const slice<_Item_t> &_Src) const noexcept {
    return !this->operator==(_Src);
  }

  inline constexpr bool operator==(const std::nullptr_t) const noexcept {
    return !this->_slice;
  }

  inline constexpr bool operator!=(const std::nullptr_t) const noexcept {
    return !this->operator==(nil);
  }

  _Item_t &operator[](const int_jnt &_Index) const {
    this->__check();
    if (this->empty() || _Index < 0 || this->len() <= _Index) {
      std::stringstream _sstream;
      __JNC_WRITE_ERROR_INDEX_OUT_OF_RANGE(_sstream, _Index);
      JNC_ID(panic)(_sstream.str().c_str());
    }
    return this->_slice[_Index];
  }

  void operator=(const slice<_Item_t> &_Src) noexcept {
    this->__dealloc();
    this->_n = _Src._n;
    this->_data = _Src._data;
    this->_slice = _Src._slice;
  }

  void operator=(const std::nullptr_t) noexcept { this->__dealloc(); }

  friend std::ostream &operator<<(std::ostream &_Stream,
                                  const slice<_Item_t> &_Src) noexcept {
    _Stream << '[';
    for (int_jnt _index{0}; _index < _Src._n;) {
      _Stream << _Src._slice[_index++];
      if (_index < _Src._n) {
        _Stream << ' ';
      }
    }
    _Stream << ']';
    return _Stream;
  }
};

#endif // !__JNC_SLICE_HPP
