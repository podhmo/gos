package prototype

import (
	"fmt"
	"strings"
	"sync"
)

type TypeBuilder interface {
	GetTypeMetadata() *TypeMetadata

	toSchemer  // to schema
	writeTyper // to string
}

func Define[T interface {
	TypeBuilder
	storeType(name string)
}](name string, typ T) T {
	typ.storeType(name)
	return typ
}

type Builder struct {
	mu          sync.Mutex
	namedTypes  []TypeBuilder
	nameToIDMap map[string][]int
}

func NewBuilder() *Builder {
	return &Builder{nameToIDMap: map[string][]int{}}
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
	val := typ.GetTypeMetadata()
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
	return &TypeRef{Name: name, rootbuilder: b}
}
func (b *Builder) Reference(typ TypeBuilder) *TypeRef {
	name := typ.GetTypeMetadata().Name
	return &TypeRef{Name: name, rootbuilder: b, _typ: typ}
}

type TypeRef struct {
	Name string
	_typ TypeBuilder

	rootbuilder *Builder
}

func (t *TypeRef) getType() TypeBuilder {
	if t._typ != nil {
		return t._typ
	}
	t._typ = t.rootbuilder.lookupType(t.Name)
	return t._typ
}
func (t *TypeRef) GetTypeMetadata() *TypeMetadata {
	return t.getType().GetTypeMetadata()
}

func (b *Builder) Array(typ TypeBuilder) *ArrayType[TypeBuilder] { // TODO: specialized
	t := &ArrayType[TypeBuilder]{ArrayBuilder: &ArrayBuilder[TypeBuilder, *ArrayType[TypeBuilder]]{
		type_:    &type_[*ArrayType[TypeBuilder]]{rootbuilder: b, metadata: &TypeMetadata{Name: "array", underlying: "array"}},
		items:    typ,
		metadata: &ArrayMetadata{},
	}}
	t.ArrayBuilder.ret = t
	return t
}

type ArrayType[T TypeBuilder] struct {
	*ArrayBuilder[T, *ArrayType[T]]
}

func (t *ArrayType[T]) GetMetadata() *ArrayMetadata {
	return t.metadata
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
		type_:    &type_[*MapType[TypeBuilder]]{rootbuilder: b, metadata: &TypeMetadata{Name: "map[string]", underlying: "map[string]"}},
		items:    valtyp,
		metadata: &MapMetadata{},
	}}
	t.MapBuilder.ret = t
	return t
}

type MapType[T TypeBuilder] struct {
	*MapBuilder[T, *MapType[T]]
}

func (t *MapType[T]) GetMetadata() *MapMetadata {
	return t.metadata
}

// string only map
type MapBuilder[V TypeBuilder, R TypeBuilder] struct {
	*type_[R]
	items    V
	metadata *MapMetadata
}

func (t *MapBuilder[T, R]) Pattern(s string) R {
	t.metadata.Pattern = s
	return t.ret
}

func (b *Builder) String() *StringType {
	t := &StringType{StringBuilder: &StringBuilder[*StringType]{
		type_:    &type_[*StringType]{rootbuilder: b, metadata: &TypeMetadata{Name: "string", underlying: "string"}},
		metadata: &StringMetadata{},
	}}
	t.StringBuilder.ret = t
	return t
}

type StringType struct {
	*StringBuilder[*StringType]
}

func (t *StringType) GetMetadata() *StringMetadata {
	return t.metadata
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
		type_:    &type_[*IntegerType]{rootbuilder: b, metadata: &TypeMetadata{Name: "integer", underlying: "integer"}},
		metadata: &IntegerMetadata{},
	}}
	t.IntegerBuilder.ret = t
	return t
}

type IntegerType struct {
	*IntegerBuilder[*IntegerType]
}

func (t *IntegerType) GetMetadata() *IntegerMetadata {
	return t.metadata
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

func (b *Builder) Object(fields ...*Field) *ObjectType {
	t := &ObjectType{
		ObjectBuilder: &ObjectBuilder[*ObjectType]{
			type_:    &type_[*ObjectType]{rootbuilder: b, metadata: &TypeMetadata{Name: "object", underlying: "object"}},
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

func (t *ObjectType) GetMetadata() *ObjectMetadata {
	return t.metadata
}

type ObjectBuilder[R TypeBuilder] struct {
	*type_[R]
	metadata *ObjectMetadata
	Fields   []*Field
}

func (b *ObjectBuilder[R]) String(v bool) R {
	b.metadata.Strict = v
	return b.ret
}

func (b *Builder) Field(name string, typ TypeBuilder) *Field {
	f := &Field{
		fieldBuilder: &fieldBuilder[*Field]{
			metadata: &FieldMetadata{Name: name, Required: true},
		},
		typ: typ,
	}
	f.fieldBuilder.ret = f
	return f
}

type type_[R TypeBuilder] struct {
	metadata *TypeMetadata
	ret      R

	rootbuilder *Builder
}

func (t *type_[R]) GetTypeMetadata() *TypeMetadata {
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
func (t *type_[R]) storeType(name string) {
	t.metadata.Name = name
	t.metadata.IsNewType = true
	t.rootbuilder.storeType(t.ret)
}

type fieldBuilder[R any] struct {
	metadata *FieldMetadata
	ret      R
}

func (t *fieldBuilder[R]) Doc(stmts ...string) R {
	t.metadata.Description = strings.Join(stmts, "\n")
	return t.ret
}

func (t *fieldBuilder[R]) Required(v bool) R {
	t.metadata.Required = v
	return t.ret
}

type Field struct {
	*fieldBuilder[*Field]
	typ TypeBuilder
}

func (f *Field) GetFieldMetadata() *FieldMetadata {
	return f.metadata
}

// ----------------------------------------
func (b *Builder) Action(inputOrOutput ...actionSignature) *ActionType {
	t := &ActionType{
		ActionBuilder: &ActionBuilder[*ActionType]{
			type_:    &type_[*ActionType]{rootbuilder: b, metadata: &TypeMetadata{Name: "", underlying: "action"}}, // need?
			metadata: &ActionMetadata{},
		},
	}
	t.ret = t
	for _, sig := range inputOrOutput {
		switch sig := sig.(type) {
		case *ActionInput:
			t.input = sig
		case *ActionOutput:
			t.output = sig
		}
	}
	// TODO: conflict check?
	return t
}

type ActionType struct {
	*ActionBuilder[*ActionType]
}

type ActionBuilder[R TypeBuilder] struct {
	*type_[R]
	metadata *ActionMetadata
	input    *ActionInput
	output   *ActionOutput
}

func (b *Builder) Input(parameters ...*Param) *ActionInput {
	t := &ActionInput{
		ActionInputBuilder: &ActionInputBuilder[*ActionInput]{
			type_:    &type_[*ActionInput]{rootbuilder: b, metadata: &TypeMetadata{Name: "", underlying: "input"}}, // need?
			metadata: &ActionInputMetadata{},
			Params:   parameters,
		},
	}
	t.ret = t
	return t
}

type ActionInput struct {
	*ActionInputBuilder[*ActionInput]
}

type ActionInputBuilder[R TypeBuilder] struct {
	*type_[R]
	metadata *ActionInputMetadata
	Params   []*Param
}

func (t *ActionInput) sig() {}

func (b *Builder) Output(typ TypeBuilder) *ActionOutput {
	p := &Param{
		parameterBuilder: &parameterBuilder[*Param]{
			metadata: &ActionParamMetadata{Name: "", Required: true},
		},
		typ: typ,
	}
	p.ret = p

	t := &ActionOutput{
		ActionOutputBuilder: &ActionOutputBuilder[*ActionOutput]{
			type_:    &type_[*ActionOutput]{rootbuilder: b, metadata: &TypeMetadata{Name: "", underlying: "output"}}, // need?
			metadata: &ActionOutputMetadata{},
			retval:   p,
		},
	}
	t.ret = t
	return t
}

type ActionOutput struct {
	*ActionOutputBuilder[*ActionOutput]
}

type ActionOutputBuilder[R TypeBuilder] struct {
	*type_[R]
	metadata *ActionOutputMetadata
	retval   *Param
}

type parameterBuilder[R any] struct {
	metadata *ActionParamMetadata
	ret      R
}

func (t *parameterBuilder[R]) Doc(stmts ...string) R {
	t.metadata.Description = strings.Join(stmts, "\n")
	return t.ret
}

func (t *parameterBuilder[R]) Required(v bool) R {
	t.metadata.Required = v
	return t.ret
}

type Param struct {
	*parameterBuilder[*Param]
	typ TypeBuilder
}

func (f *Param) GetParamMetadata() *ActionParamMetadata {
	return f.metadata
}

func (t *ActionOutput) sig() {}

type actionSignature interface {
	sig()
}

func (b *Builder) Param(name string, typ TypeBuilder) *Param {
	f := &Param{
		parameterBuilder: &parameterBuilder[*Param]{
			metadata: &ActionParamMetadata{Name: name, Required: true},
		},
		typ: typ,
	}
	f.parameterBuilder.ret = f
	return f
}
