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

	b.Constructor(
		b.Arg("Config", seed.Symbol("*Config")),
	)

	// int
	b.Type("IntEnum",
		b.Field("Default", seed.Symbol("int")).Tag(`json:"default"`),
		b.Field("Members", seed.Symbol("[]*IntValueMetadata")),
	).Constructor(
		b.Arg("Members", seed.Symbol("*IntValue")).Variadic().Transform(func(v string) string {
			return fmt.Sprintf("mapslice(%s, func(x *IntValue) *IntValueMetadata { return x.metadata})", v)
		}),
	).NeedBuilder().Underlying("int") // generate Int, IntMetadata

	b.Type("IntValue",
		b.Field("Name", seed.Symbol("string")),
		b.Field("Value", seed.Symbol("int")),
		b.Field("Doc", seed.Symbol("string")),
	).Constructor(
		b.Arg("Name", seed.Symbol("string")),
		b.Arg("Value", seed.Symbol("int")),
	).NeedBuilder() // generate IntValue, IntValueMetadata

	// string
	b.Type("StringEnum",
		b.Field("Default", seed.Symbol("string")).Tag(`json:"default"`),
		b.Field("Members", seed.Symbol("[]*StringValueMetadata")),
	).Constructor(
		b.Arg("Members", seed.Symbol("*StringValue")).Variadic().Transform(func(v string) string {
			return fmt.Sprintf("mapslice(%s, func(x *StringValue) *StringValueMetadata { return x.metadata})", v)
		}),
	).NeedBuilder().Underlying("string")
	b.Type("StringValue",
		b.Field("Name", seed.Symbol("string")),
		b.Field("Value", seed.Symbol("string")),
		b.Field("Doc", seed.Symbol("string")),
	).Constructor(
		b.Arg("Value", seed.Symbol("string")),
	).NeedBuilder()

	// for transform
	b.Footer(`
	func mapslice[S, D any](src []S, conv func(S) D) []D {
		dst := make([]D, len(src))
		for i, x := range src {
			dst[i] = conv(x)
		}
		return dst
	}	
	`)

	// emit
	if err := cmd.Do(b); err != nil {
		return fmt.Errorf("emit: %w", err)
	}
	return nil
}
