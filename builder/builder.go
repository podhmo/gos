package builder

import "strings"

func New() *Builder {
	return &Builder{}
}

type Builder struct {
	// TODO: storing types
}

func (b *Builder) Type(name string, fields ...Field) *Type {
	return &Type{
		impl: &TypeImpl{
			Name:   name,
			Fields: fields,
		},
	}
}

func (b *Builder) String(name string) *StringField {
	return &StringField{
		Field: &field[any]{
			Impl: &FieldImpl{
				Name: name,
				Type: "string",
			},
		},
	}
}
func (b *Builder) Integer(name string) *IntegerField {
	return &IntegerField{
		Field: &field[any]{
			Impl: &FieldImpl{
				Name: name,
				Type: "string",
			},
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
	Fields      []Field
	Description string
}

type FieldImpl struct {
	Name        string
	Type        string
	Description string
}

type Field interface {
	impl() *FieldImpl
}
type field[T any] struct {
	Impl *FieldImpl
}

func (f *field[T]) impl() *FieldImpl {
	return f.Impl
}

var _ Field = (*field[any])(nil)

// https://swagger.io/docs/specification/data-models/data-types/
type StringField struct {
	Field

	MinLength int64
	MaxLength int64
	Pattern   string
}
type IntegerField struct {
	Field

	// minimum ≤ value ≤ maximum
	Maximum int64
	Minimum int64
}
