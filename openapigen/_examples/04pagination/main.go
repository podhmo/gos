package main

import (
	"encoding/json"
	"os"

	"github.com/iancoleman/orderedmap"
	"github.com/podhmo/gos/openapigen"
	"github.com/podhmo/gos/pkg/maplib"
)

// with pagination
func Pagination(b *openapigen.Builder, typ openapigen.Type) *openapigen.Object {
	return b.Object(
		b.Field("totalCount", b.Int()),
		b.Field("hasMore", b.Bool()),
		b.Field("cursor", b.String()),
		b.Field("data", typ).Doc("response data of api"),
	).Doc("totalCount, hasMore, cursor fields are metadata for pagination")
}

func Pagination2(b *openapigen.Builder, typ openapigen.Type) *openapigen.Object {
	return b.Object(
		b.Field("_meta", b.Object(
			b.Field("totalCount", b.Int()),
			b.Field("hasMore", b.Bool()),
			b.Field("cursor", b.String()),
		)),
		b.Field("data", typ).Doc("response data of api"),
	)
}

func main() {
	b := openapigen.NewBuilder(openapigen.DefaultConfig())

	// types
	Name := openapigen.Define("Name", b.String()).Doc("name of something")
	openapigen.Define("DateTime", b.String().Format("date-time")) // for ReferenceByName

	Task := openapigen.Define("Task", b.Object(
		b.Field("name", b.Reference(Name)),
		b.Field("done", b.Bool()),
		b.Field("createdAt", b.ReferenceByName("DateTime")),
	))

	ListTask := b.Action("ListTask",
		b.Input(b.Param("sort", b.String().Enum([]string{"createdAt", "-createdAt"})).AsQuery()),
		b.Output(Pagination(b, b.Array(Task))),
	).Doc("paginated list task")

	ListTask2 := b.Action("ListTask2",
		b.Input(b.Param("sort", b.String().Enum([]string{"createdAt", "-createdAt"})).AsQuery()),
		b.Output(Pagination2(b, b.Array(Task))),
	).Doc("paginated list task version 2")

	// routing
	Error := openapigen.Define("Error", b.Object(
		b.Field("message", b.String()),
	)).Doc("default error")
	r := openapigen.NewRouter(Error)
	{
		r := r.Tagged("task")
		r.Get("/tasks", ListTask)
		r.Get("/tasks2", ListTask2)
	}

	// openapi data
	openapi := &openapigen.OpenAPI{
		OpenAPI: "3.0.3",
		Info: openapigen.Info{
			Title:   "task API",
			Version: "0.0.0",
			Doc:     "simple list tasks API",
		},
		Servers: []openapigen.Server{
			{
				URL: "http://localhost:8080",
				Doc: "local development",
			},
		},
	}
	doc := orderedmap.New()
	doc, err := maplib.Merge(doc, openapi)
	if err != nil {
		panic(err)
	}

	r.ToSchemaWith(b, doc)
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(doc); err != nil {
		panic(err)
	}
}
