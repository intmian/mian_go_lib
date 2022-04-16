package menu

import "fmt"

// 最朴素的输入
type normalInput struct {
}

func (n normalInput) inputWithLen(strLen int) string {
	var input string
	fmt.Scanln(&input)
	if strLen != 0 {
		// 截断
		if strLen < len(input) {
			input = input[:strLen]
		}
	}
	return input
}

func (n normalInput) input() string {
	var input string
	_, err := fmt.Scanln(input)
	if err != nil {
		return ""
	}
	return input
}

func (n normalInput) outInput(s string) string {
	fmt.Println(s)
	return n.input()
}
