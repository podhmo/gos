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

	goStringType := seed.Symbol("string") // unexported string is go-primitive

	// define
	b := seed.NewBuilder(options.PkgName)
	b.NeedReference()

	Type := b.BuildTarget("Type")
	b.TargetField("Format", goStringType, `json:"format"`)

	Int := b.Type("Int").
		NeedBuilder().Underlying("intger")
	String := b.Type("String").
		NeedBuilder().Underlying("string")

	Array := b.Type("Array", seed.TypeVar{Name: "Items", Type: seed.Symbol("TypeBuilder")}).
		NeedBuilder().Underlying("array")
	Map := b.Type("Map", seed.TypeVar{Name: "Items", Type: seed.Symbol("TypeBuilder")}).
		NeedBuilder().Underlying("map")

	Field := b.Type("Field").
		NeedBuilder().Underlying("field") //?
	Object := b.Type("Object").
		Field("Fields", seed.Symbol("[]*FieldType"), `json:"-"`).
		Constructor(seed.Arg{Name: "Fields", Type: seed.Symbol("*FieldType"), Variadic: true}).
		NeedBuilder().Underlying("object")

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

	fmt.Fprintln(os.Stderr, Type, Int, String, Array, Map, Field, Object)
	// fmt.Fprintln(os.Stderr, Param, ActionInput, ActionOutput, Action)
	// fmt.Fprintln(os.Stderr, b.Metadata.Types)

	// emit
	return cmd.Do(b)
}
