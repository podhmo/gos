//go:generate go run ./tools -write -builder -metadata -stringer -pkgname openapigen
//go:generate go fmt .
package openapigen

type Config struct {
	DisableRefLinks bool // if true, does not use $ref links

	// for to schema
	defs []TypeBuilder
	seen map[int]*TypeRef
}

func DefaultConfig() *Config {
	return &Config{
		DisableRefLinks: false,
		seen:            map[int]*TypeRef{},
	}
}
