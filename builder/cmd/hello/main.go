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
		b.String("name"),
		b.Integer("age"),
	)
	fmt.Println(Person)
}
