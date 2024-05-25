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

#ifndef __JANE_MAP_HPP
#define __JANE_MAP_HPP

#include "slice.hpp"
#include "str.hpp"
#include "types.hpp"
#include <cstddef>
#include <unordered_map>
namespace jane {
class MapKeyHasher;
template <typename Key, typename Value> class Map;

class MapKeyHasher {
public:
  size_t operator()(const jane::Str &key) const noexcept {
    size_t hash{0};
    for (jane::Int i{0}; i < key.len(); ++i) {
      hash += key[i] % 7;
    }
    return hash;
  }
  template <typename T> inline size_t operator()(const T &obj) const noexcept {
    return this->operator()(jane::to_str<T>(obj));
  }
};

template <typename Key, typename Value> class Map {
public:
  mutable std::unordered_map<Key, Value, MapKeyHasher> buffer{};
  Map<Key, Value>(void) noexcept {}
  Map<Key, Value>(const std::nullptr_t) noexcept {}
  Map<Key, Value>(
      const std::initializer_list<std::pair<Key, Value>> &src) noexcept {
    for (const std::pair<Key, Value> &pair : src) {
      this->buffer.insert(pair);
    }
  }

  inline constexpr auto begin(void) noexcept { return this->buffer.begin(); }

  inline constexpr auto end(void) noexcept { return this->buffer.end(); }

  inline constexpr auto end(void) const noexcept { return this->buffer.end(); }

  inline void clear(void) noexcept { this->buffer.clear(); }

  jane::Slice<Key> keys(void) const noexcept {
    jane::Slice<Key> keys(this->len());
    jane::Uint index{0};
    for (const auto &pair : *this) {
      keys._slice[index++] = pair.second;
    }
    return keys;
  }

  inline constexpr jane::Bool has(const Key &key) const noexcept {
    return this->buffer.find(key) != this->end();
  }

  inline jane::Int len(void) const noexcept { return this->buffer.size(); }

  inline void del(const Key &key) noexcept { this->buffer.erase(key); }

  inline jane::Bool operator==(const std::nullptr_t) const noexcept {
    return this->buffer.empty();
  }

  inline jane::Bool operator!=(const std::nullptr_t) const noexcept {
    return !this->operator==(nullptr);
  }

  Value &operator[](const Key &key) { return this->buffer[key]; }

  Value &operator[](const Key &key) const { return this->buffer[key]; }

  friend std::ostream &operator<<(std::ostream &stream,
                                  const Map<Key, Value> &src) noexcept {
    stream << '{';
    jane::Int length{src.len()};
    for (const auto pair : src) {
      stream << pair.first;
      stream << ':';
      stream << pair.second;
      if (--length > 0) {
        stream << ", ";
      }
    }
    stream << '}';
    return stream;
  }
};
} // namespace jane

#endif // __JANE_MAP_HPP