package models

type CxxEmbed struct {
	Tok     Tok
	Content string
}

func (ce CxxEmbed) String() string {
	return ce.Content
}
