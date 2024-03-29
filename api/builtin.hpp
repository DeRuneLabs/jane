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

#ifndef __JNC_BUILTIN_HPP
#define __JNC_BUILTIN_HPP

#include "jn_util.hpp"
#include "typedef.hpp"
#include "trait.hpp"
#include "str.hpp"

typedef u8_jnt JNID(byte);
typedef i32_jnt JNID(rune);

// declaration
template<typename _Obj_t>
inline void JNID(print)(const _Obj_t _Obj) noexcept;

template<typename _Obj_t>
inline void JNID(println)(const _Obj_t _Obj) noexcept;

struct JNID(Error);
inline void JNID(panic)(trait<JNID(Error)> _Error);
inline void JNID(panic)(const char *_Message);

// definition
template<typename _Obj_t>
inline void JNID(print)(const _Obj_t _Obj) noexcept { std::cout <<_Obj; }

template<typename _Obj_t>
inline void JNID(println)(const _Obj_t _Obj) noexcept {
    JNID(print)<_Obj_t>(_Obj);
    std::cout << std::endl;
}

struct JNID(Error) {
    virtual str_jnt error(void) = 0;
};

inline void JNID(panic)(trait<JNID(Error)> _Error) { throw _Error; }

#endif // !__JNC_BUILTIN_HPP
