package parser

import (
	"sync"

	"github.com/De-Rune/jane/lexer"
	"github.com/De-Rune/jane/package/io"
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
	parser := NewParser(tokens, info)
	parser.Parse()
	code := parser.String()
	info.JN_CXX += code
}
