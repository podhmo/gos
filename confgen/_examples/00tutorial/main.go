package main

import (
	"encoding/json"
	"os"

	"github.com/podhmo/gos/confgen"
)

func main() {
	b := confgen.NewBuilder(confgen.DefaultConfig())

	// https://json-schema.org/learn/getting-started-step-by-step.html

	Product := confgen.Define("Product", b.Object(
		b.Field("productId", b.Int()).Doc(`The unique identifier for a product`),
	))

	doc, err := confgen.ToJSONSchema(b, Product)
	if err != nil {
		panic(err)
	}

	w := os.Stdout
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(doc); err != nil {
		panic(err)
	}
}
