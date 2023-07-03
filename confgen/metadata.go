// Generated by github.com/podhmo/gos/confgen/tools [-write -builder -metadata -stringer -pkgname confgen]

package confgen

type TypeMetadata struct {
	id         int    // required by reference
	Name       string `json:"-"` // required by reference (and toString)
	underlying string `json:"-"` // required by toString
	goType     string `json:"-"`

	Doc     string `json:"description,omitempty"`
	Title   string `json:"title,omitempty"`
	Example string `json:"example,omitempty"`
}

type BoolMetadata struct {
	Format string `json:"format,omitempty"`

	Default bool `json:"default,omitempty"`
}

type IntMetadata struct {
	Format string `json:"format,omitempty"`

	Enum []int64 `json:"enum,omitempty"`

	Default int64 `json:"default,omitempty"`

	Maximum int64 `json:"maximum,omitempty"`

	Minimum int64 `json:"minimum,omitempty"`

	ExclusiveMin bool `json:"exclusiveMin,omitempty"`

	ExclusiveMax bool `json:"exclusiveMax,omitempty"`
}

type FloatMetadata struct {
	Format string `json:"format,omitempty"`

	Default string `json:"default,omitempty"`

	Maximum float64 `json:"maximum,omitempty"`

	Minimum float64 `json:"minimum,omitempty"`

	MultipleOf float64 `json:"multipleOf,omitempty"`

	ExclusiveMin bool `json:"exclusiveMin,omitempty"`

	ExclusiveMax bool `json:"exclusiveMax,omitempty"`
}

type StringMetadata struct {
	Format string `json:"format,omitempty"`

	Enum []string `json:"enum,omitempty"`

	Default string `json:"default,omitempty"`

	Pattern string `json:"pattern,omitempty"`

	MaxLength int64 `json:"maxLength,omitempty"`

	MinLength int64 `json:"minLength,omitempty"`
}

type ArrayMetadata struct {
	MaxItems int64 `json:"maxItems,omitempty"`

	MinItems int64 `json:"minItems,omitempty"`
}

type MapMetadata struct {
	Pattern string `json:"pattern,omitempty"`
}

type FieldMetadata struct {
	Name string `json:"-"`

	Typ Type `json:"-"`

	Doc string `json:"description,omitempty"`

	Nullable bool `json:"nullable,omitempty"`

	Required bool `json:"-"`

	ReadOnly bool `json:"readonly,omitempty"`

	WriteOnly bool `json:"writeonly,omitempty"`

	AllowEmptyValue bool `json:"allowEmptyValue,omitempty"`

	Deprecated bool `json:"deprecated,omitempty"`
}

// for x-<extension-name>
type ExtensionMetadata struct {
	Name string

	Value any
}

type ObjectMetadata struct {
	Fields []*Field `json:"-"`

	MaxProperties uint64 `json:"maxProeprties,omitempty"`

	MinProperties uint64 `json:"minProeprties,omitempty"`

	Strict bool `json:"-"`
}

// the container for allOf, anyOf, oneOf
type _ContainerMetadata struct {
	Op string `json:"-"`

	Types []Type `json:"-"`

	Discriminator string `json:"-"`
}
