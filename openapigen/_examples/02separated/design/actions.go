package design

import "github.com/podhmo/gos/openapigen"

func NewActions(b *openapigen.Builder) (Actions struct {
	Greeting struct {
		// Hello :: func(name string) string
		Hello *openapigen.Action
	}

	People struct {
		// ListPerson :: func(...) []PersonSummary
		ListPerson *openapigen.Action
		// CreatePerson :: func(...)
		CreatePerson *openapigen.Action
	}
}) {
	Definitions := NewDefinitions(b)

	Actions.Greeting.Hello = b.Action("hello",
		b.Input(
			b.Param("name", b.String()).AsPath(),
		).Doc("input"),
		b.Output(
			b.String(),
		),
	).Doc("greeting hello")

	Actions.People.ListPerson = b.Action("ListPerson",
		b.Input(
			b.Param("sort", b.String().Enum([]string{"name", "-name", "age", "-age"})).AsQuery(),
		),
		b.Output(b.Array(Definitions.PersonSummary)).Doc("list of person summary"),
	).Doc("list person")

	Actions.People.CreatePerson = b.Action("CreatePerson",
		b.Input(
			b.Param("verbose", b.Bool()).AsQuery(),
			b.Body(b.Object(
				append(Definitions.Person.IgnoreFields("id", "father", "friends"),
					b.Field("fatherId", b.String()),
					b.Field("friendIdList", b.Array(b.String())))...,
			)).Doc("person but father and friends are id"),
		),
		b.Output(nil).Status(204),
	).Doc("create person")

	return
}
