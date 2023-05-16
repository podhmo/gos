package main

import (
	"fmt"
	"os"

	"github.com/podhmo/gos/gopenapi"
)

func main() {
	b := gopenapi.NewTypeBuilder()

	Name := gopenapi.DefineType("Name", b.String().MinLength(1))

	Person := gopenapi.DefineType("Person", b.Object(
		b.Field("name", b.String()).Doc("name of person"),
		b.Field("age", b.Int().Format("int32")),
		b.Field("nickname", b.Reference(Name)).Required(false),
		b.Field("father", b.ReferenceByName("Person")).Required(false),
		b.Field("friends", b.Array(b.ReferenceByName("Person"))).Required(false),
	)).Doc("person object")

	TestScore := gopenapi.DefineType("TestScore", b.Object(
		b.Field("title", b.String()),
		b.Field("tests", b.Map(b.Int()).Pattern(`\-score$`).Doc("score (0~100)")),
	))

	// doc, err := gopenapi.ToSchema(b)
	// if err != nil {
	// 	panic(err)
	// }

	fmt.Fprintln(os.Stderr, Name, Person, TestScore)

	// enc := json.NewEncoder(os.Stdout)
	// enc.SetIndent("", "  ")
	// if err := enc.Encode(doc); err != nil {
	// 	panic(err)
	// }
}
