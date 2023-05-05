package genum

import (
	"fmt"
	"io"
	"reflect"
)

type writeCoder interface {
	writeCode(io.Writer) error
}

func WriteCode[T any](w io.Writer, b *Builder[T]) error {
	return b.EachTypes(func(typ TypeBuilder[T]) error {
		if err := typ.writeCode(w); err != nil {
			return fmt.Errorf("%s: %w", typ.GetTypeMetadata().Name, err)
		}
		fmt.Fprintln(w, "")
		return nil
	})
}

// customization

func (t *EnumType[T]) writeCode(w io.Writer) error {
	padding := t.rootbuilder.Config.Padding
	comment := t.rootbuilder.Config.Comment

	typename := t.type_.metadata.Name
	underlying := t.type_.metadata.underlying

	fmt.Fprintf(w, "%s %s", comment, typename) // nolint
	// TODO: description
	fmt.Fprintln(w, "")
	fmt.Fprintf(w, "type %s %s\n", typename, underlying) // nolint

	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "const (")

	isString := reflect.TypeOf(t.Values[0].metadata.Value).Kind() == reflect.String

	for _, v := range t.Values {
		if isString {
			fmt.Fprintf(w, "%s%s%s %s = \"%v\"", padding, typename, v.metadata.Name, typename, v.metadata.Value)
		} else {
			fmt.Fprintf(w, "%s%s%s %s = %v", padding, typename, v.metadata.Name, typename, v.metadata.Value)
		}

		if v.metadata.Default {
			fmt.Fprintf(w, "  %s default", comment)
		}
		fmt.Fprintln(w, "")
	}
	fmt.Fprintln(w, ")")
	return nil
}

func (t *type_[T, R]) writeCode(w io.Writer) error {
	panic("never")
}
