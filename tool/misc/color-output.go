package misc

import "fmt"

const (
	textBlack = iota + 30
	textRed
	textGreen
	textYellow
	textBlue
	textPurple
	textCyan
	textWhite
)

func Black(str string) string {
	if str == "" {
		return ""
	}
	return textColor(textBlack, str)
}

func Red(str string) string {
	if str == "" {
		return ""
	}
	return textColor(textRed, str)
}
func Yellow(str string) string {
	if str == "" {
		return ""
	}
	return textColor(textYellow, str)
}
func Green(str string) string {
	if str == "" {
		return ""
	}
	return textColor(textGreen, str)
}
func Cyan(str string) string {
	if str == "" {
		return ""
	}
	return textColor(textCyan, str)
}
func Blue(str string) string {
	if str == "" {
		return ""
	}
	return textColor(textBlue, str)
}
func Purple(str string) string {
	if str == "" {
		return ""
	}
	return textColor(textPurple, str)
}
func White(str string) string {
	if str == "" {
		return ""
	}
	return textColor(textWhite, str)
}

func textColor(color int, str string) string {
	return fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", color, str)
}
