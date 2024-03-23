package models

type Struct struct {
	Tok      Tok
	Id       string
	Pub      bool
	Fields   []*Var
	Generics []*GenericType
}
