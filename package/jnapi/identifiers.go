package jnapi

import (
	"fmt"
	"strings"
	"unsafe"

	"github.com/DeRuneLabs/jane/package/jnio"
)

const Ignore = "_"

func IsIgnoreId(id string) bool {
	return id == Ignore
}

func AsId(id string) string {
	return "JNID(" + id + ")"
}

func getPtrAsId(ptr unsafe.Pointer) string {
	address := fmt.Sprintf("%p", ptr)
	address = address[3:]
	for i, r := range address {
		if r != '0' {
			address = address[i:]
			break
		}
	}
	return address
}

func OutId(id string, f *jnio.File) string {
	if f != nil {
		var out strings.Builder
		out.WriteByte('f')
		out.WriteString(getPtrAsId(unsafe.Pointer(f)))
		out.WriteByte('_')
		out.WriteString(id)
		return out.String()
	}
	return AsId(id)
}

func AsTypeId(id string) string {
	return id + "_jnt"
}
