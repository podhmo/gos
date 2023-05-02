package builder

import (
	"fmt"
	"io"
	"strings"
)

func ToString(typ TypeBuilder) string {
	b := new(strings.Builder)
	if err := typ.WriteType(b); err != nil {
		return fmt.Sprintf("invalid type: %T", typ)
	}
	return b.String()
}

// customization

func (b ObjectBuilder[R]) WriteType(w io.Writer) error {
	if err := b.type_.WriteType(w); err != nil {
		return err
	}
	// if b.type_.value.IsNewType {
	// 	return nil
	// }

	io.WriteString(w, "{") // nolint
	n := len(b.Fields) - 1
	for i, f := range b.Fields {
		v := f.fieldvalue()
		io.WriteString(w, v.Name) // nolint
		if !v.Required {
			io.WriteString(w, "?") // nolint
		}
		if i < n {
			io.WriteString(w, ", ") // nolint
		}
	}
	io.WriteString(w, "}") // nolint
	return nil
}

func (t *ArrayBuilder[T, R]) WriteType(w io.Writer) error {
	if err := t.type_.WriteType(w); err != nil {
		return err
	}
	// if t.type_.value.IsNewType {
	// 	return nil
	// }

	io.WriteString(w, "[") // nolint
	if err := t.items.WriteType(w); err != nil {
		return err
	}
	io.WriteString(w, "]") // nolint
	return nil
}

func (t *MapBuilder[V, R]) WriteType(w io.Writer) error {
	if err := t.type_.WriteType(w); err != nil {
		return err
	}
	// if t.type_.value.IsNewType {
	// 	return nil
	// }

	io.WriteString(w, "[") // nolint
	if err := t.items.WriteType(w); err != nil {
		return err
	}
	io.WriteString(w, "]") // nolint
	return nil
}
