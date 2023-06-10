package main

import (
	"encoding/json"
	"os"

	"github.com/iancoleman/orderedmap"
	"github.com/podhmo/gos/openapigen"
)

func main() {
	b := openapigen.NewBuilder(openapigen.DefaultConfig())

	// https://swagger.io/docs/specification/data-models/oneof-anyof-allof-not/#oneof

	Dog := openapigen.Define("Dog", b.Object(
		b.Field("bark", b.Bool()),
		b.Field("breed", b.String().Enum([]string{"Dingo", "Husky", "Retriever", "Sheperd"})),
	))

	Cat := openapigen.Define("Cat", b.Object(
		b.Field("hunts", b.Bool()),
		b.Field("age", b.Int()),
	))

	UpdatePet := b.Action("UpdatePet",
		b.Input(b.Body(b.OneOf(Cat, Dog).Discriminator("pet_store"))),
		b.Output(nil).Doc("Updated"),
	)

	r := openapigen.NewRouter(b.Object())
	r.Patch("/pets", UpdatePet)

	doc := orderedmap.New()
	if err := r.ToSchemaWith(b, doc); err != nil {
		panic(err)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(doc); err != nil {
		panic(err)
	}
}
