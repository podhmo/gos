package enumgen_test

import (
	"os"

	"github.com/podhmo/gos/enumgen"
)

func ExampleWriteCode() {
	config := &enumgen.Config{Comment: "#", Padding: "@@"}
	b := enumgen.NewEnumBuilder(config)

	enumgen.DefineEnum("RGBColor", b.String(
		b.StringValue("R").Name("Red").Doc("red color"),
		b.StringValue("G").Name("Green").Doc("green color"),
		b.StringValue("B").Name("Blue").Doc("blue color"),
	)).Doc("rgb").Default("R")

	w := os.Stdout
	if err := enumgen.WriteCode(w, b); err != nil {
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
