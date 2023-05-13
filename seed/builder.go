package seed

import (
	"os"
	"strings"
)

type Builder struct {
	Metadata *BuilderMetadata
}

func NewBuilder(pkgname string) *Builder {
	return &Builder{
		Metadata: &BuilderMetadata{
			PkgName:     pkgname,
			SysArgs:     os.Args[1:],
			GeneratedBy: "github.com/podhmo/gos/seed",
		},
	}
}

type Symbol string

func (s Symbol) Pointer() Symbol {
	return "*" + s
}
func (s Symbol) Slice() Symbol {
	return "[]" + s
}

func (b *Builder) BuildTarget(name string) Symbol {
	b.Metadata.Target = Symbol(name)
	return Symbol(name)
}

func (b *Builder) InterfaceMethods(methods ...string) *Builder {
	b.Metadata.InterfaceMethods = append(b.Metadata.InterfaceMethods, methods...)
	return b
}

func (b *Builder) Import(path string) Symbol {
	b.Metadata.Imports = append(b.Metadata.Imports, Import{Path: path})
	parts := strings.Split(path, "/")
	return Symbol(parts[len(parts)-1])
}
func (b *Builder) NamedImport(name string, path string) Symbol {
	b.Metadata.Imports = append(b.Metadata.Imports, Import{Name: name, Path: path})
	return Symbol(name)
}

func (b *Builder) Field(name string, typ Symbol, tag string) *Builder {
	b.Metadata.Fields = append(b.Metadata.Fields, Field{Name: name, Type: typ, Tag: tag})
	return b
}
func (b *Builder) Constructor(args ...Arg) *Builder {
	b.Metadata.Constructor = &Constructor{Args: args}
	return b
}

func (b *Builder) Type(name string) *Type {
	t := &Type{
		TypeBuilder: &TypeBuilder[*Type]{Metadata: &TypeMetadata{
			Name:       Symbol(name),
			Underlying: name,
			Used:       map[string]bool{},
		}},
	}
	t.ret = t
	b.Metadata.Types = append(b.Metadata.Types, t)
	return t
}

type Type struct {
	*TypeBuilder[*Type]
}

type TypeBuilder[R any] struct {
	Metadata *TypeMetadata
	ret      R
}

func (b *TypeBuilder[R]) Field(name string, typ Symbol, tag string) R {
	b.Metadata.Fields = append(b.Metadata.Fields, Field{Name: name, Type: typ, Tag: tag})
	return b.ret
}

func (b *TypeBuilder[R]) NeedBuilder() R {
	b.Metadata.NeedBuilder = true
	return b.ret
}
func (b *TypeBuilder[R]) Underlying(v string) R {
	b.Metadata.Underlying = v
	return b.ret
}
func (b *TypeBuilder[R]) Constructor(args ...Arg) R {
	b.Metadata.Constructor = &Constructor{Args: args}
	for _, a := range args {
		b.Metadata.Used[a.Name] = true
	}
	return b.ret
}

// metadata
type BuilderMetadata struct {
	Target Symbol
	Types  []*Type

	Imports          []Import
	InterfaceMethods []string
	Constructor      *Constructor
	Fields           []Field // fields of Metadata

	SysArgs     []string // runtime os.Args[1:]
	PkgName     string   // package {{.PkgName}}}
	GeneratedBy string   // github.com/podhmo/gos/seed
}

type TypeMetadata struct {
	Name       Symbol
	Underlying string

	NeedBuilder bool
	Constructor *Constructor
	Fields      []Field // fields of Metadata

	Used map[string]bool
}

type Field struct {
	Name string
	Type Symbol
	Tag  string
}

type Constructor struct {
	Args []Arg
}

type Arg struct {
	Name     string
	Type     Symbol
	Variadic bool // as ...<type>
}

type Import struct {
	Name string
	Path string
}
