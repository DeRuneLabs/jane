// Copyright (c) 2024 arfy slowy - DeRuneLabs
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
// SOFTWARE

package build

const (
	OS_WINDOWS = "windows"
	OS_LINUX   = "linux"
	OS_DARWIN  = "darwin"
	OS_UNIX    = "unix"
)

const (
	ARCH_ARM   = "arm"
	ARCH_ARM64 = "arm64"
	ARCH_AMD64 = "amd64"
	ARCH_I386  = "i386"
	ARCH_64Bit = "64bit"
	ARCH_32Bit = "32bit"
)

var DISTOS = []string{
	OS_WINDOWS,
	OS_LINUX,
	OS_DARWIN,
}

var DISTARCH = []string{
	ARCH_ARM,
	ARCH_ARM64,
	ARCH_AMD64,
	ARCH_I386,
}

const (
	goos_windows = "windows"
	goos_darwin  = "darwin"
	goos_linux   = "linux"
)

const (
	goarch_i386  = "386"
	goarch_amd64 = "amd64"
	goarch_arm   = "arm"
	goarch_arm64 = "arm64"
)

func IsWindows(os string) bool {
	return os == goos_windows
}

func IsDarwin(os string) bool {
	return os == goos_darwin
}

func IsLinux(os string) bool {
	return os == goos_linux
}

func IsUnix(os string) bool {
	return IsDarwin(os) || IsLinux(os)
}

func IsI386(arch string) bool {
	return arch == goarch_i386
}

func IsAmd64(arch string) bool {
	return arch == goarch_amd64
}

func IsArm(arch string) bool {
	return arch == goarch_arm
}

func IsArm64(arch string) bool {
	return arch == goarch_arm64
}

func IsX32(arch string) bool {
	return IsI386(arch) || IsArm(arch)
}

func IsX64(arch string) bool {
	return IsAmd64(arch) || IsArm64(arch)
}
