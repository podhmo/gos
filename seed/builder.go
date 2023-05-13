package seed

import "os"

type Builder struct {
	Metadata *BuilderMetadata
}

func NewBuilder() *Builder {
	return &Builder{Metadata: &BuilderMetadata{Args: os.Args[1:]}}
}

type BuilderMetadata struct {
	Target Symbol
	Types  []*Type

	Args []string // runtime os.Args[1:]
}

type Symbol string

func (b *Builder) BuildTarget(name string) Symbol {
	b.Metadata.Target = Symbol(name)
	return Symbol(name)
}

func (b *Builder) Type(name string) *Type {
	t := &Type{
		TypeBuilder: &TypeBuilder[*Type]{Metadata: &TypeMetadata{Name: Symbol(name)}},
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

func (b *TypeBuilder[R]) Var(name string, typ Symbol) R {
	b.Metadata.Vars = append(b.Metadata.Vars, Var{Name: name, Type: typ})
	return b.ret
}
func (b *TypeBuilder[R]) Field(name string, typ Symbol) R {
	b.Metadata.Fields = append(b.Metadata.Fields, Var{Name: name, Type: typ})
	return b.ret
}

func (b *TypeBuilder[R]) NeedBuilder() R {
	b.Metadata.NeedBuilder = true
	return b.ret
}

// metadata

type TypeMetadata struct {
	Name Symbol

	NeedBuilder bool
	Vars        []Var
	Fields      []Var
}

type Var struct {
	Name string
	Type Symbol
}
