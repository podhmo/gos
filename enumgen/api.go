//go:generate go run ./tools -write -builder -metadata -stringer -pkgname enumgen
//go:generate go fmt .
package enumgen

type Config struct {
	Padding string // default "\t"
	Comment string // default "//"
}

func DefaultConfig() *Config {
	return &Config{Padding: "\t", Comment: "//"}
}
