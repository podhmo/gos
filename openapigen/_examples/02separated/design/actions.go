package design

import (
	"github.com/podhmo/gos/openapigen"
	"github.com/podhmo/gos/pkg/funcinfo"
)

func NewHandler(b *openapigen.Builder) *Handler {
	return &Handler{b: b, collector: funcinfo.NewCollector(1)}
}

type Handler struct {
	b         *openapigen.Builder
	collector *funcinfo.Collector
}

// greeting hello
func (h *Handler) Hello() *openapigen.Action {
	info := h.collector.Collect()
	return b.Action(info.Name,
		b.Input(
			b.Param("name", b.String()).AsPath(),
		).Doc("input"),
		b.Output(
			b.String(),
		),
	).Doc(info.Doc)
}

// list person
func (h *Handler) ListPerson() *openapigen.Action {
	info := h.collector.Collect()
	return b.Action(info.Name,
		b.Input(
			b.Param("sort", b.String().Enum([]string{"name", "-name", "age", "-age"})).AsQuery(),
		),
		b.Output(b.Array(PersonSummary)).Doc("list of person summary"),
	).Doc(info.Doc)
}

// create person
func (h *Handler) CreatePerson() *openapigen.Action {
	info := h.collector.Collect()
	return b.Action(info.Name,
		b.Input(
			b.Param("verbose", b.Bool()).AsQuery(),
			b.Body(b.Object(
				append(Person.IgnoreFields("id", "father", "friends"),
					b.Field("fatherId", b.String()),
					b.Field("friendIdList", b.Array(b.String())))...,
			)).Doc("person but father and friends are id"),
		),
		b.Output(nil).Status(204),
	).Doc(info.Doc)
}
