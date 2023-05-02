package builder

import (
	"fmt"
	"log"

	"github.com/iancoleman/orderedmap"
)

type toSchema interface {
	ToSchema(*Builder) *orderedmap.OrderedMap
}

func ToSchema(b *Builder) (*orderedmap.OrderedMap, error) {
	doc := orderedmap.New()

	components := orderedmap.New()
	doc.Set("components", components)

	schemas := orderedmap.New()
	components.Set("schemas", schemas)

	if err := b.EachTypes(func(t TypeBuilder) error {
		name := t.typevalue().Name
		if t, ok := t.(toSchema); ok {
			schemas.Set(name, t.ToSchema(b))
		} else {
			schemas.Set(name, orderedmap.New())
		}
		return nil
	}); err != nil {
		return doc, fmt.Errorf("each types -- %w", err)
	}
	return doc, nil
}

// customization
func (t *type_[R]) ToSchema(b *Builder) *orderedmap.OrderedMap {
	doc := orderedmap.New()
	doc.Set("type", t.value.underlying)
	return doc
}

func (t *ObjectBuilder[R]) ToSchema(b *Builder) *orderedmap.OrderedMap {
	doc := orderedmap.New()
	doc.Set("type", "object")
	required := make([]string, 0, len(t.Fields))
	if len(t.Fields) > 0 {
		properties := orderedmap.New()
		for _, f := range t.Fields {
			v := f.value
			name := v.Name
			if v.Required {
				required = append(required, name)
			}
			if t, ok := f.typ.(toSchema); ok {
				properties.Set(name, t.ToSchema(b))
			} else {
				properties.Set(name, orderedmap.New())
			}
		}
		doc.Set("properties", properties)
	}

	if len(required) > 0 {
		doc.Set("required", required)
	}
	doc.Set("additionalProperties", false)
	// TODO: merge object config
	return doc
}

func (t TypeRef) ToSchema(b *Builder) *orderedmap.OrderedMap {
	doc := orderedmap.New()
	typ := t.getType()
	if typ == nil {
		log.Printf("#/components/schemas/%s is not found", t.Name)
		doc.Set("$ref", fmt.Sprintf("#/components/schemas/%s", t.Name))
	} else {
		doc.Set("$ref", fmt.Sprintf("#/components/schemas/%s", typ.typevalue().Name))
	}

	return doc
}
