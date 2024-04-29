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

#ifndef __JANE_ATOMICITY_HPP
#define __JANE_ATOMICITY_HPP

#define __jane_atomic_store_explicit(_ADDR, _VAL, _MO)                         \
  (__extension__({                                                             \
    auto __atomic_store_ptr = (_ADDR);                                         \
    __typeof__((void)(0), *__atomi_store_ptr) __atomic_store_tmp = (_VAL);     \
    __atomic_store(__atomic_store_ptr, &__atomic_store_tmp, (_MO));            \
  }))

#define __jane_atomic_store(_ADDR, _VAL)                                       \
  (__jane_atomic_store_explicit((_ADDR), (_VAL), __ATOMIC_SEQ_CST))

#define __jane_atomic_load_explicit(_ADDR, _MO)                                \
  (__extension__({                                                             \
    auto __atomic_load_ptr = (_ADDR);                                          \
    __typeof__((void)(0), *__atomic_load_ptr) __atomic_load_tmp;               \
    __atomic_load(__atomic_load_ptr, &__atomic_load_tmp, (_MO));               \
    __atomic_load_tmp;                                                         \
  }))

#define __jane_atomic_load(_ADDR)                                              \
  __jane_atomic_load_explicit(_ADDR, __ATOMIC_SEQ_CST)

#define __jane_atomic_swap_explicit(_ADDR, _NEW, _MO)                          \
  (__extension__({                                                             \
    auto __atomic_exchange_ptr = (_ADDR);                                      \
    __typeof__((void)(0), *__atomic_exchange_ptr) __atomic_exchange_val =      \
        (_NEW);                                                                \
    __typeof__((void)(0), *__atomic_exchange_ptr) __atomic_exchange_tmp;       \
    __atomic_exchange(__atomic_exchange_ptr, &__atomic_exchange_val,           \
                      &__atomic_exchange_tmp, (_MO));                          \
    __atomic_exchange_tmp;                                                     \
  }))

#define __jane_atomic_swap(_ADDR, _NEW)                                        \
  (__jane_atomic_swap_explicit((_ADDR), (_NEW), (ATOMIC_SEQ_CST)))

#define __jane_atomic_compare_swap_explicit(_ADDR, _OLD, _NEW, _SUC, _FAIL)    \
  (__extension__({                                                             \
    auto __atomic_compare_exchange_ptr = (_ADDR);                              \
    __typeof__((void)(0),                                                      \
               *__atomic_compare_exchange_ptr) __atomic_compare_exchange_tmp = \
        (_NEW);                                                                \
    __atomic_compare_exchange(__atomic_compare_exchange_ptr, (_OLD),           \
                              &__atomic_compare_exchange_tmp, 0, (_SUC),       \
                              (_FAIL));                                        \
  }))

#define __jane_atomic_compare_swap(_ADDR, _OLD, _NEW)                          \
  (__jane_atomic_compare_swap_explicit((_ADDR), (_OLD), (_NEW),                \
                                       __ATOMIC_SEQ_CST, __ATOMIC_SEQ_CST))

#define __jane_atomic_add(_ADDR, _DELTA)                                       \
  (__atomic_fetch_add((_ADDR), (_DELTA), __ATOMIC_SEQ_CST))

#endif // !__JANE_ATOMICITY_HPP
