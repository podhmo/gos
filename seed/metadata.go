package seed

type BuilderMetadata struct {
	Target     *TypeMetadata
	Types      []*TypeMetadata
	UnionTypes []*UnionTypeMetadata

	NeedReference bool
	NeedStringer  bool

	Imports          []Import
	InterfaceMethods []string
	Constructor      *Constructor
	Fields           []*FieldMetadata // fields of builder

	SysArgs     []string // runtime os.Args[1:]
	PkgName     string   // package {{.PkgName}}}
	GeneratedBy string   // github.com/podhmo/gos/seed
	Footer      string
}

type TypeMetadata struct {
	Name       Symbol
	Underlying string
	TVars      TypeVarMetadataList

	Doc         string
	NeedBuilder bool

	Constructor *Constructor
	Fields      []*FieldMetadata // fields of Metadata
	Setters     []*SetterDefinition

	Used map[string]bool
}

type UnionTypeMetadata struct {
	Type          *TypeMetadata
	Args          []*TypeMetadata // symbol?
	DistinguishID string

	InterfaceMethods []string
}

type TypeVarMetadata struct { // e.g. [T any]
	Name string
	Type Symbol
}

type FieldMetadata struct {
	Name string
	Type Symbol
	Tag  string

	Doc     string
	Default string
}

type Constructor struct {
	Args []*ArgMetadata
}
type SetterDefinition struct {
	Name string
	Arg  *ArgMetadata
}

type ArgMetadata struct {
	Name     string
	Type     Symbol
	Variadic bool // as ...<type>

	BindFields []string
	Transform  func(string) string
}

type Import struct {
	Name string
	Path string
}
