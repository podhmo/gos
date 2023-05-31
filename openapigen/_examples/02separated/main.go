package main

import (
	"encoding/json"
	"os"

	"github.com/iancoleman/orderedmap"
	"github.com/podhmo/gos/openapigen"
	design "github.com/podhmo/gos/openapigen/_examples/02separated/design"
	action "github.com/podhmo/gos/openapigen/_examples/02separated/design/action"
)

func main() {
	b := design.Builder

	// routing
	Error := openapigen.DefineType("Error", b.Object(
		b.Field("message", b.String()),
	))
	r := openapigen.NewRouter(Error)
	{
		r := r.Tagged("greeting")
		r.Post("/hello/{name}", action.Hello)
	}
	{
		r := r.Tagged("people")
		r.Get("/people", action.ListPerson)
		r.Post("/people", action.CreatePerson)
	}

	// emit
	doc := orderedmap.New()
	r.ToSchemaWith(b, doc)

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(doc); err != nil {
		panic(err)
	}
}
