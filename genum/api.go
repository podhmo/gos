//go:generate go run ./tools -write -builder -metadata -stringer -pkgname genum
//go:generate go fmt .
package genum

type Config struct {
	Padding string // default "\t"
	Comment string // default "//"
}

func DefaultConfig() *Config {
	return &Config{Padding: "\t", Comment: "//"}
}
