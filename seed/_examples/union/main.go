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
	b.BuildTarget("Type")

	Param := b.Type("Param",
		b.Field("Name", seed.Symbol("string")),
		b.Field("Type", seed.Symbol("string")),
	).NeedBuilder()
	Body := b.Type("Body",
		b.Field("Type", seed.Symbol("string")),
	)

	InputArg := b.Union("inputArg", Param, Body)
	Input := b.Type("Input",
		b.Field("Params", seed.Symbol("[]*Param")),
		b.Field("Body", seed.Symbol("*Body")),
	).Constructor(
		b.Arg("Args", InputArg.GetMetadata().Name).Variadic().BindFields("Params", "Body").Transform(func(s string) string {
			return fmt.Sprintf(`func()(v1 []*Param, v2 *Body){
				for _, a := range %s {
					switch a := a.(type) {
					case *Param:
						v1 = append(v1, a)
					case *Body:
						v2 = a
					}
				}
				return
			}()`, s)
		}),
	).NeedBuilder()
	fmt.Fprintln(os.Stderr, Input)

	// emit
	return cmd.Do(b)
}
