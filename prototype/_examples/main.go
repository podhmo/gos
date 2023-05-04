package main

import (
	"encoding/json"
	"os"

	"github.com/podhmo/gos/prototype"
)

func main() {
	b := prototype.NewBuilder()

	Name := prototype.Define("Name", b.String().MinLength(1))

	prototype.Define("Person", b.Object(
		b.Field("name", b.String()).Doc("name of person"),
		b.Field("age", b.Integer().Format("int32")),
		b.Field("nickname", b.Reference(Name)).Required(false),
		b.Field("father", b.ReferenceByName("Person")).Required(false),
		b.Field("friends", b.Array(b.ReferenceByName("Person"))).Required(false),
	)).Doc("person object")

	prototype.Define("TestScore", b.Object(
		b.Field("title", b.String()),
		b.Field("tests", b.Map(b.Integer()).Pattern(`\-score$`).Doc("score (0~100)")),
	))

	doc, err := prototype.ToSchema(b)
	if err != nil {
		panic(err)
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(doc); err != nil {
		panic(err)
	}
}
