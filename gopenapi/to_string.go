// Generated by github.com/podhmo/gos/gopenapi/tools [-write -builder -metadata -to-string -pkgname gopenapi]

package gopenapi

import (
	"fmt"
	"io"
	"strings"
)

func (t *_Type[R]) String() string {
	return ToString(t.ret)
}

func ToString(t TypeBuilder) string {
	b := new(strings.Builder)
	if err := t.writeType(b); err != nil {
		return fmt.Sprintf("invalid type: %T", t)
	}
	return b.String()
}

// default implementation or write Type
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
