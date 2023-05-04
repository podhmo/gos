package genum

type TypeMetadata struct {
	id          int
	Name        string `json:"-"`
	Description string `json:"description,omitempty"`

	IsNewType bool `json:"-"` // TODO: remove

	underlying string `json:"-"`
}

// customization
type EnumMetadata[T any] struct {
}

type ValueMetadata[T any] struct {
	Name    string `json:"name"`
	Value   T      `json:"value"`
	Doc     string `json:"doc,omitempty"`
	Default bool   `json:"-"`
}
