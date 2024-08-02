// MIT License
//
// # Copyright (c) 2024 arfy slowy - DeRuneLabs
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

package types

func Is_i8_compatible(k string) bool {
	k = Real_kind_of(k)
	return k == TypeKind_I8
}

func Is_i16_compatible(k string) bool {
	k = Real_kind_of(k)
	return k == TypeKind_I8 || k == TypeKind_I16 || k == TypeKind_U8
}

func Is_i32_compatible(k string) bool {
	k = Real_kind_of(k)
	return (k == TypeKind_I8 ||
		k == TypeKind_I16 ||
		k == TypeKind_I32 ||
		k == TypeKind_U8 ||
		k == TypeKind_U16)
}

func Is_i64_compatible(k string) bool {
	k = Real_kind_of(k)
	return (k == TypeKind_I8 ||
		k == TypeKind_I16 ||
		k == TypeKind_I32 ||
		k == TypeKind_I64 ||
		k == TypeKind_U8 ||
		k == TypeKind_U16 ||
		k == TypeKind_U32)
}

func Is_u8_compatible(k string) bool {
	k = Real_kind_of(k)
	return k == TypeKind_U8
}

func Is_u16_compatible(k string) bool {
	k = Real_kind_of(k)
	return k == TypeKind_U8 || k == TypeKind_U16
}

func Is_u32_compatible(k string) bool {
	k = Real_kind_of(k)
	return k == TypeKind_U8 || k == TypeKind_U16 || k == TypeKind_U32
}

func Is_u64_compatible(k string) bool {
	k = Real_kind_of(k)
	return (k == TypeKind_U8 ||
		k == TypeKind_U16 ||
		k == TypeKind_U32 ||
		k == TypeKind_U64)
}

func Is_f32_compatible(k string) bool {
	k = Real_kind_of(k)
	return (k == TypeKind_F32 ||
		k == TypeKind_I8 ||
		k == TypeKind_I16 ||
		k == TypeKind_I32 ||
		k == TypeKind_I64 ||
		k == TypeKind_U8 ||
		k == TypeKind_U16 ||
		k == TypeKind_U32 ||
		k == TypeKind_U64)
}

func Is_f64_compatible(k string) bool {
	k = Real_kind_of(k)
	return (k == TypeKind_F64 ||
		k == TypeKind_F32 ||
		k == TypeKind_I8 ||
		k == TypeKind_I16 ||
		k == TypeKind_I32 ||
		k == TypeKind_I64 ||
		k == TypeKind_U8 ||
		k == TypeKind_U16 ||
		k == TypeKind_U32 ||
		k == TypeKind_U64)
}

func Types_are_compatible(k1 string, k2 string) bool {
	k1 = Real_kind_of(k1)
	switch k1 {
	case TypeKind_ANY:
		return true
	case TypeKind_I8:
		return Is_i8_compatible(k2)
	case TypeKind_I16:
		return Is_i16_compatible(k2)
	case TypeKind_I32:
		return Is_i32_compatible(k2)
	case TypeKind_I64:
		return Is_i64_compatible(k2)
	case TypeKind_U8:
		return Is_u8_compatible(k2)
	case TypeKind_U16:
		return Is_u16_compatible(k2)
	case TypeKind_U32:
		return Is_u32_compatible(k2)
	case TypeKind_U64:
		return Is_u64_compatible(k2)
	case TypeKind_F32:
		return Is_f32_compatible(k2)
	case TypeKind_F64:
		return Is_f64_compatible(k2)
	case TypeKind_BOOL:
		return k2 == TypeKind_BOOL
	case TypeKind_STR:
		return k2 == TypeKind_STR
	default:
		return false
	}
}

func Is_i16_greater(k string) bool {
	k = Real_kind_of(k)
	return k == TypeKind_U8
}

func Is_i32_greater(k string) bool {
	k = Real_kind_of(k)
	return k == TypeKind_I8 || k == TypeKind_I16
}

func Is_i64_greater(k string) bool {
	k = Real_kind_of(k)
	return k == TypeKind_I8 || k == TypeKind_I16 || k == TypeKind_I32
}

func Is_u8_greater(k string) bool {
	k = Real_kind_of(k)
	return k == TypeKind_I8
}

func Is_u16_greater(k string) bool {
	k = Real_kind_of(k)
	return k == TypeKind_U8 || k == TypeKind_I8 || k == TypeKind_I16
}

func Is_u32_greater(k string) bool {
	k = Real_kind_of(k)
	return (k == TypeKind_U8 ||
		k == TypeKind_U16 ||
		k == TypeKind_I8 ||
		k == TypeKind_I16 ||
		k == TypeKind_I32)
}

func Is_u64_greater(k string) bool {
	k = Real_kind_of(k)
	return (k == TypeKind_U8 ||
		k == TypeKind_U16 ||
		k == TypeKind_U32 ||
		k == TypeKind_I8 ||
		k == TypeKind_I16 ||
		k == TypeKind_I32 ||
		k == TypeKind_I64)
}

func Is_f32_greater(k string) bool {
	return k != TypeKind_F64
}

func Is_f64_greater(k string) bool {
	return true
}

func Is_greater(k1 string, k2 string) bool {
	k1 = Real_kind_of(k1)

	switch k1 {
	case TypeKind_I16:
		return Is_i16_greater(k2)
	case TypeKind_I32:
		return Is_i32_greater(k2)
	case TypeKind_I64:
		return Is_i64_greater(k2)
	case TypeKind_U16:
		return Is_u16_greater(k2)
	case TypeKind_U8:
		return Is_u8_greater(k2)
	case TypeKind_U32:
		return Is_u32_greater(k2)
	case TypeKind_U64:
		return Is_u64_greater(k2)
	case TypeKind_F32:
		return Is_f32_greater(k2)
	case TypeKind_F64:
		return Is_f64_greater(k2)
	case TypeKind_ANY:
		return true
	default:
		return false
	}
}
