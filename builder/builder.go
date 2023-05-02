package builder

import "strings"

func New() *Builder {
	return &Builder{}
}

type Builder struct {
	// TODO: storing types
}

func (b *Builder) Object(name string, fields ...FieldBuilder) *type_ {
	return &type_{
		value: &Type{
			Name:   name,
			Fields: fields,
		},
	}
}

func (b *Builder) Field(name string) *UntypedField {
	f := &UntypedField{
		field: &field[*UntypedField]{
			value: &Field{
				Name: name,
			},
		},
	}
	f.field.ret = f
	return f
}

func (b *Builder) Array(typ TypeBuilder) *ArrayType[TypeBuilder] { // TODO: specialized
	t := &ArrayType[TypeBuilder]{ArrayBuilder: &ArrayBuilder[TypeBuilder, *ArrayType[TypeBuilder]]{
		typ: &Type{},
		value: &Array[TypeBuilder]{
			Items: typ,
		},
	}}
	t.ArrayBuilder.ret = t
	return t
}

type ArrayType[T TypeBuilder] struct {
	*ArrayBuilder[T, *ArrayType[T]]
}

func (t *ArrayType[T]) typevalue() *Type {
	return t.ArrayBuilder.typ
}

func (b *Builder) String() *StringType {
	t := &StringType{StringBuilder: &StringBuilder[*StringType]{
		typ:   &Type{},
		value: &String{},
	}}
	t.StringBuilder.ret = t
	return t
}

type StringType struct {
	*StringBuilder[*StringType]
}

func (t *StringType) typevalue() *Type {
	return t.StringBuilder.typ
}
func (b *Builder) Integer() *IntegerType {
	t := &IntegerType{IntegerBuilder: &IntegerBuilder[*IntegerType]{
		typ:   &Type{},
		value: &Integer{},
	}}
	t.IntegerBuilder.ret = t
	return t
}

type IntegerType struct {
	*IntegerBuilder[*IntegerType]
}

func (t *IntegerType) typevalue() *Type {
	return t.IntegerBuilder.typ
}

type TypeBuilder interface {
	typevalue() *Type
}

type type_ struct {
	value *Type
}

func (t *type_) typevalue() *Type {
	return t.value
}
func (t *type_) Doc(stmts ...string) *type_ {
	t.value.Description = strings.Join(stmts, "\n")
	return t
}
func (t *type_) As(name string) *type_ {
	t.value.Name = name
	return t
}

type Type struct {
	Name        string
	Fields      []FieldBuilder
	Description string
}

type Field struct {
	Name        string
	Description string
	Required    bool
}

type FieldBuilder interface {
	fieldvalue() *Field
}
type field[R any] struct {
	value *Field
	ret   R
}

func (f *field[R]) fieldvalue() *Field {
	return f.value
}
func (t *field[R]) Doc(stmts ...string) R {
	t.value.Description = strings.Join(stmts, "\n")
	return t.ret
}

func (t *field[R]) Required(v bool) R {
	t.value.Required = v
	return t.ret
}

var _ FieldBuilder = (*field[any])(nil)

type UntypedField struct {
	*field[*UntypedField]
}

func (uf *UntypedField) String() *StringField {
	f := &StringField{
		field: &field[*StringField]{
			value: uf.value,
		},
		StringBuilder: &StringBuilder[*StringField]{
			typ:   &Type{},
			value: &String{},
		},
	}
	f.field.ret = f
	f.StringBuilder.ret = f
	return f
}

func (uf *UntypedField) Integer() *IntegerField {
	f := &IntegerField{
		field: &field[*IntegerField]{
			value: uf.value,
		},
		IntegerBuilder: &IntegerBuilder[*IntegerField]{
			typ:   &Type{},
			value: &Integer{},
		},
	}
	f.field.ret = f
	f.IntegerBuilder.ret = f
	return f
}

var _ FieldBuilder = (*IntegerField)(nil)

// https://swagger.io/docs/specification/data-models/data-types/
type StringField struct {
	*field[*StringField]
	*StringBuilder[*StringField]
}

var _ FieldBuilder = (*StringField)(nil)

type StringBuilder[R any] struct {
	typ   *Type
	value *String
	ret   R
}

var _ TypeBuilder = (*StringBuilder[any])(nil)

func (t *StringBuilder[R]) typevalue() *Type {
	return t.typ
}
func (t *StringBuilder[R]) MinLength(n int64) R {
	t.value.MinLength = n
	return t.ret
}
func (t *StringBuilder[R]) MaxLength(n int64) R {
	t.value.MaxLength = n
	return t.ret
}
func (t *StringBuilder[R]) Pattern(s string) R {
	t.value.Pattern = s
	return t.ret
}

type String struct {
	MinLength int64
	MaxLength int64
	Pattern   string
}

type IntegerField struct {
	*field[*IntegerField]
	*IntegerBuilder[*IntegerField]
}

var _ FieldBuilder = (*IntegerField)(nil)

type IntegerBuilder[R any] struct {
	typ   *Type
	value *Integer
	ret   R
}

var _ TypeBuilder = (*IntegerBuilder[any])(nil)

func (t *IntegerBuilder[R]) typevalue() *Type {
	return t.typ
}
func (t *IntegerBuilder[R]) Minimum(n int64) R {
	t.value.Minimum = n
	return t.ret
}
func (t *IntegerBuilder[R]) Maximum(n int64) R {
	t.value.Maximum = n
	return t.ret
}

type Integer struct {
	// minimum ≤ value ≤ maximum
	Maximum int64
	Minimum int64
}

// composite type
type ArrayBuilder[T TypeBuilder, R any] struct {
	typ   *Type
	value *Array[T]
	ret   R
}

func (t *ArrayBuilder[T, R]) typevalue() *Type {
	return t.typ
}
func (t *ArrayBuilder[T, R]) MinItems(n int64) R {
	t.value.MinItems = n
	return t.ret
}
func (t *ArrayBuilder[T, R]) MaxItems(n int64) R {
	t.value.MaxItems = n
	return t.ret
}

type Array[T TypeBuilder] struct {
	MaxItems int64
	MinItems int64

	Items T
}
