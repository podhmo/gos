//go:generate go run ./tools -write -builder -metadata -stringer -pkgname enumgen
//go:generate go fmt .
package enumgen

type Config struct {
}

func DefaultConfig() *Config {
	return &Config{}
}
