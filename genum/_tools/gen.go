package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"text/template"

	"github.com/podhmo/gos/seed"
)

var options struct {
	PkgName string

	Builder  bool
	Metadata bool

	All   bool
	Write bool
}

func main() {
	flag.StringVar(&options.PkgName, "pkgname", "M", "package {{.PkgName}}")
	flag.BoolVar(&options.Builder, "builder", false, "emit builder.go")
	flag.BoolVar(&options.Metadata, "metadata", false, "emit metadata.go")
	flag.BoolVar(&options.All, "all", false, "emit all")
	flag.BoolVar(&options.Write, "write", false, "write file")
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
	b := seed.NewBuilder(options.PkgName)

	// define
	b.BuildTarget("Enum")

	b.Type("Int").NeedBuilder()
	b.Type("String").NeedBuilder()

	// TODO: value

	// emit
	{
		tmpl := seed.Template

		t := template.Must(template.New("").Parse(tmpl))

		if options.Builder {
			fmt.Fprintln(os.Stderr, "--builder.go----------------------------------------")
			var w io.Writer = os.Stdout
			if options.Write {
				f, err := os.Create("builder.go")
				if err != nil {
					return fmt.Errorf("create builder.go: %w", err)
				}
				defer f.Close()
				w = f
			}
			if err := t.ExecuteTemplate(w, "Builder", b.Metadata); err != nil {
				return fmt.Errorf("write builder.go: %w", err)
			}
		}

		if options.Metadata {
			fmt.Fprintln(os.Stderr, "--metadata.go----------------------------------------")
			var w io.Writer = os.Stdout
			if options.Write {
				f, err := os.Create("metadata.go")
				if err != nil {
					return fmt.Errorf("create metadata.go: %w", err)
				}
				defer f.Close()
				w = f
			}
			if err := t.ExecuteTemplate(w, "Metadata", b.Metadata); err != nil {
				return fmt.Errorf("write metadata.go: %w", err)
			}
		}
	}
	return nil
}
