package parser

func IsString(value string) bool {
	return value[0] == '"'
}

func IsBoolean(value string) bool {
	return value == "true" || value == "false"
}
