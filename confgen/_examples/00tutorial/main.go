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
		b.Field("productName", b.String()).Doc(`Name of the product`),
		b.Field("price", b.Float()).Doc(`The price of the product`, "(TODO: exclsiveMinimum: 0)"),
		b.Field("tags", b.Array(b.String()).MinItems(1).UniqueItems(true)).Required(false).Doc(`tags for the product`),
	)).Doc(`A product from Acme's catalog`)

	doc, err := confgen.ToJSONSchema(b, Product)
	if err != nil {
		panic(err)
	}
	doc.Set("$id", "https://example.com/product.schema.json")

	w := os.Stdout
	enc := json.NewEncoder(w)
	enc.SetIndent("", "    ")
	if err := enc.Encode(doc); err != nil {
		panic(err)
	}
}
