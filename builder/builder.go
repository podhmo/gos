package builder

import "strings"

func New() *Builder {
	return &Builder{}
}

type Builder struct {
	// TODO: storing types
}

func (b *Builder) Type(name string, fields ...Field) *object {
	return &object{
		impl: &TypeImpl{
			Name:   name,
			Fields: fields,
		},
	}
}

func (b *Builder) Field(name string) *UntypedField {
	f := &UntypedField{
		field: &field[*UntypedField]{
			impl: &FieldImpl{
				Name: name,
			},
		},
	}
	f.field.retval = f
	return f
}

type Type interface {
	typeimpl() *TypeImpl
}

type object struct {
	impl *TypeImpl
}

func (t *object) typeimpl() *TypeImpl {
	return t.impl
}
func (t *object) Doc(stmts ...string) *object {
	t.impl.Description = strings.Join(stmts, "\n")
	return t
}
func (t *object) As(name string) *object {
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
	Required    bool
}

type Field interface {
	fieldimpl() *FieldImpl
}
type field[R any] struct {
	impl   *FieldImpl
	retval R
}

func (f *field[R]) fieldimpl() *FieldImpl {
	return f.impl
}
func (t *field[R]) Doc(stmts ...string) R {
	t.impl.Description = strings.Join(stmts, "\n")
	return t.retval
}

func (t *field[R]) Required(v bool) R {
	t.impl.Required = v
	return t.retval
}

var _ Field = (*field[any])(nil)

type UntypedField struct {
	*field[*UntypedField]
}

func (uf *UntypedField) String() *StringField {
	f := &StringField{
		field: &field[*StringField]{
			impl: uf.impl,
		},
		String: &String[*StringField]{
			typ:  &TypeImpl{},
			impl: &StringImpl{},
		},
	}
	f.field.retval = f
	f.String.retval = f
	return f
}

func (uf *UntypedField) Integer() *IntegerField {
	f := &IntegerField{
		field: &field[*IntegerField]{
			impl: uf.impl,
		},
		Integer: &Integer[*IntegerField]{
			typ:  &TypeImpl{},
			impl: &IntegerImpl{},
		},
	}
	f.field.retval = f
	f.Integer.retval = f
	return f
}

var _ Field = (*IntegerField)(nil)

// https://swagger.io/docs/specification/data-models/data-types/
type StringField struct {
	*field[*StringField]
	*String[*StringField]
}

var _ Field = (*StringField)(nil)

type String[R any] struct {
	typ    *TypeImpl
	impl   *StringImpl
	retval R
}

var _ Type = (*String[any])(nil)

func (t *String[R]) typeimpl() *TypeImpl {
	return t.typ
}
func (t *String[R]) MinLength(n int64) R {
	t.impl.MinLength = n
	return t.retval
}
func (t *String[R]) MaxLength(n int64) R {
	t.impl.MaxLength = n
	return t.retval
}
func (t *String[R]) Pattern(s string) R {
	t.impl.Pattern = s
	return t.retval
}

type StringImpl struct {
	MinLength int64
	MaxLength int64
	Pattern   string
}

type IntegerField struct {
	*field[*IntegerField]
	*Integer[*IntegerField]
}

var _ Field = (*IntegerField)(nil)

type Integer[R any] struct {
	typ    *TypeImpl
	impl   *IntegerImpl
	retval R
}

var _ Type = (*Integer[any])(nil)

func (t *Integer[R]) typeimpl() *TypeImpl {
	return t.typ
}
func (t *Integer[R]) Minimum(n int64) R {
	t.impl.Minimum = n
	return t.retval
}
func (t *Integer[R]) Maximum(n int64) R {
	t.impl.Maximum = n
	return t.retval
}

type IntegerImpl struct {
	// minimum ≤ value ≤ maximum
	Maximum int64
	Minimum int64
}
