package main

import (
	"fmt"
	"log"
	"os"

	"github.com/podhmo/gos/seed"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("!! %+v", err)
	}
}

func run() error {
	cmd := seed.NewCommand(os.Args[1:])
	options := cmd.Config

	// define
	b := seed.NewBuilder(
		options.PkgName,
	)

	b.GeneratedBy("github.com/podhmo/gos/graphqlgen/tools")
	b.NeedReference()

	b.Import("strings")
	b.Import("github.com/iancoleman/orderedmap")
	b.ImportInMetadata("github.com/iancoleman/orderedmap") // use goimports?

	Type := b.BuildTarget("Type",
		b.Field("Doc", seed.Symbol("string")).Tag(`json:"description,omitempty"`),
		b.Field("Title", seed.Symbol("string")).Tag(`json:"title,omitempty"`),
		b.Field("Example", seed.Symbol("string")).Tag(`json:"example,omitempty"`),
		b.Field("Extensions", seed.Symbol("*orderedmap.OrderedMap")).Tag(`json:"-"`),
	).Setter("Doc", b.Arg("stmts", seed.Symbol("string")).Variadic().Transform(func(stmts string) string {
		return fmt.Sprintf(`strings.Join(%s, "\n")`, stmts)
	})).Setter("Extensions", b.Arg("extensions", seed.Symbol("*Extension")).Variadic().Transform(func(extensions string) string {
		return fmt.Sprintf(`func() *orderedmap.OrderedMap {
		if %s == nil {
			return nil
		}
		data := orderedmap.New()
		for _, ext := range %s {
			m := ext.metadata
			name := m.Name
			if !strings.HasPrefix(name, "x-") {
				name = "x-" + name
			}
			data.Set(name, m.Value)
		}
		return data
		}()`, extensions, extensions)
	}))

	b.InterfaceMethods(
	//		"toSchemer // see: ./to_schema.go",
	)

	// ----------------------------------------
	// types
	// ----------------------------------------
	//

	Bool := b.Type("Bool",
		b.Field("Format", seed.Symbol("string")).Tag(`json:"format,omitempty"`),
		b.Field("Default", seed.Symbol("bool")).Tag(`json:"default,omitempty"`),
	).NeedBuilder().Underlying("boolean").GoType("bool")
	Int := b.Type("Int",
		b.Field("Format", seed.Symbol("string")).Tag(`json:"format,omitempty"`),
		b.Field("Enum", seed.Symbol("[]int64")).Tag(`json:"enum,omitempty"`),
		b.Field("Default", seed.Symbol("int64")).Tag(`json:"default,omitempty"`),
		b.Field("Maximum", seed.Symbol("int64")).Tag(`json:"maximum,omitempty"`),
		b.Field("Minimum", seed.Symbol("int64")).Tag(`json:"minimum,omitempty"`),
		b.Field("ExclusiveMin", seed.Symbol("bool")).Tag(`json:"exclusiveMin,omitempty"`),
		b.Field("ExclusiveMax", seed.Symbol("bool")).Tag(`json:"exclusiveMax,omitempty"`),
	).NeedBuilder().Underlying("integer").GoType("int64")
	Float := b.Type("Float",
		b.Field("Format", seed.Symbol("string")).Tag(`json:"format,omitempty"`),
		b.Field("Default", seed.Symbol("string")).Tag(`json:"default,omitempty"`),
		b.Field("Maximum", seed.Symbol("float64")).Tag(`json:"maximum,omitempty"`),
		b.Field("Minimum", seed.Symbol("float64")).Tag(`json:"minimum,omitempty"`),
		b.Field("MultipleOf", seed.Symbol("float64")).Tag(`json:"multipleOf,omitempty"`),
		b.Field("ExclusiveMin", seed.Symbol("bool")).Tag(`json:"exclusiveMin,omitempty"`),
		b.Field("ExclusiveMax", seed.Symbol("bool")).Tag(`json:"exclusiveMax,omitempty"`),
	).NeedBuilder().Underlying("number").GoType("float64")
	String := b.Type("String",
		b.Field("Format", seed.Symbol("string")).Tag(`json:"format,omitempty"`),
		b.Field("Enum", seed.Symbol("[]string")).Tag(`json:"enum,omitempty"`),
		b.Field("Default", seed.Symbol("string")).Tag(`json:"default,omitempty"`),
		b.Field("Pattern", seed.Symbol("string")).Tag(`json:"pattern,omitempty"`),
		b.Field("MaxLength", seed.Symbol("int64")).Tag(`json:"maxLength,omitempty"`),
		b.Field("MinLength", seed.Symbol("int64")).Tag(`json:"minLength,omitempty"`),
	).NeedBuilder().Underlying("string").GoType("string")
	Array := b.Type("Array", b.TypeVar("Items", seed.Symbol("Type")),
		// b.Field("UniqueItems", seed.Symbol("bool")).Tag(`json:"uniqueItems,omitempty"`),
		b.Field("MaxItems", seed.Symbol("int64")).Tag(`json:"maxItems,omitempty"`),
		b.Field("MinItems", seed.Symbol("int64")).Tag(`json:"minItems,omitempty"`),
	).NeedBuilder().Underlying("array")
	Map := b.Type("Map", b.TypeVar("Items", seed.Symbol("Type")),
		b.Field("Pattern", seed.Symbol("string")).Tag(`json:"pattern,omitempty"`),
	).NeedBuilder().Underlying("map")

	Field := b.Type("Field",
		b.Field("Name", seed.Symbol("string")).Tag(`json:"-"`),
		b.Field("Typ", seed.Symbol("Type")).Tag(`json:"-"`),
		b.Field("Doc", seed.Symbol("string")).Tag(`json:"description,omitempty"`),
		b.Field("Nullable", seed.Symbol("bool")).Tag(`json:"nullable,omitempty"`), // trim omitempty?
		b.Field("Required", seed.Symbol("bool")).Tag(`json:"-"`).Default("true"),
		b.Field("ReadOnly", seed.Symbol("bool")).Tag(`json:"readonly,omitempty"`),
		b.Field("WriteOnly", seed.Symbol("bool")).Tag(`json:"writeonly,omitempty"`),
		b.Field("AllowEmptyValue", seed.Symbol("bool")).Tag(`json:"allowEmptyValue,omitempty"`),
		b.Field("Deprecated", seed.Symbol("bool")).Tag(`json:"deprecated,omitempty"`),
	).Constructor(
		b.Arg("Name", seed.Symbol("string")),
		b.Arg("Typ", seed.Symbol("Type")),
	).Setter("Doc", b.Arg("stmts", seed.Symbol("string")).Variadic().Transform(func(stmts string) string {
		return fmt.Sprintf(`strings.Join(%s, "\n")`, stmts)
	})).NeedBuilder().Underlying("field") //?

	b.Type("Extension",
		b.Field("Name", seed.Symbol("string")),
		b.Field("Value", seed.Symbol("any")),
	).Constructor(
		b.Arg("Name", seed.Symbol("string")),
		b.Arg("Value", seed.Symbol("any")),
	).NeedBuilder().Underlying("extension").Doc("for x-<extension-name>")

	Object := b.Type("Object",
		b.Field("Fields", seed.Symbol("[]*Field")).Tag(`json:"-"`),
		b.Field("MaxProperties", seed.Symbol("uint64")).Tag(`json:"maxProeprties,omitempty"`),
		b.Field("MinProperties", seed.Symbol("uint64")).Tag(`json:"minProeprties,omitempty"`),
		b.Field("Strict", seed.Symbol("bool")).Tag(`json:"-"`).Default("true"),
	).Constructor(
		b.Arg("Fields", seed.Symbol("*Field")).Variadic(),
	).NeedBuilder().Underlying("object")

	b.Union("Type", Bool, Int, Float, String, Object, Array, Map).
		DistinguishID("typ").
		InterfaceMethods("TypeBuilder").
		Doc("Type is union of bool | int | float | string | object | array[T] | map[T]").
		NeedReference()

	fmt.Fprintln(os.Stderr, Type, Bool, Int, Float, String, Array, Map, Field, Object)

	// for transform
	// b.Footer(``)

	// emit
	return cmd.Do(b)
}
