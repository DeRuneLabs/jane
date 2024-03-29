package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/DeRuneLabs/jane/documenter"
	"github.com/DeRuneLabs/jane/package/jn"
	"github.com/DeRuneLabs/jane/package/jnapi"
	"github.com/DeRuneLabs/jane/package/jnio"
	"github.com/DeRuneLabs/jane/package/jnset"
	"github.com/DeRuneLabs/jane/parser"
)

type Parser = parser.Parser

const (
	commandHelp    = "help"
	commandVersion = "version"
	commandInit    = "init"
	commandDoc     = "doc"
)

const (
	localizationErrors   = "error.json"
	localizationWarnings = "warning.json"
)

var helpmap = [...][2]string{
	0: {commandHelp, "Show help."},
	1: {commandVersion, "Show version."},
	2: {commandInit, "Initialize new project here."},
	3: {commandDoc, "Documentize Jn source code."},
}

func help(cmd string) {
	if cmd != "" {
		println("This mod can only be used as single")
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
		println("This mod can only be used as single!")
		return
	}
	println("jn version", jn.Version)
}

func initProject(cmd string) {
	if cmd != "" {
		println("This module can only be used as single!")
		return
	}
	bytes, err := json.MarshalIndent(*jnset.Default, "", "\t")
	if err != nil {
		println(err)
		os.Exit(0)
	}
	err = ioutil.WriteFile(jn.SettingsFile, bytes, 0666)
	if err != nil {
		println(err.Error())
		os.Exit(0)
	}
	println("Initialized project.")
}

func doc(cmd string) {
	cmd = strings.TrimSpace(cmd)
	paths := strings.SplitN(cmd, " ", -1)
	for _, path := range paths {
		path = strings.TrimSpace(path)
		p := compile(path, false, true, true)
		if p == nil {
			continue
		}
		if printlogs(p) {
			fmt.Println(jn.GetError("doc_couldnt_generated", path))
			continue
		}
		docjson, err := documenter.Doc(p)
		if err != nil {
			fmt.Println(jn.GetError("error", err.Error()))
			continue
		}
		path = path[:len(path)-len(jn.SrcExt)]
		path = filepath.Join(jn.Set.CxxOutDir, path+jn.DocExt)
		writeOutput(path, docjson)
	}
}

func processCommand(namespace, cmd string) bool {
	switch namespace {
	case commandHelp:
		help(cmd)
	case commandVersion:
		version(cmd)
	case commandInit:
		initProject(cmd)
	case commandDoc:
		doc(cmd)
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
	if processCommand(arg[:i], arg[i:]) {
		os.Exit(0)
	}
}

func loadLangWarns(path string, infos []fs.FileInfo) {
	i := -1
	for j, f := range infos {
		if f.IsDir() || f.Name() != localizationWarnings {
			continue
		}
		i = j
		path = filepath.Join(path, f.Name())
		break
	}
	if i == -1 {
		return
	}
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		println("Language's warnings couldn't loaded (uses default);")
		println(err.Error())
		return
	}
	err = json.Unmarshal(bytes, &jn.Warnings)
	if err != nil {
		println("Language's warnings couldn't loaded (uses default);")
		println(err.Error())
		return
	}
}

func loadLangErrs(path string, infos []fs.FileInfo) {
	i := -1
	for j, f := range infos {
		if f.IsDir() || f.Name() != localizationErrors {
			continue
		}
		i = j
		path = filepath.Join(path, f.Name())
		break
	}
	if i == -1 {
		return
	}
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		println("Language's errors couldn't loaded (uses default);")
		println(err.Error())
		return
	}
	err = json.Unmarshal(bytes, &jn.Errors)
	if err != nil {
		println("Language's errors couldn't loaded (uses default);")
		println(err.Error())
		return
	}
}

func loadLang() {
	lang := strings.TrimSpace(jn.Set.Language)
	if lang == "" || lang == "default" {
		return
	}
	path := filepath.Join(jn.LangsPath, lang)
	infos, err := ioutil.ReadDir(path)
	if err != nil {
		println("Language couldn't loaded (uses default);")
		println(err.Error())
		return
	}
	loadLangWarns(path, infos)
	loadLangErrs(path, infos)
}

func checkMode() {
	lower := strings.ToLower(jn.Set.Mode)
	if lower != jnset.ModeTranspile &&
		lower != jnset.ModeCompile {
		key, _ := reflect.TypeOf(jn.Set).Elem().FieldByName("Mode")
		tag := string(key.Tag)
		tag = tag[6 : len(tag)-1]
		println(jn.GetError("invalid_value_for_key", jn.Set.Mode, tag))
		os.Exit(0)
	}
	jn.Set.Mode = lower
}

func loadJnSet() {
	info, err := os.Stat(jn.SettingsFile)
	if err != nil || info.IsDir() {
		println(`JN settings file ("` + jn.SettingsFile + `") is not found!`)
		os.Exit(0)
	}
	bytes, err := os.ReadFile(jn.SettingsFile)
	if err != nil {
		println(err.Error())
		os.Exit(0)
	}
	jn.Set, err = jnset.Load(bytes)
	if err != nil {
		println("X settings has errors;")
		println(err.Error())
		os.Exit(0)
	}
	loadLang()
	checkMode()
}

func printlogs(p *Parser) bool {
	var str strings.Builder
	for _, log := range p.Warnings {
		str.WriteString(log.String())
		str.WriteByte('\n')
	}
	for _, log := range p.Errors {
		str.WriteString(log.String())
		str.WriteByte('\n')
	}
	print(str.String())
	return len(p.Errors) > 0
}

func appendStandard(code *string) {
	year, month, day := time.Now().Date()
	hour, min, _ := time.Now().Clock()
	timeStr := fmt.Sprintf("%d/%d/%d %d.%d (DD/MM/YYYY) (HH.MM)",
		day, month, year, hour, min)
	var sb strings.Builder
	sb.WriteString("// Auto generated by JN compiler.\n")
	sb.WriteString("// JN compiler version:")
	sb.WriteString(jn.Version)
	sb.WriteByte('\n')
	sb.WriteString("// Date: ")
	sb.WriteString(timeStr)
	sb.WriteString("\n\n")
	sb.WriteString("// this file contains cxx module code which is automatically generated by JN")
	sb.WriteByte('\n')
	sb.WriteString("// compiler. generated code in this file provide cxx functions and structures")
	sb.WriteByte('\n')
	sb.WriteString("// corresponding to the definition in the JN source files")
	sb.WriteString("\n\n")
	sb.WriteString("\n\n#include \"")
	sb.WriteString(jnapi.JNCHeader)
	sb.WriteString("\"\n")
	sb.WriteString(*code)
	*code = sb.String()
}

func writeOutput(path, content string) {
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, 0o777)
	if err != nil {
		println(err.Error())
		os.Exit(0)
	}
	bytes := []byte(content)
	err = ioutil.WriteFile(path, bytes, 0o666)
	if err != nil {
		println(err.Error())
		os.Exit(0)
	}
}

func compile(path string, main, nolocal, justDefs bool) *Parser {
	loadJnSet()
	p := parser.New(nil)
	f, err := jnio.OpenJn(path)
	if err != nil {
		println(err.Error())
		return nil
	}
	if !jnio.IsUseable(path) {
		p.PushErr("file_not_useable")
		return p
	}
	inf, err := os.Stat(jn.StdlibPath)
	if err != nil || !inf.IsDir() {
		p.PushErr("no_stdlib")
		return p
	}
	p.File = f
	p.NoLocalPkg = nolocal
	p.Parsef(main, justDefs)
	return p
}

func execPostCommands() {
	for _, cmd := range jn.Set.PostCommands {
		fmt.Println(">", cmd)
		parts := strings.SplitN(cmd, " ", -1)
		err := exec.Command(parts[0], parts[1:]...).Run()
		if err != nil {
			println(err.Error())
		}
	}
}

func doSpell(path, cxx string) {
	defer execPostCommands()
	writeOutput(path, cxx)
	switch jn.Set.Mode {
	case jnset.ModeCompile:
		defer os.Remove(path)
		println("compilation is not supported yet")
	}
}

func main() {
	fpath := os.Args[0]
	p := compile(fpath, true, false, false)
	if p == nil {
		return
	}
	if printlogs(p) {
		os.Exit(0)
	}
	cxx := p.Cxx()
	appendStandard(&cxx)
	path := filepath.Join(jn.Set.CxxOutDir, jn.Set.CxxOutName)
	doSpell(path, cxx)
}
