package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/iancoleman/orderedmap"
	"github.com/podhmo/gos/enumgen"
	"github.com/podhmo/gos/openapigen"
)

func main() {

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

	fmt.Fprintln(os.Stderr, "type  \t", openapigen.ToString(Person))
	fmt.Fprintln(os.Stderr, "type  \t", openapigen.ToString(PersonSummary))
	fmt.Fprintln(os.Stderr, "action\t", openapigen.ToString(Hello))
	fmt.Fprintln(os.Stderr, "input \t", openapigen.ToString(Hello.GetMetadata().Input))
	fmt.Fprintln(os.Stderr, "output\t", openapigen.ToString(Hello.GetMetadata().Output))

	// routing
	{
		r := openapigen.NewRouter()
		mount(r)
		doc := orderedmap.New()
		r.ToSchemaWith(b, doc)
		if err := enc.Encode(doc); err != nil {
			panic(err)
		}
	}
}

func mount(r *openapigen.Router) {
	{
		r := r.Tagged("greeting")
		r.Post("/hello/{name}", Hello)
	}
	{
		r := r.Tagged("people")
		r.Get("/people", ListPerson)
		r.Post("/people", CreatePerson)
	}
}

var b = openapigen.NewTypeBuilder()

var (
	Name = openapigen.DefineType("Name", b.String().MinLength(1))

	Person = openapigen.DefineType("Person", b.Object(
		b.Field("id", b.String()),
		b.Field("name", b.String()).Doc("name of person"),
		b.Field("age", b.Int().Format("int32")),
		b.Field("nickname", b.Reference(Name)).Required(false),
		b.Field("father", b.ReferenceByName("Person")).Required(false),
		b.Field("friends", b.Array(b.ReferenceByName("Person"))).Required(false),
	)).Doc("person object")

	PersonSummary = openapigen.DefineType("PersonSummary", b.Object(
		Person.OnlyFields("name", "nickname")...,
	)).Doc("person objec summary")

	TestScore = openapigen.DefineType("TestScore", b.Object(
		b.Field("title", b.String()),
		b.Field("tests", b.Map(b.Int()).Pattern(`\-score$`).Doc("score (0~100)")),
	))

	Ordering *openapigen.String
)

func init() {
	// enum, in production, import from other package
	var orderingEnum *enumgen.String
	{
		b := enumgen.NewEnumBuilder(enumgen.DefaultConfig())
		orderingEnum = b.String(
			b.StringValue("desc").Doc("降順"),
			b.StringValue("asc").Doc("昇順"),
		).Default("desc").Doc("順序")
	}

	Ordering = openapigen.DefineType("Ordering", b.StringFromEnum(orderingEnum))
}

// actions
var (
	// Hello :: func(name string) string
	Hello = b.Action("hello",
		b.Input(
			b.Path("name", b.String()),
		).Doc("input"),
		b.Output(
			b.String(),
		),
	).Doc("greeting hello")

	ListPerson = b.Action("ListPerson",
		b.Input(
			b.Query("sort", b.String().Enum([]string{"name", "-name", "age", "-age"})),
		),
		b.Output(b.Array(PersonSummary)),
	).Doc("list person")

	CreatePerson = b.Action("CreatePerson",
		b.Input(
			b.Query("verbose", b.Bool()),
			b.Body(b.Object(
				append(Person.IgnoreFields("id", "father", "friends"),
					b.Field("fatherId", b.String()),
					b.Field("friendIdList", b.Array(b.String())))...,
			)),
		),
		b.Output(Person),
	).Doc("create person")
)
