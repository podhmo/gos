package openapigen

import (
	"fmt"

	"github.com/podhmo/gos/enumgen"
)

func (b *Builder) StringFromEnum(enum *enumgen.String) *String {
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
		typ.Doc(docValues...)
	}
	return typ
}

func (b *Builder) IntFromEnum(enum *enumgen.Int) *Int {
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
		typ.Doc(docValues...)
	}
	return typ
}

func (t *Object) OnlyFields(names ...string) []*Field {
	fields := make([]*Field, 0, len(t.metadata.Fields))
	for _, f := range t.metadata.Fields {
		for _, name := range names {
			if f.metadata.Name == name {
				fields = append(fields, f)
				break
			}
		}
	}
	return fields
}

func (t *Object) IgnoreFields(names ...string) []*Field {
	fields := make([]*Field, 0, len(t.metadata.Fields))
	for _, f := range t.metadata.Fields {
		found := false
		for _, name := range names {
			if f.metadata.Name == name {
				found = true
				break
			}
		}
		if !found {
			fields = append(fields, f)
		}
	}
	return fields
}

func (b *Builder) Query(name string, typ TypeBuilder) *Param {
	return b.Param(name, typ, "query")
}
func (b *Builder) Path(name string, typ TypeBuilder) *Param {
	return b.Param(name, typ, "path").Required(true)
}
func (b *Builder) Header(name string, typ TypeBuilder) *Param {
	return b.Param(name, typ, "header")
}
func (b *Builder) Cookie(name string, typ TypeBuilder) *Param {
	return b.Param(name, typ, "cookie")
}
