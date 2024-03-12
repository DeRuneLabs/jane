package jnapi

const Ignore = "_"

func IsIgnoreId(id string) bool {
	return id == Ignore
}

func AsId(id string) string {
	return "JNID(" + id + ")"
}
