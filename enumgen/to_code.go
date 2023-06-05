package enumgen

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
		"splitLines": func(s string) []string {
			return strings.Split(s, "\n")
		},
		"toType": func(t EnumBuilder) string {
			return t.GetEnumMetadata().underlying
		},
		"toTitle": toTitle,
	}
	t, err := template.New("").Funcs(funcMap).Parse(tmpl)
	if err != nil {
		return fmt.Errorf("load template: %w", err)
	}

	return b.EachEnums(func(tb EnumBuilder) error {
		switch tb := tb.(type) {
		case *String:
			for _, v := range tb.metadata.Members {
				if v.Name == "" {
					v.Name = toTitle(v.Value)
				}
			}
			if err := t.ExecuteTemplate(w, "String", tb); err != nil {
				return fmt.Errorf("execute template: %w", err)
			}
		case *Int:
			if err := t.ExecuteTemplate(w, "Int", tb); err != nil {
				return fmt.Errorf("execute template: %w", err)
			}
		default:
			panic(fmt.Sprintf("unexpected type: %T", tb))
		}
		return nil
	})
}

func toTitle(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
