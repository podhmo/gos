package namelib

import "strings"

// foo -> Foo
func ToTitle(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
