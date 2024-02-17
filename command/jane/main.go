package main

import (
	"fmt"
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
		println("this mod can only used as single")
		return
	}
	helpContent := [][]string{
		{"help", "showing the help"},
		{"version", "showing version"},
		{"init", "initialize new project"},
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
		println("this mod can only used as single file")
		return
	}
	println("jane programming language\n" + jane.Version)
}

func initProjecT(cmd string) {
	if cmd != "" {
		println("this mod can only be used as single")
		return
	}
	error := os.WriteFile(jane.SettingsFile, []byte(`cxx_out_dir ./ cxx_out_name jane.cpp`), 0606)
	if error != nil {
		println(error.Error())
		return
	}
	println("initialize project")
}

func proccessCommand(namespace, cmd string) bool {
	switch namespace {
	case "help":
		help(cmd)
	case "version":
		version(cmd)
	case "init":
		initProjecT(cmd)
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
	if proccessCommand(arg[:index], arg[index:]) {
		os.Exit(0)
	}
}

func LoadJaneSet() {
	info, error := os.Stat(jane.SettingsFile)
	if error != nil || info.IsDir() {
		println(`JANE settings file ("` + jane.SettingsFile + `") is not found`)
		os.Exit(0)
	}
	jane.JnSettings = jane.NewJnSet()
	bytes, error := os.ReadFile(jane.SettingsFile)
	if error != nil {
		println(error.Error())
		os.Exit(0)
	}
	errors := jane.JnSettings.Parse(bytes)
	if errors != nil {
		println("Jane settings has error;")
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

var routines *sync.WaitGroup

func main() {
	f, error := io.GetJn(os.Args[0])
	if error != nil {
		println(error.Error())
		return
	}
	LoadJaneSet()
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
	os.WriteFile(filepath.Join(jane.JnSettings.Fields["cxx_out_dir"], jane.JnSettings.Fields["cxx_out_name"]), []byte(info.JN_CXX), 0606)
}
