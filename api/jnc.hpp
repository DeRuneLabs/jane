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

#ifndef __JNC_HPP
#define __JNC_HPP

#include "any.hpp"
#include "array.hpp"
#include "atomicity.hpp"
#include "builtin.hpp"
#include "defer.hpp"
#include "fn.hpp"
#include "jn_util.hpp"
#include "map.hpp"
#include "ref.hpp"
#include "slice.hpp"
#include "str.hpp"
#include "trait.hpp"
#include "typedef.hpp"
#include "utf16.hpp"
#include "utf8.hpp"
#include <cstdlib>
#include <exception>

slice<str_jnt> __jnc_command_line_args;

inline slice<str_jnt> __jnc_get_command_line_args(void) noexcept;
inline void JNC_ID(panic)(const trait<JNC_ID(Error)> &_Error);

template <typename Type, unsigned N, unsigned Last> struct tuple_ostream;
template <typename... Types>
std::ostream &operator<<(std::ostream &_Stream,
                         const std::tuple<Types...> &_Tuple);
template <typename _Fn_t, typename _Tuple_t, size_t... _I_t>
inline auto tuple_as_args(const fn<_Fn_t> &_Function, const _Tuple_t _Tuple,
                          const std::index_sequence<_I_t...>);

template <typename _Fn_t, typename _Tuple_t>
inline auto tuple_as_args(const fn<_Fn_t> &_Function, const _Tuple_t _Tuple);
template <typename T> inline jn_ref<T> __jnc_new_structure(T *_Ptr);

template <typename _Obj_t> str_jnt __jnc_to_str(const _Obj_t &_Obj) noexcept;
void __jnc_terminate_handler(void) noexcept;

// entry point function ganerated code
void JNC_ID(main)(void);
// initialize call function
void __jnc_call_package_initializers(void);
void __jnc_setup_command_line_args(int argc, char *argv[]) noexcept;

int main(int argc, char *argv[]);

inline slice<str_jnt> __jnc_get_command_line_args(void) noexcept {
  return __jnc_command_line_args;
}

inline std::ostream &operator<<(std::ostream &_Stream,
                                const signed char _I8) noexcept {
  return _Stream << ((int)(_I8));
}

inline std::ostream &operator<<(std::ostream &_Stream,
                                const unsigned char _U8) noexcept {
  return _Stream << ((int)(_U8));
}

template <typename Type, unsigned N, unsigned Last> struct tuple_ostream {
  static void arrow(std::ostream &_Stream, const Type &_Type) {
    _Stream << std::get<N>(_Type) << ", ";
    tuple_ostream<Type, N + 1, Last>::arrow(_Stream, _Type);
  }
};

template <typename Type, unsigned N> struct tuple_ostream<Type, N, N> {
  static void arrow(std::ostream &_Stream, const Type &_Type) {
    _Stream << std::get<N>(_Type);
  }
};

template <typename... Types>
std::ostream &operator<<(std::ostream &_Stream,
                         const std::tuple<Types...> &_Tuple) {
  _Stream << '(';
  tuple_ostream<std::tuple<Types...>, 0, sizeof...(Types) - 1>::arrow(_Stream,
                                                                      _Tuple);
  _Stream << ')';
  return _Stream;
}

template <typename _Fn_t, typename _Tuple_t, size_t... _I_t>
inline auto tuple_as_args(const fn<_Fn_t> &_Function, const _Tuple_t _Tuple) {
  static constexpr auto _size{std::tuple_size<_Tuple_t>::value};
  return tuple_as_args(_Function, _Tuple, std::make_index_sequence<_size>{});
}

template <typename T> inline jn_ref<T> __jnc_new_structure(T *_Ptr) {
  if (!_Ptr) {
    JNC_ID(panic)(__JNC_ERROR_MEMORY_ALLOCATION_FAILED);
  }
  _Ptr->self._ref = new (std::nothrow) uint_jnt;
  if (!_Ptr->self._ref) {
    JNC_ID(panic)(__JNC_ERROR_MEMORY_ALLOCATION_FAILED);
  }
  *_Ptr->self._ref = 0;
  return (_Ptr->self);
}

template <typename _Obj_t> str_jnt __jnc_to_str(const _Obj_t &_Obj) noexcept {
  std::stringstream _stream;
  _stream << _Obj;
  /* return (str_jnt(_stream.str())); */
}

inline void JNC_ID(panic)(const trait<JNC_ID(Error)> &_Error) { throw(_Error); }

template <typename _Obj_t> void JNC_ID(panic)(const _Obj_t &_Expr) {
  struct panic_error : public JNC_ID(Error) {
    str_jnt _message;
    str_jnt error(void) { return (this->_message); }
  };
  struct panic_error _error;
  _error._message = __jnc_to_str(_Expr);
  throw(trait<JNC_ID(Error)>(_error));
}

void __jnc_terminate_handler(void) noexcept {
  try {
    std::rethrow_exception(std::current_exception());
  } catch (trait<JNC_ID(Error)> _error) {
    // std::cout << "panic: " << _error.get().error() << std::endl;
    std::exit(__JNC_EXIT_PANIC);
  }
}

void __jnc_setup_command_line_args(int argc, char *argv[]) noexcept {
#ifdef _WINDOWS
  const LPWSTR _cmdl{GetCommandLineW()};
  wchar_t *_wargs{_cmdl};
  const size_t _wargs_len{std::wcslen(_wargs)};
  slice<str_jnt> _args;
  int_jnt _old{0};
  for (int_jnt _i{0}; _i < _wargs_len; ++_i) {
    const wchar_t _r{_wargs[_i]};
    if (!std::iswspace(_r)) {
      continue;
    } else if (_i + 1 < _wargs_len && std::iswspace(_wargs[_i + 1])) {
      continue;
    }
    _wargs[_i] = 0;
    wchar_t *_warg{_wargs + _old};
    _old = _i + 1;
    _args.__push(__jnc_utf16_to_utf8_str(_warg, std::wcslen(_warg)));
  }
  _args.__push(
      __jnc_utf16_to_utf8_str(_wargs + _old, std::wcslen(_wargs + _old)));
  __jnc_command_line_args = _args;
#else
  // __jnc_args = slice<str_jnt>(argc);
  // for (int_jnt _i{0}; _i < argc; ++_i) {
  //   __jnc_command_line_args[_i] = argv[_i];
  // }
#endif // _WINDOWS
}

int main(int argc, char *argv[]) {
#ifdef _WINDOWS
  SetConsoleOutputCP(CP_UTF8);
  _setmode(_fileno(stdin), (0x00020000));
#else
  std::set_terminate(&__jnc_terminate_handler);
  __jnc_setup_command_line_args(argc, argv);
  __jnc_call_package_initializers();
  JNC_ID(main());

  return (EXIT_SUCCESS);
#endif // _WINDOWS
}

#endif // !__JNC_HPP
