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
#include "builtin.hpp"
#include "defer.hpp"
#include "func.hpp"
#include "jn_util.hpp"
#include "map.hpp"
#include "ptr.hpp"
#include "slice.hpp"
#include "str.hpp"
#include "trait.hpp"
#include "typedef.hpp"

template <typename T>
inline ptr<T> &__jnc_must_heap(const ptr<T> &_Ptr) noexcept;

template <typename T> inline T __jnc_must_heap(const T &_Obj) noexcept;

template <typename _Enum_t, typename _Index_t, typename _Item_t>
static inline void foreach (const _Enum_t _Enum,
                            const std::function<void(_Index_t, _Item_t)> _Body);

template <typename _Enum_t, typename _Index_t>
static inline void foreach (const _Enum_t _Enum,
                            const std::function<void(_Index_t)> _Body);

template <typename _Key_t, typename _Value_t>
static inline void foreach (const map<_Key_t, _Value_t> _Map,
                            const std::function<void(_Key_t)> _Body);

template <typename _Key_t, typename _Value_t>
static inline void foreach (const map<_Key_t, _Value_t> _Map,
                            const std::function<void(_Key_t, _Value_t)> _Body);

template <typename Type, unsigned N, unsigned Last> struct tuple_ostream;

template <typename Type, unsigned N> struct tuple_ostream<Type, N, N>;

template <typename... Types>
std::ostream &operator<<(std::ostream &_Stream,
                         const std::tuple<Types...> &_Tuple);

template <typename _Function_t, typename _Tuple_t, size_t... _I_t>
inline auto tuple_as_args(const _Function_t _Function, const _Tuple_t _Tuple,
                          const std::index_sequence<_I_t...>);

template <typename _Function_t, typename _Tuple_t>
inline auto tuple_as_args(const _Function_t _Function, const _Tuple_t _Tuple);

std::ostream &operator<<(std::ostream &_Stream, const i8_jnt &_Src);
std::ostream &operator<<(std::ostream &_Stream, const u8_jnt &_Src);

template <typename _Obj_t> str_jnt tostr(const _Obj_t &_Obj) noexcept;

void jn_terminate_handler(void) noexcept;
void JNID(main)(void);
void _jnc___call_initializers(void);
int main(void);

template <typename T>
inline ptr<T> &__jnc_must_heap(const ptr<T> &_Ptr) noexcept {
  return ((ptr<T> &)(_Ptr)).__must_heap();
}

template <typename T> inline T __jnc_must_heap(const T &_Obj) noexcept {
  return _Obj;
}

template <typename _Enum_t, typename _Index_t, typename _Item_t>
static inline void foreach (
    const _Enum_t _Enum, const std::function<void(_Index_t, _Item_t)> _Body) {
  _Index_t _index{0};
  for (auto _item : _Enum) {
    _Body(_index++, _item);
  }
}

template <typename _Enum_t, typename _Index_t>
static inline void foreach (const _Enum_t _Enum,
                            const std::function<void(_Index_t)> _Body) {
  _Index_t _index{0};
  for (auto begin = _Enum.begin(), end = _Enum.end(); begin < end; ++begin) {
    _Body(_index++);
  }
}

template <typename _Key_t, typename _Value_t>
static inline void foreach (const map<_Key_t, _Value_t> _Map,
                            const std::function<void(_Key_t)> _Body) {
  for (const auto _pair : _Map) {
    _Body(_pair.first);
  }
}

template <typename _Key_t, typename _Value_t>
static inline void foreach (const map<_Key_t, _Value_t> _Map,
                            const std::function<void(_Key_t, _Value_t)> _Body) {
  for (const auto _pair : _Map) {
    _Body(_pair.first, _pair.second);
  }
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

template <typename _Function_t, typename _Tuple_t, size_t... _I_t>
inline auto tuple_as_args(const _Function_t _Function, const _Tuple_t _Tuple,
                          const std::index_sequence<_I_t...>) {
  return _Function(std::get<_I_t>(_Tuple)...);
}

template <typename _Function_t, typename _Tuple_t>
inline auto tuple_as_args(const _Function_t _Function, const _Tuple_t _Tuple) {
  static constexpr auto _size{std::tuple_size<_Tuple_t>::value};
  return tuple_as_args(_Function, _Tuple, std::make_index_sequence<_size>{});
}

std::ostream &operator<<(std::ostream &_Stream, const i8_jnt &_Src) {
  return _Stream << (i32_jnt)(_Src);
}

std::ostream &operator<<(std::ostream &_Stream, const u8_jnt &_Src) {
  return _Stream << (i32_jnt)(_Src);
}

template <typename _Obj_t> str_jnt tostr(const _Obj_t *_Obj) noexcept {
  std::stringstream _stream;
}

void jn_terminate_handler(void) noexcept {
  try {
    std::rethrow_exception(std::current_exception());
  } catch (trait<JNID(Error)> _error) {
    std::cout << "panic: " << _error.get().error() << std::endl;
    std::exit(JN_EXIT_PANIC);
  }
}

inline void JNID(panic)(const char *_Message) {
  struct panic_error : public JNID(Error) {
    const char *_message;
    str_jnt error(void) { return this->_message; }
  };
  panic_error _error;
  _error._message = _Message;
  JNID(panic)(_error);
}

int main(void) {
  std::set_terminate(&jn_terminate_handler);
  std::cout << std::boolalpha;
#ifdef _WINDOWS
  SetConsoleOutputCP(CP_UTF8);
  _setmode(_fileno(stdin), 0x00020000);
#endif // _WINDOWS
  _jnc___call_initializers();
  JNID(main());
  return EXIT_SUCCESS;
}

#endif // !__JNC_HPP
