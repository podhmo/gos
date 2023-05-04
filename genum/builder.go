package genum

import (
	"fmt"
	"strings"
	"sync"
)

type TypeBuilder[T any] interface {
	GetTypeMetadata() *TypeMetadata

	writeCoder // to code
}

func Define[T interface {
	TypeBuilder[T]
	storeType(name string)
}](name string, typ T) T {
	typ.storeType(name)
	return typ
}

type Config struct {
	Padding string
	Comment string
}

type Builder[T any] struct {
	Config *Config

	mu          sync.Mutex
	namedTypes  []TypeBuilder[T]
	nameToIDMap map[string][]int
}

func NewBuilder[T any]() *Builder[T] {
	return &Builder[T]{
		Config:      &Config{Padding: "\t"},
		nameToIDMap: map[string][]int{},
	}
}

func (b *Builder[T]) EachTypes(fn func(TypeBuilder[T]) error) error {
	for _, t := range b.namedTypes {
		if err := fn(t); err != nil {
			return err
			// TODO: ToString()
			// return fmt.Errorf("error on %s -- %w", ToString(t), err)
		}
	}
	return nil
}

func (b *Builder[T]) storeType(typ TypeBuilder[T]) {
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

func (b *Builder[T]) lookupType(name string) TypeBuilder[T] {
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

func (b *Builder[T]) Enum(values ...*ValueType[T]) *EnumType[T] {
	var z T
	t := &EnumType[T]{
		EnumBuilder: &EnumBuilder[T, *EnumType[T]]{
			type_:    &type_[T, *EnumType[T]]{rootbuilder: b, metadata: &TypeMetadata{Name: "", underlying: fmt.Sprintf("%T", z)}}, // TODO: underlying
			metadata: &EnumMetadata[T]{},
			Values:   values,
		},
	}
	t.EnumBuilder.ret = t
	return t
}

type EnumType[T any] struct {
	*EnumBuilder[T, *EnumType[T]]
}

func (t *EnumType[T]) GetMetadata() *EnumMetadata[T] {
	return t.metadata
}

type EnumBuilder[T any, R TypeBuilder[T]] struct {
	*type_[T, R]
	metadata *EnumMetadata[T]
	Values   []*ValueType[T]
}

func (b *Builder[T]) Value(v T) *ValueType[T] {
	t := &ValueType[T]{
		ValueBuilder: &ValueBuilder[T, *ValueType[T]]{
			type_:    &type_[T, *ValueType[T]]{rootbuilder: b, metadata: &TypeMetadata{Name: "", underlying: ""}}, // TODO: remove
			metadata: &ValueMetadata[T]{Value: v},
		},
	}
	t.ValueBuilder.ret = t
	return t
}

type ValueType[T any] struct {
	*ValueBuilder[T, *ValueType[T]]
}

func (t *ValueType[T]) GetMetadata() *ValueMetadata[T] {
	return t.metadata
}

type ValueBuilder[T any, R TypeBuilder[T]] struct {
	*type_[T, R] // need?
	metadata     *ValueMetadata[T]
}

func (b *ValueBuilder[T, R]) Name(name string) R {
	b.metadata.Name = name
	return b.ret
}
func (b *ValueBuilder[T, R]) Doc(stmts ...string) R {
	b.metadata.Doc = strings.Join(stmts, "\n")
	return b.ret
}

func (b *ValueBuilder[T, R]) Default(v bool) R {
	b.metadata.Default = v
	return b.ret
}

type type_[T any, R TypeBuilder[T]] struct {
	metadata *TypeMetadata
	ret      R

	rootbuilder *Builder[T]
}

func (t *type_[T, R]) GetTypeMetadata() *TypeMetadata {
	return t.metadata
}
func (t *type_[T, R]) Doc(stmts ...string) R {
	t.metadata.Description = strings.Join(stmts, "\n")
	return t.ret
}
func (t *type_[T, R]) storeType(name string) {
	t.metadata.Name = name
	t.metadata.IsNewType = true
	t.rootbuilder.storeType(t.ret)
}
