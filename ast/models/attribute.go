package models

type Attribute struct {
	Tok Tok
	Tag string
}

func (a Attribute) String() string {
	return a.Tag
}
