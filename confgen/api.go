//go:generate go run ./tools -write -builder -metadata -stringer -pkgname confgen
//go:generate go fmt .
package confgen

import (
	"github.com/iancoleman/orderedmap"
)

type Config struct {

	// for to schema
	defs []TypeBuilder
	seen map[int]*TypeRef
}

func DefaultConfig() *Config {
	c := &Config{
		seen: map[int]*TypeRef{},
	}
	return c
}

func ToJSONSchema(b *Builder, typ Type) (*orderedmap.OrderedMap, error) {
	doc := orderedmap.New()
	useRef := false
	return ToSchemaWith(doc, b, useRef)
}
