package jn

import "github.com/De-Rune/jane/package/jnset"

const (
	Version      = `@dev_beta 0.0.1`
	SrcExt       = `.jn`
	DocExt       = ".jndoc"
	SettingsFile = "jn.set"
	Stdlib       = "lib"
	Langs        = "langs"
	Author       = "DeruneLabs"
	License      = "MIT LICENSE"

	EntryPoint = "main"
)

var (
	LangsPath  string
	StdlibPath string
	ExecPath   string
	JnSet      *jnset.JnSet
)
