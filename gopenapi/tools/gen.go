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
	b := seed.NewBuilder(options.PkgName)
	b.GeneratedBy("github.com/podhmo/gos/gopenapi/tools")
	b.NeedReference()

	b.Import("strings")

	Type := b.BuildTarget("Type",
		b.Field("Format", seed.Symbol("string")).Tag(`json:"format,omitempty"`),
		b.Field("Doc", seed.Symbol("string")).Tag(`json:"description,omitempty"`),
	).Setter("Doc", b.Arg("stmts", seed.Symbol("string")).Variadic().Transform(func(stmts string) string {
		return fmt.Sprintf(`strings.Join(%s, "\n")`, stmts)
	}))

	b.InterfaceMethods(
		"toSchemer // see: ./to_schema.go",
	)

	// ----------------------------------------
	// types
	// ----------------------------------------
	Bool := b.Type("Bool").NeedBuilder().Underlying("boolean")
	Int := b.Type("Int",
		b.Field("Enum", seed.Symbol("[]int64")).Tag(`json:"enum,omitempty"`),
		b.Field("Default", seed.Symbol("int64")).Tag(`json:"default,omitempty"`),
		b.Field("Maximum", seed.Symbol("int64")).Tag(`json:"maximum,omitempty"`),
		b.Field("Minimum", seed.Symbol("int64")).Tag(`json:"minimum,omitempty"`),
	).NeedBuilder().Underlying("integer")
	String := b.Type("String",
		b.Field("Enum", seed.Symbol("[]string")).Tag(`json:"enum,omitempty"`),
		b.Field("Default", seed.Symbol("string")).Tag(`json:"default,omitempty"`),
		b.Field("Pattern", seed.Symbol("string")).Tag(`json:"pattern,omitempty"`),
		b.Field("MaxLength", seed.Symbol("int64")).Tag(`json:"maxlength,omitempty"`),
		b.Field("MinLength", seed.Symbol("int64")).Tag(`json:"minlength,omitempty"`),
	).NeedBuilder().Underlying("string")

	Array := b.Type("Array", b.TypeVar("Items", seed.Symbol("TypeBuilder")),
		b.Field("MaxItems", seed.Symbol("int64")).Tag(`json:"maxitems,omitempty"`),
		b.Field("MinItems", seed.Symbol("int64")).Tag(`json:"minitems,omitempty"`),
	).NeedBuilder().Underlying("array")
	Map := b.Type("Map", b.TypeVar("Items", seed.Symbol("TypeBuilder")),
		b.Field("Pattern", seed.Symbol("string")).Tag(`json:"pattern,omitempty"`),
	).NeedBuilder().Underlying("map")

	Field := b.Type("Field",
		b.Field("Name", seed.Symbol("string")).Tag(`json:"-"`),
		b.Field("Typ", seed.Symbol("TypeBuilder")).Tag(`json:"-"`),
		b.Field("Description", seed.Symbol("string")).Tag(`json:"description,omitempty"`),
		b.Field("Required", seed.Symbol("bool")).Tag(`json:"-"`).Default("true"),
	).Constructor(
		b.Arg("Name", seed.Symbol("string")),
		b.Arg("Typ", seed.Symbol("TypeBuilder")),
	).NeedBuilder().Underlying("field") //?

	Object := b.Type("Object",
		b.Field("Fields", seed.Symbol("[]*Field")).Tag(`json:"-"`),
		b.Field("Strict", seed.Symbol("bool")).Tag(`json:"-"`).Default("true"),
	).Constructor(
		b.Arg("Fields", seed.Symbol("*Field")).Variadic(),
	).NeedBuilder().Underlying("object")

	// ----------------------------------------
	// action
	// ----------------------------------------
	Action := b.Type("Action",
		b.Field("Name", seed.Symbol("string")),
		b.Field("Input", "*Input"),
		b.Field("Output", "*Output"),
		b.Field("DefaultStatus", seed.Symbol("int")).Default("200"),
		b.Field("Method", seed.Symbol("string")),
		b.Field("Path", seed.Symbol("string")),
	).Constructor(
		b.Arg("Name", seed.Symbol("string")),
		b.Arg("Input", "*Input"),
		b.Arg("Output", "*Output"),
	).NeedBuilder().Underlying("action")

	Input := b.Type("Input",
		b.Field("Params", "[]*Param"),
	).Constructor(
		b.Arg("Params", "*Param").Variadic(),
	).NeedBuilder().Underlying("input")
	Output := b.Type("Output",
		b.Field("Typ", "TypeBuilder"),
	).Constructor(
		b.Arg("Typ", "TypeBuilder"),
	).NeedBuilder().Underlying("output")

	Param := b.Type("Param",
		b.Field("Name", seed.Symbol("string")).Tag(`json:"-"`),
		b.Field("In", seed.Symbol("string")).Tag(`json:"in"`),
		b.Field("Typ", seed.Symbol("TypeBuilder")).Tag(`json:"-"`),
		b.Field("Description", seed.Symbol("string")).Tag(`json:"description,omitempty"`),
		b.Field("Required", seed.Symbol("bool")).Tag(`json:"required"`).Default("true"),
	).Constructor(
		b.Arg("Name", seed.Symbol("string")),
		b.Arg("Typ", seed.Symbol("TypeBuilder")),
		b.Arg("In", seed.Symbol("string")), // query,header,path,cookie,body
	).NeedBuilder().Underlying("param")

	fmt.Fprintln(os.Stderr, Type, Bool, Int, String, Array, Map, Field, Object)
	fmt.Fprintln(os.Stderr, Action, Input, Output, Param)

	// for transform
	b.Footer(`
	// toSlice is list.map as you know.
	func toSlice[S, D any](src []S, conv func(S) D) []D {
		dst := make([]D, len(src))
		for i, x := range src {
			dst[i] = conv(x)
		}
		return dst
	}	
	`)

	// emit
	return cmd.Do(b)
}
