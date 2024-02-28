package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/De-Rune/jane/package/io"
	"github.com/De-Rune/jane/package/jn"
	"github.com/De-Rune/jane/package/jn/jnset"
	"github.com/De-Rune/jane/parser"
)

func help(cmd string) {
	if cmd != "" {
		println("this module only used a single")
		return
	}
	helpContent := [][]string{
		{"help", "show help message"},
		{"version", "show version"},
		{"init", "initialize project here"},
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
		sb.WriteString(strings.Repeat(" ", (maxlen-len(part[0]))+space))
		sb.WriteString(part[1])
		sb.WriteByte('\n')
	}
	println(sb.String()[:sb.Len()-1])
}

func version(cmd string) {
	if cmd != "" {
		println("this module can only used a single")
	}
	println("Jane Programming Language\n" + jn.Version)
}

func initProject(cmd string) {
	if cmd != "" {
		println("this module can only be used as single")
		return
	}
	err := io.WriteFileTruncate(jn.SettingsFile, []byte(`{
  "cxx_out_dir": "./dist/",
  "cxx_out_name": "jn.cxx",
  "out_name": "main"
}`))
	if err != nil {
		println(err.Error())
		os.Exit(1)
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
	jn.ExecutablePath = filepath.Dir(os.Args[0])
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
	info, err := os.Stat(jn.SettingsFile)
	if err != nil || info.IsDir() {
		println(`jn settings file ("` + jn.SettingsFile + `") not found`)
		os.Exit(1)
	}
	bytes, err := os.ReadFile(jn.SettingsFile)
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
	jn.JnSet, err = jnset.Load(bytes)
	if err != nil {
		println("jn settings has error")
		println(err.Error())
		os.Exit(1)
	}
}

func printErrors(errors []string) {
	defer os.Exit(1)
	for _, message := range errors {
		fmt.Println(message)
	}
}

func appendStandards(code *string) {
	year, month, day := time.Now().Date()
	hour, min, _ := time.Now().Clock()
	timeString := fmt.Sprintf("%d/%d/%d %d.%d (DD/MM/YYYY) (HH.MM)", day, month, year, hour, min)
	*code = `// code generate by jn compiler
// jn compiler version:` + jn.Version + `
// Date:    ` + timeString + `

#pragma region JN_STANDARD_IMPORTS
#include <iostream>
#include <string>
#include <functional>
#include <vector>
#include <locale.h>
#include <cstdint>
#pragma endregion JN_STANDARD_IMPORTS

#pragma region JN_RUNTIME_FUNCTIONS
inline void throw_exception(const std::wstring message) {
  std::wcout << message << std::endl;
  exit(1);
}
#pragma endregion JN_RUNTIME_FUNCTION

#pragma region JN_BUILTIN_TYPES
typedef int8_t i8;
typedef int16_t i16;
typedef int32_t i32;
typedef int64_t i64;
typedef uint8_t u8;
typedef uint16_t u16;
typedef uint32_t u32;
typedef uint64_t u64;
typedef float f32;
typedef double f64;
typedef wchar_t rune;

#define function std::function

class str {
public:
#pragma region FIELDS
  std::wstring string;
#pragma endregion FIELDS

#pragma region CONSTRUCTORS
  str(const std::wstring& string) {
    this->string = string;
  }

  str(const rune* string) {
    this->string = string;
  }
#pragma endregion CONSTRUCTORS

#pragma region DESTRUCTOR
  ~str() {
    this->string.clear();
  }
#pragma endregion DESTRUCTOR

#pragma region OPERATOR_OVERFLOWS
  bool operator==(const str& string) {
    return this->string == string.string;
  }

  bool operator!=(const str& string) {
    return !(this->string == string.string);
  }

  str operator+(const str& string) {
    return str(this->string + string.string);
  }

  void operator+=(const str& string) {
    this->string += string.string;
  }

  rune& operator[](const int index) {
    const u32 length = this->string.length();
    if (index < 0) {
      throw_exception(L"ERR: stackoverflow exception:\n index is less than zero");
    } else if (index >= length) {
      throw_exception(L"ERR: stackoverflow exception:\nindex overflow" + std::to_wstring(index) + L":" + std::to_wstring(length));
    }
    return this->string[index];
  }

  friend std::wostream& operator<<(std::wostream &os, const str& string) {
    os << string.string;
    return os;
  }
#pragma endregion OPERATOR_OVERFLOWS
};
#pragma endregion JN_BUILTIN_TYPES

#pragma region JN_BUILTIN_VALUES
#define nil nullptr
#pragma endregion JN_BUILTIN_VALUES

#pragma region JN_STRUCTURES
template <typename T>
class array {
public:
#pragma region FIELDS
  std::vector<T> vector;
  bool heap;
#pragma endregion FIELDS

#pragma region CONSTRUCTORS
  array() {
    this->vector = {};
    this->heap = false;
  }

  array(std::nullptr_t) : array() {}

  array(const std::vector<T>& vector, bool heap) {
    this->vector = vector;
    this->heap = heap;
  }

  array(const std::vector<T>& vector) : array(vector, false) {}
#pragma endregion CONSTRUCTORS

#pragma region DESTRUCTOR
  ~array() {
    this->vector.clear();
    if (this->heap) {
      delete this;
    }
  }
#pragma endregion DESTRUCTOR

#pragma region OPERATOR_OVERFLOWS
  bool operator==(const array& array) {
    const u32 vector_length = this->vector.size();
    const u32 array_vector_length = array.vector.size();
    if (vector_length != array_vector_length) {
      return false;
    }
    for (int index = 0; index < vector_length; ++index) {
      if (this->vector[index] != array.vector[index]) {
        return false;
      }
    }
    return true;
  }

  bool operator==(std::nullptr_t) {
    return this->vector.empty();
  }

  bool operator!=(const array& array) {
    return !(&this == array);
  }

  bool operator!=(std::nullptr_t) {
    return !this->vector.empty();
  }

  T& operator[](const int index) {
    const u32 length = this->vector.size();
    if (index < 0) {
      throw_exception(L"ERR: stackoverflow exception:\n index is less than zero");
    } else if (index >= length) {
      throw_exception(L"ERR: stackoverflow exception:\nindex overflow" + std::to_wstring(index) + L":" + std::to_wstring(length));
    }
    return this->vector[index];
  }

  friend std::wostream& operator<<(std::wostream &os, const array<T>& array) {
    os << L"[";
    const u32 size = array.vector.size();
    for (int index = 0; index < size;) {
      os << array.vector[index++];
      if (index < size) {
        os << L", ";
      }
    }
    os << L"]";
    return os;
  }
#pragma endregion OPERATOR_OVERFLOWS
};
#pragma endregion JN_STRUCTURES

#pragma region JN_BUILTIN_FUNCTIONS
template<typename any>
inline void _disp(any v) {
  std::wcout << v;
}

template <typename any>
inline void _displn(any v) {
  _disp(v);
  std::wcout << std::endl;
}
#pragma endregion JN_BUILTIN_FUNCTIONS

#pragma region TRANSPILED_JN_CODE
` + *code + `
#pragma endregion TRANSPILED_JN_CODE

#pragma region JN_ENTRY_POINT
int main() {
#pragma region JN_ENTRY_POINT_STANDARD_CODES
  setlocale(0x0, "");
#pragma endregion JN_ENTRY_POINT_STANDARD_CODES
_main();

#pragma region JN_ENTRY_POINT_END_STANDARD_CODES
  return EXIT_SUCCESS;
#pragma endregion JN_ENTRY_POINT_END_STANDARD_CODES
}
#pragma endregion JN_ENTRY_POINT`
}

func writeCxxOutput(info *parser.ParseFileInfo) {
	path := filepath.Join(jn.JnSet.CxxOutDir, jn.JnSet.CxxOutName)
	err := os.MkdirAll(jn.JnSet.CxxOutDir, 0511)
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
	err = io.WriteFileTruncate(path, []byte(info.JN_CXX))
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
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
	writeCxxOutput(info)
}
