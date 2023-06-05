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

	Friends := openapigen.DefineType("Friends", b.Array(b.ReferenceByName("Person"))).
		Doc("Friends of something")

	openapigen.DefineType("Person", b.Object(
		b.Field("name", b.Reference(Name)).Doc("name of person"),
		b.Field("age", b.Int()),
		b.Field("father", b.ReferenceByName("Person")).Nullable(true),
		b.Field("children", b.Array(b.ReferenceByName("Person"))),
		b.Field("friends", b.Reference(Friends)),
	)).Doc("person object")

	w := os.Stdout
	fmt.Fprintln(w, "package M")
	if err := openapigen.ToGocode(w, b); err != nil {
		panic(err)
	}
}
