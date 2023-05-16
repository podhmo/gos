package genum_test

import (
	"os"

	"github.com/podhmo/gos/genum"
)

func ExampleWriteCode() {
	config := &genum.Config{Comment: "#", Padding: "@@"}
	b := genum.NewEnumBuilder(config)

	genum.DefineEnum("RGBColor", b.String(
		b.StringValue("R").Name("Red").Doc("red color"),
		b.StringValue("G").Name("Green").Doc("green color"),
		b.StringValue("B").Name("Blue").Doc("blue color"),
	)).Doc("rgb").Default("R")

	w := os.Stdout
	if err := genum.WriteCode(w, b); err != nil {
		panic(err)
	}

	// Output:
	// # RGBColor : rgb
	// type RGBColor string
	//
	// const (
	// @@# "red color"
	// @@RGBColorRed RGBColor = "R"  # default
	// @@# "green color"
	// @@RGBColorGreen RGBColor = "G"
	// @@# "blue color"
	// @@RGBColorGreen RGBColor = "B"
	// )
}
