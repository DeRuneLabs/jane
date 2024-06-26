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

package parser

import (
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/DeRuneLabs/jane"
	"github.com/DeRuneLabs/jane/ast"
	"github.com/DeRuneLabs/jane/build"
	"github.com/DeRuneLabs/jane/lexer"
	"github.com/DeRuneLabs/jane/types"
)

type (
	File        = lexer.File
	TypeAlias   = ast.TypeAlias
	Var         = ast.Var
	Fn          = ast.Fn
	Arg         = ast.Arg
	Param       = ast.Param
	Type        = ast.Type
	Expr        = ast.Expr
	Enum        = ast.Enum
	Struct      = ast.Struct
	GenericType = ast.GenericType
	RetType     = ast.RetType
)

type Parser struct {
	attributes       []ast.Attribute
	doc_text         strings.Builder
	currentIter      *ast.Iter
	currentCase      *ast.Case
	wg               sync.WaitGroup
	rootBlock        *ast.Block
	nodeBlock        *ast.Block
	blockTypes       []*TypeAlias
	block_vars       []*Var
	waitingImpls     []*ast.Impl
	eval             *eval
	linked_aliases   []*ast.TypeAlias
	linked_functions []*ast.Fn
	linked_variables []*ast.Var
	linked_structs   []*ast.Struct
	allowBuiltin     bool
	package_files    *[]*Parser
	not_package      bool
	JustDefines      bool
	NoCheck          bool
	Used             *[]*ast.UseDecl
	Uses             []*ast.UseDecl
	Defines          *ast.Defmap
	Errors           []build.Log
	File             *File
}

func (p *Parser) parse_file() {
	lexer := lexer.New(p.File)
	toks := lexer.Lex()
	if lexer.Logs != nil {
		p.pusherrs(lexer.Logs...)
		return
	}

	tree, errors := get_tree(toks)
	if len(errors) > 0 {
		p.pusherrs(errors...)
		return
	}

	if !p.parseTree(tree) {
		return
	}

	p.checkParse()
	p.wg.Wait()
}

func ParseFile(path string, just_defines bool) (*Parser, string) {
	if !build.IsJane(path) {
		return nil, build.Errorf("file_not_jane", path)
	}

	p := new_parser(path)
	if !p.File.IsOk() {
		return nil, "path is not exist or inaccessible: " + path
	}

	p.not_package = true
	p.setup_package()
	p.parse_file()

	return p, ""
}

func (p *Parser) link_package(dirents []fs.DirEntry) {
	dir := filepath.Dir(p.File.Path())
	for _, info := range dirents {
		name := info.Name()
		if info.IsDir() ||
			!strings.HasSuffix(name, jane.EXT) ||
			!build.IsPassFileAnnotation(name) {
			continue
		}
		fp := new_parser(filepath.Join(dir, name))
		fp.Used = p.Used
		fp.package_files = p.package_files
		*p.package_files = append(*p.package_files, fp)
		fp.NoCheck = true
		fp.JustDefines = true

		fp.parse_file()
		fp.wg.Wait()
		if len(fp.Errors) > 0 {
			p.pusherrs(fp.Errors...)
			return
		}
	}
}

func ParsePackage(path string, just_defines bool) (*Parser, string) {
	dirents, err := os.ReadDir(path)
	if err != nil {
		return nil, err.Error()
	}
	for i, dirent := range dirents {
		name := dirent.Name()
		if dirent.IsDir() ||
			!strings.HasSuffix(name, jane.EXT) ||
			!build.IsPassFileAnnotation(name) {
			continue
		}
		p := new_parser(filepath.Join(path, name))
		p.setup_package()

		dirents = dirents[i+1:]
		p.link_package(dirents)
		p.parse_file()
		p.WrapPackage()

		return p, ""
	}
	return nil, "there is no jane source code"
}

func new_parser(path string) *Parser {
	p := new(Parser)
	p.File = lexer.NewFile(path)
	p.allowBuiltin = true
	p.Defines = new(ast.Defmap)
	p.eval = new(eval)
	p.eval.p = p
	p.Used = new([]*ast.UseDecl)
	return p
}

func (p *Parser) setup_package() {
	p.package_files = new([]*Parser)
	*p.package_files = append(*p.package_files, p)
}

func (p *Parser) pusherrtok(tok lexer.Token, key string, args ...any) {
	p.pusherrmsgtok(tok, build.Errorf(key, args...))
}

func (p *Parser) pusherrmsgtok(tok lexer.Token, msg string) {
	p.Errors = append(p.Errors, build.Log{
		Type:   build.ERR,
		Row:    tok.Row,
		Column: tok.Column,
		Path:   tok.File.Path(),
		Text:   msg,
	})
}

func (p *Parser) pusherrs(errs ...build.Log) {
	p.Errors = append(p.Errors, errs...)
}

func (p *Parser) PushErr(key string, args ...any) {
	p.pusherrmsg(build.Errorf(key, args...))
}

func (p *Parser) pusherrmsg(msg string) {
	p.Errors = append(p.Errors, build.Log{
		Type: build.FLAT_ERR,
		Text: msg,
	})
}

func get_tree(toks []lexer.Token) ([]ast.Node, []build.Log) {
	r := new_builder(toks)
	r.Build()
	return r.Tree, r.Errors
}

func (p *Parser) checkCppUsePath(use *ast.UseDecl) bool {
	if build.IsStdHeaderPath(use.Path) {
		return true
	}
	ext := filepath.Ext(use.Path)
	if !build.IsValidHeader(ext) {
		p.pusherrtok(use.Token, "invalid_header_ext", ext)
		return false
	}
	err := os.Chdir(use.Token.File.Dir())
	if err != nil {
		p.pusherrtok(use.Token, "use_not_found", use.Path)
		return false
	}
	info, err := os.Stat(use.Path)
	if err != nil || info.IsDir() {
		p.pusherrtok(use.Token, "use_not_found", use.Path)
		return false
	}
	use.Path, _ = filepath.Abs(use.Path)
	_ = os.Chdir(jane.WORKING_PATH)
	return true
}

func (p *Parser) checkPureUsePath(use *ast.UseDecl) bool {
	info, err := os.Stat(use.Path)
	if err != nil || !info.IsDir() {
		p.pusherrtok(use.Token, "use_not_found", use.Path)
		return false
	}
	return true
}

func (p *Parser) checkUsePath(use *ast.UseDecl) bool {
	if use.Cpp {
		if !p.checkCppUsePath(use) {
			return false
		}
	} else {
		if !p.checkPureUsePath(use) {
			return false
		}
	}
	return true
}

func (p *Parser) pushSelects(u *ast.UseDecl, selectors []lexer.Token) (addNs bool) {
	if len(selectors) > 0 && p.Defines.Side == nil {
		p.Defines.Side = new(ast.Defmap)
	}
	for i, id := range selectors {
		for j, jid := range selectors {
			if j >= i {
				break
			} else if jid.Kind == id.Kind {
				p.pusherrtok(id, "exist_id", id.Kind)
				i = -1
				break
			}
		}
		if i == -1 {
			break
		}
		if id.Id == lexer.ID_SELF {
			addNs = true
			continue
		}
		i, m, def_t := u.Defines.FindById(id.Kind, p.File)
		if i == -1 {
			p.pusherrtok(id, "id_not_exist", id.Kind)
			continue
		}
		switch def_t {
		case 'i':
			p.Defines.Side.Traits = append(p.Defines.Side.Traits, m.Traits[i])
		case 'f':
			p.Defines.Side.Fns = append(p.Defines.Side.Fns, m.Fns[i])
		case 'e':
			p.Defines.Side.Enums = append(p.Defines.Side.Enums, m.Enums[i])
		case 'g':
			p.Defines.Side.Globals = append(p.Defines.Side.Globals, m.Globals[i])
		case 't':
			p.Defines.Side.Types = append(p.Defines.Side.Types, m.Types[i])
		case 's':
			p.Defines.Side.Structs = append(p.Defines.Side.Structs, m.Structs[i])
		}
	}
	return
}

func (p *Parser) pushUse(use *ast.UseDecl, selectors []lexer.Token) {
	if use.FullUse {
		if p.Defines.Side == nil {
			p.Defines.Side = new(ast.Defmap)
		}
		use.Defines.PushDefines(p.Defines.Side)
	} else if len(selectors) > 0 {
		if !p.pushSelects(use, selectors) {
			return
		}
	} else if selectors != nil {
		return
	}
	identifiers := strings.SplitN(use.LinkString, lexer.KND_DBLCOLON, -1)
	src := p.pushNs(identifiers)
	src.Defines = use.Defines
}

func (p *Parser) compileCppLinkUse(ast *ast.UseDecl) *ast.UseDecl {
	ast.Cpp = true
	return ast
}

func make_use_from_ast(decl *ast.UseDecl) *ast.UseDecl {
	decl.Defines = new(ast.Defmap)
	return decl
}

func (p *Parser) WrapPackage() {
	for _, fp := range *p.package_files {
		if p == fp {
			continue
		}
		fp.Defines.PushDefines(p.Defines)
	}
}

func (p *Parser) compilePureUse(ast *ast.UseDecl) (_ *ast.UseDecl, hassErr bool) {
	dirents, err := os.ReadDir(ast.Path)
	if err != nil {
		p.pusherrmsg(err.Error())
		return nil, true
	}
	for i, dirent := range dirents {
		name := dirent.Name()
		if dirent.IsDir() ||
			!strings.HasSuffix(name, jane.EXT) ||
			!build.IsPassFileAnnotation(name) {
			continue
		}
		path := filepath.Join(ast.Path, name)
		psub := new_parser(path)
		psub.setup_package()
		psub.Used = p.Used
		dm, ok := std_builtin_defines[ast.LinkString]
		if ok {
			dm.PushDefines(psub.Defines)
		}
		dirents = dirents[i+1:]
		psub.link_package(dirents)
		if len(psub.Errors) > 0 {
			p.pusherrs(psub.Errors...)
			p.pusherrtok(ast.Token, "use_has_errors")
			return nil, true
		}

		psub.parse_file()
		psub.WrapPackage()

		u := make_use_from_ast(ast)
		psub.Defines.PushDefines(u.Defines)
		p.pusherrs(psub.Errors...)
		p.pushUse(u, ast.Selectors)

		return u, false
	}
	return nil, false
}

func (p *Parser) compileUse(ast *ast.UseDecl) (*ast.UseDecl, bool) {
	if ast.Cpp {
		return p.compileCppLinkUse(ast), false
	}
	return p.compilePureUse(ast)
}

func (p *Parser) use_decl(decl *ast.UseDecl, err *bool) {
	if !p.checkUsePath(decl) {
		*err = true
		return
	}
	for _, u := range *p.Used {
		if decl.Path == u.Path {
			old := u.FullUse
			u.FullUse = decl.FullUse
			p.pushUse(u, decl.Selectors)
			p.Uses = append(p.Uses, u)
			u.FullUse = old
			return
		}
	}
	var u *ast.UseDecl
	u, *err = p.compileUse(decl)
	if u == nil || *err {
		return
	}
	for _, pu := range p.Uses {
		if u.Path == pu.Path {
			p.pusherrtok(decl.Token, "already_uses")
			return
		}
	}
	*p.Used = append(*p.Used, u)
	p.Uses = append(p.Uses, u)
}

func (p *Parser) parseUses(tree *[]ast.Node) (err bool) {
	for i := range *tree {
		node := &(*tree)[i]
		switch node_t := node.Data.(type) {
		case ast.UseDecl:
			p.use_decl(&node_t, &err)
			if err {
				return
			}
			node.Data = nil
		case ast.Comment:
		default:
			return
		}
	}
	*tree = nil
	return
}

func node_is_ignored(node *ast.Node) bool {
	return node.Data == nil
}

func (p *Parser) parseSrcTreeNode(node ast.Node) {
	if node_is_ignored(&node) {
		return
	}
	switch node_t := node.Data.(type) {
	case ast.St:
		p.St(node_t)
	case TypeAlias:
		p.Type(node_t)
	case *Enum:
		p.Enum(node_t)
	case Struct:
		p.Struct(node_t)
	case ast.Trait:
		p.Trait(node_t)
	case ast.Impl:
		i := new(ast.Impl)
		*i = node_t
		p.waitingImpls = append(p.waitingImpls, i)
	case ast.CppLinkFn:
		p.LinkFn(node_t)
	case ast.CppLinkVar:
		p.LinkVar(node_t)
	case ast.CppLinkStruct:
		p.Link_struct(node_t)
	case ast.CppLinkAlias:
		p.Link_alias(node_t)
	case ast.Comment:
		p.Comment(node_t)
	case ast.UseDecl:
		p.pusherrtok(node.Token, "use_at_content")
	default:
		p.pusherrtok(node.Token, "invalid_syntax")
	}
}

func (p *Parser) parseSrcTree(tree []ast.Node) {
	for _, node := range tree {
		p.parseSrcTreeNode(node)
		p.checkDoc(node)
		p.checkAttribute(node)
	}
}

func (p *Parser) parseTree(tree []ast.Node) (ok bool) {
	if p.parseUses(&tree) {
		return false
	}
	p.parseSrcTree(tree)
	return true
}

func (p *Parser) checkParse() {
	if !p.NoCheck {
		p.check_package()
	}
}

func (p *Parser) checkDoc(node ast.Node) {
	if p.doc_text.Len() == 0 {
		return
	}
	switch node.Data.(type) {
	case ast.Comment, ast.Attribute, []GenericType:
		return
	}
	p.doc_text.Reset()
}

func (p *Parser) checkAttribute(node ast.Node) {
	if p.attributes == nil {
		return
	}
	switch node.Data.(type) {
	case ast.Attribute, ast.Comment, []GenericType:
		return
	}
	p.attributes = nil
}

func (p *Parser) check_generics(types []*GenericType) {
	for i, t := range types {
		if lexer.IsIgnoreId(t.Id) {
			p.pusherrtok(t.Token, "ignore_id")
			continue
		}
		for j, ct := range types {
			if j >= i {
				break
			} else if t.Id == ct.Id {
				p.pusherrtok(t.Token, "exist_id", t.Id)
				break
			}
		}
	}
}

func (p *Parser) make_type_alias(alias ast.TypeAlias) *ast.TypeAlias {
	a := new(ast.TypeAlias)
	*a = alias
	alias.Doc = p.doc_text.String()
	p.doc_text.Reset()
	return a
}

func (p *Parser) Type(alias TypeAlias) {
	if lexer.IsIgnoreId(alias.Id) {
		p.pusherrtok(alias.Token, "ignore_id")
		return
	}
	_, tok, canshadow := p.defined_by_id(alias.Id)
	if tok.Id != lexer.ID_NA && !canshadow {
		p.pusherrtok(alias.Token, "exist_id", alias.Id)
		return
	}
	p.Defines.Types = append(p.Defines.Types, p.make_type_alias(alias))
}

func (p *Parser) parse_enum_items_str(e *Enum) {
	for _, item := range e.Items {
		if lexer.IsIgnoreId(item.Id) {
			p.pusherrtok(item.Token, "ignore_id")
		} else {
			for _, checkItem := range e.Items {
				if item == checkItem {
					break
				}
				if item.Id == checkItem.Id {
					p.pusherrtok(item.Token, "exist_id", item.Id)
					break
				}
			}
		}
		if item.Expr.Tokens != nil {
			val, model := p.eval_expr(item.Expr, nil)
			if !val.constant && !p.eval.has_error {
				p.pusherrtok(item.Expr.Tokens[0], "expr_not_const")
			}
			item.ExprTag = val.expr
			item.Expr.Model = model
			assign_checker{
				p:         p,
				t:         e.DataType,
				v:         val,
				ignoreAny: true,
				errtok:    item.Token,
			}.check()
		} else {
			expr := value{constant: true, expr: item.Id}
			item.ExprTag = expr.expr
			item.Expr.Model = get_str_model(expr)
		}
		itemVar := new(Var)
		itemVar.Constant = true
		itemVar.ExprTag = item.ExprTag
		itemVar.Id = item.Id
		itemVar.DataType = e.DataType
		itemVar.Token = e.Token
		p.Defines.Globals = append(p.Defines.Globals, itemVar)
	}
}

func (p *Parser) parse_enum_items_integer(e *Enum) {
	max := types.MaxOfType(e.DataType.Id)
	for i, item := range e.Items {
		if max == 0 {
			p.pusherrtok(item.Token, "overflow_limits")
		} else {
			max--
		}
		if lexer.IsIgnoreId(item.Id) {
			p.pusherrtok(item.Token, "ignore_id")
		} else {
			for _, checkItem := range e.Items {
				if item == checkItem {
					break
				}
				if item.Id == checkItem.Id {
					p.pusherrtok(item.Token, "exist_id", item.Id)
					break
				}
			}
		}
		if item.Expr.Tokens != nil {
			val, model := p.eval_expr(item.Expr, nil)
			if !val.constant && !p.eval.has_error {
				p.pusherrtok(item.Expr.Tokens[0], "expr_not_const")
			}
			item.ExprTag = val.expr
			item.Expr.Model = model
			assign_checker{
				p:         p,
				t:         e.DataType,
				v:         val,
				ignoreAny: true,
				errtok:    item.Token,
			}.check()
		} else {
			expr := max - (max - uint64(i))
			item.ExprTag = uint64(expr)
			item.Expr.Model = exprNode{strconv.FormatUint(expr, 16)}
		}
		itemVar := new(Var)
		itemVar.Constant = true
		itemVar.ExprTag = item.ExprTag
		itemVar.Id = item.Id
		itemVar.DataType = e.DataType
		itemVar.Token = e.Token
		p.Defines.Globals = append(p.Defines.Globals, itemVar)
	}
}

func (p *Parser) Enum(e *Enum) {
	if lexer.IsIgnoreId(e.Id) {
		p.pusherrtok(e.Token, "ignore_id")
		return
	}
	_, tok, _ := p.defined_by_id(e.Id)
	if tok.Id != lexer.ID_NA {
		p.pusherrtok(e.Token, "exist_id", e.Id)
		return
	}
	e.Doc = p.doc_text.String()
	p.doc_text.Reset()
	e.DataType, _ = p.realType(e.DataType, true)
	if !types.IsPure(e.DataType) {
		p.pusherrtok(e.Token, "invalid_type_source")
		return
	}
	p.Defines.Enums = append(p.Defines.Enums, e)
	if len(e.Items) == 0 {
		p.pusherrtok(e.Token, "enum_have_not_field", e.Id)
		return
	}
	pdefs := p.Defines
	puses := p.Uses
	p.Defines = new(ast.Defmap)
	defer func() {
		p.Defines = pdefs
		p.Uses = puses
	}()
	switch {
	case e.DataType.Id == types.STR:
		p.parse_enum_items_str(e)
	case types.IsInteger(e.DataType.Id):
		p.parse_enum_items_integer(e)
	default:
		p.pusherrtok(e.Token, "invalid_type_source")
	}
}

func (p *Parser) pushField(s *Struct, f **Var, i int) {
	for _, cf := range s.Fields {
		if *f == cf {
			break
		}
		if (*f).Id == cf.Id {
			p.pusherrtok((*f).Token, "exist_id", (*f).Id)
			break
		}
	}
	if len(s.Generics) == 0 {
		p.parse_field(s, f, i)
	} else {
		p.parseNonGenericType(s.Generics, &(*f).DataType)
		param := ast.Param{Id: (*f).Id, DataType: (*f).DataType}
		param.Default.Model = exprNode{build.CPP_DEFAULT_EXPR}
		s.Constructor.Params[i] = param
	}
}

func (p *Parser) parseFields(s *Struct) {
	s.Defines.Globals = make([]*ast.Var, len(s.Fields))
	for i, f := range s.Fields {
		p.pushField(s, &f, i)
		s.Defines.Globals[i] = f
	}
}

func make_constructor(s *Struct) *ast.Fn {
	constructor := new(ast.Fn)
	constructor.Id = s.Id
	constructor.Token = s.Token
	constructor.Params = make([]ast.Param, len(s.Fields))
	constructor.RetType.DataType = Type{
		Id:    types.STRUCT,
		Kind:  s.Id,
		Token: s.Token,
		Tag:   s,
	}
	if len(s.Generics) > 0 {
		constructor.Generics = make([]*ast.GenericType, len(s.Generics))
		copy(constructor.Generics, s.Generics)
		constructor.Combines = new([][]ast.Type)
	}
	return constructor
}

func (p *Parser) make_struct(model ast.Struct) *Struct {
	s := new(Struct)
	*s = model
	s.Doc = p.doc_text.String()
	p.doc_text.Reset()
	s.Owner = p
	s.Attributes = p.attributes
	p.attributes = nil
	s.Defines = new(ast.Defmap)
	s.Constructor = make_constructor(s)
	s.Origin = s
	p.check_generics(s.Generics)
	return s
}

func (p *Parser) Struct(model Struct) {
	if lexer.IsIgnoreId(model.Id) {
		p.pusherrtok(model.Token, "ignore_id")
		return
	} else if def, _, _ := p.defined_by_id(model.Id); def != nil {
		p.pusherrtok(model.Token, "exist_id", model.Id)
		return
	}
	s := p.make_struct(model)
	p.Defines.Structs = append(p.Defines.Structs, s)
}

func (p *Parser) LinkFn(link ast.CppLinkFn) {
	if lexer.IsIgnoreId(link.Link.Id) {
		p.pusherrtok(link.Token, "ignore_id")
		return
	}
	_, def_t := p.linkById(link.Link.Id)
	if def_t != ' ' {
		p.pusherrtok(link.Token, "exist_id", link.Link.Id)
		return
	}
	linkf := link.Link
	linkf.Owner = p
	linkf.Attributes = p.attributes
	p.attributes = nil
	p.check_generics(linkf.Generics)
	p.linked_functions = append(p.linked_functions, linkf)
}

func (p *Parser) Link_alias(link ast.CppLinkAlias) {
	if lexer.IsIgnoreId(link.Link.Id) {
		p.pusherrtok(link.Token, "ignore_id")
		return
	}
	_, def_t := p.linkById(link.Link.Id)
	if def_t != ' ' {
		p.pusherrtok(link.Token, "exist_id", link.Link.Id)
		return
	}
	ta := p.make_type_alias(link.Link)
	p.linked_aliases = append(p.linked_aliases, ta)
}

func (p *Parser) Link_struct(link ast.CppLinkStruct) {
	if lexer.IsIgnoreId(link.Link.Id) {
		p.pusherrtok(link.Token, "ignore_id")
		return
	}
	_, def_t := p.linkById(link.Link.Id)
	if def_t != ' ' {
		p.pusherrtok(link.Token, "exist_id", link.Link.Id)
		return
	}
	s := p.make_struct(link.Link)
	s.CppLinked = true
	p.linked_structs = append(p.linked_structs, s)
}

func (p *Parser) LinkVar(link ast.CppLinkVar) {
	if lexer.IsIgnoreId(link.Link.Id) {
		p.pusherrtok(link.Token, "ignore_id")
		return
	}
	_, def_t := p.linkById(link.Link.Id)
	if def_t != ' ' {
		p.pusherrtok(link.Token, "exist_id", link.Link.Id)
		return
	}
	p.linked_variables = append(p.linked_variables, link.Link)
}

func (p *Parser) Trait(model ast.Trait) {
	if lexer.IsIgnoreId(model.Id) {
		p.pusherrtok(model.Token, "ignore_id")
		return
	} else if def, _, _ := p.defined_by_id(model.Id); def != nil {
		p.pusherrtok(model.Token, "exist_id", model.Id)
		return
	}
	trait := new(ast.Trait)
	*trait = model
	trait.Desc = p.doc_text.String()
	p.doc_text.Reset()
	trait.Defines = new(ast.Defmap)
	trait.Defines.Fns = make([]*Fn, len(model.Funcs))
	for i, f := range trait.Funcs {
		if lexer.IsIgnoreId(f.Id) {
			p.pusherrtok(f.Token, "ignore_id")
		}
		for j, jf := range trait.Funcs {
			if j >= i {
				break
			} else if f.Id == jf.Id {
				p.pusherrtok(f.Token, "exist_id", f.Id)
			}
		}
		_ = p.check_param_dup(f.Params)
		p.parseTypesNonGenerics(f)
		trait.Defines.Fns[i] = f
	}
	p.Defines.Traits = append(p.Defines.Traits, trait)
}

func (p *Parser) implTrait(model *ast.Impl) {
	trait_def, _, _ := p.trait_by_id(model.Base.Kind)
	if trait_def == nil {
		p.pusherrtok(model.Base, "id_not_exist", model.Base.Kind)
		return
	}
	trait_def.Used = true
	sid, _ := model.Target.KindId()
	side := p.Defines.Side
	p.Defines.Side = nil
	s, _, _ := p.struct_by_id(model.Target.Kind)
	p.Defines.Side = side
	if s == nil {
		p.pusherrtok(model.Target.Token, "id_not_exist", sid)
		return
	}
	model.Target.Tag = s
	s.Origin.Traits = append(s.Origin.Traits, trait_def)
	for _, node := range model.Tree {
		switch node_t := node.Data.(type) {
		case ast.Comment:
			p.Comment(node_t)
		case *Fn:
			if trait_def.FindFunc(node_t.Id) == nil {
				p.pusherrtok(model.Target.Token, "trait_hasnt_id", trait_def.Id, node_t.Id)
				break
			}
			i, _, _ := s.Defines.FindById(node_t.Id, nil)
			if i != -1 {
				p.pusherrtok(node_t.Token, "exist_id", node_t.Id)
				continue
			}
			node_t.Receiver.Token = s.Token
			node_t.Receiver.Tag = s
			node_t.Attributes = p.attributes
			node_t.Owner = p
			p.attributes = nil
			node_t.Doc = p.doc_text.String()
			p.doc_text.Reset()
			_ = p.check_param_dup(node_t.Params)
			p.check_ret_variables(node_t)
			node_t.Used = true
			if len(s.Generics) == 0 {
				p.parseTypesNonGenerics(node_t)
			}
			s.Defines.Fns = append(s.Defines.Fns, node_t)
		}
	}
	for _, tf := range trait_def.Defines.Fns {
		ok := false
		ds := tf.DefineString()
		sf, _, _ := s.Defines.FnById(tf.Id, nil)
		if sf != nil {
			ok = tf.Public == sf.Public && ds == sf.DefineString()
		}
		if !ok {
			p.pusherrtok(model.Target.Token, "not_impl_trait_def", trait_def.Id, ds)
		}
	}
}

func (p *Parser) check_impl_generics(s *ast.Struct, types []*GenericType) {
	if len(s.Generics) > 0 {
		for _, t := range types {
			for _, g := range s.Generics {
				if t.Id == g.Id {
					p.pusherrtok(t.Token, "exist_id", t.Id)
				}
			}
		}
	}
	p.check_generics(types)
}

func (p *Parser) implStruct(model *ast.Impl) {
	side := p.Defines.Side
	p.Defines.Side = nil
	s, _, _ := p.struct_by_id(model.Base.Kind)
	p.Defines.Side = side
	if s == nil {
		p.pusherrtok(model.Base, "id_not_exist", model.Base.Kind)
		return
	}
	for _, node := range model.Tree {
		switch node_t := node.Data.(type) {
		case ast.Comment:
			p.Comment(node_t)
		case *Fn:
			i, _, _ := s.Defines.FindById(node_t.Id, nil)
			if i != -1 {
				p.pusherrtok(node_t.Token, "exist_id", node_t.Id)
				continue
			}
			sf := new(Fn)
			*sf = *node_t
			sf.Receiver.Token = s.Token
			sf.Receiver.Tag = s
			sf.Attributes = p.attributes
			sf.Doc = p.doc_text.String()
			sf.Owner = p
			p.doc_text.Reset()
			p.attributes = nil
			_ = p.check_param_dup(sf.Params)
			p.check_ret_variables(sf)
			p.check_impl_generics(s, sf.Generics)
			for _, generic := range node_t.Generics {
				if types.FindGeneric(generic.Id, s.Generics) != nil {
					p.pusherrtok(generic.Token, "exist_id", generic.Id)
				}
			}
			if len(s.Generics) == 0 {
				p.parseTypesNonGenerics(sf)
			}
			s.Defines.Fns = append(s.Defines.Fns, sf)
		}
	}
}

func (p *Parser) Impl(impl *ast.Impl) {
	if !types.IsVoid(impl.Target) {
		p.implTrait(impl)
		return
	}
	p.implStruct(impl)
}

func (p *Parser) pushNs(identifiers []string) *ast.Namespace {
	var src *ast.Namespace
	prev := p.Defines
	for _, id := range identifiers {
		src = prev.NsById(id)
		if src == nil {
			src = new(ast.Namespace)
			src.Id = id
			src.Defines = new(ast.Defmap)
			prev.Namespaces = append(prev.Namespaces, src)
		}
		prev = src.Defines
	}
	return src
}

func (p *Parser) Comment(c ast.Comment) {
	if strings.HasPrefix(c.Content, lexer.PRAGMA_COMMENT_PREFIX) {
		p.PushAttribute(c)
		return
	}
	p.doc_text.WriteString(c.Content)
	p.doc_text.WriteByte('\n')
}

func (p *Parser) PushAttribute(c ast.Comment) {
	var attr ast.Attribute
	attr.Tag = c.Content[len(lexer.PRAGMA_COMMENT_PREFIX):]
	attr.Token = c.Token
	ok := false
	for _, kind := range build.ATTRS {
		if attr.Tag == kind {
			ok = true
			break
		}
	}
	if !ok {
		return
	}
	for _, attr2 := range p.attributes {
		if attr.Tag == attr2.Tag {
			return
		}
	}
	p.attributes = append(p.attributes, attr)
}

func (p *Parser) St(s ast.St) {
	switch data_t := s.Data.(type) {
	case Fn:
		p.function(data_t)
	case Var:
		p.global(data_t)
	default:
		p.pusherrtok(s.Token, "invalid_syntax")
	}
}

func (p *Parser) parseFnNonGenericType(generics []*GenericType, dt *Type) {
	f := dt.Tag.(*Fn)
	for i := range f.Params {
		p.parseNonGenericType(generics, &f.Params[i].DataType)
	}
	p.parseNonGenericType(generics, &f.RetType.DataType)
}

func (p *Parser) parseMultiNonGenericType(generics []*GenericType, dt *Type) {
	types := dt.Tag.([]Type)
	for i := range types {
		mt := &types[i]
		p.parseNonGenericType(generics, mt)
	}
}

func (p *Parser) parseMapNonGenericType(generics []*GenericType, dt *Type) {
	p.parseMultiNonGenericType(generics, dt)
}

func (p *Parser) parseCommonNonGenericType(generics []*GenericType, dt *Type) {
	if dt.Id == types.ID {
		id, prefix := dt.KindId()
		def, _, _ := p.defined_by_id(id)
		switch deft := def.(type) {
		case *Struct:
			deft = p.structConstructorInstance(deft)
			if dt.Tag != nil {
				deft.SetGenerics(dt.Tag.([]Type))
			}
			dt.Kind = prefix + deft.AsTypeKind()
			dt.Id = types.STRUCT
			dt.Tag = deft
			dt.Pure = true
			dt.Original = nil
			goto tagcheck
		}
	}
	if types.IsGeneric(generics, *dt) {
		return
	}
tagcheck:
	if dt.Tag != nil {
		switch t := dt.Tag.(type) {
		case *Struct:
			for _, ct := range t.GetGenerics() {
				if types.IsGeneric(generics, ct) {
					return
				}
			}
		case []Type:
			for _, ct := range t {
				if types.IsGeneric(generics, ct) {
					return
				}
			}
		}
	}
	p.fn_parse_type(dt)
}

func (p *Parser) parseNonGenericType(generics []*GenericType, dt *Type) {
	switch {
	case dt.MultiTyped:
		p.parseMultiNonGenericType(generics, dt)
	case types.IsFn(*dt):
		p.parseFnNonGenericType(generics, dt)
	case types.IsMap(*dt):
		p.parseMapNonGenericType(generics, dt)
	case types.IsArray(*dt):
		p.parseNonGenericType(generics, dt.ComponentType)
		dt.Kind = lexer.PREFIX_ARRAY + dt.ComponentType.Kind
	case types.IsSlice(*dt):
		p.parseNonGenericType(generics, dt.ComponentType)
		dt.Kind = lexer.PREFIX_SLICE + dt.ComponentType.Kind
	default:
		p.parseCommonNonGenericType(generics, dt)
	}
}

func (p *Parser) parseTypesNonGenerics(f *Fn) {
	for i := range f.Params {
		p.parseNonGenericType(f.Generics, &f.Params[i].DataType)
	}
	p.parseNonGenericType(f.Generics, &f.RetType.DataType)
}

func (p *Parser) check_ret_variables(f *Fn) {
	for i, v := range f.RetType.Identifiers {
		if lexer.IsIgnoreId(v.Kind) {
			continue
		}
		for _, generic := range f.Generics {
			if v.Kind == generic.Id {
				goto exist
			}
		}
		for _, param := range f.Params {
			if v.Kind == param.Id {
				goto exist
			}
		}
		for j, jv := range f.RetType.Identifiers {
			if j >= i {
				break
			}
			if jv.Kind == v.Kind {
				goto exist
			}
		}
		continue
	exist:
		p.pusherrtok(v, "exist_id", v.Kind)

	}
}

func (p *Parser) function(ast Fn) {
	_, tok, canshadow := p.defined_by_id(ast.Id)
	if tok.Id != lexer.ID_NA && !canshadow {
		p.pusherrtok(ast.Token, "exist_id", ast.Id)
	} else if lexer.IsIgnoreId(ast.Id) {
		p.pusherrtok(ast.Token, "ignore_id")
	}
	f := new(Fn)
	*f = ast
	f.Attributes = p.attributes
	p.attributes = nil
	f.Owner = p
	f.Doc = p.doc_text.String()
	p.doc_text.Reset()
	p.check_generics(ast.Generics)
	p.check_ret_variables(f)
	_ = p.check_param_dup(f.Params)
	f.Used = f.Id == jane.INIT_FN
	p.Defines.Fns = append(p.Defines.Fns, f)
}

func (p *Parser) global(vast Var) {
	def, _, _ := p.defined_by_id(vast.Id)
	if def != nil {
		p.pusherrtok(vast.Token, "exist_id", vast.Id)
		return
	} else {
		for _, g := range p.Defines.Globals {
			if vast.Id == g.Id {
				p.pusherrtok(vast.Token, "exist_id", vast.Id)
				return
			}
		}
	}
	vast.Doc = p.doc_text.String()
	p.doc_text.Reset()
	v := new(Var)
	*v = vast
	p.Defines.Globals = append(p.Defines.Globals, v)
}

func (p *Parser) variable(model Var) *Var {
	if lexer.IsIgnoreId(model.Id) {
		p.pusherrtok(model.Token, "ignore_id")
	}
	v := new(Var)
	*v = model
	if v.DataType.Id != types.VOID {
		vt, ok := p.realType(v.DataType, true)
		if ok {
			v.DataType = vt
		} else {
			v.DataType = ast.Type{}
		}
	}
	var val value
	switch tag_t := v.Tag.(type) {
	case value:
		val = tag_t
	default:
		if v.SetterTok.Id != lexer.ID_NA {
			val, v.Expr.Model = p.eval_expr(v.Expr, &v.DataType)
		}
	}
	if val.data.DataType.MultiTyped {
		p.pusherrtok(model.Token, "missing_multi_assign_identifiers")
		return v
	}
	if v.DataType.Id != types.VOID {
		if v.SetterTok.Id != lexer.ID_NA {
			if v.DataType.Size.AutoSized && v.DataType.Id == types.ARRAY {
				v.DataType.Size = val.data.DataType.Size
			}
			assign_checker{
				p:                p,
				t:                v.DataType,
				v:                val,
				errtok:           v.Token,
				not_allow_assign: types.IsRef(v.DataType),
			}.check()
		}
	} else {
		if v.SetterTok.Id == lexer.ID_NA {
			p.pusherrtok(v.Token, "missing_autotype_value")
		} else {
			p.eval.has_error = p.eval.has_error || val.data.Value == ""
			v.DataType = val.data.DataType
			p.check_valid_init_expr(v.Mutable, val, v.SetterTok)
			p.checkValidityForAutoType(v.DataType, v.SetterTok)
		}
	}
	if !v.IsField && v.SetterTok.Id == lexer.ID_NA {
		p.pusherrtok(v.Token, "variable_not_initialized")
	}
	if v.Constant {
		v.ExprTag = val.expr
		if !types.IsAllowForConst(v.DataType) {
			p.pusherrtok(v.Token, "invalid_type_for_const", v.DataType.Kind)
		} else if v.SetterTok.Id != lexer.ID_NA && !is_valid_for_const(val) {
			p.eval.push_err_tok(v.Token, "expr_not_const")
		}
	}
	return v
}

func (p *Parser) varsFromParams(f *Fn) []*Var {
	length := len(f.Params)
	vars := make([]*Var, length)
	for i, param := range f.Params {
		v := new(ast.Var)
		v.Owner = f.Block
		v.Mutable = param.Mutable
		v.Id = param.Id
		v.Token = param.Token
		v.DataType = param.DataType
		if param.Variadic {
			if length-i > 1 {
				p.pusherrtok(param.Token, "variadic_parameter_not_last")
			}
			v.DataType = types.ToSlice(param.DataType)
		}
		vars[i] = v
	}
	return vars
}

func (p *Parser) linked_alias_by_id(id string) *ast.TypeAlias {
	for _, fp := range *p.package_files {
		for _, link := range fp.linked_aliases {
			if link.Id == id {
				return link
			}
		}
	}
	return nil
}

func (p *Parser) linked_struct_by_id(id string) *Struct {
	for _, fp := range *p.package_files {
		for _, link := range fp.linked_structs {
			if link.Id == id {
				return link
			}
		}
	}
	return nil
}

func (p *Parser) linkedVarById(id string) *Var {
	for _, fp := range *p.package_files {
		for _, link := range fp.linked_variables {
			if link.Id == id {
				return link
			}
		}
	}
	return nil
}

func (p *Parser) linkedFnById(id string) *ast.Fn {
	for _, fp := range *p.package_files {
		for _, link := range fp.linked_functions {
			if link.Id == id {
				return link
			}
		}
	}
	return nil
}

func (p *Parser) linkById(id string) (any, byte) {
	f := p.linkedFnById(id)
	if f != nil {
		return f, 'f'
	}
	v := p.linkedVarById(id)
	if v != nil {
		return v, 'v'
	}
	s := p.linked_struct_by_id(id)
	if s != nil {
		return s, 's'
	}
	ta := p.linked_alias_by_id(id)
	if ta != nil {
		return ta, 't'
	}
	return nil, ' '
}

func (p *Parser) fn_by_id(id string) (*Fn, *ast.Defmap, bool) {
	if p.allowBuiltin {
		f, _, _ := Builtin.FnById(id, nil)
		if f != nil {
			return f, nil, false
		}
	}
	for _, fp := range *p.package_files {
		f, dm, can_shadow := fp.Defines.FnById(id, fp.File)
		if f != nil && p.is_accessible_define(fp, dm) {
			return f, dm, can_shadow
		}
	}
	return nil, nil, false
}

func (p *Parser) global_by_id(id string) (*Var, *ast.Defmap, bool) {
	for _, fp := range *p.package_files {
		g, dm, _ := fp.Defines.GlobalById(id, fp.File)
		if g != nil && p.is_accessible_define(fp, dm) {
			return g, dm, true
		}
	}
	return nil, nil, false
}

func (p *Parser) NsById(id string) *ast.Namespace {
	return p.Defines.NsById(id)
}

func (p *Parser) is_shadowed(id string) bool {
	def, _, _ := p.block_define_by_id(id)
	return def != nil
}

func (p *Parser) is_accessible_define(fp *Parser, dm *ast.Defmap) bool {
	return p == fp || dm == fp.Defines
}

func (p *Parser) type_by_id(id string) (*TypeAlias, *ast.Defmap, bool) {
	alias, canshadow := p.block_type_by_id(id)
	if alias != nil {
		return alias, nil, canshadow
	}
	if p.allowBuiltin {
		alias, _, _ = Builtin.TypeById(id, nil)
		if alias != nil {
			return alias, nil, false
		}
	}
	for _, fp := range *p.package_files {
		a, dm, can_shadow := fp.Defines.TypeById(id, fp.File)
		if a != nil && p.is_accessible_define(fp, dm) {
			return a, dm, can_shadow
		}
	}
	return nil, nil, false
}

func (p *Parser) enum_by_id(id string) (*Enum, *ast.Defmap, bool) {
	if p.allowBuiltin {
		e, _, _ := Builtin.EnumById(id, nil)
		if e != nil {
			return e, nil, false
		}
	}
	for _, fp := range *p.package_files {
		e, dm, can_shadow := fp.Defines.EnumById(id, fp.File)
		if e != nil && p.is_accessible_define(fp, dm) {
			return e, dm, can_shadow
		}
	}
	return nil, nil, false
}

func (p *Parser) struct_by_id(id string) (*Struct, *ast.Defmap, bool) {
	if p.allowBuiltin {
		s, _, _ := Builtin.StructById(id, nil)
		if s != nil {
			return s, nil, false
		}
	}
	for _, fp := range *p.package_files {
		s, dm, can_shadow := fp.Defines.StructById(id, fp.File)
		if s != nil && p.is_accessible_define(fp, dm) {
			return s, dm, can_shadow
		}
	}
	return nil, nil, false
}

func (p *Parser) trait_by_id(id string) (*ast.Trait, *ast.Defmap, bool) {
	if p.allowBuiltin {
		trait_def, _, _ := Builtin.TraitById(id, nil)
		if trait_def != nil {
			return trait_def, nil, false
		}
	}
	for _, fp := range *p.package_files {
		t, dm, can_shadow := fp.Defines.TraitById(id, fp.File)
		if t != nil && p.is_accessible_define(fp, dm) {
			return t, dm, can_shadow
		}
	}
	return nil, nil, false
}

func (p *Parser) block_type_by_id(id string) (_ *TypeAlias, can_shadow bool) {
	for i := len(p.blockTypes) - 1; i >= 0; i-- {
		alias := p.blockTypes[i]
		if alias != nil && alias.Id == id {
			return alias, !alias.Generic && alias.Owner != p.nodeBlock
		}
	}
	return nil, false
}

func (p *Parser) block_var_by_id(id string) (_ *Var, can_shadow bool) {
	for i := len(p.block_vars) - 1; i >= 0; i-- {
		v := p.block_vars[i]
		if v != nil && v.Id == id {
			return v, v.Owner != p.nodeBlock
		}
	}
	return nil, false
}

func (p *Parser) defined_by_id(id string) (def any, tok lexer.Token, canshadow bool) {
	var a *TypeAlias
	a, _, canshadow = p.type_by_id(id)
	if a != nil {
		return a, a.Token, canshadow
	}
	var e *Enum
	e, _, canshadow = p.enum_by_id(id)
	if e != nil {
		return e, e.Token, canshadow
	}
	var s *Struct
	s, _, canshadow = p.struct_by_id(id)
	if s != nil {
		return s, s.Token, canshadow
	}
	var trait *ast.Trait
	trait, _, canshadow = p.trait_by_id(id)
	if trait != nil {
		return trait, trait.Token, canshadow
	}
	var f *Fn
	f, _, canshadow = p.fn_by_id(id)
	if f != nil {
		return f, f.Token, canshadow
	}
	bv, canshadow := p.block_var_by_id(id)
	if bv != nil {
		return bv, bv.Token, canshadow
	}
	g, _, _ := p.global_by_id(id)
	if g != nil {
		return g, g.Token, true
	}
	return
}

func (p *Parser) block_define_by_id(id string) (def any, tok lexer.Token, canshadow bool) {
	bv, canshadow := p.block_var_by_id(id)
	if bv != nil {
		return bv, bv.Token, canshadow
	}
	alias, canshadow := p.block_type_by_id(id)
	if alias != nil {
		return alias, alias.Token, canshadow
	}
	return
}

func (p *Parser) precheck_package() {
	p.parse_package_aliases()
	p.parse_package_structs()
	p.parse_package_waiting_fns()
	p.parse_package_waiting_impls()
	p.parse_package_waiting_globals()
	p.check_package_cpp_links()
}

func (p *Parser) parse_package_defines() {
	for _, pf := range *p.package_files {
		pf.parse_defines()
		if p != pf {
			pf.wg.Wait()
			p.pusherrs(pf.Errors...)
		}
	}
}

func (p *Parser) parse_defines() {
	p.check_structs()
	p.check_fns()
}

func (p *Parser) check_package() {
	p.precheck_package()
	if !p.JustDefines {
		p.parse_package_defines()
	}
}

func (p *Parser) parse_struct(s *Struct) {
	p.parseFields(s)
}

func (p *Parser) parse_package_structs() {
	for _, pf := range *p.package_files {
		pf.parse_structs()
		if p != pf {
			pf.wg.Wait()
			p.pusherrs(pf.Errors...)
		}
	}
}

func (p *Parser) parse_structs() {
	types.OrderStructures(p.Defines.Structs)

	for _, s := range p.Defines.Structs {
		p.parse_struct(s)
	}
}

func (p *Parser) parse_package_linked_structs() {
	for _, pf := range *p.package_files {
		pf.parse_linked_structs()
		if p != pf {
			pf.wg.Wait()
			p.pusherrs(pf.Errors...)
		}
	}
}

func (p *Parser) check_package_linked_aliases() {
	for _, pf := range *p.package_files {
		pf.check_linked_aliases()
		if p != pf {
			pf.wg.Wait()
			p.pusherrs(pf.Errors...)
		}
	}
}

func (p *Parser) check_package_linked_vars() {
	for _, pf := range *p.package_files {
		pf.check_linked_vars()
		if p != pf {
			pf.wg.Wait()
			p.pusherrs(pf.Errors...)
		}
	}
}

func (p *Parser) check_package_linked_fns() {
	for _, pf := range *p.package_files {
		pf.check_linked_fns()
		if p != pf {
			pf.wg.Wait()
			p.pusherrs(pf.Errors...)
		}
	}
}

func (p *Parser) parse_linked_structs() {
	for _, link := range p.linked_structs {
		p.parse_struct(link)
	}
}

func (p *Parser) check_linked_aliases() {
	for _, link := range p.linked_aliases {
		link.TargetType, _ = p.realType(link.TargetType, true)
	}
}

func (p *Parser) check_linked_vars() {
	for _, link := range p.linked_variables {
		vt, ok := p.realType(link.DataType, true)
		if ok {
			link.DataType = vt
		}
	}
}

func (p *Parser) check_linked_fns() {
	for _, link := range p.linked_functions {
		if len(link.Generics) == 0 {
			p.reload_fn_types(link)
		}
	}
}

func (p *Parser) parse_aliases() {
	for i, alias := range p.Defines.Types {
		p.Defines.Types[i].TargetType, _ = p.realType(alias.TargetType, true)
	}
}

func (p *Parser) check_package_cpp_links() {
	p.check_package_linked_aliases()
	p.parse_package_linked_structs()
	p.check_package_linked_vars()
	p.check_package_linked_fns()
}

func (p *Parser) parse_package_aliases() {
	for _, pf := range *p.package_files {
		pf.parse_aliases()
		if p != pf {
			pf.wg.Wait()
			p.pusherrs(pf.Errors...)
		}
	}
}

func (p *Parser) parse_package_waiting_fns() {
	for _, pf := range *p.package_files {
		pf.ParseWaitingFns()
		if p != pf {
			pf.wg.Wait()
			p.pusherrs(pf.Errors...)
		}
	}
}

func (p *Parser) ParseWaitingFns() {
	for _, f := range p.Defines.Fns {
		owner := p
		if len(f.Generics) > 0 {
			owner.parseTypesNonGenerics(f)
		} else {
			owner.reload_fn_types(f)
		}
	}
}

func (p *Parser) parse_package_waiting_globals() {
	for _, pf := range *p.package_files {
		pf.ParseWaitingGlobals()
		if p != pf {
			pf.wg.Wait()
			p.pusherrs(pf.Errors...)
		}
	}
}

func (p *Parser) ParseWaitingGlobals() {
	for _, g := range p.Defines.Globals {
		*g = *p.variable(*g)
	}
}

func (p *Parser) parse_package_waiting_impls() {
	for _, pf := range *p.package_files {
		pf.ParseWaitingImpls()
		if p != pf {
			pf.wg.Wait()
			p.pusherrs(pf.Errors...)
		}
	}
}

func (p *Parser) ParseWaitingImpls() {
	for _, i := range p.waitingImpls {
		p.Impl(i)
	}
	p.waitingImpls = nil
}

func (p *Parser) checkParamDefaultExprWithDefault(param *Param) {
	if types.IsFn(param.DataType) {
		p.pusherrtok(param.Token, "invalid_type_for_default_arg", param.DataType.Kind)
	}
}

func (p *Parser) checkParamDefaultExpr(f *Fn, param *Param) {
	if !paramHasDefaultArg(param) || param.Token.Id == lexer.ID_NA {
		return
	}
	if param.Default.Model != nil {
		if param.Default.Model.String() == build.CPP_DEFAULT_EXPR {
			p.checkParamDefaultExprWithDefault(param)
			return
		}
	}
	v, model := p.eval_expr(param.Default, nil)
	param.Default.Model = model
	p.checkArgType(param, v, param.Token)
}

func (p *Parser) param(f *Fn, param *Param) (err bool) {
	p.checkParamDefaultExpr(f, param)
	return
}

func (p *Parser) check_param_dup(params []ast.Param) (err bool) {
	for i, param := range params {
		for j, jparam := range params {
			if j >= i {
				break
			} else if param.Id == jparam.Id {
				err = true
				p.pusherrtok(param.Token, "exist_id", param.Id)
			}
		}
	}
	return
}

func (p *Parser) params(f *Fn) (err bool) {
	hasDefaultArg := false
	for i := range f.Params {
		param := &f.Params[i]
		err = err || p.param(f, param)
		if !hasDefaultArg {
			hasDefaultArg = paramHasDefaultArg(param)
			continue
		} else if !paramHasDefaultArg(param) && !param.Variadic {
			p.pusherrtok(param.Token, "param_must_have_default_arg", param.Id)
			err = true
		}
	}
	return
}

func (p *Parser) block_variables_of_fn(f *Fn) []*Var {
	vars := p.varsFromParams(f)
	vars = append(vars, f.RetType.Vars(f.Block)...)
	if f.Receiver != nil {
		s := f.Receiver.Tag.(*Struct)
		vars = append(vars, s.SelfVar(f.Receiver))
	}
	return vars
}

func (p *Parser) parse_pure_fn(f *Fn) (err bool) {
	hasError := p.eval.has_error
	owner := f.Owner.(*Parser)
	err = owner.params(f)
	if err {
		return
	}
	owner.block_vars = owner.block_variables_of_fn(f)
	owner.check_fn(f)
	if owner != p {
		owner.wg.Wait()
		p.pusherrs(owner.Errors...)
		owner.Errors = nil
	}
	owner.blockTypes = nil
	owner.block_vars = nil
	p.eval.has_error = hasError
	return
}

func (p *Parser) parse_fn(f *Fn) (err bool) {
	if len(f.Generics) > 0 {
		return false
	}
	return p.parse_pure_fn(f)
}

func (p *Parser) check_fns() {
	err := false
	check := func(f *Fn) {
		if f.BuiltinCaller != nil || len(f.Generics) > 0 {
			return
		}
		p.check_fn_special_cases(f)
		if err {
			return
		}
		p.blockTypes = nil
		err = p.parse_fn(f)
	}
	for _, f := range p.Defines.Fns {
		check(f)
	}
}

func (p *Parser) parse_struct_fn(s *Struct, f *Fn) (err bool) {
	if len(f.Generics) > 0 {
		return
	}
	if len(s.Generics) == 0 {
		p.parseTypesNonGenerics(f)
		return p.parse_fn(f)
	}
	return
}

func (p *Parser) checkStruct(xs *Struct) (err bool) {
	for _, f := range xs.Defines.Fns {
		p.blockTypes = nil
		err = p.parse_struct_fn(xs, f)
		if err {
			break
		}
	}
	return
}

func (p *Parser) check_structs() {
	err := false
	check := func(xs *Struct) {
		if err {
			return
		}
		p.checkStruct(xs)
	}
	for _, s := range p.Defines.Structs {
		check(s)
	}
}

func (p *Parser) check_fn_special_cases(f *Fn) {
	switch f.Id {
	case jane.ENTRY_POINT, jane.INIT_FN:
		p.checkSolidFuncSpecialCases(f)
	}
}

func (p *Parser) call_fn(f *Fn, data call_data, m *expr_model) value {
	v := p.parse_fn_call_toks(f, data.generics, data.args, m)
	v.lvalue = types.IsLvalue(v.data.DataType)
	return v
}

func (p *Parser) callStructConstructor(s *Struct, argsToks []lexer.Token, m *expr_model) (v value) {
	f := s.Constructor
	s = f.RetType.DataType.Tag.(*Struct)
	v.data.DataType = f.RetType.DataType.Copy()
	v.data.DataType.Kind = s.AsTypeKind()
	v.is_type = false
	v.lvalue = false
	v.constant = false
	v.data.Value = s.Id

	argsToks[0].Kind = lexer.KND_LPAREN
	argsToks[len(argsToks)-1].Kind = lexer.KND_RPARENT

	args := p.get_args(argsToks, true)
	if s.CppLinked {
		m.append_sub(exprNode{lexer.KND_LPAREN})
		m.append_sub(exprNode{f.RetType.String()})
		m.append_sub(exprNode{lexer.KND_RPARENT})
	} else {
		m.append_sub(exprNode{f.RetType.String()})
	}
	if s.CppLinked {
		m.append_sub(exprNode{lexer.KND_LBRACE})
	} else {
		m.append_sub(exprNode{lexer.KND_LPAREN})
	}
	p.parseArgs(f, args, m, f.Token)
	if m != nil {
		m.append_sub(argsExpr{args.Src})
	}
	if s.CppLinked {
		m.append_sub(exprNode{lexer.KND_RBRACE})
	} else {
		m.append_sub(exprNode{lexer.KND_RPARENT})
	}
	return v
}

func (p *Parser) parse_field(s *Struct, f **Var, i int) {
	*f = p.variable(**f)
	v := *f
	param := ast.Param{Id: v.Id, DataType: v.DataType}
	if !types.IsPtr(v.DataType) && types.IsStruct(v.DataType) {
		ts := v.DataType.Tag.(*Struct)
		if s.IsSameBase(ts) || ts.IsDependedTo(s) {
			p.pusherrtok(v.DataType.Token, "illegal_cycle_in_declaration", s.Id)
		} else {
			s.Origin.Depends = append(s.Origin.Depends, ts)
		}
	}
	if has_expr(v.Expr) {
		param.Default = v.Expr
	} else {
		param.Default.Model = exprNode{v.DataType.InitValue()}
	}
	s.Constructor.Params[i] = param
	s.Defines.Globals[i] = v
}

func (p *Parser) structConstructorInstance(as *Struct) *Struct {
	s := new(Struct)
	*s = *as
	s.Origin = as
	s.Constructor = new(Fn)
	*s.Constructor = *as.Constructor
	s.Constructor.RetType.DataType.Tag = s
	s.Defines = as.Defines
	for i := range s.Defines.Fns {
		f := &s.Defines.Fns[i]
		nf := new(Fn)
		*nf = **f
		nf.Receiver.Tag = s
		*f = nf
	}
	return s
}

func (p *Parser) check_anon_fn(f *Fn) {
	_ = p.check_param_dup(f.Params)
	p.check_ret_variables(f)
	p.reload_fn_types(f)
	globals := p.Defines.Globals
	blockVariables := p.block_vars
	p.Defines.Globals = append(blockVariables, p.Defines.Globals...)
	p.block_vars = p.block_variables_of_fn(f)
	rootBlock := p.rootBlock
	nodeBlock := p.nodeBlock
	p.check_fn(f)
	p.rootBlock = rootBlock
	p.nodeBlock = nodeBlock
	p.Defines.Globals = globals
	p.block_vars = blockVariables
}

func (p *Parser) get_args(toks []lexer.Token, targeting bool) *ast.Args {
	toks, _ = p.get_range(lexer.KND_LPAREN, lexer.KND_RPARENT, toks)
	if toks == nil {
		toks = make([]lexer.Token, 0)
	}
	r := new(builder)
	args := r.Args(toks, targeting)
	if len(r.Errors) > 0 {
		p.pusherrs(r.Errors...)
		args = nil
	}
	return args
}

func (p *Parser) get_generics(toks []lexer.Token) (_ []Type, err bool) {
	if len(toks) == 0 {
		return nil, false
	}
	toks = toks[1 : len(toks)-1]
	parts, errs := ast.Parts(toks, lexer.ID_COMMA, true)
	generics := make([]Type, len(parts))
	p.pusherrs(errs...)
	for i, part := range parts {
		if len(part) == 0 {
			continue
		}
		r := new_builder(nil)
		j := 0
		generic, _ := r.DataType(part, &j, true)
		if j+1 < len(part) {
			p.pusherrtok(part[j+1], "invalid_syntax")
		}
		p.pusherrs(r.Errors...)
		generics[i] = generic
		ok := p.fn_parse_type(&generics[i])
		if !ok {
			err = true
		}
	}
	return generics, err
}

func (p *Parser) checkGenericsQuantity(required, given int, errTok lexer.Token) bool {
	switch {
	case required == 0 && given > 0:
		p.pusherrtok(errTok, "not_has_generics")
		return false
	case required > 0 && given == 0:
		p.pusherrtok(errTok, "has_generics")
		return false
	case required < given:
		p.pusherrtok(errTok, "generics_overflow")
		return false
	case required > given:
		p.pusherrtok(errTok, "missing_generics")
		return false
	default:
		return true
	}
}

func (p *Parser) pushGeneric(generic *GenericType, source Type, errtok lexer.Token) {
	if types.IsEnum(source) {
		p.pusherrtok(errtok, "enum_not_supports_as_generic")
	}
	alias := &TypeAlias{
		Id:         generic.Id,
		Token:      generic.Token,
		TargetType: source,
		Used:       true,
		Generic:    true,
	}
	p.blockTypes = append(p.blockTypes, alias)
}

func (p *Parser) pushGenerics(generics []*GenericType, sources []Type, errtok lexer.Token) {
	for i, generic := range generics {
		p.pushGeneric(generic, sources[i], errtok)
	}
}

func (p *Parser) fn_parse_type(t *Type) bool {
	pt, ok := p.realType(*t, true)
	if ok && types.IsArray(pt) && pt.Size.AutoSized {
		p.pusherrtok(pt.Token, "invalid_type")
		ok = false
	}
	*t = pt
	return ok
}

func (p *Parser) reload_fn_types(f *Fn) {
	for i := range f.Params {
		_ = p.fn_parse_type(&f.Params[i].DataType)
	}
	ok := p.fn_parse_type(&f.RetType.DataType)
	if ok && types.IsArray(f.RetType.DataType) {
		p.pusherrtok(f.RetType.DataType.Token, "invalid_type")
	}
}

func itsCombined(f *Fn, generics []Type) bool {
	if f.Combines == nil {
		return true
	}
	for _, combine := range *f.Combines {
		for i, gt := range generics {
			ct := combine[i]
			if types.Equals(gt, ct) {
				return true
			}
		}
	}
	return false
}

func ready_to_parse_generic_fn(f *Fn) {
	owner := f.Owner.(*Parser)
	if f.Receiver != nil {
		s := f.Receiver.Tag.(*Struct)
		owner.pushGenerics(s.Generics, s.GetGenerics(), f.Token)
	}
	owner.reload_fn_types(f)
}

func (p *Parser) parseGenericFn(f *Fn, args *ast.Args, errtok lexer.Token) {
	if !is_multi_ret_as_args(f, len(args.Src)) {
		ready_to_parse_generic_fn(f)
	}

	if f.Block == nil {
		return
	} else if itsCombined(f, args.Generics) {
		return
	}
	*f.Combines = append(*f.Combines, args.Generics)
	p.parse_pure_fn(f)
}

func (p *Parser) parseGenerics(f *Fn, args *ast.Args, errTok lexer.Token) bool {
	if len(f.Generics) > 0 && len(args.Generics) == 0 {
		for _, g := range f.Generics {
			ok := false
			for _, param := range f.Params {
				if types.HasThisGeneric(g, param.DataType) {
					ok = true
					break
				}
			}
			if !ok {
				goto check
			}
		}
		args.DynamicGenericAnnotation = true
		i := 0
		for ; i < len(f.Generics); i++ {
			args.Generics = append(args.Generics, Type{})
		}
		goto ok
	}
check:
	if !p.checkGenericsQuantity(len(f.Generics), len(args.Generics), errTok) {
		return false
	} else {
		owner := f.Owner.(*Parser)
		owner.pushGenerics(f.Generics, args.Generics, errTok)
		owner.reload_fn_types(f)
	}
ok:
	return true
}

func (p *Parser) parse_fn_call(f *Fn, args *ast.Args, m *expr_model, errTok lexer.Token) (v value) {
	args.NeedsPureType = p.rootBlock == nil || len(p.rootBlock.Func.Generics) == 0
	if f.Receiver != nil {
		switch f.Receiver.Tag.(type) {
		case *Struct:
			owner := f.Owner.(*Parser)
			s := f.Receiver.Tag.(*Struct)
			generics := s.GetGenerics()
			if len(generics) > 0 {
				owner.pushGenerics(s.Generics, generics, errTok)
				if len(f.Generics) == 0 {
					owner.reload_fn_types(f)
				}
			}
		}
	}
	if len(f.Generics) > 0 {
		params := make([]Param, len(f.Params))
		for i := range params {
			param := &params[i]
			fparam := &f.Params[i]
			*param = *fparam
			param.DataType = fparam.DataType.Copy()
		}
		retType := f.RetType.DataType.Copy()
		owner := f.Owner.(*Parser)
		rootBlock := owner.rootBlock
		nodeBlock := owner.nodeBlock
		blockVars := owner.block_vars
		blockTypes := owner.blockTypes
		defer func() {
			owner.rootBlock = rootBlock
			owner.nodeBlock = nodeBlock
			owner.block_vars = blockVars
			owner.blockTypes = blockTypes

			for i := range params {
				params[i].DataType.Generic = f.Params[i].DataType.Generic
			}
			retType.Generic = f.RetType.DataType.Generic

			f.Params = params
			f.RetType.DataType = retType
		}()
		if !p.parseGenerics(f, args, errTok) {
			return
		}
	}
	if args == nil {
		goto end
	}
	p.parseArgs(f, args, m, errTok)
	if len(args.Generics) > 0 {
		p.parseGenericFn(f, args, errTok)
	}
	if m != nil {
		model := callExpr{
			args: argsExpr{args.Src},
			f:    f,
		}
		if !is_multi_ret_as_args(f, len(args.Src)) {
			model.generics = genericsExpr{args.Generics}
		}
		m.append_sub(model)
	}
end:
	v.mutable = true
	v.data.Value = " "
	v.data.DataType = f.RetType.DataType.Copy()
	if args.NeedsPureType {
		v.data.DataType.Pure = true
		v.data.DataType.Original = nil
	}
	return
}

func (p *Parser) parse_fn_call_toks(
	f *Fn,
	genericsToks, argsToks []lexer.Token,
	m *expr_model,
) (v value) {
	var generics []Type
	var args *ast.Args
	var err bool
	generics, err = p.get_generics(genericsToks)
	if err {
		p.eval.has_error = true
		return
	}
	args = p.get_args(argsToks, false)
	args.Generics = generics
	return p.parse_fn_call(f, args, m, argsToks[0])
}

func (p *Parser) parseStructArgs(f *Fn, args *ast.Args, errTok lexer.Token) {
	sap := structArgParser{
		p:      p,
		f:      f,
		args:   args,
		errTok: errTok,
	}
	sap.parse()
}

func (p *Parser) parsePureArgs(f *Fn, args *ast.Args, m *expr_model, errTok lexer.Token) {
	pap := pureArgParser{
		p:      p,
		f:      f,
		args:   args,
		errTok: errTok,
		m:      m,
	}
	pap.parse()
}

func (p *Parser) parseArgs(f *Fn, args *ast.Args, m *expr_model, errTok lexer.Token) {
	if args.Targeted {
		p.parseStructArgs(f, args, errTok)
		return
	}
	p.parsePureArgs(f, args, m, errTok)
}

func has_expr(expr Expr) bool {
	return len(expr.Tokens) > 0 || expr.Model != nil
}

func paramHasDefaultArg(param *Param) bool {
	return has_expr(param.Default)
}

type (
	paramMap     map[string]*paramMapPair
	paramMapPair struct {
		param *Param
		arg   *Arg
	}
)

func find_generic_alias(f *Fn, g *ast.GenericType) *ast.TypeAlias {
	owner := f.Owner.(*Parser)
	alias, _ := owner.block_type_by_id(g.Id)
	return alias
}

func (p *Parser) pushGenericByFunc(f *Fn, pair *paramMapPair, args *ast.Args, gt Type) bool {
	tf := gt.Tag.(*Fn)
	cf := pair.param.DataType.Tag.(*Fn)
	if len(tf.Params) != len(cf.Params) {
		return false
	}
	for i, param := range tf.Params {
		pair := *pair
		pair.param = &cf.Params[i]
		ok := p.pushGenericByArg(f, &pair, args, param.DataType)
		if !ok {
			return ok
		}
	}
	{
		pair := *pair
		pair.param = &ast.Param{
			DataType: cf.RetType.DataType,
		}
		return p.pushGenericByArg(f, &pair, args, tf.RetType.DataType)
	}
}

func (p *Parser) pushGenericByMap(f *Fn, pair *paramMapPair, args *ast.Args, gt Type) bool {
	if !types.IsMap(pair.param.DataType) {
		return false
	}
	arg_types := gt.Tag.([]Type)
	param_types := pair.param.DataType.Tag.([]Type)
	check := func(offset int) bool {
		for _, g := range f.Generics {
			if !types.HasThisGeneric(g, param_types[offset]) {
				continue
			}
			alias := find_generic_alias(f, g)
			if alias == nil {
				pt := pair.param.DataType
				pair.param.DataType = param_types[offset]
				s := p.pushGenericByArg(f, pair, args, arg_types[offset])
				pair.param.DataType = pt
				return s
			} else {
				v := value{}
				v.data.Value = " "
				v.data.DataType = arg_types[offset]
				pt := pair.param.DataType
				pair.param.DataType = alias.TargetType
				p.checkArgType(pair.param, v, pair.arg.Token)
				pair.param.DataType = pt
			}
			return true
		}
		return false
	}
	return check(0) && check(1)
}

func (p *Parser) pushGenericByMultiTyped(f *Fn, pair *paramMapPair, args *ast.Args, gt Type) bool {
	_types := gt.Tag.([]Type)
	for _, mt := range _types {
		for _, g := range f.Generics {
			if !types.HasThisGeneric(g, pair.param.DataType) {
				continue
			}
			alias := find_generic_alias(f, g)
			if alias == nil {
				if !p.pushGenericByArg(f, pair, args, mt) {
					return false
				}
			} else {
				v := value{}
				v.data.Value = " "
				v.data.DataType = mt
				pt := pair.param.DataType
				pair.param.DataType = alias.TargetType
				p.checkArgType(pair.param, v, pair.arg.Token)
				pair.param.DataType = pt
			}
			break
		}
	}
	return true
}

func (p *Parser) pushGenericByCommonArg(f *Fn, pair *paramMapPair, args *ast.Args, t Type) bool {
	for i, g := range f.Generics {
		if !types.IsThisGeneric(g, pair.param.DataType) {
			continue
		}
		alias := find_generic_alias(f, g)
		if alias == nil {
			p.pushGenericByType(f, g, i, args, t)
		} else {
			v := value{}
			v.data.Value = " "
			v.data.DataType = t
			pt := pair.param.DataType
			pair.param.DataType = alias.TargetType
			p.checkArgType(pair.param, v, pair.arg.Token)
			pair.param.DataType = pt
		}
		return true
	}
	return false
}

func (p *Parser) pushGenericByType(f *Fn, g *GenericType, pos int, args *ast.Args, gt Type) {
	id, _ := gt.KindId()
	gt.Kind = id
	f.Owner.(*Parser).pushGeneric(g, gt, f.Token)
	args.Generics[pos] = gt
}

func (p *Parser) pushGenericByComponent(
	f *Fn,
	pair *paramMapPair,
	args *ast.Args,
	argType Type,
) bool {
	for argType.ComponentType != nil {
		argType = *argType.ComponentType
	}
	return p.pushGenericByCommonArg(f, pair, args, argType)
}

func (p *Parser) pushGenericByArg(f *Fn, pair *paramMapPair, args *ast.Args, argType Type) bool {
	_, prefix := pair.param.DataType.KindId()
	_, tprefix := argType.KindId()
	if prefix != tprefix {
		return p.pushGenericByCommonArg(f, pair, args, argType)
	}
	switch {
	case types.IsFn(argType):
		return p.pushGenericByFunc(f, pair, args, argType)
	case argType.MultiTyped:
		return p.pushGenericByMultiTyped(f, pair, args, argType)
	case types.IsMap(argType):
		return p.pushGenericByMap(f, pair, args, argType)
	case types.IsArray(argType), types.IsSlice(argType):
		return p.pushGenericByComponent(f, pair, args, argType)
	default:
		return p.pushGenericByCommonArg(f, pair, args, argType)
	}
}

func (p *Parser) check_arg(f *Fn, pair *paramMapPair, args *ast.Args, variadiced *bool, v value) {
	if variadiced != nil && !*variadiced {
		*variadiced = v.variadic
	}
	if args.DynamicGenericAnnotation &&
		types.HasGenerics(f.Generics, pair.param.DataType) {
		ok := p.pushGenericByArg(f, pair, args, v.data.DataType)
		if !ok {
			p.pusherrtok(pair.arg.Token, "dynamic_type_annotation_failed")
		}
		return
	}
	p.checkArgType(pair.param, v, pair.arg.Token)
}

func (p *Parser) parse_arg(f *Fn, pair *paramMapPair, args *ast.Args, variadiced *bool) {
	var v value
	var model ast.ExprModel
	if pair.param.Variadic {
		t := types.ToSlice(pair.param.DataType)
		v, model = p.eval_expr(pair.arg.Expr, &t)
	} else {
		v, model = p.eval_expr(pair.arg.Expr, &pair.param.DataType)
	}
	pair.arg.Expr.Model = model
	p.check_arg(f, pair, args, variadiced, v)
}

func (p *Parser) check_assign_type(real Type, val value, errTok lexer.Token) {
	assign_checker{
		p:      p,
		t:      real,
		v:      val,
		errtok: errTok,
	}.check()
}

func (p *Parser) checkArgType(param *Param, val value, errTok lexer.Token) {
	p.check_valid_init_expr(param.Mutable, val, errTok)
	p.check_assign_type(param.DataType, val, errTok)
}

func (p *Parser) get_range(open, close string, toks []lexer.Token) (_ []lexer.Token, ok bool) {
	i := 0
	toks = ast.Range(&i, open, close, toks)
	return toks, toks != nil
}

func (p *Parser) checkSolidFuncSpecialCases(f *Fn) {
	if len(f.Params) > 0 {
		p.pusherrtok(f.Token, "fn_have_parameters", f.Id)
	}
	if f.RetType.DataType.Id != types.VOID {
		p.pusherrtok(f.RetType.DataType.Token, "fn_have_ret", f.Id)
	}
	f.Attributes = nil
	if f.IsUnsafe {
		p.pusherrtok(f.Token, "fn_is_unsafe", f.Id)
	}
}

func (p *Parser) checkNewBlockCustom(b *ast.Block, oldBlockVars []*Var) {
	b.Gotos = new(ast.Gotos)
	b.Labels = new(ast.Labels)
	if p.rootBlock == nil {
		p.rootBlock = b
		p.nodeBlock = b
		defer func() {
			p.checkLabelNGoto()
			p.rootBlock = nil
			p.nodeBlock = nil
		}()
	} else {
		b.Parent = p.nodeBlock
		b.SubIndex = p.nodeBlock.SubIndex + 1
		b.Func = p.nodeBlock.Func
		oldNode := p.nodeBlock
		old_unsafe := b.IsUnsafe
		b.IsUnsafe = b.IsUnsafe || oldNode.IsUnsafe
		p.nodeBlock = b
		defer func() {
			p.nodeBlock = oldNode
			b.IsUnsafe = old_unsafe
			*p.rootBlock.Gotos = append(*p.rootBlock.Gotos, *b.Gotos...)
			*p.rootBlock.Labels = append(*p.rootBlock.Labels, *b.Labels...)
		}()
	}
	blockTypes := p.blockTypes
	p.checkNodeBlock()

	vars := p.block_vars[len(oldBlockVars):]
	aliases := p.blockTypes[len(blockTypes):]
	for _, v := range vars {
		if !v.Used {
			p.pusherrtok(v.Token, "declared_but_not_used", v.Id)
		}
	}
	for _, a := range aliases {
		if !a.Used {
			p.pusherrtok(a.Token, "declared_but_not_used", a.Id)
		}
	}
	p.block_vars = oldBlockVars
	p.blockTypes = blockTypes
}

func (p *Parser) checkNewBlock(b *ast.Block) {
	p.checkNewBlockCustom(b, p.block_vars)
}

func (p *Parser) fall_st(f *ast.Fall, i *int) {
	switch {
	case p.currentCase == nil || *i+1 < len(p.nodeBlock.Tree):
		p.pusherrtok(f.Token, "fallthrough_wrong_use")
		return
	case p.currentCase.Next == nil:
		p.pusherrtok(f.Token, "fallthrough_into_final_case")
		return
	}
	f.Case = p.currentCase
}

func (p *Parser) common_st(s *ast.St, i *int, recover bool) bool {
	switch data := s.Data.(type) {
	case ast.ExprSt:
		is_recover := p.expr_st(&data, i, recover)
		if !is_recover {
			s.Data = data
		}
	case Var:
		p.var_st(&data)
		s.Data = data
	case ast.Assign:
		p.assign(&data)
		s.Data = data
	case ast.Break:
		p.break_st(&data)
		s.Data = data
	case ast.Continue:
		p.continue_st(&data)
		s.Data = data
	case *ast.Match:
		p.matchcase(data)
	case TypeAlias:
		def, _, canshadow := p.block_define_by_id(data.Id)
		if def != nil && !canshadow {
			p.pusherrtok(data.Token, "exist_id", data.Id)
			break
		} else if lexer.IsIgnoreId(data.Id) {
			p.pusherrtok(data.Token, "ignore_id")
			break
		}
		data.TargetType, _ = p.realType(data.TargetType, true)
		p.blockTypes = append(p.blockTypes, &data)
	case *ast.Block:
		p.checkNewBlock(data)
		s.Data = data
	case ast.ConcurrentCall:
		p.concurrent_call(&data)
		s.Data = data
	case ast.Comment:
	default:
		return false
	}
	return true
}

func (p *Parser) check_st(i *int) {
	s := &p.nodeBlock.Tree[*i]
	if p.common_st(s, i, true) {
		return
	}
	switch data := s.Data.(type) {
	case ast.Iter:
		data.Parent = p.nodeBlock
		s.Data = data
		p.iter(&data)
		s.Data = data
	case ast.Fall:
		p.fall_st(&data, i)
		s.Data = data
	case ast.Conditional:
		p.conditional(&data)
		s.Data = data
	case ast.Ret:
		rc := ret_checker{p: p, ret_ast: &data, f: p.nodeBlock.Func}
		rc.check()
		s.Data = data
	case ast.Goto:
		node := new(ast.Goto)
		*node = data
		node.Index = *i
		node.Block = p.nodeBlock
		*p.nodeBlock.Gotos = append(*p.nodeBlock.Gotos, node)
	case ast.Label:
		if find_label_parent(data.Label, p.nodeBlock) != nil {
			p.pusherrtok(data.Token, "label_exist", data.Label)
			break
		}
		node := new(ast.Label)
		*node = data
		node.Index = *i
		node.Block = p.nodeBlock
		*p.nodeBlock.Labels = append(*p.nodeBlock.Labels, node)
	default:
		p.pusherrtok(s.Token, "invalid_syntax")
	}
}

func (p *Parser) checkNodeBlock() {
	for i := 0; i < len(p.nodeBlock.Tree); i++ {
		p.check_st(&i)
	}
}

func (p *Parser) recover_fn_expr_st(s *ast.ExprSt, i *int) {
	errtok := s.Expr.Tokens[0]
	callToks := s.Expr.Tokens[1:]
	args := p.get_args(callToks, false)
	handleParam := recoverFunc.Params[0]
	if len(args.Src) == 0 {
		p.pusherrtok(errtok, "missing_expr_for", handleParam.Id)
		return
	} else if len(args.Src) > 1 {
		p.pusherrtok(errtok, "argument_overflow")
	}
	v, _ := p.eval_expr(args.Src[0].Expr, nil)
	if v.data.DataType.Kind != handleParam.DataType.Kind {
		p.eval.push_err_tok(errtok, "incompatible_types",
			handleParam.DataType.Kind, v.data.DataType.Kind)
		return
	}
	rc := ast.RecoverCall{}
	rc.Try = block_from_tree(p.nodeBlock.Tree[*i+1:])
	p.checkNewBlock(rc.Try)
	rc.Handler = v.data.DataType.Tag.(*Fn)
	p.nodeBlock.Tree[*i].Data = rc
	p.nodeBlock.Tree = p.nodeBlock.Tree[:*i+1]
}

func (p *Parser) expr_st(s *ast.ExprSt, i *int, recover bool) (is_recover bool) {
	if s.Expr.IsNotBinop() {
		expr := s.Expr.Op.(ast.BinopExpr)
		tok := expr.Tokens[0]
		if tok.Id == lexer.ID_IDENT && tok.Kind == recoverFunc.Id {
			if ast.IsFnCall(s.Expr.Tokens) != nil {
				if !recover {
					p.pusherrtok(tok, "invalid_syntax")
				}
				def, _, _ := p.defined_by_id(tok.Kind)
				if def == recoverFunc {
					p.recover_fn_expr_st(s, i)
					return true
				}
			}
		}
	}
	if s.Expr.Model == nil {
		_, s.Expr.Model = p.eval_expr(s.Expr, nil)
	}
	return
}

func (p *Parser) count_match_type(m *ast.Match, t Type) int {
	n := 0
loop:
	for _, c := range m.Cases {
		for _, expr := range c.Exprs {
			if expr.Model == nil {
				break loop
			}

			if types.Equals(t, expr.Op.(Type)) {
				n++
			}
		}
	}
	return n
}

func (p *Parser) parseCase(m *ast.Match, c *ast.Case) {
	for i := range c.Exprs {
		expr := &c.Exprs[i]
		switch expr.Op.(type) {
		case Type:
			t, _ := p.realType(expr.Op.(Type), true)
			expr.Op = t
			expr.Model = exprNode{m.Expr.String() + ".__type_is<" + t.String() + ">()"}
			if p.count_match_type(m, t) > 1 {
				p.pusherrtok(c.Token, "duplicate_match_type", t.Kind)
			}
		default:
			value, model := p.eval_expr(*expr, nil)
			expr.Model = model
			assign_checker{
				p:      p,
				t:      m.ExprType,
				v:      value,
				errtok: expr.Tokens[0],
			}.check()
		}
	}
	oldCase := p.currentCase
	p.currentCase = c
	p.checkNewBlock(c.Block)
	p.currentCase = oldCase
}

func (p *Parser) cases(m *ast.Match) {
	for i := range m.Cases {
		p.parseCase(m, &m.Cases[i])
	}
}

func (p *Parser) matchcase(m *ast.Match) {
	if !m.Expr.IsEmpty() {
		value, expr_model := p.eval_expr(m.Expr, nil)
		m.Expr.Model = expr_model
		m.ExprType = value.data.DataType
	} else {
		m.ExprType.Id = types.BOOL
		m.ExprType.Kind = types.TYPE_MAP[m.ExprType.Id]
	}
	p.cases(m)
	if m.Default != nil {
		p.parseCase(m, m.Default)
	}
}

func find_label(id string, b *ast.Block) *ast.Label {
	for _, label := range *b.Labels {
		if label.Label == id {
			return label
		}
	}
	return nil
}

func (p *Parser) checkLabels() {
	labels := p.rootBlock.Labels
	for _, label := range *labels {
		if !label.Used {
			p.pusherrtok(label.Token, "declared_but_not_used", label.Label+":")
		}
	}
}

func stIsDef(s *ast.St) bool {
	switch t := s.Data.(type) {
	case Var:
		return true
	case ast.Assign:
		for _, selector := range t.Left {
			if selector.Var.New {
				return true
			}
		}
	}
	return false
}

func (p *Parser) checkSameScopeGoto(gt *ast.Goto, label *ast.Label) {
	if label.Index < gt.Index {
		return
	}
	for i := label.Index; i > gt.Index; i-- {
		s := &label.Block.Tree[i]
		if stIsDef(s) {
			p.pusherrtok(gt.Token, "goto_jumps_declarations", gt.Label)
			break
		}
	}
}

func (p *Parser) checkLabelParents(gt *ast.Goto, label *ast.Label) bool {
	block := label.Block
parent_scopes:
	if block.Parent != nil && block.Parent != gt.Block {
		block = block.Parent
		for i := 0; i < len(block.Tree); i++ {
			s := &block.Tree[i]
			switch {
			case s.Token.Row >= label.Token.Row:
				return true
			case stIsDef(s):
				p.pusherrtok(gt.Token, "goto_jumps_declarations", gt.Label)
				return false
			}
		}
		goto parent_scopes
	}
	return true
}

func (p *Parser) checkGotoScope(gt *ast.Goto, label *ast.Label) {
	for i := gt.Index; i < len(gt.Block.Tree); i++ {
		s := &gt.Block.Tree[i]
		switch {
		case s.Token.Row >= label.Token.Row:
			return
		case stIsDef(s):
			p.pusherrtok(gt.Token, "goto_jumps_declarations", gt.Label)
			return
		}
	}
}

func (p *Parser) checkDiffScopeGoto(gt *ast.Goto, label *ast.Label) {
	switch {
	case label.Block.SubIndex > 0 && gt.Block.SubIndex == 0:
		if !p.checkLabelParents(gt, label) {
			return
		}
	case label.Block.SubIndex < gt.Block.SubIndex:
		return
	}
	block := label.Block
	for i := label.Index - 1; i >= 0; i-- {
		s := &block.Tree[i]
		switch s.Data.(type) {
		case ast.Block:
			if s.Token.Row <= gt.Token.Row {
				return
			}
		}
		if stIsDef(s) {
			p.pusherrtok(gt.Token, "goto_jumps_declarations", gt.Label)
			break
		}
	}
	if block.Parent != nil && block.Parent != gt.Block {
		_ = p.checkLabelParents(gt, label)
	} else {
		p.checkGotoScope(gt, label)
	}
}

func (p *Parser) checkGoto(gt *ast.Goto, label *ast.Label) {
	switch {
	case gt.Block == label.Block:
		p.checkSameScopeGoto(gt, label)
	case label.Block.SubIndex > 0:
		p.checkDiffScopeGoto(gt, label)
	}
}

func (p *Parser) checkGotos() {
	for _, gt := range *p.rootBlock.Gotos {
		label := find_label(gt.Label, p.rootBlock)
		if label == nil {
			p.pusherrtok(gt.Token, "label_not_exist", gt.Label)
			continue
		}
		label.Used = true
		p.checkGoto(gt, label)
	}
}

func (p *Parser) checkLabelNGoto() {
	p.checkGotos()
	p.checkLabels()
}

func match_has_ret(m *ast.Match) (ok bool) {
	if m.Default == nil {
		return false
	}
	ok = true
	fall := false
	for _, c := range m.Cases {
		falled := fall
		ok, fall = has_ret(c.Block)
		if falled && !ok && !fall {
			return false
		}
		switch {
		case !ok:
			if !fall {
				return false
			}
			fallthrough
		case fall:
			if c.Next == nil {
				return false
			}
			continue
		}
		fall = false
	}
	ok, _ = has_ret(m.Default.Block)
	return ok
}

func conditional_has_ret(c *ast.Conditional) (ok bool) {
	if c.Default == nil {
		return false
	}
	ok = true
	for _, elif := range c.Elifs {
		ok, _ = has_ret(elif.Block)
		if !ok {
			return false
		}
	}
	ok, _ = has_ret(c.Default.Block)
	return ok
}

func has_ret(b *ast.Block) (ok bool, fall bool) {
	if b == nil {
		return false, false
	}
	for _, s := range b.Tree {
		switch t := s.Data.(type) {
		case *ast.Block:
			ok, fall = has_ret(t)
			if ok {
				return true, fall
			}
		case ast.Fall:
			fall = true
		case ast.Ret:
			return true, fall
		case ast.RecoverCall:
			ok, fall = has_ret(t.Try)
			if ok {
				return true, fall
			}
		case *ast.Match:
			if match_has_ret(t) {
				return true, false
			}
		case ast.Conditional:
			if conditional_has_ret(&t) {
				return true, false
			}
		}
	}
	return false, fall
}

func (p *Parser) checkRets(f *Fn) {
	ok, _ := has_ret(f.Block)
	if ok {
		return
	}
	if !types.IsVoid(f.RetType.DataType) {
		p.pusherrtok(f.Token, "missing_ret")
	}
}

func (p *Parser) check_fn(f *Fn) {
	if f.Block == nil || f.Block.Tree == nil {
		goto always
	} else {
		rootBlock := p.rootBlock
		nodeBlock := p.nodeBlock
		p.rootBlock = nil
		p.nodeBlock = nil
		f.Block.Func = f
		p.checkNewBlock(f.Block)
		p.rootBlock = rootBlock
		p.nodeBlock = nodeBlock
	}
always:
	p.checkRets(f)
}

func (p *Parser) var_st(v *Var) {
	def, _, canshadow := p.block_define_by_id(v.Id)
	if !canshadow && def != nil {
		p.pusherrtok(v.Token, "exist_id", v.Id)
		return
	}
	*v = *p.variable(*v)
	v.Owner = p.nodeBlock
	p.block_vars = append(p.block_vars, v)
}

func (p *Parser) concurrent_call(cc *ast.ConcurrentCall) {
	m := new(expr_model)
	m.nodes = make([]expr_build_node, 1)
	_, cc.Expr.Model = p.eval_expr(cc.Expr, nil)
}

func (p *Parser) check_assign(left value, errtok lexer.Token) bool {
	state := true
	if !left.lvalue {
		p.eval.push_err_tok(errtok, "assign_require_lvalue")
		state = false
	}
	if left.constant {
		p.pusherrtok(errtok, "assign_const")
		state = false
	} else if !left.mutable {
		p.pusherrtok(errtok, "assignment_to_non_mut")
	}
	switch left.data.DataType.Tag.(type) {
	case Fn:
		f, _, _ := p.fn_by_id(left.data.Token.Kind)
		if f != nil {
			p.pusherrtok(errtok, "assign_type_not_support_value")
			state = false
		}
	}
	return state
}

func (p *Parser) single_assign(assign *ast.Assign, l, r []value) {
	left := l[0]
	switch {
	case lexer.IsIgnoreId(left.data.Value):
		return
	case !p.check_assign(left, assign.Setter):
		return
	}
	right := r[0]
	if assign.Setter.Kind != lexer.KND_EQ && !lexer.IsLiteral(right.data.Value) {
		assign.Setter.Kind = assign.Setter.Kind[:len(assign.Setter.Kind)-1]
		solver := solver{
			p:  p,
			l:  left,
			r:  right,
			op: assign.Setter,
		}
		right = solver.solve()
		assign.Setter.Kind += lexer.KND_EQ
	}
	assign_checker{
		p:      p,
		t:      left.data.DataType,
		v:      right,
		errtok: assign.Setter,
	}.check()
}

func (p *Parser) assign_exprs(ast *ast.Assign) (l []value, r []value) {
	l = make([]value, len(ast.Left))
	r = make([]value, len(ast.Right))
	n := len(l)
	if n < len(r) {
		n = len(r)
	}
	for i := 0; i < n; i++ {
		var r_type *Type = nil
		if i < len(l) {
			left := &ast.Left[i]

			set_existing := func() {
				v, model := p.eval_expr(left.Expr, nil)
				left.Expr.Model = model
				l[i] = v
				r_type = &v.data.DataType
			}

			if !left.Var.New &&
				!(len(left.Expr.Tokens) == 1 && lexer.IsIgnoreId(left.Expr.Tokens[0].Kind)) {
				set_existing()
			} else {
				def, _, canshadow := p.block_define_by_id(left.Var.Id)
				if !canshadow && def != nil {
					left.Var.New = false
					set_existing()
				} else {
					l[i].data.Value = lexer.IGNORE_ID
				}
			}
		}

		if i < len(r) {
			left := &ast.Right[i]
			v, model := p.eval_expr(*left, r_type)
			left.Model = model
			r[i] = v
		}
	}
	return
}

func (p *Parser) funcMultiAssign(vsAST *ast.Assign, l, r []value) {
	types := r[0].data.DataType.Tag.([]Type)
	if len(types) > len(vsAST.Left) {
		p.pusherrtok(vsAST.Setter, "missing_multi_assign_identifiers")
		return
	} else if len(types) < len(vsAST.Left) {
		p.pusherrtok(vsAST.Setter, "overflow_multi_assign_identifiers")
		return
	}
	rights := make([]value, len(types))
	for i, t := range types {
		rights[i] = value{data: ast.Data{Token: t.Token, DataType: t}}
	}
	p.multi_assign(vsAST, l, rights)
}

func (p *Parser) check_valid_init_expr(left_mutable bool, right value, errtok lexer.Token) {
	if p.unsafe_allowed() || !lexer.IsIdentifierRune(right.data.Value) {
		return
	}
	if left_mutable && !right.mutable && types.IsMut(right.data.DataType) {
		p.pusherrtok(errtok, "assignment_non_mut_to_mut")
		return
	}
	checker := assign_checker{
		p:      p,
		v:      right,
		errtok: errtok,
	}
	_ = checker.check_validity()
}

func (p *Parser) multi_assign(assign *ast.Assign, l []value, r []value) {
	for i := range assign.Left {
		left := &assign.Left[i]
		left.Ignore = lexer.IsIgnoreId(left.Var.Id)
		right := r[i]

		if left.Var.New {
			left.Var.Tag = right
			p.var_st(&left.Var)
			continue
		}

		if left.Ignore {
			continue
		}
		leftExpr := l[i]
		if !p.check_assign(leftExpr, assign.Setter) {
			return
		}
		p.check_valid_init_expr(leftExpr.mutable, right, assign.Setter)
		assign_checker{
			p:      p,
			t:      leftExpr.data.DataType,
			v:      right,
			errtok: assign.Setter,
		}.check()
	}
}

func (p *Parser) unsafe_allowed() bool {
	return (p.rootBlock != nil && p.rootBlock.IsUnsafe) ||
		(p.nodeBlock != nil && p.nodeBlock.IsUnsafe)
}

func (p *Parser) postfix(assign *ast.Assign, l, r []value) {
	if len(r) > 0 {
		p.pusherrtok(assign.Setter, "invalid_syntax")
		return
	}
	left := l[0]
	_ = p.check_assign(left, assign.Setter)
	if types.IsExplicitPtr(left.data.DataType) {
		if !p.unsafe_allowed() {
			p.pusherrtok(assign.Left[0].Expr.Tokens[0], "unsafe_behavior_at_out_of_unsafe_scope")
		}
		return
	}
	checkType := left.data.DataType
	if types.IsRef(checkType) {
		checkType = types.Elem(checkType)
	}
	if types.IsPure(checkType) && types.IsNumeric(checkType.Id) {
		return
	}
	p.pusherrtok(
		assign.Setter,
		"operator_not_for_janetype",
		assign.Setter.Kind,
		left.data.DataType.Kind,
	)
}

func (p *Parser) assign(assign *ast.Assign) {
	ln := len(assign.Left)
	rn := len(assign.Right)
	l, r := p.assign_exprs(assign)
	switch {
	case rn == 0 && ast.IsPostfixOp(assign.Setter.Kind):
		p.postfix(assign, l, r)
		return
	case ln == 1 && !assign.Left[0].Var.New:
		p.single_assign(assign, l, r)
		return
	case assign.Setter.Kind != lexer.KND_EQ:
		p.pusherrtok(assign.Setter, "invalid_syntax")
		return
	case rn == 1:
		right := r[0]
		if right.data.DataType.MultiTyped {
			assign.MultipleRet = true
			p.funcMultiAssign(assign, l, r)
			return
		}
	}
	switch {
	case ln > rn:
		p.pusherrtok(assign.Setter, "overflow_multi_assign_identifiers")
		return
	case ln < rn:
		p.pusherrtok(assign.Setter, "missing_multi_assign_identifiers")
		return
	}
	p.multi_assign(assign, l, r)
}

func (p *Parser) whileProfile(iter *ast.Iter) {
	profile := iter.Profile.(ast.IterWhile)
	val, model := p.eval_expr(profile.Expr, nil)
	profile.Expr.Model = model
	iter.Profile = profile
	if !p.eval.has_error && val.data.Value != "" && !is_bool_expr(val) {
		p.pusherrtok(iter.Token, "iter_while_require_bool_expr")
	}
	if profile.Next.Data != nil {
		_ = p.common_st(&profile.Next, nil, false)
	}
	p.checkNewBlock(iter.Block)
}

func (p *Parser) foreachProfile(iter *ast.Iter) {
	profile := iter.Profile.(ast.IterForeach)
	val, model := p.eval_expr(profile.Expr, nil)
	profile.Expr.Model = model
	profile.ExprType = val.data.DataType
	if !p.eval.has_error && val.data.Value != "" && !is_foreach_iter_expr(val) {
		p.pusherrtok(iter.Token, "iter_foreach_require_enumerable_expr")
	} else {
		fc := foreachChecker{p, &profile, val}
		fc.check()
	}
	iter.Profile = profile
	blockVars := p.block_vars
	if !lexer.IsIgnoreId(profile.KeyA.Id) {
		p.block_vars = append(p.block_vars, &profile.KeyA)
	}
	if !lexer.IsIgnoreId(profile.KeyB.Id) {
		p.block_vars = append(p.block_vars, &profile.KeyB)
	}
	p.checkNewBlockCustom(iter.Block, blockVars)
}

func (p *Parser) iter(iter *ast.Iter) {
	oldCase := p.currentCase
	oldIter := p.currentIter
	p.currentCase = nil
	p.currentIter = iter
	switch iter.Profile.(type) {
	case ast.IterWhile:
		p.whileProfile(iter)
	case ast.IterForeach:
		p.foreachProfile(iter)
	default:
		p.checkNewBlock(iter.Block)
	}
	p.currentCase = oldCase
	p.currentIter = oldIter
}

func (p *Parser) conditional_node(node *ast.If) {
	val, model := p.eval_expr(node.Expr, nil)
	node.Expr.Model = model
	if !p.eval.has_error && val.data.Value != "" && !is_bool_expr(val) {
		p.pusherrtok(node.Token, "if_require_bool_expr")
	}
	p.checkNewBlock(node.Block)
}

func (p *Parser) conditional(model *ast.Conditional) {
	p.conditional_node(model.If)
	for _, elif := range model.Elifs {
		p.conditional_node(elif)
	}
	if model.Default != nil {
		p.checkNewBlock(model.Default.Block)
	}
}

func find_label_parent(id string, b *ast.Block) *ast.Label {
	label := find_label(id, b)
	for label == nil {
		if b.Parent == nil {
			return nil
		}
		b = b.Parent
		label = find_label(id, b)
	}
	return label
}

func (p *Parser) breakWithLabel(brk *ast.Break) {
	if p.currentIter == nil && p.currentCase == nil {
		p.pusherrtok(brk.Token, "break_at_out_of_valid_scope")
		return
	}
	var label *ast.Label
	switch {
	case p.currentCase != nil && p.currentIter != nil:
		if p.currentCase.Block.Parent.SubIndex < p.currentIter.Parent.SubIndex {
			label = find_label_parent(brk.LabelToken.Kind, p.currentIter.Parent)
			if label == nil {
				label = find_label_parent(brk.LabelToken.Kind, p.currentCase.Block.Parent)
			}
		} else {
			label = find_label_parent(brk.LabelToken.Kind, p.currentCase.Block.Parent)
			if label == nil {
				label = find_label_parent(brk.LabelToken.Kind, p.currentIter.Parent)
			}
		}
	case p.currentCase != nil:
		label = find_label_parent(brk.LabelToken.Kind, p.currentCase.Block.Parent)
	case p.currentIter != nil:
		label = find_label_parent(brk.LabelToken.Kind, p.currentIter.Parent)
	}
	if label == nil {
		p.pusherrtok(brk.LabelToken, "label_not_exist", brk.LabelToken.Kind)
		return
	} else if label.Index+1 >= len(label.Block.Tree) {
		p.pusherrtok(brk.LabelToken, "invalid_label")
		return
	}
	label.Used = true
	for i := label.Index + 1; i < len(label.Block.Tree); i++ {
		node := &label.Block.Tree[i]
		if node.Data == nil {
			continue
		}
		switch data := node.Data.(type) {
		case ast.Comment:
			continue
		case *ast.Match:
			label.Used = true
			brk.Label = data.EndLabel()
		case ast.Iter:
			label.Used = true
			brk.Label = data.EndLabel()
		default:
			p.pusherrtok(brk.LabelToken, "invalid_label")
		}
		break
	}
}

func (p *Parser) continueWithLabel(cont *ast.Continue) {
	if p.currentIter == nil {
		p.pusherrtok(cont.Token, "continue_at_out_of_valid_scope")
		return
	}
	label := find_label_parent(cont.LoopLabel.Kind, p.currentIter.Parent)
	if label == nil {
		p.pusherrtok(cont.LoopLabel, "label_not_exist", cont.LoopLabel.Kind)
		return
	} else if label.Index+1 >= len(label.Block.Tree) {
		p.pusherrtok(cont.LoopLabel, "invalid_label")
		return
	}
	label.Used = true
	for i := label.Index + 1; i < len(label.Block.Tree); i++ {
		node := &label.Block.Tree[i]
		if node.Data == nil {
			continue
		}
		switch data := node.Data.(type) {
		case ast.Comment:
			continue
		case ast.Iter:
			label.Used = true
			cont.Label = data.NextLabel()
		default:
			p.pusherrtok(cont.LoopLabel, "invalid_label")
		}
		break
	}
}

func (p *Parser) break_st(ast *ast.Break) {
	switch {
	case ast.LabelToken.Id != lexer.ID_NA:
		p.breakWithLabel(ast)
	case p.currentCase != nil:
		ast.Label = p.currentCase.Match.EndLabel()
	case p.currentIter != nil:
		ast.Label = p.currentIter.EndLabel()
	default:
		p.pusherrtok(ast.Token, "break_at_out_of_valid_scope")
	}
}

func (p *Parser) continue_st(ast *ast.Continue) {
	switch {
	case p.currentIter == nil:
		p.pusherrtok(ast.Token, "continue_at_out_of_valid_scope")
	case ast.LoopLabel.Id != lexer.ID_NA:
		p.continueWithLabel(ast)
	default:
		ast.Label = p.currentIter.NextLabel()
	}
}

func (p *Parser) checkValidityForAutoType(expr_t Type, errtok lexer.Token) {
	if p.eval.has_error {
		return
	}
	switch expr_t.Id {
	case types.NIL:
		p.pusherrtok(errtok, "nil_for_autotype")
	case types.VOID:
		p.pusherrtok(errtok, "void_for_autotype")
	}
}

func (p *Parser) typeSourceOfMultiTyped(dt Type, err bool) (Type, bool) {
	types := dt.Tag.([]Type)
	ok := false
	for i, mt := range types {
		mt, ok = p.typeSource(mt, err)
		types[i] = mt
	}
	dt.Tag = types
	return dt, ok
}

func (p *Parser) typeSourceIsAlias(dt Type, alias *TypeAlias, err bool) (Type, bool) {
	original := dt.Original
	old := dt
	dt = alias.TargetType
	dt.Token = alias.Token
	dt.Generic = alias.Generic
	dt.Original = original
	dt, ok := p.typeSource(dt, err)
	dt.Pure = false
	if ok && old.Tag != nil && !types.IsStruct(alias.TargetType) {
		p.pusherrtok(dt.Token, "invalid_type_source")
	}
	return dt, ok
}

func (p *Parser) typeSourceIsEnum(e *Enum, tag any) (dt Type, _ bool) {
	dt.Id = types.ENUM
	dt.Kind = e.Id
	dt.Tag = e
	dt.Token = e.Token
	dt.Pure = true
	if tag != nil {
		p.pusherrtok(dt.Token, "invalid_type_source")
	}
	return dt, true
}

func (p *Parser) typeSourceIsFn(dt Type, err bool) (Type, bool) {
	f := dt.Tag.(*Fn)
	if len(f.Generics) > 0 {
		p.pusherrtok(dt.Token, "genericed_fn_as_anonymous_fn")
		return dt, false
	}
	p.reload_fn_types(f)
	dt.Kind = f.TypeKind()
	return dt, true
}

func (p *Parser) typeSourceIsMap(dt Type, err bool) (Type, bool) {
	types := dt.Tag.([]Type)
	key := &types[0]
	*key, _ = p.realType(*key, err)
	value := &types[1]
	*value, _ = p.realType(*value, err)
	dt.Kind = dt.MapKind()
	return dt, true
}

func (p *Parser) typeSourceIsStruct(s *Struct, st Type) (dt Type, _ bool) {
	generics := s.GetGenerics()
	if len(generics) > 0 {
		if !p.checkGenericsQuantity(len(s.Generics), len(generics), st.Token) {
			goto end
		}
		for i, g := range generics {
			var ok bool
			g, ok = p.realType(g, true)
			generics[i] = g
			if !ok {
				goto end
			}
		}
		*s.Constructor.Combines = append(*s.Constructor.Combines, generics)
		owner := s.Owner.(*Parser)
		blockTypes := owner.blockTypes
		owner.blockTypes = nil
		owner.pushGenerics(s.Generics, generics, st.Token)
		for i, f := range s.Fields {
			owner.parse_field(s, &f, i)
		}
		if len(s.Defines.Fns) > 0 {
			for _, f := range s.Defines.Fns {
				if len(f.Generics) == 0 {
					blockVars := owner.block_vars
					blockTypes := owner.blockTypes
					owner.reload_fn_types(f)
					_ = p.parse_pure_fn(f)
					owner.block_vars = blockVars
					owner.blockTypes = blockTypes
				}
			}
		}
		if owner != p {
			owner.wg.Wait()
			p.pusherrs(owner.Errors...)
			owner.Errors = nil
		}
		owner.blockTypes = blockTypes
	} else if len(s.Generics) > 0 {
		p.pusherrtok(st.Token, "has_generics")
	}
end:
	dt.Id = types.STRUCT
	dt.Kind = s.AsTypeKind()
	dt.Tag = s
	dt.Token = s.Token
	return dt, true
}

func (p *Parser) typeSourceIsTrait(
	trait_def *ast.Trait,
	tag any,
	errTok lexer.Token,
) (dt Type, _ bool) {
	if tag != nil {
		p.pusherrtok(errTok, "invalid_type_source")
	}
	trait_def.Used = true
	dt.Id = types.TRAIT
	dt.Kind = trait_def.Id
	dt.Tag = trait_def
	dt.Token = trait_def.Token
	dt.Pure = true
	return dt, true
}

func (p *Parser) tokenizeDataType(id string) []lexer.Token {
	parts := strings.SplitN(id, lexer.KND_DBLCOLON, -1)
	var toks []lexer.Token
	for i, part := range parts {
		toks = append(toks, lexer.Token{
			Id:   lexer.ID_IDENT,
			Kind: part,
			File: p.File,
		})
		if i < len(parts)-1 {
			toks = append(toks, lexer.Token{
				Id:   lexer.ID_DBLCOLON,
				Kind: lexer.KND_DBLCOLON,
				File: p.File,
			})
		}
	}
	return toks
}

func (p *Parser) typeSourceIsArrayType(arr_t *Type) (ok bool) {
	ok = true
	arr_t.Original = nil
	arr_t.Pure = true
	*arr_t.ComponentType, ok = p.realType(*arr_t.ComponentType, true)
	if !ok {
		return
	} else if types.IsArray(*arr_t.ComponentType) && arr_t.ComponentType.Size.AutoSized {
		p.pusherrtok(arr_t.Token, "invalid_type")
	}
	modifiers := arr_t.Modifiers()
	arr_t.Kind = modifiers + lexer.PREFIX_ARRAY + arr_t.ComponentType.Kind
	if arr_t.Size.AutoSized || arr_t.Size.Expr.Model != nil {
		return
	}
	val, model := p.eval_expr(arr_t.Size.Expr, nil)
	arr_t.Size.Expr.Model = model
	if val.constant {
		arr_t.Size.N = ast.Size(to_num_unsigned(val.expr))
	} else {
		p.eval.push_err_tok(arr_t.Token, "expr_not_const")
	}
	assign_checker{
		p:      p,
		t:      Type{Id: types.UINT, Kind: types.TYPE_MAP[types.UINT]},
		v:      val,
		errtok: arr_t.Size.Expr.Tokens[0],
	}.check()
	return
}

func (p *Parser) typeSourceIsSliceType(slc_t *Type) (ok bool) {
	*slc_t.ComponentType, ok = p.realType(*slc_t.ComponentType, true)
	if ok && types.IsArray(*slc_t.ComponentType) && slc_t.ComponentType.Size.AutoSized {
		p.pusherrtok(slc_t.Token, "invalid_type")
	}
	modifiers := slc_t.Modifiers()
	slc_t.Kind = modifiers + lexer.PREFIX_SLICE + slc_t.ComponentType.Kind
	return
}

func (p *Parser) check_type_validity(expr_t Type, errtok lexer.Token) {
	modifiers := expr_t.Modifiers()
	if strings.Contains(modifiers, "&&") ||
		(strings.Contains(modifiers, "*") && strings.Contains(modifiers, "&")) {
		p.pusherrtok(expr_t.Token, "invalid_type")
		return
	}
	if types.IsRef(expr_t) && !types.ValidForRef(types.Elem(expr_t)) {
		p.pusherrtok(errtok, "invalid_type")
		return
	}
	if expr_t.Id == types.UNSAFE {
		n := len(expr_t.Kind) - len(lexer.KND_UNSAFE) - 1
		if n < 0 || expr_t.Kind[n] != '*' {
			p.pusherrtok(errtok, "invalid_type")
		}
	}
}

func (p *Parser) get_define(id string, cpp_linked bool) any {
	var def any = nil
	if cpp_linked {
		def, _ = p.linkById(id)
	} else if strings.Contains(id, lexer.KND_DBLCOLON) {
		toks := p.tokenizeDataType(id)
		defs := p.eval.get_ns(&toks)
		if defs == nil {
			return nil
		}
		i, m, def_t := defs.FindById(toks[0].Kind, p.File)
		switch def_t {
		case 't':
			def = m.Types[i]
		case 's':
			def = m.Structs[i]
		case 'e':
			def = m.Enums[i]
		case 'i':
			def = m.Traits[i]
		}
	} else {
		def, _, _ = p.defined_by_id(id)
	}
	return def
}

func (p *Parser) typeSource(dt Type, err bool) (ret Type, ok bool) {
	if dt.Kind == "" {
		return dt, true
	}
	original := dt.Original
	defer func() {
		ret.CppLinked = (original != nil && original.(Type).CppLinked) || dt.CppLinked
		ret.Original = original
		p.check_type_validity(ret, dt.Token)
	}()
	dt.SetToOriginal()
	switch {
	case dt.MultiTyped:
		return p.typeSourceOfMultiTyped(dt, err)
	case dt.Id == types.MAP:
		return p.typeSourceIsMap(dt, err)
	case dt.Id == types.ARRAY:
		ok = p.typeSourceIsArrayType(&dt)
		return dt, ok
	case dt.Id == types.SLICE:
		ok = p.typeSourceIsSliceType(&dt)
		return dt, ok
	}
	switch dt.Id {
	case types.STRUCT:
		_, prefix := dt.KindId()
		ret, ok = p.typeSourceIsStruct(dt.Tag.(*Struct), dt)
		ret.Kind = prefix + ret.Kind
		return
	case types.ID:
		id, prefix := dt.KindId()
		defer func() { ret.Kind = prefix + ret.Kind }()
		def := p.get_define(id, dt.CppLinked)
		switch def := def.(type) {
		case *TypeAlias:
			def.Used = true
			return p.typeSourceIsAlias(dt, def, err)
		case *Enum:
			def.Used = true
			return p.typeSourceIsEnum(def, dt.Tag)
		case *Struct:
			def.Used = true
			def = p.structConstructorInstance(def)
			switch tagt := dt.Tag.(type) {
			case []ast.Type:
				def.SetGenerics(tagt)
			}
			return p.typeSourceIsStruct(def, dt)
		case *ast.Trait:
			def.Used = true
			return p.typeSourceIsTrait(def, dt.Tag, dt.Token)
		default:
			if err {
				p.pusherrtok(dt.Token, "invalid_type_source")
			}
			return dt, false
		}
	case types.FN:
		return p.typeSourceIsFn(dt, err)
	}
	return dt, true
}

func (p *Parser) realType(dt Type, err bool) (ret Type, _ bool) {
	original := dt.Original
	defer func() {
		ret.CppLinked = (original != nil && original.(Type).CppLinked) || dt.CppLinked
		ret.Original = original
	}()
	dt.SetToOriginal()
	return p.typeSource(dt, err)
}

func (p *Parser) checkMultiType(real, check Type, ignoreAny bool, errTok lexer.Token) {
	if real.MultiTyped != check.MultiTyped {
		p.pusherrtok(errTok, "incompatible_types", real.Kind, check.Kind)
		return
	}
	realTypes := real.Tag.([]Type)
	checkTypes := real.Tag.([]Type)
	if len(realTypes) != len(checkTypes) {
		p.pusherrtok(errTok, "incompatible_types", real.Kind, check.Kind)
		return
	}
	for i := 0; i < len(realTypes); i++ {
		realType := realTypes[i]
		checkType := checkTypes[i]
		p.check_type(realType, checkType, ignoreAny, true, errTok)
	}
}

func (p *Parser) check_type(real, check Type, ignoreAny, allow_assign bool, errTok lexer.Token) {
	if types.IsVoid(check) {
		p.eval.push_err_tok(errTok, "incompatible_types", real.Kind, check.Kind)
		return
	}
	if !ignoreAny && real.Id == types.ANY {
		return
	}
	if real.MultiTyped || check.MultiTyped {
		p.checkMultiType(real, check, ignoreAny, errTok)
		return
	}
	checker := types.Checker{
		ErrTok:      errTok,
		L:           real,
		R:           check,
		IgnoreAny:   ignoreAny,
		AllowAssign: allow_assign,
	}
	ok := checker.Check()
	if ok || checker.ErrorLogged {
		p.pusherrs(checker.Errors...)
		return
	}
	if real.Kind != check.Kind {
		p.pusherrtok(errTok, "incompatible_types", real.Kind, check.Kind)
	} else if types.IsArray(real) || types.IsArray(check) {
		if types.IsArray(real) != types.IsArray(check) {
			p.pusherrtok(errTok, "incompatible_types", real.Kind, check.Kind)
			return
		}
		realKind := strings.Replace(real.Kind, lexer.MARK_ARRAY, strconv.Itoa(real.Size.N), 1)
		checkKind := strings.Replace(check.Kind, lexer.MARK_ARRAY, strconv.Itoa(check.Size.N), 1)
		p.pusherrtok(errTok, "incompatible_types", realKind, checkKind)
	}
}

func (p *Parser) eval_expr(expr Expr, prefix *ast.Type) (value, ast.ExprModel) {
	p.eval.has_error = false
	p.eval.type_prefix = prefix
	return p.eval.eval_expr(expr)
}

func (p *Parser) evalToks(toks []lexer.Token, prefix *ast.Type) (value, ast.ExprModel) {
	p.eval.has_error = false
	p.eval.type_prefix = prefix
	return p.eval.eval_toks(toks)
}
