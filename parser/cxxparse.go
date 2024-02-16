package parser

import (
	"fmt"
	"strings"

	"github.com/De-Rune/jane/lexer"
	"github.com/De-Rune/jane/package/jane"
)

type CxxParser struct {
	Functions []*Function
	Position  int
	Tokens    []lexer.Token
	PFI       *ParseFileInfo
}

func NewParser(tokens []lexer.Token, PFI *ParseFileInfo) *CxxParser {
	parser := new(CxxParser)
	parser.Tokens = tokens
	parser.PFI = PFI
	parser.Position = 0
	return parser
}

func (cp *CxxParser) PushError(err string) {
	message := jane.Errors[err]
	cp.PFI.Errors = append(cp.PFI.Errors, fmt.Sprintf("%s %s: %d", cp.PFI.File.Path, message, cp.Tokens[cp.Position].Line))
}

func (cp CxxParser) String() string {
	var sb strings.Builder
	for _, function := range cp.Functions {
		sb.WriteString(function.String())
	}
	return sb.String()
}

func (cp CxxParser) Ended() bool {
	return cp.Position >= len(cp.Tokens)
}

func (cp *CxxParser) Parse() {
	for cp.Position != -1 && !cp.Ended() {
		firsToken := cp.Tokens[cp.Position]
		switch firsToken.Type {
		case lexer.Name:
			cp.processName()
		default:
			cp.PushError("invalid_syntax")
		}
	}
}

func (cp *CxxParser) ParseFunction() {
	defineToken := cp.Tokens[cp.Position]
	function := new(Function)
	function.Name = defineToken.Value
	function.Line = defineToken.Line
	function.FILE = defineToken.File
	function.ReturnType = jane.Void

	cp.Position += 3
	if cp.Ended() {
		cp.Position--
		cp.PushError("function_body")
		cp.Position = -1
		return
	}
	token := cp.Tokens[cp.Position]
	if token.Type == lexer.Type {
		function.ReturnType = typeFromName(token.Value)
		cp.Position++
		if cp.Ended() {
			cp.Position--
			cp.PushError("function_body")
			cp.Position = -1
			return
		}
		token = cp.Tokens[cp.Position]
	}
	switch token.Type {
	case lexer.Brace:
		if token.Value != "{" {
			cp.PushError("invalid_syntax")
			cp.Position = -1
			return
		}
		cp.Position += 3
	default:
		cp.PushError("invalid_syntax")
		cp.Position = -1
		return
	}
	cp.Functions = append(cp.Functions, function)
}

func (cp *CxxParser) processName() {
	cp.Position++
	if cp.Ended() {
		cp.Position--
		cp.PushError("invalid_syntax")
		return
	}
	cp.Position--
	secondToken := cp.Tokens[cp.Position+1]
	switch secondToken.Type {
	case lexer.Brace:
		switch secondToken.Value {
		case "(":
			cp.ParseFunction()
		default:
			cp.PushError("invalid_syntax")
		}
	}
}
