package main

import (
	"encoding/json"
	"os"

	"github.com/podhmo/gos/builder"
)

func main() {
	b := builder.New()

	Name := builder.Define("Name", b.String().MinLength(1))

	builder.Define("Person", b.Object(
		b.Field("name", b.String()).Doc("name of person"),
		b.Field("age", b.Integer().Format("int32")),
		b.Field("nickname", b.Reference(Name)).Required(false),
		b.Field("father", b.ReferenceByName("Person")).Required(false),
		b.Field("friends", b.Array(b.ReferenceByName("Person"))).Required(false),
	)).Doc("person object")

	builder.Define("TestScore", b.Object(
		b.Field("title", b.String()),
		b.Field("tests", b.Map(b.Integer()).
			PatternProperties(`\-score$`, b.Integer().Doc("score (0~100)")).
			PatternProperties(`\-grade$`, b.String().Doc("grade (A,B,C,D,E,F)"))),
	))

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
