package misc

import (
	"regexp"
	"strconv"
)

func ReplaceUnicodeEscapes(input string) string {
	re := regexp.MustCompile(`\\u([0-9A-Fa-f]{4})`)

	// 使用 ReplaceAllStringFunc 函数将所有 Unicode 转义字符替换为对应的实际字符
	result := re.ReplaceAllStringFunc(input, func(match string) string {
		// 将 Unicode 转义字符转换为十进制整数
		unicodeInt, _ := strconv.ParseInt(match[2:], 16, 32)

		// 将整数转换为对应的 Unicode 字符
		return string(rune(unicodeInt))
	})

	return result
}
