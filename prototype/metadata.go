package prototype

// metadata (options)
//
// - primitive types
// - composite types
// - object with fields (product type)

// subset of OAS component schemas definition (not strict)
// https://swagger.io/docs/specification/data-models/data-types/

type TypeMetadata struct {
	id         int    // required by reference
	Name       string `json:"-"` // required by reference (and toString)
	underlying string `json:"-"` // required by toString

	Description string `json:"description,omitempty"`
	Format      string `json:"format,omitempty"`
}

type ObjectMetadata struct {
	Strict bool `json:"-"`
}

type FieldMetadata struct {
	Name        string `json:"-"`
	Description string `json:"description,omitempty"`
	Required    bool   `json:"-"`
}

// ----------------------------------------
// primitive types
// ----------------------------------------

type StringMetadata struct {
	MinLength int64  `json:"minlength,omitempty"`
	MaxLength int64  `json:"maxlength,omitempty"`
	Pattern   string `json:"pattern,omitempty"`
}

type IntegerMetadata struct {
	// minimum ≤ value ≤ maximum
	Maximum int64 `json:"maximum,omitempty"`
	Minimum int64 `json:"minimum,omitempty"`
}

// ----------------------------------------
// composite types
// ----------------------------------------

type ArrayMetadata struct {
	MaxItems int64 `json:"maxitems,omitempty"`
	MinItems int64 `json:"minitems,omitempty"`
}

type MapMetadata struct {
	Pattern string `json:"-,omitempty"`
}

// ----------------------------------------
// actions
// ----------------------------------------

type ActionMetadata struct {
	DefaultStatus int
}
type ActionInputMetadata struct {
}
type ActionOutputMetadata struct {
}
type ActionParamMetadata struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Required    bool   `json:"required"`
	In          string `json:"in"` // enum
}
