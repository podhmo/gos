package builder

import (
	"fmt"
	"io"
	"strings"
)

type writeTyper interface {
	WriteType(io.Writer) error
}

func ToString(typ TypeBuilder) string {
	b := new(strings.Builder)
	if err := typ.WriteType(b); err != nil {
		return fmt.Sprintf("invalid type: %T", typ)
	}
	return b.String()
}

func (t *TypeRef) WriteType(w io.Writer) error {
	return t.getType().WriteType(w)
}

func (t *type_[R]) WriteType(w io.Writer) error {
	if _, err := io.WriteString(w, t.value.Name); err != nil {
		return err
	}
	return nil
}

// customization
func (b *ObjectType) WriteType(w io.Writer) error {
	if err := b.type_.WriteType(w); err != nil {
		return err
	}
	// if b.type_.value.IsNewType {
	// 	return nil
	// }

	io.WriteString(w, "{") // nolint
	n := len(b.Fields) - 1
	for i, f := range b.Fields {
		v := f.value
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

func (t *ArrayType[T]) WriteType(w io.Writer) error {
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

func (t *MapType[T]) WriteType(w io.Writer) error {
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
