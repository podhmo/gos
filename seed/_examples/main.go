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

	fmt.Fprintln(os.Stderr, Type, Int, String, Array, Map, Field, Object, Param, ActionInput, ActionOutput, Action)
	fmt.Fprintln(os.Stderr, b.Metadata.Types)

	// emit
	return cmd.Do(b)
}
