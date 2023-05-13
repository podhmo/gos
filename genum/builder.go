// Generated by github.com/podhmo/gos/seed [-write -builder -metadata -pkgname genum]
package genum

import (
	"fmt"
	"sync"
	"strings"
)

type EnumBuilder interface {
	GetEnumMetadata() *EnumMetadata
}

// DefineEnum names Enum value.
func DefineEnum[T interface {
	EnumBuilder
	storeEnum(name string)
}](name string, t T) T {
	t.storeEnum(name)
	return t
}

type Builder struct {
	mu          	sync.Mutex
	namedEnums  []EnumBuilder
	nameToIDMap map[string][]int
}

func NewEnumBuilder() *Builder {
	return &Builder{nameToIDMap: map[string][]int{}}
}

// EachEnum iterates named Enum.
func (b *Builder) EachEnums(fn func(EnumBuilder) error) error {
	for _, t := range b.namedEnums {
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
func (b *Builder) Int(members ...IntValue,) *IntEnum {
	t := &IntEnum{
		IntBuilder: &IntBuilder[*IntEnum]{
			_Enum:    &_Enum[*IntEnum]{rootbuilder: b, metadata: &EnumMetadata{Name: "", underlying: "Int"}},
			metadata: &IntMetadata{
				Members: members, 
			},
		},
	}
	t.ret = t
	return t
}

type IntEnum struct {
	*IntBuilder[*IntEnum]
}

func (t *IntEnum) GetMetadata() *IntMetadata {
	return t.metadata
}

type IntBuilder[R EnumBuilder] struct {
	*_Enum[R]
	metadata *IntMetadata
	ret R
}

// begin setter of Int ----------------------------------------

// Default set Metadata.Default
func (b *IntBuilder[R]) Default(value int) R {
	b.metadata.Default = value
	return b.ret
}

// end setter of Int ----------------------------------------







// String builds Enum for String
func (b *Builder) String(members ...StringValue,) *StringEnum {
	t := &StringEnum{
		StringBuilder: &StringBuilder[*StringEnum]{
			_Enum:    &_Enum[*StringEnum]{rootbuilder: b, metadata: &EnumMetadata{Name: "", underlying: "String"}},
			metadata: &StringMetadata{
				Members: members, 
			},
		},
	}
	t.ret = t
	return t
}

type StringEnum struct {
	*StringBuilder[*StringEnum]
}

func (t *StringEnum) GetMetadata() *StringMetadata {
	return t.metadata
}

type StringBuilder[R EnumBuilder] struct {
	*_Enum[R]
	metadata *StringMetadata
	ret R
}

// begin setter of String ----------------------------------------

// Default set Metadata.Default
func (b *StringBuilder[R]) Default(value string) R {
	b.metadata.Default = value
	return b.ret
}

// end setter of String ----------------------------------------





// internal Enum

type _Enum[R EnumBuilder] struct {
	metadata *EnumMetadata
	ret      R

	rootbuilder *Builder
}

func (t *_Enum[R]) GetEnumMetadata() *EnumMetadata {
	return t.metadata
}
func (t *_Enum[R]) Doc(stmts ...string) R {
	t.metadata.Doc = strings.Join(stmts, "\n")
	return t.ret
}
func (t *_Enum[R]) storeEnum(name string) {
	t.metadata.Name = name
	t.rootbuilder.storeEnum(t.ret)
}
