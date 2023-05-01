package builder

import "strings"

func New() *Builder {
	return &Builder{}
}

type Builder struct {
	// TODO: storing types
}

func (b *Builder) Type(name string, fields ...*Field) *Type {
	return &Type{
		impl: &TypeImpl{
			Name:   name,
			Fields: fields,
		},
	}
}

func (b *Builder) Field(name string, typ *Type) *Field {
	return &Field{}
}

func (b *Builder) String(name string) *Field {
	return &Field{
		impl: &FieldImpl{
			Name: name,
			Type: "string",
		},
	}
}
func (b *Builder) Integer(name string) *Field {
	return &Field{
		impl: &FieldImpl{
			Name: name,
			Type: "integer",
		},
	}
}

type Type struct {
	impl *TypeImpl
}

func (t *Type) Doc(stmts ...string) *Type {
	t.impl.Description = strings.Join(stmts, "\n")
	return t
}
func (t *Type) As(name string) *Type {
	t.impl.Name = name
	return t
}

type TypeImpl struct {
	Name        string
	Fields      []*Field
	Description string
}

type FieldImpl struct {
	Name        string
	Type        string
	Description string
}

type Field struct {
	impl *FieldImpl
}

func (f *Field) Description(stmts ...string) *Field {
	f.impl.Description = strings.Join(stmts, "\n")
	return f
}

// https://swagger.io/docs/specification/data-models/data-types/
