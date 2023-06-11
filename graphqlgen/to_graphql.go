package graphqlgen

import (
	_ "embed"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"

	"github.com/podhmo/gos/pkg/namelib"
)

//go:embed graphql.tmpl
var tmpl string

func graphqlType(t Type) string {
	return t.GetTypeMetadata().underlying
}

func ToGraphql(w io.Writer, b *Builder) error {
	funcMap := template.FuncMap{
		"toTitle": namelib.ToTitle,
		"toGraphqlType": func(t Type) string {
			return graphqlType(t)
		},
		"splitLines": func(s string) []string {
			return strings.Split(s, "\n")
		},
	}
	t, err := template.New("").Funcs(funcMap).Parse(tmpl)
	if err != nil {
		return err
	}

	w = os.Stderr
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
