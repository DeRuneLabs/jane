package models

import "fmt"

type Preprocessor struct {
	Tok     Tok
	Command any
}

func (pp Preprocessor) String() string {
	return fmt.Sprint(pp.Command)
}

type Directive struct {
	Command any
}

func (d Directive) String() string {
	return fmt.Sprint(d.Command)
}

type DirectiveEnofi struct{}
