package jane

const (
	Version      = `@dev_ver 0.0.1`
	Extension    = `.jn`
	SettingsFile = "jn.set"
	EntryPoint   = "main"
)

var (
	ExecutablePath string
	JnSettings     *JnSet
)
