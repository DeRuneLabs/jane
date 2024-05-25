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

#ifndef __JANE_PANIC_HPP
#define __JANE_PANIC_HPP

#include <exception>
#include <sstream>
#include <string>

namespace jane {
class Exception;

template <typename T> void panic(const T &expr);

class Exception : public std::exception {
private:
  std::string message;

public:
  Exception(void) noexcept {}

  Exception(char *message) noexcept { this->message = message; }

  Exception(const std::string &message) noexcept { this->message = message; }

  char *what(void) noexcept { return (char *)this->message.c_str(); }

  const char *what(void) const noexcept { return this->message.c_str(); }
};

template <typename T> void panic(const T &expr) {
  std::stringstream sstream;
  sstream << expr;
  jane::Exception exception(sstream.str());
  throw exception;
}
} // namespace jane

#endif // __JANE_PANIC_HPP