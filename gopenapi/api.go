//go:generate go run ./tools -write -builder -metadata -to-string -pkgname gopenapi
//go:generate go fmt .
package gopenapi

type Config struct {
}

func DefaultConfig() *Config {
	return &Config{}
}
