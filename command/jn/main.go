package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/De-Rune/jane/documenter"
	"github.com/De-Rune/jane/package/jn"
	"github.com/De-Rune/jane/package/jnio"
	"github.com/De-Rune/jane/package/jnlog"
	"github.com/De-Rune/jane/package/jnset"
	"github.com/De-Rune/jane/parser"
)

func help(cmd string) {
	if cmd != "" {
		println("this module can only be used on a single")
	}
	helpmap := [][]string{
		{"help", "show help message"},
		{"version", "show version"},
		{"init", "initialize project"},
		{"doc", "documentize jn source code"},
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
		println("this module can only be used on a single")
		return
	}
	println("Jane Programming Language\n" + jn.Version)
}

func initProject(cmd string) {
	if cmd != "" {
		println("this module can only be used on a single")
		return
	}
	txt := []byte(`{
  "cxx_out_dir": "./dist/",
  "cxx_out_name": "jn.cxx",
  "out_name": "main",
  "language": "default"
}`)
	err := ioutil.WriteFile(jn.SettingsFile, txt, 0666)
	if err != nil {
		println(err.Error())
		os.Exit(0)
	}
	println("project initialized")
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
		path = filepath.Join(jn.JnSet.CxxOutDir, path+jn.DocExt)
		writeOuput(path, docjson)
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
	execp, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		println(err.Error())
		os.Exit(0)
	}
	jn.ExecPath = execp
	jn.StdlibPath = filepath.Join(jn.ExecPath, jn.Stdlib)
	jn.LangsPath = filepath.Join(jn.ExecPath, jn.Langs)

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
	loadJnSet()
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
		println("language's warning couldn't loaded (using default)")
		println(err.Error())
		return
	}
	err = json.Unmarshal(bytes, &jn.Warns)
	if err != nil {
		println("language's warning couldn't loaded (using default)")
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
		println("language's error couldn't loaded (using default)")
		println(err.Error())
		return
	}
	err = json.Unmarshal(bytes, &jn.Errs)
	if err != nil {
		println("language's error couldn't loaded (using default)")
		println(err.Error())
		return
	}
}

func loadLang() {
	lang := strings.TrimSpace(jn.JnSet.Language)
	if lang == "" || lang == "default" {
		return
	}
	path := filepath.Join(jn.LangsPath, lang)
	infos, err := ioutil.ReadDir(path)
	if err != nil {
		println("language's couldn't loaded (using default)")
		println(err.Error())
		return
	}
	loadLangWarns(path, infos)
	loadLangErrs(path, infos)
}

func loadJnSet() {
	info, err := os.Stat(jn.SettingsFile)
	if err != nil || info.IsDir() {
		println(`jn settings file ("` + jn.SettingsFile + `") not found`)
		os.Exit(0)
	}
	bytes, err := os.ReadFile(jn.SettingsFile)
	if err != nil {
		println(err.Error())
		os.Exit(0)
	}
	jn.JnSet, err = jnset.Load(bytes)
	if err != nil {
		println("jn settings errors;")
		println(err.Error())
		os.Exit(0)
	}
	loadLang()
}

func printlogs(p *parser.Parser) bool {
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
			str.WriteString("ERROR: ")
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
	timeStr := fmt.Sprintf("%d/%d/%d %d.%d (DD/MM/YYYY) (HH.MM)", day, month, year, hour, min)
	*code = `// JN compiler version:		` + jn.Version + `
// Date: 					` + timeStr + `
// Author:  				` + jn.Author + `
// License: 				` + jn.License + `

// this file contains cxx module code which is automatically generated by JN
// compiler. generated code in this file provide cxx functions and structures
// corresponding to the definition in the JN source files

// region JN_STANDARD_IMPORTS
#include <codecvt>
#include <cstdint>
#include <functional>
#include <iostream>
#include <locale>
#include <string>
#include <type_traits>
#include <vector>
// endregion JN_STANDARD_IMPORTS

// region JN_CXX_API
// region JN_BUILTIN_VALUES
#define nil nullptr
// endregion JN_BUILTIN_VALUES

// region JN_MISC
class exception : public std::exception {
private:
  std::basic_string<char> _buffer;

public:
  exception(const char *_Str) { this->_buffer = _Str; }
  const char *what() const throw() { return this->_buffer.c_str(); }
};

#define JNALLOC(_Alloc) new (std::nothrow) _Alloc
#define JNTHROW(_Msg) throw exception(_Msg)

template <typename _Enum_t, typename _Index_t, typename _Item_t>
static inline void foreach (
    const _Enum_t _Enum, const std::function<void(_Index_t, _Item_t)> _Body) {
  _Index_t _index{0};
  for (auto _item : _Enum) {
    _Body(_index++, _item);
  }
}

template <typename _Enum_t, typename _Index_t>
static inline void foreach (const _Enum_t _Enum,
                            const std::function<void(_Index_t)> _Body) {
  _Index_t _index{0};
  for (auto begin = _Enum.begin(), end = _Enum.end(); begin < end; ++begin) {
    _Body(_index++);
  }
}
// endregion JN_MISC

// region JN_BUILTIN_TYPES
typedef int8_t i8;
typedef int16_t i16;
typedef int32_t i32;
typedef int64_t i64;
typedef ssize_t ssize;
typedef uint8_t u8;
typedef uint16_t u16;
typedef uint32_t u32;
typedef uint64_t u64;
typedef size_t size;
typedef float f32;
typedef double f64;
typedef wchar_t rune;
// endregion JN_BUILTIN_TYPES

class str : public std::basic_string<rune> {
public:
  // region CONSTRUCTOR
  str(void) : str(L"") {}
  str(const rune *_Str) { this->assign(_Str); }
  str(const std::basic_string<rune> _Src) : str(_Src.c_str()) {}
  // endregion CONSTRUCTOR
};
// endregion JN_BUILTIN_TYPES

// region JN_STRUCTURES
template <typename _Item_t> class array {
public:
  // region FIELDS
  std::vector<_Item_t> _buffer;
  // endregion FIELDS

  // region CONSTRUCTORS
  array<_Item_t>(void) noexcept { this->_buffer = {}; }
  array<_Item_t>(const std::vector<_Item_t> &_Src) noexcept {
    this->_buffer = _Src;
  }
  array<_Item_t>(const std::nullptr_t) noexcept : array<_Item_t>() {}
  array<_Item_t>(const array<_Item_t> &_Src) noexcept
      : array<_Item_t>(_Src._buffer) {}

  array<_Item_t>(const str _Str) {
    if (std::is_same<_Item_t, rune>::value) {
      this->_buffer = std::vector<_Item_t>(_Str.begin(), _Str.end());
      return;
    }
    if (std::is_same<_Item_t, u8>::value) {
      std::wstring_convert<std::codecvt_utf8_utf16<rune>> _conv;
      std::string _bytes = _conv.to_bytes(_Str);
      this->_buffer = std::vector<_Item_t>(_bytes.begin(), _bytes.end());
      return;
    }
  }
  // endregion CONSTRUCTORS

  // region DESTRUCTOR
  ~array<_Item_t>(void) noexcept { this->_buffer.clear(); }
  // endregion DESTRUCTOR

  // region FOREACH_SUPPORT
  typedef _Item_t *iterator;
  typedef const _Item_t *const_iterator;
  iterator begin(void) noexcept { return &this->_buffer[0]; }
  const_iterator begin(void) const noexcept { return &this->_buffer[0]; }
  iterator end(void) noexcept { return &this->_buffer[this->_buffer.size()]; }
  const_iterator end(void) const noexcept {
    return &this->_buffer[this->_buffer.size()];
  }
  // endregion FOREACH_SUPPORT

  // region OPERATOR_OVERFLOWS
  operator str(void) const noexcept {
    if (std::is_same<_Item_t, rune>::value) {
      return str(std::basic_string<rune>(this->begin(), this->end()));
    }
    if (std::is_same<_Item_t, u8>::value) {
      std::wstring_convert<std::codecvt_utf8_utf16<rune>> _conv;
      const std::string _bytes(this->begin(), this->end());
      return str(_conv.from_bytes(_bytes));
    }
  }

  bool operator==(const array<_Item_t> &_Src) const noexcept {
    const size _length = this->_buffer.size();
    const size _Src_length = _Src._buffer.size();
    if (_length != _Src_length) {
      return false;
    }
    for (size _index = 0; _index < _length; ++_index) {
      if (this->_buffer[_index] != _Src._buffer[_index]) {
        return false;
      }
    }
    return true;
  }

  bool operator==(const std::nullptr_t) const noexcept {
    return this->_buffer.empty();
  }
  bool operator!=(const array<_Item_t> &_Src) const noexcept {
    return !(*this == _Src);
  }
  bool operator!=(const std::nullptr_t) const noexcept {
    return !this->_buffer.empty();
  }
  _Item_t &operator[](const size _Index) { return this->_buffer[_Index]; }

  friend std::wostream &operator<<(std::wostream &_Stream,
                                   const array<_Item_t> &_Src) {
    _Stream << L"[";
    const size _length = _Src._buffer.size();
    for (size _index = 0; _index < _length;) {
      _Stream << _Src._buffer[_index++];
      if (_index < _length) {
        _Stream << L", ";
      }
    }
    _Stream << L"]";
    return _Stream;
  }
  // endregion OPERATOR_OVERFLOWS
};
// endregion JN_STRUCTURES

// region JN_BUILTIN_FUNCTIONS
template <typename _Obj_t> static inline void _print(_Obj_t _Obj) {
  std::wcout << _Obj;
}

template <typename _Obj_t> static inline void _println(_Obj_t _Obj) {
  _print<_Obj_t>(_Obj);
  std::wcout << std::endl;
}
// endregion JN_BUILTIN_FUNCTIONS
// endregion JN_CXX_API

// region TRANSPILED_JN_CODE
` + *code + `
// endregion TRANSPILED_JN_CODE

// region JN_ENTRY_POINT
int main() {
// region JN_ENTRY_POINT_STANDARD_CODES
  std::locale::global(std::locale());
// endregion JN_ENTRY_POINT_STANDARD_CODES
  _main();
// region JN_ENTRY_POINT_END_STANDARD_CODES
  return EXIT_SUCCESS;
// endregion JN_ENTRY_POINT_END_STANDARD_CODES
}
// endergion JN_ENTRY_POINT`
}

func writeOuput(path, content string) {
	err := os.MkdirAll(jn.JnSet.CxxOutDir, 0777)
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

func compile(path string, main, justDefs bool) *parser.Parser {
	p := parser.New(nil)
	inf, err := os.Stat(jn.StdlibPath)
	if err != nil || !inf.IsDir() {
		p.Errs = append(p.Errs, jnlog.CompilerLog{
			Type: jnlog.FlatErr,
			Msg:  "standard library directory not found",
		})
		return p
	}
	f, err := jnio.OpenJn(path)
	if err != nil {
		println(err.Error())
		return nil
	}
	p.File = f
	p.Parsef(true, false)
	return p
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
	path := filepath.Join(jn.JnSet.CxxOutDir, jn.JnSet.CxxOutName)
	writeOuput(path, cxx)
}
