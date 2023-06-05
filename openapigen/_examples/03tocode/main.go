package main

import (
	"os"

	"github.com/podhmo/gos/openapigen"
)

func main() {
	b := openapigen.NewTypeBuilder(openapigen.DefaultConfig())

	openapigen.DefineType("Person", b.Object(
		b.Field("name", b.String()).Doc("name of person"),
		b.Field("age", b.Int()),
	)).Doc("person object")

	w := os.Stdout
	if err := openapigen.ToGocode(w, b); err != nil {
		panic(err)
	}
}
