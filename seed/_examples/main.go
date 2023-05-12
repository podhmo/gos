package main

import (
	"fmt"
	"os"
	"text/template"

	"github.com/podhmo/gos/seed"
)

func main() {
	b := seed.NewBuilder()

	Type := b.BuildTarget("Type")
	goStringType := seed.Symbol("string") // unexported string is go-primitive

	Int := b.Type("Int").NeedBuilder()
	String := b.Type("String").NeedBuilder()

	Array := b.Type("Array").Var("typ", Type).NeedBuilder()
	Map := b.Type("Map").Var("valtype", Type).NeedBuilder()

	Field := b.Type("Field")
	FieldType := Field.Metadata.Name
	Object := b.Type("Object").Field("Fields", FieldType).NeedBuilder()

	Param := b.Type("Param")
	ParamType := Param.Metadata.Name
	ActionInput := b.Type("ActionInput").Field("Params", ParamType)
	ActionOutput := b.Type("ActionOutput").Field("Return", Type)
	Action := b.Type("Action").Field("Name", goStringType).Field("Input", ActionInput.Metadata.Name).Field("Output", ActionOutput.Metadata.Name).NeedBuilder()

	// emit
	{
		fmt.Fprintln(os.Stderr, Type, Int, String, Array, Map, Field, Object, Param, ActionInput, ActionOutput, Action)
		fmt.Fprintln(os.Stderr, b.Metadata.Types)
		t := template.Must(template.New("").Parse(tmpl))
		if err := t.Execute(os.Stdout, b.Metadata); err != nil {
			panic(err)
		}
	}
}

const tmpl = `
package M

type Builder struct {}

{{range $_, $t := .Types}}{{with $name := ($t.Metadata.Name)}}
{{if $t.Metadata.NeedBuilder }}
func (b *Builder) {{$name}}() *{{$name}} {
	t := &{{$name}}{
		{{$name}}Builder: &{{$name}}Builder[*{{$name}}]{Metadata: &{{$name}}Metadata{}},
	}
	t.ret = t
	return t
}
type {{$name}} struct {
	*{{$name}}Builder[*{{$name}}]
}

type {{$name}}Builder[R any] struct {
	Metadata *{{$name}}Metadata
	ret R
}
{{else}}
type {{$name}} struct {
}
{{end}}
{{end}}{{end}}
`