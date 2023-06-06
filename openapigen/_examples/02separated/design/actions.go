package design

import "github.com/podhmo/gos/openapigen"

func NewHandler(b *openapigen.Builder) *Handler {
	return &Handler{b: b}
}

type Handler struct {
	b *openapigen.Builder
}

// Hello :: func(name string) string
func (h *Handler) Hello() *openapigen.Action {
	return b.Action("hello",
		b.Input(
			b.Param("name", b.String()).AsPath(),
		).Doc("input"),
		b.Output(
			b.String(),
		),
	).Doc("greeting hello")
}

// ListPerson :: func(...) []PersonSummary
func (h *Handler) ListPerson() *openapigen.Action {
	return b.Action("ListPerson",
		b.Input(
			b.Param("sort", b.String().Enum([]string{"name", "-name", "age", "-age"})).AsQuery(),
		),
		b.Output(b.Array(PersonSummary)).Doc("list of person summary"),
	).Doc("list person")
}

// CreatePerson :: func(...) ()
func (h *Handler) CreatePerson() *openapigen.Action {
	return b.Action("CreatePerson",
		b.Input(
			b.Param("verbose", b.Bool()).AsQuery(),
			b.Body(b.Object(
				append(Person.IgnoreFields("id", "father", "friends"),
					b.Field("fatherId", b.String()),
					b.Field("friendIdList", b.Array(b.String())))...,
			)).Doc("person but father and friends are id"),
		),
		b.Output(nil).Status(204),
	).Doc("create person")
}
