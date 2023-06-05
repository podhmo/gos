package design

import (
	"github.com/podhmo/gos/enumgen"
	"github.com/podhmo/gos/openapigen"
)

var b = openapigen.NewTypeBuilder(openapigen.DefaultConfig())
var Builder = b // export

var (
	Name = openapigen.Define("Name", b.String().MinLength(1))

	Person = openapigen.Define("Person", b.Object(
		b.Field("id", b.String()),
		b.Field("name", b.String()).Doc("name of person"),
		b.Field("age", b.Int().Format("int32")),
		b.Field("nickname", b.Reference(Name)).Required(false),
		b.Field("father", b.ReferenceByName("Person")).Required(false),
		b.Field("friends", b.Array(b.ReferenceByName("Person"))).Required(false),
	)).Doc("person object")

	PersonSummary = openapigen.Define("PersonSummary", b.Object(
		Person.OnlyFields("name", "nickname")...,
	)).Doc("person objec summary")

	TestScore = openapigen.Define("TestScore", b.Object(
		b.Field("title", b.String()),
		b.Field("tests", b.Map(b.Int()).Pattern(`\-score$`).Doc("score (0~100)")),
	))

	Ordering *openapigen.String
)

func init() {
	// enum, in production, import from other package
	var orderingEnum *enumgen.String
	{
		b := enumgen.NewBuilder(enumgen.DefaultConfig())
		orderingEnum = b.String(
			b.StringValue("desc").Doc("降順"),
			b.StringValue("asc").Doc("昇順"),
		).Default("desc").Doc("順序")
	}

	Ordering = openapigen.Define("Ordering", b.StringFromEnum(orderingEnum))
}
