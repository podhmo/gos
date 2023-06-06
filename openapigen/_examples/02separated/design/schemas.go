package design

import (
	"github.com/podhmo/gos/openapigen"
	"github.com/podhmo/gos/openapigen/_examples/02separated/enum"
)

var Builder = openapigen.NewBuilder(openapigen.DefaultConfig()) // export
var b = Builder

var (
	Error = openapigen.Define("Error", b.Object(
		b.Field("message", b.String()),
	))
)

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

	// enumgen.String ->  openapigen.String
	Ordering = openapigen.Define("Ordering", b.StringFromEnum(enum.Ordering))
)
