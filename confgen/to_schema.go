package confgen

import (
	"fmt"
	"log"

	"github.com/iancoleman/orderedmap"
	"github.com/podhmo/gos/pkg/maplib"
)

type toSchemer interface {
	toSchema(b *Builder, useRef bool) *orderedmap.OrderedMap
}

func ToSchemaWith(doc *orderedmap.OrderedMap, b *Builder, t Type, useRef bool) (*orderedmap.OrderedMap, error) {
	name := t.GetTypeMetadata().Name
	doc.Set("title", name)
	return maplib.Merge(doc, t.toSchema(b, useRef))
}

func _toRefSchemaIfNamed[R TypeBuilder](b *Builder, t *_Type[R], useRef bool) (doc *orderedmap.OrderedMap, cached bool) {
	if !useRef {
		return nil, false
	}
	id := t.metadata.id
	if named := id > 0; !named { // if named by Define(), id > 0
		return nil, false
	}

	if ref, cached := b.Config.seen[id]; cached {
		return ref.toSchemaInternal(b), true
	}

	b.Config.defs = append(b.Config.defs, t.ret) // enqueue definitions
	ref := &TypeRef{Name: t.metadata.Name, rootbuilder: b, _Type: t.ret}
	b.Config.seen[id] = ref
	return ref.toSchemaInternal(b), false
}

// customization
func (t *_Type[R]) toSchema(b *Builder, useRef bool) *orderedmap.OrderedMap {
	doc := orderedmap.New()
	doc.Set("type", t.metadata.underlying)
	doc, err := maplib.Merge(doc, t.metadata)
	if err != nil {
		panic(err)
	}
	return doc
}

func (t *String) toSchema(b *Builder, useRef bool) *orderedmap.OrderedMap {
	if doc, _ := _toRefSchemaIfNamed(b, t._Type, useRef); doc != nil {
		return doc
	}

	doc := t._Type.toSchema(b, useRef)
	doc, err := maplib.Merge(doc, t.metadata)
	if err != nil {
		panic(err)
	}
	return doc
}
func (t *Int) toSchema(b *Builder, useRef bool) *orderedmap.OrderedMap {
	if doc, _ := _toRefSchemaIfNamed(b, t._Type, useRef); doc != nil {
		return doc
	}

	doc := t._Type.toSchema(b, useRef)
	doc, err := maplib.Merge(doc, t.metadata)
	if err != nil {
		panic(err)
	}
	return doc
}
func (t *Float) toSchema(b *Builder, useRef bool) *orderedmap.OrderedMap {
	if doc, _ := _toRefSchemaIfNamed(b, t._Type, useRef); doc != nil {
		return doc
	}

	doc := t._Type.toSchema(b, useRef)
	doc, err := maplib.Merge(doc, t.metadata)
	if err != nil {
		panic(err)
	}
	return doc
}
func (t *Array[T]) toSchema(b *Builder, useRef bool) *orderedmap.OrderedMap {
	doc, cached := _toRefSchemaIfNamed(b, t._Type, useRef)
	if doc != nil {
		if !cached {
			t.items.toSchema(b, true /* useRef */)
		}
		return doc
	}

	doc = t._Type.toSchema(b, useRef)
	doc.Set("items", t.items.toSchema(b, true)) // treating sub schema as always the ref.
	doc, err := maplib.Merge(doc, t.metadata)
	if err != nil {
		panic(err)
	}
	return doc
}
func (t *Map[T]) toSchema(b *Builder, useRef bool) *orderedmap.OrderedMap {
	doc, cached := _toRefSchemaIfNamed(b, t._Type, useRef)
	if doc != nil {
		if !cached {
			t.items.toSchema(b, true /* useRef */)
		}
		return doc
	}

	doc = t._Type.toSchema(b, useRef)
	doc.Set("type", "object")
	if t.metadata.Pattern == "" {
		doc.Set("additionalProperties", t.items.toSchema(b, true)) // treating sub schema as always the ref.
	} else {
		props := orderedmap.New()
		props.Set(t.metadata.Pattern, t.items.toSchema(b, true)) // treating sub schema as always the ref.
		doc.Set("patternProperties", props)
		doc.Set("additionalProperties", false)
	}

	doc, err := maplib.Merge(doc, t.metadata)
	if err != nil {
		panic(err)
	}
	return doc
}

func (t *Object) toSchema(b *Builder, useRef bool) *orderedmap.OrderedMap {
	doc, cached := _toRefSchemaIfNamed(b, t._Type, useRef)
	if doc != nil {
		if !cached {
			for _, v := range t.metadata.Fields {
				v.metadata.Typ.toSchema(b, true /* true */)
			}
		}
		return doc
	}

	doc = t._Type.toSchema(b, useRef)
	required := make([]string, 0, len(t.metadata.Fields))
	if len(t.metadata.Fields) > 0 {
		useRef := true // treating sub schema as always the ref.
		properties := orderedmap.New()
		for _, v := range t.metadata.Fields {
			name := v.metadata.Name
			if v.metadata.Required {
				required = append(required, name)
			}

			def := v.metadata.Typ.toSchema(b, useRef)
			def, err := maplib.Merge(def, v.metadata)
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

func (t TypeRef) toSchema(b *Builder, useRef bool) *orderedmap.OrderedMap {
	return t.getType().toSchema(b, useRef)
}
func (t TypeRef) toSchemaInternal(b *Builder) *orderedmap.OrderedMap {
	doc := orderedmap.New()
	typ := t.getType()

	// in draft-2020-12 $defs
	if typ == nil {
		log.Printf("#/definitions/%s is not found", t.Name)
		doc.Set("$ref", fmt.Sprintf("#/definitions/%s", t.Name))
	} else {
		doc.Set("$ref", fmt.Sprintf("#/definitions/%s", typ.GetTypeMetadata().Name))
	}
	return doc
}

func (t *_Container) toSchema(b *Builder, useRef bool) *orderedmap.OrderedMap {
	doc, cached := _toRefSchemaIfNamed(b, t._Type, useRef)
	if doc != nil {
		if !cached {
			for _, typ := range t.metadata.Types {
				typ.toSchema(b, true /* useRef */)
			}
		}
		return doc
	}

	doc = orderedmap.New()
	types := make([]*orderedmap.OrderedMap, len(t.metadata.Types))
	for i, typ := range t.metadata.Types {
		types[i] = typ.toSchema(b, useRef)
	}
	doc.Set(t.metadata.Op, types)
	if discriminator := t.metadata.Discriminator; discriminator != "" {
		v := orderedmap.New()
		v.Set("propertyName", discriminator)
		doc.Set("discriminator", v)
	}
	doc, err := maplib.Merge(doc, t.metadata)
	if err != nil {
		panic(err)
	}
	return doc
}
