package models

type Args struct {
	Src                      []Arg
	Targeted                 bool
	Generics                 []DataType
	DynamicGenericAnnotation bool
	NeedsPureType            bool
}
