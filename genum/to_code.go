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
	typename := t.type_.metadata.Name
	underlying := t.type_.metadata.underlying
	fmt.Fprintf(w, "// %s", typename) // nolint
	// TODO: description
	fmt.Fprintln(w, "")
	fmt.Fprintf(w, "type %s %s\n", typename, underlying) // nolint

	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "const (")

	isString := false
	if reflect.TypeOf(t.Values[0].metadata.Value).Kind() == reflect.String {
		isString = true
	}

	for _, v := range t.Values {
		if isString {
			fmt.Fprintf(w, "\t %s%s %s = \"%v\"", typename, v.metadata.Name, typename, v.metadata.Value)
		} else {
			fmt.Fprintf(w, "\t %s%s %s = %v", typename, v.metadata.Name, typename, v.metadata.Value)
		}

		if v.metadata.Default {
			fmt.Fprint(w, "  // default")
		}
		fmt.Fprintln(w, "")
	}
	fmt.Fprintln(w, ")")
	return nil
}

func (t *ValueType[T]) writeCode(w io.Writer) error {
	panic("hmm")
}
