// Copyright (c) 2024 - DeRuneLabs
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
// SOFTWARE.

package jnset

import (
	"encoding/json"
	"runtime"
)

const (
	ModeTranspile = "transpile"
	ModeCompile   = "compile"
)

// set struct
type Set struct {
	CppOutDir    string   `json:"cpp_out_dir"`
	CppOutName   string   `json:"cpp_out_name"`
	OutName      string   `json:"out_name"`
	Language     string   `json:"language"`
	Mode         string   `json:"mode"`
	PostCommands []string `json:"post_commands"`
	Indent       string   `json:"indent"`
	IndentCount  int      `json:"indent_count"`
	Compiler     string   `json:"compiler"`
	CompilerPath string   `json:"compiler_path"`
}

// default instance of jn.set
var Default = &Set{
	CppOutDir:    "./dist",
	CppOutName:   "jn.cpp",
	OutName:      "main",
	Language:     "",
	Mode:         "transpile",
	Indent:       "\t",
	IndentCount:  1,
	Compiler:     "",
	CompilerPath: "",
	PostCommands: []string{},
}

// load set from json string
func Load(bytes []byte) (*Set, error) {
	set := *Default
	err := json.Unmarshal(bytes, &set)
	if err != nil {
		return nil, err
	}
	return &set, nil
}

// initialize jn.set
func init() {
	if runtime.GOOS == "windows" {
		Default.Compiler = "gcc"
		Default.CompilerPath = "g++"
	} else {
		Default.Compiler = "clang"
		Default.CompilerPath = "clang++"
	}
}
