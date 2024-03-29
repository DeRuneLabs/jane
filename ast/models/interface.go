package models

type CompiledStruct interface {
	OutId() string
	Generics() []DataType
	SetGenerics([]DataType)
}

type Genericable interface {
	Generics() []DataType
	SetGenerics([]DataType)
}

type IterProfile interface {
	String(iter Iter) string
}

type IExprModel interface {
	String() string
}
