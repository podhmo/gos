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

	Type := b.BuildTarget("Type")
	b.TargetField("Format", seed.Symbol("string"), `json:"format"`)

	Bool := b.Type("Bool").
		NeedBuilder().Underlying("boolean")
	Int := b.Type("Int").
		Field("Maximum", seed.Symbol("int64"), `json:"maximum,omitempty"`).
		Field("Minimum", seed.Symbol("int64"), `json:"minimum,omitempty"`).
		NeedBuilder().Underlying("intger")
	// String := b.Type("String").
	// 	Field("Pattern", seed.Symbol("string"), `json:"pattern,omitempty"`).
	// 	Field("MaxLength", seed.Symbol("int64"), `json:"maxlength,omitempty"`).
	// 	Field("MinLength", seed.Symbol("int64"), `json:"minlength,omitempty"`).
	// 	NeedBuilder().Underlying("string")

	// Array := b.Type("Array", seed.TypeVar{Name: "Items", Type: seed.Symbol("TypeBuilder")}).
	// 	Field("MaxItems", seed.Symbol("int64"), `json:"maxitems,omitempty"`).
	// 	Field("MinItems", seed.Symbol("int64"), `json:"minitems,omitempty"`).
	// 	NeedBuilder().Underlying("array")
	// Map := b.Type("Map", seed.TypeVar{Name: "Items", Type: seed.Symbol("TypeBuilder")}).
	// 	Field("Pattern", seed.Symbol("string"), `json:"pattern,omitempty"`).
	// 	NeedBuilder().Underlying("map")

	// Field := b.Type("Field").
	// 	Field("Name", seed.Symbol("string"), `json:"-"`).
	// 	Field("Typ", seed.Symbol("TypeBuilder"), `json:"-"`).
	// 	Field("Description", seed.Symbol("string"), `json:"description,omitempty"`).
	// 	Field("Required", seed.Symbol("bool"), `json:"-"`).
	// 	Constructor(
	// 		seed.Arg{Name: "Name", Type: seed.Symbol("string")},
	// 		seed.Arg{Name: "Typ", Type: seed.Symbol("TypeBuilder")},
	// 	).
	// 	NeedBuilder().Underlying("field") //?
	// Object := b.Type("Object").
	// 	Field("Fields", seed.Symbol("[]*FieldType"), `json:"-"`).
	// 	Field("Strict", seed.Symbol("bool"), `json:"-"`).
	// 	Constructor(seed.Arg{Name: "Fields", Type: seed.Symbol("*FieldType"), Variadic: true}).
	// 	NeedBuilder().Underlying("object")

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

	fmt.Fprintln(os.Stderr, Type, Bool, Int)
	// fmt.Fprintln(os.Stderr, Type, Bool, Int, String, Array, Map, Field, Object)
	// fmt.Fprintln(os.Stderr, Param, ActionInput, ActionOutput, Action)
	// fmt.Fprintln(os.Stderr, b.Metadata.Types)

	// emit
	return cmd.Do(b)
}
