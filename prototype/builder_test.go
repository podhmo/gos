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
