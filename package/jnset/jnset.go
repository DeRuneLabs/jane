package jnset

import "encoding/json"

const (
	ModeTranspile = "transpile"
	ModeCompile   = "compile"
)

type JnSet struct {
	CxxOutDir    string   `json:"cxx_out_dir"`
	CxxOutName   string   `json:"cxx_out_name"`
	OutName      string   `json:"out_name"`
	Language     string   `json:"language"`
	Mode         string   `json:"mode"`
	PostCommands []string `json:"post_commands"`
	Indent       string   `json:"indent"`
	IndentCount  int      `json:"indent_count"`
}

var Default = &JnSet{
	CxxOutDir:    "./dist",
	CxxOutName:   "jn.cxx",
	OutName:      "main",
	Language:     "",
	Mode:         "transpile",
	Indent:       "\t",
	IndentCount:  1,
	PostCommands: []string{},
}

func Load(bytes []byte) (*JnSet, error) {
	set := *Default
	err := json.Unmarshal(bytes, &set)
	if err != nil {
		return nil, err
	}
	return &set, nil
}
