package preprocessor

import "github.com/DeRuneLabs/jane/ast/models"

func TrimEnofi(tree *Tree) {
	for i, obj := range *tree {
		switch t := obj.Value.(type) {
		case models.Preprocessor:
			switch t := t.Command.(type) {
			case models.Directive:
				switch t.Command.(type) {
				case models.DirectiveEnofi:
					*tree = (*tree)[:i]
					return
				}
			}
		}
	}
}
