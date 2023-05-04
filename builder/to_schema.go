package builder

import (
	"fmt"
	"log"

	"github.com/iancoleman/orderedmap"
	"github.com/podhmo/gos/builder/maplib"
)

type toSchemer interface {
	toSchema(*Builder) *orderedmap.OrderedMap
}

func ToSchema(b *Builder) (*orderedmap.OrderedMap, error) {
	doc := orderedmap.New()

	components := orderedmap.New()
	doc.Set("components", components)

	schemas := orderedmap.New()
	components.Set("schemas", schemas)

	if err := b.EachTypes(func(t TypeBuilder) error {
		name := t.typemetadata().Name
		if t, ok := t.(toSchemer); ok {
			schemas.Set(name, t.toSchema(b))
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
func (t *type_[R]) toSchema(b *Builder) *orderedmap.OrderedMap {
	doc := orderedmap.New()
	doc.Set("type", t.metadata.underlying)
	doc, err := maplib.Merge(doc, t.metadata)
	if err != nil {
		panic(err)
	}
	return doc
}

func (t *StringType) toSchema(b *Builder) *orderedmap.OrderedMap {
	doc := t.type_.toSchema(b)
	doc, err := maplib.Merge(doc, t.metadata)
	if err != nil {
		panic(err)
	}
	return doc
}
func (t *IntegerType) toSchema(b *Builder) *orderedmap.OrderedMap {
	doc := t.type_.toSchema(b)
	doc, err := maplib.Merge(doc, t.metadata)
	if err != nil {
		panic(err)
	}
	return doc
}
func (t *ArrayType[T]) toSchema(b *Builder) *orderedmap.OrderedMap {
	doc := t.type_.toSchema(b)
	doc.Set("items", t.items.toSchema(b))
	doc, err := maplib.Merge(doc, t.metadata)
	if err != nil {
		panic(err)
	}
	return doc
}
func (t *MapType[T]) toSchema(b *Builder) *orderedmap.OrderedMap {
	doc := t.type_.toSchema(b)
	doc.Set("type", "object")
	if t.metadata.PatternProperties == nil {
		doc.Set("additionalProperties", t.items.toSchema(b))
	} else {
		props := orderedmap.New()
		for k, typ := range t.metadata.PatternProperties {
			props.Set(k, typ.toSchema(b))
		}
		doc.Set("patternProperties", props)
	}

	doc, err := maplib.Merge(doc, t.metadata)
	if err != nil {
		panic(err)
	}
	return doc
}

func (t *ObjectType) toSchema(b *Builder) *orderedmap.OrderedMap {
	doc := t.type_.toSchema(b)
	required := make([]string, 0, len(t.Fields))
	if len(t.Fields) > 0 {
		properties := orderedmap.New()
		for _, f := range t.Fields {
			v := f.metadata
			name := v.Name
			if v.Required {
				required = append(required, name)
			}

			var def *orderedmap.OrderedMap
			if t, ok := f.typ.(toSchemer); ok {
				def = t.toSchema(b)
			} else {
				def = orderedmap.New()
			}
			def, err := maplib.Merge(def, f.metadata)
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
	doc, err := maplib.Merge(doc, t.metadata)
	if err != nil {
		panic(err)
	}
	return doc
}

func (t TypeRef) toSchema(b *Builder) *orderedmap.OrderedMap {
	doc := orderedmap.New()
	typ := t.getType()
	if typ == nil {
		log.Printf("#/components/schemas/%s is not found", t.Name)
		doc.Set("$ref", fmt.Sprintf("#/components/schemas/%s", t.Name))
	} else {
		doc.Set("$ref", fmt.Sprintf("#/components/schemas/%s", typ.typemetadata().Name))
	}

	return doc
}
