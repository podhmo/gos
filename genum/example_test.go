package genum_test

import (
	"os"

	"github.com/podhmo/gos/genum"
)

func ExampleWriteCode() {
	config := &genum.Config{Comment: "#", Padding: "@@"}
	b := genum.NewEnumBuilder(config)

	genum.DefineEnum("RGBColor", b.String(
		genum.StringValue{Name: "Red", Value: "R", Doc: "red color"},
		genum.StringValue{Name: "Green", Value: "G", Doc: "green color"},
		genum.StringValue{Name: "Blue", Value: "B", Doc: "blue color"},
	)).Doc("rgb").Default("R")

	w := os.Stdout
	if err := genum.WriteCode(w, b); err != nil {
		panic(err)
	}

	// Output:
	// # RGBColor
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
