package builder

// metadata (options)
//
// - primitive types
// - composite types
// - object with fields (product type)

// subset of OAS component schemas definition (not strict)
// https://swagger.io/docs/specification/data-models/data-types/

type Type struct {
	id          int
	Name        string `json:"-"`
	Description string `json:"description,omitempty"`
	Format      string `json:"format,omitempty"`

	IsNewType bool `json:"-"`

	underlying string `json:"-"`
}

type Object struct {
	Strict bool `json:"-"`
}

type Field struct {
	Name        string `json:"-"`
	Description string `json:"description,omitempty"`
	Required    bool   `json:"-"`
}

// ----------------------------------------
// primitive types
// ----------------------------------------

type String struct {
	MinLength int64  `json:"minlength,omitempty"`
	MaxLength int64  `json:"maxlength,omitempty"`
	Pattern   string `json:"pattern,omitempty"`
}

type Integer struct {
	// minimum ≤ value ≤ maximum
	Maximum int64 `json:"maximum,omitempty"`
	Minimum int64 `json:"minimum,omitempty"`
}

// ----------------------------------------
// composite types
// ----------------------------------------

type Array struct {
	MaxItems int64 `json:"maxitems,omitempty"`
	MinItems int64 `json:"minitems,omitempty"`
}

type Map struct {
	PatternProperties map[string]TypeBuilder `json:"-,omitempty"`
}
