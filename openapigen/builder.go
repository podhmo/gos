// Generated by github.com/podhmo/gos/openapigen/tools [-write -builder -metadata -stringer -pkgname openapigen]

package openapigen

import (
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/iancoleman/orderedmap"
)

type TypeBuilder interface {
	GetTypeMetadata() *TypeMetadata
	writeType(io.Writer) error // see: ./stringer.go
	toSchemer                  // see: ./to_schema.go
}

// Define names Type value.
func Define[T interface {
	TypeBuilder
	storeType(name string)
}](name string, t T) T {
	t.storeType(name)
	return t
}

type Builder struct {
	mu          sync.Mutex
	namedTypes  []TypeBuilder
	nameToIDMap map[string][]int

	Config *Config
}

func NewBuilder(config *Config) *Builder {
	return &Builder{
		nameToIDMap: map[string][]int{},
		namedTypes:  []TypeBuilder{nil}, // nil is sentinel (id<=0 is unnamed)
		Config:      config,
	}
}

// EachType iterates named Type.
func (b *Builder) EachTypes(fn func(TypeBuilder) error) error {
	for _, t := range b.namedTypes {
		if t == nil {
			continue
		}
		if err := fn(t); err != nil {
			return fmt.Errorf("error on %v -- %w", t, err) // TODO: use ToString()
		}
	}
	return nil
}

func (b *Builder) storeType(typ TypeBuilder) {
	val := typ.GetTypeMetadata()
	val.id = -1
	if val.Name == "" {
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
	return &TypeRef{Name: name, rootbuilder: b, _Type: typ}
}

type TypeRef struct {
	Name  string
	_Type TypeBuilder

	rootbuilder *Builder
}

func (t *TypeRef) getType() TypeBuilder {
	if t._Type != nil {
		return t._Type
	}
	t._Type = t.rootbuilder.lookupType(t.Name)
	return t._Type
}
func (t *TypeRef) GetTypeMetadata() *TypeMetadata {
	return t.getType().GetTypeMetadata()
}

// Bool builds Type for Bool
func (b *Builder) Bool() *Bool {
	t := &Bool{
		_BoolBuilder: &_BoolBuilder[*Bool]{
			_Type:    &_Type[*Bool]{rootbuilder: b, metadata: &TypeMetadata{Name: "", underlying: "boolean", goType: "bool"}},
			metadata: &BoolMetadata{},
		},
	}
	t.ret = t
	return t
}

type Bool struct {
	*_BoolBuilder[*Bool]
}

func (t *Bool) GetMetadata() *BoolMetadata {
	return t.metadata
}

type _BoolBuilder[R TypeBuilder] struct {
	*_Type[R]
	metadata *BoolMetadata
}

// begin setter of Bool --------------------

// Format set Metadata.Format
func (b *_BoolBuilder[R]) Format(value string) R {
	b.metadata.Format = value
	return b.ret
}

// Default set Metadata.Default
func (b *_BoolBuilder[R]) Default(value bool) R {
	b.metadata.Default = value
	return b.ret
}

// end setter of Bool --------------------

// Int builds Type for Int
func (b *Builder) Int() *Int {
	t := &Int{
		_IntBuilder: &_IntBuilder[*Int]{
			_Type:    &_Type[*Int]{rootbuilder: b, metadata: &TypeMetadata{Name: "", underlying: "integer", goType: "int64"}},
			metadata: &IntMetadata{},
		},
	}
	t.ret = t
	return t
}

type Int struct {
	*_IntBuilder[*Int]
}

func (t *Int) GetMetadata() *IntMetadata {
	return t.metadata
}

type _IntBuilder[R TypeBuilder] struct {
	*_Type[R]
	metadata *IntMetadata
}

// begin setter of Int --------------------

// Format set Metadata.Format
func (b *_IntBuilder[R]) Format(value string) R {
	b.metadata.Format = value
	return b.ret
}

// Enum set Metadata.Enum
func (b *_IntBuilder[R]) Enum(value []int64) R {
	b.metadata.Enum = value
	return b.ret
}

// Default set Metadata.Default
func (b *_IntBuilder[R]) Default(value int64) R {
	b.metadata.Default = value
	return b.ret
}

// Maximum set Metadata.Maximum
func (b *_IntBuilder[R]) Maximum(value int64) R {
	b.metadata.Maximum = value
	return b.ret
}

// Minimum set Metadata.Minimum
func (b *_IntBuilder[R]) Minimum(value int64) R {
	b.metadata.Minimum = value
	return b.ret
}

// ExclusiveMin set Metadata.ExclusiveMin
func (b *_IntBuilder[R]) ExclusiveMin(value bool) R {
	b.metadata.ExclusiveMin = value
	return b.ret
}

// ExclusiveMax set Metadata.ExclusiveMax
func (b *_IntBuilder[R]) ExclusiveMax(value bool) R {
	b.metadata.ExclusiveMax = value
	return b.ret
}

// end setter of Int --------------------

// Float builds Type for Float
func (b *Builder) Float() *Float {
	t := &Float{
		_FloatBuilder: &_FloatBuilder[*Float]{
			_Type:    &_Type[*Float]{rootbuilder: b, metadata: &TypeMetadata{Name: "", underlying: "number", goType: "float64"}},
			metadata: &FloatMetadata{},
		},
	}
	t.ret = t
	return t
}

type Float struct {
	*_FloatBuilder[*Float]
}

func (t *Float) GetMetadata() *FloatMetadata {
	return t.metadata
}

type _FloatBuilder[R TypeBuilder] struct {
	*_Type[R]
	metadata *FloatMetadata
}

// begin setter of Float --------------------

// Format set Metadata.Format
func (b *_FloatBuilder[R]) Format(value string) R {
	b.metadata.Format = value
	return b.ret
}

// Default set Metadata.Default
func (b *_FloatBuilder[R]) Default(value string) R {
	b.metadata.Default = value
	return b.ret
}

// Maximum set Metadata.Maximum
func (b *_FloatBuilder[R]) Maximum(value float64) R {
	b.metadata.Maximum = value
	return b.ret
}

// Minimum set Metadata.Minimum
func (b *_FloatBuilder[R]) Minimum(value float64) R {
	b.metadata.Minimum = value
	return b.ret
}

// MultipleOf set Metadata.MultipleOf
func (b *_FloatBuilder[R]) MultipleOf(value float64) R {
	b.metadata.MultipleOf = value
	return b.ret
}

// ExclusiveMin set Metadata.ExclusiveMin
func (b *_FloatBuilder[R]) ExclusiveMin(value bool) R {
	b.metadata.ExclusiveMin = value
	return b.ret
}

// ExclusiveMax set Metadata.ExclusiveMax
func (b *_FloatBuilder[R]) ExclusiveMax(value bool) R {
	b.metadata.ExclusiveMax = value
	return b.ret
}

// end setter of Float --------------------

// String builds Type for String
func (b *Builder) String() *String {
	t := &String{
		_StringBuilder: &_StringBuilder[*String]{
			_Type:    &_Type[*String]{rootbuilder: b, metadata: &TypeMetadata{Name: "", underlying: "string", goType: "string"}},
			metadata: &StringMetadata{},
		},
	}
	t.ret = t
	return t
}

type String struct {
	*_StringBuilder[*String]
}

func (t *String) GetMetadata() *StringMetadata {
	return t.metadata
}

type _StringBuilder[R TypeBuilder] struct {
	*_Type[R]
	metadata *StringMetadata
}

// begin setter of String --------------------

// Format set Metadata.Format
func (b *_StringBuilder[R]) Format(value string) R {
	b.metadata.Format = value
	return b.ret
}

// Enum set Metadata.Enum
func (b *_StringBuilder[R]) Enum(value []string) R {
	b.metadata.Enum = value
	return b.ret
}

// Default set Metadata.Default
func (b *_StringBuilder[R]) Default(value string) R {
	b.metadata.Default = value
	return b.ret
}

// Pattern set Metadata.Pattern
func (b *_StringBuilder[R]) Pattern(value string) R {
	b.metadata.Pattern = value
	return b.ret
}

// MaxLength set Metadata.MaxLength
func (b *_StringBuilder[R]) MaxLength(value int64) R {
	b.metadata.MaxLength = value
	return b.ret
}

// MinLength set Metadata.MinLength
func (b *_StringBuilder[R]) MinLength(value int64) R {
	b.metadata.MinLength = value
	return b.ret
}

// end setter of String --------------------

// Array builds Type for Array
func (b *Builder) Array(items Type) *Array[Type] {
	t := &Array[Type]{
		_ArrayBuilder: &_ArrayBuilder[Type, *Array[Type]]{
			_Type:    &_Type[*Array[Type]]{rootbuilder: b, metadata: &TypeMetadata{Name: "", underlying: "array", goType: "Array"}},
			metadata: &ArrayMetadata{},
			items:    items,
		},
	}
	t.ret = t
	return t
}

type Array[Items Type] struct {
	*_ArrayBuilder[Items, *Array[Items]]
}

func (t *Array[Items]) GetMetadata() *ArrayMetadata {
	return t.metadata
}

type _ArrayBuilder[Items Type, R TypeBuilder] struct {
	*_Type[R]
	metadata *ArrayMetadata
	items    Items
}

// begin setter of Array --------------------

// MaxItems set Metadata.MaxItems
func (b *_ArrayBuilder[Items, R]) MaxItems(value int64) R {
	b.metadata.MaxItems = value
	return b.ret
}

// MinItems set Metadata.MinItems
func (b *_ArrayBuilder[Items, R]) MinItems(value int64) R {
	b.metadata.MinItems = value
	return b.ret
}

// end setter of Array --------------------

// Map builds Type for Map
func (b *Builder) Map(items Type) *Map[Type] {
	t := &Map[Type]{
		_MapBuilder: &_MapBuilder[Type, *Map[Type]]{
			_Type:    &_Type[*Map[Type]]{rootbuilder: b, metadata: &TypeMetadata{Name: "", underlying: "map", goType: "Map"}},
			metadata: &MapMetadata{},
			items:    items,
		},
	}
	t.ret = t
	return t
}

type Map[Items Type] struct {
	*_MapBuilder[Items, *Map[Items]]
}

func (t *Map[Items]) GetMetadata() *MapMetadata {
	return t.metadata
}

type _MapBuilder[Items Type, R TypeBuilder] struct {
	*_Type[R]
	metadata *MapMetadata
	items    Items
}

// begin setter of Map --------------------

// Pattern set Metadata.Pattern
func (b *_MapBuilder[Items, R]) Pattern(value string) R {
	b.metadata.Pattern = value
	return b.ret
}

// end setter of Map --------------------

// Field builds Type for Field
func (b *Builder) Field(name string, typ Type) *Field {
	t := &Field{
		_FieldBuilder: &_FieldBuilder[*Field]{
			_Type: &_Type[*Field]{rootbuilder: b, metadata: &TypeMetadata{Name: "", underlying: "field", goType: "Field"}},
			metadata: &FieldMetadata{
				Name: name, Typ: typ,
				Required: true,
			},
		},
	}
	t.ret = t
	return t
}

type Field struct {
	*_FieldBuilder[*Field]
}

func (t *Field) GetMetadata() *FieldMetadata {
	return t.metadata
}

type _FieldBuilder[R TypeBuilder] struct {
	*_Type[R]
	metadata *FieldMetadata
}

// begin setter of Field --------------------

// Nullable set Metadata.Nullable
func (b *_FieldBuilder[R]) Nullable(value bool) R {
	b.metadata.Nullable = value
	return b.ret
}

// Required set Metadata.Required
func (b *_FieldBuilder[R]) Required(value bool) R {
	b.metadata.Required = value
	return b.ret
}

// ReadOnly set Metadata.ReadOnly
func (b *_FieldBuilder[R]) ReadOnly(value bool) R {
	b.metadata.ReadOnly = value
	return b.ret
}

// WriteOnly set Metadata.WriteOnly
func (b *_FieldBuilder[R]) WriteOnly(value bool) R {
	b.metadata.WriteOnly = value
	return b.ret
}

// AllowEmptyValue set Metadata.AllowEmptyValue
func (b *_FieldBuilder[R]) AllowEmptyValue(value bool) R {
	b.metadata.AllowEmptyValue = value
	return b.ret
}

// Deprecated set Metadata.Deprecated
func (b *_FieldBuilder[R]) Deprecated(value bool) R {
	b.metadata.Deprecated = value
	return b.ret
}

func (b *_FieldBuilder[R]) Doc(stmts ...string) R {
	b.metadata.Doc = strings.Join(stmts, "\n")
	return b.ret
}

// end setter of Field --------------------

// Extension builds Type for Extension
func (b *Builder) Extension(name string, value any) *Extension {
	t := &Extension{
		_ExtensionBuilder: &_ExtensionBuilder[*Extension]{
			_Type: &_Type[*Extension]{rootbuilder: b, metadata: &TypeMetadata{Name: "", underlying: "extension", goType: "Extension"}},
			metadata: &ExtensionMetadata{
				Name: name, Value: value,
			},
		},
	}
	t.ret = t
	return t
}

type Extension struct {
	*_ExtensionBuilder[*Extension]
}

func (t *Extension) GetMetadata() *ExtensionMetadata {
	return t.metadata
}

type _ExtensionBuilder[R TypeBuilder] struct {
	*_Type[R]
	metadata *ExtensionMetadata
}

// begin setter of Extension --------------------

// end setter of Extension --------------------

// Object builds Type for Object
func (b *Builder) Object(fields ...*Field) *Object {
	t := &Object{
		_ObjectBuilder: &_ObjectBuilder[*Object]{
			_Type: &_Type[*Object]{rootbuilder: b, metadata: &TypeMetadata{Name: "", underlying: "object", goType: "Object"}},
			metadata: &ObjectMetadata{
				Fields: fields,
				Strict: true,
			},
		},
	}
	t.ret = t
	return t
}

type Object struct {
	*_ObjectBuilder[*Object]
}

func (t *Object) GetMetadata() *ObjectMetadata {
	return t.metadata
}

type _ObjectBuilder[R TypeBuilder] struct {
	*_Type[R]
	metadata *ObjectMetadata
}

// begin setter of Object --------------------

// MaxProperties set Metadata.MaxProperties
func (b *_ObjectBuilder[R]) MaxProperties(value uint64) R {
	b.metadata.MaxProperties = value
	return b.ret
}

// MinProperties set Metadata.MinProperties
func (b *_ObjectBuilder[R]) MinProperties(value uint64) R {
	b.metadata.MinProperties = value
	return b.ret
}

// Strict set Metadata.Strict
func (b *_ObjectBuilder[R]) Strict(value bool) R {
	b.metadata.Strict = value
	return b.ret
}

// end setter of Object --------------------

// _Container builds Type for _Container
func (b *Builder) _Container() *_Container {
	t := &_Container{
		__ContainerBuilder: &__ContainerBuilder[*_Container]{
			_Type:    &_Type[*_Container]{rootbuilder: b, metadata: &TypeMetadata{Name: "", underlying: "container", goType: "_Container"}},
			metadata: &_ContainerMetadata{},
		},
	}
	t.ret = t
	return t
}

type _Container struct {
	*__ContainerBuilder[*_Container]
}

func (t *_Container) GetMetadata() *_ContainerMetadata {
	return t.metadata
}

type __ContainerBuilder[R TypeBuilder] struct {
	*_Type[R]
	metadata *_ContainerMetadata
}

// begin setter of _Container --------------------

// Op set Metadata.Op
func (b *__ContainerBuilder[R]) Op(value string) R {
	b.metadata.Op = value
	return b.ret
}

// Types set Metadata.Types
func (b *__ContainerBuilder[R]) Types(value []Type) R {
	b.metadata.Types = value
	return b.ret
}

// Discriminator set Metadata.Discriminator
func (b *__ContainerBuilder[R]) Discriminator(value string) R {
	b.metadata.Discriminator = value
	return b.ret
}

// end setter of _Container --------------------

// Action builds Type for Action
func (b *Builder) Action(name string, inputoroutput ...InputOrOutput) *Action {
	t := &Action{
		_ActionBuilder: &_ActionBuilder[*Action]{
			_Type: &_Type[*Action]{rootbuilder: b, metadata: &TypeMetadata{Name: "", underlying: "action", goType: "Action"}},
			metadata: &ActionMetadata{
				Name: name,
			},
		},
	}

	t.metadata.Input, t.metadata.Outputs = func() (v1 *Input, v2 []*Output) {
		for _, x := range inputoroutput {
			switch x := x.(type) {
			case *Input:
				v1 = x
			case *Output:
				v2 = append(v2, x) // TODO: status conflict check
			default:
				panic(fmt.Sprintf("unexpected Type: %T", x))
			}
		}
		return
	}()
	t.ret = t
	return t
}

type Action struct {
	*_ActionBuilder[*Action]
}

func (t *Action) GetMetadata() *ActionMetadata {
	return t.metadata
}

type _ActionBuilder[R TypeBuilder] struct {
	*_Type[R]
	metadata *ActionMetadata
}

// begin setter of Action --------------------

// DefaultError set Metadata.DefaultError
func (b *_ActionBuilder[R]) DefaultError(value Type) R {
	b.metadata.DefaultError = value
	return b.ret
}

// Method set Metadata.Method
func (b *_ActionBuilder[R]) Method(value string) R {
	b.metadata.Method = value
	return b.ret
}

// Path set Metadata.Path
func (b *_ActionBuilder[R]) Path(value string) R {
	b.metadata.Path = value
	return b.ret
}

// Tags set Metadata.Tags
func (b *_ActionBuilder[R]) Tags(value []string) R {
	b.metadata.Tags = value
	return b.ret
}

// end setter of Action --------------------

// Param builds Type for Param
func (b *Builder) Param(name string, typ Type) *Param {
	t := &Param{
		_ParamBuilder: &_ParamBuilder[*Param]{
			_Type: &_Type[*Param]{rootbuilder: b, metadata: &TypeMetadata{Name: "", underlying: "param", goType: "Param"}},
			metadata: &ParamMetadata{
				Name: name, Typ: typ,
				In: "query", Required: true,
			},
		},
	}
	t.ret = t
	return t
}

type Param struct {
	*_ParamBuilder[*Param]
}

func (t *Param) GetMetadata() *ParamMetadata {
	return t.metadata
}

type _ParamBuilder[R TypeBuilder] struct {
	*_Type[R]
	metadata *ParamMetadata
}

// begin setter of Param --------------------

// In set Metadata.In
func (b *_ParamBuilder[R]) In(value string) R {
	b.metadata.In = value
	return b.ret
}

// Required set Metadata.Required
func (b *_ParamBuilder[R]) Required(value bool) R {
	b.metadata.Required = value
	return b.ret
}

// Deprecated set Metadata.Deprecated
func (b *_ParamBuilder[R]) Deprecated(value bool) R {
	b.metadata.Deprecated = value
	return b.ret
}

// AllowEmptyValue set Metadata.AllowEmptyValue
func (b *_ParamBuilder[R]) AllowEmptyValue(value bool) R {
	b.metadata.AllowEmptyValue = value
	return b.ret
}

func (b *_ParamBuilder[R]) Doc(stmts ...string) R {
	b.metadata.Doc = strings.Join(stmts, "\n")
	return b.ret
}

// end setter of Param --------------------

// Body builds Type for Body
func (b *Builder) Body(typ Type) *Body {
	t := &Body{
		_BodyBuilder: &_BodyBuilder[*Body]{
			_Type: &_Type[*Body]{rootbuilder: b, metadata: &TypeMetadata{Name: "", underlying: "Body", goType: "Body"}},
			metadata: &BodyMetadata{
				Typ: typ,
			},
		},
	}
	t.ret = t
	return t
}

type Body struct {
	*_BodyBuilder[*Body]
}

func (t *Body) GetMetadata() *BodyMetadata {
	return t.metadata
}

type _BodyBuilder[R TypeBuilder] struct {
	*_Type[R]
	metadata *BodyMetadata
}

// begin setter of Body --------------------

// end setter of Body --------------------

// Input builds Type for Input
func (b *Builder) Input(params ...paramOrBody) *Input {
	t := &Input{
		_InputBuilder: &_InputBuilder[*Input]{
			_Type:    &_Type[*Input]{rootbuilder: b, metadata: &TypeMetadata{Name: "", underlying: "input", goType: "Input"}},
			metadata: &InputMetadata{},
		},
	}

	t.metadata.Params, t.metadata.Body = func() (v1 []*Param, v2 *Body) {
		for _, a := range params {
			switch a := a.(type) {
			case *Param:
				v1 = append(v1, a)
			case *Body:
				v2 = a
			default:
				panic(fmt.Sprintf("unexpected Type: %T", a))
			}
		}
		return
	}()
	t.ret = t
	return t
}

type Input struct {
	*_InputBuilder[*Input]
}

func (t *Input) GetMetadata() *InputMetadata {
	return t.metadata
}

type _InputBuilder[R TypeBuilder] struct {
	*_Type[R]
	metadata *InputMetadata
}

// begin setter of Input --------------------

// end setter of Input --------------------

// Output builds Type for Output
func (b *Builder) Output(typ Type) *Output {
	t := &Output{
		_OutputBuilder: &_OutputBuilder[*Output]{
			_Type: &_Type[*Output]{rootbuilder: b, metadata: &TypeMetadata{Name: "", underlying: "output", goType: "Output"}},
			metadata: &OutputMetadata{
				Typ:    typ,
				Status: 200,
			},
		},
	}
	t.ret = t
	return t
}

type Output struct {
	*_OutputBuilder[*Output]
}

func (t *Output) GetMetadata() *OutputMetadata {
	return t.metadata
}

type _OutputBuilder[R TypeBuilder] struct {
	*_Type[R]
	metadata *OutputMetadata
}

// begin setter of Output --------------------

// Status set Metadata.Status
func (b *_OutputBuilder[R]) Status(value int) R {
	b.metadata.Status = value
	return b.ret
}

// IsDefault set Metadata.IsDefault
func (b *_OutputBuilder[R]) IsDefault(value bool) R {
	b.metadata.IsDefault = value
	return b.ret
}

// end setter of Output --------------------

// internal Type

type _Type[R TypeBuilder] struct {
	metadata *TypeMetadata
	ret      R

	rootbuilder *Builder
}

func (t *_Type[R]) GetTypeMetadata() *TypeMetadata {
	return t.metadata
}

// begin setter of Type --------------------

// Title set Metadata.Title
func (t _Type[R]) Title(value string) R {
	t.metadata.Title = value
	return t.ret
}

// Example set Metadata.Example
func (t _Type[R]) Example(value string) R {
	t.metadata.Example = value
	return t.ret
}

func (t _Type[R]) Doc(stmts ...string) R {
	t.metadata.Doc = strings.Join(stmts, "\n")
	return t.ret
}

func (t _Type[R]) Extensions(extensions ...*Extension) R {
	t.metadata.Extensions = func() *orderedmap.OrderedMap {
		if extensions == nil {
			return nil
		}
		data := orderedmap.New()
		for _, ext := range extensions {
			m := ext.metadata
			name := m.Name
			if !strings.HasPrefix(name, "x-") {
				name = "x-" + name
			}
			data.Set(name, m.Value)
		}
		return data
	}()
	return t.ret
}

// end setter of Type --------------------
func (t *_Type[R]) storeType(name string) {
	t.metadata.Name = name
	t.rootbuilder.storeType(t.ret)
}

// Type is the one of pseudo sum types (union).
// Type is union of bool | int | float | string | object | array[T] | map[T]
type Type interface {
	typ()

	TypeBuilder
}

func (t *Bool) typ() {}

func (t *Int) typ() {}

func (t *Float) typ() {}

func (t *String) typ() {}

func (t *Object) typ() {}

func (t *Array[Items]) typ() {}

func (t *Map[Items]) typ() {}

func (t *_Container) typ() {}

func (t *TypeRef) typ() {}

// paramOrBody is the one of pseudo sum types (union).

type paramOrBody interface {
	paramorbody()
}

func (t *Param) paramorbody() {}

func (t *Body) paramorbody() {}

// InputOrOutput is the one of pseudo sum types (union).

type InputOrOutput interface {
	inputoroutput()
}

func (t *Input) inputoroutput() {}

func (t *Output) inputoroutput() {}
