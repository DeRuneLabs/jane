# Copyright (c) 2024 - DeRuneLabs
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
# SOFTWARE.

MAIN_LOCATION = command/jane/main.go
BIN_LOCATION = ../bin/jane

HOST_SYSTEM = $(shell uname | cut -f 1 -d_)
SYSTEM ?= $(HOST_SYSTEM)

ifeq ($(SYSTEM),MSYS)
	SYSTEM = MINGW32
endif
ifeq ($(SYSTEM), MINGW4)
	SYSTEM = MINGW64
endif

HAS_GCC = $(shell which gcc > /dev/null 2> /dev/null && echo true || echo false)
HAS_GO = $(shell which go > /dev/null 2> /dev/null && echo true || echo false)

ifeq ($(HAS_GCC), true)
	DEFAULT_CXX = g++
else
	@echo "please install gnu c and c++"
endif

ifeq ($(HAS_GO), true)
	DEFAULT_GO = go
else
	@echo "please install go compiler"
endif


.PHONY: all
all:
	$(MAKE) build

.PHONY: build
build:
	@echo "build project"
	$(DEFAULT_GO) build -o $(BIN_LOCATION) -v $(MAIN_LOCATION)

# NOTE: DO NOT CHANGE THIS ONE MAKEFILE
.PHONY: clean
clean:
	@echo "clean project"
	rm -rf $(DIST_FOLDER) $(SET_FILE) $(BIN_LOCATION)
