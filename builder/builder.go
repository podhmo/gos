package builder

import (
	"fmt"
	"io"
	"strings"
	"sync"
)

func New() *Builder {
	return &Builder{nameToIDMap: map[string][]int{}}
}

type TypeBuilder interface {
	typevalue() *Type
	WriteType(io.Writer) error
}

type Builder struct {
	mu          sync.Mutex
	namedTypes  []TypeBuilder
	nameToIDMap map[string][]int
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
	val.id = -1
	if !val.IsNewType {
		return
	}

	b.mu.Lock()
	defer b.mu.Unlock()
	id := len(b.namedTypes)
	val.id = id
	b.namedTypes = append(b.namedTypes, typ)
	b.nameToIDMap[val.Name] = append(b.nameToIDMap[val.Name], id)
	// TODO: name conflict check
}

func (b *Builder) lookupType(name string) TypeBuilder {
	b.mu.Lock()
	defer b.mu.Unlock()
	ids, ok := b.nameToIDMap[name]
	if !ok {
		b.nameToIDMap[name] = nil
		return nil
	}
	if len(ids) == 0 {
		return nil
	}

	// TODO: name conflict check
	return b.namedTypes[ids[0]]
}

func (b *Builder) ReferenceByName(name string) *TypeRef {
	return &TypeRef{Name: name, builder: b}
}
func (b *Builder) Reference(typ TypeBuilder) *TypeRef {
	name := typ.typevalue().Name
	return &TypeRef{Name: name, builder: b, _typ: typ}
}

type TypeRef struct {
	Name string
	_typ TypeBuilder

	builder *Builder
}

func (t *TypeRef) getType() TypeBuilder {
	if t._typ != nil {
		return t._typ
	}
	t._typ = t.builder.lookupType(t.Name)
	return t._typ
}
func (t *TypeRef) typevalue() *Type {
	return t.getType().typevalue()
}
func (t *TypeRef) WriteType(w io.Writer) error {
	return t.getType().WriteType(w)
}

func (b *Builder) Object(fields ...*TypedField) *ObjectType {
	t := &ObjectType{
		ObjectBuilder: &ObjectBuilder[*ObjectType]{
			type_:  &type_[*ObjectType]{builder: b, value: &Type{Name: "object", underlying: "object"}},
			Fields: fields,
			value:  &Object{},
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
	Name        string `json:"-"`
	Description string `json:"description,omitempty"`
	Format      string `json:"format,omitempty"`

	IsNewType bool `json:"-"`

	underlying string `json:"-"`
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

func (t *field[R]) Doc(stmts ...string) R {
	t.value.Description = strings.Join(stmts, "\n")
	return t.ret
}

func (t *field[R]) Required(v bool) R {
	t.value.Required = v
	return t.ret
}

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
	value  *Object
	Fields []*TypedField
}

func (b *ObjectBuilder[R]) String(v bool) R {
	b.value.Strict = v
	return b.ret
}

type Object struct {
	Strict bool
}

type ArrayBuilder[T TypeBuilder, R TypeBuilder] struct {
	*type_[R]
	items T
	value *Array
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

func (t *MapBuilder[T, R]) PatternProperties(s string) R {
	t.value.PatternProperties = s
	return t.ret
}

type Map struct {
	PatternProperties string
}
