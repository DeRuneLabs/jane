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

#ifndef __JANE_SIGNAL_HPP
#define __JANE_SIGNAL_HPP

#include <csignal>

#define __JANE_SIG constexpr int

#if defined(_WINDOWS)
__JANE_SIG __JANE_SIGHUP{0x1};
__JANE_SIG __JANE_SIGINT{0x2};
__JANE_SIG __JANE_SIGQUIT{0x3};
__JANE_SIG __JANE_SIGILL{0x4};
__JANE_SIG __JANE_SIGTRAP{0x5};
__JANE_SIG __JANE_SIGABRT{0x6};
__JANE_SIG __JANE_SIGBUS{0x7};
__JANE_SIG __JANE_SIGFPE{0x8};
__JANE_SIG __JANE_SIGKILL{0x9};
__JANE_SIG __JANE_SIGSEGV{0xb};
__JANE_SIG __JANE_SIGPIPE{0xd};
__JANE_SIG __JANE_SIGALRM{0xe};
__JANE_SIG __JANE_SIGTERM{0xf};

#elif defined(_DARWIN)

__JANE_SIG __JANE_SIGABRT{0x6};
__JANE_SIG __JANE_SIGALRM{0xe};
__JANE_SIG __JANE_SIGBUS{0xa};
__JANE_SIG __JANE_SIGCHLD{0x14};
__JANE_SIG __JANE_SIGCONT{0x13};
__JANE_SIG __JANE_SIGEMT{0x7};
__JANE_SIG __JANE_SIGFPE{0x8};
__JANE_SIG __JANE_SIGHUP{0x1};
__JANE_SIG __JANE_SIGILL{0x4};
__JANE_SIG __JANE_SIGINFO{0x1d};
__JANE_SIG __JANE_SIGINT{0x2};
__JANE_SIG __JANE_SIGIO{0x17};
__JANE_SIG __JANE_SIGIOT{0x6};
__JANE_SIG __JANE_SIGKILL{0x9};
__JANE_SIG __JANE_SIGPIPE{0xd};
__JANE_SIG __JANE_SIGPROF{0x1b};
__JANE_SIG __JANE_SIGQUIT{0x3};
__JANE_SIG __JANE_SIGSEGV{0xb};
__JANE_SIG __JANE_SIGSTOP{0x11};
__JANE_SIG __JANE_SIGSYS{0xc};
__JANE_SIG __JANE_SIGTERM{0xf};
__JANE_SIG __JANE_SIGTRAP{0x5};
__JANE_SIG __JANE_SIGTSTP{0x12};
__JANE_SIG __JANE_SIGTTIN{0x15};
__JANE_SIG __JANE_SIGTTOU{0x16};
__JANE_SIG __JANE_SIGURG{0x10};
__JANE_SIG __JANE_SIGUSR1{0x1e};
__JANE_SIG __JANE_SIGUSR2{0x1f};
__JANE_SIG __JANE_SIGVTALRM{0x1a};
__JANE_SIG __JANE_SIGWINCH{0x1c};
__JANE_SIG __JANE_SIGXCPU{0x18};
__JANE_SIG __JANE_SIGXFSZ{0x19};

#elif defined(_LINUX)

__JANE_SIG ___JANE_SIGABRT{0x6};
__JANE_SIG ___JANE_SIGALRM{0xe};
__JANE_SIG ___JANE_SIGBUS{0x7};
__JANE_SIG ___JANE_SIGCHLD{0x11};
__JANE_SIG ___JANE_SIGCLD{0x11};
__JANE_SIG ___JANE_SIGCONT{0x12};
__JANE_SIG ___JANE_SIGFPE{0x8};
__JANE_SIG ___JANE_SIGHUP{0x1};
__JANE_SIG ___JANE_SIGILL{0x4};
__JANE_SIG ___JANE_SIGINT{0x2};
__JANE_SIG ___JANE_SIGIO{0x1d};
__JANE_SIG ___JANE_SIGIOT{0x6};
__JANE_SIG ___JANE_SIGKILL{0x9};
__JANE_SIG ___JANE_SIGPIPE{0xd};
__JANE_SIG ___JANE_SIGPOLL{0x1d};
__JANE_SIG ___JANE_SIGPROF{0x1b};
__JANE_SIG ___JANE_SIGPWR{0x1e};
__JANE_SIG ___JANE_SIGQUIT{0x3};
__JANE_SIG ___JANE_SIGSEGV{0xb};
__JANE_SIG ___JANE_SIGSTKFLT{0x10};
__JANE_SIG ___JANE_SIGSTOP{0x13};
__JANE_SIG ___JANE_SIGSYS{0x1f};
__JANE_SIG ___JANE_SIGTERM{0xf};
__JANE_SIG ___JANE_SIGTRAP{0x5};
__JANE_SIG ___JANE_SIGTSTP{0x14};
__JANE_SIG ___JANE_SIGTTIN{0x15};
__JANE_SIG ___JANE_SIGTTOU{0x16};
__JANE_SIG ___JANE_SIGUNUSED{0x1f};
__JANE_SIG ___JANE_SIGURG{0x17};
__JANE_SIG ___JANE_SIGUSR1{0xa};
__JANE_SIG ___JANE_SIGUSR2{0xc};
__JANE_SIG ___JANE_SIGVTALRM{0x1a};
__JANE_SIG ___JANE_SIGWINCH{0x1c};
__JANE_SIG ___JANE_SIGXCPU{0x18};
__JANE_SIG ___JANE_SIGXFSZ{0x19};

#endif // #define(_WINDOWS)

void __jane_set_sig_handler(void (*_Handler)(int _Sig)) noexcept;
void __jane_set_sig_handler(void (*_Handler)(int _Sig)) noexcept {
#if defined(_WINDOWS)
  signal(__JANE_SIGHUP, _Handler);
  signal(__JANE_SIGINT, _Handler);
  signal(__JANE_SIGQUIT, _Handler);
  signal(__JANE_SIGILL, _Handler);
  signal(__JANE_SIGTRAP, _Handler);
  signal(__JANE_SIGABRT, _Handler);
  signal(__JANE_SIGBUS, _Handler);
  signal(__JANE_SIGFPE, _Handler);
  signal(__JANE_SIGKILL, _Handler);
  signal(__JANE_SIGSEGV, _Handler);
  signal(__JANE_SIGPIPE, _Handler);
  signal(__JANE_SIGALRM, _Handler);
  signal(__JANE_SIGTERM, _Handler);
#elif defined(_DARWIN)
  signal(__JANE_SIGABRT, _Handler);
  signal(__JANE_SIGALRM, _Handler);
  signal(__JANE_SIGBUS, _Handler);
  signal(__JANE_SIGCHLD, _Handler);
  signal(__JANE_SIGCONT, _Handler);
  signal(__JANE_SIGEMT, _Handler);
  signal(__JANE_SIGFPE, _Handler);
  signal(__JANE_SIGHUP, _Handler);
  signal(__JANE_SIGILL, _Handler);
  signal(__JANE_SIGINFO, _Handler);
  signal(__JANE_SIGINT, _Handler);
  signal(__JANE_SIGIO, _Handler);
  signal(__JANE_SIGIOT, _Handler);
  signal(__JANE_SIGKILL, _Handler);
  signal(__JANE_SIGPIPE, _Handler);
  signal(__JANE_SIGPROF, _Handler);
  signal(__JANE_SIGQUIT, _Handler);
  signal(__JANE_SIGSEGV, _Handler);
  signal(__JANE_SIGSTOP, _Handler);
  signal(__JANE_SIGSYS, _Handler);
  signal(__JANE_SIGTERM, _Handler);
  signal(__JANE_SIGTRAP, _Handler);
  signal(__JANE_SIGTSTP, _Handler);
  signal(__JANE_SIGTTIN, _Handler);
  signal(__JANE_SIGTTOU, _Handler);
  signal(__JANE_SIGURG, _Handler);
  signal(__JANE_SIGUSR1, _Handler);
  signal(__JANE_SIGUSR2, _Handler);
  signal(__JANE_SIGVTALRM, _Handler);
  signal(__JANE_SIGWINCH, _Handler);
  signal(__JANE_SIGXCPU, _Handler);
  signal(__JANE_SIGXFSZ, _Handler);
#elif defined(_LINUX)
  signal(___JANE_SIGABRT, _Handler);
  signal(___JANE_SIGALRM, _Handler);
  signal(___JANE_SIGBUS, _Handler);
  signal(___JANE_SIGCHLD, _Handler);
  signal(___JANE_SIGCLD, _Handler);
  signal(___JANE_SIGCONT, _Handler);
  signal(___JANE_SIGFPE, _Handler);
  signal(___JANE_SIGHUP, _Handler);
  signal(___JANE_SIGILL, _Handler);
  signal(___JANE_SIGINT, _Handler);
  signal(___JANE_SIGIO, _Handler);
  signal(___JANE_SIGIOT, _Handler);
  signal(___JANE_SIGKILL, _Handler);
  signal(___JANE_SIGPIPE, _Handler);
  signal(___JANE_SIGPOLL, _Handler);
  signal(___JANE_SIGPROF, _Handler);
  signal(___JANE_SIGPWR, _Handler);
  signal(___JANE_SIGQUIT, _Handler);
  signal(___JANE_SIGSEGV, _Handler);
  signal(___JANE_SIGSTKFLT, _Handler);
  signal(___JANE_SIGSTOP, _Handler);
  signal(___JANE_SIGSYS, _Handler);
  signal(___JANE_SIGTERM, _Handler);
  signal(___JANE_SIGTRAP, _Handler);
  signal(___JANE_SIGTSTP, _Handler);
  signal(___JANE_SIGTTIN, _Handler);
  signal(___JANE_SIGTTOU, _Handler);
  signal(___JANE_SIGUNUSED, _Handler);
  signal(___JANE_SIGURG, _Handler);
  signal(___JANE_SIGUSR1, _Handler);
  signal(___JANE_SIGUSR2, _Handler);
  signal(___JANE_SIGVTALRM, _Handler);
  signal(___JANE_SIGWINCH, _Handler);
  signal(___JANE_SIGXCPU, _Handler);
  signal(___JANE_SIGXFSZ, _Handler);
#endif // defined(_WINDOWS)
}

#endif // !__JANE_SIGNAL_HPP
