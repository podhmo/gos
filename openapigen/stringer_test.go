package openapigen_test

import (
	"testing"

	"github.com/podhmo/gos/openapigen"
)

func TestToString(t *testing.T) {
	b := openapigen.NewBuilder(openapigen.DefaultConfig())

	// Hello :: func(name string) string
	Hello := b.Action("hello",
		b.Input(
			b.Param("name", b.String()).AsPath(),
		).Doc("input"),
		b.Output(
			b.String(),
		),
	).Doc("greeting hello")

	tests := []struct {
		name string
		typ  openapigen.TypeBuilder
		want string
	}{
		{"primitive-0", b.String(), "string"},
		{"primitive-1", b.Int(), "integer"},
		{"new-type-primitive", openapigen.Define("Name", b.String()), "Name"},
		{"array-string", b.Array(b.String()), "array[string]"},
		{"array-array-string", b.Array(b.Array(b.String())), "array[array[string]]"},
		{"object", b.Object(
			b.Field("name", b.String()),
			b.Field("age", b.String()).Required(false),
		), "object{name, age?}"},
		{"new-type-object", openapigen.Define("Person", b.Object(
			b.Field("name", b.String()),
			b.Field("age", b.String()).Required(false),
		)), "Person{name, age?}"},
		{"action", Hello, "hello :: (string) => string"},
		{"action input", Hello.GetMetadata().Input, "(string)"},
		{"action output", Hello.GetMetadata().Outputs[0], " => string"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := openapigen.ToString(tt.typ)
			if got != tt.want {
				t.Errorf("ToString() = %v, but want is %v", got, tt.want)
			}
		})
	}
}
