package jnset

import (
	"encoding/json"
)

type JnSet struct {
	CxxOutDir  string `json:"cxx_out_dir"`
	CxxOutName string `json:"cxx_out_name"`
	OutName    string `json:"out_name"`
	Language   string `json:"language"`
}

func Load(bytes []byte) (*JnSet, error) {
	set := JnSet{
		CxxOutDir:  "./dist",
		CxxOutName: "jn.cxx",
		OutName:    "main",
		Language:   "",
	}
	err := json.Unmarshal(bytes, &set)
	if err != nil {
		return nil, err
	}
	return &set, nil
}
