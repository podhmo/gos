package openapigen

import (
	_ "embed"
	"fmt"
	"io"
	"strings"
	"text/template"

	"github.com/podhmo/gos/pkg/namelib"
)

//go:embed to_code.tmpl
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

func ToGocode(w io.Writer, b *Builder) error {
	funcMap := template.FuncMap{
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
