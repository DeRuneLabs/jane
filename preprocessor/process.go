package preprocessor

import "github.com/DeRuneLabs/jane/ast/models"

type Tree = []models.Object

func Process(tree *Tree, includeEnofi bool) {
	if includeEnofi {
		TrimEnofi(tree)
	}
}
