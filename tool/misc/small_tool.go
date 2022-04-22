package misc

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

//Str2list 将json风格的数组字符串s转换为字符串
func Str2list(s string) []string {
	return strings.Split(s[1:len(s)-1], ",")
}

//MinMaxInt 求一个int切片的最大最小值
func MinMaxInt(array []int) (min int, minIndex int, max int, maxIndex int) {
	max = array[0]
	min = array[0]
	for i, value := range array {
		if max < value {
			max = value
			maxIndex = i
		}
		if min > value {
			min = value
			minIndex = i
		}
	}
	return
}

//MinMaxUInt 求一个int切片的最大最小值
func MinMaxUInt(array []uint32) (min uint32, minIndex int, max uint32, maxIndex int) {
	max = array[0]
	min = array[0]
	for i, value := range array {
		if max < value {
			max = value
			maxIndex = i
		}
		if min > value {
			min = value
			minIndex = i
		}
	}
	return
}

type CanCompare interface {
	Less(CanCompare) bool
}

//MinMax 求一个CanCompare切片的最大最小值
func MinMax(array []CanCompare) (min CanCompare, minIndex int, max CanCompare, maxIndex int) {
	max = array[0]
	min = array[0]
	for i, value := range array {
		if max.Less(value) {
			max = value
			maxIndex = i
		}
		if min.Less(value) {
			min = value
			minIndex = i
		}
	}
	return
}

//GetTimeStr 返回当前的时间，格式为2006-01-02 15:04:05
func GetTimeStr() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func InputStr(len int) string {
	re := ""
	for i := 0; i < len; i++ {
		re += "_"
	}
	for i := 0; i < len; i++ {
		re += "\b"
	}
	return re
}

func Input(hint string, len int, a ...interface{}) error {
	print(hint + InputStr(len))
	_, err := fmt.Scan(a...)
	return err
}

func Stop() {
	fmt.Printf("输入任意键继续...")
	ClearIOBuffer()
	b := make([]byte, 1)
	// 不知道为什么清空缓冲区后，还是有残留一个ascii为10的字符。。。但是goland里面调试时好的。。
	_, err := os.Stdin.Read(b)
	_, err = os.Stdin.Read(b)
	if err != nil {
		return
	}
}

func ClearIOBuffer() {
	myReader := bufio.NewReader(nil)
	myReader.Reset(os.Stdin)
}

//IsLegalOutURL 判断是否是合法的外链
func IsLegalOutURL(url string) bool{
	if len(url) < 6 {
		return false
	}
	if strings.Index(url,"://") == -1 {
		return false
	}

	return true
}

//Clear 清空屏幕
func clear() {
	cmd := exec.Command("cmd.exe", "/c", "cls")
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		print("clear fail")
	}
}