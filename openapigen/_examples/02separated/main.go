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
	r := openapigen.NewRouter(design.Error)
	{
		r := r.Tagged("greeting")

		hello := b.Action("hello",
			b.Input(
				b.Param("name", b.String()).AsPath(),
			).Doc("input"),
			b.Output(
				b.String(),
			),
		).Doc("greeting hello")
		r.Post("/hello/{name}", hello)
	}

	{

		r := r.Tagged("people")
		r.Get("/people", design.ListPerson())
		r.Post("/people", design.CreatePerson())
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
