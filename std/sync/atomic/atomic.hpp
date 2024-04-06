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

#include "../../../api/ptr.hpp"
#include "../../../api/typedef.hpp"

#ifndef __JNC_STD_SYNC_ATOMIC_ATOMIC_HPP
#define __JNC_STD_SYNC_ATOMIC_ATOMIC_HPP

#define __jnc_atomic_store_explicit(ADDR, VAL, MO)                             \
  __extension__({                                                              \
    auto __atomic_store_ptr = (ADDR);                                          \
    __typeof__((void)0, *__atomic_store_ptr) __atomic_store_tmp = (VAL);       \
    __atomic_store(__atomic_store_ptr, &__atomic_store_tmp, (MO));             \
  })

#define __jnc_atomic_store(ADDR, VAL)                                          \
  __jnc_atomic_store_explicit(ADDR, VAL, __ATOMIC_SEQ_CST)

#define __jnc_atomic_load_explicit(ADDR, MO)                                          \
  __extension__({                                                              \
    auto __atomic_load_ptr = (ADDR);                                            \
    __typeof__((void)0, *__atomic_load_ptr) __atomic_load_tmp;                 \
    __atomic_load(__atomic_load_ptr, &__atomic_load_tmp, (MO));                \
    __atomic_load_tmp;                                                         \
  })

#define __jnc_atomic_load(ADDR) __jnc_atomic_load_explicit(ADDR, __ATOMIC_SEQ_CST)

#define __jnc_atomic_swap_explicit(ADDR, NEW, MO)                              \
  __extension__({                                                              \
    auto __atomic_exchange_ptr = (ADDR);                                       \
    __typeof__((void)0, *__atomic_exchange_ptr) __atomic_exchange_val = (NEW); \
    __typeof__((void)0, *__atomic_exchange_ptr) __atomic_exchange_tmp;         \
    __atomic_exchange(__atomic_exchange_ptr, &__atomic_exchange_val,           \
                      &__atomic_exchange_tmp, (MO));                           \
    __atomic_exchange_tmp;                                                     \
  })

#define __jnc_atomic_swap(ADDR, NEW)                                           \
  __jnc_atomic_swap_explicit(ADDR, NEW, __ATOMIC_SEQ_CST)

#define __jnc_atomic_compare_swap_explicit(ADDR, OLD, NEW, SUC, FAIL)          \
  __extension__({                                                              \
    auto __atomic_compare_exchange_ptr = (ADDR);                               \
    __typeof__((void)0,                                                        \
               *__atomic_compare_exchange_ptr) __atomic_compare_exchange_tmp = \
        (NEW);                                                                 \
    __atomic_compare_exchange(__atomic_compare_exchange_ptr, (OLD),            \
                              &__atomic_compare_exchange_tmp, 0, (SUC),        \
                              (FAIL));                                         \
  })

#define __jnc_atomic_compare_swap(ADDR, OLD, NEW)                              \
  __jnc_atomic_compare_swap_explicit(ADDR, OLD, NEW, __ATOMIC_SEQ_CST,         \
                                     __ATOMIC_SEQ_CST)

#define __jnc_atomic_add(ADDR, DELTA)                                          \
  __atomic_fetch_add((ADDR), (DELTA), __ATOMIC_SEQ_CST)

inline i32_jnt __jnc_atomic_swap_i32(const ptr<i32_jnt> &_Addr,
                                     const i32_jnt &_New) noexcept;
inline i64_jnt __jnc_atomic_swap_i64(const ptr<i64_jnt> &_Addr,
                                     const i64_jnt &_New) noexcept;

inline u32_jnt __jnc_atomic_swap_u32(const ptr<u32_jnt> &_Addr,
                                     const u32_jnt &_New) noexcept;
inline u64_jnt __jnc_atomic_swap_u64(const ptr<u64_jnt> &_Addr,
                                     const u64_jnt &_New) noexcept;

inline uintptr_jnt __jnc_atomic_swap_uintptr(const ptr<uintptr_jnt> &_Addr,
                                             const uintptr_jnt &_New) noexcept;

inline bool __jnc_atomic_compare_swap_i32(const ptr<i32_jnt> &_Addr,
                                          const i32_jnt &_Old,
                                          const i32_jnt &_New) noexcept;
inline bool __jnc_atomic_compare_swap_i64(const ptr<i64_jnt> &_Addr,
                                          const i64_jnt &_Old,
                                          const i64_jnt &_New) noexcept;

inline bool __jnc_atomic_compare_swap_u64(const ptr<u64_jnt> &_Addr,
                                          const u64_jnt &_Old,
                                          const u64_jnt &_New) noexcept;
inline bool __jnc_atomic_compare_swap_uintptr(const ptr<uintptr_jnt> &_Addr,
                                              const uintptr_jnt &_Old,
                                              const uintptr_jnt &_New) noexcept;

inline i32_jnt __jnc_atomic_add_i32(const ptr<i32_jnt> &_Addr,
                                    const i32_jnt &_Delta) noexcept;
inline i64_jnt __jcn_atomic_add_i64(const ptr<i64_jnt> &_Addr,
                                    const i64_jnt &_Delta) noexcept;

inline u32_jnt __jnc_atomic_add_u32(const ptr<u32_jnt> &_Addr,
                                    const u32_jnt &_Delta) noexcept;
inline u64_jnt __jnc_atomic_add_u64(const ptr<u64_jnt> &_Addr,
                                    const u64_jnt &_Delta) noexcept;

inline uintptr_jnt __jnc_atomic_add_uintptr(const ptr<uintptr_jnt> &_Addr,
                                            const uintptr_jnt &_Delta) noexcept;

inline i32_jnt __jnc_atomic_load_i32(const ptr<i32_jnt> &_Addr) noexcept;
inline i64_jnt __jnc_atomic_load_i64(const ptr<i64_jnt> &_Addr) noexcept;
inline u32_jnt __jnc_atomic_load_u32(const ptr<u32_jnt> &_Addr) noexcept;
inline u64_jnt __jnc_atomic_load_u64(const ptr<u64_jnt> &_Addr) noexcept;

inline uintptr_jnt
__jnc_atomic_load_uintptr(const ptr<uintptr_jnt> &_Addr) noexcept;

inline void __jnc_atomic_store_i32(const ptr<i32_jnt> &_Addr,
                                   const i32_jnt &_Val) noexcept;
inline void __jnc_atomic_store_i64(const ptr<i64_jnt> &_Addr,
                                   const i32_jnt &_Val) noexcept;
inline void __jnc_atomic_store_u32(const ptr<u32_jnt> &_Addr,
                                   const u32_jnt &_Val) noexcept;
inline void __jnc_atomic_store_u64(const ptr<u64_jnt> &_Addr,
                                   const u64_jnt &_Val) noexcept;
inline void __jnc_atomic_store_uintptr(const ptr<uintptr_jnt> &_Addr,
                                       const uintptr_jnt &_Val) noexcept;

inline i32_jnt __jnc_atomic_swap_i32(const ptr<i32_jnt> &_Addr,
                                     const i32_jnt &_New) noexcept {
  return __jnc_atomic_swap(_Addr._ptr, _New);
}

inline i64_jnt __jnc_atomic_swap_i64(const ptr<i64_jnt> &_Addr,
                                     const i64_jnt &_New) noexcept {
  return __jnc_atomic_swap(_Addr._ptr, _New);
}

inline u32_jnt __jnc_atomic_swap_u32(const ptr<u32_jnt> &_Addr,
                                     const u32_jnt &_New) noexcept {
  return __jnc_atomic_swap(_Addr._ptr, _New);
}

inline u64_jnt __jnc_atomic_swap_u64(const ptr<u64_jnt> &_Addr,
                                     const u64_jnt &_New) noexcept {
  return __jnc_atomic_swap(_Addr._ptr, _New);
}

inline uintptr_jnt __jnc_atomic_swap_uintptr(const ptr<uintptr_jnt> &_Addr,
                                             const uintptr_jnt &_New) noexcept {
  return __jnc_atomic_swap(_Addr._ptr, _New);
}

inline bool __jnc_atomic_compare_swap_i32(const ptr<i32_jnt> &_Addr,
                                          const i32_jnt &_Old,
                                          const i32_jnt &_New) noexcept {
  return __jnc_atomic_compare_swap((i32_jnt *)(_Addr._ptr), (i32_jnt *)(&_Old),
                                   _New);
}

inline bool __jnc_atomic_compare_swap_i64(const ptr<i64_jnt> &_Addr,
                                          const i64_jnt &_Old,
                                          const i64_jnt &_New) noexcept {
  return __jnc_atomic_compare_swap((i64_jnt *)(_Addr._ptr), (i64_jnt *)(&_Old),
                                   _New);
}

inline bool __jnc_atomic_compare_swap_u32(const ptr<u32_jnt> &_Addr,
                                          const u32_jnt &_Old,
                                          const u32_jnt &_New) noexcept {
  return __jnc_atomic_compare_swap((u32_jnt *)(_Addr._ptr), (u32_jnt *)(&_Old),
                                   _New);
}

inline bool __jnc_atomic_compare_swap_u64(const ptr<u64_jnt> &_Addr,
                                          const u64_jnt &_Old,
                                          const u64_jnt &_New) noexcept {
  return __jnc_atomic_compare_swap((u64_jnt *)(_Addr._ptr), (u64_jnt *)(&_Old),
                                   _New);
}

inline bool
__jnc_atomic_compare_swap_uintptr(const ptr<uintptr_jnt> &_Addr,
                                  const uintptr_jnt &_Old,
                                  const uintptr_jnt &_New) noexcept {
  return __jnc_atomic_compare_swap((uintptr_jnt *)(_Addr._ptr),
                                   (uintptr_jnt *)(&_Old), _New);
}

inline i32_jnt __jnc_atomic_add_i32(const ptr<i32_jnt> &_Addr,
                                    const i32_jnt &_Delta) noexcept {
  return __jnc_atomic_add(_Addr._ptr, _Delta);
}

inline i64_jnt __jnc_atomic_add_i64(const ptr<i64_jnt> &_Addr,
                                    const i64_jnt &_Delta) noexcept {
  return __jnc_atomic_add(_Addr._ptr, _Delta);
}

inline u32_jnt __jnc_atomic_add_u32(const ptr<u32_jnt> &_Addr,
                                    const u32_jnt &_Delta) noexcept {
  return __jnc_atomic_add(_Addr._ptr, _Delta);
}

inline u64_jnt __jnc_atomic_add_u64(const ptr<u64_jnt> &_Addr,
                                    const u64_jnt &_Delta) noexcept {
  return __jnc_atomic_add(_Addr._ptr, _Delta);
}

inline uintptr_jnt
__jnc_atomic_add_uintptr(const ptr<uintptr_jnt> &_Addr,
                         const uintptr_jnt &_Delta) noexcept {
  return __jnc_atomic_add(_Addr._ptr, _Delta);
}

inline i32_jnt __jnc_atomic_load_i32(const ptr<i32_jnt> &_Addr) noexcept {
  return __jnc_atomic_load(_Addr._ptr);
}

inline i64_jnt __jnc_atomic_load_i64(const ptr<i64_jnt> &_Addr) noexcept {
  return __jnc_atomic_load(_Addr._ptr);
}

inline u32_jnt __jnc_atomic_load_u32(const ptr<u32_jnt> &_Addr) noexcept {
  return __jnc_atomic_load(_Addr._ptr);
}

inline u64_jnt __jnc_atomic_load_u64(const ptr<u64_jnt> &_Addr) noexcept {
  return __jnc_atomic_load(_Addr._ptr);
}

inline uintptr_jnt
__jnc_atomic_load_uintptr(const ptr<uintptr_jnt> &_Addr) noexcept {
  return __jnc_atomic_load(_Addr._ptr);
}

inline void __jnc_atomic_store_i32(const ptr<i32_jnt> &_Addr,
                                   const i32_jnt &_Val) noexcept {
  __jnc_atomic_store(_Addr._ptr, _Val);
}

inline void __jnc_atomic_store_i64(const ptr<i64_jnt> &_Addr,
                                   const i64_jnt &_Val) noexcept {
  __jnc_atomic_store(_Addr._ptr, _Val);
}

inline void __jnc_atomic_store_u32(const ptr<u32_jnt> &_Addr,
                                   const u32_jnt &_Val) noexcept {
  __jnc_atomic_store(_Addr._ptr, _Val);
}

inline void __jnc_atomic_store_u64(const ptr<u64_jnt> &_Addr,
                                   const u64_jnt &_Val) noexcept {
  __jnc_atomic_store(_Addr._ptr, _Val);
}

inline void __jnc_atomic_store_uintptr(const ptr<uintptr_jnt> &_Addr,
                                       const uintptr_jnt &_Val) noexcept {
  __jnc_atomic_store(_Addr._ptr, _Val);
}

#endif // !__JNC_STD_SYNC_ATOMIC_ATOMIC_HPP
