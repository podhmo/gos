// Generated by github.com/podhmo/gos/gopenapi/tools [-write -builder -metadata -pkgname gopenapi]

package gopenapi

type TypeMetadata struct {
	id         int    // required by reference
	Name       string `json:"-"` // required by reference (and toString)
	underlying string `json:"-"` // required by toString

	Format string `json:"format"`
	Doc    string `json:"description"`
}

type BoolMetadata struct {
}

type IntMetadata struct {
	Enum    []int64 `json:"enum,omitempty"`
	Default int64   `json:"enum,omitempty"`
	Maximum int64   `json:"maximum,omitempty"`
	Minimum int64   `json:"minimum,omitempty"`
}

type StringMetadata struct {
	Enum      []string `json:"string,omitempty"`
	Default   string   `json:"enum,omitempty"`
	Pattern   string   `json:"pattern,omitempty"`
	MaxLength int64    `json:"maxlength,omitempty"`
	MinLength int64    `json:"minlength,omitempty"`
}

type ArrayMetadata struct {
	MaxItems int64 `json:"maxitems,omitempty"`
	MinItems int64 `json:"minitems,omitempty"`
}

type MapMetadata struct {
	Pattern string `json:"pattern,omitempty"`
}

type FieldMetadata struct {
	Name        string      `json:"-"`
	Typ         TypeBuilder `json:"-"`
	Description string      `json:"description,omitempty"`
	Required    bool        `json:"-"`
}

type ObjectMetadata struct {
	Fields []*FieldMetadata `json:"-"`
	Strict bool             `json:"-"`
}

type ActionMetadata struct {
	Name          string
	Input         *Input
	Output        *Output
	DefaultStatus int
	Method        string
	Path          string
}

type InputMetadata struct {
	Params []*Param
}

type OutputMetadata struct {
	Typ TypeBuilder
}

type ParamMetadata struct {
	Name        string      `json:"-"`
	In          string      `json:"in"`
	Typ         TypeBuilder `json:"-"`
	Description string      `json:"description,omitempty"`
	Required    bool        `json:"required"`
}
