package models

type Use struct {
	Tok        Tok
	Path       string
	LinkString string
	FullUse    bool
	Selectors  []Tok
}
