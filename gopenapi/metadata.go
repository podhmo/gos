// Generated by github.com/podhmo/gos/gopenapi/tools [-write -builder -metadata -pkgname gopenapi]

package gopenapi

type TypeMetadata struct {
	id         int    // required by reference
	Name       string `json:"-"` // required by reference (and toString)
	underlying string `json:"-"` // required by toString

	Doc    string `json:"description"`
	Format string `json:"format"`
}

type BoolMetadata struct {
}

type IntMetadata struct {
	Maximum int64 `json:"maximum,omitempty"`

	Minimum int64 `json:"minimum,omitempty"`
}

type StringMetadata struct {
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

	Typ TypeBuilder `json:"-"`

	Description string `json:"description,omitempty"`

	Required bool `json:"-"`
}

type ObjectMetadata struct {
	Fields []*FieldMetadata `json:"-"`

	Strict bool `json:"-"`
}
