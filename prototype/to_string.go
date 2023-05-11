package prototype

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

func (t *type_[R]) writeType(w io.Writer) error {
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
func (b *ObjectType) writeType(w io.Writer) error {
	if err := b.type_.writeType(w); err != nil {
		return err
	}
	// if b.type_.metadata.IsNewType {
	// 	return nil
	// }

	io.WriteString(w, "{") // nolint
	n := len(b.Fields) - 1
	for i, f := range b.Fields {
		v := f.metadata
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

func (t *ArrayType[T]) writeType(w io.Writer) error {
	if err := t.type_.writeType(w); err != nil {
		return err
	}
	// if t.type_.metadata.IsNewType {
	// 	return nil
	// }

	io.WriteString(w, "[") // nolint
	if err := t.items.writeType(w); err != nil {
		return err
	}
	io.WriteString(w, "]") // nolint
	return nil
}

func (t *MapType[T]) writeType(w io.Writer) error {
	if err := t.type_.writeType(w); err != nil {
		return err
	}
	// if t.type_.metadata.IsNewType {
	// 	return nil
	// }

	io.WriteString(w, "[") // nolint
	if err := t.items.writeType(w); err != nil {
		return err
	}
	io.WriteString(w, "]") // nolint
	return nil
}

func (t *ActionType) writeType(w io.Writer) error {
	if err := t.type_.writeType(w); err != nil {
		return err
	}
	io.WriteString(w, " ")
	if err := t.input.writeType(w); err != nil {
		return err
	}
	if err := t.output.writeType(w); err != nil {
		return err
	}
	return nil
}

func (t *ActionInput) writeType(w io.Writer) error {
	io.WriteString(w, "(") // nolint
	if t != nil {
		for i, p := range t.Params {
			if err := p.typ.writeType(w); err != nil {
				return err
			}
			if i < len(t.Params)-1 {
				io.WriteString(w, ", ")
			}
		}
	}
	io.WriteString(w, ")") // nolint
	return nil
}

func (t *ActionOutput) writeType(w io.Writer) error {
	if t != nil {
		if err := t.retval.typ.writeType(w); err != nil {
			return err
		}
	}
	return nil
}
