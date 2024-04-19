// Copyright (c) 2024 - DeRuneLabs
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package transpiler

import (
	"strings"

	"github.com/DeRuneLabs/jane/ast/models"
	"github.com/DeRuneLabs/jane/package/jnapi"
)

type trait struct {
	Ast     *models.Trait
	Defines *DefineMap
	Used    bool
	Desc    string
}

func (t *trait) has_reference_receiver() bool {
	for _, f := range t.Defines.Funcs {
		if typeIsRef(f.Ast.Receiver.Type) {
			return true
		}
	}
	return false
}

func (t *trait) FindFunc(id string) *Fn {
	for _, f := range t.Defines.Funcs {
		if f.Ast.Id == id {
			return f
		}
	}
	return nil
}

func (t *trait) OutId() string {
	return jnapi.OutId(t.Ast.Id, t.Ast.Token.File)
}

func (t *trait) String() string {
	var cpp strings.Builder
	cpp.WriteString("struct ")
	outid := t.OutId()
	cpp.WriteString(outid)
	cpp.WriteString(" {\n")
	models.AddIndent()
	is := models.IndentString()
	cpp.WriteString(is)
	cpp.WriteString("virtual ~")
	cpp.WriteString(outid)
	cpp.WriteString("(void) noexcept {}\n\n")
	for _, f := range t.Ast.Funcs {
		cpp.WriteString(is)
		cpp.WriteString("virtual ")
		cpp.WriteString(f.RetType.String())
		cpp.WriteByte(' ')
		cpp.WriteString(f.Id)
		cpp.WriteString(paramsToCpp(f.Params))
		cpp.WriteString(" {")
		if !typeIsVoid(f.RetType.Type) {
			cpp.WriteString(" return {}; ")
		}
		cpp.WriteString("}\n")
	}
	models.DoneIndent()
	cpp.WriteString("};")
	return cpp.String()
}
