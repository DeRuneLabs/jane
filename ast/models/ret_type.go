package models

import "github.com/DeRuneLabs/jane/package/jnapi"

type RetType struct {
	Type        DataType
	Identifiers Toks
}

func (rt RetType) String() string {
	return rt.Type.String()
}

func (rt *RetType) AnyVar() bool {
	for _, tok := range rt.Identifiers {
		if !jnapi.IsIgnoreId(tok.Kind) {
			return true
		}
	}
	return false
}

func (rt *RetType) Vars() []*Var {
	if !rt.Type.MultiTyped {
		return nil
	}
	types := rt.Type.Tag.([]DataType)
	var vars []*Var
	for i, tok := range rt.Identifiers {
		if jnapi.IsIgnoreId(tok.Kind) {
			continue
		}
		variable := new(Var)
		variable.IdTok = tok
		variable.Id = tok.Kind
		variable.Type = types[i]
		vars = append(vars, variable)
	}
	return vars
}
