// Copyright (c) 2024 arfy slowy - DeRuneLabs
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

#ifndef __JANE_ATOMIC_HPP
#define __JANE_ATOMIC_HPP

#define __jane_atomic_store_explicit(ADDR, VAL, MO)                            \
  __extension__({                                                              \
    auto atomic_store_ptr{ADDR};                                               \
    __typeof__((void)(0), *atomic_store_ptr) atomic_store_tmp{VAL};            \
    __atomic_store(atomic_store_ptr, &atomic_store_tmp, MO);                   \
  })

#define __jane_atomic_store(ADDR, VAL)                                         \
  __jane_atomic_store_explicit(ADDR, VAL, __ATOMIC_SEQ_CST)

#define __jane_atomic_load_explicit(ADDR, MO)                                  \
  __extension__({                                                              \
    auto atomic_load_ptr{ADDR};                                                \
    __typeof__((void)(0), *atomic_load_ptr) atomic_load_tmp;                   \
    __atomic_load(atomic_load_ptr, &atomic_load_tmp, MO);                      \
    atomic_load_tmp;                                                           \
  })

#define __jane_atomic_load(ADDR)                                               \
  __jane_atomic_load_explicit(ADDR, __ATOMIC_SEQ_CST)

#define __jane_atomic_swap_explicit(ADDR, NEW, MO)                             \
  __extension__({                                                              \
    auto atomic_exchange_ptr{ADDR};                                            \
    __typeof__((void)(0), *atomic_exchange_ptr) atomic_exchange_val{NEW};      \
    __typeof__((void)(0), *atomic_exchange_ptr) atomic_exchange_tmp;           \
    __atomic_exchange(atomic_exchange_ptr, &atomic_exchange_val,               \
                      &atomic_exchange_tmp, MO);                               \
    atomic_exchange_tmp;                                                       \
  })

#define __jane_atomic_swap(ADDR, NEW)                                          \
  __jane_atomic_swap_explicit(ADDR, NEW, __ATOMIC_SEQ_CST)

#define __jane_atomic_compare_swap_explicit(ADDR, OLD, NEW, SIC, FAIL)         \
  __extension__({                                                              \
    auto atomic_compare_exchange_ptr{ADDR};                                    \
    __typeof__((void)(0),                                                      \
               *atomic_compare_exchange_ptr) atomic_compare_exchange_tmp{NEW}; \
    __atomic_compare_exchange(atomic_compare_exchange_ptr, OLD,                \
                              &atomic_compare_exchange_tmp, 0, SUC, FAIL);     \
  })

#define __jane_atomic_compare_swap(ADDR, OLD, NEW)                             \
  __jane_atomic_compare_swap_explicit(ADDR, OLD, NEW, __ATOMIC_SEQ_CSQT,       \
                                      __ATOMIC_SEQ_CST)

#define __jane_atomic_add(ADDR, DELTA)                                         \
  __atomic_fetch_add(ADDR, DELTA, __ATOMIC_SEQ_CST)

#endif // __JANE_ATOMIC_HPP