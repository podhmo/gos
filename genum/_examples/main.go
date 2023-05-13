package main

import (
	"os"

	"github.com/podhmo/gos/genum"
)

func main() {
	w := os.Stdout
	b := genum.NewEnumBuilder()

	genum.DefineEnum("Ordering", b.String(
		genum.StringValue{Value: "desc", Doc: "降順"},
		genum.StringValue{Value: "asc", Doc: "昇順"},
	)).Default("desc")

	genum.DefineEnum("Season", b.Int(
		genum.IntValue{Name: "Spring", Value: 0},
		genum.IntValue{Name: "Summer", Value: 1},
		genum.IntValue{Name: "Autumn", Value: 2},
		genum.IntValue{Name: "Wrinter", Value: 3},
	))

	if err := genum.WriteCode(w, b); err != nil {
		panic(err)
	}
}
