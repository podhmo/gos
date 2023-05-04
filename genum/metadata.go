package genum

type EnumMetadata struct {
	Name string `json:"name"`
}

type ValueMetadata[T any] struct {
	Name  string `json:"name"`
	Value T      `json:"value"`
	Doc   string `json:"doc,omitempty"`
}
