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

#ifndef __JANE_PROCESS_HPP
#define __JANE_PROCESS_HPP

#include <cstring>
#include <sys/types.h>

#include "platform.hpp"
#include "slice.hpp"
#include "str.hpp"
#if defined(OS_DARWIN)
#include <climits>
#include <mach-o/dyld.h>
#elif defined(OS_WINDOWS)
#include <windows.h>
#elif defined(OS_LINUX)
#include <linux/limits.h>
#include <unistd.h>
#endif

namespace jane {
jane::Slice<jane::Str> command_line_args;
void setup_command_line_args(int argc, char *argv[]) noexcept;
jane::Str executable(void) noexcept;

void setup_command_line_args(int argc, char *argv[]) noexcept {
#ifdef OS_WINDOWS
  const LPWSTR cmdl{GetCommandLineW()};
  LPWSTR *argvw{CommandLineToArgvW(cmdl, &argc)};
#endif

  jane::command_line_args = jane::Slice<jane::Str>::alloc(argc);
  for (jane::Int i{0}; i < argc; ++i) {
#ifdef OS_WINDOWS
    const LPWSTR warg{argvw[i]};
    jane::comand_line_args[i] =
        jane::utf16_to_utf8_str(warg, std::wcslen(warg));
#else
    jane::command_line_args[i] = argv[i];
#endif
  }
#ifdef OS_WINDOWS
  LocalFree(argvw);
  argvw = nullptr;
#endif
}
jane::Str executable(void) noexcept {
#if defined(OS_DARWIN)
  char buff[PATH_MAX];
  uint32_t buff_size{PATH_MAX};
  if (!_NSGetExecutablePath(buff, &buff_size)) {
    return jane::Str(buff);
  }
  return jane::Str();
#elif defined(OS_WINDOWS)
  wchar_t buffer[MAX_PATH];
  const DWORD n{GetModuleFileName(NULL, buffer, MAX_PATH)};
  if (n) {
    return jane::utf16_to_utf8_str(&buffer[0], n);
  }
  return jane::Str();
#elif defined(OS_LINUX)
  char result[PATH_MAX];
  const ssize_t count{readlink("/proc/self/exe", result, PATH_MAX)};
  if (count != -1) {
    return jane::Str(result);
  }
  return jane::Str();
#endif
}
} // namespace jane

#endif //__JANE_PROCESS_HPP