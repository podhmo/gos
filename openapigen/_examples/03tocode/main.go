package main

import (
	"fmt"
	"os"

	"github.com/podhmo/gos/openapigen"
)

func main() {
	b := openapigen.NewTypeBuilder(openapigen.DefaultConfig())

	Name := openapigen.DefineType("Name", b.String()).
		Doc("name of something")

	openapigen.DefineType("Person", b.Object(
		b.Field("name", b.Reference(Name)).Doc("name of person"),
		b.Field("age", b.Int()),
	)).Doc("person object")

	w := os.Stdout
	fmt.Fprintln(w, "package M")
	if err := openapigen.ToGocode(w, b); err != nil {
		panic(err)
	}
}
