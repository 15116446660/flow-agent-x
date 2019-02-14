package dao

import (
	"strings"
	"unicode"
)

// change camel string to string with '_'
func flatCamelString(v string) string {
	var builder strings.Builder
	builder.Grow(len(v) + 5)

	for i, c := range v {
		r := rune(c)

		if unicode.IsUpper(r) {
			r = unicode.ToLower(r)

			if i > 0 {
				builder.WriteByte('_')
			}
		}

		builder.WriteByte(byte(r))
	}

	return builder.String()
}

func capitalFirstChar(v string) string {
	bytes := []byte(v)
	bytes[0] = byte(unicode.ToUpper(rune(bytes[0])))
	return string(bytes)
}
