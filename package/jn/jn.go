package jn

import "github.com/DeRuneLabs/jane/package/jnset"

const (
	Version       = `@dev_beta 0.0.1`
	SrcExt        = `.jn`
	DocExt        = ".jndoc"
	SettingsFile  = "jn.set"
	Stdlib        = "lib"
	Localizations = "localization"

	EntryPoint = "main"

	Anonymous = "<anonymous>"

	DocPrefix = "doc:"

	PlatformWindows = "windows"
	PlatformLinux   = "linux"
	PlatformDarwin  = "darwin"

	ArchArm   = "arm"
	ArchArm64 = "arm64"
	ArchAmd64 = "amd64"
	ArchI386  = "i386"
)

var (
	LangsPath  string
	StdlibPath string
	ExecPath   string
	Set        *jnset.JnSet
)
