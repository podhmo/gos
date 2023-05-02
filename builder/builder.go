package builder

import (
	"fmt"
	"io"
	"strings"
	"sync"
)

func New() *Builder {
	return &Builder{}
}

type TypeBuilder interface {
	typevalue() *Type
	WriteType(io.Writer) error
}
type FieldBuilder interface {
	fieldvalue() *Field
}

type Builder struct {
	mu         sync.Mutex
	namedTypes []TypeBuilder
}

func (b *Builder) EachTypes(fn func(TypeBuilder) error) error {
	for _, t := range b.namedTypes {
		if err := fn(t); err != nil {
			return fmt.Errorf("error on %s -- %w", ToString(t), err)
		}
	}
	return nil
}

func (b *Builder) storeType(typ TypeBuilder) {
	val := typ.typevalue()
	if !val.IsNewType {
		return
	}

	b.mu.Lock()
	defer b.mu.Unlock()
	id := len(b.namedTypes) + 1
	val.id = id
	b.namedTypes = append(b.namedTypes, typ)
	// TODO: name conflict check
}

func (b *Builder) Object(fields ...FieldBuilder) *ObjectType {
	t := &ObjectType{
		ObjectBuilder: &ObjectBuilder[*ObjectType]{
			type_: &type_[*ObjectType]{builder: b, value: &Type{Name: "object", underlying: "object"}},
			value: &Object{Fields: fields},
		},
	}
	t.ObjectBuilder.ret = t
	return t
}

type ObjectType struct {
	*ObjectBuilder[*ObjectType]
}

func (b *Builder) Field(name string, typ TypeBuilder) *TypedField {
	f := &TypedField{
		field: &field[*TypedField]{
			value: &Field{Name: name, Required: true},
		},
		typ: typ,
	}
	f.field.ret = f
	return f
}

func (b *Builder) Array(typ TypeBuilder) *ArrayType[TypeBuilder] { // TODO: specialized
	t := &ArrayType[TypeBuilder]{ArrayBuilder: &ArrayBuilder[TypeBuilder, *ArrayType[TypeBuilder]]{
		type_: &type_[*ArrayType[TypeBuilder]]{builder: b, value: &Type{Name: "array", underlying: "array"}},
		items: typ,
		value: &Array{},
	}}
	t.ArrayBuilder.ret = t
	return t
}

type ArrayType[T TypeBuilder] struct {
	*ArrayBuilder[T, *ArrayType[T]]
}

func (b *Builder) Map(valtyp TypeBuilder) *MapType[TypeBuilder] { // TODO: specialized
	t := &MapType[TypeBuilder]{MapBuilder: &MapBuilder[TypeBuilder, *MapType[TypeBuilder]]{
		type_: &type_[*MapType[TypeBuilder]]{builder: b, value: &Type{Name: "map[string]", underlying: "map[string]"}},
		items: valtyp,
		value: &Map{},
	}}
	t.MapBuilder.ret = t
	return t
}

type MapType[T TypeBuilder] struct {
	*MapBuilder[T, *MapType[T]]
}

func (b *Builder) String() *StringType {
	t := &StringType{StringBuilder: &StringBuilder[*StringType]{
		type_: &type_[*StringType]{builder: b, value: &Type{Name: "string", underlying: "string"}},
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
		type_: &type_[*IntegerType]{builder: b, value: &Type{Name: "integer", underlying: "integer"}},
		value: &Integer{},
	}}
	t.IntegerBuilder.ret = t
	return t
}

type IntegerType struct {
	*IntegerBuilder[*IntegerType]
}

type type_[R TypeBuilder] struct {
	value *Type
	ret   R

	builder *Builder
}

func (t *type_[R]) typevalue() *Type {
	return t.value
}
func (t *type_[R]) WriteType(w io.Writer) error {
	if _, err := io.WriteString(w, t.value.Name); err != nil {
		return err
	}
	return nil
}
func (t *type_[R]) Doc(stmts ...string) R {
	t.value.Description = strings.Join(stmts, "\n")
	return t.ret
}
func (t *type_[R]) Format(v string) R {
	t.value.Format = v
	return t.ret
}
func (t *type_[R]) As(name string) R {
	t.value.Name = name
	t.value.IsNewType = true
	t.builder.storeType(t.ret)
	return t.ret
}

type Type struct {
	id          int
	Name        string
	Description string
	Format      string

	IsNewType bool

	underlying string
}

type Field struct {
	Name        string
	Description string
	Required    bool
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

var _ FieldBuilder = (*field[TypeBuilder])(nil)

type TypedField struct {
	*field[*TypedField]
	typ TypeBuilder
}

// https://swagger.io/docs/specification/data-models/data-types/

type StringBuilder[R TypeBuilder] struct {
	*type_[R]
	value *String
}

var _ TypeBuilder = (*StringBuilder[TypeBuilder])(nil)

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

type IntegerBuilder[R TypeBuilder] struct {
	*type_[R]
	value *Integer
}

var _ TypeBuilder = (*IntegerBuilder[TypeBuilder])(nil)

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
type ObjectBuilder[R TypeBuilder] struct {
	*type_[R]
	value *Object
}

func (b ObjectBuilder[R]) WriteType(w io.Writer) error {
	if err := b.type_.WriteType(w); err != nil {
		return err
	}
	if b.type_.value.IsNewType {
		return nil
	}

	io.WriteString(w, "{") // nolint
	n := len(b.value.Fields) - 1
	for i, f := range b.value.Fields {
		v := f.fieldvalue()
		io.WriteString(w, v.Name) // nolint
		if !v.Required {
			io.WriteString(w, "?") // nolint
		}
		if i < n {
			io.WriteString(w, ", ") // nolint
		}
	}
	io.WriteString(w, "}") // nolint
	return nil
}

func (b *ObjectBuilder[R]) String(v bool) R {
	b.value.Strict = v
	return b.ret
}

type Object struct {
	Strict bool

	Fields []FieldBuilder
}

type ArrayBuilder[T TypeBuilder, R TypeBuilder] struct {
	*type_[R]
	items T
	value *Array
}

func (t *ArrayBuilder[T, R]) WriteType(w io.Writer) error {
	if err := t.type_.WriteType(w); err != nil {
		return err
	}
	if t.type_.value.IsNewType {
		return nil
	}

	io.WriteString(w, "[") // nolint
	if err := t.items.WriteType(w); err != nil {
		return err
	}
	io.WriteString(w, "]") // nolint
	return nil
}

func (t *ArrayBuilder[T, R]) MinItems(n int64) R {
	t.value.MinItems = n
	return t.ret
}
func (t *ArrayBuilder[T, R]) MaxItems(n int64) R {
	t.value.MaxItems = n
	return t.ret
}

type Array struct {
	MaxItems int64
	MinItems int64
}

// string only map
type MapBuilder[V TypeBuilder, R TypeBuilder] struct {
	*type_[R]
	items V
	value *Map
}

func (t *MapBuilder[V, R]) WriteType(w io.Writer) error {
	if err := t.type_.WriteType(w); err != nil {
		return err
	}
	if t.type_.value.IsNewType {
		return nil
	}

	io.WriteString(w, "[") // nolint
	if err := t.items.WriteType(w); err != nil {
		return err
	}
	io.WriteString(w, "]") // nolint
	return nil
}

func (t *MapBuilder[T, R]) PatternProperties(s string) R {
	t.value.PatternProperties = s
	return t.ret
}

type Map struct {
	PatternProperties string
}

func ToString(typ TypeBuilder) string {
	b := new(strings.Builder)
	if err := typ.WriteType(b); err != nil {
		return fmt.Sprintf("invalid type: %T", typ)
	}
	return b.String()
}
