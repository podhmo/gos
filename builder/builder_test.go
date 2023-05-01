package builder

import (
	"fmt"
	"testing"
)

func TestIt(*testing.T) {
	b := New()

	Person := b.Type("Person",
		b.Field("name").String().MinLength(1).MaxLength(255),
		b.Field("age").Integer().Minimum(0).Required(true).Doc("hoho"),
	).Doc(
		"this is summary",
		"",
		"this is long description\nhehehhe",
	)
	fmt.Println(Person)
}
