package main

import (
	"io"
	"os"

	"github.com/podhmo/gos/graphqlgen"
)

func main() {
	b := graphqlgen.NewBuilder()

	// https://github.com/99designs/gqlgen/blob/4a78eb0c9be84793df821616b4ebfa4bfb42a49c/_examples/todo/schema.graphql
	graphqlgen.Define("Todo", b.Object(
		b.Field("id", b.String()), // TODO: ID
		b.Field("text", b.String()),
		b.Field("done", b.Bool()), // TODO: @hasRole(role: Owner) # only the owner can set if a todo is done
	))

	var w io.Writer = os.Stdout
	if err := graphqlgen.ToGraphql(w, b); err != nil {
		panic(err)
	}

}
