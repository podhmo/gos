//go:generate go run ./tools -write -builder -metadata -stringer -pkgname openapigen
//go:generate go fmt .
package openapigen

type Config struct {
}

func DefaultConfig() *Config {
	return &Config{}
}
