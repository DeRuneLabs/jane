package models

type Use struct {
	Tok        Tok
	Path       string
	Cpp        bool
	LinkString string
	FullUse    bool
	Selectors  []Tok
}
