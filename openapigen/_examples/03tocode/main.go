package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"github.com/podhmo/gos/openapigen"
)

func main() {
	b := openapigen.NewBuilder(openapigen.DefaultConfig())

	Name := openapigen.Define("Name", b.String()).
		Doc("name of something")

	Friends := openapigen.Define("Friends", b.Array(b.ReferenceByName("Person"))).
		Doc("Friends of something")

	Tag := openapigen.Define("Tag", b.Object(
		b.Field("name", b.Reference(Name)).Doc("name of tag"),
		b.Field("doc", b.String()),
	))

	openapigen.Define("Person", b.Object(
		b.Field("name", b.Reference(Name)).Doc("name of person"),
		b.Field("age", b.Int()),
		b.Field("father", b.ReferenceByName("Person")).Nullable(true),
		b.Field("children", b.Array(b.ReferenceByName("Person"))),
		b.Field("friends", b.Reference(Friends)),
		b.Field("tags", b.Map(b.Reference(Tag))),
	)).Doc("person object")

	// emit
	_, file, _, _ := runtime.Caller(0)
	cwd := filepath.Dir(file)
	{
		outname := "testdata/generated.go"
		fmt.Fprintln(os.Stderr, "write", outname)
		f, err := os.Create(filepath.Join(cwd, outname))
		if err != nil {
			panic(err)
		}
		defer f.Close()
		var w io.Writer = f
		fmt.Fprintln(w, "package M")
		if err := openapigen.ToGocode(w, b); err != nil {
			panic(err)
		}
	}
	{
		outname := "testdata/openapi.json"
		fmt.Fprintln(os.Stderr, "write", outname)
		f, err := os.Create(filepath.Join(cwd, outname))
		if err != nil {
			panic(err)
		}
		defer f.Close()

		doc, err := openapigen.ToSchema(b)
		if err != nil {
			panic(err)
		}

		enc := json.NewEncoder(f)
		enc.SetIndent("", "  ")
		if err := enc.Encode(doc); err != nil {
			panic(err)
		}
	}
}
