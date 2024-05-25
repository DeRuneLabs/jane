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

#ifndef __JANE_SIGNAL_HPP
#define __JANE_SIGNAL_HPP

#include "platform.hpp"
#include <csignal>
#include <cstdlib>

namespace jane {
typedef int Signal;

void set_sig_handler(void (*handler)(int sig)) noexcept;
void signal_handler(int signal) noexcept;

#if defined(OS_WINDOWS)
constexpr jane::Signal SIG_HUP{0x1};
constexpr jane::Signal SIG_INT{0x2};
constexpr jane::Signal SIG_QUIT{0x3};
constexpr jane::Signal SIG_ILL{0x4};
constexpr jane::Signal SIG_TRAP{0x5};
constexpr jane::Signal SIG_ABRT{0x6};
constexpr jane::Signal SIG_BUS{0x7};
constexpr jane::Signal SIG_FPE{0x8};
constexpr jane::Signal SIG_KILL{0x9};
constexpr jane::Signal SIG_SEGV{0xb};
constexpr jane::Signal SIG_PIPE{0xd};
constexpr jane::Signal SIG_ALRM{0xe};
constexpr jane::Signal SIG_TERM{0xf};
#elif defined(OS_DARWIN)
constexpr jane::Signal SIG_ABRT{0x6};
constexpr jane::Signal SIG_ALRM{0xe};
constexpr jane::Signal SIG_BUS{0xa};
constexpr jane::Signal SIG_CHLD{0x14};
constexpr jane::Signal SIG_CONT{0x13};
constexpr jane::Signal SIG_EMT{0x7};
constexpr jane::Signal SIG_FPE{0x8};
constexpr jane::Signal SIG_HUP{0x1};
constexpr jane::Signal SIG_ILL{0x4};
constexpr jane::Signal SIG_INFO{0x1d};
constexpr jane::Signal SIG_INT{0x2};
constexpr jane::Signal SIG_IO{0x17};
constexpr jane::Signal SIG_IOT{0x6};
constexpr jane::Signal SIG_KILL{0x9};
constexpr jane::Signal SIG_PIPE{0xd};
constexpr jane::Signal SIG_PROF{0x1b};
constexpr jane::Signal SIG_QUIT{0x3};
constexpr jane::Signal SIG_SEGV{0xb};
constexpr jane::Signal SIG_STOP{0x11};
constexpr jane::Signal SIG_SYS{0xc};
constexpr jane::Signal SIG_TERM{0xf};
constexpr jane::Signal SIG_TRAP{0x5};
constexpr jane::Signal SIG_TSTP{0x12};
constexpr jane::Signal SIG_TTIN{0x15};
constexpr jane::Signal SIG_TTOU{0x16};
constexpr jane::Signal SIG_URG{0x10};
constexpr jane::Signal SIG_USR1{0x1e};
constexpr jane::Signal SIG_USR2{0x1f};
constexpr jane::Signal SIG_VTALRM{0x1a};
constexpr jane::Signal SIG_WINCH{0x1c};
constexpr jane::Signal SIG_XCPU{0x18};
constexpr jane::Signal SIG_XFSZ{0x19};
#elif defined(OS_LINUX)
constexpr jane::Signal SIG_ABRT{0x6};
constexpr jane::Signal SIG_ALRM{0xe};
constexpr jane::Signal SIG_BUS{0x7};
constexpr jane::Signal SIG_CHLD{0x11};
constexpr jane::Signal SIG_CLD{0x11};
constexpr jane::Signal SIG_CONT{0x12};
constexpr jane::Signal SIG_FPE{0x8};
constexpr jane::Signal SIG_HUP{0x1};
constexpr jane::Signal SIG_ILL{0x4};
constexpr jane::Signal SIG_INT{0x2};
constexpr jane::Signal SIG_IO{0x1d};
constexpr jane::Signal SIG_IOT{0x6};
constexpr jane::Signal SIG_KILL{0x9};
constexpr jane::Signal SIG_PIPE{0xd};
constexpr jane::Signal SIG_POLL{0x1d};
constexpr jane::Signal SIG_PROF{0x1b};
constexpr jane::Signal SIG_PWR{0x1e};
constexpr jane::Signal SIG_QUIT{0x3};
constexpr jane::Signal SIG_SEGV{0xb};
constexpr jane::Signal SIG_STKFLT{0x10};
constexpr jane::Signal SIG_STOP{0x13};
constexpr jane::Signal SIG_SYS{0x1f};
constexpr jane::Signal SIG_TERM{0xf};
constexpr jane::Signal SIG_TRAP{0x5};
constexpr jane::Signal SIG_TSTP{0x14};
constexpr jane::Signal SIG_TTIN{0x15};
constexpr jane::Signal SIG_TTOU{0x16};
constexpr jane::Signal SIG_UNUSED{0x1f};
constexpr jane::Signal SIG_URG{0x17};
constexpr jane::Signal SIG_USR1{0xa};
constexpr jane::Signal SIG_USR2{0xc};
constexpr jane::Signal SIG_VTALRM{0x1a};
constexpr jane::Signal SIG_WINCH{0x1c};
constexpr jane::Signal SIG_XCPU{0x18};
constexpr jane::Signal SIG_XFSZ{0x19};
#endif

void set_sig_handler(void (*handler)(int _sig)) noexcept {
#if defined(OS_WINDOWS)
  std::signal(jane::SIG_HUP, handler);
  std::signal(jane::SIG_INT, handler);
  std::signal(jane::SIG_QUIT, handler);
  std::signal(jane::SIG_ILL, handler);
  std::signal(jane::SIG_TRAP, handler);
  std::signal(jane::SIG_ABRT, handler);
  std::signal(jane::SIG_BUS, handler);
  std::signal(jane::SIG_FPE, handler);
  std::signal(jane::SIG_KILL, handler);
  std::signal(jane::SIG_SEGV, handler);
  std::signal(jane::SIG_PIPE, handler);
  std::signal(jane::SIG_ALRM, handler);
  std::signal(jane::SIG_TERM, handler);
#elif defined(OS_DARWIN)
  std::signal(jane::SIG_ABRT, handler);
  std::signal(jane::SIG_ALRM, handler);
  std::signal(jane::SIG_BUS, handler);
  std::signal(jane::SIG_CHLD, handler);
  std::signal(jane::SIG_CONT, handler);
  std::signal(jane::SIG_EMT, handler);
  std::signal(jane::SIG_FPE, handler);
  std::signal(jane::SIG_HUP, handler);
  std::signal(jane::SIG_ILL, handler);
  std::signal(jane::SIG_INFO, handler);
  std::signal(jane::SIG_INT, handler);
  std::signal(jane::SIG_IO, handler);
  std::signal(jane::SIG_IOT, handler);
  std::signal(jane::SIG_KILL, handler);
  std::signal(jane::SIG_PIPE, handler);
  std::signal(jane::SIG_PROF, handler);
  std::signal(jane::SIG_QUIT, handler);
  std::signal(jane::SIG_SEGV, handler);
  std::signal(jane::SIG_STOP, handler);
  std::signal(jane::SIG_SYS, handler);
  std::signal(jane::SIG_TERM, handler);
  std::signal(jane::SIG_TRAP, handler);
  std::signal(jane::SIG_TSTP, handler);
  std::signal(jane::SIG_TTIN, handler);
  std::signal(jane::SIG_TTOU, handler);
  std::signal(jane::SIG_URG, handler);
  std::signal(jane::SIG_USR1, handler);
  std::signal(jane::SIG_USR2, handler);
  std::signal(jane::SIG_VTALRM, handler);
  std::signal(jane::SIG_WINCH, handler);
  std::signal(jane::SIG_XCPU, handler);
  std::signal(jane::SIG_XFSZ, handler);
#elif defined(OS_LINUX)
  std::signal(jane::SIG_ABRT, handler);
  std::signal(jane::SIG_ALRM, handler);
  std::signal(jane::SIG_BUS, handler);
  std::signal(jane::SIG_CHLD, handler);
  std::signal(jane::SIG_CLD, handler);
  std::signal(jane::SIG_CONT, handler);
  std::signal(jane::SIG_FPE, handler);
  std::signal(jane::SIG_HUP, handler);
  std::signal(jane::SIG_ILL, handler);
  std::signal(jane::SIG_INT, handler);
  std::signal(jane::SIG_IO, handler);
  std::signal(jane::SIG_IOT, handler);
  std::signal(jane::SIG_KILL, handler);
  std::signal(jane::SIG_PIPE, handler);
  std::signal(jane::SIG_POLL, handler);
  std::signal(jane::SIG_PROF, handler);
  std::signal(jane::SIG_PWR, handler);
  std::signal(jane::SIG_QUIT, handler);
  std::signal(jane::SIG_SEGV, handler);
  std::signal(jane::SIG_STKFLT, handler);
  std::signal(jane::SIG_STOP, handler);
  std::signal(jane::SIG_SYS, handler);
  std::signal(jane::SIG_TERM, handler);
  std::signal(jane::SIG_TRAP, handler);
  std::signal(jane::SIG_TSTP, handler);
  std::signal(jane::SIG_TTIN, handler);
  std::signal(jane::SIG_TTOU, handler);
  std::signal(jane::SIG_UNUSED, handler);
  std::signal(jane::SIG_URG, handler);
  std::signal(jane::SIG_USR1, handler);
  std::signal(jane::SIG_USR2, handler);
  std::signal(jane::SIG_VTALRM, handler);
  std::signal(jane::SIG_WINCH, handler);
  std::signal(jane::SIG_XCPU, handler);
  std::signal(jane::SIG_XFSZ, handler);
#endif
}

void signal_handler(int signal) noexcept {
  jane::print("program terminated with signal: ");
  jane::println(signal);
  std::exit(signal);
}

} // namespace jane

#endif // __JANE_SIGNAL_HPP