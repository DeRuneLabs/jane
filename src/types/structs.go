// copyright (c) 2024 arfy slowy - derunelabs
//
// permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "software"), to deal
// in the software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the software, and to permit persons to whom the software is
// furnished to do so, subject to the following conditions:
//
// the above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the software.
//
// the software is provided "as is", without warranty of any kind, express or
// implied, including but not limited to the warranties of merchantability,
// fitness for a particular purpose and noninfringement. in no event shall the
// authors or copyright holders be liable for any claim, damages or other
// liability, whether in an action of contract, tort or otherwise, arising from,
// out of or in connection with the software or the use or other dealings in the
// software.

package types

import "github.com/DeRuneLabs/jane/ast"

func IsStructOrdered(s *ast.Struct) bool {
	for _, d := range s.Origin.Depends {
		if d.Origin.Order > s.Origin.Order {
			return false
		}
	}
	return true
}

func OrderStructures(structures []*ast.Struct) {
	for i, s := range structures {
		s.Order = i
	}

	n := len(structures)
	for i := 0; i < n; i++ {
		swapped := false
		for j := 0; j < n-i-1; j++ {
			curr := &structures[j]
			if !IsStructOrdered(*curr) {
				(*curr).Origin.Order = j + 1
				next := &structures[j+1]
				(*next).Origin.Order = j
				*curr, *next = *next, *curr
				swapped = true
			}
		}
		if !swapped {
			break
		}
	}
}
