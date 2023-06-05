package main

import (
	"encoding/json"
	"os"

	"github.com/iancoleman/orderedmap"
	"github.com/podhmo/gos/openapigen"
	design "github.com/podhmo/gos/openapigen/_examples/02separated/design"
	"github.com/podhmo/gos/pkg/maplib"
)

func main() {
	b := design.Builder

	// routing
	Error := openapigen.Define("Error", b.Object(
		b.Field("message", b.String()),
	))
	r := openapigen.NewRouter(Error)
	{
		r := r.Tagged("greeting")
		r.Post("/hello/{name}", design.Actions.Hello)
	}
	{
		r := r.Tagged("people")
		r.Get("/people", design.Actions.ListPerson)
		r.Post("/people", design.Actions.CreatePerson)
	}

	// emit
	doc, _ := maplib.Merge(orderedmap.New(), &openapigen.OpenAPI{
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
	})
	r.ToSchemaWith(b, doc)

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(doc); err != nil {
		panic(err)
	}
}
