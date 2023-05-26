package main

import (
	"os"

	"github.com/podhmo/gos/enumgen"
)

func main() {
	w := os.Stdout
	b := enumgen.NewEnumBuilder(enumgen.DefaultConfig())

	enumgen.DefineEnum("Ordering", b.String(
		b.StringValue("desc").Doc("降順"),
		b.StringValue("asc").Doc("昇順"),
	)).Default("desc").Doc("順序")

	enumgen.DefineEnum("Season", b.Int(
		b.IntValue(0, "Spring"),
		b.IntValue(1, "Summer"),
		b.IntValue(2, "Autumn"),
		b.IntValue(3, "Winter"),
	))

	if err := enumgen.WriteCode(w, b); err != nil {
		panic(err)
	}
}