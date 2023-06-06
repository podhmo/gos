package design

import (
	"github.com/podhmo/gos/openapigen"
	"github.com/podhmo/gos/pkg/callerinfo"
)

var collector = callerinfo.NewCollector(1 /* depth */)

func ListPerson() *openapigen.Action {
	return b.Action(collector.CallerName(), // extract function name
		b.Input(
			b.Param("sort", b.String().Enum([]string{"name", "-name", "age", "-age"})).AsQuery(),
		),
		b.Output(b.Array(PersonSummary)).Doc("list of person summary"),
	).Doc("list person")
}

// create person
func CreatePerson() *openapigen.Action {
	info := collector.Info()
	return b.Action(info.CallerName,
		b.Input(
			b.Param("verbose", b.Bool()).AsQuery(),
			b.Body(b.Object(
				append(Person.IgnoreFields("id", "father", "friends"),
					b.Field("fatherId", b.String()),
					b.Field("friendIdList", b.Array(b.String())))...,
			)).Doc("person but father and friends are id"),
		),
		b.Output(nil).Status(204),
	).Doc(info.CallerDoc) // extract doc from doc string of this function
}
