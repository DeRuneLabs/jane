package models

type Value struct {
	Tok  Tok
	Data string
	Type DataType
}

func (v Value) String() string {
	return v.Data
}
