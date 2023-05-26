//go:generate go run ./tools -write -builder -metadata -stringer -pkgname openapigen
//go:generate go fmt .
package openapigen

type Config struct {
	DisableRefLinks bool // if true, does not use $ref links

	defs []toSchemer
	seen map[toSchemer]bool // TODO: name conflict?
}

func DefaultConfig() *Config {
	return &Config{
		DisableRefLinks: false,
		seen:            map[toSchemer]bool{},
	}
}
