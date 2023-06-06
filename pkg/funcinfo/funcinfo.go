package funcinfo

import (
	"go/ast"
	"go/parser"
	"go/token"
	"runtime"
	"strings"
)

type Collector struct {
	astCache map[string]*ast.File
	Fset     *token.FileSet
	Depth    int
}

func NewCollector(depth int) *Collector {
	return &Collector{
		astCache: map[string]*ast.File{},
		Fset:     token.NewFileSet(),
		Depth:    depth,
	}
}

type FuncInfo struct {
	FuncName string
	File     string
	Lineno   int
	FuncDoc  string
}

func (c *Collector) FuncName() string {
	depth := c.Depth

	pc, _, _, _ := runtime.Caller(depth)
	rfunc := runtime.FuncForPC(pc)
	parts := strings.Split(rfunc.Name(), ".")
	return parts[len(parts)-1]
}

func (c *Collector) Info() FuncInfo {
	fset := c.Fset
	astCache := c.astCache
	depth := c.Depth

	pc, _, _, _ := runtime.Caller(depth)
	rfunc := runtime.FuncForPC(pc)
	file, lineno := rfunc.FileLine(pc)
	f, ok := astCache[file]
	if !ok {
		tree, err := parser.ParseFile(fset, file, nil, parser.ParseComments|parser.SkipObjectResolution)
		if err != nil {
			panic(err)
		}
		astCache[file] = tree
		f = tree
	}

	pos := fset.File(f.Pos()).LineStart(lineno)
	doc := ""
	for _, decl := range f.Decls {
		decl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		if decl.Pos() <= pos && pos <= decl.End() {
			doc = strings.TrimSpace(decl.Doc.Text())
			break
		}
	}
	parts := strings.Split(rfunc.Name(), ".")
	return FuncInfo{
		FuncName: parts[len(parts)-1],
		File:     file,
		Lineno:   lineno,
		FuncDoc:  doc,
	}
}
