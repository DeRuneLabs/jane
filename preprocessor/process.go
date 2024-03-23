package preprocessor

import "github.com/DeRuneLabs/jane/ast/models"

type Tree = []models.Object

func Process(tree *Tree) {
	TrimEnofi(tree)
}
