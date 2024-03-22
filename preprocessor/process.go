package preprocessor

import (
	"github.com/DeRuneLabs/jane/ast"
)

type Tree = []ast.Obj

func Process(tree *Tree) {
	TrimEnofi(tree)
}
