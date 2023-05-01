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
		field: &field[string]{
			impl: &FieldImpl{
				Name: name,
			},
		},
		impl: &StringFieldImpl{},
	}
}
func (b *Builder) Integer(name string) *IntegerField {
	return &IntegerField{
		field: &field[int]{
			impl: &FieldImpl{
				Name: name,
			},
		},
		impl: &IntegerFieldImpl{},
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
	Description string
}

type Field interface {
	fieldimpl() *FieldImpl
}
type field[T any] struct {
	impl *FieldImpl
}

func (f *field[T]) fieldimpl() *FieldImpl {
	return f.impl
}
func (t *field[T]) Doc(stmts ...string) *field[T] {
	t.impl.Description = strings.Join(stmts, "\n")
	return t
}

var _ Field = (*field[any])(nil)

// https://swagger.io/docs/specification/data-models/data-types/
type StringField struct {
	*field[string]
	impl *StringFieldImpl
}

var _ Field = (*StringField)(nil)

func (f *StringField) MinLength(n int64) *StringField {
	f.impl.MinLength = n
	return f
}
func (f *StringField) MaxLength(n int64) *StringField {
	f.impl.MaxLength = n
	return f
}
func (f *StringField) Pattern(s string) *StringField {
	f.impl.Pattern = s
	return f
}

type StringFieldImpl struct {
	MinLength int64
	MaxLength int64
	Pattern   string
}

type IntegerField struct {
	*field[int]
	impl *IntegerFieldImpl
}

func (f *IntegerField) Minimum(n int64) *IntegerField {
	f.impl.Minimum = n
	return f
}
func (f *IntegerField) Maximum(n int64) *IntegerField {
	f.impl.Maximum = n
	return f
}

type IntegerFieldImpl struct {
	// minimum ≤ value ≤ maximum
	Maximum int64
	Minimum int64
}
