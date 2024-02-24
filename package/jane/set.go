package jane

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
	"unicode"
)

type JnSet struct {
	Fields map[string]string
}

func NewJnSet() *JnSet {
	jnset := new(JnSet)
	jnset.Fields = make(map[string]string)
	jnset.Fields["out_name"] = ""
	jnset.Fields["cxx_out_dir"] = ""
	jnset.Fields["cxx_out_name"] = ""
	return jnset
}

func splitLines(content string) []string {
	if runtime.GOOS == "windows" {
		return strings.SplitN(string(content), "\n", 1)
	}
	return strings.SplitN(string(content), "\n", -1)
}

func (jnset *JnSet) checkUnset() []error {
	var errs []error
	for key := range jnset.Fields {
		if jnset.Fields[key] == "" {
			errs = append(errs, errors.New("\""+key+"\" is not define"))
		}
	}
	return errs
}

func (jnset *JnSet) Parse(content []byte) []error {
	lines := splitLines(string(content))
	for index, line := range lines {
		line = strings.TrimFunc(line, unicode.IsSpace)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, " ", -1)
		if len(parts) < 2 {
			return []error{fmt.Errorf("invalid syntax at line %d", index+1)}
		}
		key := parts[0]
		value := strings.TrimFunc(parts[1], unicode.IsSpace)
		_, ok := jnset.Fields[key]
		if !ok {
			return []error{fmt.Errorf("invalid field at line %d", index+1)}
		}
		switch key {
		case "out_name":
			if len(parts) < 2 {
				return []error{fmt.Errorf("invalid value at line %d", index+1)}
			}
		}
		jnset.Fields[key] = value
	}
	return jnset.checkUnset()
}
