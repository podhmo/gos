package main

import (
	"encoding/json"
	"os"

	"github.com/iancoleman/orderedmap"
	"github.com/podhmo/gos/openapigen"
)

func main() {
	b := openapigen.NewTypeBuilder(openapigen.DefaultConfig())

	Name := openapigen.DefineType("Name", b.String()).Doc("name of something")
	openapigen.DefineType("DateTime", b.String().Format("date-time")) // for ReferenceByName

	Task := openapigen.DefineType("Task", b.Object(
		b.Field("name", b.Reference(Name)),
		b.Field("done", b.Bool()),
		b.Field("createdAt", b.ReferenceByName("DateTime")),
	))

	ListTask := b.Action("ListTask",
		b.Input(b.Param("sort", b.String().Enum([]string{"createdAt", "-createdAt"})).AsQuery()),
		b.Output(b.Array(Task)),
	)

	Error := openapigen.DefineType("Error", b.Object(
		b.Field("message", b.String()),
	)).Doc("default error")
	r := openapigen.NewRouter(Error)
	{
		r := r.Tagged("task")
		r.Get("/tasks", ListTask)
	}

	doc := orderedmap.New()
	r.ToSchemaWith(b, doc)
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(doc); err != nil {
		panic(err)
	}
}
