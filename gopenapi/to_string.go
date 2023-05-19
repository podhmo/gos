package gopenapi

import (
	"fmt"
	"io"
	"strings"
)

type writeTyper interface {
	writeType(io.Writer) error
}

func ToString(typ TypeBuilder) string {
	b := new(strings.Builder)
	if err := typ.writeType(b); err != nil {
		return fmt.Sprintf("invalid type: %T", typ)
	}
	return b.String()
}

func (t *TypeRef) writeType(w io.Writer) error {
	return t.getType().writeType(w)
}

func (t *_Type[R]) writeType(w io.Writer) error {
	if t.metadata.Name != "" {
		if _, err := io.WriteString(w, t.metadata.Name); err != nil {
			return err
		}
	} else {
		if _, err := io.WriteString(w, t.metadata.underlying); err != nil {
			return err
		}
	}
	return nil
}

// customization
func (t *Object) writeType(w io.Writer) error {
	if err := t._Type.writeType(w); err != nil {
		return err
	}
	// if t._Type.metadata.Name != "" {
	// 	return nil
	// }

	io.WriteString(w, "{") // nolint
	n := len(t.metadata.Fields) - 1
	for i, f := range t.metadata.Fields {
		io.WriteString(w, f.metadata.Name) // nolint
		if !f.metadata.Required {
			io.WriteString(w, "?") // nolint
		}
		if i < n {
			io.WriteString(w, ", ") // nolint
		}
	}
	io.WriteString(w, "}") // nolint
	return nil
}

func (t *Array[T]) writeType(w io.Writer) error {
	if err := t._Type.writeType(w); err != nil {
		return err
	}
	// if t._Type.metadata.Name == "" {
	// 	return nil
	// }

	io.WriteString(w, "[") // nolint
	if err := t.items.writeType(w); err != nil {
		return err
	}
	io.WriteString(w, "]") // nolint
	return nil
}

func (t *Map[T]) writeType(w io.Writer) error {
	if err := t._Type.writeType(w); err != nil {
		return err
	}
	// if t._Type.metadata.Name == "" {
	// 	return nil
	// }

	io.WriteString(w, "[") // nolint
	if err := t.items.writeType(w); err != nil {
		return err
	}
	io.WriteString(w, "]") // nolint
	return nil
}

func (t *Action) writeType(w io.Writer) error {
	io.WriteString(w, t.metadata.Name)
	io.WriteString(w, " :: ")
	if err := t.metadata.Input.writeType(w); err != nil {
		return err
	}
	if err := t.metadata.Output.writeType(w); err != nil {
		return err
	}
	return nil
}

func (t *Input) writeType(w io.Writer) error {
	io.WriteString(w, "(") // nolint
	if t != nil {
		for i, p := range t.metadata.Params {
			if err := p.metadata.Typ.writeType(w); err != nil {
				return err
			}
			if i < len(t.metadata.Params)-1 {
				io.WriteString(w, ", ")
			}
		}
	}
	io.WriteString(w, ")") // nolint
	return nil
}

func (t *Output) writeType(w io.Writer) error {
	if t != nil {
		io.WriteString(w, " => ") // nolint
		if err := t.metadata.Typ.writeType(w); err != nil {
			return err
		}
	}
	return nil
}
