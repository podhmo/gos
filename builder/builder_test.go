package builder_test

import (
	"fmt"
	"testing"

	"github.com/podhmo/gos/builder"
)

func TestIt(*testing.T) {
	b := builder.New()

	Person := b.Object(
		b.Field("name", b.String().MinLength(1).MaxLength(255)),
		b.Field("age", b.Integer().Minimum(0).Doc("hoho")).Required(true).Doc("haha"),
		b.Field("skills", b.Array(b.String().Doc("yaya").MinLength(1)).MinItems(1)).Required(false),
	).As("Person").Doc(
		"this is summary",
		"",
		"this is long description\nhehehhe",
	)
	fmt.Println(Person)
}

func TestToString(t *testing.T) {
	b := builder.New()

	tests := []struct {
		name string
		typ  builder.TypeBuilder
		want string
	}{
		{"primitive-0", b.String(), "string"},
		{"primitive-1", b.Integer(), "integer"},
		{"new-type-primitive", b.String().As("Name"), "Name"},
		{"array-string", b.Array(b.String()), "array[string]"},
		{"array-array-string", b.Array(b.Array(b.String())), "array[array[string]]"},
		{"object", b.Object(
			b.Field("name", b.String()),
			b.Field("age", b.String()).Required(false),
		), "object{name, age?}"},
		{"new-type-object", b.Object(
			b.Field("name", b.String()),
			b.Field("age", b.String()).Required(false),
		).As("Person"), "Person"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := builder.ToString(tt.typ)
			if got != tt.want {
				t.Errorf("ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}
