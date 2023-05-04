package genum_test

import (
	"os"

	"github.com/podhmo/gos/genum"
)

func ExampleWriteCode() {
	w := os.Stdout
	padding := "@@"
	comment := "#"

	b := genum.NewBuilder[string]()
	b.Config.Padding = padding
	b.Config.Comment = comment

	genum.Define("RGBColor", b.Enum(
		b.Value("R").Name("Red").Doc("red color").Default(true),
		b.Value("G").Name("Green").Doc("green color"),
		b.Value("B").Name("Blue").Doc("blue color"),
	).Doc("rgb"))

	if err := genum.WriteCode(w, b); err != nil {
		panic(err)
	}

	// Output:
	// # RGBColor
	// type RGBColor string
	//
	// const (
	//@@RGBColorRed RGBColor = "R"  # default
	//@@RGBColorGreen RGBColor = "G"
	//@@RGBColorBlue RGBColor = "B"
	// )
}
