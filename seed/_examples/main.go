package main

import (
	"fmt"
	"log"
	"os"

	"github.com/podhmo/gos/seed"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("!! %+v", err)
	}
}

func run() error {
	cmd := seed.NewCommand(os.Args[1:])
	options := cmd.Config

	// define
	b := seed.NewBuilder(options.PkgName)
	b.NeedReference()

	Type := b.BuildTarget("Type",
		b.Field("Format", seed.Symbol("string")).Tag(`json:"format"`),
	)

	Bool := b.Type("Bool").
		NeedBuilder().Underlying("boolean")
	Int := b.Type("Int",
		b.Field("Maximum", seed.Symbol("int64")).Tag(`json:"maximum,omitempty"`),
		b.Field("Minimum", seed.Symbol("int64")).Tag(`json:"minimum,omitempty"`),
	).NeedBuilder().Underlying("integer")
	String := b.Type("String",
		b.Field("Pattern", seed.Symbol("string")).Tag(`json:"pattern,omitempty"`),
		b.Field("MaxLength", seed.Symbol("int64")).Tag(`json:"maxlength,omitempty"`),
		b.Field("MinLength", seed.Symbol("int64")).Tag(`json:"minlength,omitempty"`),
	).NeedBuilder().Underlying("string")

	Array := b.Type("Array",
		b.TypeVar("Items", seed.Symbol("TypeBuilder")),
		b.Field("MaxItems", seed.Symbol("int64")).Tag(`json:"maxitems,omitempty"`),
		b.Field("MinItems", seed.Symbol("int64")).Tag(`json:"minitems,omitempty"`),
	).NeedBuilder().Underlying("array")
	Map := b.Type("Map", b.TypeVar("Items", seed.Symbol("TypeBuilder")),
		b.Field("Pattern", seed.Symbol("string")).Tag(`json:"pattern,omitempty"`),
	).NeedBuilder().Underlying("map")

	Field := b.Type("Field",
		b.Field("Name", seed.Symbol("string")).Tag(`json:"-"`),
		b.Field("Typ", seed.Symbol("TypeBuilder")).Tag(`json:"-"`),
		b.Field("Description", seed.Symbol("string")).Tag(`json:"description,omitempty"`),
		b.Field("Required", seed.Symbol("bool")).Tag(`json:"-"`),
	).Constructor(
		b.Arg("Name", seed.Symbol("string")),
		b.Arg("Typ", seed.Symbol("TypeBuilder")),
	).NeedBuilder().Underlying("field") //?

	Object := b.Type("Object",
		b.Field("Fields", seed.Symbol("[]*FieldType")).Tag(`json:"-"`),
		b.Field("Strict", seed.Symbol("bool")).Tag(`json:"-"`),
	).Constructor(
		b.Arg("Fields", seed.Symbol("*FieldType")).Variadic(),
	).NeedBuilder().Underlying("object")

	// Param := b.Type("Param")
	// ParamType := Param.Metadata.Name
	// ActionInput := b.Type("ActionInput").
	// 	Field("Params", ParamType, "")
	// ActionOutput := b.Type("ActionOutput").
	// 	Field("Return", Type, "")
	// Action := b.Type("Action").
	// 	Field("Name", goStringType, "").
	// 	Field("Input", ActionInput.Metadata.Name, "").
	// 	Field("Output", ActionOutput.Metadata.Name, "").
	// 	NeedBuilder()

	fmt.Fprintln(os.Stderr, Type, Bool, Int, String, Array, Map, Field, Object)
	// fmt.Fprintln(os.Stderr, Param, ActionInput, ActionOutput, Action)
	// fmt.Fprintln(os.Stderr, b.Metadata.Types)

	// emit
	return cmd.Do(b)
}
