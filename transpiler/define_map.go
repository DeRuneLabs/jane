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
	"github.com/DeRuneLabs/jane/package/jnio"
	"github.com/DeRuneLabs/jane/package/jntype"
)

func isAccessable(finder, target *File, defIsPub bool) bool {
	return defIsPub || finder == nil || target == nil || finder.Dir == target.Dir
}

type DefineMap struct {
	Namespaces []*namespace
	Enums      []*Enum
	Structs    []*structure
	Traits     []*trait
	Types      []*TypeAlias
	Funcs      []*Fn
	Globals    []*Var
	side       *DefineMap
}

func (dm *DefineMap) findNsById(id string) int {
	for i, t := range dm.Namespaces {
		if t != nil && t.Id == id {
			return i
		}
	}
	return -1
}

func (dm *DefineMap) nsById(id string) *namespace {
	i := dm.findNsById(id)
	if i == -1 {
		return nil
	}
	return dm.Namespaces[i]
}

func (dm *DefineMap) findStructById(id string, f *File) (int, *DefineMap, bool) {
	for i, s := range dm.Structs {
		if s != nil && s.Ast.Id == id {
			if isAccessable(f, s.Ast.Token.File, s.Ast.Pub) {
				return i, dm, false
			}
		}
	}
	if dm.side != nil {
		i, m, _ := dm.side.findStructById(id, f)
		return i, m, true
	}
	return -1, nil, false
}

func (dm *DefineMap) structById(id string, f *File) (*structure, *DefineMap, bool) {
	i, m, canshadow := dm.findStructById(id, f)
	if i == -1 {
		return nil, nil, false
	}
	return m.Structs[i], m, canshadow
}

func (dm *DefineMap) findTraitById(id string, f *File) (int, *DefineMap, bool) {
	for i, t := range dm.Traits {
		if t != nil && t.Ast.Id == id {
			if isAccessable(f, t.Ast.Token.File, t.Ast.Pub) {
				return i, dm, false
			}
		}
	}
	if dm.side != nil {
		i, m, _ := dm.side.findTraitById(id, f)
		return i, m, true
	}
	return -1, nil, false
}

func (dm *DefineMap) traitById(id string, f *File) (*trait, *DefineMap, bool) {
	i, m, canshadow := dm.findTraitById(id, f)
	if i == -1 {
		return nil, nil, false
	}
	return m.Traits[i], m, canshadow
}

func (dm *DefineMap) findEnumById(id string, f *File) (int, *DefineMap, bool) {
	for i, e := range dm.Enums {
		if e != nil && e.Id == id {
			if isAccessable(f, e.Token.File, e.Pub) {
				return i, dm, false
			}
		}
	}
	if dm.side != nil {
		i, m, _ := dm.side.findEnumById(id, f)
		return i, m, true
	}
	return -1, nil, false
}

func (dm *DefineMap) enumById(id string, f *File) (*Enum, *DefineMap, bool) {
	i, m, canshadow := dm.findEnumById(id, f)
	if i == -1 {
		return nil, nil, false
	}
	return m.Enums[i], m, canshadow
}

func (dm *DefineMap) findTypeById(id string, f *File) (int, *DefineMap, bool) {
	for i, t := range dm.Types {
		if t != nil && t.Id == id {
			if isAccessable(f, t.Token.File, t.Pub) {
				return i, dm, false
			}
		}
	}
	if dm.side != nil {
		i, m, _ := dm.side.findTypeById(id, f)
		return i, m, true
	}
	return -1, nil, false
}

func (dm *DefineMap) typeById(id string, f *File) (*TypeAlias, *DefineMap, bool) {
	i, m, canshadow := dm.findTypeById(id, f)
	if i == -1 {
		return nil, nil, false
	}
	return m.Types[i], m, canshadow
}

func (dm *DefineMap) findFuncById(id string, f *File) (int, *DefineMap, bool) {
	for i, fn := range dm.Funcs {
		if fn != nil && fn.Ast.Id == id {
			if isAccessable(f, fn.Ast.Token.File, fn.Ast.Pub) {
				return i, dm, false
			}
		}
	}
	if dm.side != nil {
		i, m, _ := dm.side.findFuncById(id, f)
		return i, m, true
	}
	return -1, nil, false
}

func (dm *DefineMap) funcById(id string, f *File) (*Fn, *DefineMap, bool) {
	i, m, canshadow := dm.findFuncById(id, f)
	if i == -1 {
		return nil, nil, false
	}
	return m.Funcs[i], m, canshadow
}

func (dm *DefineMap) findGlobalById(id string, f *File) (int, *DefineMap, bool) {
	for i, g := range dm.Globals {
		if g != nil && g.Type.Id != jntype.Void && g.Id == id {
			if isAccessable(f, g.Token.File, g.Pub) {
				return i, dm, false
			}
		}
	}
	if dm.side != nil {
		i, m, _ := dm.side.findGlobalById(id, f)
		return i, m, true
	}
	return -1, nil, false
}

func (dm *DefineMap) globalById(id string, f *File) (*Var, *DefineMap, bool) {
	i, m, canshadow := dm.findGlobalById(id, f)
	if i == -1 {
		return nil, nil, false
	}
	return m.Globals[i], m, canshadow
}

func (dm *DefineMap) findById(id string, f *File) (int, *DefineMap, byte) {
	var finders = map[byte]func(string, *jnio.File) (int, *DefineMap, bool){
		'g': dm.findGlobalById,
		'f': dm.findFuncById,
		'e': dm.findEnumById,
		's': dm.findStructById,
		't': dm.findTypeById,
		'i': dm.findTraitById,
	}
	for code, finder := range finders {
		i, m, _ := finder(id, f)
		if i != -1 {
			return i, m, code
		}
	}
	return -1, nil, ' '
}

func pushDefines(dest, src *DefineMap) {
	dest.Types = append(dest.Types, src.Types...)
	dest.Traits = append(dest.Traits, src.Traits...)
	dest.Structs = append(dest.Structs, src.Structs...)
	dest.Enums = append(dest.Enums, src.Enums...)
	dest.Globals = append(dest.Globals, src.Globals...)
	dest.Funcs = append(dest.Funcs, src.Funcs...)
}
