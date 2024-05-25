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

#ifndef __JANE_TERMINATE_HPP
#define __JANE_TERMINATE_HPP

#include "builtin.hpp"
#include "error.hpp"
#include "panic.hpp"
#include "str.hpp"
#include "trait.hpp"
#include "types.hpp"
#include <cstdlib>
#include <exception>

namespace jane {
struct Error {
  virtual jane::Str error(void) { return jane::Str(); }
  virtual ~Error(void) noexcept {}

  jane::Bool operator==(const Error &) { return false; }
  jane::Bool operator!=(const Error &src) { return !this->operator==(src); }

  friend std::ostream &operator<<(std::ostream &stream, Error error) noexcept {
    return stream << error.error();
  }
};

void terminate_handler(void) noexcept;

jane::Trait<Error> exception_to_error(const jane::Exception &exception);

void terminate_handler(void) noexcept {
  try {
    std::rethrow_exception(std::current_exception());
  } catch (const jane::Exception &e) {
    jane::println(std::string("panic: ") + std::string(e.what()));
    std::exit(jane::EXIT_PANIC);
  }
}

jane::Trait<Error> exception_to_error(const jane::Exception &exception) {
  struct PanicError : public Error {
    jane::Str message;
    jane::Str error(void) { return this->message; }
  };
  struct PanicError error;
  error.message = jane::to_str(exception.what());
  return jane::Trait<Error>(error);
}
} // namespace jane

#endif // __JANE_TERMINATE_HPP