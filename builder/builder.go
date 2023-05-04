package builder

import (
	"fmt"
	"strings"
	"sync"
)

func New() *Builder {
	return &Builder{nameToIDMap: map[string][]int{}}
}

type TypeBuilder interface {
	typemetadata() *TypeMetadata

	toSchemer  // to schema
	writeTyper // to string
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
	val := typ.typemetadata()
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
	name := typ.typemetadata().Name
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
func (t *TypeRef) typemetadata() *TypeMetadata {
	return t.getType().typemetadata()
}

func (b *Builder) Array(typ TypeBuilder) *ArrayType[TypeBuilder] { // TODO: specialized
	t := &ArrayType[TypeBuilder]{ArrayBuilder: &ArrayBuilder[TypeBuilder, *ArrayType[TypeBuilder]]{
		type_:    &type_[*ArrayType[TypeBuilder]]{builder: b, metadata: &TypeMetadata{Name: "array", underlying: "array"}},
		items:    typ,
		metadata: &ArrayMetadata{},
	}}
	t.ArrayBuilder.ret = t
	return t
}

type ArrayType[T TypeBuilder] struct {
	*ArrayBuilder[T, *ArrayType[T]]
}

type ArrayBuilder[T TypeBuilder, R TypeBuilder] struct {
	*type_[R]
	items    T
	metadata *ArrayMetadata
}

func (t *ArrayBuilder[T, R]) MinItems(n int64) R {
	t.metadata.MinItems = n
	return t.ret
}
func (t *ArrayBuilder[T, R]) MaxItems(n int64) R {
	t.metadata.MaxItems = n
	return t.ret
}

func (b *Builder) Map(valtyp TypeBuilder) *MapType[TypeBuilder] { // TODO: specialized
	t := &MapType[TypeBuilder]{MapBuilder: &MapBuilder[TypeBuilder, *MapType[TypeBuilder]]{
		type_:    &type_[*MapType[TypeBuilder]]{builder: b, metadata: &TypeMetadata{Name: "map[string]", underlying: "map[string]"}},
		items:    valtyp,
		metadata: &MapMetadata{},
	}}
	t.MapBuilder.ret = t
	return t
}

type MapType[T TypeBuilder] struct {
	*MapBuilder[T, *MapType[T]]
}

// string only map
type MapBuilder[V TypeBuilder, R TypeBuilder] struct {
	*type_[R]
	items    V
	metadata *MapMetadata
}

func (t *MapBuilder[T, R]) PatternProperties(s string, typ TypeBuilder) R {
	if t.metadata.PatternProperties == nil {
		t.metadata.PatternProperties = map[string]TypeBuilder{}
	}
	t.metadata.PatternProperties[s] = typ
	return t.ret
}

func (b *Builder) String() *StringType {
	t := &StringType{StringBuilder: &StringBuilder[*StringType]{
		type_:    &type_[*StringType]{builder: b, metadata: &TypeMetadata{Name: "string", underlying: "string"}},
		metadata: &StringMetadata{},
	}}
	t.StringBuilder.ret = t
	return t
}

type StringType struct {
	*StringBuilder[*StringType]
}

type StringBuilder[R TypeBuilder] struct {
	*type_[R]
	metadata *StringMetadata
}

var _ TypeBuilder = (*StringBuilder[TypeBuilder])(nil)

func (t *StringBuilder[R]) MinLength(n int64) R {
	t.metadata.MinLength = n
	return t.ret
}
func (t *StringBuilder[R]) MaxLength(n int64) R {
	t.metadata.MaxLength = n
	return t.ret
}
func (t *StringBuilder[R]) Pattern(s string) R {
	t.metadata.Pattern = s
	return t.ret
}

func (b *Builder) Integer() *IntegerType {
	t := &IntegerType{IntegerBuilder: &IntegerBuilder[*IntegerType]{
		type_:    &type_[*IntegerType]{builder: b, metadata: &TypeMetadata{Name: "integer", underlying: "integer"}},
		metadata: &IntegerMetadata{},
	}}
	t.IntegerBuilder.ret = t
	return t
}

type IntegerType struct {
	*IntegerBuilder[*IntegerType]
}

type IntegerBuilder[R TypeBuilder] struct {
	*type_[R]
	metadata *IntegerMetadata
}

var _ TypeBuilder = (*IntegerBuilder[TypeBuilder])(nil)

func (t *IntegerBuilder[R]) Minimum(n int64) R {
	t.metadata.Minimum = n
	return t.ret
}
func (t *IntegerBuilder[R]) Maximum(n int64) R {
	t.metadata.Maximum = n
	return t.ret
}

func (b *Builder) Object(fields ...*TypedField) *ObjectType {
	t := &ObjectType{
		ObjectBuilder: &ObjectBuilder[*ObjectType]{
			type_:    &type_[*ObjectType]{builder: b, metadata: &TypeMetadata{Name: "object", underlying: "object"}},
			Fields:   fields,
			metadata: &ObjectMetadata{},
		},
	}
	t.ObjectBuilder.ret = t
	return t
}

type ObjectType struct {
	*ObjectBuilder[*ObjectType]
}

type ObjectBuilder[R TypeBuilder] struct {
	*type_[R]
	metadata *ObjectMetadata
	Fields   []*TypedField
}

func (b *ObjectBuilder[R]) String(v bool) R {
	b.metadata.Strict = v
	return b.ret
}

func (b *Builder) Field(name string, typ TypeBuilder) *TypedField {
	f := &TypedField{
		field: &field[*TypedField]{
			metadata: &FieldMetadata{Name: name, Required: true},
		},
		typ: typ,
	}
	f.field.ret = f
	return f
}

type type_[R TypeBuilder] struct {
	metadata *TypeMetadata
	ret      R

	builder *Builder
}

func (t *type_[R]) typemetadata() *TypeMetadata {
	return t.metadata
}
func (t *type_[R]) Doc(stmts ...string) R {
	t.metadata.Description = strings.Join(stmts, "\n")
	return t.ret
}
func (t *type_[R]) Format(v string) R {
	t.metadata.Format = v
	return t.ret
}
func (t *type_[R]) As(name string) R {
	t.metadata.Name = name
	t.metadata.IsNewType = true
	t.builder.storeType(t.ret)
	return t.ret
}

type field[R any] struct {
	metadata *FieldMetadata
	ret      R
}

func (t *field[R]) Doc(stmts ...string) R {
	t.metadata.Description = strings.Join(stmts, "\n")
	return t.ret
}

func (t *field[R]) Required(v bool) R {
	t.metadata.Required = v
	return t.ret
}

type TypedField struct {
	*field[*TypedField]
	typ TypeBuilder
}
