package main

import (
	"fmt"
	"log"
	"os"

	"github.com/podhmo/gos/seed"
)

func main() {
	cmd := seed.NewCommand(os.Args[1:])
	if err := run(cmd); err != nil {
		log.Fatalf("!! %+v", err)
	}
}

func run(cmd *seed.Command) error {
	options := cmd.Config
	b := seed.NewBuilder(options.PkgName,
		seed.Root.Field("Config", seed.Symbol("*Config")),
	)
	b.GeneratedBy("github.com/podhmo/gos/genum/tools")

	// define
	b.BuildTarget("Enum")
	b.InterfaceMethods(`writeCoder // see: to_code.go`)

	b.Constructor(
		b.Arg("Config", seed.Symbol("*Config")),
	)

	b.Type("Int",
		b.Field("Default", seed.Symbol("int")).Tag(`json:"default"`),
		b.Field("Members", seed.Symbol("[]IntValue")),
	).Constructor(
		b.Arg("Members", seed.Symbol("IntValue")).Variadic(),
	).NeedBuilder().Underlying("int")
	b.Type("IntValue",
		b.Field("Name", seed.Symbol("string")),
		b.Field("Value", seed.Symbol("int")),
		b.Field("Doc", seed.Symbol("string")),
	)

	b.Type("String",
		b.Field("Default", seed.Symbol("string")).Tag(`json:"default"`),
		b.Field("Members", seed.Symbol("[]StringValue")),
	).Constructor(
		b.Arg("Members", seed.Symbol("StringValue")).Variadic(),
	).NeedBuilder().Underlying("string")
	b.Type("StringValue",
		b.Field("Name", seed.Symbol("string")),
		b.Field("Value", seed.Symbol("string")),
		b.Field("Doc", seed.Symbol("string")),
	)

	// emit
	if err := cmd.Do(b); err != nil {
		return fmt.Errorf("emit: %w", err)
	}
	return nil
}
