package seed

type Builder struct {
}

type Symbol string

func (b *Builder) BuildTarget(name string) Symbol {
	return Symbol(name)
}

func (b *Builder) Type(name string) *Type {
	t := &Type{
		TypeBuilder: &TypeBuilder[*Type]{Metadata: &TypeMetadata{Name: Symbol(name)}},
	}
	t.ret = t
	return t
}

type Type struct {
	*TypeBuilder[*Type]
}

type TypeBuilder[R any] struct {
	Metadata *TypeMetadata
	ret      R
}

func (b *TypeBuilder[R]) NeedBuilder() R {
	b.Metadata.NeedBuilder = true
	return b.ret
}
func (b *TypeBuilder[R]) Var(name string, typ Symbol) R {
	b.Metadata.Vars = append(b.Metadata.Vars, Var{Name: name, Type: typ})
	return b.ret
}
func (b *TypeBuilder[R]) Field(name string, typ Symbol) R {
	b.Metadata.Fields = append(b.Metadata.Fields, Var{Name: name, Type: typ})
	return b.ret
}

func NewBuilder() *Builder {
	return &Builder{}
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
