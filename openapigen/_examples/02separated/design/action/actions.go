package action

import (
	design "github.com/podhmo/gos/openapigen/_examples/02separated/design"
)

var b = design.Builder

var (
	// Hello :: func(name string) string
	Hello = b.Action("hello",
		b.Input(
			b.Param("name", b.String()).AsPath(),
		).Doc("input"),
		b.Output(
			b.String(),
		),
	).Doc("greeting hello")

	// ListPerson :: func(...) []PersonSummary
	ListPerson = b.Action("ListPerson",
		b.Input(
			b.Param("sort", b.String().Enum([]string{"name", "-name", "age", "-age"})).AsQuery(),
		),
		b.Output(b.Array(design.PersonSummary)).Doc("list of person summary"),
	).Doc("list person")

	// CreatePerson :: func(...)
	CreatePerson = b.Action("CreatePerson",
		b.Input(
			b.Param("verbose", b.Bool()).AsQuery(),
			b.Body(b.Object(
				append(design.Person.IgnoreFields("id", "father", "friends"),
					b.Field("fatherId", b.String()),
					b.Field("friendIdList", b.Array(b.String())))...,
			)).Doc("person but father and friends are id"),
		),
		b.Output(nil).Status(204),
	).Doc("create person")
)
