package jn

import "fmt"

var Warns = map[string]string{
	`doc_ignored`:         `documentation is ignored because object isn't supports documentations`,
	`exist_undefined_doc`: `your source code has undefined documentations (some documentations isn't document anything)`,
}

func GetWarn(key string, args ...interface{}) string {
	return fmt.Sprintf(Warns[key], args...)
}
