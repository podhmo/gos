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
	b := seed.NewBuilder(options.PkgName)

	b.Import("io")

	// define
	b.BuildTarget("Enum")
	b.InterfaceMethods(`writeCode(io.Writer) error`)

	b.Field("Config", seed.Symbol("*Config"), "")
	b.Constructor(seed.Arg{Name: "Config", Type: seed.Symbol("*Config")})

	b.Type("Int").NeedBuilder().Underlying("int").
		Field("Default", seed.Symbol("int"), `json:"default"`).
		Field("Members", seed.Symbol("[]IntValue"), "").
		Constructor(seed.Arg{Name: "Members", Type: seed.Symbol("IntValue"), Variadic: true})
	b.Type("IntValue").
		Field("Name", seed.Symbol("string"), "").
		Field("Value", seed.Symbol("int"), "").
		Field("Doc", seed.Symbol("string"), "")

	b.Type("String").NeedBuilder().Underlying("string").
		Field("Default", seed.Symbol("string"), `json:"default"`).
		Field("Members", seed.Symbol("[]StringValue"), "").
		Constructor(seed.Arg{Name: "Members", Type: seed.Symbol("StringValue"), Variadic: true})
	b.Type("StringValue").
		Field("Name", seed.Symbol("string"), "").
		Field("Value", seed.Symbol("string"), "").
		Field("Doc", seed.Symbol("string"), "")

	// emit
	if err := cmd.Do(b); err != nil {
		return fmt.Errorf("emit: %w", err)
	}
	return nil
}
