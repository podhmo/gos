package openapigen

import (
	"fmt"
	"log"
	"strconv"

	"github.com/iancoleman/orderedmap"
	"github.com/podhmo/gos/pkg/maplib"
)

type toSchemer interface {
	toSchema(*Builder) *orderedmap.OrderedMap
}

func ToSchemaWith(b *Builder, doc *orderedmap.OrderedMap) (*orderedmap.OrderedMap, error) {
	components := orderedmap.New()
	doc.Set("components", components)

	schemas := orderedmap.New()
	components.Set("schemas", schemas)

	if err := b.EachTypes(func(t TypeBuilder) error {
		name := t.GetTypeMetadata().Name
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
func ToSchema(b *Builder) (*orderedmap.OrderedMap, error) {
	doc := orderedmap.New()
	return ToSchemaWith(b, doc)
}

// customization
func (t *_Type[R]) toSchema(b *Builder) *orderedmap.OrderedMap {
	doc := orderedmap.New()
	doc.Set("type", t.metadata.underlying)
	doc, err := maplib.Merge(doc, t.metadata)
	if err != nil {
		panic(err)
	}
	return doc
}

func (t *String) toSchema(b *Builder) *orderedmap.OrderedMap {
	doc := t._Type.toSchema(b)
	doc, err := maplib.Merge(doc, t.metadata)
	if err != nil {
		panic(err)
	}
	return doc
}
func (t *Int) toSchema(b *Builder) *orderedmap.OrderedMap {
	doc := t._Type.toSchema(b)
	doc, err := maplib.Merge(doc, t.metadata)
	if err != nil {
		panic(err)
	}
	return doc
}
func (t *Array[T]) toSchema(b *Builder) *orderedmap.OrderedMap {
	doc := t._Type.toSchema(b)
	doc.Set("items", t.items.toSchema(b))
	doc, err := maplib.Merge(doc, t.metadata)
	if err != nil {
		panic(err)
	}
	return doc
}
func (t *Map[T]) toSchema(b *Builder) *orderedmap.OrderedMap {
	doc := t._Type.toSchema(b)
	doc.Set("type", "object")
	if t.metadata.Pattern == "" {
		doc.Set("additionalProperties", t.items.toSchema(b))
	} else {
		props := orderedmap.New()
		props.Set(t.metadata.Pattern, t.items.toSchema(b))
		doc.Set("patternProperties", props)
		doc.Set("additionalProperties", false)
	}

	doc, err := maplib.Merge(doc, t.metadata)
	if err != nil {
		panic(err)
	}
	return doc
}

func (t *Object) toSchema(b *Builder) *orderedmap.OrderedMap {
	doc := t._Type.toSchema(b)
	required := make([]string, 0, len(t.metadata.Fields))
	if len(t.metadata.Fields) > 0 {
		properties := orderedmap.New()
		for _, v := range t.metadata.Fields {
			name := v.metadata.Name
			if v.metadata.Required {
				required = append(required, name)
			}

			def := v.metadata.Typ.toSchema(b)
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

func (t TypeRef) toSchema(b *Builder) *orderedmap.OrderedMap {
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

func (t *Action) toSchema(b *Builder) *orderedmap.OrderedMap {
	doc := orderedmap.New()
	doc.Set("operationId", t.metadata.Name)

	if t.metadata.Input != nil {
		if params := t.metadata.Input.metadata.Params; len(params) > 0 {
			parameters := make([]*orderedmap.OrderedMap, len(params))
			for i, p := range params {
				doc := orderedmap.New()
				doc.Set("name", p.metadata.Name)
				doc.Set("in", p.metadata.In)
				if p.metadata.Description != "" {
					doc.Set("description", p.metadata.Description)
				}
				doc.Set("required", p.metadata.Required)
				doc.Set("schema", p.metadata.Typ.toSchema(b))
				parameters[i] = doc
			}
		}
		if body := t.metadata.Input.metadata.Body; body != nil {
			requestBody := orderedmap.New()
			requestBody.Set("required", true)
			doc.Set("requestBody", requestBody)
			content := orderedmap.New()
			requestBody.Set("content", content)
			appjson := orderedmap.New()
			content.Set("application/json", appjson)
			appjson.Set("schema", body.metadata.Typ.toSchema(b))
		}
	}

	responses := orderedmap.New()
	responses.Set("responses", responses)
	res := orderedmap.New()
	doc.Set(strconv.Itoa(t.metadata.DefaultStatus), res)
	res.Set("description", t._Type.metadata.Doc)
	content := orderedmap.New()
	res.Set("content", content)
	appjson := orderedmap.New()
	content.Set("applicatioin/json", appjson)
	if t.metadata.Output != nil {
		appjson.Set("schema", t.metadata.Output.metadata.Typ.toSchema(b))
	}
	return doc
}
