// Generated by github.com/podhmo/gos/enumgen/tools [-write -builder -metadata -stringer -pkgname enumgen]

package enumgen

import (
	"fmt"
	"io"
	"strings"
	"sync"
)

type EnumBuilder interface {
	GetEnumMetadata() *EnumMetadata
	writeEnum(io.Writer) error // see: ./stringer.go
}

// Define names Enum value.
func Define[T interface {
	EnumBuilder
	storeEnum(name string)
}](name string, t T) T {
	t.storeEnum(name)
	return t
}

type Builder struct {
	mu          sync.Mutex
	namedEnums  []EnumBuilder
	nameToIDMap map[string][]int

	Config *Config
}

func NewBuilder(config *Config) *Builder {
	return &Builder{
		nameToIDMap: map[string][]int{},
		namedEnums:  []EnumBuilder{nil}, // nil is sentinel (id<=0 is unnamed)
		Config:      config,
	}
}

// EachEnum iterates named Enum.
func (b *Builder) EachEnums(fn func(EnumBuilder) error) error {
	for _, t := range b.namedEnums {
		if t == nil {
			continue
		}
		if err := fn(t); err != nil {
			return fmt.Errorf("error on %v -- %w", t, err) // TODO: use ToString()
		}
	}
	return nil
}

func (b *Builder) storeEnum(typ EnumBuilder) {
	val := typ.GetEnumMetadata()
	val.id = -1
	if val.Name == "" {
		return
	}

	b.mu.Lock()
	defer b.mu.Unlock()
	id := len(b.namedEnums)
	val.id = id
	b.namedEnums = append(b.namedEnums, typ)
	b.nameToIDMap[val.Name] = append(b.nameToIDMap[val.Name], id)
	// TODO: name conflict check
}

func (b *Builder) lookupEnum(name string) EnumBuilder {
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
	return b.namedEnums[ids[0]]
}

// Int builds Enum for Int
func (b *Builder) Int(members ...*IntValue) *Int {
	t := &Int{
		_IntBuilder: &_IntBuilder[*Int]{
			_Enum: &_Enum[*Int]{rootbuilder: b, metadata: &EnumMetadata{Name: "", underlying: "int", goType: "Int"}},
			metadata: &IntMetadata{
				Members: members,
			},
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

type _IntBuilder[R EnumBuilder] struct {
	*_Enum[R]
	metadata *IntMetadata
}

// begin setter of Int --------------------

// Default set Metadata.Default
func (b *_IntBuilder[R]) Default(value int) R {
	b.metadata.Default = value
	return b.ret
}

// end setter of Int --------------------

// IntValue builds Enum for IntValue
func (b *Builder) IntValue(value int, name string) *IntValue {
	t := &IntValue{
		_IntValueBuilder: &_IntValueBuilder[*IntValue]{
			_Enum: &_Enum[*IntValue]{rootbuilder: b, metadata: &EnumMetadata{Name: "", underlying: "IntValue", goType: "IntValue"}},
			metadata: &IntValueMetadata{
				Value: value, Name: name,
			},
		},
	}
	t.ret = t
	return t
}

type IntValue struct {
	*_IntValueBuilder[*IntValue]
}

func (t *IntValue) GetMetadata() *IntValueMetadata {
	return t.metadata
}

type _IntValueBuilder[R EnumBuilder] struct {
	*_Enum[R]
	metadata *IntValueMetadata
}

// begin setter of IntValue --------------------

// Doc set Metadata.Doc
func (b *_IntValueBuilder[R]) Doc(value string) R {
	b.metadata.Doc = value
	return b.ret
}

// end setter of IntValue --------------------

// String builds Enum for String
func (b *Builder) String(members ...*StringValue) *String {
	t := &String{
		_StringBuilder: &_StringBuilder[*String]{
			_Enum: &_Enum[*String]{rootbuilder: b, metadata: &EnumMetadata{Name: "", underlying: "string", goType: "String"}},
			metadata: &StringMetadata{
				Members: members,
			},
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

type _StringBuilder[R EnumBuilder] struct {
	*_Enum[R]
	metadata *StringMetadata
}

// begin setter of String --------------------

// Default set Metadata.Default
func (b *_StringBuilder[R]) Default(value string) R {
	b.metadata.Default = value
	return b.ret
}

// end setter of String --------------------

// StringValue builds Enum for StringValue
func (b *Builder) StringValue(value string) *StringValue {
	t := &StringValue{
		_StringValueBuilder: &_StringValueBuilder[*StringValue]{
			_Enum: &_Enum[*StringValue]{rootbuilder: b, metadata: &EnumMetadata{Name: "", underlying: "StringValue", goType: "StringValue"}},
			metadata: &StringValueMetadata{
				Value: value,
			},
		},
	}
	t.ret = t
	return t
}

type StringValue struct {
	*_StringValueBuilder[*StringValue]
}

func (t *StringValue) GetMetadata() *StringValueMetadata {
	return t.metadata
}

type _StringValueBuilder[R EnumBuilder] struct {
	*_Enum[R]
	metadata *StringValueMetadata
}

// begin setter of StringValue --------------------

// Name set Metadata.Name
func (b *_StringValueBuilder[R]) Name(value string) R {
	b.metadata.Name = value
	return b.ret
}

// Doc set Metadata.Doc
func (b *_StringValueBuilder[R]) Doc(value string) R {
	b.metadata.Doc = value
	return b.ret
}

// end setter of StringValue --------------------

// internal Enum

type _Enum[R EnumBuilder] struct {
	metadata *EnumMetadata
	ret      R

	rootbuilder *Builder
}

func (t *_Enum[R]) GetEnumMetadata() *EnumMetadata {
	return t.metadata
}

// begin setter of Enum --------------------

func (t _Enum[R]) Doc(stmts ...string) R {
	t.metadata.Doc = strings.Join(stmts, "\n")
	return t.ret
}

// end setter of Enum --------------------
func (t *_Enum[R]) storeEnum(name string) {
	t.metadata.Name = name
	t.rootbuilder.storeEnum(t.ret)
}
