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
	"github.com/DeRuneLabs/jane/ast/models"
	"github.com/DeRuneLabs/jane/package/jn"
)

func (p *Parser) getFieldMap(f *Func) *paramMap {
	pmap := new(paramMap)
	*pmap = paramMap{}
	s := f.RetType.Type.Tag.(*jnstruct)
	for i, g := range s.Defs.Globals {
		if isAccessable(p.File, g.Token.File, g.Pub) {
			param := &f.Params[i]
			(*pmap)[param.Id] = &paramMapPair{param, nil}
		}
	}
	return pmap
}

type structArgParser struct {
	p      *Parser
	fmap   *paramMap
	f      *Func
	args   *models.Args
	i      int
	arg    Arg
	errTok Tok
}

func (sap *structArgParser) buildArgs() {
	sap.args.Src = make([]models.Arg, len(*sap.fmap))
	for i, p := range sap.f.Params {
		pair := (*sap.fmap)[p.Id]
		switch {
		case pair.arg != nil:
			sap.args.Src[i] = *pair.arg
		case paramHasDefaultArg(pair.param):
			arg := Arg{Expr: pair.param.Default}
			sap.args.Src[i] = arg
		case pair.param.Variadic:
			model := sliceExpr{pair.param.Type, nil}
			model.dataType.Kind = jn.Prefix_Slice + model.dataType.Kind
			arg := Arg{Expr: Expr{Model: model}}
			sap.args.Src[i] = arg
		}
	}
}

func (sap *structArgParser) pushArg() {
	defer func() { sap.i++ }()
	if sap.arg.TargetId == "" {
		sap.p.pusherrtok(sap.arg.Tok, "argument_must_target_to_parameter")
		return
	}
	pair, ok := (*sap.fmap)[sap.arg.TargetId]
	if !ok {
		sap.p.pusherrtok(sap.arg.Tok, "id_noexist", sap.arg.TargetId)
		return
	} else if pair.arg != nil {
		sap.p.pusherrtok(sap.arg.Tok, "already_has_expr", sap.arg.TargetId)
		return
	}
	arg := sap.arg
	pair.arg = &arg
	sap.p.parseArg(sap.f, pair, sap.args, nil)
}

func (sap *structArgParser) checkPasses() {
	for _, pair := range *sap.fmap {
		if pair.arg == nil &&
			!pair.param.Variadic &&
			!paramHasDefaultArg(pair.param) {
			sap.p.pusherrtok(sap.errTok, "missing_expr_for", pair.param.Id)
		}
	}
}

func (sap *structArgParser) parse() {
	sap.fmap = sap.p.getFieldMap(sap.f)
	argCount := 0
	for sap.i, sap.arg = range sap.args.Src {
		if sap.arg.TargetId != "" {
			break
		}
		if argCount >= len(sap.f.Params) {
			sap.p.pusherrtok(sap.errTok, "argument_overflow")
			return
		}
		argCount++
		param := &sap.f.Params[sap.i]
		arg := sap.arg
		(*sap.fmap)[param.Id].arg = &arg
		sap.p.parseArg(sap.f, (*sap.fmap)[param.Id], sap.args, nil)
	}
	for sap.i < len(sap.args.Src) {
		sap.arg = sap.args.Src[sap.i]
		sap.pushArg()
	}
	sap.checkPasses()
	sap.buildArgs()
}
