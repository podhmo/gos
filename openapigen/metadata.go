// Generated by github.com/podhmo/gos/openapigen/tools [-write -builder -metadata -stringer -pkgname openapigen]

package openapigen

type TypeMetadata struct {
	id         int    // required by reference
	Name       string `json:"-"` // required by reference (and toString)
	underlying string `json:"-"` // required by toString

	Doc     string `json:"description,omitempty"`
	Title   string `json:"title,omitempty"`
	Format  string `json:"format,omitempty"`
	Example string `json:"example,omitempty"`
}

type BoolMetadata struct {
	Default bool `json:"default,omitempty"`
}

type IntMetadata struct {
	Enum []int64 `json:"enum,omitempty"`

	Default int64 `json:"default,omitempty"`

	Maximum int64 `json:"maximum,omitempty"`

	Minimum int64 `json:"minimum,omitempty"`

	ExclusiveMin bool `json:"exclusiveMin,omitempty"`

	ExclusiveMax bool `json:"exclusiveMax,omitempty"`
}

type FloatMetadata struct {
	Default string `json:"default,omitempty"`

	Maximum float64 `json:"maximum,omitempty"`

	Minimum float64 `json:"minimum,omitempty"`

	MultipleOf float64 `json:"multipleOf,omitempty"`

	ExclusiveMin bool `json:"exclusiveMin,omitempty"`

	ExclusiveMax bool `json:"exclusiveMax,omitempty"`
}

type StringMetadata struct {
	Enum []string `json:"enum,omitempty"`

	Default string `json:"default,omitempty"`

	Pattern string `json:"pattern,omitempty"`

	MaxLength int64 `json:"maxlength,omitempty"`

	MinLength int64 `json:"minlength,omitempty"`
}

type ArrayMetadata struct {
	MaxItems int64 `json:"maxitems,omitempty"`

	MinItems int64 `json:"minitems,omitempty"`
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

type ObjectMetadata struct {
	Fields []*Field `json:"-"`

	MaxProperties uint64 `json:"maxProeprties,omitempty"`

	MinProperties uint64 `json:"minProeprties,omitempty"`

	Strict bool `json:"-"`
}

type ActionMetadata struct {
	Name string `json:"-"`

	Input *Input `json:"-"`

	Outputs []*Output `json:"-"`

	DefaultError Type

	Method string `json:"-"`

	Path string `json:"-"`

	Tags []string `json:"tags"`
}

type ParamMetadata struct {
	Name string `json:"name"`
	// openapi's in parameter {query, header, path, cookie} (default is query)
	In string `json:"in"`

	Typ Type `json:"-"`

	Doc string `json:"description,omitempty"`

	Required bool `json:"required"`

	Deprecated bool `json:"deprecated,omitempty"`

	AllowEmptyValue bool `json:"allowEmptyValue,omitempty"`
}

type BodyMetadata struct {
	Typ Type `json:"-"`
}

type InputMetadata struct {
	Params []*Param

	Body *Body
}

type OutputMetadata struct {
	Typ Type

	Status int `json:"-"`

	IsDefault bool `json:"-"`
}

type Contact struct {
	Name string `json:"name"`

	URL string `json:"url"`

	Email string `json:"email"`
}

type License struct {

	// required
	Name string `json:"name"`

	Identifier string `json:"identifier"`

	URL string `json:"url"`
}

type Server struct {

	// required
	URL string `json:"url"`

	Doc string `json:"description"`
	// todo: typed
	Variables map[string]any `json:"variables,omitempty"`
}

type Info struct {

	// required
	Title string `json:"title"`

	Summary string `json:"summary,omitempty"`

	Doc string `json:"description,omitempty"`

	TermsOfService string `json:"termOfService,omitempty"`

	Contact *Contact `json:"contact,omitempty"`

	License *License `json:"license,omitempty"`
	// required
	Version string `json:"version"`
}

type OpenAPI struct {

	// required
	OpenAPI string `json:"openapi"`
	// required
	Info Info `json:"info"`

	Servers []Server `json:"servers"`
}
