package enum

import "github.com/podhmo/gos/enumgen"

var b = enumgen.NewBuilder(enumgen.DefaultConfig())

var (
	Ordering = b.String(
		b.StringValue("desc").Doc("降順"),
		b.StringValue("asc").Doc("昇順"),
	).Default("desc").Doc("順序")
)
