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

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/DeRuneLabs/jane/documenter"
	"github.com/DeRuneLabs/jane/package/jn"
	"github.com/DeRuneLabs/jane/package/jnapi"
	"github.com/DeRuneLabs/jane/package/jnio"
	"github.com/DeRuneLabs/jane/package/jnset"
	"github.com/DeRuneLabs/jane/transpiler"
)

const (
	command_help    = "help"
	command_version = "version"
	command_init    = "init"
	command_doc     = "doc"
	command_bug     = "bug"
	command_tool    = "tool"
)

var helpmap = [...][2]string{
	0: {command_help, "Show help"},
	1: {command_version, "Show version"},
	2: {command_init, "Initialize new project here"},
	3: {command_doc, "Documentize Jane source code"},
	4: {command_bug, "Start a new bug report"},
	5: {command_tool, "Tools for effective Jane"},
}

func help(cmd string) {
	if cmd != "" {
		println("can only be used as single")
		return
	}
	max := len(helpmap[0][0])
	for _, key := range helpmap {
		len := len(key[0])
		if len > max {
			max = len
		}
	}
	var sb strings.Builder
	const space = 5
	for _, part := range helpmap {
		sb.WriteString(part[0])
		sb.WriteString(strings.Repeat(" ", (max-len(part[0]))+space))
		sb.WriteString(part[1])
		sb.WriteByte('\n')
	}
	println(sb.String()[:sb.Len()-1])
}

func version(cmd string) {
	if cmd != "" {
		println("can only be used as single")
		return
	}
	println("jn version: ", jn.Version)
}

func init_project(cmd string) {
	if cmd != "" {
		println("can only be used as single")
		return
	}
	bytes, err := json.MarshalIndent(*jnset.Default, "", "\t")
	if err != nil {
		println(err)
		os.Exit(0)
	}
	err = os.WriteFile(jn.SettingsFile, bytes, 0666)
	if err != nil {
		println(err.Error())
		os.Exit(0)
	}
	println("initialize jane project")
}

func doc(cmd string) {
	fmt_json := false
	cmd = strings.TrimSpace(cmd)
	if strings.HasPrefix(cmd, "--json") {
		cmd = strings.TrimSpace(cmd[len("--json"):])
		fmt_json = true
	}
	paths := strings.SplitN(cmd, " ", -1)
	for _, path := range paths {
		path = strings.TrimSpace(path)
		t := compile(path, false, true, true)
		if t == nil {
			continue
		}
		if print_logs(t) {
			fmt.Println(jn.GetError("doc_couldnt_generated", path))
			continue
		}
		docjson, err := documenter.Doc(t, fmt_json)
		if err != nil {
			fmt.Println(jn.GetError("error", err.Error()))
			continue
		}
		path = path[:len(path)-len(jn.SrcExt)]
		path = filepath.Join(jn.Set.CppOutDir, path+jn.DocExt)
		write_output(path, docjson)
	}
}

func open_url(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default:
		cmd = "xdg-open"
	}
	args = append(args, url)
	command := exec.Command(cmd, args...)
	return command.Start()
}

func bug(cmd string) {
	if cmd != "" {
		println("can only be used as single")
		return
	}
	err := open_url("https://github.com/DeRuneLabs/jane/issues/new")
	if err != nil {
		fmt.Println(err.Error())
	}
}

func list_horizontal_slice(s []string) string {
	lst := fmt.Sprint(s)
	return lst[1 : len(lst)-1]
}

func tool(cmd string) {
	if cmd == "" {
		println(`tooling command:
distos display all supported operating system
distarch display all support architecture`)
		return
	}
	switch cmd {
	case "distos":
		print("supported operating system list:\n")
		println(list_horizontal_slice(jn.Distos))
	case "distarch":
		print("supported architecture:\n")
		println(list_horizontal_slice(jn.Distarch))
	default:
		println("undefine command: " + cmd)
	}
}

func process_command(namespace, cmd string) bool {
	cmd = strings.TrimSpace(cmd)
	switch namespace {
	case command_help:
		help(cmd)
	case command_version:
		version(cmd)
	case command_init:
		init_project(cmd)
	case command_doc:
		doc(cmd)
	case command_bug:
		bug(cmd)
	case command_tool:
		tool(cmd)
	default:
		return false
	}
	return true
}

func init() {
	execp, err := os.Executable()
	if err != nil {
		println(err.Error())
		os.Exit(0)
	}
	jn.WorkingPath, err = os.Getwd()
	if err != nil {
		println(err.Error())
		os.Exit(0)
	}
	execp = filepath.Dir(execp)
	jn.ExecPath = execp
	jn.StdlibPath = filepath.Join(jn.ExecPath, jn.Stdlib)
	jnapi.JNCHeader = filepath.Join(jn.ExecPath, "api")
	jnapi.JNCHeader = filepath.Join(jnapi.JNCHeader, "jnc.hpp")
	jn.LangsPath = filepath.Join(jn.ExecPath, jn.Localizations)
	if len(os.Args) < 2 {
		os.Exit(0)
	}
	var sb strings.Builder
	for _, arg := range os.Args[1:] {
		sb.WriteString(" " + arg)
	}
	os.Args[0] = sb.String()[1:]
	arg := os.Args[0]
	i := strings.Index(arg, " ")
	if i == -1 {
		i = len(arg)
	}
	if process_command(arg[:i], arg[i:]) {
		os.Exit(0)
	}
}

func load_localization() {
	lang := strings.TrimSpace(jn.Set.Language)
	if lang == "" || lang == "default" {
		return
	}
	path := filepath.Join(jn.LangsPath, lang+".json")
	bytes, err := os.ReadFile(path)
	if err != nil {
		println("language not loaded (using default);")
		println(err.Error())
	}
	err = json.Unmarshal(bytes, &jn.Errors)
	if err != nil {
		println("language error couldn't load (using default);")
		println(err.Error())
		return
	}
}

func check_mode() {
	mode := jn.Set.Mode
	if mode != jnset.ModeTranspile && mode != jnset.ModeCompile {
		println(jn.GetError("invalid_value_for_key", mode, "mode"))
		os.Exit(0)
	}
}

func check_compiler() {
	c := jn.Set.Compiler
	if c != jn.CompilerGCC && c != jn.CompilerClang {
		println(jn.GetError("invalid_value_for_key", c, "compiler"))
		os.Exit(0)
	}
}

func load_jnset() {
	info, err := os.Stat(jn.SettingsFile)
	if err != nil || info.IsDir() {
		jn.Set = new(jnset.Set)
		*jn.Set = *jnset.Default
		return
	}
	bytes, err := os.ReadFile(jn.SettingsFile)
	if err != nil {
		println(err.Error())
		os.Exit(0)
	}
	jn.Set, err = jnset.Load(bytes)
	if err != nil {
		println("jane settings has errors;")
		println(err.Error())
		os.Exit(0)
	}
	load_localization()
	check_mode()
	check_compiler()
}

func print_logs(t *transpiler.Transpiler) bool {
	var str strings.Builder
	for _, log := range t.Warnings {
		str.WriteString(log.String())
		str.WriteByte('\n')
	}
	for _, log := range t.Errors {
		str.WriteString(log.String())
		str.WriteByte('\n')
	}
	print(str.String())
	return len(t.Errors) > 0
}

func append_standard(code *string) {
	year, month, day := time.Now().Date()
	hour, min, _ := time.Now().Clock()
	timeStr := fmt.Sprintf("%d/%d/%d %d.%d (DD/MM/YYYY) (HH.MM)",
		day, month, year, hour, min)
	var sb strings.Builder
	sb.WriteString("// Auto generated by JN compiler.\n")
	sb.WriteString("// JN version:")
	sb.WriteString(jn.Version)
	sb.WriteByte('\n')
	sb.WriteString("// Date: ")
	sb.WriteString(timeStr)
	sb.WriteString("\n\n")
	sb.WriteString("// this file contains cpp module code which is automatically generated by JN")
	sb.WriteByte('\n')
	sb.WriteString("// compiler. generated code in this file provide cpp functions and structures")
	sb.WriteByte('\n')
	sb.WriteString("// corresponding to the definition in the JN source files")
	sb.WriteString("\n\n")
	sb.WriteString("\n\n#include \"")
	sb.WriteString(jnapi.JNCHeader)
	sb.WriteString("\"\n\n")
	sb.WriteString(*code)
	*code = sb.String()
}

func write_output(path, content string) {
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, 0o777)
	if err != nil {
		println(err.Error())
		os.Exit(0)
	}
	f, err := os.Create(path)
	if err != nil {
		println(err.Error())
		os.Exit(0)
	}
	_, err = f.WriteString(content)
	if err != nil {
		println(err.Error())
		os.Exit(0)
	}
}

func compile(path string, main, nolocal, justDefs bool) *transpiler.Transpiler {
	load_jnset()
	t := transpiler.New(nil)
	inf, err := os.Stat(jn.StdlibPath)
	if err != nil || !inf.IsDir() {
		t.PushErr("stdlib_not_exist")
		return t
	}
	f, err := jnio.OpenJn(path)
	if err != nil {
		println(err.Error())
		return nil
	}
	if !jnio.IsPassFileAnnotation(path) {
		t.PushErr("file_not_useable")
		return t
	}
	t.File = f
	t.NoLocalPkg = nolocal
	t.Parsef(main, justDefs)
	return t
}

func exec_post_commands() {
	for _, cmd := range jn.Set.PostCommands {
		fmt.Println(">", cmd)
		parts := strings.SplitN(cmd, " ", -1)
		err := exec.Command(parts[0], parts[1:]...).Run()
		if err != nil {
			println(err.Error())
		}
	}
}

func generate_compile_command(source_path string) (c, cmd string) {
	var cpp strings.Builder
	cpp.WriteString("-g -O0 ")
	cpp.WriteString(source_path)
	return jn.Set.CompilerPath, cpp.String()
}

func do_spell(cpp string) {
	defer exec_post_commands()
	path := filepath.Join(jn.WorkingPath, jn.Set.CppOutDir)
	path = filepath.Join(path, jn.Set.CppOutName)
	write_output(path, cpp)
	switch jn.Set.Mode {
	case jnset.ModeCompile:
		c, cmd := generate_compile_command(path)
		println(c + " " + cmd)
		command := exec.Command(c, strings.SplitN(cmd, " ", -1)...)
		err := command.Start()
		if err != nil {
			println(err.Error())
		}
		err = command.Wait()
		if err != nil {
			println(err.Error())
		}
	}
}

func main() {
	fpath := os.Args[0]
	t := compile(fpath, true, false, false)
	if t == nil {
		return
	}
	if print_logs(t) {
		os.Exit(0)
	}
	cpp := t.Cpp()
	append_standard(&cpp)
	do_spell(cpp)
}
