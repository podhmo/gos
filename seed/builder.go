package seed

import (
	"fmt"
	"os"
	"strings"
)

type Builder struct {
	Metadata *BuilderMetadata
}

func NewBuilder(pkgname string, fields ...*Field) *Builder {
	metadata := make([]*FieldMetadata, len(fields))
	for i, f := range fields {
		metadata[i] = f.Metadata
	}

	return &Builder{
		Metadata: &BuilderMetadata{
			PkgName:     pkgname,
			SysArgs:     os.Args[1:],
			GeneratedBy: "github.com/podhmo/gos/seed",
			Fields:      metadata,
		},
	}
}

var Root = NewBuilder("")

type Symbol string

func (s Symbol) Pointer() Symbol {
	return "*" + s
}
func (s Symbol) Slice() Symbol {
	return "[]" + s
}

func (b *Builder) BuildTarget(name string, fields ...*Field) Symbol {
	b.Metadata.Target = Symbol(name)

	metadata := make([]*FieldMetadata, len(fields))
	for i, f := range fields {
		metadata[i] = f.Metadata
	}
	b.Metadata.TargetFields = metadata
	return Symbol(name)
}
func (b *Builder) NeedReference() *Builder {
	b.Metadata.NeedReference = true
	return b
}

func (b *Builder) InterfaceMethods(methods ...string) *Builder {
	b.Metadata.InterfaceMethods = append(b.Metadata.InterfaceMethods, methods...)
	return b
}

func (b *Builder) Import(path string) Symbol {
	b.Metadata.Imports = append(b.Metadata.Imports, Import{Path: path})
	parts := strings.Split(path, "/")
	return Symbol(parts[len(parts)-1])
}
func (b *Builder) NamedImport(name string, path string) Symbol {
	b.Metadata.Imports = append(b.Metadata.Imports, Import{Name: name, Path: path})
	return Symbol(name)
}

func (b *Builder) TypeVar(name string, typ Symbol) *TypeVar {
	t := &TypeVar{
		typeVarBuilder: &typeVarBuilder[*TypeVar]{
			Metadata: &TypeVarMetadata{
				Name: name,
				Type: typ,
			},
		},
	}
	t.ret = t
	return t
}
func (b *Builder) Field(name string, typ Symbol) *Field {
	t := &Field{
		fieldBuilder: &fieldBuilder[*Field]{
			Metadata: &FieldMetadata{
				Name: name,
				Type: typ,
			},
		},
	}
	t.ret = t
	return t
}
func (b *Builder) Arg(name string, typ Symbol) *Arg {
	t := &Arg{
		argBuilder: &argBuilder[*Arg]{
			Metadata: &ArgMetadata{
				Name: name,
				Type: typ,
			},
		},
	}
	t.ret = t
	return t
}
func (b *Builder) Constructor(args ...*Arg) *Builder {
	metadata := make([]*ArgMetadata, len(args))
	for i, a := range args {
		metadata[i] = a.Metadata
	}
	b.Metadata.Constructor = &Constructor{Args: metadata}
	return b
}

func (b *Builder) Type(name string, typeVarOrFieldList ...typeAttr) *Type {
	tvars := make([]*TypeVarMetadata, 0, len(typeVarOrFieldList))
	fields := make([]*FieldMetadata, 0, len(typeVarOrFieldList))
	for _, tattr := range typeVarOrFieldList {
		switch t := tattr.(type) {
		case *TypeVar:
			tvars = append(tvars, t.Metadata)
		case *Field:
			fields = append(fields, t.Metadata)
		default:
			panic(fmt.Sprintf("unexpected type: %T", tattr))
		}
	}

	t := &Type{
		TypeBuilder: &TypeBuilder[*Type]{Metadata: &TypeMetadata{
			Name:       Symbol(name),
			Underlying: name,
			TVars:      tvars,
			Fields:     fields,
			Used:       map[string]bool{},
		}},
	}
	t.ret = t
	b.Metadata.Types = append(b.Metadata.Types, t)
	return t
}

type Type struct {
	*TypeBuilder[*Type]
}

type TypeBuilder[R any] struct {
	Metadata *TypeMetadata
	ret      R
}

func (b *TypeBuilder[R]) NeedBuilder() R {
	b.Metadata.NeedBuilder = true
	return b.ret
}
func (b *TypeBuilder[R]) Underlying(v string) R {
	b.Metadata.Underlying = v
	return b.ret
}
func (b *TypeBuilder[R]) Constructor(args ...*Arg) R {
	metadata := make([]*ArgMetadata, len(args))
	for i, a := range args {
		metadata[i] = a.Metadata
	}
	b.Metadata.Constructor = &Constructor{Args: metadata}
	for _, a := range metadata {
		b.Metadata.Used[a.Name] = true
	}
	return b.ret
}

// metadata
type BuilderMetadata struct {
	Target       Symbol
	TargetFields []*FieldMetadata // fields of Metadata

	Types []*Type

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
	TVars      []*TypeVarMetadata

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
	Metadata *FieldMetadata
	ret      R
}

func (b *fieldBuilder[R]) Tag(v string) R {
	b.Metadata.Tag = v
	return b.ret
}

type TypeVar struct {
	*typeVarBuilder[*TypeVar]
}
type typeVarBuilder[R any] struct {
	Metadata *TypeVarMetadata
	ret      R
}
type Arg struct {
	*argBuilder[*Arg]
}
type argBuilder[R any] struct {
	Metadata *ArgMetadata
	ret      R
}

func (b *argBuilder[R]) Variadic() R {
	b.Metadata.Variadic = true
	return b.ret
}

type typeAttr interface {
	typeattr()
}

func (t *TypeVar) typeattr() {}
func (t *Field) typeattr()   {}
