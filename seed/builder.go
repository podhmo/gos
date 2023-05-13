package seed

import (
	"os"
)

type Builder struct {
	Metadata *BuilderMetadata
}

func NewBuilder(pkgname string) *Builder {
	return &Builder{
		Metadata: &BuilderMetadata{Args: os.Args[1:], PkgName: pkgname},
	}
}

type BuilderMetadata struct {
	Target Symbol
	Types  []*Type

	Args    []string // runtime os.Args[1:]
	PkgName string   // package {{.PkgName}}}
}

type Symbol string

func (b *Builder) BuildTarget(name string) Symbol {
	b.Metadata.Target = Symbol(name)
	return Symbol(name)
}

func (b *Builder) Type(name string) *Type {
	t := &Type{
		TypeBuilder: &TypeBuilder[*Type]{Metadata: &TypeMetadata{
			Name: Symbol(name),
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
	b.Metadata.Fields = append(b.Metadata.Fields, Var{Name: name, Type: typ, Tag: tag})
	return b.ret
}

func (b *TypeBuilder[R]) NeedBuilder() R {
	b.Metadata.NeedBuilder = true
	return b.ret
}
func (b *TypeBuilder[R]) Constructor(args ...Arg) R {
	b.Metadata.Constructor = &Constructor{Args: args}
	return b.ret
}

// metadata

type TypeMetadata struct {
	Name Symbol

	NeedBuilder bool
	Constructor *Constructor
	Fields      []Var // fields of Metadata
}

type Var struct {
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
