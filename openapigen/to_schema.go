package openapigen

import (
	"fmt"
	"log"
	"strconv"

	"github.com/iancoleman/orderedmap"
	"github.com/podhmo/gos/pkg/maplib"
)

type toSchemer interface {
	toSchema(b *Builder, useRef bool) *orderedmap.OrderedMap
}

func ToSchemaWith(doc *orderedmap.OrderedMap, b *Builder, useRef bool) (*orderedmap.OrderedMap, error) {
	components := orderedmap.New()
	doc.Set("components", components)

	schemas := orderedmap.New()
	components.Set("schemas", schemas)

	if err := b.EachTypes(func(t TypeBuilder) error {
		name := t.GetTypeMetadata().Name
		if t, ok := t.(toSchemer); ok {
			schemas.Set(name, t.toSchema(b, useRef))
		} else {
			schemas.Set(name, orderedmap.New())
		}
		return nil
	}); err != nil {
		return doc, fmt.Errorf("each types -- %w", err)
	}
	return doc, nil
}
func ToSchema(b *Builder) (*orderedmap.OrderedMap, error) {
	doc := orderedmap.New()
	useRef := false
	return ToSchemaWith(doc, b, useRef)
}

func _toRefSchemaIfNamed[R TypeBuilder](b *Builder, t *_Type[R], useRef bool) (doc *orderedmap.OrderedMap, cached bool) {
	if !useRef {
		return nil, false
	}
	id := t.metadata.id
	if named := id > 0; !named { // if named by DefineType(), id > 0
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
	doc.Set("items", t.items.toSchema(b, useRef))
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
		doc.Set("additionalProperties", t.items.toSchema(b, useRef))
	} else {
		props := orderedmap.New()
		props.Set(t.metadata.Pattern, t.items.toSchema(b, useRef))
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
	if typ == nil {
		log.Printf("#/components/schemas/%s is not found", t.Name)
		doc.Set("$ref", fmt.Sprintf("#/components/schemas/%s", t.Name))
	} else {
		doc.Set("$ref", fmt.Sprintf("#/components/schemas/%s", typ.GetTypeMetadata().Name))
	}

	return doc
}

func (t *Action) toSchema(b *Builder, useRef bool) *orderedmap.OrderedMap {
	doc := orderedmap.New()
	doc.Set("operationId", t.metadata.Name)
	doc.Set("description", t._Type.metadata.Doc)

	if input := t.metadata.Input; input != nil {
		if params := input.metadata.Params; len(params) > 0 {
			parameters := make([]*orderedmap.OrderedMap, len(params))
			for i, p := range params {
				doc := orderedmap.New()
				doc, err := maplib.Merge(doc, p.metadata)
				if err != nil {
					panic(err)
				}
				doc.Set("schema", p.metadata.Typ.toSchema(b, useRef))
				parameters[i] = doc
			}
			doc.Set("parameters", parameters)
		}
		if body := input.metadata.Body; body != nil {
			requestBody := orderedmap.New()
			requestBody.Set("required", true)
			doc.Set("requestBody", requestBody)
			content := orderedmap.New()
			requestBody.Set("content", content)
			appjson := orderedmap.New()
			content.Set("application/json", appjson)
			appjson.Set("schema", body.metadata.Typ.toSchema(b, useRef))
			
			description := body.GetTypeMetadata().Doc
			if description == "" {
				description = input.GetTypeMetadata().Doc
			}
			if description != "" {
				requestBody.Set("description", description)
			}
		}
	}

	responses := orderedmap.New()
	doc.Set("responses", responses)

	for _, output := range t.metadata.Outputs {
		if output.metadata.IsDefault {
			k := "default"
			typ := t.metadata.DefaultError
			description := output.GetTypeMetadata().Doc
			if description == "" {
				description = "default error"
			}
			res := orderedmap.New()
			responses.Set(k, res)
			if typ != nil {
				res.Set("description", description)
			}
			content := orderedmap.New()
			res.Set("content", content)
			appjson := orderedmap.New()
			content.Set("application/json", appjson)
			if typ != nil {
				appjson.Set("schema", output.metadata.Typ.toSchema(b, useRef))
			}
		} else {
			k := strconv.Itoa(output.metadata.Status)
			typ := output.metadata.Typ
			description := output.GetTypeMetadata().Doc
			if typ != nil && description == "" {
				description = typ.GetTypeMetadata().Doc
			}
			res := orderedmap.New()
			responses.Set(k, res)
			res.Set("description", description)
			content := orderedmap.New()
			res.Set("content", content)
			appjson := orderedmap.New()
			content.Set("application/json", appjson)
			if typ != nil {
				appjson.Set("schema", output.metadata.Typ.toSchema(b, useRef))
			}
		}
	}
	doc, err := maplib.Merge(doc, t.metadata)
	if err != nil {
		panic(err)
	}
	return doc
}
