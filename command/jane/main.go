package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/De-Rune/jane/package/io"
	"github.com/De-Rune/jane/package/jane"
	"github.com/De-Rune/jane/parser"
)

func help(cmd string) {
	if cmd != "" {
		println("this module can be only used as single!")
		return
	}
	helpContent := [][]string{
		{"help", "show help message"},
		{"version", "show version"},
		{"init", "intialize new project jane"},
	}
	maxlen := len(helpContent[0][0])
	for _, part := range helpContent {
		length := len(part[0])
		if length > maxlen {
			maxlen = length
		}
	}
	var sb strings.Builder
	const space = 5
	for _, part := range helpContent {
		sb.WriteString(part[0])
		mlchc := (maxlen - len(part[0])) + space
		for mlchc > 0 {
			sb.WriteByte(' ')
			mlchc--
		}
		sb.WriteString(part[1])
		sb.WriteByte('\n')
	}
	println(sb.String()[:sb.Len()-1])
}

func version(cmd string) {
	if cmd != "" {
		println("this module can only be used as single")
	}
	println("The Jane Programming Language\n" + jane.Version)
}

func initProject(cmd string) {
	if cmd != "" {
		println("this module can only be used as single")
		return
	}
	err := os.WriteFile(jane.SettingFile, []byte(`out_name main
cxx_out_dir ./
cxx_out_name jane.cxx`), 0606)
	if err != nil {
		println(err.Error())
		return
	}
	println("initialized project")
}

func processCommand(namespace, cmd string) bool {
	switch namespace {
	case "help":
		help(cmd)
	case "version":
		version(cmd)
	case "init":
		initProject(cmd)
	default:
		return false
	}
	return true
}

func init() {
	jane.ExecutablePath = filepath.Dir(os.Args[0])
	if len(os.Args) < 2 {
		os.Exit(0)
	}
	var sb strings.Builder
	for _, arg := range os.Args[1:] {
		sb.WriteString(" " + arg)
	}
	os.Args[0] = sb.String()[1:]
	arg := os.Args[0]
	index := strings.Index(arg, " ")
	if index == -1 {
		index = len(arg)
	}
	if processCommand(arg[:index], arg[index:]) {
		os.Exit(0)
	}
}

func loadJnSet() {
	info, err := os.Stat(jane.SettingFile)
	if err != nil || info.IsDir() {
		println(`Jane settings file ("` + jane.SettingFile + `") not found`)
		os.Exit(0)
	}
	jane.JaneSettings = jane.NewJnSet()
	bytes, err := os.ReadFile(jane.SettingFile)
	if err != nil {
		println(err.Error())
		os.Exit(0)
	}
	errors := jane.JaneSettings.Parse(bytes)
	if errors != nil {
		println("Jane settings has errors;")
		for _, err := range errors {
			println(err.Error())
		}
		os.Exit(0)
	}
}

func printErrors(errors []string) {
	defer os.Exit(0)
	for _, message := range errors {
		fmt.Println(message)
	}
}

func appendStandards(code *string) {
	*code = `// testing jane
#pragma region JANE_STANDARD_IMPORTS
#include <iostream>
#include <locale.h>
template <typename any>
#pragma endregion JANE_STANDARD_IMPORTS

#pragma region JANE_BUILTIN_FUNCTIONS
inline void print(any v) {
  std::wcout << v;
}

template <typename any>
inline void println(any v) {
  print(v);
  std::wcout << std::endl;
}

#pragma endregion JANE_BUILTIN_FUNCTIONS

` + *code
}

var routines *sync.WaitGroup

func main() {
	f, err := io.GetJn(os.Args[0])
	if err != nil {
		println(err.Error())
		return
	}
	loadJnSet()
	routines = new(sync.WaitGroup)
	info := new(parser.ParseFileInfo)
	info.File = f
	info.Routines = routines
	routines.Add(1)
	go parser.ParseFile(info)
	routines.Wait()
	if info.Errors != nil {
		printErrors(info.Errors)
	}
	appendStandards(&info.JN_CXX)
	os.WriteFile(
		filepath.Join(
			jane.JaneSettings.Fields["cxx_out_dir"],
			jane.JaneSettings.Fields["cxx_out_name"]), []byte(info.JN_CXX), 0606)
}
