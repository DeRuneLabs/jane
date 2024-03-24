package models

type Data struct {
	Tok   Tok
	Value string
	Type  DataType
}

func (d Data) String() string {
	return d.Value
}
