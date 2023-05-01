package builder

import (
	"fmt"
	"testing"
)

func TestIt(*testing.T) {
	b := New()

	Person := b.Type("Person",
		b.String("name").MinLength(1).MaxLength(255),
		b.Integer("age").Minimum(0).Required(true).Doc("hoho"),
	).Doc(
		"this is summary",
		"",
		"this is long description\nhehehhe",
	)
	fmt.Println(Person)
}
