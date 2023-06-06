package design

import (
	"go/ast"
	"go/parser"
	"go/token"
	"runtime"
	"strings"

	"github.com/podhmo/gos/openapigen"
)

func NewHandler(b *openapigen.Builder) *Handler {
	return &Handler{b: b}
}

type Handler struct {
	b *openapigen.Builder
}

// greeting hello
func (h *Handler) Hello() *openapigen.Action {
	info := callerInfo()
	return b.Action(info.Name,
		b.Input(
			b.Param("name", b.String()).AsPath(),
		).Doc("input"),
		b.Output(
			b.String(),
		),
	).Doc(info.Doc)
}

// list person
func (h *Handler) ListPerson() *openapigen.Action {
	info := callerInfo()
	return b.Action(info.Name,
		b.Input(
			b.Param("sort", b.String().Enum([]string{"name", "-name", "age", "-age"})).AsQuery(),
		),
		b.Output(b.Array(PersonSummary)).Doc("list of person summary"),
	).Doc(info.Doc)
}

// create person
func (h *Handler) CreatePerson() *openapigen.Action {
	info := callerInfo()
	return b.Action(info.Name,
		b.Input(
			b.Param("verbose", b.Bool()).AsQuery(),
			b.Body(b.Object(
				append(Person.IgnoreFields("id", "father", "friends"),
					b.Field("fatherId", b.String()),
					b.Field("friendIdList", b.Array(b.String())))...,
			)).Doc("person but father and friends are id"),
		),
		b.Output(nil).Status(204),
	).Doc(info.Doc)
}

var (
	astCache = map[string]*ast.File{}
	fset     = token.NewFileSet()
)

func callerInfo() (info struct {
	Name   string
	File   string
	Lineno int
	Doc    string
}) {
	pc, _, _, _ := runtime.Caller(1)
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
	for _, decl := range f.Decls {
		decl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		if decl.Pos() <= pos && pos <= decl.End() {
			info.Doc = strings.TrimSpace(decl.Doc.Text())
			break
		}
	}
	parts := strings.Split(rfunc.Name(), ".")
	info.Name = parts[len(parts)-1]

	info.File = file
	info.Lineno = lineno
	return
}
