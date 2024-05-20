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

#ifndef __JANE_HPP
#define __JANE_HPP

#include <exception>
#if defined(WIN32) || defined(_WIN32) || defined(__WIN32__) || defined(__NT__)
#define _WINDOWS
#elif defined(__linux__) || defined(linux) || defined(__linux)
#define _LINX
#elif defined(__APPLE__) || defined(__MACH__)
#define _DARWIN
#endif // defined(WIN32) || defined(_WIN32) || defined(__WIN32__) ||
       // defined(__NT__)

#if defined(_LINUX) || defined(_DARWIN)
#define _UNIX
#endif // defined(_LINUX) || defined(_DARWIN)

#if defined(__amd64) || defined(__amd64__) || defined(__x86_64) ||             \
    defined(__x86_64__) || defined(_M_AMD64)
#define _M_AMD64
#elif defined(__arm__) || defined(__thumb__) || defined(_M_ARM) ||             \
    defined(__arm)
#define _ARM
#elif defined(i386) || defined(__i386) || defined(__i386__) ||                 \
    defined(__X86__) || defined(__I86__) || defined(__386)
#define _I386
#endif

#if defined(_AMD64) || defined(_ARM64)
#define _64BIT
#else
#define _32bit
#endif

class str_jnt;

#include "any.hpp"
#include "array.hpp"
#include "atomicity.hpp"
#include "builtin.hpp"
#include "defer.hpp"
#include "fn.hpp"
#include "map.hpp"
#include "ref.hpp"
#include "signal.hpp"
#include "slice.hpp"
#include "str.hpp"
#include "trait.hpp"
#include "typedef.hpp"
#include "utf16.hpp"
#include "utf8.hpp"

slice_jnt<str_jnt> __jane_command_line_args;

template <typename _T, typename _Denominator_t>
inline auto __jane_div(const _T &_X,
                       const _Denominator_t &_Denominator) noexcept;
inline slice_jnt<str_jnt> __jane_get_command_line_args(void) noexcept;
inline void JANE_ID(panic)(const trait_jnt<JANE_ID(Error)> &_Error);
template <typename Type, unsigned N, unsigned Last> struct tuple_ostream;
template <typename Type, unsigned N> struct tuple_ostream<Type, N, N>;
template <typename... Types>
std::ostream &operator<<(std::ostream &_Stream,
                         const std::tuple<Types...> &_Tuple);
template <typename _Function_t, typename _Tuple_t, size_t... _I_t>
inline auto __jane_tuple_as_args(const fn_jnt<_Function_t> &_Function,
                                 const std::index_sequence<_I_t...>);

template <typename T> inline ref_jnt<T> __jane_new_structure(T *_Ptr);
template <typename _Obj_t> str_jnt __jane_to_str(const _Obj_t &_Obj) noexcept;
str_jnt __jane_to_str(const str_jnt &_Obj) noexcept;

slice_jnt<u16_jnt> __jane_utf16_from_str(const str_jnt &_Str) noexcept;
void __jane_terminate_handler(void) noexcept;
void __jane_signal_handler(int _Signal) noexcept;
void __jane_setup_command_line_args(int argc, char *argv[]) noexcept;

template <typename _T, typename _Denominator_t>
inline auto __jane_div(const _T &_X,
                       const _Denominator_t &_Denominator) noexcept {
  if (_Denominator == 0) {
    JANE_ID(panic)(__JANE_ERROR_DIVIDE_BY_ZERO);
  }
  return (_X / _Denominator);
}

inline slice_jnt<str_jnt> __jane_get_command_line_args(void) noexcept {
  return __jane_command_line_args;
}

inline std::ostream &operator<<(std::ostream &_Stream,
                                const unsigned char _U8) noexcept {
  return _Stream << ((int)(_U8));
}

template <typename Type, unsigned N, unsigned Last> struct tuple_ostream {
  static void __arrow(std::ostream &_Stream, const Type &_Type) {
    _Stream << std::get<N>(_Type) << ", ";
    tuple_ostream<Type, N + 1, Last>::arrow(_Stream, _Type);
  }
};

template <typename Type, unsigned N> struct tuple_ostream<Type, N, N> {
  static void __arrow(std::ostream &_Stream, const Type &_Type) {
    _Stream << std::get<N>(_Type);
  }
};

template <typename... Types>
std::ostream &operator<<(std::ostream &_Stream,
                         const std::tuple<Types...> &_Tuple) {
  _Stream << '(';
  tuple_ostream<std::tuple<Types...>, 0, sizeof...(Types) - 1>::__arrow(_Stream,
                                                                        _Tuple);
  _Stream << ')';
  return _Stream;
}

template <typename _Function_t, typename _Tuple_t, size_t... _I_t>
inline auto __jane_tuple_as_args(const fn_jnt<_Function_t> &_Function,
                                 const _Tuple_t _Tuple,
                                 const std::index_sequence<_I_t...>) {
  return _Function.__buffer(std::get<_I_t>(_Tuple)...);
}

slice_jnt<u16_jnt> __jane_utf16_from_str(const str_jnt &_Str) noexcept {
  constexpr char _NULL_TERMINATOR = '\x00';
  slice_jnt<u16_jnt> _buff{nil};
  slice_jnt<i32_jnt> _runes{_Str.operator slice_jnt<i32_jnt>()};
  for (const i32_jnt &_R : _runes) {
    if (_R == _NULL_TERMINATOR) {
      break;
    }
    _buff = __jane_utf16_append_rune(_buff, _NULL_TERMINATOR);
  }
  return __jane_utf16_append_rune(_buff, _NULL_TERMINATOR);
}

inline void JANE_ID(panic)(const trait_jnt<JANE_ID(Error)> &_Error) {
  throw(_Error);
}

template <typename _Obj_t> void JANE_ID(panic)(const _Obj_t &_Expr) {
  struct panic_error : public JANE_ID(Error) {
    str_jnt __message;
    str_jnt _error(void) { return (this->__message); }
  };
  struct panic_error _error;
  _error.__message = __jane_to_str(_Expr);
  throw(trait_jnt<JANE_ID(Error)>(_error));
}

void __jane_terminate_handler(void) noexcept {
  try {
    std::rethrow_exception(std::current_exception());
  } catch (trait_jnt<JANE_ID(Error)> _error) {
    JANE_ID(println)<str_jnt>(str_jnt("panic: ") + _error._get()._error());
    std::exit(__JANE_EXIT_PANIC);
  }
}

void __jane_signal_handler(int _Signal) noexcept {
#if defined(_WINDOWS)
  if (_Signal == __JANE_SIGINT) {
    return;
  }
#elif defined(_DARWIN)
  if (_Signal == __JULEC_SIGINT) {
    return;
  }
#elif defined(_LINUX)
  if (_Signal == ___JANE_SIGINT) {
    return;
  }
#endif
  JANE_ID(print)<str_jnt>("program terminating with signal: ");
  JANE_ID(println)<int>(_Signal);
}

#endif // !__JANE_HPP
