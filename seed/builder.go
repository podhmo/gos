package seed

import (
	"fmt"
	"os"
	"strings"
)

type Builder struct {
	metadata *BuilderMetadata
}

func NewBuilder(pkgname string, fields ...*Field) *Builder {
	metadata := make([]*FieldMetadata, len(fields))
	for i, f := range fields {
		metadata[i] = f.metadata
	}

	return &Builder{
		metadata: &BuilderMetadata{
			PkgName:     pkgname,
			SysArgs:     os.Args[1:],
			GeneratedBy: "github.com/podhmo/gos/seed",
			Fields:      metadata,
		},
	}
}

func (b *Builder) GetMetadata() *BuilderMetadata {
	return b.metadata
}

var Root = NewBuilder("")

type Symbol string

func (s Symbol) Pointer() Symbol {
	return "*" + s
}
func (s Symbol) Slice() Symbol {
	return "[]" + s
}

// BuildTarget is setter method for setting the name of your root builder's type
func (b *Builder) BuildTarget(name string, fields ...*Field) Symbol {
	b.metadata.Target = Symbol(name)

	metadata := make([]*FieldMetadata, len(fields))
	for i, f := range fields {
		metadata[i] = f.metadata
	}
	b.metadata.TargetFields = metadata
	return Symbol(name)
}

// NeedReference is setter method if you need reference functions in generated code.
func (b *Builder) NeedReference() *Builder {
	b.metadata.NeedReference = true
	return b
}

// InterfaceMethods is setter method if your builder requires more interfaces
// (e.g. "String() string", "fmt.Stringer", ...)
func (b *Builder) InterfaceMethods(methods ...string) *Builder {
	b.metadata.InterfaceMethods = append(b.metadata.InterfaceMethods, methods...)
	return b
}

// Import is setter method adding import path in generated code.
func (b *Builder) Import(path string) Symbol {
	b.metadata.Imports = append(b.metadata.Imports, Import{Path: path})
	parts := strings.Split(path, "/")
	return Symbol(parts[len(parts)-1])
}

// Import is setter method adding import path with name in generated code.
func (b *Builder) NamedImport(name string, path string) Symbol {
	b.metadata.Imports = append(b.metadata.Imports, Import{Name: name, Path: path})
	return Symbol(name)
}

// TypeVar is factory method for TypeVar builder.
func (b *Builder) TypeVar(name string, typ Symbol) *TypeVar {
	t := &TypeVar{
		typeVarBuilder: &typeVarBuilder[*TypeVar]{
			metadata: &TypeVarMetadata{
				Name: name,
				Type: typ,
			},
		},
	}
	t.ret = t
	return t
}

// TypeVar is factory method for Field builder.
func (b *Builder) Field(name string, typ Symbol) *Field {
	t := &Field{
		fieldBuilder: &fieldBuilder[*Field]{
			metadata: &FieldMetadata{
				Name: name,
				Type: typ,
			},
		},
	}
	t.ret = t
	return t
}

// Arg is factory method for Arg builder.
func (b *Builder) Arg(name string, typ Symbol) *Arg {
	t := &Arg{
		argBuilder: &argBuilder[*Arg]{
			metadata: &ArgMetadata{
				Name: name,
				Type: typ,
			},
		},
	}
	t.ret = t
	return t
}

// Constructor is setter method customize your root builder's factory.
func (b *Builder) Constructor(args ...*Arg) *Builder {
	metadata := make([]*ArgMetadata, len(args))
	for i, a := range args {
		metadata[i] = a.metadata
	}
	b.metadata.Constructor = &Constructor{Args: metadata}
	return b
}

// Type is factory method for Type builder.
func (b *Builder) Type(name string, typeVarOrFieldList ...typeAttr) *Type {
	tvars := make([]*TypeVarMetadata, 0, len(typeVarOrFieldList))
	fields := make([]*FieldMetadata, 0, len(typeVarOrFieldList))
	for _, tattr := range typeVarOrFieldList {
		switch t := tattr.(type) {
		case *TypeVar:
			tvars = append(tvars, t.metadata)
		case *Field:
			fields = append(fields, t.metadata)
		default:
			panic(fmt.Sprintf("unexpected type: %T", tattr))
		}
	}

	t := &Type{
		typeBuilder: &typeBuilder[*Type]{metadata: &TypeMetadata{
			Name:       Symbol(name),
			Underlying: name,
			TVars:      tvars,
			Fields:     fields,
			Used:       map[string]bool{},
		}},
	}
	t.ret = t
	b.metadata.Types = append(b.metadata.Types, t.metadata)
	return t
}

// GeneratedBy is setter method for auto generated comment.
// (default value is "github.com/podhmo/gos/seed" )
func (b *Builder) GeneratedBy(v string) *Builder {
	b.metadata.GeneratedBy = v
	return b
}

type Type struct {
	*typeBuilder[*Type]
}

type typeBuilder[R any] struct {
	metadata *TypeMetadata
	ret      R
}

func (t *typeBuilder[R]) GetMetadata() *TypeMetadata {
	return t.metadata
}

// NeedBuilder is setter method for the generated go type needs builder implementation.
func (b *typeBuilder[R]) NeedBuilder() R {
	b.metadata.NeedBuilder = true
	return b.ret
}

// Underlying is setter method if you can set underlying type (default is same as TypeName)
func (b *typeBuilder[R]) Underlying(v string) R {
	b.metadata.Underlying = v
	return b.ret
}

// Constructor is setter method for cusotmization of builder factory
func (b *typeBuilder[R]) Constructor(args ...*Arg) R {
	metadata := make([]*ArgMetadata, len(args))
	for i, a := range args {
		metadata[i] = a.metadata
	}
	b.metadata.Constructor = &Constructor{Args: metadata}
	for _, a := range metadata {
		b.metadata.Used[a.Name] = true
	}
	return b.ret
}

// metadata
type BuilderMetadata struct {
	Target       Symbol
	TargetFields []*FieldMetadata // fields of Metadata

	Types []*TypeMetadata

	NeedReference bool

	Imports          []Import
	InterfaceMethods []string
	Constructor      *Constructor
	Fields           []*FieldMetadata // fields of builder

	SysArgs     []string // runtime os.Args[1:]
	PkgName     string   // package {{.PkgName}}}
	GeneratedBy string   // github.com/podhmo/gos/seed
}

type TypeMetadata struct {
	Name       Symbol
	Underlying string
	TVars      TypeVarMetadataList

	NeedBuilder bool
	Constructor *Constructor
	Fields      []*FieldMetadata // fields of Metadata

	Used map[string]bool
}

type TypeVarMetadata struct { // e.g. [T any]
	Name string
	Type Symbol
}

type FieldMetadata struct {
	Name string
	Type Symbol
	Tag  string
}

type Constructor struct {
	Args []*ArgMetadata
}

type ArgMetadata struct {
	Name     string
	Type     Symbol
	Variadic bool // as ...<type>

	Transform func(string) string
}

type Import struct {
	Name string
	Path string
}

// ----------------------------------------

type Field struct {
	*fieldBuilder[*Field]
}
type fieldBuilder[R any] struct {
	metadata *FieldMetadata
	ret      R
}

func (b *fieldBuilder[R]) GetMetadata() *FieldMetadata {
	return b.metadata
}

// Tag is setter method for set metadata.Tag
func (b *fieldBuilder[R]) Tag(v string) R {
	b.metadata.Tag = v
	return b.ret
}

type TypeVar struct {
	*typeVarBuilder[*TypeVar]
}
type typeVarBuilder[R any] struct {
	metadata *TypeVarMetadata
	ret      R
}

func (b *typeVarBuilder[R]) GetMetadata() *TypeVarMetadata {
	return b.metadata
}

type Arg struct {
	*argBuilder[*Arg]
}
type argBuilder[R any] struct {
	metadata *ArgMetadata
	ret      R
}

func (b *argBuilder[R]) GetMetadata() *ArgMetadata {
	return b.metadata
}

// Variadic is setter method for set metadata.Variadic is true
func (b *argBuilder[R]) Variadic() R {
	b.metadata.Variadic = true
	return b.ret
}

// Transform is setter method for setting transform function.
// (This method is typically used when the parent data has children of type XXXMetadata.)
func (b *argBuilder[R]) Transform(fn func(string) string) R {
	b.metadata.Transform = fn
	return b.ret
}

type typeAttr interface {
	typeattr()
}

func (t *TypeVar) typeattr() {}
func (t *Field) typeattr()   {}
