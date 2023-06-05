package openapigen

import (
	_ "embed"
	"fmt"
	"io"
	"strings"
	"text/template"
)

//go:embed to_code.tmpl
var tmpl string

func ToGocode(w io.Writer, b *Builder) error {
	funcMap := template.FuncMap{
		"toTitle": func(s string) string { // TODO: goify
			if s == "" {
				return ""
			}
			return strings.ToUpper(s[:1]) + s[1:]
		},
		"toType": func(t Type) string {
			metadata := t.GetTypeMetadata()
			if named := metadata.id > 0; !named {
				return metadata.goType
			}
			return metadata.Name
		},
		"splitLines": func(s string) []string {
			return strings.Split(s, "\n")
		},
	}
	t, err := template.New("").Funcs(funcMap).Parse(tmpl)
	if err != nil {
		return fmt.Errorf("load template: %w", err)
	}

	return b.EachTypes(func(tb TypeBuilder) error {
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
	})
}
