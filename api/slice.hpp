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

#ifndef __JANE_SLICE_HPP
#define __JANE_SLICE_HPP

#include "ref.hpp"
#include "typedef.hpp"

template <typename _Item_t> class slice_jnt;

template <typename _Item_t> class slice_jnt {
public:
  ref_jnt<_Item_t> __data{};
  _Item_t *__slice{nil};
  uint_jnt __len{0};
  uint_jnt __cap{0};

  slice_jnt<_Item_t>(void) noexcept {}
  slice_jnt<_Item_t>(const std::nullptr_t) noexcept {}

  slice_jnt<_Item_t>(const uint_jnt &_N) noexcept {
    const uint_jnt _n{_N < 0 ? 0 : _N};
    if (_n == 0) {
      return;
    }
    this->__alloc_new(_n);
  }

  slice_jnt<_Item_t>(const slice_jnt<_Item_t> &_Src) noexcept {
    this->operator=(_Src);
  }

  slice_jnt<_Item_t>(const std::initializer_list<_Item_t> &_Src) noexcept {
    if (_Src.size() == 0) {
      return;
    }
    this->__alloc_new(_Src.size());
    const auto _Src_begin{_Src.begin()};
    for (int_jnt _i{0}; _i < this->__len; ++_i) {
      this->__data.__alloc[_i] = *(_Item_t *)(_Src_begin + _i);
    }
  }

  ~slice_jnt<_Item_t>(void) noexcept { this->__dealloc(); }

  inline void __check(void) const noexcept {
    if (this->operator==(nil)) {
      JANE_ID(panic)(__JANE_ERROR_INVALID_MEMORY);
    }
  }

  void __dealloc(void) noexcept {
    this->__len = 0;
    this->__cap = 0;
    if (!this->__data.__ref) {
      this->__data.__alloc = nil;
      return;
    }
    if ((this->__data.__get_ref_n()) != __JANE_REFERENCE_DELTA) {
      this->__data.__alloc = nil;
      return;
    }
    delete this->__data.__ref;
    this->__data.__ref = nil;
    delete[] this->__data.__alloc;
    this->__data.__alloc = nil;
    this->__data.__ref = nil;
    this->__slice = nil;
  }

  void __alloc_new(const int_jnt _N) noexcept {
    this->__dealloc();
    _Item_t *_alloc{new (std::nothrow) _Item_t[_N]{_Item_t()}};
    if (!_alloc) {
      JANE_ID(panic)(__JANE_ERROR_MEMORY_ALLOCATION_FAILED);
    }
    this->__data = ref_jnt<_Item_t>::make(_alloc);
    this->__len = _N;
    this->__cap = _N;
    this->__slice = &_alloc[0];
  }

  typedef _Item_t *iterator;
  typedef const _Item_t *const_iterator;

  inline constexpr iterator begin(void) noexcept { return &this->__slice[0]; }

  inline constexpr const_iterator begin(void) const noexcept {
    return &this->__slice[0];
  }

  inline constexpr iterator end(void) noexcept {
    return &this->__slice[this->__len];
  }

  inline constexpr const_iterator end(void) const noexcept {
    return &this->__slice[this->__len];
  }

  inline slice_jnt<_Item_t> ___slice(const int_jnt &_Start,
                                     const int_jnt &_End) const noexcept {
    this->__check();
    if (_Start < 0 || _End < 0 || _Start > _End || _End > this->_cap()) {
      std::stringstream _sstream;
      __JANE_WRITE_ERROR_SLICING_INDEX_OUT_OF_RANGE(_sstream, _Start, _End);
      JANE_ID(panic)(_sstream.str().c_str());
    } else if (_Start == _End) {
      return slice_jnt<_Item_t>();
    }
    slice_jnt<_Item_t> _slice;
    _slice.__data = this->__data;
    _slice.__slice = &this->__slice[_Start];
    _slice.__len = _End - _Start;
    _slice.__cap = this->__cap - _Start;
    return _slice;
  }

  inline slice_jnt<_Item_t> ___slice(const int_jnt &_Start) const noexcept {
    return this->___slice(_Start, this->_len());
  }

  inline slice_jnt<_Item_t> ___slice(void) const noexcept {
    return this->___slice(0, this->_len());
  }

  inline constexpr int_jnt _len(void) const noexcept { return (this->__len); }

  inline constexpr int_jnt _cap(void) const noexcept { return (this->__cap); }

  inline bool _empty(void) const noexcept {
    return (!this->__slice || this->__len == 0 || this->__cap == 0);
  }

  void __push(const _Item_t &_Item) noexcept {
    if (this->__len == this->__cap) {
      _Item_t *_new{new (std::nothrow) _Item_t[this->__len + 1]};
      if (!_new) {
        JANE_ID(panic)(__JANE_ERROR_MEMORY_ALLOCATION_FAILED);
      }
      for (int_jnt _index{0}; _index < this->__len; ++_index) {
        _new[_index] = this->__data.__alloc[_index];
      }
      _new[this->__len] = _Item;
      delete[] this->__data.__alloc;
      this->__data.__alloc = nil;
      this->__data.__alloc = _new;
      this->__slice = this->__data.__alloc;
      ++this->__cap;
    } else {
      this->__slice[this->__len] = _Item;
    }
    ++this->__len;
  }

  bool operator==(const slice_jnt<_Item_t> &_Src) const noexcept {
    if (this->__len != _Src.__len) {
      return false;
    }
    for (int_jnt _index{0}; _index < this->__len; ++_index) {
      if (this->__slice[_index] != _Src.__slice[_index]) {
        return (false);
      }
    }
    return (true);
  }

  inline constexpr bool
  operator!=(const slice_jnt<_Item_t> &_Src) const noexcept {
    return !this->operator==(_Src);
  }

  inline constexpr bool operator==(const std::nullptr_t) const noexcept {
    return !this->__slice;
  }

  inline constexpr bool operator!=(const std::nullptr_t) const noexcept {
    return !this->operator==(nil);
  }

  _Item_t &operator[](const int_jnt &_Index) const {
    this->__check();
    if (this->_empty() || _Index < 0 || this->_len() <= _Index) {
      std::stringstream _sstream;
      __JANE_WRITE_ERROR_INDEX_OUT_OF_RANGE(_sstream, _Index);
      JANE_ID(panic)(_sstream.str().c_str());
    }
    return this->__slice[_Index];
  }

  void operator=(const slice_jnt<_Item_t> &_Src) noexcept {
    this->__dealloc();
    if (_Src.operator==(nil)) {
      return;
    }
    this->__len = _Src.__len;
    this->__cap = _Src.__cap;
    this->__data = _Src.__data;
    this->__slice = _Src.__slice;
  }

  void operator=(const std::nullptr_t) noexcept { this->__dealloc(); }

  friend std::ostream &operator<<(std::ostream &_Stream,
                                  const slice_jnt<_Item_t> &_Src) noexcept {
    if (_Src._empty()) {
      return (_Stream << "[]");
    }
    _Stream << '[';
    for (int_jnt _index{0}; _index < _Src.__len;) {
      _Stream << _Src.__slice[_index++];
      if (_index < _Src.__len) {
        _Stream << ' ';
      }
      _Stream << ']';
      return (_Stream);
    }
  }
};

#endif // !__JANE_SLICE_HPP
