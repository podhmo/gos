package main

import (
	_ "embed"
	"io"
	"os"
	"strings"
)

//go:embed testdata/gen.go.tmpl
var sourceFile string

func main() {
	var w io.Writer = os.Stdout
	if _, err := io.Copy(w, strings.NewReader(sourceFile)); err != nil {
		panic(err)
	}
}
