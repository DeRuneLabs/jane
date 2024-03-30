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
#include "typedef.hpp"

template <typename _Item_t> class slice;

template <typename _Item_t> class slice {
public:
  _Item_t *_buffer{nil};
  mutable uint_jnt *_ref{nil};
  int_jnt _size{0};

  slice<_Item_t>(void) noexcept {}
  slice<_Item_t>(const std::nullptr_t) noexcept {}

  slice<_Item_t>(const uint_jnt &_N) noexcept { this->__alloc_new(_N); }

  slice<_Item_t>(const slice<_Item_t> &_Src) noexcept { this->operator=(_Src); }

  slice<_Item_t>(const std::initializer_list<_Item_t> &_Src) noexcept {
    this->__alloc_new(_Src.size());
    const auto _Src_begin{_Src.begin()};
    for (int_jnt _index{0}; _index < this->_size; ++_index) {
      this->_buffer[_index] = *(_Item_t *)(_Src_begin + _index);
    }
  }

  ~slice<_Item_t>(void) noexcept { this->__dealloc(); }

  void __check(void) const noexcept {
    if (!this->_buffer) {
      JNID(panic)("invalid memory address or nil pointer deference");
    }
  }

  void __dealloc(void) noexcept {
    if (!this->_ref) {
      return;
    }
    (*this->_ref)--;
    if ((*this->_ref) != 0) {
      return;
    }
    delete this->_ref;
    this->_ref = nil;
    delete[] this->_buffer;
    this->_buffer = nil;
    this->_size = 0;
  }

  void __alloc_new(const int_jnt _N) noexcept {
    this->__dealloc();
    this->_buffer = new (std::nothrow) _Item_t[_N];
    if (!this->_buffer) {
      JNID(panic)("memory allocation failed");
    }
    this->_ref = new (std::nothrow) uint_jnt{1};
    if (!this->_ref) {
      JNID(panic)("memory allocation failed");
    }
    this->_size = _N;
  }

  typedef _Item_t *iterator;
  typedef const _Item_t *const_iterator;

  inline constexpr iterator begin(void) noexcept { return &this->_buffer[0]; }

  inline constexpr const_iterator begin(void) const noexcept {
    return &this->_buffer[0];
  }

  inline constexpr iterator end(void) noexcept {
    return &this->_buffer[this->_size];
  }

  inline constexpr const_iterator end(void) const noexcept {
    return &this->_buffer[this->_size];
  }

  inline slice<_Item_t> ___slice(const int_jnt &_Start,
                                 const int_jnt &_End) const noexcept {
    this->__check();
    if (_Start < 0 || _End < 0 || _Start > _End) {
      std::stringstream _sstream;
      _sstream << "index out of range [" << _Start << ':' << _End << ']';
      JNID(panic)(_sstream.str().c_str());
    } else if (_Start == _End) {
      return slice<_Item_t>();
    }
    slice<_Item_t> _slice;
    _slice._buffer = &this->_buffer[_Start];
    _slice._size = _End - _Start;
    return _slice;
  }

  inline slice<_Item_t> ___slice(const int_jnt &_Start) const noexcept {
    return this->___slice(_Start, this->len());
  }

  inline slice<_Item_t> ___slice(void) const noexcept {
    return this->___slice(0, this->len());
  }

  inline constexpr int_jnt len(void) const noexcept { return this->_size; }

  inline bool empty(void) const noexcept {
    return !this->_buffer || this->_size == 0;
  }

  void __push(const _Item_t &_Item) noexcept {
    _Item_t *_new = new (std::nothrow) _Item_t[this->_size + 1];
    if (!_new) {
      JNID(panic)("memory allocation failed");
    }
    for (int_jnt _index{0}; _index < this->_size; ++_index) {
      _new[_index] = this->_buffer[_index];
    }
    _new[this->_size] = _Item;
    delete[] this->_buffer;
    this->_buffer = nil;
    this->_buffer = _new;
    ++this->_size;
  }

  bool operator==(const slice<_Item_t> &_Src) const noexcept {
    if (this->_size != _Src._size) {
      return false;
    }
    for (int_jnt _index{0}; _index < this->_size; ++_index) {
      if (this->_buffer[_index] != _Src._buffer[_index]) {
        return false;
      }
    }
    return true;
  }

  inline constexpr bool operator!=(const slice<_Item_t> &_Src) const noexcept {
    return !this->operator==(_Src);
  }

  inline constexpr bool operator==(const std::nullptr_t) const noexcept {
    return !this->_buffer;
  }

  inline constexpr bool operator!=(const std::nullptr_t) const noexcept {
    return !this->operator==(nil);
  }

  _Item_t &operator[](const int_jnt &_Index) {
    this->__check();
    if (this->empty() || _Index < 0 || this->len() <= _Index) {
      std::stringstream _sstream;
      _sstream << "index out of range [" << _Index << ']';
      JNID(panic)(_sstream.str().c_str());
    }
    return this->_buffer[_Index];
  }

  void operator=(const slice<_Item_t> &_Src) noexcept {
    this->__dealloc();
    if (_Src._ref) {
      (*_Src._ref)++;
    }
    this->_size = _Src._size;
    this->_ref = _Src._ref;
    this->_buffer = _Src._buffer;
  }

  void operator=(const std::nullptr_t) noexcept {
    if (!this->_ref) {
      this->_ptr = nil;
      return;
    }
    this->__dealloc();
  }

  friend std::ostream &operator<<(std::ostream &_Stream,
                                  const slice<_Item_t> &_Src) noexcept {
    _Stream << '[';
    for (int_jnt _index{0}; _index < _Src._size;) {
      _Stream << _Src._buffer[_index++];
      if (_index < _Src._size) {
        _Stream << ", ";
      }
    }
    _Stream << ']';
    return _Stream;
  }
};

#endif // !__JNC_SLICE_HPP
