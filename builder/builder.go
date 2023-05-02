package builder

import "strings"

func New() *Builder {
	return &Builder{}
}

type Builder struct {
	// TODO: storing types
}

type ObjectType struct {
	*ObjectBuilder[*ObjectType]
}

func (b *Builder) Object(name string, fields ...FieldBuilder) *ObjectType {
	t := &ObjectType{
		ObjectBuilder: &ObjectBuilder[*ObjectType]{
			type_: &type_{value: &Type{}},
			value: &Object{Fields: fields},
		},
	}
	t.ObjectBuilder.ret = t
	return t
}

func (b *Builder) Field(name string, typ TypeBuilder) *TypedField {
	f := &TypedField{
		field: &field[*TypedField]{
			value: &Field{
				Name: name,
			},
		},
		typ: typ,
	}
	f.field.ret = f
	return f
}

func (b *Builder) Array(typ TypeBuilder) *ArrayType[TypeBuilder] { // TODO: specialized
	t := &ArrayType[TypeBuilder]{ArrayBuilder: &ArrayBuilder[TypeBuilder, *ArrayType[TypeBuilder]]{
		type_: &type_{value: &Type{}},
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

func (b *Builder) String() *StringType {
	t := &StringType{StringBuilder: &StringBuilder[*StringType]{
		type_: &type_{value: &Type{}},
		value: &String{},
	}}
	t.StringBuilder.ret = t
	return t
}

type StringType struct {
	*StringBuilder[*StringType]
}

func (b *Builder) Integer() *IntegerType {
	t := &IntegerType{IntegerBuilder: &IntegerBuilder[*IntegerType]{
		type_: &type_{value: &Type{}},
		value: &Integer{},
	}}
	t.IntegerBuilder.ret = t
	return t
}

type IntegerType struct {
	*IntegerBuilder[*IntegerType]
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

type TypedField struct {
	*field[*TypedField]
	typ TypeBuilder
}

// https://swagger.io/docs/specification/data-models/data-types/

type StringBuilder[R any] struct {
	*type_
	value *String
	ret   R
}

var _ TypeBuilder = (*StringBuilder[any])(nil)

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

type IntegerBuilder[R any] struct {
	*type_
	value *Integer
	ret   R
}

var _ TypeBuilder = (*IntegerBuilder[any])(nil)

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
type ObjectBuilder[R any] struct {
	*type_
	value *Object
	ret   R
}

func (b *ObjectBuilder[R]) String(v bool) R {
	b.value.Strict = v
	return b.ret
}

type Object struct {
	Strict bool

	Fields []FieldBuilder
}

type ArrayBuilder[T TypeBuilder, R any] struct {
	*type_
	value *Array[T]
	ret   R
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

// string only map
type MapBuilder[T TypeBuilder, R any] struct {
	*type_
	value *Map[T]
	ret   R
}

func (t *MapBuilder[T, R]) PatternProperties(s string) R {
	t.value.PatternProperties = s
	return t.ret
}

type Map[T TypeBuilder] struct {
	PatternProperties string
	Items             T
}
