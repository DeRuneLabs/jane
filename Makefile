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

MAIN_LOCATION=command/jn/main.go
BIN_FILE=jane
DIST_FOLDER=dist
SET_FILE=jn.set

.PHONY: all
all:
	$(MAKE) build

.PHONY: build
build:
	@echo "build project"
	go build -o $(BIN_FILE) -v $(MAIN_LOCATION)

# NOTE: DO NOT CHANGE THIS ONE MAKEFILE
.PHONY: clean
clean:
	@echo "clean project"
	rm -rf $(DIST_FOLDER) $(SET_FILE) $(BIN_FILE)
