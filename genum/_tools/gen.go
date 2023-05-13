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

	// define
	b.BuildTarget("Enum")

	b.Type("Int").NeedBuilder()
	b.Type("String").NeedBuilder()

	// TODO: value

	// emit
	if err := cmd.Do(b); err != nil {
		return fmt.Errorf("emit: %w", err)
	}
	return nil
}
