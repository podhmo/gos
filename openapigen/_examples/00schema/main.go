package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/podhmo/gos/enumgen"
	"github.com/podhmo/gos/openapigen"
)

func main() {
	b := openapigen.NewTypeBuilder(openapigen.DefaultConfig())

	Name := openapigen.DefineType("Name", b.String().MinLength(1))

	Person := openapigen.DefineType("Person", b.Object(
		b.Field("id", b.String()),
		b.Field("name", b.String()).Doc("name of person"),
		b.Field("age", b.Int().Format("int32")),
		b.Field("nickname", b.Reference(Name)).Required(false),
		b.Field("father", b.ReferenceByName("Person")).Required(false),
		b.Field("friends", b.Array(b.ReferenceByName("Person"))).Required(false),
	)).Doc("person object")

	PersonSummary := openapigen.DefineType("PersonSummary", b.Object(
		Person.OnlyFields("name", "nickname")...,
	)).Doc("person objec summary")

	TestScore := openapigen.DefineType("TestScore", b.Object(
		b.Field("title", b.String()),
		b.Field("tests", b.Map(b.Int()).Pattern(`\-score$`).Doc("score (0~100)")),
	))

	// enum, in production, import from other package
	var orderingEnum *enumgen.String
	{
		b := enumgen.NewEnumBuilder(enumgen.DefaultConfig())
		orderingEnum = b.String(
			b.StringValue("desc").Doc("降順"),
			b.StringValue("asc").Doc("昇順"),
		).Default("desc").Doc("順序")
	}
	Ordering := openapigen.DefineType("Ordering", b.StringFromEnum(orderingEnum))

	doc, err := openapigen.ToSchema(b)
	if err != nil {
		panic(err)
	}

	fmt.Fprintln(os.Stderr, Name, Person, PersonSummary, TestScore, Ordering)
	// fmt.Fprintln(os.Stderr)

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(doc); err != nil {
		panic(err)
	}
}
