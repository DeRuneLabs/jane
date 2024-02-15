package jane

import (
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
	jnset.Fields["out_dir"] = ""
	jnset.Fields["out_name"] = ""
	return jnset
}

func (jnset *JnSet) Parse(contet []byte) error {
	var lines []string
	if runtime.GOOS == "windows" {
		lines = strings.SplitN(string(contet), "\n", -1)
	} else {
		lines = strings.SplitN(string(contet), "\n\r", -1)
	}
	for index, line := range lines {
		line = strings.TrimFunc(line, unicode.IsSpace)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, " ", -1)
		if len(parts) < 2 {
			return fmt.Errorf("invalid syntax at line %d", index+1)
		}
		key, value := parts[0], parts[1]
		_, ok := jnset.Fields[key]
		if !ok {
			return fmt.Errorf("invalid field at line %d", index+1)
		}
		switch key {
		case "out_name":
			if len(parts) > 2 {
				return fmt.Errorf("invalid value at line %d", index+1)
			}
		}
		jnset.Fields[value] = value
	}
	return nil
}
