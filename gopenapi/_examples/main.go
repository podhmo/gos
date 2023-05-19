package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/podhmo/gos/genum"
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

	PersonSummary := gopenapi.DefineType("PersonSummary", b.Object(
		Person.OnlyFields("name", "nickname")...
	)).Doc("person objec summary")

	TestScore := gopenapi.DefineType("TestScore", b.Object(
		b.Field("title", b.String()),
		b.Field("tests", b.Map(b.Int()).Pattern(`\-score$`).Doc("score (0~100)")),
	))

	// enum, in production, import from other package
	var orderingEnum *genum.String
	{
		b := genum.NewEnumBuilder(genum.DefaultConfig())
		orderingEnum = b.String(
			b.StringValue("desc").Doc("降順"),
			b.StringValue("asc").Doc("昇順"),
		).Default("desc").Doc("順序")
	}
	Ordering := gopenapi.DefineType("Ordering", b.StringFromEnum(orderingEnum))

	// TODO:
	// Hello :: func(name string) string
	Hello := b.Action("hello",
		b.Input(
			b.Param("name", b.String(), "path"),
		).Doc("input"),
		b.Output(
			b.String(),
		),
	).Method("GET").Path("/hello")

	doc, err := gopenapi.ToSchema(b)
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

	fmt.Fprintln(os.Stderr, "type  \t", gopenapi.ToString(Person))
	fmt.Fprintln(os.Stderr, "type  \t", gopenapi.ToString(PersonSummary))
	fmt.Fprintln(os.Stderr, "action\t", gopenapi.ToString(Hello))
	fmt.Fprintln(os.Stderr, "input \t", gopenapi.ToString(Hello.GetMetadata().Input))
	fmt.Fprintln(os.Stderr, "output\t", gopenapi.ToString(Hello.GetMetadata().Output))
}
