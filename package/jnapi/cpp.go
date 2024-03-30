package jnapi

var CppHeaderExtensions = []string{
	".h",
	".hpp",
	".hxx",
	".hh",
}

func IsValidHeader(ext string) bool {
	for _, validExt := range CppHeaderExtensions {
		if ext == validExt {
			return true
		}
	}
	return false
}
