package parser

import (
	"fmt"
	"sync"

	"github.com/slowy07/jane/lexer"
	"github.com/slowy07/jane/package/io"
)

type ParseFileInfo struct {
	JN_CXX   string
	Errors   []string
	File     *io.FILE
	Routines *sync.WaitGroup
}

func ParseFile(info *ParseFileInfo) {
	defer info.Routines.Done()
	info.JN_CXX = ""
	jnlexer := lexer.New(info.File)
	tokens := jnlexer.Tokenize()
	if jnlexer.Errors != nil {
		info.Errors = jnlexer.Errors
		return
	}
	for _, token := range tokens {
		fmt.Println("'"+token.Value+"'", " ")
	}
}
