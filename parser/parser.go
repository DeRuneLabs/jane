package parser

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/De-Rune/jane/ast"
	"github.com/De-Rune/jane/lexer"
	"github.com/De-Rune/jane/package/jn"
	"github.com/De-Rune/jane/package/jnapi"
	"github.com/De-Rune/jane/package/jnbits"
	"github.com/De-Rune/jane/package/jnio"
	"github.com/De-Rune/jane/package/jnlog"
	"github.com/De-Rune/jane/preprocessor"
)

type use struct {
	Path string
	defs *defmap
}

var used []*use

type Parser struct {
	attributes     []ast.Attribute
	docText        strings.Builder
	iterCount      int
	wg             sync.WaitGroup
	justDefs       bool
	main           bool
	isLocalPackage bool

	Embeds         strings.Builder
	Uses           []*use
	Defs           *defmap
	waitingGlobals []ast.Var
	BlockVars      []ast.Var
	Errs           []jnlog.CompilerLog
	Warns          []jnlog.CompilerLog
	File           *jnio.File
}

func New(f *jnio.File) *Parser {
	p := new(Parser)
	p.File = f
	p.isLocalPackage = false
	p.Defs = new(defmap)
	return p
}

func Parset(tree []ast.Obj, main, justDefs bool) *Parser {
	p := New(nil)
	p.Parset(tree, main, justDefs)
	return p
}

func (p *Parser) pusherrtok(token lexer.Token, key string, args ...interface{}) {
	p.pusherrmsgtok(token, jn.GetErr(key, args...))
}

func (p *Parser) pusherrmsgtok(token lexer.Token, msg string) {
	p.Errs = append(p.Errs, jnlog.CompilerLog{
		Type:   jnlog.Err,
		Row:    token.Row,
		Column: token.Column,
		Path:   token.File.Path,
		Msg:    msg,
	})
}

func (p *Parser) pushwarntok(token lexer.Token, key string, args ...interface{}) {
	p.Warns = append(p.Warns, jnlog.CompilerLog{
		Type:   jnlog.Warn,
		Row:    token.Row,
		Column: token.Column,
		Path:   token.File.Path,
		Msg:    jn.GetWarn(key, args...),
	})
}

func (p *Parser) pusherrs(errs ...jnlog.CompilerLog) {
	p.Errs = append(p.Errs, errs...)
}

func (p *Parser) pusherr(key string, args ...interface{}) {
	p.pusherrmsg(jn.GetErr(key, args...))
}

func (p *Parser) pusherrmsg(msg string) {
	p.Errs = append(p.Errs, jnlog.CompilerLog{
		Type: jnlog.FlatErr,
		Msg:  msg,
	})
}

func (p *Parser) pushwarn(key string, args ...interface{}) {
	p.Warns = append(p.Warns, jnlog.CompilerLog{
		Type: jnlog.FlatWarn,
		Msg:  jn.GetWarn(key, args...),
	})
}

func (p *Parser) CxxEmbeds() string {
	var cxx strings.Builder
	cxx.WriteString("// region EMBEDS\n")
	cxx.WriteString(p.Embeds.String())
	cxx.WriteString("// endregion EMBEDS")
	return cxx.String()
}

func (p *Parser) CxxPrototypes() string {
	var cxx strings.Builder
	cxx.WriteString("// region PROTOTYPES\n")
	for _, use := range used {
		for _, f := range use.defs.Funcs {
			cxx.WriteString(f.Prototype())
			cxx.WriteByte('\n')
		}
	}
	for _, f := range p.Defs.Funcs {
		cxx.WriteString(f.Prototype())
		cxx.WriteByte('\n')
	}
	cxx.WriteString("// endregion PROTOTYPES")
	return cxx.String()
}

func (p *Parser) CxxGlobals() string {
	var cxx strings.Builder
	cxx.WriteString("// region GLOBALS\n")
	for _, use := range used {
		for _, v := range use.defs.Globals {
			cxx.WriteString(v.String())
			cxx.WriteByte('\n')
		}
	}
	for _, v := range p.Defs.Globals {
		cxx.WriteString(v.String())
		cxx.WriteByte('\n')
	}
	cxx.WriteString("// endregion GLOBALS")
	return cxx.String()
}

func (p *Parser) CxxFuncs() string {
	var cxx strings.Builder
	cxx.WriteString("// region FUNCTIONS\n")
	for _, use := range used {
		for _, f := range use.defs.Funcs {
			cxx.WriteString(f.String())
			cxx.WriteString("\n\n")
		}
	}
	for _, f := range p.Defs.Funcs {
		cxx.WriteString(f.String())
		cxx.WriteString("\n\n")
	}
	cxx.WriteString("// endregion FUNCTIONS")
	return cxx.String()
}

func (p *Parser) Cxx() string {
	var cxx strings.Builder
	cxx.WriteString(p.CxxEmbeds())
	cxx.WriteString("\n\n")
	cxx.WriteString(p.CxxPrototypes())
	cxx.WriteString("\n\n")
	cxx.WriteString(p.CxxGlobals())
	cxx.WriteString("\n\n")
	cxx.WriteString(p.CxxFuncs())
	return cxx.String()
}

func getTree(tokens []lexer.Token, errs *[]jnlog.CompilerLog) []ast.Obj {
	b := ast.NewBuilder(tokens)
	b.Build()
	if len(b.Errs) > 0 {
		if errs != nil {
			*errs = append(*errs, b.Errs...)
		}
		return nil
	}
	return b.Tree
}

func (p *Parser) checkUsePath(use *ast.Use) bool {
	info, err := os.Stat(use.Path)
	if err != nil || !info.IsDir() {
		p.pusherrtok(use.Token, "use_not_found", use.Path)
		return false
	}

	for _, puse := range p.Uses {
		if use.Path == puse.Path {
			p.pusherrtok(use.Token, "already_uses")
			return false
		}
	}
	return true
}

func (p *Parser) compileUse(useAST *ast.Use) *use {
	infos, err := ioutil.ReadDir(useAST.Path)
	if err != nil {
		p.pusherrmsg(err.Error())
		return nil
	}
	for _, info := range infos {
		name := info.Name()
		if info.IsDir() || !strings.HasSuffix(name, jn.SrcExt) {
			continue
		}
		f, err := jnio.OpenJn(filepath.Join(useAST.Path, name))
		if err != nil {
			p.pusherrmsg(err.Error())
			continue
		}
		psub := New(f)
		psub.Parsef(false, false)
		if psub.Errs != nil {
			p.pusherrtok(useAST.Token, "use_has_errors")
		}
		use := new(use)
		use.defs = new(defmap)
		use.Path = useAST.Path
		p.pusherrs(psub.Errs...)
		p.Warns = append(p.Warns, psub.Warns...)
		p.pushUseDefs(use, psub.Defs)
		return use
	}
	return nil
}

func (p *Parser) pushUseTypes(use *use, dm *defmap) {
	for _, t := range dm.Types {
		def := p.typeById(t.Id)
		if def != nil {
			p.pusherrmsgtok(def.Token,
				fmt.Sprintf(`"%s" identifier is already defined in this source`, t.Id))
		} else {
			use.defs.Types = append(use.defs.Types, t)
		}
	}
}

func (p *Parser) pushUseGlobals(use *use, dm *defmap) {
	for _, g := range dm.Globals {
		def := p.Defs.globalById(g.Id)
		if def != nil {
			p.pusherrmsgtok(def.IdToken,
				fmt.Sprintf(`"%s" identifier is already defined in this source`, g.Id))
		} else {
			use.defs.Globals = append(use.defs.Globals, g)
		}
	}
}

func (p *Parser) pushUseFuncs(use *use, dm *defmap) {
	for _, f := range dm.Funcs {
		def := p.Defs.funcById(f.Ast.Id)
		if def != nil {
			p.pusherrmsgtok(def.Ast.Token,
				fmt.Sprintf(`"%s" identifier is already defined in this source`, f.Ast.Id))
		} else {
			use.defs.Funcs = append(use.defs.Funcs, f)
		}
	}
}

func (p *Parser) pushUseDefs(use *use, dm *defmap) {
	p.pushUseTypes(use, dm)
	p.pushUseGlobals(use, dm)
	p.pushUseFuncs(use, dm)
}

func (p *Parser) use(useAST *ast.Use) {
	if !p.checkUsePath(useAST) {
		return
	}
	for _, use := range used {
		if useAST.Path == use.Path {
			p.Uses = append(p.Uses, use)
			return
		}
	}
	use := p.compileUse(useAST)
	if use == nil {
		return
	}
	exist := false
	for _, guse := range used {
		if guse.Path == use.Path {
			exist = true
			break
		}
	}
	if !exist {
		used = append(used, use)
	}
	p.Uses = append(p.Uses, use)
}

func (p *Parser) parseUses(tree *[]ast.Obj) {
	for i, obj := range *tree {
		switch t := obj.Value.(type) {
		case ast.Use:
			p.use(&t)
		default:
			*tree = (*tree)[i:]
			return
		}
	}
	*tree = nil
}

func (p *Parser) parseSrcTreeObj(obj ast.Obj) {
	switch t := obj.Value.(type) {
	case ast.Attribute:
		p.PushAttribute(t)
	case ast.Statement:
		p.Statement(t)
	case ast.Type:
		p.Type(t)
	case ast.CxxEmbed:
		p.Embeds.WriteString(t.String())
		p.Embeds.WriteByte('\n')
	case ast.Comment:
		p.Comment(t)
	case ast.Use:
		p.pusherrtok(obj.Token, "use_at_content")
	case ast.Preprocessor:
	default:
		p.pusherrtok(obj.Token, "invalid_syntax")
	}
}

func (p *Parser) parseSrcTree(tree []ast.Obj) {
	for _, obj := range tree {
		p.parseSrcTreeObj(obj)
		p.checkDoc(obj)
		p.checkAttribute(obj)
	}
}

func (p *Parser) parseTree(tree []ast.Obj) {
	p.parseUses(&tree)
	p.parseSrcTree(tree)
}

func (p *Parser) checkParse() {
	if p.docText.Len() > 0 {
		p.pushwarn("exist_undefined_doc")
	}
	p.wg.Add(1)
	go p.checkAsync()
}

func (p *Parser) useLocalPackage(tree *[]ast.Obj) {
	if p.File == nil {
		return
	}
	dir := filepath.Dir(p.File.Path)
	infos, err := ioutil.ReadDir(dir)
	if err != nil {
		p.pusherrmsg(err.Error())
		return
	}
	_, mainName := filepath.Split(p.File.Path)
	for _, info := range infos {
		name := info.Name()
		if info.IsDir() ||
			!strings.HasSuffix(name, jn.SrcExt) ||
			name == mainName {
			continue
		}
		f, err := jnio.OpenJn(filepath.Join(dir, name))
		if err != nil {
			p.pusherrmsg(err.Error())
			continue
		}
		lexer := lexer.NewLexer(f)
		toks := lexer.Lexer()
		if lexer.Logs != nil {
			p.pusherrs(lexer.Logs...)
			continue
		}
		subtree := getTree(toks, &p.Errs)
		if subtree == nil {
			continue
		}
		preprocessor.TrimEnofi(&subtree)
		p.parseUses(&subtree)
		*tree = append(*tree, subtree...)
	}
}

func (p *Parser) Parset(tree []ast.Obj, main, justDefs bool) {
	p.main = main
	p.justDefs = justDefs
	if !p.isLocalPackage {
		p.useLocalPackage(&tree)
	}
	if !main {
		preprocessor.TrimEnofi(&tree)
	}
	p.parseTree(tree)
	p.checkParse()
	p.wg.Wait()
}

func (p *Parser) Parse(tokens []lexer.Token, main, justDefs bool) {
	tree := getTree(tokens, &p.Errs)
	if tree == nil {
		return
	}
	p.Parset(tree, main, justDefs)
}

func (p *Parser) Parsef(main, justDefs bool) {
	lexer := lexer.NewLexer(p.File)
	toks := lexer.Lexer()
	if lexer.Logs != nil {
		p.pusherrs(lexer.Logs...)
		return
	}
	p.Parse(toks, main, justDefs)
}

func (p *Parser) checkDoc(obj ast.Obj) {
	if p.docText.Len() == 0 {
		return
	}
	switch obj.Value.(type) {
	case ast.Comment, ast.Attribute:
		return
	}
	p.pushwarntok(obj.Token, "doc_ignored")
	p.docText.Reset()
}

func (p *Parser) checkAttribute(obj ast.Obj) {
	if p.attributes == nil {
		return
	}
	switch obj.Value.(type) {
	case ast.Attribute, ast.Comment:
		return
	}
	p.pusherrtok(obj.Token, "attribute_not_supports")
	p.attributes = nil
}

func (p *Parser) Type(t ast.Type) {
	if p.existid(t.Id).Id != lexer.NA {
		p.pusherrtok(t.Token, "exist_id", t.Id)
		return
	} else if jnapi.IsIgnoreId(t.Id) {
		p.pusherrtok(t.Token, "ignore_id")
		return
	}
	t.Desc = p.docText.String()
	p.docText.Reset()
	p.Defs.Types = append(p.Defs.Types, t)
}

func (p *Parser) Comment(c ast.Comment) {
	c.Content = strings.TrimSpace(c.Content)
	if p.docText.Len() == 0 {
		if strings.HasPrefix(c.Content, "doc:") {
			c.Content = c.Content[4:]
			if c.Content == "" {
				c.Content = " "
			}
			goto write
		}
		return
	}
	p.docText.WriteByte('\n')
write:
	p.docText.WriteString(c.Content)
}

func (p *Parser) PushAttribute(attribute ast.Attribute) {
	switch attribute.Tag.Kind {
	case "inline":
	default:
		p.pusherrtok(attribute.Tag, "undefined_tag")
	}
	for _, attr := range p.attributes {
		if attr.Tag.Kind == attribute.Tag.Kind {
			p.pusherrtok(attribute.Tag, "attribute_repeat")
			return
		}
	}
	p.attributes = append(p.attributes, attribute)
}

func (p *Parser) Statement(s ast.Statement) {
	switch t := s.Val.(type) {
	case ast.Func:
		p.Func(t)
	case ast.Var:
		p.Global(t)
	default:
		p.pusherrtok(s.Token, "invalid_syntax")
	}
}

func (p *Parser) Func(fast ast.Func) {
	if p.existid(fast.Id).Id != lexer.NA {
		p.pusherrtok(fast.Token, "exist_id", fast.Id)
	} else if jnapi.IsIgnoreId(fast.Id) {
		p.pusherrtok(fast.Token, "ignore_id")
	}
	fast.RetType, _ = p.readyType(fast.RetType, true)
	for i, param := range fast.Params {
		fast.Params[i].Type, _ = p.readyType(param.Type, true)
	}
	f := new(function)
	f.Ast = fast
	f.Attributes = p.attributes
	f.Desc = p.docText.String()
	p.attributes = nil
	p.docText.Reset()
	p.checkFuncAttributes(f.Attributes)
	p.Defs.Funcs = append(p.Defs.Funcs, f)
}

func (p *Parser) Global(vast ast.Var) {
	if p.existid(vast.Id).Id != lexer.NA {
		p.pusherrtok(vast.IdToken, "exist_id", vast.Id)
		return
	}
	vast.Desc = p.docText.String()
	p.docText.Reset()
	p.waitingGlobals = append(p.waitingGlobals, vast)
}

func (p *Parser) Var(vast ast.Var) ast.Var {
	if jnapi.IsIgnoreId(vast.Id) {
		p.pusherrtok(vast.IdToken, "ignore_id")
	}
	var val value
	switch t := vast.Tag.(type) {
	case value:
		val = t
	default:
		if vast.SetterToken.Id != lexer.NA {
			val, vast.Val.Model = p.evalExpr(vast.Val)
		}
	}
	if vast.Type.Id != jn.Void {
		if vast.SetterToken.Id != lexer.NA {
			p.wg.Add(1)
			go assignChecker{
				p,
				vast.Const,
				vast.Type,
				val,
				false,
				vast.IdToken,
			}.checkAssignTypeAsync()
		} else {
			dt, ok := p.readyType(vast.Type, true)
			if ok {
				var valTok lexer.Token
				valTok.Id = lexer.Value
				valTok.Kind = p.defaultValueOfType(dt)
				valToks := []lexer.Token{valTok}
				processes := [][]lexer.Token{valToks}
				vast.Val = ast.Expr{Tokens: valToks, Processes: processes}
				_, vast.Val.Model = p.evalExpr(vast.Val)
			}
		}
	} else {
		if vast.SetterToken.Id == lexer.NA {
			p.pusherrtok(vast.IdToken, "missing_autotype_value")
		} else {
			vast.Type = val.ast.Type
			p.checkValidityForAutoType(vast.Type, vast.SetterToken)
			p.checkAssignConst(vast.Const, vast.Type, val, vast.SetterToken)
		}
	}
	if vast.Const {
		if vast.SetterToken.Id == lexer.NA {
			p.pusherrtok(vast.IdToken, "missing_const_value")
		}
	}
	return vast
}

func (p *Parser) checkFuncAttributes(attributes []ast.Attribute) {
	for _, attribute := range attributes {
		switch attribute.Tag.Kind {
		case "inline":
		default:
			p.pusherrtok(attribute.Token, "invalid_attribute")
		}
	}
}

func (p *Parser) varsFromParams(params []ast.Parameter) []ast.Var {
	var vars []ast.Var
	length := len(params)
	for i, param := range params {
		var vast ast.Var
		vast.Id = param.Id
		vast.IdToken = param.Token
		vast.Type = param.Type
		vast.Const = param.Const
		vast.Volatile = param.Volatile
		if param.Variadic {
			if length-i > 1 {
				p.pusherrtok(param.Token, "variadic_parameter_notlast")
			}
			vast.Type.Val = "[]" + vast.Type.Val
		}
		vars = append(vars, vast)
	}
	return vars
}

func (p *Parser) FuncById(id string) *function {
	for _, f := range builtinFuncs {
		if f.Ast.Id == id {
			return f
		}
	}
	for _, use := range p.Uses {
		f := use.defs.funcById(id)
		if f != nil && f.Ast.Pub {
			return f
		}
	}
	return p.Defs.funcById(id)
}

func (p *Parser) varById(id string) *ast.Var {
	for _, v := range p.BlockVars {
		if v.Id == id {
			return &v
		}
	}
	return p.globalById(id)
}

func (p *Parser) globalById(id string) *ast.Var {
	for _, use := range p.Uses {
		g := use.defs.globalById(id)
		if g != nil && g.Pub {
			return g
		}
	}
	return p.Defs.globalById(id)
}

func (p *Parser) typeById(id string) *ast.Type {
	for _, use := range p.Uses {
		t := use.defs.typeById(id)
		if t != nil && t.Pub {
			return t
		}
	}
	return p.Defs.typeById(id)
}

func (p *Parser) existIdf(id string, exceptGlobals bool) lexer.Token {
	t := p.typeById(id)
	if t != nil {
		return t.Token
	}
	f := p.FuncById(id)
	if f != nil {
		return f.Ast.Token
	}
	for _, v := range p.BlockVars {
		if v.Id == id {
			return v.IdToken
		}
	}
	if !exceptGlobals {
		v := p.globalById(id)
		if v != nil {
			return v.IdToken
		}
		for _, v := range p.waitingGlobals {
			if v.Id == id {
				return v.IdToken
			}
		}
	}
	return lexer.Token{}
}

func (p *Parser) existid(id string) lexer.Token {
	return p.existIdf(id, false)
}

func (p *Parser) checkAsync() {
	defer func() { p.wg.Done() }()
	if p.main && !p.justDefs {
		if p.FuncById(jn.EntryPoint) == nil {
			p.pusherr("no_entry_point")
		}
	}
	p.wg.Add(1)
	go p.checkTypesAsync()
	p.WaitingGlobals()
	p.waitingGlobals = nil
	if !p.justDefs {
		p.wg.Add(1)
		go p.checkFuncsAsync()
	}
}

func (p *Parser) checkTypesAsync() {
	defer func() { p.wg.Done() }()
	for _, t := range p.Defs.Types {
		_, _ = p.readyType(t.Type, true)
	}
}

func (p *Parser) WaitingGlobals() {
	for _, varAST := range p.waitingGlobals {
		variable := p.Var(varAST)
		p.Defs.Globals = append(p.Defs.Globals, variable)
	}
}

func (p *Parser) checkFuncsAsync() {
	defer func() { p.wg.Done() }()
	for _, f := range p.Defs.Funcs {
		p.BlockVars = p.varsFromParams(f.Ast.Params)
		p.wg.Add(1)
		go p.checkFuncSpecialCasesAsync(f)
		p.checkFunc(&f.Ast)
	}
}

func (p *Parser) checkFuncSpecialCasesAsync(fun *function) {
	defer func() { p.wg.Done() }()
	switch fun.Ast.Id {
	case jn.EntryPoint:
		p.checkEntryPointSpecialCases(fun)
	}
}

type value struct {
	ast      ast.Value
	constant bool
	volatile bool
	lvalue   bool
	variadic bool
}

func eliminateProcesses(processes *[][]lexer.Token, i, to int) {
	for i < to {
		(*processes)[i] = nil
		i++
	}
}

func (p *Parser) evalProcesses(processes [][]lexer.Token) (v value, e iExpr) {
	if processes == nil {
		return
	}
	m := newExprModel(processes)
	e = m
	if len(processes) == 1 {
		v = p.evalExprPart(processes[0], m)
		return
	}
	process := solver{p: p, model: m}
	boolean := false
	for i := p.nextOperator(processes); i != -1 && !noData(processes); i = p.nextOperator(processes) {
		if !boolean {
			boolean = v.ast.Type.Id == jn.Bool
		}
		if boolean {
			v.ast.Type.Id = jn.Bool
		}
		m.index = i
		process.operator = processes[m.index][0]
		m.appendSubNode(exprNode{process.operator.Kind})
		if processes[i-1] == nil {
			process.leftVal = v.ast
			m.index = i + 1
			process.right = processes[m.index]
			process.rightVal = p.evalExprPart(process.right, m).ast
			v.ast = process.Solve()
			eliminateProcesses(&processes, i, i+2)
			continue
		} else if processes[i+1] == nil {
			m.index = i - 1
			process.left = processes[m.index]
			process.leftVal = p.evalExprPart(process.left, m).ast
			process.rightVal = v.ast
			v.ast = process.Solve()
			eliminateProcesses(&processes, i-1, i+1)
			continue
		} else if isOperator(processes[i-1]) {
			process.leftVal = v.ast
			m.index = i + 1
			process.right = processes[m.index]
			process.rightVal = p.evalExprPart(process.right, m).ast
			v.ast = process.Solve()
			eliminateProcesses(&processes, i, i+1)
			continue
		}
		m.index = i - 1
		process.left = processes[m.index]
		process.leftVal = p.evalExprPart(process.left, m).ast
		m.index = i + 1
		process.right = processes[m.index]
		process.rightVal = p.evalExprPart(process.right, m).ast
		solvedv := process.Solve()
		if v.ast.Type.Id != jn.Void {
			process.operator.Kind = "+"
			process.leftVal = v.ast
			process.right = processes[i+1]
			process.rightVal = solvedv
			solvedv = process.Solve()
		}
		v.ast = solvedv
		eliminateProcesses(&processes, i-1, i+2)
	}
	return
}

func noData(processes [][]lexer.Token) bool {
	for _, p := range processes {
		if !isOperator(p) && p != nil {
			return false
		}
	}
	return true
}

func isOperator(process []lexer.Token) bool {
	return len(process) == 1 && process[0].Id == lexer.Operator
}

func (p *Parser) nextOperator(processes [][]lexer.Token) int {
	precedence5 := -1
	precedence4 := -1
	precedence3 := -1
	precedence2 := -1
	precedence1 := -1
	for i, process := range processes {
		if !isOperator(process) {
			continue
		}
		if processes[i-1] == nil && processes[i+1] == nil {
			continue
		}
		switch process[0].Kind {
		case "*", "/", "%", "<<", ">>", "&":
			precedence5 = i
		case "+", "-", "|", "^":
			precedence4 = i
		case "==", "!=", "<", "<=", ">", ">=":
			precedence3 = i
		case "&&":
			precedence2 = i
		case "||":
			precedence1 = i
		default:
			p.pusherrtok(process[0], "invalid_operator")
		}
	}
	switch {
	case precedence5 != -1:
		return precedence5
	case precedence4 != -1:
		return precedence4
	case precedence3 != -1:
		return precedence3
	case precedence2 != -1:
		return precedence2
	default:
		return precedence1
	}
}

func (p *Parser) evalToks(tokens []lexer.Token) (value, iExpr) {
	return p.evalExpr(new(ast.Builder).Expr(tokens))
}

func (p *Parser) evalExpr(ex ast.Expr) (value, iExpr) {
	processes := make([][]lexer.Token, len(ex.Processes))
	copy(processes, ex.Processes)
	return p.evalProcesses(processes)
}

func toRawStrLiteral(literal string) string {
	literal = literal[1 : len(literal)-1]
	literal = `"(` + literal + `)"`
	literal = jnapi.ToRawStr(literal)
	return literal
}

type valueEvaluator struct {
	token  lexer.Token
	model  *exprModel
	parser *Parser
}

func (p *valueEvaluator) str() value {
	var v value
	v.ast.Data = p.token.Kind
	v.ast.Type.Id = jn.Str
	v.ast.Type.Val = "str"
	if israwstr(p.token.Kind) {
		p.model.appendSubNode(exprNode{toRawStrLiteral(p.token.Kind)})
	} else {
		p.model.appendSubNode(exprNode{jnapi.ToStr(p.token.Kind)})
	}
	return v
}

func (p *valueEvaluator) rune() value {
	var v value
	v.ast.Data = p.token.Kind
	v.ast.Type.Id = jn.Rune
	v.ast.Type.Val = "rune"
	p.model.appendSubNode(exprNode{jnapi.ToRune(p.token.Kind)})
	return v
}

func (p *valueEvaluator) bool() value {
	var v value
	v.ast.Data = p.token.Kind
	v.ast.Type.Id = jn.Bool
	v.ast.Type.Val = "bool"
	p.model.appendSubNode(exprNode{p.token.Kind})
	return v
}

func (p *valueEvaluator) nil() value {
	var v value
	v.ast.Data = p.token.Kind
	v.ast.Type.Id = jn.Nil
	p.model.appendSubNode(exprNode{p.token.Kind})
	return v
}

func (p *valueEvaluator) num() value {
	var v value
	v.ast.Data = p.token.Kind
	p.model.appendSubNode(exprNode{p.token.Kind})
	if strings.Contains(p.token.Kind, ".") ||
		strings.ContainsAny(p.token.Kind, "eE") {
		v.ast.Type.Id = jn.F64
		v.ast.Type.Val = "f64"
	} else {
		v.ast.Type.Id = jn.I32
		v.ast.Type.Val = "i32"
		ok := jnbits.CheckBitInt(p.token.Kind, 32)
		if !ok {
			v.ast.Type.Id = jn.I64
			v.ast.Type.Val = "i64"
		}
	}
	return v
}

func (p *valueEvaluator) id() (v value, ok bool) {
	id := p.token.Kind
	if variable := p.parser.varById(id); variable != nil {
		v.ast.Data = id
		v.ast.Type = variable.Type
		v.constant = variable.Const
		v.volatile = variable.Volatile
		v.ast.Token = variable.IdToken
		v.lvalue = true
		p.model.appendSubNode(exprNode{jnapi.AsId(id)})
		ok = true
	} else if fun := p.parser.FuncById(id); fun != nil {
		v.ast.Data = id
		v.ast.Type.Id = jn.Func
		v.ast.Type.Tag = fun.Ast
		v.ast.Type.Val = fun.Ast.DataTypeString()
		v.ast.Token = fun.Ast.Token
		p.model.appendSubNode(exprNode{jnapi.AsId(id)})
		ok = true
	} else {
		p.parser.pusherrtok(p.token, "id_noexist", id)
	}
	return
}

type solver struct {
	p        *Parser
	left     []lexer.Token
	leftVal  ast.Value
	right    []lexer.Token
	rightVal ast.Value
	operator lexer.Token
	model    *exprModel
}

func (s solver) ptr() (v ast.Value) {
	ok := false
	switch {
	case s.leftVal.Type.Val == s.rightVal.Type.Val:
		ok = true
	case typeIsSingle(s.leftVal.Type):
		switch {
		case s.leftVal.Type.Id == jn.Nil,
			jn.IsIntegerType(s.leftVal.Type.Id):
			ok = true
		}
	case typeIsSingle(s.rightVal.Type):
		switch {
		case s.rightVal.Type.Id == jn.Nil,
			jn.IsIntegerType(s.rightVal.Type.Id):
			ok = true
		}
	}
	if !ok {
		s.p.pusherrtok(s.operator, "incompatible_datatype",
			s.rightVal.Type.Val, s.leftVal.Type.Val)
		return
	}
	switch s.operator.Kind {
	case "+", "-":
		if typeIsPtr(s.leftVal.Type) && typeIsPtr(s.rightVal.Type) {
			s.p.pusherrtok(s.operator, "incompatible_datatype",
				s.rightVal.Type.Val, s.leftVal.Type.Val)
			return
		}
		if typeIsPtr(s.leftVal.Type) {
			v.Type = s.leftVal.Type
		} else {
			v.Type = s.rightVal.Type
		}
	case "!=", "==":
		v.Type.Id = jn.Bool
	default:
		s.p.pusherrtok(s.operator, "operator_notfor_pointer")
	}
	return
}

func (s solver) str() (v ast.Value) {
	if s.leftVal.Type.Id != s.rightVal.Type.Id {
		s.p.pusherrtok(s.operator, "incompatible_datatype",
			s.leftVal.Type.Val, s.rightVal.Type.Val)
		return
	}
	switch s.operator.Kind {
	case "+":
		v.Type.Id = jn.Str
	case "==", "!=":
		v.Type.Id = jn.Bool
	default:
		s.p.pusherrtok(s.operator, "operator_notfor_string")
	}
	return
}

func (s solver) any() (v ast.Value) {
	switch s.operator.Kind {
	case "!=", "==":
		v.Type.Id = jn.Bool
	default:
		s.p.pusherrtok(s.operator, "operator_notfor_any")
	}
	return
}

func (s solver) bool() (v ast.Value) {
	if !typesAreCompatible(s.leftVal.Type, s.rightVal.Type, true) {
		s.p.pusherrtok(s.operator, "incompatible_datatype",
			s.rightVal.Type.Val, s.leftVal.Type.Val)
		return
	}
	switch s.operator.Kind {
	case "!=", "==":
		v.Type.Id = jn.Bool
	default:
		s.p.pusherrtok(s.operator, "operator_notfor_bool")
	}
	return
}

func (s solver) float() (v ast.Value) {
	if !typesAreCompatible(s.leftVal.Type, s.rightVal.Type, true) {
		if !isConstNum(s.leftVal.Data) &&
			!isConstNum(s.rightVal.Data) {
			s.p.pusherrtok(s.operator, "incompatible_datatype",
				s.rightVal.Type.Val, s.leftVal.Type.Val)
			return
		}
	}
	switch s.operator.Kind {
	case "!=", "==", "<", ">", ">=", "<=":
		v.Type.Id = jn.Bool
	case "+", "-", "*", "/":
		v.Type.Id = jn.F32
		if s.leftVal.Type.Id == jn.F64 || s.rightVal.Type.Id == jn.F64 {
			v.Type.Id = jn.F64
		}
	default:
		s.p.pusherrtok(s.operator, "operator_notfor_float")
	}
	return
}

func (s solver) signed() (v ast.Value) {
	if !typesAreCompatible(s.leftVal.Type, s.rightVal.Type, true) {
		if !isConstNum(s.leftVal.Data) &&
			!isConstNum(s.rightVal.Data) {
			s.p.pusherrtok(s.operator, "incompatible_datatype",
				s.rightVal.Type.Val, s.leftVal.Type.Val)
			return
		}
	}
	switch s.operator.Kind {
	case "!=", "==", "<", ">", ">=", "<=":
		v.Type.Id = jn.Bool
	case "+", "-", "*", "/", "%", "&", "|", "^":
		v.Type = s.leftVal.Type
		if jn.TypeGreaterThan(s.rightVal.Type.Id, v.Type.Id) {
			v.Type = s.rightVal.Type
		}
	case ">>", "<<":
		v.Type = s.leftVal.Type
		if !jn.IsUnsignedNumericType(s.rightVal.Type.Id) &&
			!checkIntBit(s.rightVal, jnbits.BitsizeType(jn.U64)) {
			s.p.pusherrtok(s.rightVal.Token, "bitshift_must_unsigned")
		}
	default:
		s.p.pusherrtok(s.operator, "operator_notfor_int")
	}
	return
}

func (s solver) unsigned() (v ast.Value) {
	if !typesAreCompatible(s.leftVal.Type, s.rightVal.Type, true) {
		if !isConstNum(s.leftVal.Data) &&
			!isConstNum(s.rightVal.Data) {
			s.p.pusherrtok(s.operator, "incompatible_datatype",
				s.rightVal.Type.Val, s.leftVal.Type.Val)
			return
		}
		return
	}
	switch s.operator.Kind {
	case "!=", "==", "<", ">", ">=", "<=":
		v.Type.Id = jn.Bool
	case "+", "-", "*", "/", "%", "&", "|", "^":
		v.Type = s.leftVal.Type
		if jn.TypeGreaterThan(s.rightVal.Type.Id, v.Type.Id) {
			v.Type = s.rightVal.Type
		}
	default:
		s.p.pusherrtok(s.operator, "operator_notfor_uint")
	}
	return
}

func (s solver) logical() (v ast.Value) {
	v.Type.Id = jn.Bool
	if s.leftVal.Type.Id != jn.Bool {
		s.p.pusherrtok(s.leftVal.Token, "logical_not_bool")
	}
	if s.rightVal.Type.Id != jn.Bool {
		s.p.pusherrtok(s.rightVal.Token, "logical_not_bool")
	}
	return
}

func (s solver) rune() (v ast.Value) {
	if !typesAreCompatible(s.leftVal.Type, s.rightVal.Type, true) {
		s.p.pusherrtok(s.operator, "incompatible_datatype",
			s.rightVal.Type.Val, s.leftVal.Type.Val)
		return
	}
	switch s.operator.Kind {
	case "!=", "==", ">", "<", ">=", "<=":
		v.Type.Id = jn.Bool
	case "+", "-", "*", "/", "^", "&", "%", "|":
		v.Type.Id = jn.Rune
	default:
		s.p.pusherrtok(s.operator, "operator_notfor_rune")
	}
	return
}

func (s solver) array() (v ast.Value) {
	if !typesAreCompatible(s.leftVal.Type, s.rightVal.Type, true) {
		s.p.pusherrtok(s.operator, "incompatible_datatype",
			s.rightVal.Type.Val, s.leftVal.Type.Val)
		return
	}
	switch s.operator.Kind {
	case "!=", "==":
		v.Type.Id = jn.Bool
	default:
		s.p.pusherrtok(s.operator, "operator_notfor_array")
	}
	return
}

func (s solver) nil() (v ast.Value) {
	if !typesAreCompatible(s.leftVal.Type, s.rightVal.Type, false) {
		s.p.pusherrtok(s.operator, "incompatible_datatype",
			s.rightVal.Type.Val, s.leftVal.Type.Val)
		return
	}
	switch s.operator.Kind {
	case "!=", "==":
		v.Type.Id = jn.Bool
	default:
		s.p.pusherrtok(s.operator, "operator_notfor_nil")
	}
	return
}

func (s solver) Solve() (v ast.Value) {
	switch s.operator.Kind {
	case "+", "-", "*", "/", "%", ">>",
		"<<", "&", "|", "^", "==", "!=", ">", "<", ">=", "<=":
		break
	case "&&", "||":
		return s.logical()
	default:
		s.p.pusherrtok(s.operator, "invalid_operator")
	}
	switch {
	case typeIsArray(s.leftVal.Type) || typeIsArray(s.rightVal.Type):
		return s.array()
	case typeIsPtr(s.leftVal.Type) || typeIsPtr(s.rightVal.Type):
		return s.ptr()
	case s.leftVal.Type.Id == jn.Nil || s.rightVal.Type.Id == jn.Nil:
		return s.nil()
	case s.leftVal.Type.Id == jn.Rune || s.rightVal.Type.Id == jn.Rune:
		return s.rune()
	case s.leftVal.Type.Id == jn.Any || s.rightVal.Type.Id == jn.Any:
		return s.any()
	case s.leftVal.Type.Id == jn.Bool || s.rightVal.Type.Id == jn.Bool:
		return s.bool()
	case s.leftVal.Type.Id == jn.Str || s.rightVal.Type.Id == jn.Str:
		return s.str()
	case jn.IsFloatType(s.leftVal.Type.Id) ||
		jn.IsFloatType(s.rightVal.Type.Id):
		return s.float()
	case jn.IsSignedNumericType(s.leftVal.Type.Id) ||
		jn.IsSignedNumericType(s.rightVal.Type.Id):
		return s.signed()
	case jn.IsUnsignedNumericType(s.leftVal.Type.Id) ||
		jn.IsUnsignedNumericType(s.rightVal.Type.Id):
		return s.unsigned()
	}
	return
}

func (p *Parser) evalSingleExpr(token lexer.Token, m *exprModel) (v value, ok bool) {
	eval := valueEvaluator{token, m, p}
	v.ast.Type.Id = jn.Void
	v.ast.Token = token
	switch token.Id {
	case lexer.Value:
		ok = true
		switch {
		case isstr(token.Kind):
			v = eval.str()
		case isrune(token.Kind):
			v = eval.rune()
		case isbool(token.Kind):
			v = eval.bool()
		case isnil(token.Kind):
			v = eval.nil()
		default:
			v = eval.num()
		}
	case lexer.Id:
		v, ok = eval.id()
	default:
		p.pusherrtok(token, "invalid_syntax")
	}
	return
}

type operatorProcessor struct {
	token  lexer.Token
	tokens []lexer.Token
	model  *exprModel
	parser *Parser
}

func (p *operatorProcessor) unary() value {
	v := p.parser.evalExprPart(p.tokens, p.model)
	if !typeIsSingle(v.ast.Type) {
		p.parser.pusherrtok(p.token, "invalid_data_unary")
	} else if !jn.IsNumericType(v.ast.Type.Id) {
		p.parser.pusherrtok(p.token, "invalid_data_unary")
	}
	if isConstNum(v.ast.Data) {
		v.ast.Data = "-" + v.ast.Data
	}
	return v
}

func (p *operatorProcessor) plus() value {
	v := p.parser.evalExprPart(p.tokens, p.model)
	if !typeIsSingle(v.ast.Type) {
		p.parser.pusherrtok(p.token, "invalid_data_plus")
	} else if !jn.IsNumericType(v.ast.Type.Id) {
		p.parser.pusherrtok(p.token, "invalid_data_plus")
	}
	return v
}

func (p *operatorProcessor) tilde() value {
	v := p.parser.evalExprPart(p.tokens, p.model)
	if !typeIsSingle(v.ast.Type) {
		p.parser.pusherrtok(p.token, "invalid_data_tilde")
	} else if !jn.IsIntegerType(v.ast.Type.Id) {
		p.parser.pusherrtok(p.token, "invalid_data_tilde")
	}
	return v
}

func (p *operatorProcessor) logicalNot() value {
	v := p.parser.evalExprPart(p.tokens, p.model)
	if !isBoolExpr(v) {
		p.parser.pusherrtok(p.token, "invalid_data_logical_not")
	}
	v.ast.Type.Val = "bool"
	v.ast.Type.Id = jn.Bool
	return v
}

func (p *operatorProcessor) star() value {
	v := p.parser.evalExprPart(p.tokens, p.model)
	v.lvalue = true
	if !typeIsPtr(v.ast.Type) {
		p.parser.pusherrtok(p.token, "invalid_data_star")
	} else {
		v.ast.Type.Val = v.ast.Type.Val[1:]
	}
	return v
}

func (p *operatorProcessor) amper() value {
	v := p.parser.evalExprPart(p.tokens, p.model)
	v.lvalue = true
	if !canGetPointer(v) {
		p.parser.pusherrtok(p.token, "invalid_data_amper")
	}
	v.ast.Type.Val = "*" + v.ast.Type.Val
	return v
}

func (p *Parser) evalOperatorExprPart(tokens []lexer.Token, m *exprModel) value {
	var v value

	exprToks := tokens[1:]
	processor := operatorProcessor{tokens[0], exprToks, m, p}
	m.appendSubNode(exprNode{processor.token.Kind})
	if processor.tokens == nil {
		p.pusherrtok(processor.token, "invalid_syntax")
		return v
	}
	switch processor.token.Kind {
	case "-":
		v = processor.unary()
	case "+":
		v = processor.plus()
	case "~":
		v = processor.tilde()
	case "!":
		v = processor.logicalNot()
	case "*":
		v = processor.star()
	case "&":
		v = processor.amper()
	default:
		p.pusherrtok(processor.token, "invalid_syntax")
	}
	v.ast.Token = processor.token
	return v
}

func canGetPointer(v value) bool {
	if v.ast.Type.Id == jn.Func {
		return false
	}
	return v.ast.Token.Id == lexer.Id
}

func (p *Parser) evalHeapAllocExpr(tokens []lexer.Token, m *exprModel) (v value) {
	if len(tokens) == 1 {
		p.pusherrtok(tokens[0], "invalid_syntax_keyword_new")
		return
	}
	v.lvalue = true
	v.ast.Token = tokens[0]
	tokens = tokens[1:]
	b := new(ast.Builder)
	i := new(int)
	dt, ok := b.DataType(tokens, i, true)
	m.appendSubNode(newHeapAllocExpr{dt})
	dt.Val = "*" + dt.Val
	v.ast.Type = dt
	if !ok {
		p.pusherrtok(tokens[0], "fail_build_heap_allocation_type", dt.Val)
		return
	}
	if *i < len(tokens)-1 {
		p.pusherrtok(tokens[*i+1], "invalid_syntax")
	}
	return
}

func (p *Parser) evalExprPart(tokens []lexer.Token, m *exprModel) (v value) {
	if len(tokens) == 1 {
		val, ok := p.evalSingleExpr(tokens[0], m)
		if ok {
			v = val
			return
		}
	}
	tok := tokens[0]
	switch tok.Id {
	case lexer.Operator:
		return p.evalOperatorExprPart(tokens, m)
	case lexer.New:
		return p.evalHeapAllocExpr(tokens, m)
	case lexer.Brace:
		switch tok.Kind {
		case "(":
			val, ok := p.evalTryCastExpr(tokens, m)
			if ok {
				v = val
				return
			}
			val, ok = p.evalTryAssignExpr(tokens, m)
			if ok {
				v = val
				return
			}
		}
	}
	tok = tokens[len(tokens)-1]
	switch tok.Id {
	case lexer.Id:
		return p.evalIdExprPart(tokens, m)
	case lexer.Operator:
		return p.evalOperatorExprPartRight(tokens, m)
	case lexer.Brace:
		switch tok.Kind {
		case ")":
			return p.evalParenthesesRangeExpr(tokens, m)
		case "}":
			return p.evalBraceRangeExpr(tokens, m)
		case "]":
			return p.evalBracketRangeExpr(tokens, m)
		}
	default:
		p.pusherrtok(tokens[0], "invalid_syntax")
	}
	return
}

func (p *Parser) evalStrSubId(val value, idTok lexer.Token, m *exprModel) (v value) {
	i, t := strDefs.defById(idTok.Kind)
	if i == -1 {
		p.pusherrtok(idTok, "obj_have_not_id", val.ast.Type.Val)
		return
	}
	v = val
	m.appendSubNode(exprNode{"."})
	switch t {
	case 'g':
		g := &strDefs.Globals[i]
		m.appendSubNode(exprNode{g.Tag.(string)})
		v.ast.Type = g.Type
		v.lvalue = true
		v.constant = g.Const
	default:
	}
	return v
}

func (p *Parser) evalArraySubId(val value, idTok lexer.Token, m *exprModel) (v value) {
	i, t := arrDefs.defById(idTok.Kind)
	if i == -1 {
		p.pusherrtok(idTok, "obj_have_not_id", val.ast.Type.Val)
		return
	}
	v = val
	m.appendSubNode(exprNode{"."})
	switch t {
	case 'g':
		g := &arrDefs.Globals[i]
		m.appendSubNode(exprNode{g.Tag.(string)})
		v.ast.Type = g.Type
		v.lvalue = true
		v.constant = g.Const
	default:
	}
	return v
}

func (p *Parser) evalIdExprPart(tokens []lexer.Token, m *exprModel) (v value) {
	i := len(tokens) - 1
	tok := tokens[i]
	if i <= 0 {
		v, _ = p.evalSingleExpr(tok, m)
		return
	}
	i--
	if i == 0 || tokens[i].Id != lexer.Dot {
		p.pusherrtok(tokens[i], "invalid_syntax")
		return
	}
	idTok := tokens[i+1]
	valTok := tokens[i]
	tokens = tokens[:i]
	val := p.evalExprPart(tokens, m)
	switch {
	case typeIsSingle(val.ast.Type) && val.ast.Type.Id == jn.Str:
		return p.evalStrSubId(val, idTok, m)
	case typeIsArray(val.ast.Type):
		return p.evalArraySubId(val, idTok, m)
	}
	p.pusherrtok(valTok, "obj_not_support_sub_fields", val.ast.Type.Val)
	return
}

func (p *Parser) evalTryCastExpr(tokens []lexer.Token, m *exprModel) (v value, _ bool) {
	braceCount := 0
	errTok := tokens[0]
	for i, tok := range tokens {
		if tok.Id == lexer.Brace {
			switch tok.Kind {
			case "(", "[", "{":
				braceCount++
				continue
			default:
				braceCount--
			}
		}
		if braceCount > 0 {
			continue
		}
		astb := ast.NewBuilder(nil)
		dtindex := 0
		typeToks := tokens[1:i]
		dt, ok := astb.DataType(typeToks, &dtindex, false)
		if !ok {
			return
		}
		dt, ok = p.readyType(dt, false)
		if !ok {
			return
		}
		if dtindex+1 < len(typeToks) {
			return
		}
		if i+1 >= len(tokens) {
			p.pusherrtok(tok, "casting_missing_expr")
			return
		}
		exprToks := tokens[i+1:]
		m.appendSubNode(exprNode{"(" + dt.String() + ")"})
		val := p.evalExprPart(exprToks, m)
		val = p.evalCast(val, dt, errTok)
		return val, true
	}
	return
}

func (p *Parser) evalTryAssignExpr(tokens []lexer.Token, m *exprModel) (v value, ok bool) {
	b := ast.NewBuilder(nil)
	tokens = tokens[1 : len(tokens)-1]
	assign, ok := b.AssignExpr(tokens, true)
	if !ok {
		return
	}
	ok = true
	if len(b.Errs) > 0 {
		p.pusherrs(b.Errs...)
		return
	}
	p.checkAssign(&assign)
	m.appendSubNode(assignExpr{assign})
	v, _ = p.evalExpr(assign.SelectExprs[0].Expr)
	return
}

func (p *Parser) evalCast(v value, t ast.DataType, errtok lexer.Token) value {
	switch {
	case typeIsPtr(t):
		p.checkCastPtr(v.ast.Type, errtok)
	case typeIsArray(t):
		p.checkCastArray(t, v.ast.Type, errtok)
	case typeIsSingle(t):
		v.lvalue = false
		p.checkCastSingle(t, v.ast.Type, errtok)
	default:
		p.pusherrtok(errtok, "type_notsupports_casting", t.Val)
	}
	v.ast.Type = t
	v.constant = false
	v.volatile = false
	return v
}

func (p *Parser) checkCastSingle(t, vt ast.DataType, errtok lexer.Token) {
	switch t.Id {
	case jn.Str:
		p.checkCastStr(vt, errtok)
		return
	}
	switch {
	case jn.IsIntegerType(t.Id):
		p.checkCastInteger(t, vt, errtok)
	case jn.IsNumericType(t.Id):
		p.checkCastNumeric(t, vt, errtok)
	default:
		p.pusherrtok(errtok, "type_notsupports_casting", t.Val)
	}
}

func (p *Parser) checkCastStr(vt ast.DataType, errtok lexer.Token) {
	if !typeIsArray(vt) {
		p.pusherrtok(errtok, "type_notsupports_casting", vt.Val)
		return
	}
	vt.Val = vt.Val[2:]
	if !typeIsSingle(vt) || (vt.Id != jn.Rune && vt.Id != jn.U8) {
		p.pusherrtok(errtok, "type_notsupports_casting", vt.Val)
	}
}

func (p *Parser) checkCastInteger(t, vt ast.DataType, errtok lexer.Token) {
	if typeIsPtr(vt) {
		return
	}
	if typeIsSingle(vt) && jn.IsNumericType(vt.Id) {
		return
	}
	p.pusherrtok(errtok, "type_notsupports_casting_to", vt.Val, t.Val)
}

func (p *Parser) checkCastNumeric(t, vt ast.DataType, errtok lexer.Token) {
	if typeIsSingle(vt) && jn.IsNumericType(vt.Id) {
		return
	}
	p.pusherrtok(errtok, "type_notsupports_casting_to", vt.Val, t.Val)
}

func (p *Parser) checkCastPtr(vt ast.DataType, errtok lexer.Token) {
	if typeIsPtr(vt) {
		return
	}
	if typeIsSingle(vt) && jn.IsIntegerType(vt.Id) {
		return
	}
	p.pusherrtok(errtok, "type_notsupports_casting", vt.Val)
}

func (p *Parser) checkCastArray(t, vt ast.DataType, errtok lexer.Token) {
	if !typeIsSingle(vt) || vt.Id != jn.Str {
		p.pusherrtok(errtok, "type_notsupports_casting", vt.Val)
		return
	}
	t.Val = t.Val[2:]
	if !typeIsSingle(t) || (t.Id != jn.Rune && t.Id != jn.U8) {
		p.pusherrtok(errtok, "type_notsupports_casting", vt.Val)
	}
}

func (p *Parser) evalOperatorExprPartRight(tokens []lexer.Token, m *exprModel) (v value) {
	tok := tokens[len(tokens)-1]
	switch tok.Kind {
	case "...":
		tokens = tokens[:len(tokens)-1]
		return p.evalVariadicExprPart(tokens, m, tok)
	default:
		p.pusherrtok(tok, "invalid_syntax")
	}
	return
}

func (p *Parser) evalVariadicExprPart(
	tokens []lexer.Token,
	m *exprModel,
	errtok lexer.Token,
) (v value) {
	v = p.evalExprPart(tokens, m)
	if !typeIsVariadicable(v.ast.Type) {
		p.pusherrtok(errtok, "variadic_with_nonvariadicable", v.ast.Type.Val)
		return
	}
	v.ast.Type.Val = v.ast.Type.Val[2:]
	v.variadic = true
	return
}

func (p *Parser) evalParenthesesRangeExpr(tokens []lexer.Token, m *exprModel) (v value) {
	var valueToks []lexer.Token
	braceCount := 0
	for i := len(tokens) - 1; i >= 0; i-- {
		tok := tokens[i]
		if tok.Id != lexer.Brace {
			continue
		}
		switch tok.Kind {
		case ")", "}", "]":
			braceCount++
		case "(", "{", "[":
			braceCount--
		}
		if braceCount > 0 {
			continue
		}
		valueToks = tokens[:i]
		break
	}
	if len(valueToks) == 0 && braceCount == 0 {
		m.appendSubNode(exprNode{"("})
		defer m.appendSubNode(exprNode{")"})

		tk := tokens[0]
		tokens = tokens[1 : len(tokens)-1]
		if len(tokens) == 0 {
			p.pusherrtok(tk, "invalid_syntax")
		}
		val, model := p.evalToks(tokens)
		v = val
		m.appendSubNode(model)
		return
	}
	v = p.evalExprPart(valueToks, m)

	m.appendSubNode(exprNode{"("})
	defer m.appendSubNode(exprNode{")"})

	switch v.ast.Type.Id {
	case jn.Func:
		fun := v.ast.Type.Tag.(ast.Func)
		p.parseFuncCall(fun, tokens[len(valueToks):], m)
		v.ast.Type = fun.RetType
		v.lvalue = typeIsLvalue(v.ast.Type)
	default:
		p.pusherrtok(tokens[len(valueToks)], "invalid_syntax")
	}
	return
}

func (p *Parser) evalBraceRangeExpr(tokens []lexer.Token, m *exprModel) (v value) {
	var exprToks []lexer.Token
	braceCount := 0
	for i := len(tokens) - 1; i >= 0; i-- {
		tok := tokens[i]
		if tok.Id != lexer.Brace {
			continue
		}
		switch tok.Kind {
		case "}", "]", ")":
			braceCount++
		case "{", "(", "[":
			braceCount--
		}
		if braceCount > 0 {
			continue
		}
		exprToks = tokens[:i]
		break
	}
	valToksLen := len(exprToks)
	if valToksLen == 0 || braceCount > 0 {
		p.pusherrtok(tokens[0], "invalid_syntax")
		return
	}
	switch exprToks[0].Id {
	case lexer.Brace:
		switch exprToks[0].Kind {
		case "[":
			ast := ast.NewBuilder(nil)
			dt, ok := ast.DataType(exprToks, new(int), true)
			if !ok {
				p.pusherrs(ast.Errs...)
				return
			}
			exprToks = tokens[len(exprToks):]
			var model iExpr
			v, model = p.buildArray(p.buildEnumerableParts(exprToks),
				dt, exprToks[0])
			m.appendSubNode(model)
			return
		case "(":
			b := ast.NewBuilder(tokens)
			f := b.Func(b.Tokens, true)
			if len(b.Errs) > 0 {
				p.pusherrs(b.Errs...)
				return
			}
			p.checkAnonFunc(&f)
			v.ast.Type.Tag = f
			v.ast.Type.Id = jn.Func
			v.ast.Type.Val = f.DataTypeString()
			m.appendSubNode(anonFunc{f})
			return
		default:
			p.pusherrtok(exprToks[0], "invalid_syntax")
		}
	default:
		p.pusherrtok(exprToks[0], "invalid_syntax")
	}
	return
}

func (p *Parser) evalBracketRangeExpr(tokens []lexer.Token, m *exprModel) (v value) {
	var exprToks []lexer.Token
	braceCount := 0
	for i := len(tokens) - 1; i >= 0; i-- {
		tok := tokens[i]
		if tok.Id != lexer.Brace {
			continue
		}
		switch tok.Kind {
		case "}", "]", ")":
			braceCount++
		case "{", "(", "[":
			braceCount--
		}
		if braceCount > 0 {
			continue
		}
		exprToks = tokens[:i]
		break
	}
	valToksLen := len(exprToks)
	if valToksLen == 0 || braceCount > 0 {
		p.pusherrtok(tokens[0], "invalid_syntax")
		return
	}
	var model iExpr
	v, model = p.evalToks(exprToks)
	m.appendSubNode(model)
	tokens = tokens[len(exprToks)+1 : len(tokens)-1]
	m.appendSubNode(exprNode{"["})
	selectv, model := p.evalToks(tokens)
	m.appendSubNode(model)
	m.appendSubNode(exprNode{"]"})
	return p.evalEnumerableSelect(v, selectv, tokens[0])
}

func (p *Parser) evalEnumerableSelect(enumv, selectv value, errtok lexer.Token) (v value) {
	switch {
	case typeIsArray(enumv.ast.Type):
		return p.evalArraySelect(enumv, selectv, errtok)
	case typeIsSingle(enumv.ast.Type):
		return p.evalStrSelect(enumv, selectv, errtok)
	}
	p.pusherrtok(errtok, "not_enumerable")
	return
}

func (p *Parser) evalArraySelect(arrv, selectv value, errtok lexer.Token) value {
	arrv.lvalue = true
	arrv.ast.Type = typeOfArrayElements(arrv.ast.Type)
	if !typeIsSingle(selectv.ast.Type) ||
		!jn.IsIntegerType(selectv.ast.Type.Id) {
		p.pusherrtok(errtok, "notint_array_select")
	}
	return arrv
}

func (p *Parser) evalStrSelect(strv, selectv value, errtok lexer.Token) value {
	strv.lvalue = true
	strv.ast.Type.Id = jn.Rune
	if !typeIsSingle(selectv.ast.Type) ||
		!jn.IsIntegerType(selectv.ast.Type.Id) {
		p.pusherrtok(errtok, "notint_string_select")
	}
	return strv
}

func (p *Parser) buildEnumerableParts(tokens []lexer.Token) [][]lexer.Token {
	tokens = tokens[1 : len(tokens)-1]
	braceCount := 0
	lastComma := -1
	var parts [][]lexer.Token
	for i, tok := range tokens {
		if tok.Id == lexer.Brace {
			switch tok.Kind {
			case "{", "[", "(":
				braceCount++
			default:
				braceCount--
			}
		}
		if braceCount > 0 {
			continue
		}
		if tok.Id == lexer.Comma {
			if i-lastComma-1 == 0 {
				p.pusherrtok(tok, "missing_expr")
				lastComma = i
				continue
			}
			parts = append(parts, tokens[lastComma+1:i])
			lastComma = i
		}
	}
	if lastComma+1 < len(tokens) {
		parts = append(parts, tokens[lastComma+1:])
	}
	return parts
}

func (p *Parser) buildArray(
	parts [][]lexer.Token,
	t ast.DataType,
	errtok lexer.Token,
) (value, iExpr) {
	var v value
	v.ast.Type = t
	model := arrayExpr{dataType: t}
	elemType := typeOfArrayElements(t)
	for _, part := range parts {
		partVal, expModel := p.evalToks(part)
		model.expr = append(model.expr, expModel)
		p.wg.Add(1)
		go assignChecker{
			p,
			false,
			elemType,
			partVal,
			false,
			part[0],
		}.checkAssignTypeAsync()
	}
	return v, model
}

func (p *Parser) checkAnonFunc(f *ast.Func) {
	globals := p.Defs.Globals
	blockVariables := p.BlockVars
	p.Defs.Globals = append(blockVariables, p.Defs.Globals...)
	p.BlockVars = p.varsFromParams(f.Params)
	p.checkFunc(f)
	p.Defs.Globals = globals
	p.BlockVars = blockVariables
}

func (p *Parser) parseFuncCall(f ast.Func, tokens []lexer.Token, m *exprModel) {
	errTok := tokens[0]
	tokens, _ = p.getRange("(", ")", tokens)
	if tokens == nil {
		tokens = make([]lexer.Token, 0)
	}
	b := new(ast.Builder)
	args := b.Args(tokens)
	if len(b.Errs) > 0 {
		p.pusherrs(b.Errs...)
	}
	p.parseArgs(f.Params, &args, errTok, m)
	if m != nil {
		m.appendSubNode(argsExpr{args})
	}
}

func (p *Parser) parseArgs(
	params []ast.Parameter,
	args *[]ast.Arg,
	errTok lexer.Token,
	m *exprModel,
) {
	parsedArgs := make([]ast.Arg, 0)
	if len(params) > 0 && params[len(params)-1].Variadic {
		if len(*args) == 0 && len(params) == 1 {
			return
		} else if len(*args) < len(params)-1 {
			p.pusherrtok(errTok, "missing_argument")
			goto argParse
		} else if len(*args) <= len(params)-1 {
			goto argParse
		}
		variadicArgs := (*args)[len(params)-1:]
		variadicParam := params[len(params)-1]
		*args = (*args)[:len(params)-1]
		params = params[:len(params)-1]
		defer func() {
			model := arrayExpr{variadicParam.Type, nil}
			model.dataType.Val = "[]" + model.dataType.Val
			variadiced := false
			for _, arg := range variadicArgs {
				p.parseArg(variadicParam, &arg, &variadiced)
				model.expr = append(model.expr, arg.Expr.Model.(iExpr))
			}
			if variadiced && len(variadicArgs) > 1 {
				p.pusherrtok(errTok, "more_args_with_varidiced")
			}
			arg := ast.Arg{Expr: ast.Expr{Model: model}}
			parsedArgs = append(parsedArgs, arg)
			*args = parsedArgs
		}()
	}
	if len(*args) == 0 && len(params) == 0 {
		return
	} else if len(*args) < len(params) {
		p.pusherrtok(errTok, "missing_argument")
	} else if len(*args) > len(params) {
		p.pusherrtok(errTok, "argument_overflow")
		return
	}
argParse:
	for i, arg := range *args {
		p.parseArg(params[i], &arg, nil)
		parsedArgs = append(parsedArgs, arg)
	}
	*args = parsedArgs
}

func (p *Parser) parseArg(param ast.Parameter, arg *ast.Arg, variadiced *bool) {
	value, model := p.evalExpr(arg.Expr)
	arg.Expr.Model = model
	if variadiced != nil && !*variadiced {
		*variadiced = value.variadic
	}
	p.wg.Add(1)
	go p.checkArgTypeAsync(param, value, false, arg.Token)
}

func (p *Parser) checkArgTypeAsync(
	param ast.Parameter,
	val value,
	ignoreAny bool,
	errTok lexer.Token,
) {
	defer func() { p.wg.Done() }()
	p.wg.Add(1)
	go assignChecker{
		p,
		param.Const,
		param.Type,
		val,
		false,
		errTok,
	}.checkAssignTypeAsync()
}

func (p *Parser) getRange(open, close string, tokens []lexer.Token) (_ []lexer.Token, ok bool) {
	braceCount := 0
	start := 1
	if tokens[0].Id != lexer.Brace {
		return nil, false
	}
	for i, tok := range tokens {
		if tok.Id != lexer.Brace {
			continue
		}
		if tok.Kind == open {
			braceCount++
		} else if tok.Kind == close {
			braceCount--
		}
		if braceCount > 0 {
			continue
		}
		return tokens[start:i], true
	}
	return nil, false
}

func (p *Parser) checkEntryPointSpecialCases(fun *function) {
	if len(fun.Ast.Params) > 0 {
		p.pusherrtok(fun.Ast.Token, "entrypoint_have_parameters")
	}
	if fun.Ast.RetType.Id != jn.Void {
		p.pusherrtok(fun.Ast.RetType.Token, "entrypoint_have_return")
	}
	if fun.Attributes != nil {
		p.pusherrtok(fun.Ast.Token, "entrypoint_have_attributes")
	}
}

func (p *Parser) checkBlock(b *ast.BlockAST) {
	for i := 0; i < len(b.Tree); i++ {
		model := &b.Tree[i]
		switch t := model.Val.(type) {
		case ast.ExprStatement:
			_, t.Expr.Model = p.evalExpr(t.Expr)
			model.Val = t
		case ast.Var:
			p.checkVarStatement(&t, false)
			model.Val = t
		case ast.Assign:
			p.checkAssign(&t)
			model.Val = t
		case ast.Free:
			p.checkFreeStatement(&t)
			model.Val = t
		case ast.Iter:
			p.checkIterExpr(&t)
			model.Val = t
		case ast.Break:
			p.checkBreakStatement(&t)
		case ast.Continue:
			p.checkContinueStatement(&t)
		case ast.If:
			p.checkIfExpr(&t, &i, b.Tree)
			model.Val = t
		case ast.CxxEmbed:
		case ast.Comment:
		case ast.Ret:
		default:
			p.pusherrtok(model.Token, "invalid_syntax")
		}
	}
}

type retChecker struct {
	p        *Parser
	retAST   *ast.Ret
	fun      *ast.Func
	expModel multiRetExpr
	values   []value
}

func (rc *retChecker) pushval(last, current int, errTk lexer.Token) {
	if current-last == 0 {
		rc.p.pusherrtok(errTk, "missing_expr")
		return
	}
	toks := rc.retAST.Expr.Tokens[last:current]
	val, model := rc.p.evalToks(toks)
	rc.expModel.models = append(rc.expModel.models, model)
	rc.values = append(rc.values, val)
}

func (rc *retChecker) checkepxrs() {
	braceCount := 0
	last := 0
	for i, tok := range rc.retAST.Expr.Tokens {
		if tok.Id == lexer.Brace {
			switch tok.Kind {
			case "(", "{", "[":
				braceCount++
			default:
				braceCount--
			}
		}
		if braceCount > 0 || tok.Id != lexer.Comma {
			continue
		}
		rc.pushval(last, i, tok)
		last = i + 1
	}
	length := len(rc.retAST.Expr.Tokens)
	if last < length {
		if last == 0 {
			rc.pushval(0, length, rc.retAST.Token)
		} else {
			rc.pushval(last, length, rc.retAST.Expr.Tokens[last-1])
		}
	}
	if !typeIsVoidRet(rc.fun.RetType) {
		rc.checkExprTypes()
	}
}

func (rc *retChecker) checkExprTypes() {
	valLength := len(rc.values)
	if !rc.fun.RetType.MultiTyped {
		rc.retAST.Expr.Model = rc.expModel.models[0]
		if valLength > 1 {
			rc.p.pusherrtok(rc.retAST.Token, "overflow_return")
		}
		rc.p.wg.Add(1)
		go assignChecker{
			p:         rc.p,
			constant:  false,
			t:         rc.fun.RetType,
			v:         rc.values[0],
			ignoreAny: false,
			errtok:    rc.retAST.Token,
		}.checkAssignTypeAsync()
		return
	}
	rc.retAST.Expr.Model = rc.expModel
	types := rc.fun.RetType.Tag.([]ast.DataType)
	if valLength == 1 {
		rc.p.pusherrtok(rc.retAST.Token, "missing_multi_return")
	} else if valLength > len(types) {
		rc.p.pusherrtok(rc.retAST.Token, "overflow_return")
	}
	for i, t := range types {
		if i >= valLength {
			break
		}
		rc.p.wg.Add(1)
		go assignChecker{
			p:         rc.p,
			constant:  false,
			t:         t,
			v:         rc.values[i],
			ignoreAny: false,
			errtok:    rc.retAST.Token,
		}.checkAssignTypeAsync()
	}
}

func (rc *retChecker) check() {
	exprToksLen := len(rc.retAST.Expr.Tokens)
	if exprToksLen == 0 && !typeIsVoidRet(rc.fun.RetType) {
		rc.p.pusherrtok(rc.retAST.Token, "require_return_value")
		return
	}
	if exprToksLen > 0 && typeIsVoidRet(rc.fun.RetType) {
		rc.p.pusherrtok(rc.retAST.Token, "void_function_return_value")
	}
	rc.checkepxrs()
}

func (p *Parser) checkRets(fun *ast.Func) {
	missed := true
	for i, s := range fun.Block.Tree {
		switch t := s.Val.(type) {
		case ast.Ret:
			rc := retChecker{p: p, retAST: &t, fun: fun}
			rc.check()
			fun.Block.Tree[i].Val = t
			missed = false
		}
	}
	if missed && !typeIsVoidRet(fun.RetType) {
		p.pusherrtok(fun.Token, "missing_ret")
	}
}

func (p *Parser) checkFunc(f *ast.Func) {
	p.checkBlock(&f.Block)
	p.checkRets(f)
}

func (p *Parser) checkVarStatement(vast *ast.Var, noParse bool) {
	if p.existIdf(vast.Id, true).Id != lexer.NA {
		p.pusherrtok(vast.IdToken, "exist_id", vast.Id)
	}
	if !noParse {
		*vast = p.Var(*vast)
	}
	p.BlockVars = append(p.BlockVars, *vast)
}

func (p *Parser) checkAssignment(selected value, errtok lexer.Token) bool {
	state := true
	if !selected.lvalue {
		p.pusherrtok(errtok, "assign_nonlvalue")
		state = false
	}
	if selected.constant {
		p.pusherrtok(errtok, "assign_const")
		state = false
	}
	switch selected.ast.Type.Tag.(type) {
	case ast.Func:
		if p.FuncById(selected.ast.Token.Kind) != nil {
			p.pusherrtok(errtok, "assign_type_not_support_value")
			state = false
		}
	}
	return state
}

func (p *Parser) checkSingleAssign(assign *ast.Assign) {
	sexpr := &assign.SelectExprs[0].Expr
	if len(sexpr.Tokens) == 1 && jnapi.IsIgnoreId(sexpr.Tokens[0].Kind) {
		return
	}
	selected, _ := p.evalExpr(*sexpr)
	if !p.checkAssignment(selected, assign.Setter) {
		return
	}
	vexpr := &assign.ValueExprs[0]
	val, model := p.evalExpr(*vexpr)
	*vexpr = model.(*exprModel).Expr()
	if assign.Setter.Kind != "=" {
		assign.Setter.Kind = assign.Setter.Kind[:len(assign.Setter.Kind)-1]
		solver := solver{
			p:        p,
			left:     sexpr.Tokens,
			leftVal:  selected.ast,
			right:    vexpr.Tokens,
			rightVal: val.ast,
			operator: assign.Setter,
		}
		val.ast = solver.Solve()
		assign.Setter.Kind += "="
	}
	p.wg.Add(1)
	go assignChecker{
		p,
		selected.constant,
		selected.ast.Type,
		val,
		false,
		assign.Setter,
	}.checkAssignTypeAsync()
}

func (p *Parser) assignExprs(vsAST *ast.Assign) []value {
	vals := make([]value, len(vsAST.ValueExprs))
	for i, expr := range vsAST.ValueExprs {
		val, model := p.evalExpr(expr)
		vsAST.ValueExprs[i].Model = model
		vals[i] = val
	}
	return vals
}

func (p *Parser) processFuncMultiAssign(vsAST *ast.Assign, funcVal value) {
	types := funcVal.ast.Type.Tag.([]ast.DataType)
	if len(types) != len(vsAST.SelectExprs) {
		p.pusherrtok(vsAST.Setter, "missing_multiassign_identifiers")
		return
	}
	vals := make([]value, len(types))
	for i, t := range types {
		vals[i] = value{
			ast: ast.Value{
				Token: t.Token,
				Type:  t,
			},
		}
	}
	p.processMultiAssign(vsAST, vals)
}

func (p *Parser) processMultiAssign(assign *ast.Assign, vals []value) {
	for i := range assign.SelectExprs {
		selector := &assign.SelectExprs[i]
		selector.Ignore = jnapi.IsIgnoreId(selector.Var.Id)
		val := vals[i]
		if !selector.Var.New {
			if selector.Ignore {
				continue
			}
			selected, _ := p.evalExpr(selector.Expr)
			if !p.checkAssignment(selected, assign.Setter) {
				return
			}
			p.wg.Add(1)
			go assignChecker{
				p,
				selected.constant,
				selected.ast.Type,
				val,
				false,
				assign.Setter,
			}.checkAssignTypeAsync()
			continue
		}
		selector.Var.Tag = val
		p.checkVarStatement(&selector.Var, false)
	}
}

func (p *Parser) checkAssign(assign *ast.Assign) {
	selectLength := len(assign.SelectExprs)
	valueLength := len(assign.ValueExprs)
	if selectLength == 1 && !assign.SelectExprs[0].Var.New {
		p.checkSingleAssign(assign)
		return
	} else if assign.Setter.Kind != "=" {
		p.pusherrtok(assign.Setter, "invalid_syntax")
		return
	}
	if valueLength == 1 {
		firstVal, _ := p.evalExpr(assign.ValueExprs[0])
		if firstVal.ast.Type.MultiTyped {
			assign.MultipleRet = true
			p.processFuncMultiAssign(assign, firstVal)
			return
		}
	}
	switch {
	case selectLength > valueLength:
		p.pusherrtok(assign.Setter, "overflow_multiassign_identifiers")
		return
	case selectLength < valueLength:
		p.pusherrtok(assign.Setter, "missing_multiassign_identifiers")
		return
	}
	p.processMultiAssign(assign, p.assignExprs(assign))
}

func (p *Parser) checkFreeStatement(freeAST *ast.Free) {
	val, model := p.evalExpr(freeAST.Expr)
	freeAST.Expr.Model = model
	if !typeIsPtr(val.ast.Type) {
		p.pusherrtok(freeAST.Token, "free_nonpointer")
	}
}

func (p *Parser) checkWhileProfile(iter *ast.Iter) {
	profile := iter.Profile.(ast.WhileProfile)
	val, model := p.evalExpr(profile.Expr)
	profile.Expr.Model = model
	iter.Profile = profile
	if !isBoolExpr(val) {
		p.pusherrtok(iter.Token, "iter_while_notbool_expr")
	}
	p.checkBlock(&iter.Block)
}

type foreachTypeChecker struct {
	p       *Parser
	profile *ast.ForeachProfile
	value   value
}

func (frc *foreachTypeChecker) array() {
	if !jnapi.IsIgnoreId(frc.profile.KeyA.Id) {
		keyA := &frc.profile.KeyA
		if keyA.Type.Id == jn.Void {
			keyA.Type.Id = jn.Size
			keyA.Type.Val = jn.CxxTypeIdFromType(keyA.Type.Id)
		} else {
			var ok bool
			keyA.Type, ok = frc.p.readyType(keyA.Type, true)
			if ok {
				if !typeIsSingle(keyA.Type) || !jn.IsNumericType(keyA.Type.Id) {
					frc.p.pusherrtok(keyA.IdToken, "incompatible_datatype",
						keyA.Type.Val, jn.NumericTypeStr)
				}
			}
		}
	}
	if !jnapi.IsIgnoreId(frc.profile.KeyB.Id) {
		elementType := frc.profile.ExprType
		elementType.Val = elementType.Val[2:]
		keyB := &frc.profile.KeyB
		if keyB.Type.Id == jn.Void {
			keyB.Type = elementType
		} else {
			frc.p.wg.Add(1)
			go frc.p.checkTypeAsync(elementType, frc.profile.KeyB.Type, true, frc.profile.InToken)
		}
	}
}

func (frc *foreachTypeChecker) keyA() {
	if jnapi.IsIgnoreId(frc.profile.KeyA.Id) {
		return
	}
	keyA := &frc.profile.KeyA
	if keyA.Type.Id == jn.Void {
		keyA.Type.Id = jn.Size
		keyA.Type.Val = jn.CxxTypeIdFromType(keyA.Type.Id)
		return
	}
	var ok bool
	keyA.Type, ok = frc.p.readyType(keyA.Type, true)
	if ok {
		if !typeIsSingle(keyA.Type) || !jn.IsNumericType(keyA.Type.Id) {
			frc.p.pusherrtok(keyA.IdToken, "incompatible_datatype",
				keyA.Type.Val, jn.NumericTypeStr)
		}
	}
}

func (frc *foreachTypeChecker) keyB() {
	if jnapi.IsIgnoreId(frc.profile.KeyB.Id) {
		return
	}
	runeType := ast.DataType{
		Id:  jn.Rune,
		Val: jn.CxxTypeIdFromType(jn.Rune),
	}
	keyB := &frc.profile.KeyB
	if keyB.Type.Id == jn.Void {
		keyB.Type = runeType
		return
	}
	frc.p.wg.Add(1)
	go frc.p.checkTypeAsync(runeType, frc.profile.KeyB.Type, true, frc.profile.InToken)
}

func (frc *foreachTypeChecker) str() {
	frc.keyA()
	frc.keyB()
}

func (ftc *foreachTypeChecker) check() {
	switch {
	case typeIsArray(ftc.value.ast.Type):
		ftc.array()
	case ftc.value.ast.Type.Id == jn.Str:
		ftc.str()
	}
}

func (p *Parser) checkForeachProfile(iter *ast.Iter) {
	profile := iter.Profile.(ast.ForeachProfile)
	val, model := p.evalExpr(profile.Expr)
	profile.Expr.Model = model
	profile.ExprType = val.ast.Type
	if !isForeachIterExpr(val) {
		p.pusherrtok(iter.Token, "iter_foreach_nonenumerable_expr")
	} else {
		checker := foreachTypeChecker{p, &profile, val}
		checker.check()
	}
	iter.Profile = profile
	blockVariables := p.BlockVars
	if profile.KeyA.New {
		if jnapi.IsIgnoreId(profile.KeyA.Id) {
			p.pusherrtok(profile.KeyA.IdToken, "ignore_id")
		}
		p.checkVarStatement(&profile.KeyA, true)
	}
	if profile.KeyB.New {
		if jnapi.IsIgnoreId(profile.KeyB.Id) {
			p.pusherrtok(profile.KeyB.IdToken, "ignore_id")
		}
		p.checkVarStatement(&profile.KeyB, true)
	}
	p.checkBlock(&iter.Block)
	p.BlockVars = blockVariables
}

func (p *Parser) checkIterExpr(iter *ast.Iter) {
	p.iterCount++
	if iter.Profile != nil {
		switch iter.Profile.(type) {
		case ast.WhileProfile:
			p.checkWhileProfile(iter)
		case ast.ForeachProfile:
			p.checkForeachProfile(iter)
		}
	}
	p.iterCount--
}

func (p *Parser) checkIfExpr(ifast *ast.If, i *int, statements []ast.Statement) {
	val, model := p.evalExpr(ifast.Expr)
	ifast.Expr.Model = model
	statement := statements[*i]
	if !isBoolExpr(val) {
		p.pusherrtok(ifast.Token, "if_notbool_expr")
	}
	p.checkBlock(&ifast.Block)
node:
	if statement.WithTerminator {
		return
	}
	*i++
	if *i >= len(statements) {
		*i--
		return
	}
	statement = statements[*i]
	switch t := statement.Val.(type) {
	case ast.ElseIf:
		val, model := p.evalExpr(t.Expr)
		t.Expr.Model = model
		if !isBoolExpr(val) {
			p.pusherrtok(t.Token, "if_notbool_expr")
		}
		p.checkBlock(&t.Block)
		goto node
	case ast.Else:
		p.checkElseBlock(&t)
		statement.Val = t
	default:
		*i--
	}
}

func (p *Parser) checkElseBlock(elseast *ast.Else) {
	p.checkBlock(&elseast.Block)
}

func (p *Parser) checkBreakStatement(breakAST *ast.Break) {
	if p.iterCount == 0 {
		p.pusherrtok(breakAST.Token, "break_at_outiter")
	}
}

func (p *Parser) checkContinueStatement(continueAST *ast.Continue) {
	if p.iterCount == 0 {
		p.pusherrtok(continueAST.Token, "continue_at_outiter")
	}
}

func (p *Parser) checkValidityForAutoType(t ast.DataType, err lexer.Token) {
	switch t.Id {
	case jn.Nil:
		p.pusherrtok(err, "nil_for_autotype")
	case jn.Void:
		p.pusherrtok(err, "void_for_autotype")
	}
}

func (p *Parser) defaultValueOfType(t ast.DataType) string {
	if typeIsPtr(t) || typeIsArray(t) {
		return "nil"
	}
	return jn.DefaultValOfType(t.Id)
}

func (p *Parser) readyType(dt ast.DataType, err bool) (_ ast.DataType, ok bool) {
	if dt.Val == "" {
		return dt, true
	}
	if dt.MultiTyped {
		types := dt.Tag.([]ast.DataType)
		for i, t := range types {
			t, okr := p.readyType(t, err)
			types[i] = t
			if ok {
				ok = okr
			}
		}
		dt.Tag = types
		return dt, ok
	}
	switch dt.Id {
	case jn.Id:
		t := p.typeById(dt.Token.Kind)
		if t == nil {
			if err {
				p.pusherrtok(dt.Token, "invalid_type_source")
			}
			return dt, false
		}
		t.Type.Val = dt.Val[:len(dt.Val)-len(dt.Token.Kind)] + t.Type.Val
		return p.readyType(t.Type, err)
	case jn.Func:
		f := dt.Tag.(ast.Func)
		for i, param := range f.Params {
			f.Params[i].Type, _ = p.readyType(param.Type, err)
		}
		f.RetType, _ = p.readyType(f.RetType, err)
		dt.Val = dt.Tag.(ast.Func).DataTypeString()
	}
	return dt, true
}

func (p *Parser) checkMultiTypeAsync(real, check ast.DataType, ignoreAny bool, errTok lexer.Token) {
	defer func() { p.wg.Done() }()
	if real.MultiTyped != check.MultiTyped {
		p.pusherrtok(errTok, "incompatible_datatype", real.Val, check.Val)
		return
	}
	realTypes := real.Tag.([]ast.DataType)
	checkTypes := real.Tag.([]ast.DataType)
	if len(realTypes) != len(checkTypes) {
		p.pusherrtok(errTok, "incompatible_datatype", real.Val, check.Val)
		return
	}
	for i := 0; i < len(realTypes); i++ {
		realType := realTypes[i]
		checkType := checkTypes[i]
		p.checkTypeAsync(realType, checkType, ignoreAny, errTok)
	}
}

func (p *Parser) checkAssignConst(constant bool, t ast.DataType, val value, errTok lexer.Token) {
	if typeIsMut(t) && val.constant && !constant {
		p.pusherrtok(errTok, "constant_assignto_nonconstant")
	}
}

type assignChecker struct {
	p         *Parser
	constant  bool
	t         ast.DataType
	v         value
	ignoreAny bool
	errtok    lexer.Token
}

func (ac assignChecker) checkAssignTypeAsync() {
	defer func() { ac.p.wg.Done() }()
	ac.p.checkAssignConst(ac.constant, ac.t, ac.v, ac.errtok)
	if typeIsSingle(ac.t) && isConstNum(ac.v.ast.Data) {
		switch {
		case jn.IsSignedIntegerType(ac.t.Id):
			if jnbits.CheckBitInt(ac.v.ast.Data, jnbits.BitsizeType(ac.t.Id)) {
				return
			}
			ac.p.pusherrtok(ac.errtok, "incompatible_datatype", ac.t, ac.v.ast.Type)
			return
		case jn.IsFloatType(ac.t.Id):
			if checkFloatBit(ac.v.ast, jnbits.BitsizeType(ac.t.Id)) {
				return
			}
			ac.p.pusherrtok(ac.errtok, "incompatible_datatype", ac.t, ac.v.ast.Type)
			return
		case jn.IsUnsignedNumericType(ac.t.Id):
			if jnbits.CheckBitUInt(ac.v.ast.Data, jnbits.BitsizeType(ac.t.Id)) {
				return
			}
			ac.p.pusherrtok(ac.errtok, "incompatible_datatype", ac.t, ac.v.ast.Type)
			return
		}
	}
	ac.p.wg.Add(1)
	go ac.p.checkTypeAsync(ac.t, ac.v.ast.Type, ac.ignoreAny, ac.errtok)
}

func (p *Parser) checkTypeAsync(real, check ast.DataType, ignoreAny bool, errTok lexer.Token) {
	defer func() { p.wg.Done() }()
	if !ignoreAny && real.Id == jn.Any {
		return
	}
	if real.MultiTyped || check.MultiTyped {
		p.wg.Add(1)
		go p.checkMultiTypeAsync(real, check, ignoreAny, errTok)
		return
	}
	if typeIsSingle(real) && typeIsSingle(check) {
		if !typesAreCompatible(real, check, ignoreAny) {
			p.pusherrtok(errTok, "incompatible_datatype", real.Val, check.Val)
		}
		return
	}
	if (typeIsPtr(real) || typeIsArray(real)) && check.Id == jn.Nil {
		return
	}
	if real.Val != check.Val {
		p.pusherrtok(errTok, "incompatible_datatype", real.Val, check.Val)
	}
}
