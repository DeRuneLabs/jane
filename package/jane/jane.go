package jane

const (
	Version      = `@dev_beta 1.0.0`
	Extension    = `.jn`
	SettingsFile = "jane.set"
	EntryPoint   = "main"
)

var (
	ExecutablePath string
	JaneSettings   *JnSet
)
