//go:generate go run ./tools -write -builder -metadata -stringer -pkgname enumgen
//go:generate go fmt .
package enumgen

import (
	_ "embed"
	"fmt"
	"io"
	"strings"
	"text/template"

	"github.com/podhmo/gos/pkg/namelib"
)

type Config struct {
	FuncMap  template.FuncMap
	Template *template.Template
}

func DefaultConfig() *Config {
	c := &Config{
		FuncMap: template.FuncMap{
			"splitLines": func(s string) []string {
				return strings.Split(s, "\n")
			},
			"toType": func(t EnumBuilder) string {
				return t.GetEnumMetadata().underlying
			},
			"toTitle": namelib.ToTitle,
		},
	}
	c.Template = template.Must(template.New("").Funcs(c.FuncMap).Parse(tmpl))
	return c
}

//go:embed gocode.tmpl
var tmpl string

func (c *Config) ToGoCode(w io.Writer, tb EnumBuilder) error {
	switch tb := tb.(type) {
	case *String:
		for _, v := range tb.metadata.Members {
			if v.metadata.Name == "" {
				v.metadata.Name = namelib.ToTitle(v.metadata.Value)
			}
		}
		if err := c.Template.ExecuteTemplate(w, "String", tb); err != nil {
			return fmt.Errorf("execute template: %w", err)
		}
	case *Int:
		if err := c.Template.ExecuteTemplate(w, "Int", tb); err != nil {
			return fmt.Errorf("execute template: %w", err)
		}
	default:
		panic(fmt.Sprintf("unexpected type: %T", tb))
	}
	return nil
}

func ToGocode(w io.Writer, b *Builder) error {
	return b.EachEnums(func(tb EnumBuilder) error {
		return b.Config.ToGoCode(w, tb)
	})
}
