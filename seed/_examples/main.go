package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"text/template"

	"github.com/podhmo/gos/seed"
)

var options struct {
	Builder  bool
	Metadata bool

	All bool
}

func main() {
	flag.BoolVar(&options.Builder, "builder", false, "emit builder.go")
	flag.BoolVar(&options.Metadata, "metadata", false, "emit builder.go")
	flag.BoolVar(&options.All, "all", false, "emit all")
	flag.Parse()

	if options.All {
		options.Builder = true
		options.Metadata = true
	}

	if err := run(); err != nil {
		log.Fatalf("!! %+v", err)
	}
}

func run() error {
	b := seed.NewBuilder()

	Type := b.BuildTarget("Type")
	goStringType := seed.Symbol("string") // unexported string is go-primitive

	Int := b.Type("Int").NeedBuilder()
	String := b.Type("String").NeedBuilder()

	Array := b.Type("Array").Var("typ", Type).NeedBuilder()
	Map := b.Type("Map").Var("valtype", Type).NeedBuilder()

	Field := b.Type("Field")
	FieldType := Field.Metadata.Name
	Object := b.Type("Object").Field("Fields", FieldType).NeedBuilder()

	Param := b.Type("Param")
	ParamType := Param.Metadata.Name
	ActionInput := b.Type("ActionInput").Field("Params", ParamType)
	ActionOutput := b.Type("ActionOutput").Field("Return", Type)
	Action := b.Type("Action").Field("Name", goStringType).Field("Input", ActionInput.Metadata.Name).Field("Output", ActionOutput.Metadata.Name).NeedBuilder()

	// emit
	{
		tmpl := seed.Template

		fmt.Fprintln(os.Stderr, Type, Int, String, Array, Map, Field, Object, Param, ActionInput, ActionOutput, Action)
		fmt.Fprintln(os.Stderr, b.Metadata.Types)
		t := template.Must(template.New("").Parse(tmpl))

		if options.Builder {
			fmt.Fprintln(os.Stderr, "--builder.go----------------------------------------")
			if err := t.ExecuteTemplate(os.Stdout, "Builder", b.Metadata); err != nil {
				panic(err)
			}
		}

		if options.Metadata {
			fmt.Fprintln(os.Stderr, "--metadata.go----------------------------------------")
			if err := t.ExecuteTemplate(os.Stdout, "Metadata", b.Metadata); err != nil {
				panic(err)
			}
		}
	}
	return nil
}
