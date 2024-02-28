package jnset

import (
	"encoding/json"
)

type JnSet struct {
	CxxOutDir  string `json:"cxx_out_dir"`
	CxxOutName string `json:"cxx_out_name"`
	OutName    string `json:"out_name"`
}

func Load(jsonbytes []byte) (*JnSet, error) {
	set := JnSet{
		CxxOutDir:  "./dist",
		CxxOutName: "jn.cxx",
		OutName:    "main",
	}
	err := json.Unmarshal(jsonbytes, &set)
	if err != nil {
		return nil, err
	}
	return &set, nil
}
