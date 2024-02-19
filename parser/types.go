package parser

import "github.com/De-Rune/jane/package/jane"

func cxxTypeNameFromType(typeCode uint8) string {
	switch typeCode {
	case jane.Void:
		return "void"
	case jane.Int8:
		return "signed char"
	case jane.Int16:
		return "short"
	case jane.Int32:
		return "int"
	case jane.Int64:
		return "long"
	case jane.UInt8:
		return "unsigned char"
	case jane.UInt16:
		return "unsigned short"
	case jane.UInt32:
		return "unsigned int"
	case jane.UInt64:
		return "unsigned long"
	case jane.Boolean:
		return "bool"
	case jane.Float32:
		return "float"
	case jane.Float64:
		return "double"
	}
	return ""
}
