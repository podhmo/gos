// Generated by github.com/podhmo/gos/gopenapi/tools [-write -builder -metadata -stringer -pkgname gopenapi]

package gopenapi

import (
	"fmt"
	"io"
	"strings"
	"sync"
)

type TypeBuilder interface {
	GetTypeMetadata() *TypeMetadata
	writeType(io.Writer) error // see: ./to_string.go
	toSchemer                  // see: ./to_schema.go
}

// DefineType names Type value.
func DefineType[T interface {
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
}

func NewTypeBuilder() *Builder {
	return &Builder{
		nameToIDMap: map[string][]int{},
	}
}

// EachType iterates named Type.
func (b *Builder) EachTypes(fn func(TypeBuilder) error) error {
	for _, t := range b.namedTypes {
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
		BoolBuilder: &BoolBuilder[*Bool]{
			_Type:    &_Type[*Bool]{rootbuilder: b, metadata: &TypeMetadata{Name: "", underlying: "boolean"}},
			metadata: &BoolMetadata{},
		},
	}
	t.ret = t
	return t
}

type Bool struct {
	*BoolBuilder[*Bool]
}

func (t *Bool) GetMetadata() *BoolMetadata {
	return t.metadata
}

type BoolBuilder[R TypeBuilder] struct {
	*_Type[R]
	metadata *BoolMetadata
}

// Int builds Type for Int
func (b *Builder) Int() *Int {
	t := &Int{
		IntBuilder: &IntBuilder[*Int]{
			_Type:    &_Type[*Int]{rootbuilder: b, metadata: &TypeMetadata{Name: "", underlying: "integer"}},
			metadata: &IntMetadata{},
		},
	}
	t.ret = t
	return t
}

type Int struct {
	*IntBuilder[*Int]
}

func (t *Int) GetMetadata() *IntMetadata {
	return t.metadata
}

type IntBuilder[R TypeBuilder] struct {
	*_Type[R]
	metadata *IntMetadata
}

// begin setter of Int --------------------

// Enum set Metadata.Enum
func (b *IntBuilder[R]) Enum(value []int64) R {
	b.metadata.Enum = value
	return b.ret
}

// Default set Metadata.Default
func (b *IntBuilder[R]) Default(value int64) R {
	b.metadata.Default = value
	return b.ret
}

// Maximum set Metadata.Maximum
func (b *IntBuilder[R]) Maximum(value int64) R {
	b.metadata.Maximum = value
	return b.ret
}

// Minimum set Metadata.Minimum
func (b *IntBuilder[R]) Minimum(value int64) R {
	b.metadata.Minimum = value
	return b.ret
}

// end setter of Int --------------------

// String builds Type for String
func (b *Builder) String() *String {
	t := &String{
		StringBuilder: &StringBuilder[*String]{
			_Type:    &_Type[*String]{rootbuilder: b, metadata: &TypeMetadata{Name: "", underlying: "string"}},
			metadata: &StringMetadata{},
		},
	}
	t.ret = t
	return t
}

type String struct {
	*StringBuilder[*String]
}

func (t *String) GetMetadata() *StringMetadata {
	return t.metadata
}

type StringBuilder[R TypeBuilder] struct {
	*_Type[R]
	metadata *StringMetadata
}

// begin setter of String --------------------

// Enum set Metadata.Enum
func (b *StringBuilder[R]) Enum(value []string) R {
	b.metadata.Enum = value
	return b.ret
}

// Default set Metadata.Default
func (b *StringBuilder[R]) Default(value string) R {
	b.metadata.Default = value
	return b.ret
}

// Pattern set Metadata.Pattern
func (b *StringBuilder[R]) Pattern(value string) R {
	b.metadata.Pattern = value
	return b.ret
}

// MaxLength set Metadata.MaxLength
func (b *StringBuilder[R]) MaxLength(value int64) R {
	b.metadata.MaxLength = value
	return b.ret
}

// MinLength set Metadata.MinLength
func (b *StringBuilder[R]) MinLength(value int64) R {
	b.metadata.MinLength = value
	return b.ret
}

// end setter of String --------------------

// Array builds Type for Array
func (b *Builder) Array(items TypeBuilder) *Array[TypeBuilder] {
	t := &Array[TypeBuilder]{
		ArrayBuilder: &ArrayBuilder[TypeBuilder, *Array[TypeBuilder]]{
			_Type:    &_Type[*Array[TypeBuilder]]{rootbuilder: b, metadata: &TypeMetadata{Name: "", underlying: "array"}},
			metadata: &ArrayMetadata{},
			items:    items,
		},
	}
	t.ret = t
	return t
}

type Array[Items TypeBuilder] struct {
	*ArrayBuilder[Items, *Array[Items]]
}

func (t *Array[Items]) GetMetadata() *ArrayMetadata {
	return t.metadata
}

type ArrayBuilder[Items TypeBuilder, R TypeBuilder] struct {
	*_Type[R]
	metadata *ArrayMetadata
	items    Items
}

// begin setter of Array --------------------

// MaxItems set Metadata.MaxItems
func (b *ArrayBuilder[Items, R]) MaxItems(value int64) R {
	b.metadata.MaxItems = value
	return b.ret
}

// MinItems set Metadata.MinItems
func (b *ArrayBuilder[Items, R]) MinItems(value int64) R {
	b.metadata.MinItems = value
	return b.ret
}

// end setter of Array --------------------

// Map builds Type for Map
func (b *Builder) Map(items TypeBuilder) *Map[TypeBuilder] {
	t := &Map[TypeBuilder]{
		MapBuilder: &MapBuilder[TypeBuilder, *Map[TypeBuilder]]{
			_Type:    &_Type[*Map[TypeBuilder]]{rootbuilder: b, metadata: &TypeMetadata{Name: "", underlying: "map"}},
			metadata: &MapMetadata{},
			items:    items,
		},
	}
	t.ret = t
	return t
}

type Map[Items TypeBuilder] struct {
	*MapBuilder[Items, *Map[Items]]
}

func (t *Map[Items]) GetMetadata() *MapMetadata {
	return t.metadata
}

type MapBuilder[Items TypeBuilder, R TypeBuilder] struct {
	*_Type[R]
	metadata *MapMetadata
	items    Items
}

// begin setter of Map --------------------

// Pattern set Metadata.Pattern
func (b *MapBuilder[Items, R]) Pattern(value string) R {
	b.metadata.Pattern = value
	return b.ret
}

// end setter of Map --------------------

// Field builds Type for Field
func (b *Builder) Field(name string, typ TypeBuilder) *Field {
	t := &Field{
		FieldBuilder: &FieldBuilder[*Field]{
			_Type: &_Type[*Field]{rootbuilder: b, metadata: &TypeMetadata{Name: "", underlying: "field"}},
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
	*FieldBuilder[*Field]
}

func (t *Field) GetMetadata() *FieldMetadata {
	return t.metadata
}

type FieldBuilder[R TypeBuilder] struct {
	*_Type[R]
	metadata *FieldMetadata
}

// begin setter of Field --------------------

// Description set Metadata.Description
func (b *FieldBuilder[R]) Description(value string) R {
	b.metadata.Description = value
	return b.ret
}

// Required set Metadata.Required
func (b *FieldBuilder[R]) Required(value bool) R {
	b.metadata.Required = value
	return b.ret
}

// end setter of Field --------------------

// Object builds Type for Object
func (b *Builder) Object(fields ...*Field) *Object {
	t := &Object{
		ObjectBuilder: &ObjectBuilder[*Object]{
			_Type: &_Type[*Object]{rootbuilder: b, metadata: &TypeMetadata{Name: "", underlying: "object"}},
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
	*ObjectBuilder[*Object]
}

func (t *Object) GetMetadata() *ObjectMetadata {
	return t.metadata
}

type ObjectBuilder[R TypeBuilder] struct {
	*_Type[R]
	metadata *ObjectMetadata
}

// begin setter of Object --------------------

// Strict set Metadata.Strict
func (b *ObjectBuilder[R]) Strict(value bool) R {
	b.metadata.Strict = value
	return b.ret
}

// end setter of Object --------------------

// Action builds Type for Action
func (b *Builder) Action(name string, input *Input, output *Output) *Action {
	t := &Action{
		ActionBuilder: &ActionBuilder[*Action]{
			_Type: &_Type[*Action]{rootbuilder: b, metadata: &TypeMetadata{Name: "", underlying: "action"}},
			metadata: &ActionMetadata{
				Name: name, Input: input, Output: output,
				DefaultStatus: 200,
			},
		},
	}
	t.ret = t
	return t
}

type Action struct {
	*ActionBuilder[*Action]
}

func (t *Action) GetMetadata() *ActionMetadata {
	return t.metadata
}

type ActionBuilder[R TypeBuilder] struct {
	*_Type[R]
	metadata *ActionMetadata
}

// begin setter of Action --------------------

// DefaultStatus set Metadata.DefaultStatus
func (b *ActionBuilder[R]) DefaultStatus(value int) R {
	b.metadata.DefaultStatus = value
	return b.ret
}

// Method set Metadata.Method
func (b *ActionBuilder[R]) Method(value string) R {
	b.metadata.Method = value
	return b.ret
}

// Path set Metadata.Path
func (b *ActionBuilder[R]) Path(value string) R {
	b.metadata.Path = value
	return b.ret
}

// Tags set Metadata.Tags
func (b *ActionBuilder[R]) Tags(value []string) R {
	b.metadata.Tags = value
	return b.ret
}

// end setter of Action --------------------

// Param builds Type for Param
func (b *Builder) Param(name string, typ TypeBuilder, in string) *Param {
	t := &Param{
		ParamBuilder: &ParamBuilder[*Param]{
			_Type: &_Type[*Param]{rootbuilder: b, metadata: &TypeMetadata{Name: "", underlying: "param"}},
			metadata: &ParamMetadata{
				Name: name, Typ: typ, In: in,
				Required: true,
			},
		},
	}
	t.ret = t
	return t
}

type Param struct {
	*ParamBuilder[*Param]
}

func (t *Param) GetMetadata() *ParamMetadata {
	return t.metadata
}

type ParamBuilder[R TypeBuilder] struct {
	*_Type[R]
	metadata *ParamMetadata
}

// begin setter of Param --------------------

// Description set Metadata.Description
func (b *ParamBuilder[R]) Description(value string) R {
	b.metadata.Description = value
	return b.ret
}

// Required set Metadata.Required
func (b *ParamBuilder[R]) Required(value bool) R {
	b.metadata.Required = value
	return b.ret
}

// end setter of Param --------------------

// Body builds Type for Body
func (b *Builder) Body(typ TypeBuilder) *Body {
	t := &Body{
		BodyBuilder: &BodyBuilder[*Body]{
			_Type: &_Type[*Body]{rootbuilder: b, metadata: &TypeMetadata{Name: "", underlying: "Body"}},
			metadata: &BodyMetadata{
				Typ: typ,
			},
		},
	}
	t.ret = t
	return t
}

type Body struct {
	*BodyBuilder[*Body]
}

func (t *Body) GetMetadata() *BodyMetadata {
	return t.metadata
}

type BodyBuilder[R TypeBuilder] struct {
	*_Type[R]
	metadata *BodyMetadata
}

// begin setter of Body --------------------

// end setter of Body --------------------

// Input builds Type for Input
func (b *Builder) Input(params ...paramOrBody) *Input {
	t := &Input{
		InputBuilder: &InputBuilder[*Input]{
			_Type:    &_Type[*Input]{rootbuilder: b, metadata: &TypeMetadata{Name: "", underlying: "input"}},
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
			}
		}
		return
	}()
	t.ret = t
	return t
}

type Input struct {
	*InputBuilder[*Input]
}

func (t *Input) GetMetadata() *InputMetadata {
	return t.metadata
}

type InputBuilder[R TypeBuilder] struct {
	*_Type[R]
	metadata *InputMetadata
}

// begin setter of Input --------------------

// end setter of Input --------------------

// Output builds Type for Output
func (b *Builder) Output(typ TypeBuilder) *Output {
	t := &Output{
		OutputBuilder: &OutputBuilder[*Output]{
			_Type: &_Type[*Output]{rootbuilder: b, metadata: &TypeMetadata{Name: "", underlying: "output"}},
			metadata: &OutputMetadata{
				Typ: typ,
			},
		},
	}
	t.ret = t
	return t
}

type Output struct {
	*OutputBuilder[*Output]
}

func (t *Output) GetMetadata() *OutputMetadata {
	return t.metadata
}

type OutputBuilder[R TypeBuilder] struct {
	*_Type[R]
	metadata *OutputMetadata
}

// begin setter of Output --------------------

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

// Format set Metadata.Format
func (t _Type[R]) Format(value string) R {
	t.metadata.Format = value
	return t.ret
}

func (t _Type[R]) Doc(stmts ...string) R {
	t.metadata.Doc = strings.Join(stmts, "\n")
	return t.ret
}

// end setter of Type --------------------
func (t *_Type[R]) storeType(name string) {
	t.metadata.Name = name
	t.rootbuilder.storeType(t.ret)
}

// paramOrBody is the one of pseudo sum types (union).
type paramOrBody interface {
	paramorbody()
}

func (t *Param) paramorbody() {}

func (t *Body) paramorbody() {}

// footer. ----
// toSlice is list.map as you know.
func toSlice[S, D any](src []S, conv func(S) D) []D {
	dst := make([]D, len(src))
	for i, x := range src {
		dst[i] = conv(x)
	}
	return dst
}
