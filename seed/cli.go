package seed

import (
	_ "embed"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"
)

//go:embed builder.tmpl
var Template string

type Command struct {
	*Config

	Template string
	FuncMap  template.FuncMap
	fs       *flag.FlagSet
}

type Config struct {
	PkgName string

	Builder  bool
	Metadata bool
	Stringer bool

	All   bool
	Write bool
}

func NewCommand(args []string) *Command {
	config := Config{}

	fs := flag.NewFlagSet("seed", flag.PanicOnError)
	fs.StringVar(&config.PkgName, "pkgname", "M", "package {{.PkgName}}")
	fs.BoolVar(&config.Builder, "builder", false, "emit builder.go")
	fs.BoolVar(&config.Metadata, "metadata", false, "emit metadata.go")
	fs.BoolVar(&config.Stringer, "stringer", false, "emit stringer.go")
	fs.BoolVar(&config.All, "all", false, "emit all")
	fs.BoolVar(&config.Write, "write", false, "write file")

	fs.Parse(args)
	if config.All {
		config.Builder = true
		config.Metadata = true
		config.Stringer = true
	}

	funcMap := template.FuncMap{
		"toLower": strings.ToLower,
		"toUpper": strings.ToUpper,
	}
	return &Command{Config: &config, Template: Template, FuncMap: funcMap, fs: fs}
}

func (c *Command) Do(b *Builder) error {
	options := c.Config
	b.metadata.NeedStringer = options.Stringer

	t := template.Must(template.New("").Funcs(c.FuncMap).Parse(c.Template))

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
		if err := t.ExecuteTemplate(w, "Builder", b.metadata); err != nil {
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
		if err := t.ExecuteTemplate(w, "Metadata", b.metadata); err != nil {
			return fmt.Errorf("write metadata.go: %w", err)
		}
	}

	if options.Stringer {
		fmt.Fprintln(os.Stderr, "--stringer.go----------------------------------------")
		var w io.Writer = os.Stdout
		if options.Write {
			f, err := os.Create("stringer.go")
			if err != nil {
				return fmt.Errorf("create stringer.go: %w", err)
			}
			defer f.Close()
			w = f
		}
		if err := t.ExecuteTemplate(w, "Stringer", b.metadata); err != nil {
			return fmt.Errorf("write stringer.go: %w", err)
		}
	}
	return nil
}
