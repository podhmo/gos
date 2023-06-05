//go:generate go run ./tools -write -builder -metadata -stringer -pkgname openapigen
//go:generate go fmt .
package openapigen

import (
	"fmt"
	"io"
	"strings"
	"text/template"

	_ "embed"

	"github.com/podhmo/gos/pkg/namelib"
)

type Config struct {
	DisableRefLinks bool // if true, does not use $ref links

	FuncMap  template.FuncMap
	Template *template.Template

	// for to schema
	defs []TypeBuilder
	seen map[int]*TypeRef
}

func DefaultConfig() *Config {
	c := &Config{
		DisableRefLinks: false,
		seen:            map[int]*TypeRef{},
		FuncMap: template.FuncMap{
			"toTitle": namelib.ToTitle,
			"toType": func(t Type) string {
				return typeString(t, false)
			},
			"toTypeInternal": func(t Type) string {
				return typeString(t, true)
			},
			"splitLines": func(s string) []string {
				return strings.Split(s, "\n")
			},
		},
	}
	c.Template = template.Must(template.New("").Funcs(c.FuncMap).Parse(tmpl))
	return c
}

//go:embed gocode.tmpl
var tmpl string

func typeString(t Type, internal bool) string {
	metadata := t.GetTypeMetadata()
	if !internal {
		if named := metadata.id > 0; named {
			return metadata.Name
		}
	}

	switch impl := t.(type) {
	case *Array[Type]:
		return fmt.Sprintf("[]%s", typeString(impl.items, false))
	case *Map[Type]:
		return fmt.Sprintf("map[string]%s", typeString(impl.items, false))
	default:
		return metadata.goType
	}
}

func (c *Config) ToGoCode(w io.Writer, tb TypeBuilder) error {
	t := c.Template
	switch tb := tb.(type) {
	case *Object:
		if err := t.ExecuteTemplate(w, "Object", tb); err != nil {
			return fmt.Errorf("execute template: %w", err)
		}
	default:
		if err := t.ExecuteTemplate(w, "Type", tb); err != nil {
			return fmt.Errorf("execute template: %w", err)
		}
	}
	return nil
}

func ToGocode(w io.Writer, b *Builder) error {
	c := b.Config
	return b.EachTypes(func(tb TypeBuilder) error {
		return c.ToGoCode(w, tb)
	})
}
