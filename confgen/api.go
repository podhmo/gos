//go:generate go run ./tools -write -builder -metadata -stringer -pkgname confgen
//go:generate go fmt .
package confgen

import "io"

type Config struct {
}

func DefaultConfig() *Config {
	c := &Config{}
	return c
}

func EmitSchema(w io.Writer, t Type) error {
	return nil
}
