package models

type Attribute struct {
	Tok Tok
	Tag Tok
}

func (a Attribute) String() string {
	return a.Tag.Kind
}
