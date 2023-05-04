package main

import (
	"os"

	"github.com/podhmo/gos/genum"
)

func main() {
	{
		b := genum.NewBuilder[int]()

		// simple
		genum.Define("OneTwo", b.Enum(
			b.Value(1).Name("One"),
			b.Value(2).Name("Two"),
		).Default(1))

		w := os.Stdout
		if err := genum.WriteCode(w, b); err != nil {
			panic(err)
		}
	}

	{
		b := genum.NewBuilder[string]()

		// complex
		genum.Define("RGBColor", b.Enum(
			b.Value("R").Name("Red").Doc("red color"),
			b.Value("G").Name("Green").Doc("green color"),
			b.Value("B").Name("Blue").Doc("blue color"),
		)).Default("R").Doc("rgb")
	}
}
