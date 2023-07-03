package main

import (
	"os"

	"github.com/podhmo/gos/confgen"
)

func main() {
	b := confgen.NewBuilder(confgen.DefaultConfig())

	w := os.Stdout

	// https://json-schema.org/learn/getting-started-step-by-step.html

	Product := confgen.Define("Product", b.Object(
		b.Field("productId", b.Int()).Doc(`The unique identifier for a product`),
	))

	if err := confgen.EmitSchema(w, Product); err != nil {
		panic(err)
	}
}
