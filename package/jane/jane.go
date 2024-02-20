package jane

const (
	Version     = `@dev_beta 0.0.1`
	Extension   = `.jn`
	SettingFile = "jane.set"
	EntryPoint  = "main"
)

var (
	ExecutablePath string
	JaneSettings   *JnSet
)
