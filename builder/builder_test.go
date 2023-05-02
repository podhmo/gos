package builder_test

import (
	"fmt"
	"testing"

	"github.com/podhmo/gos/builder"
)

func TestIt(*testing.T) {
	b := builder.New()

	Person := b.Object("Person",
		b.Field("name", b.String().MinLength(1).MaxLength(255)),
		b.Field("age", b.Integer().Minimum(0).Doc("hoho")).Required(true).Doc("haha"),
		b.Field("skills", b.Array(b.String().Doc("yaya").MinLength(1)).MinItems(1)).Required(false),
	).Doc(
		"this is summary",
		"",
		"this is long description\nhehehhe",
	)
	fmt.Println(Person)
}

func TestType(t *testing.T) {
	b := builder.New()

	b.String().MaxLength(10).MinLength(1).Pattern("^xxx$")
	b.Array(b.Array(b.String().MaxLength(10)).MaxItems(10).MinItems(1))
}
