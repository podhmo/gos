package main

import (
	"encoding/json"
	"os"

	"github.com/podhmo/gos/builder"
)

func main() {
	b := builder.New()

	// TODO: minlength
	Name := b.String().MinLength(1).As("Name")

	b.Object(
		b.Field("name", b.String()).Doc("name of person"),
		b.Field("age", b.Integer().Format("int32")),
		b.Field("nickname", b.Reference(Name)).Required(false),
		b.Field("father", b.ReferenceByName("Person")).Required(false),
	).As("Person").Doc("person object")

	doc, err := builder.ToSchema(b)
	if err != nil {
		panic(err)
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(doc); err != nil {
		panic(err)
	}
}
