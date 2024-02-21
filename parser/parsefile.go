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
	lexer := lexer.New(info.File)
	tokens := lexer.Tokenize()
	if lexer.Errors != nil {
		info.Errors = lexer.Errors
		return
	}
	parser := NewParser(tokens, info)
	parser.Parse()
	info.JN_CXX += parser.Cxx()
}
