package openapigen_test

import (
	"testing"

	"github.com/podhmo/gos/openapigen"
)

func TestToString(t *testing.T) {
	b := openapigen.NewTypeBuilder()

	tests := []struct {
		name string
		typ  openapigen.TypeBuilder
		want string
	}{
		{"primitive-0", b.String(), "string"},
		{"primitive-1", b.Int(), "integer"},
		{"new-type-primitive", openapigen.DefineType("Name", b.String()), "Name"},
		{"array-string", b.Array(b.String()), "array[string]"},
		{"array-array-string", b.Array(b.Array(b.String())), "array[array[string]]"},
		{"object", b.Object(
			b.Field("name", b.String()),
			b.Field("age", b.String()).Required(false),
		), "object{name, age?}"},
		{"new-type-object", openapigen.DefineType("Person", b.Object(
			b.Field("name", b.String()),
			b.Field("age", b.String()).Required(false),
		)), "Person{name, age?}"},
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
