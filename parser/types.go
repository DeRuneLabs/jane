package parser

import "github.com/De-Rune/jane/package/jane"

func typeFromName(name string) uint {
	switch name {
	case "int8":
		return jane.Int8
	case "int16":
		return jane.Int16
	case "int32":
		return jane.Int32
	case "int64":
		return jane.Int64
	case "uint8":
		return jane.UInt8
	case "uint16":
		return jane.UInt16
	case "uint32":
		return jane.UInt32
	case "uint64":
		return jane.UInt64

	}
	return 0
}

func cxxTypeNameFromType(typecode uint) string {
	switch typecode {
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
	}
	return ""
}
