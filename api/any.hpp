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

#ifndef __JNC_ANY_HPP
#define __JNC_ANY_HPP

#include "jn_util.hpp"

struct any_jnt;

struct any_jnt {
  public:
    std::any _expr;

    any_jnt(void) noexcept {}

    template<typename T>
    any_jnt(const T &_Expr) noexcept
    { this->operator=(_Expr); }

    ~any_jnt(void) noexcept
    { this->_delete(); }

    inline void _delete(void) noexcept
    { this->_expr.reset(); }

    inline bool _isnil(void) const noexcept
    { return !this->_expr.has_value(); }

    template<typename T>
    inline bool type_is(void) const noexcept {
        if (std::is_same<T, std::nullptr_t>::value) { return false; }
        if (this->_isnil()) { return false; }
        return std::strcmp(this->_expr.type().name(), typeid(T).name()) == 0;
    }

    template<typename T>
    void operator=(const T &_Expr) noexcept {
        this->_delete();
        this->_expr = _Expr;
    }

    inline void operator=(const std::nullptr_t) noexcept
    { this->_delete(); }

    template<typename T>
    operator T(void) const noexcept {
        if (this->_isnil()) { JNID(panic)("invalid memory address or nil pointer deference"); }
        if (!this->type_is<T>()) { JNID(panic)("incompatible type"); }
        return std::any_cast<T>(this->_expr);
    }

    template<typename T>
    inline bool operator==(const T &_Expr) const noexcept
    { return this->type_is<T>() && this->operator T() == _Expr; }

    template<typename T>
    inline constexpr
    bool operator!=(const T &_Expr) const noexcept
    { return !this->operator==(_Expr); }

    inline bool operator==(const any_jnt &_Any) const noexcept {
        if (this->_isnil() && _Any._isnil()) { return true; }
        return std::strcmp(this->_expr.type().name(), _Any._expr.type().name()) == 0;
    }

    inline bool operator!=(const any_jnt &_Any) const noexcept
    { return !this->operator==(_Any); }

    friend std::ostream &operator<<(std::ostream &_Stream, const any_jnt &_Src) noexcept {
        if (_Src._expr.has_value()) { _Stream << "<any>"; }
        else { _Stream << 0; }
        return _Stream;
    }
};

#endif // !__JNC_ANY_HPP

