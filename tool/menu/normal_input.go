package menu

import "fmt"

// 最朴素的输入
type normalInput struct {
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
