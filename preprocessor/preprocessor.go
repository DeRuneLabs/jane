package preprocessor

import "github.com/De-Rune/jane/ast"

func TrimEnofi(tree *[]ast.Obj) {
	for i, obj := range *tree {
		switch t := obj.Value.(type) {
		case ast.Preprocessor:
			switch t := t.Command.(type) {
			case ast.Directive:
				switch t.Command.(type) {
				case ast.EnofiDirective:
					*tree = (*tree)[:i]
					return
				}
			}
		}
	}
}
