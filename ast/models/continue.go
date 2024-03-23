package models

type Continue struct {
	Tok Tok
}

func (c Continue) String() string {
	return "continue;"
}
