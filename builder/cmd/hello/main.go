package main

import (
	"fmt"

	"github.com/podhmo/gos/builder"
)

// TODO: required,
// TODO: completion
// TODO: distinct, type and field?
// TODO: construct type (e.g. b.Union(b.String(), b.Integer()))
// TODO: recursive definition
// TODO: action(input,output)
// TODO: namelib (goify)
// TODO: enum

func main() {
	b := builder.New()

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
