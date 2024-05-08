// Copyright (c) 2024 arfy slowy - DeRuneLabs
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

package ast

import "github.com/DeRuneLabs/jane/lexer"

func IsAccessable(finder *lexer.File, target *lexer.File, defIsPub bool) bool {
	return defIsPub || finder == nil || target == nil || finder.Dir() == target.Dir()
}

type Defmap struct {
	Namespaces []*Namespace
	Enums      []*Enum
	Structs    []*Struct
	Traits     []*Trait
	Types      []*TypeAlias
	Fns        []*Fn
	Globals    []*Var
	Side       *Defmap
}

func (dm *Defmap) FindNsById(id string) int {
	for i, t := range dm.Namespaces {
		if t != nil && t.Id == id {
			return i
		}
	}
	return -1
}

func (dm *Defmap) NsById(id string) *Namespace {
	i := dm.FindNsById(id)
	if i == -1 {
		return nil
	}
	return dm.Namespaces[i]
}

func (dm *Defmap) FindStructById(id string, f *lexer.File) (int, *Defmap, bool) {
	for i, s := range dm.Structs {
		if s != nil && s.Id == id {
			if IsAccessable(f, s.Token.File, s.Pub) {
				return i, dm, false
			}
		}
	}
	if dm.Side != nil {
		i, m, _ := dm.Side.FindStructById(id, f)
		return i, m, true
	}
	return -1, nil, false
}

func (dm *Defmap) StructById(id string, f *lexer.File) (*Struct, *Defmap, bool) {
	i, m, canshadow := dm.FindStructById(id, f)
	if i == -1 {
		return nil, nil, false
	}
	return m.Structs[i], m, canshadow
}

func (dm *Defmap) FindTraitById(id string, f *lexer.File) (int, *Defmap, bool) {
	for i, t := range dm.Traits {
		if t != nil && t.Id == id {
			if IsAccessable(f, t.Token.File, t.Pub) {
				return i, dm, false
			}
		}
	}
	if dm.Side != nil {
		i, m, _ := dm.Side.FindTraitById(id, f)
		return i, m, true
	}
	return -1, nil, false
}

func (dm *Defmap) TraitById(id string, f *lexer.File) (*Trait, *Defmap, bool) {
	i, m, canshadow := dm.FindTraitById(id, f)
	if i == -1 {
		return nil, nil, false
	}
	return m.Traits[i], m, canshadow
}

func (dm *Defmap) FindEnumById(id string, f *lexer.File) (int, *Defmap, bool) {
	for i, e := range dm.Enums {
		if e != nil && e.Id == id {
			if IsAccessable(f, e.Token.File, e.Pub) {
				return i, dm, false
			}
		}
	}
	if dm.Side != nil {
		i, m, _ := dm.Side.FindEnumById(id, f)
		return i, m, true
	}
	return -1, nil, false
}

func (dm *Defmap) EnumById(id string, f *lexer.File) (*Enum, *Defmap, bool) {
	i, m, canshadow := dm.FindEnumById(id, f)
	if i == -1 {
		return nil, nil, false
	}
	return m.Enums[i], m, canshadow
}

func (dm *Defmap) FindTypeById(id string, f *lexer.File) (int, *Defmap, bool) {
	for i, t := range dm.Types {
		if t != nil && t.Id == id {
			if IsAccessable(f, t.Token.File, t.Pub) {
				return i, dm, false
			}
		}
	}
	if dm.Side != nil {
		i, m, _ := dm.Side.FindTypeById(id, f)
		return i, m, true
	}
	return -1, nil, false
}

func (dm *Defmap) TypeById(id string, f *lexer.File) (*TypeAlias, *Defmap, bool) {
	i, m, canshadow := dm.FindTypeById(id, f)
	if i == -1 {
		return nil, nil, false
	}
	return m.Types[i], m, canshadow
}

func (dm *Defmap) FindFnById(id string, f *lexer.File) (int, *Defmap, bool) {
	for i, fn := range dm.Fns {
		if fn != nil && fn.Id == id {
			if IsAccessable(f, fn.Token.File, fn.Public) {
				return i, dm, false
			}
		}
	}
	if dm.Side != nil {
		i, m, _ := dm.Side.FindFnById(id, f)
		return i, m, true
	}
	return -1, nil, false
}

func (dm *Defmap) FnById(id string, f *lexer.File) (*Fn, *Defmap, bool) {
	i, m, canshadow := dm.FindFnById(id, f)
	if i == -1 {
		return nil, nil, false
	}
	return m.Fns[i], m, canshadow
}

func (dm *Defmap) FindGlobalById(id string, f *lexer.File) (int, *Defmap, bool) {
	for i, g := range dm.Globals {
		if g != nil && g.DataType.Id != void_t && g.Id == id {
			if IsAccessable(f, g.Token.File, g.Public) {
				return i, dm, false
			}
		}
	}
	if dm.Side != nil {
		i, m, _ := dm.Side.FindGlobalById(id, f)
		return i, m, true
	}
	return -1, nil, false
}

func (dm *Defmap) GlobalById(id string, f *lexer.File) (*Var, *Defmap, bool) {
	i, m, canshadow := dm.FindGlobalById(id, f)
	if i == -1 {
		return nil, nil, false
	}
	return m.Globals[i], m, canshadow
}

func (dm *Defmap) FindById(id string, f *lexer.File) (int, *Defmap, byte) {
	var finders = map[byte]func(string, *lexer.File) (int, *Defmap, bool){
		'g': dm.FindGlobalById,
		'f': dm.FindFnById,
		'e': dm.FindEnumById,
		's': dm.FindStructById,
		't': dm.FindTypeById,
		'i': dm.FindTraitById,
	}
	for code, finder := range finders {
		i, m, _ := finder(id, f)
		if i != -1 {
			return i, m, code
		}
	}
	return -1, nil, ' '
}

func (dm *Defmap) PushDefines(dest *Defmap) {
	dest.Types = append(dest.Types, dm.Types...)
	dest.Traits = append(dest.Traits, dm.Traits...)
	dest.Structs = append(dest.Structs, dm.Structs...)
	dest.Enums = append(dest.Enums, dm.Enums...)
	dest.Globals = append(dest.Globals, dm.Globals...)
	dest.Fns = append(dest.Fns, dm.Fns...)
}
