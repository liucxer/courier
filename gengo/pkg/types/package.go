package types

import (
	"go/ast"
	"go/token"
	"go/types"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/packages"
)

type Package interface {
	// Pkg of go package
	Pkg() *types.Package
	// SourceDir code source absolute dir
	SourceDir() string
	// Files ast files of package
	Files() []*ast.File
	// Doc leading comments for pos
	Doc(pos token.Pos) []string
	// Comment trailing comments for pos
	Comment(pos token.Pos) []string

	// Eval eval expr in package
	Eval(expr ast.Expr) (types.TypeAndValue, error)

	// Constant get constant by name
	Constant(name string) *types.Const
	// Constants get all constants of package
	Constants() map[string]*types.Const
	// Type get type by name
	Type(name string) *types.TypeName
	// Types get all types of package
	Types() map[string]*types.TypeName
	// Function get function by name
	Function(name string) *types.Func
	// Functions get all signatures of package
	Functions() map[string]*types.Func
	// MethodsOf get methods of types.TypeName
	MethodsOf(n *types.Named, canPtr bool) []*types.Func
	// ResultsOf get possible resolveFuncResults of function
	ResultsOf(tpe *types.Func) (results Results, resultN int)
	// Position get position of pos
	Position(pos token.Pos) token.Position
}

func newPkg(pkg *packages.Package, u Universe) Package {
	p := &pkgInfo{
		u: u,

		Package: pkg,

		endLineToCommentGroup:         map[fileLine]*ast.CommentGroup{},
		endLineToTrailingCommentGroup: map[fileLine]*ast.CommentGroup{},

		signatures:  map[*types.Signature]ast.Node{},
		funcResults: map[*types.Signature][]TypeAndValues{},

		constants: map[string]*types.Const{},
		types:     map[string]*types.TypeName{},
		funcs:     map[string]*types.Func{},

		methods: map[*types.Named][]*types.Func{},
	}

	fileLineFor := func(pos token.Pos) fileLine {
		position := p.Fset.Position(pos)
		return fileLine{position.Filename, position.Line}
	}

	collectCommentGroup := func(c *ast.CommentGroup, isTrailing bool) {
		if c != nil {
			fl := fileLineFor(c.End())

			if isTrailing {
				p.endLineToTrailingCommentGroup[fl] = c
			} else {
				p.endLineToCommentGroup[fl] = c
			}
		}
	}

	for i := range p.Syntax {
		f := p.Syntax[i]

		ast.Inspect(f, func(node ast.Node) bool {
			switch x := node.(type) {
			case *ast.CallExpr:
				// signature will be from other package
				// stored p.TypesInfo.Uses[*ast.Ident].(*types.PkgName)
				fn := p.TypesInfo.TypeOf(x.Fun)
				if fn != nil {
					if s, ok := fn.(*types.Signature); ok {
						if n, ok := p.signatures[s]; ok {
							switch n.(type) {
							case *ast.FuncDecl, *ast.FuncLit:
								// skip declared functions
							default:
								p.signatures[s] = x.Fun
							}
						} else {
							p.signatures[s] = x.Fun
						}
					}
				}
			case *ast.FuncDecl:
				fn := p.TypesInfo.TypeOf(x.Name)
				if fn != nil {
					p.signatures[fn.(*types.Signature)] = x
				}
			case *ast.FuncLit:
				fn := p.TypesInfo.TypeOf(x)
				if fn != nil {
					p.signatures[fn.(*types.Signature)] = x
				}
			case *ast.CommentGroup:
				collectCommentGroup(x, false)
			case *ast.ValueSpec:
				collectCommentGroup(x.Comment, true)
			case *ast.ImportSpec:
				collectCommentGroup(x.Comment, true)
			case *ast.TypeSpec:
				collectCommentGroup(x.Comment, true)
			case *ast.Field:
				collectCommentGroup(x.Comment, true)
			}
			return true
		})
	}

	for ident := range p.TypesInfo.Defs {
		switch x := p.TypesInfo.Defs[ident].(type) {
		case *types.Func:
			s := x.Type().(*types.Signature)

			if r := s.Recv(); r != nil {
				var named *types.Named

				switch t := r.Type().(type) {
				case *types.Pointer:
					named = t.Elem().(*types.Named)
				case *types.Named:
					named = t
				}

				if named != nil {
					p.methods[named] = append(p.methods[named], x)
				}
			} else {
				p.funcs[x.Name()] = x
			}
		case *types.TypeName:
			p.types[x.Name()] = x
		case *types.Const:
			p.constants[x.Name()] = x
		}
	}

	return p
}

type pkgInfo struct {
	u Universe

	*packages.Package

	constants map[string]*types.Const
	types     map[string]*types.TypeName
	funcs     map[string]*types.Func
	methods   map[*types.Named][]*types.Func

	endLineToCommentGroup         map[fileLine]*ast.CommentGroup
	endLineToTrailingCommentGroup map[fileLine]*ast.CommentGroup

	signatures  map[*types.Signature]ast.Node
	funcResults map[*types.Signature][]TypeAndValues
}

func (pi *pkgInfo) SourceDir() string {
	if pi.PkgPath == pi.Module.Path {
		return pi.Module.Dir
	}
	return filepath.Join(pi.Module.Dir, pi.PkgPath[len(pi.Module.Path):])
}

func (pi *pkgInfo) Pkg() *types.Package {
	return pi.Package.Types
}

func (pi *pkgInfo) Files() []*ast.File {
	return pi.Package.Syntax
}

func (pi *pkgInfo) Eval(expr ast.Expr) (types.TypeAndValue, error) {
	return types.Eval(pi.Fset, pi.Package.Types, expr.Pos(), StringifyNode(pi.Fset, expr))
}

func (pi *pkgInfo) Constant(n string) *types.Const {
	return pi.constants[n]
}

func (pi *pkgInfo) Constants() map[string]*types.Const {
	return pi.constants
}

func (pi *pkgInfo) Type(n string) *types.TypeName {
	return pi.types[n]
}

func (pi *pkgInfo) Types() map[string]*types.TypeName {
	return pi.types
}

func (pi *pkgInfo) Function(n string) *types.Func {
	return pi.funcs[n]
}

func (pi *pkgInfo) Functions() map[string]*types.Func {
	return pi.funcs
}

func (pi *pkgInfo) MethodsOf(n *types.Named, ptr bool) []*types.Func {
	funcs, _ := pi.methods[n]

	if ptr {
		return funcs
	}

	notPtrMethods := make([]*types.Func, 0)

	for i := range funcs {
		s := funcs[i].Type().(*types.Signature)

		if _, ok := s.Recv().Type().(*types.Pointer); !ok {
			notPtrMethods = append(notPtrMethods, funcs[i])
		}
	}

	return notPtrMethods
}

func (pi *pkgInfo) Position(pos token.Pos) token.Position {
	return pi.Fset.Position(pos)
}

func (pi *pkgInfo) Doc(pos token.Pos) []string {
	c1 := pi.priorCommentLines(pos, 1)
	if c1 == nil {
		return commentLinesFrom(pi.priorCommentLines(pos, 2))
	}
	return commentLinesFrom(pi.priorCommentLines(pos, 1))
}

func (pi *pkgInfo) Comment(pos token.Pos) []string {
	return commentLinesFrom(pi.priorCommentLines(pos, 0))
}

func (pi *pkgInfo) priorCommentLines(pos token.Pos, lines int) *ast.CommentGroup {
	position := pi.Fset.Position(pos)
	key := fileLine{position.Filename, position.Line - lines}
	if lines != 0 {
		// should ignore trailing comments
		// when lines eq 0 means find trailing comments
		if _, ok := pi.endLineToTrailingCommentGroup[key]; ok {
			return nil
		}
	}
	return pi.endLineToCommentGroup[key]
}

type fileLine struct {
	file string
	line int
}

func commentLinesFrom(commentGroups ...*ast.CommentGroup) (comments []string) {
	if len(commentGroups) == 0 {
		return nil
	}

	for _, commentGroup := range commentGroups {
		if commentGroup == nil {
			continue
		}

		for _, line := range strings.Split(strings.TrimSpace(commentGroup.Text()), "\n") {
			// skip go: prefix
			if strings.HasPrefix(line, "go:") {
				continue
			}
			comments = append(comments, line)
		}
	}
	return comments
}
