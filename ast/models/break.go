package models

type Break struct {
	Tok Tok
}

func (b Break) String() string {
	return "break;"
}
