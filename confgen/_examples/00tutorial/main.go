package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/podhmo/gos/confgen"
)

func main() {
	c := confgen.DefaultConfig()
	b := confgen.NewBuilder(c)

	// https://json-schema.org/learn/getting-started-step-by-step.html

	Product := confgen.Define("Product", b.Object(
		b.Field("productId", b.Int()).Doc(`The unique identifier for a product`),
		b.Field("productName", b.String()).Doc(`Name of the product`),
		b.Field("price", b.Float().ExclusiveMin(0)).Doc(`The price of the product`),
		b.Field("tags", b.Array(b.String()).MinItems(1).UniqueItems(true)).Required(false).Doc(`tags for the product`),
		b.Field("dimensions", b.Object(
			b.Field("length", b.Float()),
			b.Field("width", b.Float()),
			b.Field("height", b.Float()),
		)).Required(false),
	)).Doc(`A product from Acme's catalog`)

	doc, err := confgen.ToJSONSchema(b, Product)
	if err != nil {
		panic(err)
	}

	w := os.Stdout
	enc := json.NewEncoder(w)
	enc.SetIndent("", "    ")
	if err := enc.Encode(doc); err != nil {
		panic(err)
	}

	{
		w := os.Stderr
		fmt.Fprintln(os.Stderr, "----------------------------------------")
		fmt.Fprintln(w, "package M")
		if err := c.ToGoCode(w, Product); err != nil {
			panic(err)
		}
	}
}
