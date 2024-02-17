package ast

const (
	NA         uint8 = 0
	Identifier uint8 = 1
	Statement  uint8 = 2
	Range      uint8 = 3
	Block      uint8 = 4
	Type       uint8 = 5
)

const (
	IdentifierNA   uint8 = 0
	IdentifierName uint8 = 1
)

const (
	StatementNA       uint8 = 0
	StatementFunction uint8 = 1
	StatementReturn   uint8 = 2
)

const (
	RangeNA          uint8 = 0
	RangeBrace       uint8 = 1
	RangeParentheses uint8 = 2
)
