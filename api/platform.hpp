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

#ifndef __JANE_PLATFORM_HPP
#define __JANE_PLATFORM_HPP

#if defined(WIN32) || defined(_WIN32) || defined(__WIN32__) || defined(__NT__)
#define OS_WINDOWS
#elif defined(__linux__) || defined(linux) || defined(__linux)
#define OS_LINUX
#elif defined(__APPLE__) || defined(__MACH__)
#define OS_DARWIN
#endif

#if defined(OS_LINUX) || defined(OS_DARWIN)
#define OS_UNIX
#endif

#if defined(__amd64) || defined(__amd64__) || defined(__x86_64) ||             \
    defined(__x86_64__) || defined(_M_AMD64)
#define ARCH_AMD64
#elif defined(__arm__) || defined(__thumb__) || defined(_M_ARM) ||             \
    defined(__arm)
#define ARCH_ARM
#elif defined(__aarch64__)
#define ARCH_ARM64
#elif defined(i386) || defined(__i386) || defined(__i386__) ||                 \
    defined(_X86_) || defined(__I86__) || defined(__386)
#define ARCH_I386
#endif

#if defined(ARCH_AMD64) || defined(ARCH_ARM64)
#define ARCH_64BIT
#else
#define ARCH_32BIT
#endif

#endif // __JANE_PLATFORM_HPP