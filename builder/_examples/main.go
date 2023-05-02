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
		b.Field("name", b.String()),
		b.Field("age", b.Integer()),
		b.Field("nickname", b.Reference(Name)).Required(false),
		b.Field("father", b.ReferenceByName("Person")).Required(false),
	).As("Person")

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
