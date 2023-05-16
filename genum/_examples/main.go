package main

import (
	"os"

	"github.com/podhmo/gos/genum"
)

func main() {
	w := os.Stdout
	b := genum.NewEnumBuilder(genum.DefaultConfig())

	genum.DefineEnum("Ordering", b.String(
		b.StringValue("desc").Doc("降順"),
		b.StringValue("asc").Doc("昇順"),
	)).Default("desc").Doc("順序")

	genum.DefineEnum("Season", b.Int(
		b.IntValue(0, "Spring"),
		b.IntValue(1, "Summer"),
		b.IntValue(2, "Autumn"),
		b.IntValue(3, "Winter"),
	))

	if err := genum.WriteCode(w, b); err != nil {
		panic(err)
	}
}
