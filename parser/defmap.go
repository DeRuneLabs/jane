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

package parser

import (
	"github.com/DeRuneLabs/jane/package/jnio"
	"github.com/DeRuneLabs/jane/package/jntype"
)

func isAccessable(finder, target *File, defIsPub bool) bool {
	return defIsPub || finder == nil || finder.Dir == target.Dir
}

type Defmap struct {
	Namespaces []*namespace
	Enums      []*Enum
	Structs    []*jnstruct
	Traits     []*trait
	Types      []*Type
	Funcs      []*function
	Globals    []*Var
	side       *Defmap
}

func (dm *Defmap) findNsById(id string) int {
	for i, t := range dm.Namespaces {
		if t != nil && t.Id == id {
			return i
		}
	}
	return -1
}

func (dm *Defmap) nsById(id string) *namespace {
	i := dm.findNsById(id)
	if i == -1 {
		return nil
	}
	return dm.Namespaces[i]
}

func (dm *Defmap) findStructById(id string, f *File) (int, *Defmap, bool) {
	for i, s := range dm.Structs {
		if s != nil && s.Ast.Id == id {
			if isAccessable(f, s.Ast.Tok.File, s.Ast.Pub) {
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

func (dm *Defmap) structById(id string, f *File) (*jnstruct, *Defmap, bool) {
	i, m, canshadow := dm.findStructById(id, f)
	if i == -1 {
		return nil, nil, false
	}
	return m.Structs[i], m, canshadow
}

func (dm *Defmap) findTraitById(id string, f *File) (int, *Defmap, bool) {
	for i, t := range dm.Traits {
		if t != nil && t.Ast.Id == id {
			if isAccessable(f, t.Ast.Tok.File, t.Ast.Pub) {
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

func (dm *Defmap) traitById(id string, f *File) (*trait, *Defmap, bool) {
	i, m, canshadow := dm.findTraitById(id, f)
	if i == -1 {
		return nil, nil, false
	}
	return m.Traits[i], m, canshadow
}

func (dm *Defmap) findEnumById(id string, f *File) (int, *Defmap, bool) {
	for i, e := range dm.Enums {
		if e != nil && e.Id == id {
			if isAccessable(f, e.Tok.File, e.Pub) {
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

func (dm *Defmap) enumById(id string, f *File) (*Enum, *Defmap, bool) {
	i, m, canshadow := dm.findEnumById(id, f)
	if i == -1 {
		return nil, nil, false
	}
	return m.Enums[i], m, canshadow
}

func (dm *Defmap) findTypeById(id string, f *File) (int, *Defmap, bool) {
	for i, t := range dm.Types {
		if t != nil && t.Id == id {
			if isAccessable(f, t.Tok.File, t.Pub) {
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

func (dm *Defmap) typeById(id string, f *File) (*Type, *Defmap, bool) {
	i, m, canshadow := dm.findTypeById(id, f)
	if i == -1 {
		return nil, nil, false
	}
	return m.Types[i], m, canshadow
}

func (dm *Defmap) findFuncById(id string, f *File) (int, *Defmap, bool) {
	for i, fn := range dm.Funcs {
		if fn != nil && fn.Ast.Id == id {
			if isAccessable(f, fn.Ast.Tok.File, fn.Ast.Pub) {
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

func (dm *Defmap) funcById(id string, f *File) (*function, *Defmap, bool) {
	i, m, canshadow := dm.findFuncById(id, f)
	if i == -1 {
		return nil, nil, false
	}
	return m.Funcs[i], m, canshadow
}

func (dm *Defmap) findGlobalById(id string, f *File) (int, *Defmap, bool) {
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

func (dm *Defmap) globalById(id string, f *File) (*Var, *Defmap, bool) {
	i, m, canshadow := dm.findGlobalById(id, f)
	if i == -1 {
		return nil, nil, false
	}
	return m.Globals[i], m, canshadow
}

func (dm *Defmap) findById(id string, f *File) (int, *Defmap, byte) {
	finders := map[byte]func(string, *jnio.File) (int, *Defmap, bool){
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
