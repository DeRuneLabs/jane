package jn

import "github.com/De-Rune/jane/package/jn/jnset"

const (
	Version      = `@dev_beta 0.0.1`
	Extension    = `.jn`
	SettingsFile = "jn.set"
	EntryPoint   = "main"
)

var (
	ExecutablePath string
	JnSet          *jnset.JnSet
)

func IsIgnoreName(name string) bool {
	return name == "__"
}
