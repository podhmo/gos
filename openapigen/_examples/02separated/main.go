package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/iancoleman/orderedmap"
	"github.com/podhmo/gos/openapigen"
	design "github.com/podhmo/gos/openapigen/_examples/02separated/design"
	action "github.com/podhmo/gos/openapigen/_examples/02separated/design/action"
)

func main() {
	b := design.Builder

	fmt.Fprintln(os.Stderr, "type  \t", design.Person)
	fmt.Fprintln(os.Stderr, "type  \t", design.PersonSummary)
	fmt.Fprintln(os.Stderr, "action\t", action.Hello)
	fmt.Fprintln(os.Stderr, "input \t", action.Hello.GetMetadata().Input)
	fmt.Fprintln(os.Stderr, "output\t", action.Hello.GetMetadata().Output)

	// routing
	r := openapigen.NewRouter()
	mount(r)

	doc := orderedmap.New()
	r.ToSchemaWith(b, doc)
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(doc); err != nil {
		panic(err)
	}
}

func mount(r *openapigen.Router) {
	{
		r := r.Tagged("greeting")
		r.Post("/hello/{name}", action.Hello)
	}
	{
		r := r.Tagged("people")
		r.Get("/people", action.ListPerson)
		r.Post("/people", action.CreatePerson)
	}
}
