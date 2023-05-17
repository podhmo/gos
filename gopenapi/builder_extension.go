package gopenapi

import (
	"fmt"
	"strings"

	"github.com/podhmo/gos/genum"
)

func (b *Builder) StringFromEnum(enum *genum.String) *String {
	typ := b.String()

	var docValues []string
	if doc := enum.GetEnumMetadata().Doc; doc != "" {
		docValues = append(docValues, doc, "")
	}

	metadata := enum.GetMetadata()
	values := make([]string, len(metadata.Members))
	for i, v := range metadata.Members {
		values[i] = v.Value
		if v.Doc != "" {
			docValues = append(docValues, fmt.Sprintf("* %s %s", v.Name, v.Doc))
		}
	}

	typ.Enum(values)
	if metadata.Default != "" {
		typ.Default(metadata.Default)
	}
	if len(docValues) > 0 {
		typ.Doc(strings.Join(docValues, "\n"))
	}
	return typ
}

func (b *Builder) IntFromEnum(enum *genum.Int) *Int {
	typ := b.Int()

	var docValues []string
	if doc := enum.GetEnumMetadata().Doc; doc != "" {
		docValues = append(docValues, doc)
	}

	metadata := enum.GetMetadata()
	values := make([]int64, len(metadata.Members))
	for i, v := range metadata.Members {
		values[i] = int64(v.Value)
		if v.Doc != "" {
			docValues = append(docValues, fmt.Sprintf("* %s %s", v.Name, v.Doc))
		}
	}
	typ.Enum(values)
	if metadata.Default != 0 {
		typ.Default(int64(metadata.Default))
	}
	if len(docValues) > 0 {
		typ.Doc(strings.Join(docValues, "\n"))
	}
	return typ
}
