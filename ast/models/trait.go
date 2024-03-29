package models

type Trait struct {
	Pub   bool
	Tok   Tok
	Id    string
	Desc  string
	Used  bool
	Funcs []*Func
}
