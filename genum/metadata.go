// Generated by github.com/podhmo/gos/seed [-write -builder -metadata -pkgname genum]
package genum

type EnumMetadata struct {
	id         int    // required by reference
	Name       string `json:"-"` // required by reference (and toString)
	underlying string `json:"-"` // required by toString

	Doc string `json:"Doc"`
}

type IntMetadata struct {
	Default int `json:"default"`

	Members []IntValue
}

type IntValue struct {
	Name string

	Value int

	Doc string
}

type StringMetadata struct {
	Default string `json:"default"`

	Members []StringValue
}

type StringValue struct {
	Name string

	Value string

	Doc string
}
