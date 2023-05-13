package genum

import (
	"fmt"
	"io"
	"strings"
)

type writeCoder interface {
	writeCode(io.Writer) error
}

func WriteCode(w io.Writer, b *Builder) error {
	return b.EachEnums(func(typ EnumBuilder) error {
		if err := typ.writeCode(w); err != nil {
			return fmt.Errorf("%s: %w", typ.GetEnumMetadata().Name, err)
		}
		fmt.Fprintln(w, "")
		return nil
	})
}

// customization
func (t *StringEnum) writeCode(w io.Writer) error {
	// padding := t.rootbuilder.Config.Padding
	// comment := t.rootbuilder.Config.Comment
	padding := "\t"
	comment := "//"

	typename := t._Enum.metadata.Name
	underlying := t._Enum.metadata.underlying

	fmt.Fprintf(w, "%s %s", comment, typename) // nolint
	// TODO: description
	fmt.Fprintln(w, "")
	fmt.Fprintf(w, "type %s %s\n", typename, underlying) // nolint

	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "const (")

	for _, v := range t.metadata.Members {
		name := v.Name
		if name == "" {
			name = toTitle(v.Value)
		}

		if v.Doc != "" {
			fmt.Fprintf(w, "%s%s %q\n", padding, comment, v.Doc)
		}
		fmt.Fprintf(w, "%s%s%s %s = %q", padding, typename, name, typename, v.Value)
		if t.metadata.Default == v.Value {
			fmt.Fprintf(w, "  %s default", comment)
		}
		fmt.Fprintln(w, "")
	}
	fmt.Fprintln(w, ")")
	return nil
}

func (t *IntEnum) writeCode(w io.Writer) error {
	// padding := t.rootbuilder.Config.Padding
	// comment := t.rootbuilder.Config.Comment
	padding := "\t"
	comment := "//"

	typename := t._Enum.metadata.Name
	underlying := t._Enum.metadata.underlying

	fmt.Fprintf(w, "%s %s", comment, typename) // nolint
	// TODO: description
	fmt.Fprintln(w, "")
	fmt.Fprintf(w, "type %s %s\n", typename, underlying) // nolint

	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "const (")

	for _, v := range t.metadata.Members {
		name := v.Name

		if v.Doc != "" {
			fmt.Fprintf(w, "%s%s %q\n", padding, comment, v.Doc)
		}
		fmt.Fprintf(w, "%s%s%s %s = %v", padding, typename, name, typename, v.Value)
		if t.metadata.Default == v.Value {
			fmt.Fprintf(w, "  %s default", comment)
		}
		fmt.Fprintln(w, "")
	}
	fmt.Fprintln(w, ")")
	return nil
}

func toTitle(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
