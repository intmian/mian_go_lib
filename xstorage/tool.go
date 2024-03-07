package xstorage

import "strings"

func Join(src ...string) string {
	var buffer strings.Builder
	for i, s := range src {
		buffer.WriteString(s)
		if i != len(src)-1 {
			buffer.WriteString(".")
		}
	}
	return buffer.String()
}

func Split(src string) []string {
	return strings.Split(src, ".")
}
