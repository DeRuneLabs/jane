// Copyright (c) 2024 - DeRuneLabs
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

#ifndef __JNC_STR_HPP
#define __JNC_STR_HPP

#include "jn_util.hpp"
#include "typedef.hpp"
#include "slice.hpp"

class str_jnt;

class str_jnt {
  public:
    std::basic_string<u8_jnt> _buffer{};

    str_jnt(void) noexcept {}

    str_jnt(const char *_Src) noexcept {
        if (!_Src) { return; }
        this->_buffer = std::basic_string<u8_jnt>(&_Src[0], &_Src[std::strlen(_Src)]);
    }

    str_jnt(const std::initializer_list<u8_jnt> &_Src) noexcept
    { this->_buffer = _Src; }

    str_jnt(const std::basic_string<u8_jnt> &_Src) noexcept
    { this->_buffer = _Src; }

    str_jnt(const std::string &_Src) noexcept
    { this->_buffer = std::basic_string<u8_jnt>(_Src.begin(), _Src.end()); }

    str_jnt(const str_jnt &_Src) noexcept
    { this->_buffer = _Src._buffer; }

    str_jnt(const uint_jnt &_N) noexcept
    { this->_buffer = std::basic_string<u8_jnt>(0, _N); }

    str_jnt(const slice<u8_jnt> &_Src) noexcept
    { this->_buffer = std::basic_string<u8_jnt>(_Src.begin(), _Src.end()); }

    typedef u8_jnt       *iterator;
    typedef const u8_jnt *const_iterator;

    inline iterator begin(void) noexcept
    { return (iterator)(&this->_buffer[0]); }

    inline const_iterator begin(void) const noexcept
    { return (const_iterator)(&this->_buffer[0]); }

    inline iterator end(void) noexcept
    { return (iterator)(&this->_buffer[this->len()]); }

    inline const_iterator end(void) const noexcept
    { return (const_iterator)(&this->_buffer[this->len()]); }

    inline str_jnt ___slice(const int_jnt &_Start,
                           const int_jnt &_End) const noexcept {
        if (_Start < 0 || _End < 0 || _Start > _End) {
            std::stringstream _sstream;
            _sstream << "index out of range [" << _Start << ':' << _End << ']';
            JNID(panic)(_sstream.str().c_str());
        } else if (_Start == _End) { return str_jnt(); }
        const int_jnt _n{_End-_Start};
        return this->_buffer.substr(_Start, _n);
    }

    inline str_jnt ___slice(const int_jnt &_Start) const noexcept
    { return this->___slice(_Start, this->len()); }

    inline str_jnt ___slice(void) const noexcept
    { return this->___slice(0, this->len()); }

    inline int_jnt len(void) const noexcept
    { return this->_buffer.length(); }

    inline bool empty(void) const noexcept
    { return this->_buffer.empty(); }

    inline bool has_prefix(const str_jnt &_Sub) const noexcept {
        return this->len() >= _Sub.len() &&
                this->_buffer.substr(0, _Sub.len()) == _Sub._buffer;
    }

    inline bool has_suffix(const str_jnt &_Sub) const noexcept {
        return this->len() >= _Sub.len() &&
            this->_buffer.substr(this->len()-_Sub.len()) == _Sub._buffer;
    }

    inline int_jnt find(const str_jnt &_Sub) const noexcept
    { return (int_jnt)(this->_buffer.find(_Sub._buffer)); }

    inline int_jnt rfind(const str_jnt &_Sub) const noexcept
    { return (int_jnt)(this->_buffer.rfind(_Sub._buffer)); }

    inline const char* cstr(void) const noexcept
    { return (const char*)(this->_buffer.c_str()); }

    str_jnt trim(const str_jnt &_Bytes) const noexcept {
        const_iterator _it{this->begin()};
        const const_iterator _end{this->end()};
        const_iterator _begin{this->begin()};
        for (; _it < _end; ++_it) {
            bool exist{false};
            const_iterator _bytes_it{_Bytes.begin()};
            const const_iterator _bytes_end{_Bytes.end()};
            for (; _bytes_it < _bytes_end; ++_bytes_it)
            { if ((exist = *_it == *_bytes_it)) { break; } }
            if (!exist) { return this->_buffer.substr(_it-_begin); }
        }
        return str_jnt{""};
    }

    str_jnt rtrim(const str_jnt &_Bytes) const noexcept {
        const_iterator _it{this->end()-1};
        const const_iterator _begin{this->begin()};
        for (; _it >= _begin; --_it) {
            bool exist{false};
            const_iterator _bytes_it{_Bytes.begin()};
            const const_iterator _bytes_end{_Bytes.end()};
            for (; _bytes_it < _bytes_end; ++_bytes_it)
            { if ((exist = *_it == *_bytes_it)) { break; } }
            if (!exist) { return this->_buffer.substr(0, _it-_begin+1); }
        }
        return str_jnt{""};
    }

    slice<str_jnt> split(const str_jnt &_Sub, const i64_jnt &_N) const noexcept {
        slice<str_jnt> _parts;
        if (_N == 0) { return _parts; }
        const const_iterator _begin{this->begin()};
        std::basic_string<u8_jnt> _s{this->_buffer};
        uint_jnt _pos{std::string::npos};
        if (_N < 0) {
            while ((_pos = _s.find(_Sub._buffer)) != std::string::npos) {
                _parts.__push(_s.substr(0, _pos));
                _s = _s.substr(_pos+_Sub.len());
            }
            if (!_parts.empty()) { _parts.__push(str_jnt{_s}); }
        } else {
            uint_jnt _n{0};
            while ((_pos = _s.find(_Sub._buffer)) != std::string::npos) {
                _parts.__push(_s.substr(0, _pos));
                _s = _s.substr(_pos+_Sub.len());
                if (++_n >= _N) { break; }
            }
            if (!_parts.empty() && _n < _N) { _parts.__push(str_jnt{_s}); }
        }
        return _parts;
    }

    str_jnt replace(const str_jnt &_Sub,
                   const str_jnt &_New,
                   const i64_jnt &_N) const noexcept {
        if (_N == 0) { return *this; }
        std::basic_string<u8_jnt> _s{this->_buffer};
        uint_jnt start_pos{0};
        if (_N < 0) {
            while((start_pos = _s.find(_Sub._buffer, start_pos)) != std::string::npos) {
                _s.replace(start_pos, _Sub.len(), _New._buffer);
                start_pos += _New.len();
            }
        } else {
            uint_jnt _n{0};
            while((start_pos = _s.find(_Sub._buffer, start_pos)) != std::string::npos) {
                _s.replace(start_pos, _Sub.len(), _New._buffer);
                start_pos += _New.len();
                if (++_n >= _N) { break; }
            }
        }
        return str_jnt{_s};
    }

    operator slice<u8_jnt>(void) const noexcept {
        slice<u8_jnt> _slice(this->len());
        for (int_jnt _index{0}; _index < this->len(); ++_index)
        { _slice[_index] = this->operator[](_index);  }
        return _slice;
    }

    u8_jnt &operator[](const int_jnt &_Index) {
        if (this->empty() || _Index < 0 || this->len() <= _Index) {
            std::stringstream _sstream;
            _sstream << "index out of range [" << _Index << ']';
            JNID(panic)(_sstream.str().c_str());
        }
        return this->_buffer[_Index];
    }

    inline u8_jnt operator[](const uint_jnt &_Index) const
    { return (*this)._buffer[_Index]; }

    inline void operator+=(const str_jnt &_Str) noexcept
    { this->_buffer += _Str._buffer; }

    inline str_jnt operator+(const str_jnt &_Str) const noexcept
    { return str_jnt{this->_buffer + _Str._buffer}; }

    inline bool operator==(const str_jnt &_Str) const noexcept
    { return this->_buffer == _Str._buffer; }

    inline bool operator!=(const str_jnt &_Str) const noexcept
    { return !this->operator==(_Str); }

    friend std::ostream &operator<<(std::ostream &_Stream, const str_jnt &_Src) noexcept {
        for (const u8_jnt &_byte: _Src)
        { _Stream << _byte; }
        return _Stream;
    }
};

#endif // !__JNC_STR_HPP
