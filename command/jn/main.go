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
	"github.com/DeRuneLabs/jane/package/jnlog"
	"github.com/DeRuneLabs/jane/package/jnset"
	"github.com/DeRuneLabs/jane/parser"
)

type Parser = parser.Parser

func help(cmd string) {
	if cmd != "" {
		println("this module can only be using as single")
		return
	}
	helpmap := [][]string{
		{"help", "show help"},
		{"version", "show version"},
		{"init", "initialize jane module project"},
		{"doc", "documenting jane source code"},
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
		println("this module can only be using as single")
		return
	}
	println("jane version", jn.Version)
}

func initProject(cmd string) {
	if cmd != "" {
		println("this module can only be using as single")
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
	println("intialized project")
}

func doc(cmd string) {
	cmd = strings.TrimSpace(cmd)
	paths := strings.SplitN(cmd, " ", -1)
	for _, path := range paths {
		path = strings.TrimSpace(path)
		p := compile(path, false, true)
		if p == nil {
			continue
		}
		if printlogs(p) {
			fmt.Println(jn.GetErr("doc_couldnt_generated", path))
			continue
		}
		docjson, err := documenter.Documentize(p)
		if err != nil {
			fmt.Println(jn.GetErr("error", err.Error()))
			continue
		}
		path = path[len(filepath.Dir(path)):]
		path = filepath.Join(jn.Set.CxxOutDir, path+jn.DocExt)
		writeOutput(path, docjson)
	}
}

func processCommand(namespace, cmd string) bool {
	switch namespace {
	case "help":
		help(cmd)
	case "version":
		version(cmd)
	case "init":
		initProject(cmd)
	case "doc":
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
		if f.IsDir() || f.Name() != "warns.json" {
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
		println("language warning's couldn't loaded (use default)")
		println(err.Error())
		return
	}
	err = json.Unmarshal(bytes, &jn.Warns)
	if err != nil {
		println("language warning's couldn't loaded (use default)")
		println(err.Error())
		return
	}
}

func loadLangErrs(path string, infos []fs.FileInfo) {
	i := -1
	for j, f := range infos {
		if f.IsDir() || f.Name() != "errs.json" {
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
		println("language warning's couldn't loaded (use default)")
		println(err.Error())
		return
	}
	err = json.Unmarshal(bytes, &jn.Errs)
	if err != nil {
		println("language warning's couldn't loaded (use default)")
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
		println("language warning's couldn't loaded (use default)")
		println(err.Error())
		return
	}
	loadLangWarns(path, infos)
	loadLangErrs(path, infos)
}

func checkMode() {
	lower := strings.ToLower(jn.Set.Mode)
	if lower != jnset.ModeTranspile && lower != jnset.ModeCompile {
		key, _ := reflect.TypeOf(jn.Set).Elem().FieldByName("Mode")
		tag := string(key.Tag)
		tag = tag[6 : len(tag)-1]
		println(jn.GetErr("invalid_value_for_key", jn.Set.Mode, tag))
		os.Exit(0)
	}
	jn.Set.Mode = lower
}

func loadJnSet() {
	info, err := os.Stat(jn.SettingsFile)
	if err != nil || info.IsDir() {
		println(`Jane settings file ("` + jn.SettingsFile + `") is not found!`)
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
	for _, log := range p.Warns {
		switch log.Type {
		case jnlog.FlatWarn:
			str.WriteString("WARNING: ")
			str.WriteString(log.Msg)
		case jnlog.Warn:
			str.WriteString("WARNING: ")
			str.WriteString(log.Path)
			str.WriteByte(':')
			str.WriteString(fmt.Sprint(log.Row))
			str.WriteByte(':')
			str.WriteString(fmt.Sprint(log.Column))
			str.WriteByte(' ')
			str.WriteString(log.Msg)
		}
		str.WriteByte('\n')
	}
	for _, log := range p.Errs {
		switch log.Type {
		case jnlog.FlatErr:
			str.WriteString("ERROR: ")
			str.WriteString(log.Msg)
		case jnlog.Err:
			str.WriteString(log.Path)
			str.WriteByte(':')
			str.WriteString(fmt.Sprint(log.Row))
			str.WriteByte(':')
			str.WriteString(fmt.Sprint(log.Column))
			str.WriteByte(' ')
			str.WriteString(log.Msg)
		}
		str.WriteByte('\n')
	}
	print(str.String())
	return len(p.Errs) > 0
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
	sb.WriteString("\n")
	sb.WriteString("// compiler. generated code in this file provide cxx functions and structures")
	sb.WriteString("\n")
	sb.WriteString("// corresponding to the definition in the JN source files")
	sb.WriteString("\n\n")
	sb.WriteString(jnapi.CxxDefault)
	sb.WriteString("\n\n// region TRANSPILED_JN_CODE\n")
	sb.WriteString(*code)
	sb.WriteString("\n// endregion TRANSPILED_JN_CODE\n\n")
	sb.WriteString(jnapi.CxxMain)
	*code = sb.String()
}

func writeOutput(path, content string) {
	err := os.MkdirAll(jn.Set.CxxOutDir, 0777)
	if err != nil {
		println(err.Error())
		os.Exit(0)
	}
	bytes := []byte(content)
	err = ioutil.WriteFile(path, bytes, 0666)
	if err != nil {
		println(err.Error())
		os.Exit(0)
	}
}

func loadBuiltin() bool {
	f, err := jnio.OpenJn(filepath.Join(jn.StdlibPath, "lib.jn"))
	p := parser.New(f)
	if err != nil {
		println(err.Error())
		return false
	}
	p.Defs = parser.Builtin
	p.Parsef(false, false)
	return true
}

func compile(path string, main, justDefs bool) *Parser {
	loadJnSet()
	p := parser.New(nil)
	f, err := jnio.OpenJn(path)
	if err != nil {
		println(err.Error())
		return nil
	}
	if !jnio.IsUseable(path) {
		p.Errs = append(p.Errs, jnlog.CompilerLog{
			Type: jnlog.FlatErr,
			Msg:  "file is not useable for this platform",
		})
		return p
	}
	inf, err := os.Stat(jn.StdlibPath)
	if err != nil || !inf.IsDir() {
		p.Errs = append(p.Errs, jnlog.CompilerLog{
			Type: jnlog.FlatErr,
			Msg:  "standard library directory not found",
		})
		return p
	}
	if !loadBuiltin() {
		return nil
	}
	p.File = f
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
	p := compile(fpath, true, false)
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
