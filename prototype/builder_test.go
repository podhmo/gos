package prototype_test

import (
	"testing"

	"github.com/podhmo/gos/prototype"
)

func TestIt(t *testing.T) {
	b := prototype.NewBuilder()

	typeName := "Person"
	prototype.Define(typeName, b.Object(
		b.Field("name", b.String().MinLength(1).MaxLength(255)),
		b.Field("age", b.Integer().Minimum(0).Doc("hoho")).Required(true).Doc("haha"),
		b.Field("skills", b.Array(b.String().Doc("yaya").MinLength(1)).MinItems(1)).Required(false), // composite
		b.Field("father", b.ReferenceByName(typeName)).Required(false),                              // recursive
	)).Doc(
		"this is summary",
		"",
		"this is long description\nhehehhe",
	)

	b.EachTypes(func(typ prototype.TypeBuilder) error {
		got := prototype.ToString(typ)
		want := `Person{name, age, skills?, father?}`
		if want != got {
			t.Errorf("want %q, but got %q", want, got)
		}
		return nil
	})
}

func TestToString(t *testing.T) {
	b := prototype.NewBuilder()

	tests := []struct {
		name string
		typ  prototype.TypeBuilder
		want string
	}{
		{"primitive-0", b.String(), "string"},
		{"primitive-1", b.Integer(), "integer"},
		{"new-type-primitive", prototype.Define("Name", b.String()), "Name"},
		{"array-string", b.Array(b.String()), "array[string]"},
		{"array-array-string", b.Array(b.Array(b.String())), "array[array[string]]"},
		{"object", b.Object(
			b.Field("name", b.String()),
			b.Field("age", b.String()).Required(false),
		), "object{name, age?}"},
		{"new-type-object", prototype.Define("Person", b.Object(
			b.Field("name", b.String()),
			b.Field("age", b.String()).Required(false),
		)), "Person{name, age?}"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := prototype.ToString(tt.typ)
			if got != tt.want {
				t.Errorf("ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}
