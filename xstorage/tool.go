package xstorage

import "strings"

func Join(src ...string) string {
	var result string
	for i, s := range src {
		if i == 0 {
			result = s
		} else {
			result += "." + s
		}
	}
	return result
}

func Split(src string) []string {
	return strings.Split(src, ".")
}
