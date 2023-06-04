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
		seed.Root.Field("Config", seed.Symbol("*Config")),
	)
	b.Constructor(
		b.Arg("Config", seed.Symbol("*Config")),
	)

	b.GeneratedBy("github.com/podhmo/gos/openapigen/tools")
	b.NeedReference()

	b.Import("strings")

	Type := b.BuildTarget("Type",
		b.Field("Doc", seed.Symbol("string")).Tag(`json:"description,omitempty"`),
		b.Field("Title", seed.Symbol("string")).Tag(`json:"title,omitempty"`),
		b.Field("Format", seed.Symbol("string")).Tag(`json:"format,omitempty"`),
		b.Field("Example", seed.Symbol("string")).Tag(`json:"example,omitempty"`),
	).Setter("Doc", b.Arg("stmts", seed.Symbol("string")).Variadic().Transform(func(stmts string) string {
		return fmt.Sprintf(`strings.Join(%s, "\n")`, stmts)
	}))

	b.InterfaceMethods(
		"toSchemer // see: ./to_schema.go",
	)

	// ----------------------------------------
	// types
	// ----------------------------------------
	//
	// - todo: oneOf, anyOf, allOf, not
	// - see: https://github.com/getkin/kin-openapi/blob/master/openapi3/schema.go

	Bool := b.Type("Bool",
		b.Field("Default", seed.Symbol("bool")).Tag(`json:"default,omitempty"`),
	).NeedBuilder().Underlying("boolean")
	Int := b.Type("Int",
		b.Field("Enum", seed.Symbol("[]int64")).Tag(`json:"enum,omitempty"`),
		b.Field("Default", seed.Symbol("int64")).Tag(`json:"default,omitempty"`),
		b.Field("Maximum", seed.Symbol("int64")).Tag(`json:"maximum,omitempty"`),
		b.Field("Minimum", seed.Symbol("int64")).Tag(`json:"minimum,omitempty"`),
		b.Field("ExclusiveMin", seed.Symbol("bool")).Tag(`json:"exclusiveMin,omitempty"`),
		b.Field("ExclusiveMax", seed.Symbol("bool")).Tag(`json:"exclusiveMax,omitempty"`),
	).NeedBuilder().Underlying("integer")
	Float := b.Type("Float",
		b.Field("Default", seed.Symbol("string")).Tag(`json:"default,omitempty"`),
		b.Field("Maximum", seed.Symbol("float64")).Tag(`json:"maximum,omitempty"`),
		b.Field("Minimum", seed.Symbol("float64")).Tag(`json:"minimum,omitempty"`),
		b.Field("MultipleOf", seed.Symbol("float64")).Tag(`json:"multipleOf,omitempty"`),
		b.Field("ExclusiveMin", seed.Symbol("bool")).Tag(`json:"exclusiveMin,omitempty"`),
		b.Field("ExclusiveMax", seed.Symbol("bool")).Tag(`json:"exclusiveMax,omitempty"`),
	).NeedBuilder().Underlying("number")
	String := b.Type("String",
		b.Field("Enum", seed.Symbol("[]string")).Tag(`json:"enum,omitempty"`),
		b.Field("Default", seed.Symbol("string")).Tag(`json:"default,omitempty"`),
		b.Field("Pattern", seed.Symbol("string")).Tag(`json:"pattern,omitempty"`),
		b.Field("MaxLength", seed.Symbol("int64")).Tag(`json:"maxlength,omitempty"`),
		b.Field("MinLength", seed.Symbol("int64")).Tag(`json:"minlength,omitempty"`),
	).NeedBuilder().Underlying("string")
	Array := b.Type("Array", b.TypeVar("Items", seed.Symbol("TypeBuilder")),
		// b.Field("UniqueItems", seed.Symbol("bool")).Tag(`json:"uniqueItems,omitempty"`),
		b.Field("MaxItems", seed.Symbol("int64")).Tag(`json:"maxitems,omitempty"`),
		b.Field("MinItems", seed.Symbol("int64")).Tag(`json:"minitems,omitempty"`),
	).NeedBuilder().Underlying("array")
	Map := b.Type("Map", b.TypeVar("Items", seed.Symbol("TypeBuilder")),
		b.Field("Pattern", seed.Symbol("string")).Tag(`json:"pattern,omitempty"`),
	).NeedBuilder().Underlying("map")

	Field := b.Type("Field",
		b.Field("Name", seed.Symbol("string")).Tag(`json:"-"`),
		b.Field("Typ", seed.Symbol("TypeBuilder")).Tag(`json:"-"`),
		b.Field("Doc", seed.Symbol("string")).Tag(`json:"description,omitempty"`),
		b.Field("Nullable", seed.Symbol("bool")).Tag(`json:"nullable,omitempty"`), // trim omitempty?
		b.Field("Required", seed.Symbol("bool")).Tag(`json:"-"`).Default("true"),
		b.Field("ReadOnly", seed.Symbol("bool")).Tag(`json:"readonly,omitempty"`),
		b.Field("WriteOnly", seed.Symbol("bool")).Tag(`json:"writeonly,omitempty"`),
		b.Field("AllowEmptyValue", seed.Symbol("bool")).Tag(`json:"allowEmptyValue,omitempty"`),
		b.Field("Deprecated", seed.Symbol("bool")).Tag(`json:"deprecated,omitempty"`),
	).Constructor(
		b.Arg("Name", seed.Symbol("string")),
		b.Arg("Typ", seed.Symbol("TypeBuilder")),
	).Setter("Doc", b.Arg("stmts", seed.Symbol("string")).Variadic().Transform(func(stmts string) string {
		return fmt.Sprintf(`strings.Join(%s, "\n")`, stmts)
	})).NeedBuilder().Underlying("field") //?

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
		InterfaceMethods("TypeBuilder")

	// ----------------------------------------
	// action
	// ----------------------------------------
	Action := b.Type("Action",
		b.Field("Name", seed.Symbol("string")).Tag(`json:"-"`),
		b.Field("Input", "*Input").Tag(`json:"-"`),
		b.Field("Output", "*Output").Tag(`json:"-"`),
		b.Field("Method", seed.Symbol("string")).Tag(`json:"-"`),
		b.Field("Path", seed.Symbol("string")).Tag(`json:"-"`),
		b.Field("Tags", seed.Symbol("[]string")).Tag(`json:"tags"`),
		b.Field("DefaultStatus", seed.Symbol("int")).Tag(`json:"-"`).Default("200"),
	).Constructor(
		b.Arg("Name", seed.Symbol("string")),
		b.Arg("Input", "*Input"),
		b.Arg("Output", "*Output"),
	).NeedBuilder().Underlying("action")

	Param := b.Type("Param",
		b.Field("Name", seed.Symbol("string")).Tag(`json:"name"`),
		b.Field("In", seed.Symbol("string")).Tag(`json:"in"`).Default(`"query"`).Doc("openapi's in parameter {query, header, path, cookie} (default is query)"),
		b.Field("Typ", seed.Symbol("TypeBuilder")).Tag(`json:"-"`),
		b.Field("Doc", seed.Symbol("string")).Tag(`json:"description,omitempty"`),
		b.Field("Required", seed.Symbol("bool")).Tag(`json:"required"`).Default("true"),
		b.Field("Deprecated", seed.Symbol("bool")).Tag(`json:"deprecated,omitempty"`),
		b.Field("AllowEmptyValue", seed.Symbol("bool")).Tag(`json:"allowEmptyValue,omitempty"`),
	).Constructor(
		b.Arg("Name", seed.Symbol("string")),
		b.Arg("Typ", seed.Symbol("TypeBuilder")),
	).Setter("Doc", b.Arg("stmts", seed.Symbol("string")).Variadic().Transform(func(stmts string) string {
		return fmt.Sprintf(`strings.Join(%s, "\n")`, stmts)
	})).NeedBuilder().Underlying("param")

	Body := b.Type("Body",
		b.Field("Typ", seed.Symbol("TypeBuilder")).Tag(`json:"-"`),
	).Constructor(
		b.Arg("Typ", seed.Symbol("TypeBuilder")),
	).NeedBuilder()
	paramOrBody := b.Union("paramOrBody", Param, Body)

	Input := b.Type("Input",
		b.Field("Params", "[]*Param"),
		b.Field("Body", "*Body"),
	).Constructor(
		b.Arg("Params", paramOrBody.GetMetadata().Name).Variadic().BindFields("Params", "Body").Transform(func(s string) string {
			return fmt.Sprintf(`func()(v1 []*Param, v2 *Body){
				for _, a := range %s {
					switch a := a.(type) {
					case *Param:
						v1 = append(v1, a)
					case *Body:
						v2 = a
					}
				}
				return
			}()`, s)
		}),
	).NeedBuilder().Underlying("input")
	Output := b.Type("Output",
		b.Field("Typ", "TypeBuilder"),
		b.Field("DefaultError", seed.Symbol("TypeBuilder")),
	).Constructor(
		b.Arg("Typ", "TypeBuilder"),
	).NeedBuilder().Underlying("output")

	// openapi root info: https://swagger.io/specification/
	Contact := b.Type("Contact",
		b.Field("Name", seed.Symbol("string")).Tag(`json:"name"`),
		b.Field("URL", seed.Symbol("string")).Tag(`json:"url"`),
		b.Field("Email", seed.Symbol("string")).Tag(`json:"email"`),
	)
	License := b.Type("License",
		b.Field("Name", seed.Symbol("string")).Tag(`json:"name"`).Doc("required"),
		b.Field("Identifier", seed.Symbol("string")).Tag(`json:"identifier"`),
		b.Field("URL", seed.Symbol("string")).Tag(`json:"url"`),
	)
	Server := b.Type("Server",
		b.Field("URL", seed.Symbol("string")).Tag(`json:"url"`).Doc("required"),
		b.Field("Doc", seed.Symbol("string")).Tag(`json:"description"`),
		b.Field("Variables", seed.Symbol("map[string]any")).Tag(`json:"variables,omitempty"`).Doc("todo: typed"),
	).Setter("Doc", b.Arg("stmts", seed.Symbol("string")).Variadic().Transform(func(stmts string) string {
		return fmt.Sprintf(`strings.Join(%s, "\n")`, stmts)
	}))
	Info := b.Type("Info",
		b.Field("Title", seed.Symbol("string")).Tag(`json:"title"`).Doc("required"),
		b.Field("Summary", seed.Symbol("string")).Tag(`json:"summary,omitempty"`),
		b.Field("Doc", seed.Symbol("string")).Tag(`json:"description,omitempty"`),
		b.Field("TermsOfService", seed.Symbol("string")).Tag(`json:"termOfService,omitempty"`),
		b.Field("Contact", Contact.GetMetadata().Name.Pointer()).Tag(`json:"contact,omitempty"`),
		b.Field("License", License.GetMetadata().Name.Pointer()).Tag(`json:"license,omitempty"`),
		b.Field("Version", seed.Symbol("string")).Tag(`json:"version"`).Default(`"0.0.0"`).Doc("required"),
	).Setter("Doc", b.Arg("stmts", seed.Symbol("string")).Variadic().Transform(func(stmts string) string {
		return fmt.Sprintf(`strings.Join(%s, "\n")`, stmts)
	}))

	OpenAPI := b.Type("OpenAPI",
		b.Field("OpenAPI", seed.Symbol("string")).Tag(`json:"openapi"`).Default(`"3.0.3"`).Doc("required"),
		b.Field("Info", Info.GetMetadata().Name).Tag(`json:"info"`).Doc("required"),
		b.Field("Servers", Server.GetMetadata().Name.Slice()).Tag(`json:"servers"`),
	)

	fmt.Fprintln(os.Stderr, Type, Bool, Int, Float, String, Array, Map, Field, Object)
	fmt.Fprintln(os.Stderr, Action, Input, Output, Param)
	fmt.Fprintln(os.Stderr, Contact, License, Server, Info, OpenAPI)

	// for transform
	// b.Footer(``)

	// emit
	return cmd.Do(b)
}
