package builder

import (
	"fmt"
	"log"

	"github.com/iancoleman/orderedmap"
	"github.com/podhmo/gos/builder/maplib"
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
	doc, err := maplib.Merge(doc, t.value)
	if err != nil {
		panic(err)
	}
	return doc
}

func (t *StringBuilder[R]) ToSchema(b *Builder) *orderedmap.OrderedMap {
	doc := t.type_.ToSchema(b)
	doc, err := maplib.Merge(doc, t.value)
	if err != nil {
		panic(err)
	}
	return doc
}
func (t *IntegerBuilder[R]) ToSchema(b *Builder) *orderedmap.OrderedMap {
	doc := t.type_.ToSchema(b)
	doc, err := maplib.Merge(doc, t.value)
	if err != nil {
		panic(err)
	}
	return doc
}
func (t *ArrayBuilder[T, R]) ToSchema(b *Builder) *orderedmap.OrderedMap {
	doc := t.type_.ToSchema(b)
	doc.Set("items", t.items.ToSchema(b))
	doc, err := maplib.Merge(doc, t.value)
	if err != nil {
		panic(err)
	}
	return doc
}
func (t *MapBuilder[V, R]) ToSchema(b *Builder) *orderedmap.OrderedMap {
	doc := t.type_.ToSchema(b)
	doc.Set("type", "object")
	if t.value.PatternProperties == nil {
		doc.Set("additionalProperties", t.items.ToSchema(b))
	} else {
		props := orderedmap.New()
		for k, typ := range t.value.PatternProperties {
			props.Set(k, typ.ToSchema(b))
		}
		doc.Set("patternProperties", props)
	}

	doc, err := maplib.Merge(doc, t.value)
	if err != nil {
		panic(err)
	}
	return doc
}

func (t *ObjectBuilder[R]) ToSchema(b *Builder) *orderedmap.OrderedMap {
	doc := t.type_.ToSchema(b)
	required := make([]string, 0, len(t.Fields))
	if len(t.Fields) > 0 {
		properties := orderedmap.New()
		for _, f := range t.Fields {
			v := f.value
			name := v.Name
			if v.Required {
				required = append(required, name)
			}

			var def *orderedmap.OrderedMap
			if t, ok := f.typ.(toSchema); ok {
				def = t.ToSchema(b)
			} else {
				def = orderedmap.New()
			}
			def, err := maplib.Merge(def, f.value)
			if err != nil {
				panic(err)
			}
			properties.Set(name, def)
		}
		doc.Set("properties", properties)
	}

	if len(required) > 0 {
		doc.Set("required", required)
	}
	doc.Set("additionalProperties", false)
	doc, err := maplib.Merge(doc, t.value)
	if err != nil {
		panic(err)
	}
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